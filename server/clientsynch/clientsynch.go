package clientsynch

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"../fileoperations"
)

const BUFFERSIZE = 32768

func ReceiveFile(connection net.Conn, serverPath string) {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, err := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	if err != nil {
		panic(err)
	}

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

func CreateDir(serverPath string, connection net.Conn) {

	dirNameBuffer := make([]byte, 256)
	connection.Read(dirNameBuffer)
	dirName := strings.Trim(string(dirNameBuffer), ":")

	fileoperations.CreateDir(serverPath + "/" + dirName)
}

func RemoveFile(serverPath string, connection net.Conn) {
	fileNameBuffer := make([]byte, 256)
	connection.Read(fileNameBuffer)
	fileName := strings.Trim(string(fileNameBuffer), ":")

	fileoperations.Remove(serverPath + "/" + fileName)
}