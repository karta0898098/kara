package db

import (
	"github.com/karta0898098/kara/db"
)

type Config struct {
	Read  orm.Database `mapstructure:"read"`
	Write orm.Database `mapstructure:"write"`
}
