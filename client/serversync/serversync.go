package serversync

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	BUFFERSIZE = 1400
	SYNC_PORT  = "32001"
)

func fillString(returnString string, toLength int) string {
	for {
		lengthString := len(returnString)
		if lengthString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}

func TransmitFile(filepath string, filename string) {
	connection, err := net.Dial("tcp", "localhost:"+SYNC_PORT)
	if err != nil {
		return
	}
	defer connection.Close()

	wPacketType := fillString("File", 10)
	connection.Write([]byte(wPacketType))

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

func NewDir(dirPath string) {
	connection, err := net.Dial("tcp", "localhost:"+SYNC_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	wPacketType := fillString("New Dir", 10)
	connection.Write([]byte(wPacketType))

	wDirPath := fillString(dirPath, 256)
	connection.Write([]byte(wDirPath))
}

func Remove(filePath string) {
	connection, err := net.Dial("tcp", "localhost:"+SYNC_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	wPacketType := fillString("Remove", 10)
	connection.Write([]byte(wPacketType))

	wFilePath := fillString(filePath, 256)
	connection.Write([]byte(wFilePath))
}

func InitialSynch(clientPath string, dirReadPath string) {
	files, err := ioutil.ReadDir(dirReadPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		pathFromClientDir := strings.SplitAfter(dirReadPath+"/"+file.Name(), clientPath)[1]

		if file.IsDir() {
			NewDir(pathFromClientDir)
			InitialSynch(clientPath, dirReadPath+"/"+file.Name())

		} else {
			TransmitFile(clientPath, pathFromClientDir)
		}

	}
}
