package onepassword

import (
	"fmt"
	"sync"
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

func (o *OnePassword) FindCertificates() error {
	items, err := getListOfItems()

	if err != nil {
		return err
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

	for result := range results {
		if result.potentialError != nil {
			fmt.Println("Error:", result.potentialError)
			continue
		}
		fmt.Println(result.itemWithFields.Title)
	}

	return nil
}

func (o *OnePassword) retrieveItemAndCache(item Item) (ItemWithFields, error) {

	itemWithFields, err := getItemDetails(item.Id)
	if err != nil {
		return ItemWithFields{}, err
	}

	return itemWithFields, nil
}
