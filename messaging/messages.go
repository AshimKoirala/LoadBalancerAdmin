package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
)

type Message struct {
	Name string      `json:"name"`
	Body interface{} `json:"body"`
}

type ReplicaAdded struct {
	URL string `json:"url"`
}

type ReplicaStatisticsParameters struct {
	SuccessfulRequests int
	FailedRequests     int
}

type statDataArr struct {
	ReplicaName string                      `json:"replica_name"`
	Statistics  ReplicaStatisticsParameters `json:"statistics"`
}

type StatMessage struct {
	Name string        `json:"name"`
	Body []statDataArr `json:"body"`
}

type Messages struct {
	Name string      `json:"name"`
	Body interface{} `json:"body"`
}

func handleReplicaAdded(body string) {
	log.Print("Replica added: ", body)
	db.UpdateStatusByUrl(body, "active")
	ctx := context.Background()
	replica, err := db.GetReplicaByUrl(ctx, body)

	if err != nil {
		log.Printf("Failed to get replica by URL: %s", err)
		return
	}

	err = db.LogActivity(ctx, "success", fmt.Sprintf("Replica %v is now active", replica.Name), &replica.Id)

	if err != nil {
		log.Printf("Failed to log activity: %s", err)
		return
	}
	// log.Printf("Replica added: Name=%s, URL=%s", data.Name, data.URL)
}

func handleReplicaFailed(body string) {
	log.Print("Replica removed: ", body)
	db.UpdateStatusByUrl(body, "inactive")
	ctx := context.Background()
	replica, err := db.GetReplicaByUrl(ctx, body)

	if err != nil {
		log.Printf("Failed to get replica by URL: %s", err)
		return
	}

	err = db.LogActivity(ctx, "error", fmt.Sprintf("Replica %v is unavailable and set to inactive", replica.Name), &replica.Id)

	if err != nil {
		log.Printf("Failed to log activity: %s", err)
		return
	}
}

func handleReplicaRemoved(body string) {
	log.Print("Replica removed: ", body)
	db.UpdateStatusByUrl(body, "disabled")
	ctx := context.Background()
	replica, err := db.GetReplicaByUrl(ctx, body)

	if err != nil {
		log.Printf("Failed to get replica by URL: %s", err)
		return
	}

	err = db.LogActivity(ctx, "error", fmt.Sprintf("Replica %v is disabled", replica.Name), &replica.Id)

	if err != nil {
		log.Printf("Failed to log activity: %s", err)
		return
	}
}

func handleParametersUpdated(body []byte) {
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	updatedFields, ok := msg.Body.([]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-updated: %v", msg.Body)
		return
	}

	log.Printf("Parameters updated successfully: %v", updatedFields)
}

func handleParametersUpdateFailed(body []byte) {
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	errorMessage, ok := msg.Body.(string)
	if !ok {
		log.Printf("Invalid error message for parameters-update-failed: %v", msg.Body)
		return
	}

	log.Printf("Failed to update parameters: %s", errorMessage)
}

func handleStatistics(body []statDataArr) {
	var statisticsDatum []db.StatisticsData

	for _, replica := range body {
		data := db.StatisticsData{
			URL:                replica.ReplicaName,
			SuccessfulRequests: int64(replica.Statistics.SuccessfulRequests),
			FailedRequests:     int64(replica.Statistics.FailedRequests),
		}

		statisticsDatum = append(statisticsDatum, data)
	}

	err := db.BatchAddStatistics(&statisticsDatum)

	if err != nil {
		log.Printf("Failed to update statistics: %s", err)
		return
	}
}

func messageDemo() {
	// Example message publishing
	msg := Message{
		Name: "REMOVE_REPLICA",
		Body: "some message in byte",
	}

	PublishMessage(PUBLISHING_QUEUE, &msg)
}
