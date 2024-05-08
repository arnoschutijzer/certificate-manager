package onepassword

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqliteCache struct {
	db *gorm.DB
}

type DatabaseItem struct {
	Id        string `gorm:"primaryKey"`
	Title     string
	UpdatedAt time.Time
	Content   string
}

func NewSqliteCache() (*SqliteCache, error) {
	db, err := gorm.Open(sqlite.Open("cache.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	db.AutoMigrate(&DatabaseItem{})

	if err != nil {
		return nil, err
	}

	return &SqliteCache{
		db: db,
	}, nil
}

func (s *SqliteCache) SaveItem(item ItemWithFields) error {
	field, err := item.findContentField()

	if err != nil {
		return err
	}

	dbItem := DatabaseItem{
		Id:        item.Id,
		Title:     item.Title,
		Content:   field.Value,
		UpdatedAt: item.UpdatedAt,
	}
	result := s.db.Create(&dbItem)

	return result.Error
}

func (s *SqliteCache) RetrieveItem(id string) (ItemWithFields, error) {
	var dbItem DatabaseItem
	result := s.db.Where(&DatabaseItem{Id: id}).First(&dbItem)

	if result.Error != nil {
		return ItemWithFields{}, result.Error
	}

	field := Field{
		Value: dbItem.Content,
	}

	item := Item{
		Id:        dbItem.Id,
		Title:     dbItem.Title,
		UpdatedAt: dbItem.UpdatedAt,
	}

	return ItemWithFields{
		Item: item,
		Fields: []Field{
			field,
		},
	}, result.Error
}

// func (s *SqliteCache) UpdateSecret(item) error {
// 	result := s.db.Save(&secret)
// 	return result.Error
// }
