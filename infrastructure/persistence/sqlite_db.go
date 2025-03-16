package persistence

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}
