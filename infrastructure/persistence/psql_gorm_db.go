package persistence

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"my-project/infrastructure/logger"
)

func NewPsqlGormDb() *gorm.DB {
	// Create a new connection to the database postgres using gorm
	// The connection string is defined in the .env file
	dsn := "host=localhost user=project password=MyPassword_123 dbname=social_media_wrap port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.GetLogger().WithField("package", "infrastructure/persistence/psql_gorm_db").Error(err)
	}

	return db
}
