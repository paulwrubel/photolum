package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	// Phantom import for SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

var createImageTableQueryString string = `
	CREATE TABLE image (
		scene_id TEXT NOT NULL,
		file_type TEXT NOT NULL,
		created_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		modified_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		accessed_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		image_data BLOB NOT NULL,
		PRIMARY KEY (scene_id, file_type)
	)
`

var createSceneTableQueryString string = `
	CREATE TABLE scene (
		scene_id TEXT PRIMARY KEY,
		render_status TEXT NOT NULL,
		created_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		modified_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		accessed_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		image_width INTEGER NOT NULL,
		image_height INTEGER NOT NULL,
		image_file_types TEXT NOT NULL
	)
`

func InitDB() (*sql.DB, error) {
	fmt.Println("Initializing DB...")

	// Remove db if exists
	err := os.Remove("/app/photolum.db")
	if err != nil {
		fmt.Printf("Error removing db file (we probably don't care): %s\n", err.Error())
	}

	// Create db
	db, err := sql.Open("sqlite3", "/app/photolum.db")
	if err != nil {
		return nil, err
	}

	// Create the DB schema
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// create scene table
	sceneStmt, err := tx.Prepare(createSceneTableQueryString)
	if err != nil {
		return nil, err
	}
	defer sceneStmt.Close()
	_, err = sceneStmt.Exec()
	if err != nil {
		return nil, err
	}

	// create image table
	imageStmt, err := tx.Prepare(createImageTableQueryString)
	if err != nil {
		return nil, err
	}
	defer imageStmt.Close()
	_, err = imageStmt.Exec()
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	fmt.Println("DB initialized, returning...")
	return db, nil
}