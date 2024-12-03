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
	mux.HandleFunc("/admin/register", AuthRegister)
	mux.HandleFunc("/admin/login", AuthLogin)
	mux.Handle("/admin/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedRoute)))
	mux.HandleFunc("/admin/users", GetUsers)
	mux.HandleFunc("/admin/update", UpdateUser)
	mux.HandleFunc("/admin/forgot-password", ForgotPassword)
	mux.HandleFunc("/admin/reset-password", ResetPassword)
	mux.HandleFunc("/admin/add-replica", AddReplica)
	mux.HandleFunc("/admin/get-replica", GetReplicas)
	mux.HandleFunc("/admin/remove-replica", RemoveReplica)
	mux.HandleFunc("/admin/change-status", ChangeStatus)
	mux.HandleFunc("/admin/status", Status)
	mux.HandleFunc("/admin/activity-logs", GetActivityLogs)
	mux.HandleFunc("/admin/update-prequal-parameters", AddPrequalParameters)
	mux.HandleFunc("/admin/get-prequal-parameters", GetPrequalParameters)
	mux.HandleFunc("/admin/get-statistics", DummyStatistics)
	
	// Wrap the entire mux with CORS
	handlerWithCORS := middleware.CORS(mux)

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on : %s", port)
	if err := http.ListenAndServe(":"+port, handlerWithCORS); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
