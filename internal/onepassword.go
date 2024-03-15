package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const ONEPASSWORD_EXECUTABLE = "op"

type OnePasswordStore struct{}

func (s *OnePasswordStore) FindCertificatesThatAreOutdated() []string {
	getListOfItems()
	return nil
}

func NewOnePasswordStore() *OnePasswordStore {
	return &OnePasswordStore{}
}

func getListOfItems() {
	listItems := createCommand(findListOfItems(), withTags("automation", "certificate"), withFormat("json"))
	cmd := exec.Command(ONEPASSWORD_EXECUTABLE, listItems...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("%s\n", stderr.String())
	}

	fmt.Printf("%s\n", out.String())
}

func createCommand(commands ...[]string) []string {
	allCommands := []string{}

	for _, v := range commands {
		allCommands = append(allCommands, v...)
	}

	return allCommands
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

func withFormat(format string) []string {
	return []string{
		"--format",
		format,
	}
}
