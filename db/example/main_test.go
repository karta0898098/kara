package main

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/karta0898098/kara/db"
	"github.com/karta0898098/kara/logging"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
)

type testSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestEndpoint(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupTest() {
	logging.Setup(logging.Config{
		Env:   "local",
		App:   "database",
		Level: logging.DebugLevel,
		Debug: true,
	})

	newLogger := db.NewLogger(dbLogger.Config{
		SlowThreshold: time.Second,   // Slow SQL threshold
		LogLevel:      dbLogger.Info, // Log level
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

func (s *testSuite) TestLog() {
	var (
		book Book
		ctx  context.Context
	)
	ctx = context.Background()
	// ctx = context.WithValue(ctx,"trace_id","123")

	logger := log.With().Str("trace_id", "123").Logger()
	ctx = logger.WithContext(ctx)

	err := s.db.WithContext(ctx).Model(&Book{}).Where("id = ?", 1).Find(&book).Error
	s.Equal(nil, err)
}
