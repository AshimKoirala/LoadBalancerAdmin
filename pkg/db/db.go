package db

import (
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var db *bun.DB

func InitDB() error {
	dsn := "postgres://postgres:postgres@localhost:5432/prequal?sslmode=disable"
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return err
	}

	// Wrap sql.DB with Bun
	db = bun.NewDB(sqldb, pgdialect.New())

	// Run migrations uncomment when first creating
	// if err := migrations.RunMigrations(sqldb); err != nil {
	// 	return err
	// }

	log.Println("Database initialized and migrations applied.")
	return nil
}
