package main

import (
	"log"
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/handlers"
	"github.com/AshimKoirala/load-balancer-admin/middleware"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

func main() {
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Route setup
	http.HandleFunc("/admin/register", handlers.AuthRegister)
	http.HandleFunc("/admin/login", handlers.AuthLogin)
	http.Handle("/admin/protected", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProtectedRoute)))
	http.HandleFunc("/admin/users", handlers.GetUsers)

	// Start the server
	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
