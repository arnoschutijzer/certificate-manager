package onepassword

import (
	"errors"
	"sync"
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type OnePassword struct {
}

func NewOnePassword() (*OnePassword, error) {
	return &OnePassword{}, nil
}

type fetchResult struct {
	itemWithFields ItemWithFields
	potentialError error
}

var ErrNoCertificatesFound = errors.New("no certificates found")

func (o *OnePassword) FindCertificates(after time.Time) ([]internal.Certificate, error) {
	itemsWithFields, err := o.retrieveItemsWithFields()

	certificates, err := fetchCertificatesFromItemsWithFields(itemsWithFields)

	if err != nil {
		return nil, err
	}

	olderCertificates := internal.Filter(certificates, func(certificate internal.Certificate) bool {
		return !certificate.IsValid(after)
	})

	return olderCertificates, nil
}

func fetchCertificatesFromItemsWithFields(itemsWithFields []ItemWithFields) ([]internal.Certificate, error) {
	return internal.FlatMap(itemsWithFields, func(itemWithFields ItemWithFields) ([]internal.Certificate, error) {
		content, err := itemWithFields.findContentField()
		if err != nil {
			return nil, err
		}
		certificates := internal.GetCertificatesFromString(content.Value, itemWithFields.Title)
		return certificates, nil
	})
}

func (o *OnePassword) retrieveItemsWithFields() ([]ItemWithFields, error) {
	items, err := getListOfItems()

	if err != nil {
		return nil, nil
	}

	results := make(chan fetchResult)

	var wg sync.WaitGroup
	wg.Add(len(items))

	go func() {
		wg.Wait()
		close(results)
	}()

	for _, v := range items {
		go func(item Item, results chan<- fetchResult) {
			defer wg.Done()

			itemWithFields, err := o.retrieveItemAndCache(item)

			result := fetchResult{
				itemWithFields: itemWithFields,
				potentialError: err,
			}

			results <- result

		}(v, results)
	}

	itemsWithFields := []ItemWithFields{}
	errorsForFields := []error{}

	for result := range results {
		if result.potentialError != nil {
			errorsForFields = append(errorsForFields, result.potentialError)
			continue
		}

		itemsWithFields = append(itemsWithFields, result.itemWithFields)
	}

	return itemsWithFields, errors.Join(errorsForFields...)
}

func (o *OnePassword) retrieveItemAndCache(item Item) (ItemWithFields, error) {

	itemWithFields, err := getItemDetails(item.Id)
	if err != nil {
		return ItemWithFields{}, err
	}

	return itemWithFields, nil
}
