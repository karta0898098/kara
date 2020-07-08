package db

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"strconv"
	"time"
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

type Connection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func NewConnection(cfg Config) (*Connection, error) {

	readDB, err := setupDatabase(&cfg.Read)
	if err != nil {
		return nil, err
	}
	writeDB, err := setupDatabase(&cfg.Write)
	if err != nil {
		return nil, err
	}
	readDB.LogMode(cfg.Read.Debug)
	writeDB.LogMode(cfg.Write.Debug)

	return &Connection{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}, nil
}

func setupDatabase(database *Database) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	if database.WriteTimeout == "" {
		database.WriteTimeout = "10s"
	}

	if database.ReadTimeout == "" {
		database.ReadTimeout = "10s"
	}

	driver := ""
	dsn := ""
	switch database.Type {
	case MySQL:
		driver = string(MySQL)
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true&readTimeout=%s&writeTimeout=%s", database.User, database.Password, database.Host+":"+strconv.Itoa(database.Port), database.Name, database.ReadTimeout, database.WriteTimeout)
	case Postgres:
		driver = string(Postgres)
		dsn = fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable `, database.User, database.Password, database.Host, database.Port, database.Name)
	default:
		return nil, errors.New("Not support driver")
	}

	var conn *gorm.DB

	// 嘗試重新連線database
	err := backoff.Retry(func() error {
		db, err := gorm.Open(driver, dsn)
		if err != nil {
			return err
		}
		err = db.DB().Ping()
		if err != nil {
			return err
		}

		conn = db
		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	// set default idle conn
	if database.MaxIdleConns == 0 {
		database.MaxIdleConns = 10
	}

	if database.MaxOpenConns == 0 {
		database.MaxOpenConns = 20
	}

	if database.MaxLifetimeSec == 0 {
		database.MaxLifetimeSec = 14400
	}

	conn.DB().SetMaxIdleConns(database.MaxIdleConns)
	conn.DB().SetMaxOpenConns(database.MaxOpenConns)
	conn.DB().SetMaxIdleConns(database.MaxLifetimeSec)

	return conn, nil
}
