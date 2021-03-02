package main

import (
	"github.com/karta0898098/kara/db"
	"github.com/karta0898098/kara/zlog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type Book struct {
	ID   int
	Name string
}

func main() {
	zlog.New(&zlog.Config{
		Env:   "local",
		AppID: "db_test",
		Debug: true,
	})

	conn, err := orm.NewConnection(&orm.Config{
		DB: orm.Database{
			Debug:    true,
			Host:     "127.0.0.1",
			User:     "rode",
			Port:     3306,
			Password: "rode@3306",
			Name:     "library",
			Type:     "mysql",
		},
	})

	if err != nil {
		log.Error().Msgf("db error: reason: %v", err)
	}

	ctx := context.WithValue(context.Background(), "request_id", "test")
	var book Book
	err = conn.WithContext(ctx).Model(&Book{}).Where("id = ?", 1).Find(&book).Error
	if err != nil {
		log.Error().Msgf("db error: reason: %v", err)
	}

	var book2 Book
	err = conn.WithContext(ctx).Model(&Book{}).Where("id = ?", 10).Find(&book2).Error
	if err != nil {
		log.Error().Msgf("db error: reason: %v", err)
	}
}
