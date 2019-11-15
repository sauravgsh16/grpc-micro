package v1

import (
	"database/sql"
	"fmt"
)

func CreateDB(db *sql.DB) error {
	// Create DB
	_, err := db.Exec("CREATE DATABASE `ToDoDB`")
	if err != nil {
		return fmt.Errorf("Failed to create db: %v", err)
	}

}
