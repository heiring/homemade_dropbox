package network2

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	BUFFERSIZE = 1024
	FILE_PORT  = "27001"
)

func ListenForFiles(serverRoot string) {
	server, err := net.Listen("tcp", "localhost:"+FILE_PORT)
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		receiveFileFromClient(connection, serverRoot)
	}
}

func receiveFileFromClient(connection net.Conn, serverRoot string) {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)

	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(serverRoot + "/" + fileName)
	if err != nil {
		return
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}
