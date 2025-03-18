package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

const (
	MaxIDLETime            = 1 * time.Hour
	ConnectionMaxLifetime  = 24 * time.Hour
	MaxIDLEConnectionCount = 10
	MaxOpenConnectionCount = 20
)

var (
	connection *gorm.DB
	once       sync.Once
	err        error
)

func NewPgSQLConnection(cfg *Config) (*gorm.DB, error) {
	var database *sql.DB

	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)

		connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err != nil {
			return
		}

		database, err = connection.DB()

		if err != nil {
			return
		}

		database.SetConnMaxIdleTime(MaxIDLETime)
		database.SetConnMaxLifetime(ConnectionMaxLifetime)
		database.SetMaxIdleConns(MaxIDLEConnectionCount)
		database.SetMaxOpenConns(MaxOpenConnectionCount)

		if err = goose.Up(database, "./migrations", goose.WithAllowMissing()); err != nil {
			return
		}
	})

	return connection, err
}
