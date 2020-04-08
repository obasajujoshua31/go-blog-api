package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-blog-api/config"
)

type DB struct {
	*gorm.DB
}

func NewDB(config *config.AppConfig) (*DB, error) {
	db, err := gorm.Open(config.Dialect, getConnString(config))
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func getConnString(appConfig *config.AppConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		appConfig.DBHost, appConfig.DBPort,
		appConfig.DBUser, appConfig.DBName,
		appConfig.DBPassword,
		appConfig.SSLMode,
	)
}

func (d *DB) CreateTables(in interface{}) error {
	d.DropTableIfExists(in)
	if err := d.Debug().AutoMigrate(in).Error; err != nil {
		return err
	}

	return nil
}
