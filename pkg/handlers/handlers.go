package handlers

import (
	"log"
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/middleware"
	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
)

func Handler() {

	if err := db.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	// Routes setup with CORS
	mux := http.NewServeMux()

	// Route setup
	mux.HandleFunc("/admin/register", AuthRegister)
	mux.HandleFunc("/admin/login", AuthLogin)
	mux.Handle("/admin/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedRoute)))
	mux.HandleFunc("/admin/users", GetUsers)
	mux.HandleFunc("/admin/update", UpdateUser)

	// Wrap the entire mux with CORS
	handlerWithCORS := middleware.CORS(mux)

	// Start the server
	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
