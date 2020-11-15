package main

import (
	"net"
	"os"
	"strconv"
	"strings"

	"./fileoperations"
	"./network2"
)

const (
	FILE_PORT  = "27001"
	EVENT_PORT = "27002"
	REMOVE     = 4
)

func receiveFile(connectionEstablished <-chan net.Conn, serverRoot string) {
	for {
		connection := <-connectionEstablished
		if network2.ExtractRemotePort(connection) == FILE_PORT {
			network2.CreateFileFromSocket(connection, serverRoot)
		}
	}
}

func receiveEvent(serverRoot string) {
	server, err := net.Listen("tcp", "localhost:"+EVENT_PORT)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	defer server.Close()
	for {
		connection, err := server.Accept()
		if err != nil {
			panic(err)
			os.Exit(1)
		}
		defer connection.Close()
		buffer := make([]byte, 64)

		connection.Read(buffer)
		event := string(buffer)
		eventSplit := strings.SplitAfter(event, "-")

		eventName := strings.TrimSuffix(eventSplit[0], "-")

		sOp := strings.TrimSuffix(eventSplit[1], "-")
		op, err := strconv.Atoi(sOp)
		if err != nil {
			panic(err)
		}

		sIsNewDir := eventSplit[2]
		isNewDir, err := strconv.ParseBool(sIsNewDir[0:4])
		if err != nil {
			if sIsNewDir[0:4] != "fals" {
				panic(err)
			} else {
				isNewDir = false
			}
		}

		if isNewDir {
			fileoperations.CreateDir(serverRoot, eventName)
		}
		if op == REMOVE {
			fileoperations.Remove(serverRoot, eventName)
		}
	}
}

func main() {
	serverRoot := os.Args[1]

	done := make(chan bool)

	connectionEstablished := make(chan net.Conn)
	go network2.EstablishFileConnection(connectionEstablished)
	go receiveFile(connectionEstablished, serverRoot)
	go receiveEvent(serverRoot)
	<-done
}
