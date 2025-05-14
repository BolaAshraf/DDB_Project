package replication

import (
	"distributed-mysql/pkg/network"
	"log"
)

// ReplicateQuery sends a query to all registered slaves
func ReplicateQuery(query string, slaves []string) {
	for _, slaveIP := range slaves {
		err := network.SendTCPMessage(slaveIP, "8080", network.Message{
			Type:  "REPLICATION",
			Query: query,
		})
		if err != nil {
			log.Printf("Failed to replicate to slave %s: %v\n", slaveIP, err)
		}
	}
}
