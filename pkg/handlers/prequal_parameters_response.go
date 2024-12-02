package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

// to fetch latest PrequalParametersResponse
func GetPrequalParameters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	response, err := db.GetPrequalParametersResponse(r.Context())
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to fetch the latest entry"})
		return
	}

	utils.NewSuccessResponse(w, response)
}

// to create a new PrequalParametersResponse
func AddPrequalParameters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	var payload db.PrequalParametersResponse
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	err := db.AddPrequalParametersResponse(r.Context(), payload)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to create new entry"})
		return
	}

	utils.NewSuccessResponse(w, "Entry created successfully")
}
