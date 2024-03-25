package main

import (
	"os"
	"time"

	"github.com/algleymi/certificate-manager/internal"
	"github.com/algleymi/certificate-manager/internal/vaults/caches"
)

func main() {
	cache, err := caches.NewSqliteCache()

	if err != nil {
		panic(err)
	}

	certificateAsString, err := os.ReadFile("./internal/test_fixtures/RootCA.pem")
	certificate := internal.NewCertificate(string(certificateAsString), "aName")
	vaultItem := caches.Secret{
		Id:           "an-id",
		Title:        "A VaultItem",
		UpdatedAt:    time.Now(),
		Certificates: []caches.Certificate{caches.ToDbCertificate("an-id", certificate)},
	}

	cache.SaveSecret(vaultItem)

	defer cache.Cleanup()
}
