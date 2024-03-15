package internal

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

const ONEPASSWORD_EXECUTABLE = "op"

type OnePasswordStore struct{}

type Item struct {
	Id     string  `json:"id"`
	Fields []Field `json:"fields"`
}
type Field struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (s *OnePasswordStore) FindCertificatesThatAreOutdated() []string {
	items, _ := getListOfItems()

	for _, v := range items {
		getItemDetails(v.Id)
	}

	return nil
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
