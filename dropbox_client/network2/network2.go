package network2

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
)

const (
	BUFFERSIZE = 1024
	FILE_PORT  = "27001"
	EVENT_PORT = "27002"
)

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

func TransmitEvent(eventName string, op fsnotify.Op, isNewDir bool) {
	connection, err := net.Dial("tcp", "localhost:"+EVENT_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	//connection.Write([]byte(eventName))

	sOp := strconv.FormatUint(uint64(op), 16)
	//fmt.Println("transmit op: " + sOp)
	//connection.Write([]byte(sOp))

	sIsNewDir := strconv.FormatBool(isNewDir)
	//fmt.Println("transmit isnewdie: " + sIsNewDir)
	connection.Write([]byte(eventName + "/" + sOp + "/" + sIsNewDir))
}
