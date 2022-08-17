package postgresql

// import (
// 	"time"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/pkg/errors"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/yurttasutkan/alarmservice/internal/config"
// )

// // Integration implements a PostgreSQL integration.
// type Integration struct {
// 	db *sqlx.DB
// }

// // New creates a new PostgreSQL integration.
// func New(conf config.IntegrationPostgreSQLConfig) (*Integration, error) {
// 	log.Info("integration/postgresql: connecting to PostgreSQL database")
// 	d, err := sqlx.Open("postgres", conf.DSN)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "integration/postgresql: PostgreSQL connection error")
// 	}
// 	for {
// 		if err := d.Ping(); err != nil {
// 			log.WithError(err).Warning("integration/postgresql: ping PostgreSQL database error, will retry in 2s")
// 			time.Sleep(2 * time.Second)
// 		} else {
// 			break
// 		}
// 	}

// 	d.SetMaxOpenConns(conf.MaxOpenConnections)
// 	d.SetMaxIdleConns(conf.MaxIdleConnections)

// 	return &Integration{
// 		db: d,
// 	}, nil
// }

// // Close closes the integration.
// func (i *Integration) Close() error {
// 	if err := i.db.Close(); err != nil {
// 		return errors.Wrap(err, "close database error")
// 	}
// 	return nil
// }