package messaging

import (
	"encoding/json"
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
	case ADDED_REPLICA:
		return
		// handleReplicaAdded(msg.Body)
	case REMOVED_REPLICA:
		handleReplicaRemoved(msg.Body)
	case PARAMETERS_UPDATED:
		handleParametersUpdated(msg.Body)
	case PARAMETERS_UPDATE_FAILED:
		handleParametersUpdateFailed(msg.Body)
	case STATISTICS:
		handleStatistics(msg.Body.(string))
	default:
		log.Printf("Unknown message type: %s", msg.Name)
	}
}
