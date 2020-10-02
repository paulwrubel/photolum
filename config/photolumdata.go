package config

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulwrubel/photolum/constants"
	"github.com/paulwrubel/photolum/database"
	"github.com/sirupsen/logrus"
)

type PhotolumData struct {
	DB *pgxpool.Pool
}

func InitPhotolumData(log *logrus.Logger) (*PhotolumData, error) {
	log.Info("initializing PhotolumData")
	photolumData := new(PhotolumData)
	pgHost, isSet := os.LookupEnv(constants.PostgresHostnameEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresHostnameEnvironmentKey)
	}
	pgPass, isSet := os.LookupEnv(constants.PostgresPasswordEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresPasswordEnvironmentKey)
	}

	db, err := database.InitDB(log, pgHost, pgPass)
	if err != nil {
		return nil, fmt.Errorf("error initializing db: %s", err.Error())
	}

	err = database.InitSchema(log, db)
	if err != nil {
		return nil, fmt.Errorf("error initializing db schema: %s", err.Error())
	}

	photolumData.DB = db
	log.Info("PhotolumData initialized")
	return photolumData, nil
}
