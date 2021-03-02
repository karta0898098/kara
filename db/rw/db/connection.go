package db

import (
	"github.com/karta0898098/kara/db"
	"gorm.io/gorm"
)

type Connection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func NewConnection(config Config) (*Connection, error) {
	readDB, err := orm.SetupDatabase(&config.Read)
	if err != nil {
		return nil, err
	}
	writeDB, err := orm.SetupDatabase(&config.Write)
	if err != nil {
		return nil, err
	}

	return &Connection{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}, nil
}
