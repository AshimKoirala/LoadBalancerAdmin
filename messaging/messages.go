package messaging

import (
	"encoding/json"
	"fmt"
	"log"
)

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

func handleParametersApproved(body interface{}) {

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
		log.Printf("Invalid message body for parameters-approved: %v", msg.Body)
		return
	}

	id, _ := data["id"].(float64)
	log.Printf("Parameters approved for ID=%d", int64(id))
}

func handleParametersModified(body interface{}) {

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
		log.Printf("Invalid message body for parameters-modified: %v", msg.Body)
		return
	}

	modifiedFields, _ := data["fields"].([]interface{})
	log.Printf("Parameters modified: %v", modifiedFields)
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

func handleStatistics(body interface{}) {

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
		log.Printf("Invalid message body for statistics: %v", msg.Body)
		return
	}

	stats, _ := data["statistics"].(map[string]interface{})
	log.Printf("Statistics received: %v", stats)
}

func messageDemo() {
	// sending message to admin server when removing replica``
	msg := Message{
		Name: REMOVE_REPLICA,
		Body: "some message in byte",
	}

	PublishMessage(PUBLISHING_QUEUE, &msg)
}
