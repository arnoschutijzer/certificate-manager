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

	certificateAsString, _ := os.ReadFile("../../internal/domain/test_fixtures/RootCA.pem")

	certificate := domain.NewCertificate(string(certificateAsString), "aName")
	secret := domain.Secret{
		Id:           "an-id",
		Title:        "A Secret",
		UpdatedAt:    time.Now(),
		Certificates: []domain.Certificate{certificate},
	}

	cache.SaveSecret(secret)

	defer cache.Cleanup()
}
