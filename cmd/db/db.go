package main

import (
	"os"
	"time"

	"github.com/algleymi/certificate-manager/internal/caches"
	"github.com/algleymi/certificate-manager/internal/domain"
)

func main() {
	cache, err := caches.NewSqliteCache()

	if err != nil {
		panic(err)
	}

	certificateAsString, err := os.ReadFile("./internal/test_fixtures/RootCA.pem")
	certificate := domain.NewCertificate(string(certificateAsString), "aName")
	vaultItem := caches.Secret{
		Id:           "an-id",
		Title:        "A VaultItem",
		UpdatedAt:    time.Now(),
		Certificates: []caches.Certificate{caches.ToDbCertificate("an-id", certificate)},
	}

	cache.SaveSecret(vaultItem)

	defer cache.Cleanup()
}
