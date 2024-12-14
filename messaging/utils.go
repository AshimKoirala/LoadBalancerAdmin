package messaging

import (
	"encoding/json"
	"log"
)

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
		// return
		handleReplicaAdded(body)
	case REMOVED_REPLICA:
		handleReplicaRemoved(body)
	case PARAMETERS_UPDATED:
		handleParametersUpdated(body)
	case PARAMETERS_UPDATE_FAILED:
		handleParametersUpdateFailed(body)
	case STATISTICS:
		var stmsg StatMessage
		if err := json.Unmarshal(body, &stmsg); err != nil {
			log.Printf("Failed to unmarshal stat message: %v", err)
			return
		}
		handleStatistics(stmsg.Body)
	default:
		log.Printf("Unknown message type: %s", msg.Name)
	}
}
