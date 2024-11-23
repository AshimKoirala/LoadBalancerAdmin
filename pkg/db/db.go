package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

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


	log.Println("Database initialized and migrations applied.")
	return nil
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) //6 digit code
}

func SetPasswordResetOTP(email string) (string, error) {
	otp := GenerateOTP()
	expiry := time.Now().Add(5 * time.Minute) //for 5 min

	_, err := db.NewUpdate().
		Model(&User{}).
		Set("password_reset_token = ?", otp).
		Set("otp_expiry = ?", expiry).
		Where("LOWER(email) = ?", email).
		Exec(ctx)
	if err != nil {
		return "", err
	}
	return otp, nil
}