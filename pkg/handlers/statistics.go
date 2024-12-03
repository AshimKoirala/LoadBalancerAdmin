package handlers

import (
	"encoding/json"
	"net/http"
)

type StatisticsResponse struct {
	Success bool `json:"success"`
	Total   struct {
		SuccessfulRequests int `json:"successful_requests"`
		FailedRequests     int `json:"failed_requests"`
	} `json:"total"`
	Data []ReplicaStatistics `json:"data"`
}

type ReplicaStatistics struct {
	Name               string `json:"name"`
	SuccessfulRequests int    `json:"successful_requests"`
	FailedRequests     int    `json:"failed_requests"`
}

func DummyStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Static data
	response := StatisticsResponse{
		Success: true,
	}
	response.Total.SuccessfulRequests = 1000
	response.Total.FailedRequests = 200
	response.Data = []ReplicaStatistics{
		{Name: "replica 1", SuccessfulRequests: 500, FailedRequests: 20},
		{Name: "replica 2", SuccessfulRequests: 500, FailedRequests: 180},
	}

	// Set response header and write JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
