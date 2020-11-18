package main

import (
	"net"
	"os"
	"strings"

	"./clientcommunication"
)

const TCP_PORT = "32001"

func main() {
	serverPath := os.Args[1]

	server, err := net.Listen("tcp", "localhost:"+TCP_PORT)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			panic(err)
		}
		defer connection.Close()

		packetTypeBuffer := make([]byte, 10)
		connection.Read(packetTypeBuffer)
		packetType := strings.Trim(string(packetTypeBuffer), ":")

		switch packetType {

		case "File":
			clientcommunication.ReceiveFile(connection, serverPath)

		case "Event":
			clientcommunication.ReceiveEvent(serverPath, connection)
		}
	}
}
