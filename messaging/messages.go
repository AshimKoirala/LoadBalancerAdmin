package messaging

import (
	"log"
)

func handleReplicaAdded(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for replica-added: %v", body)
		return
	}

	replicaName, _ := data["name"].(string)
	replicaURL, _ := data["url"].(string)

	log.Printf("Replica added: Name=%s, URL=%s", replicaName, replicaURL)

}

func handleReplicaRemoved(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for replica-removed: %v", body)
		return
	}

	replicaName, _ := data["name"].(string)
	replicaURL, _ := data["url"].(string)

	log.Printf("Replica removed: Name=%s, URL=%s", replicaName, replicaURL)

}

func handleParametersApproved(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-approved: %v", body)
		return
	}

	id, _ := data["id"].(float64)
	log.Printf("Parameters approved for ID=%d", int64(id))
}

func handleParametersModified(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-modified: %v", body)
		return
	}

	modifiedFields, _ := data["fields"].([]interface{})
	log.Printf("Parameters modified: %v", modifiedFields)
}

func handleParametersUpdated(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-updated: %v", body)
		return
	}

	updatedFields, _ := data["fields"].([]interface{})
	log.Printf("Parameters updated successfully: %v", updatedFields)
}

func handleParametersUpdateFailed(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for parameters-update-failed: %v", body)
		return
	}

	errorMessage, _ := data["error"].(string)
	log.Printf("Failed to update parameters: %s", errorMessage)
}

func handleStatistics(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for statistics: %v", body)
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
