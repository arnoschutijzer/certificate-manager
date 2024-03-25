package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/algleymi/certificate-manager/internal/caches"
	vaults "github.com/algleymi/certificate-manager/internal/vaults"
)

func main() {
	cache, err := caches.NewSqliteCache()

	if err != nil {
		panic(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigc

		cache.Cleanup()
	}()

	store := vaults.NewOnePasswordStore(cache)

	firstOfJune2024UTCAt1AM := time.Date(2030, time.June, 1, 0, 0, 0, 0, time.UTC)
	certificates, err := store.FindCertificatesOlderThanDate2(firstOfJune2024UTCAt1AM)

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
