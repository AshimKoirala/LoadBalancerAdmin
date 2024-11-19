package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"sync"
	"time"
	"unicode"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

var (
	users        = make(map[string]db.User)
	usersMutex   sync.Mutex
	emailRegex   = regexp.MustCompile(`^[a-z0-9._+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	psymbolRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	var user db.User
	var validationErrors []string

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		validationErrors = append(validationErrors, "Invalid request payload")
		utils.NewErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	// Validate username's length
	if len(user.Username) < 3 || len(user.Username) > 32 {
		validationErrors = append(validationErrors, "Username must be at least 3-32 characters long")
	}

	// Validate email format
	if !emailRegex.MatchString(user.Email) {
		validationErrors = append(validationErrors, "Invalid email format")
	}

	// Validate password
	if len(user.Password) < 8 || len(user.Password) > 32 {
		validationErrors = append(validationErrors, "Password must be at least 8 characters long")
	}

	if !containsCapitalLetter(user.Password) {
		validationErrors = append(validationErrors, "Password must contain at least one capital letter")
	}

	if !containsSymbol(user.Password) {
		validationErrors = append(validationErrors, "Password must contain at least one symbol")
	}

	// Check for errors
	if len(validationErrors) > 0 {
		utils.NewErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	// Check if user already exists
	if _, exists := users[user.Username]; exists {
		validationErrors = append(validationErrors, "User already exists")
		utils.NewErrorResponse(w, http.StatusConflict, validationErrors)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error hashing password"})
		return
	}

	// Save user with hashed password
	user.Password = hashedPassword

	// Insert the user into the database
	if err := db.InsertUser(user); err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error inserting user into the database"})
		return
	}

	utils.NewSuccessResponse(w, "User registered successfully")
}

func containsCapitalLetter(password string) bool {
	for _, ch := range password {
		if unicode.IsUpper(ch) {
			return true
		}
	}
	return false
}

func containsSymbol(password string) bool {
	return psymbolRegex.MatchString(password)
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	var creds db.Credentials
	var validationErrors []string

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		validationErrors = append(validationErrors, "Invalid request payload")
		utils.NewErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	// Retrieve user from the database
	var user db.User
	err = db.GetUserByUsername(creds.Username, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			validationErrors = append(validationErrors, "Could not find user")
			utils.NewErrorResponse(w, http.StatusUnauthorized, validationErrors)
			return
		}
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error fetching user"})
		return
	}

	// Validate password
	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		validationErrors = append(validationErrors, "Invalid credentials")
		utils.NewErrorResponse(w, http.StatusUnauthorized, validationErrors)
		return
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error generating token"})
		return
	}

	// Set JWT as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
	})

	// Send success response
	utils.NewSuccessResponse(w, "Login successful")
}

func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	w.Write([]byte("Welcome to load balancer , " + username))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	// Fetch the list of users from the database
	users, err := db.GetUsersinfo()
	if err != nil {
		// Send error response in case of database fetching error
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error fetching users"})
		return
	}

	// If no users are found, return an empty array
	if len(users) == 0 {
		utils.NewSuccessResponse(w, []string{})
		return
	}

	// Send the user data
	utils.NewSuccessResponse(w, users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updatedUser db.User
	var validationErrors []string

	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		validationErrors = append(validationErrors, "Invalid request payload")
		utils.NewErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	user, err := db.GetUserByID(updatedUser.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			validationErrors = append(validationErrors, "User not found")
			utils.NewErrorResponse(w, http.StatusNotFound, validationErrors)
			return
		}
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error fetching user"})
		return
	}

	// Validate and update fields
	if updatedUser.Username != "" {
		user.Username = updatedUser.Username
	}

	if updatedUser.Email != "" {
		if !emailRegex.MatchString(updatedUser.Email) {
			validationErrors = append(validationErrors, "Invalid email format")
		} else {
			user.Email = updatedUser.Email
		}
	}

	// Update password
	if updatedUser.Password != "" {
		hashedPassword, err := utils.HashPassword(updatedUser.Password)
		if err != nil {
			utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error hashing password"})
			return
		}
		user.Password = hashedPassword
	}

	// If there are validation errors, send an error response
	if len(validationErrors) > 0 {
		utils.NewErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	// Update the timestamp
	user.UpdatedAt = time.Now()

	// Save the updated user
	err = db.UpdateUser(user)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Error updating user"})
		return
	}

	utils.NewSuccessResponse(w, "User information updated successfully")
}
