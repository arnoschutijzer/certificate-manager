package caches

import (
	"github.com/algleymi/certificate-manager/internal"
	"github.com/algleymi/certificate-manager/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteCache struct {
	db *gorm.DB
}

var _ Cache = &SqliteCache{}

func NewSqliteCache() (*SqliteCache, error) {
	db, err := gorm.Open(sqlite.Open("cache.db"), &gorm.Config{})

	db.AutoMigrate(&Secret{})
	db.AutoMigrate(&Certificate{})

	if err != nil {
		return nil, err
	}

	return &SqliteCache{
		db: db,
	}, nil
}

func (s *SqliteCache) SaveSecret(secret domain.Secret) error {
	certificates, _ := internal.Map(secret.Certificates, func(certificate domain.Certificate) (Certificate, error) {
		return Certificate{
			Fingerprint: certificate.Fingerprint,
			Subject:     certificate.Subject,
			CustomName:  certificate.CustomName,
			NotAfter:    certificate.NotAfter,
			NotBefore:   certificate.NotBefore,
		}, nil
	})
	dbSecret := Secret{
		Id:           secret.Id,
		Title:        secret.Title,
		UpdatedAt:    secret.UpdatedAt,
		Certificates: certificates,
	}
	result := s.db.Create(&dbSecret)
	return result.Error
}

func (s *SqliteCache) RetrieveSecret(id string) (domain.Secret, error) {
	var secret Secret
	result := s.db.Where(&Secret{Id: id}).Preload("Certificates").First(&secret)

	if result.Error != nil {
		return domain.Secret{}, result.Error
	}

	certificates, _ := internal.Map(secret.Certificates, func(certificate Certificate) (domain.Certificate, error) {
		return domain.Certificate{
			Fingerprint: certificate.Fingerprint,
			Subject:     certificate.Subject,
			CustomName:  certificate.CustomName,
			NotAfter:    certificate.NotAfter,
			NotBefore:   certificate.NotBefore,
		}, nil
	})

	return domain.Secret{
		Id:           secret.Id,
		Title:        secret.Title,
		UpdatedAt:    secret.UpdatedAt,
		Certificates: certificates,
	}, result.Error
}

func (s *SqliteCache) UpdateSecret(secret Secret) error {
	result := s.db.Save(&secret)
	return result.Error
}

func (s *SqliteCache) Cleanup() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}
