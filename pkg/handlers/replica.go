package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/AshimKoirala/load-balancer-admin/messaging"
	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

func AddReplica(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}
	var payload struct {
		Name                string `json:"name"`
		URL                 string `json:"url"`
		HealthcheckEndpoint string `json:"healthcheck_endpoint"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)

	switch {
	case payload.Name == "" || payload.URL == "" || payload.HealthcheckEndpoint == "":
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"all fields (name, url, healthCheckEndpoint) must be provided"})
		return
		// case !strings.HasPrefix(payload.URL, os.Getenv("REPLICA_URL")):
		// 	return fmt.Errorf("malicious URL")
	}

	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}
	url, err := url.Parse(payload.URL)

	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Malformed url"})
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s://%s/%s", url.Scheme, url.Host, payload.HealthcheckEndpoint))

	if resp.StatusCode != http.StatusOK {
		// log.Printf("received non-200 response: %d", resp.StatusCode)

		if err != nil {
			utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Replica did not pass the healthcheck"})
			return
		}
	}
	defer resp.Body.Close()

	err = db.AddReplica(r.Context(), payload.Name, payload.URL, payload.HealthcheckEndpoint)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	// Publish message to RabbitMQ
	message := &messaging.Message{
		Name: "replica-added",
		Body: map[string]string{
			"name": payload.Name,
			"url":  payload.URL,
		},
	}

	if err := messaging.PublishMessage(utils.PUBLISHING_QUEUE, message); err != nil {
		log.Printf("Failed to publish message: %v", err)
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to publish message"})
		return
	}

	utils.NewSuccessResponse(w, "Replica added successfully")
}

func Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Replica status checking..")
}

func GetReplicas(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}
	replicas, err := db.GetReplicas(r.Context())
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to fetch replicas"})
		return
	}

	utils.NewSuccessResponse(w, replicas)
}

func RemoveReplica(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}
	var payload struct {
		ID int64 `json:"id"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	replica, err := db.GetReplicaByID(payload.ID)

	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, []string{"Could not find replica"})
		return
	}

	// Remove replica from the database
	err = db.RemoveReplica(r.Context(), payload.ID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to remove replica"})
		return
	}

	message := &messaging.Message{
		Name: "replica-removed",
		Body: map[string]string{
			"name": replica.Name,
			"url":  replica.URL,
		},
	}

	if err := messaging.PublishMessage(utils.PUBLISHING_QUEUE, message); err != nil {
		log.Printf("Failed to publish message: %v", err)
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to publish message"})
		return
	}

	utils.NewSuccessResponse(w, "Replica removed successfully")
}

func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
	    utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}
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
