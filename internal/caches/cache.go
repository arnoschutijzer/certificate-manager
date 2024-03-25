package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal/domain"
)

type Cache interface {
	SaveSecret(secret domain.Secret) error
	RetrieveSecret(id string) (domain.Secret, error)
	UpdateSecret(secret Secret) error
	Cleanup() error
}

type Secret struct {
	Id           string `gorm:"primaryKey"`
	Title        string
	UpdatedAt    time.Time
	Certificates []Certificate `gorm:"foreignKey:SecretId"`
}

type Certificate struct {
	SecretId    string
	NotBefore   time.Time
	NotAfter    time.Time
	Subject     string
	CustomName  string
	Fingerprint string `gorm:"primaryKey"`
}
