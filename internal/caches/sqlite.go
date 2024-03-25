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

	db.AutoMigrate(&Secret{})
	db.AutoMigrate(&Certificate{})

	if err != nil {
		return nil, err
	}

	return &SqliteCache{
		db: db,
	}, nil
}

func (s *SqliteCache) SaveSecret(secret Secret) error {
	result := s.db.Create(&secret)
	return result.Error
}

func (s *SqliteCache) RetrieveSecret(id string) (Secret, error) {
	var secret Secret
	result := s.db.Where(&Secret{Id: id}).Preload("Certificates").First(&secret)
	return secret, result.Error
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
