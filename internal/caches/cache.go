package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type Cache interface {
	SaveVaultItem(vaultItem VaultItem, updatedAt time.Time) error
	RetrieveVaultItem(id string)
	Cleanup() error
}

// Decouple vaults from db schema
type VaultItem struct {
	Id           string
	Title        string
	UpdatedAt    time.Time
	Certificates []Certificate `gorm:"foreignKey:VaultItem"`
}

type Certificate struct {
	VaultItem   string
	NotBefore   time.Time
	NotAfter    time.Time
	Subject     string
	CustomName  string
	Fingerprint string
}

func ToDbCertificate(id string, certificate internal.Certificate) Certificate {
	return Certificate{
		VaultItem:   id,
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}
