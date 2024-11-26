package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/messaging"
	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

func AddReplica(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name                string `json:"name"`
		URL                 string `json:"url"`
		HealthcheckEndpoint string `json:"healthcheck_endpoint"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	// check health and add replica to db
	err = db.AddReplica(r.Context(), payload.Name, payload.URL, payload.HealthcheckEndpoint)
	if err != nil {
		// If there's an error i.e health check fail then return an appropriate error response
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	// Publish message to RabbitMQ
	message := &messaging.Message{
		Name: "ReplicaAdded",
		Body: map[string]string{
			"name": payload.Name,
			"url":  payload.URL,
		},
	}
	if err := messaging.PublishMessage("replica-events", message); err != nil {
		log.Printf("Failed to publish message: %v", err)
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to publish message"})
		return
	}

	// Return success message if replica is added successfully
	utils.NewSuccessResponse(w, "Replica added successfully")
}

func Status(w http.ResponseWriter, r *http.Request){
	log.Println("Replica status checking..")
}


func RemoveReplica(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID int64 `json:"id"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	// Remove replica from the database
	err = db.RemoveReplica(r.Context(), payload.ID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to remove replica"})
		return
	}

	utils.NewSuccessResponse(w, "Replica removed successfully")
}

func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID        int64  `json:"id"`
		NewStatus string `json:"new_status"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	// Change status of the replica
	err = db.ChangeStatus(r.Context(), payload.ID, payload.NewStatus)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to change replica status"})
		return
	}

	utils.NewSuccessResponse(w, "Replica status updated successfully")
}
