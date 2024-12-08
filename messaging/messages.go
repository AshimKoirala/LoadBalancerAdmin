package messaging

import (
	"context"
	"log"

	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
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
		log.Printf("Invalid message body for prequal-approval: %v", body)
		return
	}

	//extract id
	id, ok := data["id"].(float64)
	if !ok {
		log.Printf("Missing or invalid ID in prequal-approval message: %v", body)
		return
	}
	var response db.PrequalParametersResponse
	var someInt int64 = 0
	err := db.AddPrequalParametersResponse(context.Background(), response, &someInt)
	if err != nil {
		log.Printf("Failed to activate PrequalParametersResponse with ID %d: %v", int(id), err)
	} else {
		log.Printf("Successfully activated PrequalParametersResponse with ID %d", int(id))
	}
}

func messageDemo() {
	// sending message to admin server when removing replica``
	msg := Message{
		Name: REMOVE_REPLICA,
		Body: "some message in byte",
	}

	PublishMessage(PUBLISHING_QUEUE, &msg)
}
