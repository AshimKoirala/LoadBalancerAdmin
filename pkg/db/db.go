package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var db *bun.DB

func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return err
	}

	// Wrap sql.DB with Bun
	db = bun.NewDB(sqldb, pgdialect.New())


	log.Println("Database initialized and migrations applied.")
	return nil
}
