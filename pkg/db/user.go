package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

var ctx = context.Background()

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id                   int64     `json:"id" bun:"id,pk,autoincrement"`
	Username             string    `json:"username" bun:"username,unique,notnull"`
	Email                string    `json:"email" bun:"email,unique,notnull"`
	Password             string    `json:"password" bun:"password,notnull"`
	CreatedAt            time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt            time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
	Password_Reset_Token string    `bun:"password_reset_token"`
	token_expires_at     time.Time `bun:"token_expires_at"`
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
	query := db.NewUpdate().Model(&user).Column("updated_at")

	if user.Username != "" {
		query = query.Column("username")
	}

	if user.Password != "" {
		query = query.Column("password")
	}

	_, err := query.Where("id = ?", user.Id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func GetUserById(id int64) (User, error) {
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

func GetUserByOtp(otp string) (User, error) {
	var user User
	err := db.NewSelect().
		Model(&user).
		Where("password_reset_token = ?", otp).
		Where("token_expires_at > ?", time.Now()). // Ensure OTP is still valid
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
		Set("password_reset_token = NULL"). // Clear the OTP
		Set("token_expires_at = NULL").
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

func GenerateOtp() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000)) //6 digit code
}

func SetPasswordResetOtp(email string) (string, error) {
	otp := GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute) //for 5 min

	_, err := db.NewUpdate().
		Model(&User{}).
		Set("password_reset_token = ?", otp).
		Set("token_expires_at = ?", expiry).
		Where("LOWER(email) = ?", email).
		Exec(ctx)
	if err != nil {
		return "", err
	}
	return otp, nil
}
