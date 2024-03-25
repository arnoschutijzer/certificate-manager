package vaults

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/algleymi/certificate-manager/internal"
	"github.com/algleymi/certificate-manager/internal/caches"
	"github.com/algleymi/certificate-manager/internal/domain"
)

const (
	ONEPASSWORD_EXECUTABLE = "op"
	NOTES_PURPOSE          = "NOTES"
	STRING_TYPE            = "STRING"
)

type Item struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemWithFields struct {
	Item
	Fields []Field `json:"fields"`
}

var errNoContentFieldFound = errors.New("no content field found")

func (i *ItemWithFields) findContentField() (*Field, error) {
	for _, v := range i.Fields {
		if v.Type == STRING_TYPE && v.Purpose == NOTES_PURPOSE {
			return &v, nil
		}
	}

	return nil, errNoContentFieldFound
}

type Field struct {
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
	Value   string `json:"value"`
}

type OnePasswordStore struct {
	cache caches.Cache
}

func NewOnePasswordStore(cache caches.Cache) Vault {
	return &OnePasswordStore{cache: cache}
}

func (s *OnePasswordStore) FindCertificatesThatAreOutdated() ([]domain.Certificate, error) {
	return s.FindCertificatesOlderThanDate(time.Now())
}

func (s *OnePasswordStore) FindCertificatesOlderThanDate(date time.Time) ([]domain.Certificate, error) {
	items, err := getListOfItems()
	if err != nil {
		return nil, err
	}

	secrets, err := internal.Map(items, s.retrieveMaybeCachedSecret)

	certificates, _ := internal.FlatMap(secrets, func(secret domain.Secret) ([]domain.Certificate, error) {
		return internal.Map(secret.Certificates, func(certificate domain.Certificate) (domain.Certificate, error) {
			return domain.Certificate{
				Fingerprint: certificate.Fingerprint,
				Subject:     certificate.Subject,
				CustomName:  certificate.CustomName,
				NotBefore:   certificate.NotBefore,
				NotAfter:    certificate.NotAfter,
			}, nil
		})
	})

	invalidCertificates := internal.Filter(certificates, func(certificate domain.Certificate) bool {
		return !certificate.IsValid(date)
	})

	return deduplicate(invalidCertificates), nil
}

func deduplicate(certificates []domain.Certificate) []domain.Certificate {
	m := make(map[string]domain.Certificate)

	for _, v := range certificates {
		m[v.Fingerprint] = v
	}

	uniq := []domain.Certificate{}
	for _, v := range m {
		uniq = append(uniq, v)
	}

	return uniq
}

func (s *OnePasswordStore) retrieveMaybeCachedSecret(item Item) (domain.Secret, error) {
	cached, err := s.cache.RetrieveSecret(item.Id)

	if err != nil {
		secret, err := s.retrieveSecretAndCache(item.Id)

		if err != nil {
			return domain.Secret{}, err
		}

		return secret, nil
	}

	if !item.UpdatedAt.After(cached.UpdatedAt) {
		return cached, nil
	}

	secret, err := s.retrieveSecretAndCache(item.Id)

	if err != nil {
		return domain.Secret{}, err
	}

	return secret, nil
}

func (s *OnePasswordStore) retrieveSecretAndCache(id string) (domain.Secret, error) {
	itemWithFields, err := getItemDetails(id)
	if err != nil {
		return domain.Secret{}, err
	}

	secret := mapOnePasswordToInternal(itemWithFields)
	s.cache.SaveSecret(secret)

	return secret, nil
}

func mapOnePasswordToInternal(item ItemWithFields) domain.Secret {
	field, _ := item.findContentField()
	certificates := internal.GetCertificatesFromString(field.Value, item.Title)

	cacheCertificates := []domain.Certificate{}
	for _, certificate := range certificates {
		cacheCertificates = append(cacheCertificates, certificate)
	}

	return domain.Secret{
		Id:           item.Id,
		Title:        item.Title,
		UpdatedAt:    item.UpdatedAt,
		Certificates: cacheCertificates,
	}
}

func execute[T any](cmd *exec.Cmd) (*T, error) {
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var response T
	json.Unmarshal(out.Bytes(), &response)

	return &response, nil
}

func createCommand(commands ...[]string) *exec.Cmd {
	allCommands := []string{}

	for _, v := range commands {
		allCommands = append(allCommands, v...)
	}

	return exec.Command(ONEPASSWORD_EXECUTABLE, allCommands...)
}

func getListOfItems() ([]Item, error) {
	cmd := createCommand(
		listItemsCommand(),
		// withTags("certificate"),
		// withCategories("SecureNote"),
		withJsonFormat(),
	)
	items, err := execute[[]Item](cmd)

	if err != nil {
		panic(err)
	}

	return *items, nil
}

func listItemsCommand() []string {
	return []string{"item", "list"}
}

func getItemDetails(id string) (ItemWithFields, error) {
	cmd := createCommand(listItemDetailsCommand(id), withJsonFormat())
	item, err := execute[ItemWithFields](cmd)
	if err != nil {
		return ItemWithFields{}, err
	}

	return *item, nil
}

func listItemDetailsCommand(id string) []string {
	return []string{"item", "get", id}
}

func withTags(tags ...string) []string {
	commandTags := []string{
		"--tags",
	}
	commandTags = append(commandTags, strings.Join(tags, ","))
	return commandTags
}

func withCategories(categories ...string) []string {
	commandCategories := []string{
		"--categories",
	}

	commandCategories = append(commandCategories, strings.Join(categories, ","))
	return commandCategories
}

func withJsonFormat() []string {
	return []string{
		"--format",
		"json",
	}
}
