package network

import (
	"encoding/json"
	"net"
)

type Message struct {
	Type     string `json:"type"`  // QUERY, REPLICATION
	Query    string `json:"query"` // SQL query
	SenderIP string `json:"slaveIP"`
}

func StartTCPServer(port string, handler func(Message)) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			var msg Message
			json.NewDecoder(c).Decode(&msg)
			handler(msg)
		}(conn)
	}
}
