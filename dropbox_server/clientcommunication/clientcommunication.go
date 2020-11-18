package clientcommunication

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"../fileoperations"
)

const (
	BUFFERSIZE = 32768
	REMOVE     = 4
)

func ReceiveFile(connection net.Conn, serverPath string) {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)

	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(serverPath + "/" + fileName)
	if err != nil {
		return
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
}

func ReceiveEvent(serverPath string, connection net.Conn) {
	opBuffer := make([]byte, 10)
	connection.Read(opBuffer)
	op, err := strconv.Atoi(strings.Trim(string(opBuffer), ":"))
	if err != nil {
		panic(err)
	}

	isNewDirBuffer := make([]byte, 10)
	connection.Read(isNewDirBuffer)
	isNewDir, err := strconv.ParseBool(strings.Trim(string(isNewDirBuffer), ":"))
	if err != nil {
		panic(err)
	}

	eventNameBuffer := make([]byte, 256)
	connection.Read(eventNameBuffer)
	eventName := strings.Trim(string(eventNameBuffer), ":")

	if isNewDir {
		fileoperations.CreateDir(serverPath + "/" + eventName)
	}
	if op == REMOVE {
		fileoperations.Remove(serverPath + "/" + eventName)
	}
}
