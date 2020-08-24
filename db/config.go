package db

import (
	"github.com/jinzhu/gorm"
)

type Config struct {
	DB      Database `mapstructure:"db"`
	Secrets string   `mapstructure:"secrets"`
}

func NewConnection(config *Config) (*gorm.DB, error) {
	db, err := SetupDatabase(&config.DB)
	if err != nil {
		return nil, err
	}

	db.LogMode(config.DB.Debug)
	db.LogMode(config.DB.Debug)

	return db, nil
}
