package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	ONEPASSWORD_EXECUTABLE = "op"
	NOTES_PURPOSE          = "NOTES"
	STRING_TYPE            = "STRING"
)

type Item struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
}

var errNoContentFieldFound = errors.New("no content field found")

func (i *Item) findContentField() (*Field, error) {
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

type OnePasswordStore struct{}

func (s *OnePasswordStore) FindCertificatesThatAreOutdated() ([]Certificate, error) {
	return s.FindCertificatesOlderThanDate(time.Now())
}

func (s *OnePasswordStore) FindCertificatesOlderThanDate(date time.Time) ([]Certificate, error) {
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

	outdatedCertificates := []Certificate{}
	for _, v := range itemsWithCertificates {
		field, _ := v.findContentField()
		for _, certificate := range GetCertificatesFromString(field.Value, v.Title) {
			if !certificate.IsValid(date) {
				outdatedCertificates = append(outdatedCertificates, certificate)
			}
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

	return DoesSecretContainAnyCertificate(field.Value)
}

func NewOnePasswordStore() Store {
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

func createCommand(commands ...[]string) *exec.Cmd {
	allCommands := []string{}

	for _, v := range commands {
		allCommands = append(allCommands, v...)
	}

	return exec.Command(ONEPASSWORD_EXECUTABLE, allCommands...)
}

func getListOfItemsWithDetails() ([]Item, error) {
	items, _ := getListOfItems()

	totalItems := len(items)

	itemsWithDetails := make([]Item, len(items))
	for i, v := range items {
		fmt.Printf("Retrieving item details %d/%d\n", i+1, totalItems)

		before := time.Now()
		item, err := getItemDetails(v.Id)
		fmt.Printf("Retrieved item details in %f\n", time.Since(before).Seconds())

		itemsWithDetails = append(itemsWithDetails, item)

		if err != nil {
			return []Item{}, err
		}
	}

	return itemsWithDetails, nil
}

func getListOfItems() ([]Item, error) {
	cmd := createCommand(listItemsCommand(), withCategories("SecureNote"), withJsonFormat())
	items, err := execute[[]Item](cmd)

	if err != nil {
		panic(err)
	}

	return *items, nil
}

func listItemsCommand() []string {
	return []string{"item", "list"}
}

func getItemDetails(id string) (Item, error) {
	cmd := createCommand(listItemDetailsCommand(id), withJsonFormat())
	item, err := execute[Item](cmd)
	if err != nil {
		return Item{}, err
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
