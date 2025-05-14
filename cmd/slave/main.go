package main

import (
	"bufio"
	"distributed-mysql/pkg/db"
	"distributed-mysql/pkg/network"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	masterIP = "192.168.0.108" // Master's IP address
	slaveIP  = "192.168.0.107" // THIS slave's own IP (CHANGE THIS FOR EACH SLAVE!)
)

func main() {
	// Initialize local MySQL
	db, err := db.New("127.0.0.1", 3306, "root", "rootroot", "")
	if err != nil {
		log.Fatal(err)
	}

	// Start TCP server to listen for replication
	go network.StartTCPServer("8080", func(msg network.Message) {
		switch msg.Type {
		case "REPLICATION":
			log.Printf("Executing replicated query from %s: %s\n", msg.SenderIP, msg.Query)
			if err := db.ExecQuery(msg.Query); err != nil {
				log.Printf("Replication failed: %v\n", err)
			}
		default:
			log.Printf("Received unexpected message type: %s\n", msg.Type)
		}
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Slave> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if isMasterQuery(query) {
			fmt.Println("🚫 Access denied: Only the master can execute CREATE or DROP operations.")
			continue
		}
		if selectAll(query, db) {
			continue
		}
		// 1. Execute locally first
		if err := db.ExecQuery(query); err != nil {
			log.Println("Local execution error:", err)
			continue
		}

		// 2. Send to master for replication (only for DML/DDL queries)
		if shouldReplicate(query) {
			log.Printf("Sending query to master: %s\n", query)
			err := network.SendTCPMessage(masterIP, "8080", network.Message{
				Type:     "QUERY_FROM_SLAVE", // Changed from "QUERY" to "QUERY_FROM_SLAVE"
				Query:    query,
				SenderIP: slaveIP, // Identify which slave sent this
			})
			if err != nil {
				log.Printf("Failed to send to master: %v\n", err)
			}
		}
	}
}

func shouldReplicate(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	ignoredPrefixes := []string{"SELECT", "SHOW", "USE", "SET"}
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(query, prefix) {
			return false
		}
	}
	return true
}
func isMasterQuery(query string) bool {
	query = strings.TrimSpace(strings.ToUpper(query))
	return strings.HasPrefix(query, "CREATE") || strings.HasPrefix(query, "DROP")
}
func selectAll(query string, db *db.Database) bool {
	query = strings.TrimSpace(strings.ToUpper(query))
	if !strings.HasPrefix(query, "SELECT") {
		return false
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("❌ Error executing query:", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("❌ Error getting columns:", err)
	}

	fmt.Println("📃 Query Result:")
	fmt.Println(strings.Join(columns, " | "))

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal("❌ Error scanning row:", err)
		}
		for _, val := range values {
			if b, ok := val.([]byte); ok {
				fmt.Printf("%s\t", string(b))
			} else {
				fmt.Printf("%v\t", val)
			}
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Fatal("❌ Row error:", err)
	}
	return true
}
