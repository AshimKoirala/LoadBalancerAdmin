package db

import (
	"context"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

var ctx = context.Background()

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Username  string    `json:"username" bun:"username,unique,notnull"`
	Email     string    `json:"email" bun:"email,unique,notnull"`
	Password  string    `json:"password" bun:"password,notnull"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
	Password_Reset_Token   string    `bun:"password_reset_token"`       
    OTPExpiry  time.Time `bun:"otp_expiry"` 
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func InsertUser(user User) error {
	_, err := db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}
	return nil
}

func GetUsersinfo() ([]User, error) {
	var users []User
	err := db.NewSelect().Model(&users).Order("id ASC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
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
		Model(user).
		Where("username = ?", username).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return err
	}
	return nil
}

func GetUserByOTP(otp string) (User, error) {
	var user User
	err := db.NewSelect().
		Model(&user).
		Where("password_reset_token = ?", otp).
		Where("otp_expiry > ?", time.Now()). // Ensure OTP is still valid
		Limit(1).
		Scan(ctx)
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdatePassword(userID int64, hashedPassword string) error {
	_, err := db.NewUpdate().
		Model(&User{}).
		Set("password = ?", hashedPassword).
		Set("password_reset_token = NULL").  // Clear the OTP
		Set("otp_expiry = NULL").
		Where("id = ?", userID).
		Exec(ctx)
	return err
}
