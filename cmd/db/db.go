package main

import (
	"os"

	"github.com/algleymi/certificate-manager/internal"
	"github.com/algleymi/certificate-manager/internal/caches"
)

func main() {
	cache, err := caches.NewSqliteCache()

	if err != nil {
		panic(err)
	}

	certificateAsString, err := os.ReadFile("./internal/test_fixtures/RootCA.pem")
	certificate := internal.NewCertificate(string(certificateAsString), "aName")
	cache.SaveCertificate(certificate)

	defer cache.Cleanup()
}
