package caches

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteCache struct {
	db *gorm.DB
}

func NewSqliteCache() (*SqliteCache, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	db.AutoMigrate(&VaultItem{})
	db.AutoMigrate(&Certificate{})

	if err != nil {
		return nil, err
	}

	return &SqliteCache{
		db: db,
	}, nil
}

func (s *SqliteCache) SaveVaultItem(vaultItem *VaultItem, updatedAt time.Time) error {
	result := s.db.Create(&vaultItem)
	return result.Error
}

func (s *SqliteCache) RetrieveVaultItem(fingerprint string) {

}

func (s *SqliteCache) Cleanup() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}
