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

func (s *SqliteCache) SaveCertificate(certificate internal.Certificate) error {
	dbCertificate := ToDatabaseCertificate(certificate)
	result := s.db.Create(&dbCertificate)
	return result.Error
}

func (s *SqliteCache) RetrieveCertificate(fingerprint string) {

}

func (s *SqliteCache) Cleanup() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}
