package database

import (
	"context"
	"fmt"
	"io/ioutil"

	// Phantom import for PostgreSQL driver

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitDB(log *logrus.Logger, pgHost string, pgPassword string) (*pgxpool.Pool, error) {
	log.Debug("initializing database")

	// initialize configuration
	connectionString := fmt.Sprintf("host=%s password=%s", pgHost, pgPassword)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// initialize connection pool
	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
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
