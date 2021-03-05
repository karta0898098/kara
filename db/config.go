package db

import (
	"gorm.io/gorm"
)

type Config struct {
	DB      Database `mapstructure:"db"`
	Secrets string   `mapstructure:"secrets"`
}

func NewConnection(config Config) (*gorm.DB, error) {
	db, err := SetupDatabase(&config.DB)
	if err != nil {
		return nil, err
	}

	return db, nil
}
