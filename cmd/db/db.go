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
	vaultItem := caches.CachedItem{
		VaultId:      "an-id",
		Title:        "A VaultItem",
		UpdatedAt:    time.Now(),
		Certificates: []caches.CachedCertificate{caches.ToDbCertificate("an-id", certificate)},
	}

	cache.SaveVaultItem(vaultItem)

	defer cache.Cleanup()
}
