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
	Read    Database `mapstructure:"read"`
	Write   Database `mapstructure:"write"`
	Secrets string   `mapstructure:"secrets"`
}

type Database struct {
	Debug          bool         `mapstructure:"debug"`
	Host           string       `mapstructure:"host"`
	User           string       `mapstructure:"user"`
	Port           int          `mapstructure:"port"`
	Password       string       `mapstructure:"password"`
	Name           string       `mapstructure:"name"`
	Type           DatabaseType `mapstructure:"type"`
	MaxIdleConns   int          `mapstructure:"max_idle_conns"`
	MaxOpenConns   int          `mapstructure:"max_open_conns"`
	MaxLifetimeSec int          `mapstructure:"max_lifetime"`
	ReadTimeout    string       `mapstructure:"read_timeout"`
	WriteTimeout   string       `mapstructure:"write_timeout"`
}

func (c *Config) New() *Config {
	return c
}
