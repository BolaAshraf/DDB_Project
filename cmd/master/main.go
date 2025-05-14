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

var (
	masterIP = "192.168.0.107"           // Master's own IP
	slaves   = []string{"192.168.0.106"} // Known slave IPs
)

func main() {
	// Initialize local MySQL
	db, err := db.New("127.0.0.1", 3306, "root", "Leo10Messi#229", "")
	if err != nil {
		log.Fatal(err)
	}

	// Start TCP server with proper message handling
	go network.StartTCPServer("8080", func(msg network.Message) {
		switch msg.Type {
		case "QUERY_FROM_SLAVE":
			log.Printf("Received query from slave %s: %s", msg.SenderIP, msg.Query)

			// 1. Execute the query locally on master
			if err := db.ExecQuery(msg.Query); err != nil {
				log.Printf("Failed to execute slave query: %v", err)
				return
			}

			// 2. Replicate to other slaves (excluding the sender)
			for _, slave := range slaves {
				if slave != msg.SenderIP {
					log.Printf("Replicating to slave %s", slave)
					err := network.SendTCPMessage(slave, "8080", network.Message{
						Type:     "REPLICATION",
						Query:    msg.Query,
						SenderIP: masterIP,
					})
					if err != nil {
						log.Printf("Failed to replicate to slave %s: %v", slave, err)
					}
				}
			}

		case "QUERY": // Local queries from master's console
			log.Printf("Processing local query: %s", msg.Query)
			if shouldReplicate(msg.Query) {
				for _, slave := range slaves {
					log.Printf("Replicating to slave %s", slave)
					err := network.SendTCPMessage(slave, "8080", network.Message{
						Type:     "REPLICATION",
						Query:    msg.Query,
						SenderIP: masterIP,
					})
					if err != nil {
						log.Printf("Failed to replicate to slave %s: %v", slave, err)
					}
				}
			}
		}
	})

	// CLI for master console input
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Master> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if selectAll(query, db) {
			continue
		}
		// Execute locally first
		if err := db.ExecQuery(query); err != nil {
			log.Println("Local execution error:", err)
			continue
		}

		// Send to TCP handler for replication
		network.SendTCPMessage("127.0.0.1", "8080", network.Message{
			Type:  "QUERY",
			Query: query,
		})
	}
}

func shouldReplicate(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	ignoredPrefixes := []string{"SELECT", "SHOW", "SET"}
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(query, prefix) {
			return false
		}
	}
	return true
}

func selectAll(query string, db *db.Database) bool {
	query = strings.TrimSpace(strings.ToUpper(query))
	if !strings.HasPrefix(query, "SELECT") {
		return false
	}
	// فتح قاعدة البيانات المحلية

	//defer db.Close()

	// تنفيذ الاستعلام
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("❌ Error executing query:", err)
	}
	defer rows.Close()

	// الحصول على أسماء الأعمدة
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("❌ Error getting columns:", err)
	}

	// طباعة أسماء الأعمدة
	fmt.Println("📃 Query Result:")
	fmt.Println(strings.Join(columns, " | "))

	// معالجة الصفوف وطباعة البيانات
	for rows.Next() {
		// تجهيز مصفوفة لحفظ القيم
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// قراءة الصف
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal("❌ Error scanning row:", err)
		}

		// طباعة الصف
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
