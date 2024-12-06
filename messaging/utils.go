package messaging

import (
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	Name string      `json:"name"`
	Body interface{} `json:"body"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func processMessage(body []byte) {
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	switch msg.Name {
	case "replica-added":
		handleReplicaAdded(msg.Body)
	case "replica-removed":
		handleReplicaRemoved(msg.Body)
	case "parameters-modified":
		fmt.Println("Parameters modified")
	case "statistics":
		fmt.Println("Save statistics")
	default:
		log.Printf("Unknown message type: %s", msg.Name)
	}
}
