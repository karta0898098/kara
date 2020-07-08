package db

import "gorm.io/gorm"

type DatabaseType string

const (
	// MySQL ...
	MySQL DatabaseType = "mysql"
	// Postgres ...
	Postgres DatabaseType = "postgres"
)

type Config struct {
	Read    Database
	Write   Database
	Secrets string
	Color   bool
}

type Database struct {
	Host           string
	User           string
	Port           int
	Password       string
	Database       string
	Type           DatabaseType
	MaxIdleConn    int
	MaxLifetimeSec int
	ReadTimeout    string
	WriteTimeout   string
}

type Connection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}
