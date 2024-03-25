package caches

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteCache struct {
	db *gorm.DB
}

var _ Cache = &SqliteCache{}

func NewSqliteCache() (*SqliteCache, error) {
	db, err := gorm.Open(sqlite.Open("cache.db"), &gorm.Config{})

	db.AutoMigrate(&CachedItem{})
	db.AutoMigrate(&CachedCertificate{})

	if err != nil {
		return nil, err
	}

	return &SqliteCache{
		db: db,
	}, nil
}

func (s *SqliteCache) SaveVaultItem(vaultItem CachedItem) error {
	result := s.db.Create(&vaultItem)
	return result.Error
}

func (s *SqliteCache) RetrieveVaultItem(id string) (CachedItem, error) {
	var vaultItem CachedItem
	result := s.db.Where(&CachedItem{VaultId: id}).Preload("Certificates").First(&vaultItem)
	return vaultItem, result.Error
}

func (s *SqliteCache) UpdateVaultItem(vaultItem CachedItem) error {
	result := s.db.Save(&vaultItem)
	return result.Error
}

func (s *SqliteCache) Cleanup() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}
