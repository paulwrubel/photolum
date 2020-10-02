package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	// Phantom import for PostgreSQL driver

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitDB(log *logrus.Logger, pgHost, pgUser, pgPassword string) (*pgxpool.Pool, error) {
	log.Debug("initializing database")

	// initialize configuration
	connectionString := fmt.Sprintf("host=%s user=%s password=%s", pgHost, pgUser, pgPassword)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// initialize connection pool
	var db *pgxpool.Pool
	connectionAttempts := 0
	for db, err = pgxpool.ConnectConfig(context.Background(), poolConfig); ; {
		if err == nil {
			break
		}
		connectionAttempts++
		if connectionAttempts >= 5 {
			return nil, err
		}
		// retry db
		log.WithError(err).Error("database connection attempt failed, waiting 5s then retrying")
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		return nil, err
	}

	log.Debug("database initialized")
	return db, nil
}

func InitSchema(log *logrus.Logger, db *pgxpool.Pool) error {
	log.Debug("initializing database schema")

	// get queries from sql file
	schemaFileBytes, err := ioutil.ReadFile("/app/schema.sql")
	if err != nil {
		return err
	}

	// execute queries in file
	_, err = db.Exec(context.Background(), string(schemaFileBytes))
	if err != nil {
		return err
	}

	log.Debug("database schema initialized")
	return nil
}
