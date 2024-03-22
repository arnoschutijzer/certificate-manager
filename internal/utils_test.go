package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiltersStrings(t *testing.T) {
	strings := []string{
		"hello",
		"world",
	}
	predicate := func(value string) bool {
		return value == "hello"
	}
	filtered := Filter(strings, predicate)
	assert.Len(t, filtered, 1)
}

func TestFiltersNumbers(t *testing.T) {
	strings := []int{
		1,
		2,
		3,
		4,
	}
	predicate := func(value int) bool {
		return value%2 == 0
	}
	filtered := Filter(strings, predicate)
	assert.Len(t, filtered, 2)
}

func TestFiltersStructs(t *testing.T) {
	type Person struct {
		Name string
	}
	strings := []Person{
		{Name: "Simpson"},
		{Name: "Humperdinck"},
	}
	predicate := func(value Person) bool {
		return value.Name == "Humperdinck"
	}
	filtered := Filter(strings, predicate)
	assert.Len(t, filtered, 1)
}

func TestMapsStructs(t *testing.T) {
	type Ditto struct {
		NickName string
	}
	type Mewtwo struct {
		NickName string
	}

	dittos := []Ditto{
		{NickName: "ADitto"},
	}
	predicate := func(value Ditto) (Mewtwo, error) {
		return Mewtwo{
			NickName: value.NickName,
		}, nil
	}
	mapped, err := Map(dittos, predicate)
	assert.Equal(t, []Mewtwo{{NickName: "ADitto"}}, mapped)
	assert.Nil(t, err)
}

func TestMapsStructsWithErrors(t *testing.T) {
	type Foo struct {
		ANumber int
	}
	type Bar struct {
		AnotherNumber int
	}
	errSomeKindOfError := errors.New("some kind of error")
	predicate := func(value Foo) (Bar, error) {
		if value.ANumber > 0 {
			return Bar{}, errSomeKindOfError
		}
		return Bar{
			AnotherNumber: value.ANumber,
		}, nil
	}
	foos := []Foo{
		{ANumber: 0},
		{ANumber: 1},
	}
	_, err := Map(foos, predicate)
	assert.True(t, errors.Is(err, errSomeKindOfError))
}
