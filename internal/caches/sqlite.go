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

func (s *SqliteCache) SaveVaultItem(vaultItem VaultItem) error {
	result := s.db.Create(&vaultItem)
	return result.Error
}

func (s *SqliteCache) RetrieveVaultItem(id string) {

}

func (s *SqliteCache) Cleanup() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}
