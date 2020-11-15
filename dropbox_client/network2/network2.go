package network2

import (
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

	defer connection.Close()
	file, err := os.Open(filepath + "/" + filename)
	if err != nil {
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(filename, 64)

	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)

	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}

	return
}

func TransmitFile2(filepath string, filename string) {
	connection, err := net.Dial("tcp", "localhost:"+FILE_PORT)
	if err != nil {
		return
	}
	defer connection.Close()
	sendFileToServer(connection, filepath, filename)
}

func TransmitEvent(eventName string, op fsnotify.Op, isNewDir bool) {
	connection, err := net.Dial("tcp", "localhost:"+EVENT_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	sOp := strconv.FormatUint(uint64(op), 16)

	sIsNewDir := strconv.FormatBool(isNewDir)

	connection.Write([]byte(eventName + "-" + sOp + "-" + sIsNewDir))
}
