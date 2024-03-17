package main

import (
	"fmt"
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

func main() {
	GetThem(internal.NewOnePasswordStore())
}

func GetThem(store internal.Store) {
	firstOfJune2024UTCAt1AM := time.Date(2030, time.June, 1, 0, 0, 0, 0, time.UTC)
	certificates, err := store.FindCertificatesOlderThanDate(firstOfJune2024UTCAt1AM)

	if err != nil {
		panic(err)
	}

	if len(certificates) == 0 {
		fmt.Println("no outdated certificates, nice!")
		return
	}

	for _, v := range certificates {
		fmt.Printf("%s\n", v.Title)
	}
}
