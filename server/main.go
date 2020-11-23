package main

import (
	"net"
	"os"
	"strings"

	"./clientsync"
)

const SYNC_PORT = "32001"

func main() {
	serverPath := os.Args[1]

	server, err := net.Listen("tcp", "localhost:"+SYNC_PORT)
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
			clientsync.ReceiveFile(connection, serverPath)

		case "New Dir":
			clientsync.CreateDir(serverPath, connection)

		case "Remove":
			clientsync.RemoveFile(serverPath, connection)
		}
	}
}
