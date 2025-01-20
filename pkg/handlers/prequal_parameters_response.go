package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Printf("Error fetching latest prequal parameters response: %v", err)
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

	var payload db.AddPrequalParametersType
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Printf("Error decoding request payload: %v", err)
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	if payload.MaxLifeTime <= 0 || payload.ProbeRemoveFactor <= 0 || payload.PoolSize <= 10 || payload.Mu <= 0 || payload.ProbeFactor <= 0 {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid probe parameters. Ensure that max_life_time, probe_factor, probe_remove_factor and mu are greater than 0 and Pool size is greater than 10"})
		return
	}

	// Call the database function with the payload
	_, err := db.AddPrequalParametersResponse(r.Context(), payload)
	if err != nil {
		log.Println(err)
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to create or activate entry"})
		return
	}

	message := "Prequal Parameter added successfully"
	// if data.Id != nil {
	// 	message += fmt.Sprintf(" and entry with ID %d was activated", *payload.ActivateId)
	// }

	// rabbitmessage := &messaging.Message{
	// 	Name: messaging.NEW_PARAMETERS,
	// 	Body: map[string]interface{}{
	// 		"data":        payload,
	// 		"id": payload.ActivateId,
	// 	},
	// }

	// Publish the message to RabbitMQ
	// err = messaging.PublishMessage(messaging.PUBLISHING_QUEUE, rabbitmessage)
	// if err != nil {
	// 	log.Printf("Failed to publish change parameters message")
	// 	// utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to publish message to RabbitMQ"})
	// 	// return
	// }

	utils.NewSuccessResponse(w, message)
}
