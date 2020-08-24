package db

import (
	"github.com/jinzhu/gorm"
	"github.com/karta0898098/kara/db"
)

type Connection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func NewConnection(config *Config) (*Connection, error) {
	readDB, err := db.SetupDatabase(&config.Read)
	if err != nil {
		return nil, err
	}
	writeDB, err := db.SetupDatabase(&config.Write)
	if err != nil {
		return nil, err
	}

	readDB.LogMode(config.Read.Debug)
	writeDB.LogMode(config.Write.Debug)

	return &Connection{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}, nil
}


