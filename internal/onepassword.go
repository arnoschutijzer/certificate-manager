package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"
)

const ONEPASSWORD_EXECUTABLE = "op"

type OnePasswordStore struct{}

const STRING_TYPE = "STRING"

type Item struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Vault  Vault   `json:"vault"`
	Fields []Field `json:"fields"`
}

type Vault struct {
	Name string `json:"name"`
}

var errNoContentFieldFound = errors.New("no content field found")

func (i *Item) findContentField() (*Field, error) {
	for _, v := range i.Fields {
		if v.Type == STRING_TYPE {
			return &v, nil
		}
	}

	return nil, errNoContentFieldFound
}

type Field struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (s *OnePasswordStore) FindCertificatesThatAreOutdated() ([]Item, error) {
	return s.FindCertificatesOlderThanDate(time.Now())
}

func (s *OnePasswordStore) FindCertificatesOlderThanDate(date time.Time) ([]Item, error) {
	items, err := getListOfItemsWithDetails()

	if err != nil {
		return nil, err
	}

	itemsWithCertificates := []Item{}
	for _, v := range items {
		if doesItemContainAtLeastOneCertificate(v) {
			itemsWithCertificates = append(itemsWithCertificates, v)
		}
	}

	outdatedCertificates := []Item{}
	for _, v := range itemsWithCertificates {
		field, _ := v.findContentField()

		if !IsValidCertificate([]byte(field.Value), date) {
			outdatedCertificates = append(outdatedCertificates, v)
		}
	}

	return outdatedCertificates, nil
}

func doesItemContainAtLeastOneCertificate(item Item) bool {
	field, err := item.findContentField()

	if errors.Is(err, errNoContentFieldFound) {
		return false
	}

	if err != nil {
		return false
	}

	certificates := GetCertificatesFromString(field.Value)

	return len(certificates) > 0
}

func NewOnePasswordStore() *OnePasswordStore {
	return &OnePasswordStore{}
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

func getListOfItemsWithDetails() ([]Item, error) {
	items, _ := getListOfItems()

	itemsWithDetails := make([]Item, len(items))
	for _, v := range items {
		item, err := getItemDetails(v.Id)

		itemsWithDetails = append(itemsWithDetails, item)

		if err != nil {
			return []Item{}, err
		}
	}

	return itemsWithDetails, nil
}

func getListOfItems() ([]Item, error) {
	cmd := createCommand(findListOfItems(), withTags("automation", "certificate"), withJsonFormat())
	items, err := execute[[]Item](cmd)

	if err != nil {
		panic(err)
	}

	return *items, nil
}

func findListOfItems() []string {
	return []string{"item", "list"}
}

func withTags(tags ...string) []string {
	commandTags := []string{
		"--tags",
	}
	commandTags = append(commandTags, strings.Join(tags, ","))
	return commandTags
}

func getItemDetails(id string) (Item, error) {
	cmd := createCommand(findItemDetails(id), withJsonFormat())
	item, err := execute[Item](cmd)
	if err != nil {
		return Item{}, err
	}

	return *item, nil
}

func findItemDetails(id string) []string {
	return []string{"item", "get", id}
}

func createCommand(commands ...[]string) *exec.Cmd {
	allCommands := []string{}

	for _, v := range commands {
		allCommands = append(allCommands, v...)
	}

	return exec.Command(ONEPASSWORD_EXECUTABLE, allCommands...)
}

func withJsonFormat() []string {
	return []string{
		"--format",
		"json",
	}
}
