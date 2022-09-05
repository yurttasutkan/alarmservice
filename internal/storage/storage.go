package storage

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/yurttasutkan/alarmservice/internal/config"
)

var (
	db *sqlx.DB
)

// DB returns the DB instance.
func DB() *sqlx.DB {
	return db
}

// Setup configures the storage package.
func Setup(conf *config.Config) error {
	log.Info("storage: connecting to PostgreSQL database")
	d, err := sqlx.Open("postgres", conf.PostgreSQL.DSN)
	if err != nil {
		return fmt.Errorf("open postgresql connection error: %w", err)
	}
	d.SetMaxOpenConns(conf.PostgreSQL.MaxOpenConnections)
	d.SetMaxIdleConns(conf.PostgreSQL.MaxIdleConnections)
	for {
		if err := d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(time.Second * 2)
		} else {
			log.Infof("Connected to PostgreSQL")
			break
		}
	}

	db = d

	return nil
}

