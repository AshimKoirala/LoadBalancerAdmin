package db

import (
	"database/sql"
	"log"

	"github.com/AshimKoirala/load-balancer-admin/migrations"
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

	// Run migrations
	if err := migrations.RunMigrations(sqldb); err != nil {
		return err
	}

	log.Println("Database initialized and migrations applied.")
	return nil
}

func UpdateUser(user User) error {
	_, err := db.NewUpdate().
		Model(&user).
		Column("email", "password", "updated_at").
		Where("username = ?", user.Username).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(id int64) (User, error) {
	var user User
	err := db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return user, err
	}
	return user, nil
}

func GetUserByUsername(username string, user *User) error {
	err := db.NewSelect().
		Model(&user).
		Where("username = ?", username).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return err
	}
	return nil
}
