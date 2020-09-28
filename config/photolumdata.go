package config

import (
	"database/sql"
	"fmt"

	"github.com/paulwrubel/photolum/database"
)

type PhotolumData struct {
	DB *sql.DB
}

func InitPhotolumData() (*PhotolumData, error) {
	fmt.Println("Initializing PhotolumData...")
	photolumData := new(PhotolumData)
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}
	photolumData.DB = db
	fmt.Println("PhotolumData initialized, returning...")
	return photolumData, nil
}
