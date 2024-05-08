package onepassword

import (
	"github.com/algleymi/certificate-manager/internal"
)

type OnePassword struct {
	cache *SqliteCache
}

func NewOnePassword() (*OnePassword, error) {
	cache, err := NewSqliteCache()

	if err != nil {
		return nil, err
	}

	return &OnePassword{
		cache,
	}, nil
}

func (o *OnePassword) FindCertificates() error {
	items, err := getListOfItems()

	if err != nil {
		return err
	}

	_, err = internal.Map(items, o.retrieveItemAndCache)

	if err != nil {
		return err
	}

	return nil
}

func (o *OnePassword) retrieveItemAndCache(item Item) (ItemWithFields, error) {
	cached, err := o.cache.RetrieveItem(item.Id)

	if err != nil {
		itemWithFields, err := o.retrieveSecretAndCache(item.Id)

		if err != nil {
			return ItemWithFields{}, err
		}

		return itemWithFields, nil
	}

	if !item.UpdatedAt.After(cached.UpdatedAt) {
		return cached, nil
	}

	itemWithFields, err := o.retrieveSecretAndCache(item.Id)

	if err != nil {
		return ItemWithFields{}, err
	}

	return itemWithFields, nil
}

func (s *OnePassword) retrieveSecretAndCache(id string) (ItemWithFields, error) {
	itemWithFields, err := getItemDetails(id)
	if err != nil {
		return ItemWithFields{}, err
	}

	s.cache.SaveItem(itemWithFields)

	return itemWithFields, nil
}
