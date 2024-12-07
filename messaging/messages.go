package messaging

import "log"

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

func messageDemo() {
	// sending message to admin server when removing replica``
	msg := Message{
		Name: REMOVE_REPLICA,
		Body: "some message in byte",
	}

	PublishMessage(PUBLISHING_QUEUE, &msg)
}
