package main

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/karta0898098/kara/db"
	"github.com/karta0898098/kara/zlog"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

type testSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestEndpoint(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupTest() {
	zlog.Setup(&zlog.Config{
		Env:   "local",
		AppID: "database",
		Level: int8(zerolog.DebugLevel),
		Debug: true,
	})

	newLogger := db.NewLogger(logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logger.Info, // Log level
		Colorful:      true,        // Disable color
	})

	returnValue := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	mockConn, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mock.
		ExpectQuery("SELECT * FROM `books` WHERE id = ?").
		WithArgs(1).
		WillReturnRows(returnValue)

	s.db, _ = gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockConn,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: newLogger,
	})
}

func (s *testSuite) Test_log() {
	var (
		book Book
		ctx  context.Context
	)
	ctx = context.Background()

	err := s.db.WithContext(ctx).Model(&Book{}).Where("id = ?", 1).Find(&book).Error
	s.Equal(nil, err)
}
