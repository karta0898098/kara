package db

import (
	"github.com/karta0898098/kara/db"
)

type Config struct {
	Read    db.Database `mapstructure:"read"`
	Write   db.Database `mapstructure:"write"`
	Secrets string      `mapstructure:"secrets"`
}


