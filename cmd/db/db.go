package main

import (
	"os"
	"time"

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
	vaultItem := &caches.VaultItem{
		Id:           "an-id",
		Title:        "A VaultItem",
		UpdatedAt:    time.Now(),
		Certificates: []caches.Certificate{caches.ToDbCertificate("an-id", certificate)},
	}

	cache.SaveVaultItem(vaultItem, time.Now())

	defer cache.Cleanup()
}
