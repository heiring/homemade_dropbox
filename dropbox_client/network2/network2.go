package network2

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
)

const BUFFERSIZE = 1024

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func sendFileToServer(connection net.Conn, filepath string, filename string) {
	fmt.Println("A server has connected!")
	defer connection.Close()
	file, err := os.Open(filepath + "/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}

func sendEventToServer(connection net.Conn, event fsnotify.Event, isDir bool) {
	fmt.Println("A server has connected!")
	defer connection.Close()

	connection.Write([]byte(event.Name))
	connection.Write([]byte(strconv.FormatUint(uint64(event.Op), 16)))
	connection.Write([]byte(strconv.FormatBool(isDir)))

	return
}

func TransmitFile(filepath string, filename string) {
	server, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	transmissionComplete := false
	for !transmissionComplete {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		sendFileToServer(connection, filepath, filename)
		transmissionComplete = true
	}
}

func TransmitEvent(event fsnotify.Event, isDir bool) {
	server, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	transmissionComplete := false
	for !transmissionComplete {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		sendEventToServer(connection, event, isDir)
		transmissionComplete = true
	}
}
