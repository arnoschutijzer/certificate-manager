package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/algleymi/certificate-manager/internal"
	"github.com/algleymi/certificate-manager/internal/caches"
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

func (s *OnePasswordStore) FindCertificatesThatAreOutdated() ([]internal.Certificate, error) {
	return s.FindCertificatesOlderThanDate(time.Now())
}

func (s *OnePasswordStore) FindCertificatesOlderThanDate(date time.Time) ([]internal.Certificate, error) {
	items, err := getListOfItems()
	if err != nil {
		return nil, err
	}

	vaultItems, err := internal.Map(items, s.retrieveMaybeCachedVaultItem)

	certificates, _ := internal.FlatMap(vaultItems, func(vaultItem caches.CachedItem) ([]internal.Certificate, error) {
		return internal.Map(vaultItem.Certificates, func(certificate caches.CachedCertificate) (internal.Certificate, error) {
			return caches.ToDomainCertificate(certificate), nil
		})
	})

	return internal.Filter(certificates, func(certificate internal.Certificate) bool {
		return !certificate.IsValid(date)
	}), nil
}

func (s *OnePasswordStore) retrieveMaybeCachedVaultItem(item Item) (caches.CachedItem, error) {
	cached, err := s.cache.RetrieveVaultItem(item.Id)

	if err != nil {
		vaultItem, err := s.retrieveVaultItemAndCache(item.Id)

		if err != nil {
			return caches.CachedItem{}, err
		}

		return vaultItem, nil
	}

	if !item.UpdatedAt.After(cached.UpdatedAt) {
		return cached, nil
	}

	vaultItem, err := s.retrieveVaultItemAndCache(item.Id)

	if err != nil {
		return caches.CachedItem{}, err
	}

	return vaultItem, nil
}

func (s *OnePasswordStore) retrieveVaultItemAndCache(id string) (caches.CachedItem, error) {
	itemWithFields, err := getItemDetails(id)
	if err != nil {
		return caches.CachedItem{}, err
	}

	vaultItem := mapOnePasswordToInternal(itemWithFields)
	s.cache.SaveVaultItem(vaultItem)

	return vaultItem, nil
}

func mapOnePasswordToInternal(item ItemWithFields) caches.CachedItem {
	field, _ := item.findContentField()
	certificates := internal.GetCertificatesFromString(field.Value, item.Title)

	cacheCertificates := []caches.CachedCertificate{}
	for _, certificate := range certificates {
		cacheCertificates = append(cacheCertificates, caches.ToDbCertificate(item.Id, certificate))
	}

	return caches.CachedItem{
		VaultId:      item.Id,
		Title:        item.Title,
		UpdatedAt:    item.UpdatedAt,
		Certificates: cacheCertificates,
	}
}

func doesItemContainAtLeastOneCertificate(item ItemWithFields) bool {
	field, err := item.findContentField()

	if errors.Is(err, errNoContentFieldFound) {
		return false
	}

	if err != nil {
		return false
	}

	return internal.DoesSecretContainAnyCertificate(field.Value)
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
		withTags("certificate"),
		withCategories("SecureNote"),
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
