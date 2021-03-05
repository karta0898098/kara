package db

import (
	"github.com/karta0898098/kara/db"
	"gorm.io/gorm"
)

type Connection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func NewConnection(config Config) (Connection, error) {
	var (
		conn Connection
	)
	readDB, err := db.SetupDatabase(&config.Read)
	if err != nil {
		return conn, err
	}
	writeDB, err := db.SetupDatabase(&config.Write)
	if err != nil {
		return conn, err
	}

	return Connection{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}, nil
}
