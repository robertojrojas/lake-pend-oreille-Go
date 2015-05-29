package db

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	fmt.Printf("Initializing db...")
}

func GetDBConnection() (*sql.DB, error) {

	//TODO: Need to change the database location from an environment var.
	db, err := sql.Open("sqlite3", "./lpo.db")
	return db, err
}