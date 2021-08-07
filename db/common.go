package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	defaultDatabase *gorm.DB

	Logger logger.Interface
)

type CommonConfig struct {
	DBEndpoints string `json:"DB_ENDPOINTS" env:"DB_ENDPOINTS,required"`
	DBType      string `json:"DB_TYPE" env:"DB_TYPE,required"`
}

func InitDefaultDatabase(databaseType, endPoint string) error {
	var db *gorm.DB
	var err error
	switch databaseType {
	case "sqlite":
		// NOTE: In memory -> file::memory:?cache=shared
		db, err = gorm.Open(sqlite.Open(endPoint), &gorm.Config{
			Logger: Logger,
		})
		break
	case "mysql":
		db, err = gorm.Open(mysql.Open(endPoint), &gorm.Config{
			Logger: Logger,
		})
	case "postgres":
		db, err = gorm.Open(postgres.Open(endPoint), &gorm.Config{
			Logger: Logger,
		})
		break
	}
	if err == nil && db != nil {
		defaultDatabase = db
	}

	return err
}

func SetLogger(logger logger.Interface) {
	Logger = logger
	defaultDatabase.Logger = logger
}

func GetDefaultDatabase() *gorm.DB {
	return defaultDatabase
}
