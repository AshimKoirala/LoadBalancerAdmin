package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
)

type Message struct {
	Name string      `json:"name"`
	Body interface{} `json:"body"`
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
	fmt.Print(body)
	// bodyBytes, ok := body.([]byte)
	// if !ok {
	// 	log.Printf("Invalid body type: expected []byte, got %T", body)
	// 	return
	// }

	// var msg Messages
	// if err := json.Unmarshal(bodyBytes, &msg); err != nil {
	// 	log.Printf("Failed to unmarshal message: %v", err)
	// 	return
	// }

	// data, ok := msg.Body.(map[string]interface{})
	// if !ok {
	// 	log.Printf("Invalid message body for replica-added: %v", msg.Body)
	// 	return
	// }

	// replicaName, _ := data["name"].(string)
	// replicaURL, _ := data["url"].(string)

	// log.Printf("Replica added: Name=%s, URL=%s", replicaName, replicaURL)
}

func handleReplicaRemoved(body interface{}) {

	bodyBytes, ok := body.([]byte)
	if !ok {
		log.Printf("Invalid body type: expected []byte, got %T", body)
		return
	}

	var msg Messages
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	data, ok := msg.Body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for replica-removed: %v", msg.Body)
		return
	}

	replicaName, _ := data["name"].(string)
	replicaURL, _ := data["url"].(string)

	log.Printf("Replica removed: Name=%s, URL=%s", replicaName, replicaURL)
}

func handleParametersUpdated(body interface{}) {

	bodyBytes, ok := body.([]byte)
	if !ok {
		log.Printf("Invalid body type: expected []byte, got %T", body)
		return
	}

	var msg Message
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	data, ok := msg.Body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-updated: %v", msg.Body)
		return
	}

	updatedFields, _ := data["fields"].([]interface{})
	log.Printf("Parameters updated successfully: %v", updatedFields)
}

func handleParametersUpdateFailed(body interface{}) {

	bodyBytes, ok := body.([]byte)
	if !ok {
		log.Printf("Invalid body type: expected []byte, got %T", body)
		return
	}

	var msg Messages
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	data, ok := msg.Body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-update-failed: %v", msg.Body)
		return
	}

	errorMessage, _ := data["error"].(string)
	log.Printf("Failed to update parameters: %s", errorMessage)
}

func handleStatistics(body []statDataArr) {
	// decodedBytes, err := base64.StdEncoding.DecodeString(body)
	// if err != nil {
	// 	log.Fatalf("Failed to decode Base64: %v", err)
	// }

	// type ReplicaStatistics struct {
	// 	SuccessfulRequests int `json:"successful_requests"`
	// 	FailedRequests     int `json:"failed_requests"`
	// }

	// type ReplicaStatisticsData struct {
	// 	ReplicaName string            `json:"replica_name"`
	// 	Statistics  ReplicaStatistics `json:"statistics"`
	// }

	// var replicaStatisticsMessages []ReplicaStatisticsData

	// err = json.Unmarshal(decodedBytes, &replicaStatisticsMessages)
	// if err != nil {
	// 	fmt.Printf("Failed to unmarshal JSON: %v\n", err)
	// 	return
	// }

	var statisticsDatum []db.StatisticsData

	for _, replica := range body {
		data := db.StatisticsData{
			URL:                replica.ReplicaName,
			SuccessfulRequests: int64(replica.Statistics.SuccessfulRequests),
			FailedRequests:     int64(replica.Statistics.FailedRequests),
		}

		statisticsDatum = append(statisticsDatum, data)
		fmt.Printf("Replica: %s, SuccessfulRequests: %d, FailedRequests: %d\n",
			replica.ReplicaName, replica.Statistics.SuccessfulRequests, replica.Statistics.FailedRequests)
	}

	err := db.BatchAddStatistics(&statisticsDatum)

	if err != nil {
		log.Printf("Failed to update statistics: %s", err)
		return
	}
	// bodyBytes, ok := body.([]byte)
	// if !ok {
	// 	log.Printf("Invalid body type: expected []byte, got %T", body)
	// 	return
	// }

	// var msg Messages
	// if err := json.Unmarshal(bodyBytes, &msg); err != nil {
	// 	log.Printf("Failed to unmarshal message: %v", err)
	// 	return
	// }

	// data, ok := msg.Body.(map[string]interface{})
	// if !ok {
	// 	log.Printf("Invalid message body for statistics: %v", msg.Body)
	// 	return
	// }

	// stats, _ := data["statistics"].(map[string]interface{})
	// log.Printf("Statistics received: %v", stats)
}

func messageDemo() {
	// sending message to admin server when removing replica``
	msg := Message{
		Name: REMOVE_REPLICA,
		Body: "some message in byte",
	}

	PublishMessage(PUBLISHING_QUEUE, &msg)
}
