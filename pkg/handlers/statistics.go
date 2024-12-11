package handlers

import (
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
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

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := db.GetStatistics(r.Context())

	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to fetch statistics"})
		return
	}

	utils.NewSuccessResponse(w, stats)
}
