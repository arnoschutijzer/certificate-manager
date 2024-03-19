package main

import (
	"os"

	"github.com/algleymi/certificate-manager/internal"
	"github.com/algleymi/certificate-manager/internal/caches"
)

func main() {
	db, err := caches.NewSqliteCache()

	if err != nil {
		panic(err)
	}

	certificateAsString, err := os.ReadFile("./internal/test_fixtures/RootCA.pem")
	certificate := internal.NewCertificate(string(certificateAsString), "aName")
	db.SaveCertificate(certificate)
}
