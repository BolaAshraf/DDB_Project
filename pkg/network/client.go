package network

import (
	"encoding/json"
	"net"
)

func SendTCPMessage(ip, port string, msg Message) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()
	return json.NewEncoder(conn).Encode(msg)
}
