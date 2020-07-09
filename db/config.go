package db

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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
}

type Database struct {
	Debug          bool
	Host           string
	User           string
	Port           int
	Password       string
	Name           string
	Type           DatabaseType
	MaxIdleConns   int
	MaxOpenConns   int
	MaxLifetimeSec int
	ReadTimeout    string
	WriteTimeout   string
}

func (c *Config) New() Config {
	return *c
}
