package main

import (
	"cmp"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/algleymi/certificate-manager/internal"
	vaults "github.com/algleymi/certificate-manager/internal/vaults"
	"github.com/algleymi/certificate-manager/internal/vaults/caches"
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

	inTwoMonths := time.Now().AddDate(0, 2, 0)

	fmt.Println("Finding certificates...")
	before := time.Now()
	certificates, err := store.FindCertificatesOlderThanDate(inTwoMonths)
	fmt.Printf("Treated all certificates, took %f\n", time.Since(before).Seconds())

	if err != nil {
		panic(err)
	}

	numberOfCertificates := len(certificates)

	if numberOfCertificates == 0 {
		fmt.Println("No outdated certificates, nice!")
		return
	}

	slices.SortFunc(certificates, func(a, b internal.Certificate) int {
		return cmp.Compare(strings.ToLower(a.CustomName), strings.ToLower(b.CustomName))
	})

	fmt.Printf("Found %d outdated certificates\n", numberOfCertificates)
	for _, v := range certificates {
		fmt.Printf("%s - %s\n", v.CustomName, v.Subject)
	}
}
