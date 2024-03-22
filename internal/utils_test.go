package internal

import (
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
	predicate := func(value Ditto) Mewtwo {
		return Mewtwo{
			NickName: value.NickName,
		}
	}
	mapped := Map(dittos, predicate)
	assert.Equal(t, []Mewtwo{{NickName: "ADitto"}}, mapped)
}
