package handlers

import (
	"log"
	"net/http"
	"os"

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
	mux.HandleFunc("POST /admin/register", AuthRegister)
	mux.HandleFunc("POST /admin/login", AuthLogin)
	mux.Handle("GET /admin/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedRoute)))
	mux.Handle("GET /admin/users", middleware.AuthMiddleware(http.HandlerFunc(GetUsers)))
	mux.Handle("PATCH /admin/update/{id}", middleware.AuthMiddleware(http.HandlerFunc(UpdateUser)))
	mux.HandleFunc("/admin/forgot-password", ForgotPassword)
	mux.HandleFunc("/admin/reset-password", ResetPassword)
	mux.Handle("POST /admin/add-replica", middleware.AuthMiddleware(http.HandlerFunc(AddReplica)))
	mux.HandleFunc("GET /admin/get-replica", GetReplicas)
	mux.Handle("DELETE /admin/remove-replica", middleware.AuthMiddleware(http.HandlerFunc(RemoveReplica)))
	mux.Handle("PATCH /admin/change-status", middleware.AuthMiddleware(http.HandlerFunc(ChangeStatus)))
	mux.Handle("GET /admin/activity-logs", middleware.AuthMiddleware(http.HandlerFunc(GetActivityLogs)))
	mux.Handle("POST /admin/update-prequal-parameters", middleware.AuthMiddleware(http.HandlerFunc(AddPrequalParameters)))
	mux.Handle("GET /admin/get-prequal-parameters", middleware.AuthMiddleware(http.HandlerFunc(GetPrequalParameters)))
	mux.Handle("GET /admin/get-statistics", middleware.AuthMiddleware(http.HandlerFunc(GetStatistics)))

	// Wrap the entire mux with CORS
	// handlerWithCORS := middleware.CORS(mux)

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on : %s", port)
	if err := http.ListenAndServe(":"+port, middleware.CORS(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
