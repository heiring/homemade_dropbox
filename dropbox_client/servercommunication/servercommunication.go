package servercommunication

import (
	"io"
	"net"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
)

const (
	BUFFERSIZE = 32768
	TCP_PORT   = "32001"
)

func fillString(returnString string, toLength int) string {
	for {
		lengtString := len(returnString)
		if lengtString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
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

func TransmitFile(filepath string, filename string) {
	connection, err := net.Dial("tcp", "localhost:"+TCP_PORT)
	if err != nil {
		return
	}
	defer connection.Close()

	wPacketType := fillString("File", 10)
	connection.Write([]byte(wPacketType))

	sendFileToServer(connection, filepath, filename)
}

func TransmitEvent(eventName string, op fsnotify.Op, isNewDir bool) {
	connection, err := net.Dial("tcp", "localhost:"+TCP_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	wPacketType := fillString("Event", 10)
	connection.Write([]byte(wPacketType))

	sOp := fillString(strconv.FormatUint(uint64(op), 16), 10)
	connection.Write([]byte(sOp))

	sIsNewDir := fillString(strconv.FormatBool(isNewDir), 10)
	connection.Write([]byte(sIsNewDir))

	wEventName := fillString(eventName, 256)
	connection.Write([]byte(wEventName))
}
