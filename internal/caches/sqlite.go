package caches

import (
	"github.com/algleymi/certificate-manager/internal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteCache struct {
	db *gorm.DB
}

func NewSqliteCache() (*SqliteCache, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	db.AutoMigrate(&DatabaseCertificate{})

	if err != nil {
		return nil, err
	}

	return &SqliteCache{
		db: db,
	}, nil
}

// SaveCertificate(certificate i.Certificate)
// RetrieveCertificate(fingerprint string)

func (s *SqliteCache) SaveCertificate(certificate internal.Certificate) error {
	dbCertificate := ToDatabaseCertificate(certificate)

	tx := s.db.Begin()
	result := s.db.Create(&dbCertificate)
	tx.Commit()

	return result.Error
}

func (s *SqliteCache) RetrieveCertificate(fingerprint string) {

}
