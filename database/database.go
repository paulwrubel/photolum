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
	connectionAttempts := 0
	var db *pgxpool.Pool
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		db, err = pgxpool.ConnectConfig(ctx, poolConfig)
		if err == nil {
			break
		}
		connectionAttempts++
		if connectionAttempts >= 10 {
			return nil, err
		}
		cancel()
		// retry db
		log.WithError(err).Error("database connection attempt failed, waiting 5s then retrying")
		time.Sleep(time.Second * 5)
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
