package main

import (
	"fmt"
	"time"

	vaults "github.com/algleymi/certificate-manager/internal/vaults"
)

func main() {
	store := vaults.NewOnePasswordStore()

	firstOfJune2024UTCAt1AM := time.Date(2030, time.June, 1, 0, 0, 0, 0, time.UTC)
	certificates, err := store.FindCertificatesOlderThanDate(firstOfJune2024UTCAt1AM)

	if err != nil {
		panic(err)
	}

	numberOfCertificates := len(certificates)

	if numberOfCertificates == 0 {
		fmt.Println("no outdated certificates, nice!")
		return
	}

	fmt.Printf("found %d outdated certificates\n", numberOfCertificates)
	for _, v := range certificates {
		fmt.Printf("%s\n", v.CustomName)
	}
}
