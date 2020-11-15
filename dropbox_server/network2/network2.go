package network2

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func CreateFileFromSocket(connection net.Conn, filepath string) {
	// connection, err := net.Dial("tcp", "localhost:27001")
	// if err != nil {
	// 	panic(err)
	// }
	defer connection.Close()

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)

	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(filepath + "/" + fileName)
	if err != nil {
		//panic(err)
		return //XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXx
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
func ReceiveEvent(connection net.Conn) (string, uint32, bool) {
	// connection, err := net.Dial("tcp", "localhost:27001")
	// if err != nil {
	// 	panic(err)
	// }

	defer connection.Close()

	bufferEventName := make([]byte, 64)
	connection.Read(bufferEventName)
	eventName := string(bufferEventName)

	bufferOp := make([]byte, 64)
	connection.Read(bufferOp)
	op, err := strconv.ParseUint(string(bufferOp), 16, 32)
	if err != nil {
		panic(err)
	}

	bufferIsDir := make([]byte, 64)
	connection.Read(bufferIsDir)
	isDir, err := strconv.ParseBool(string(bufferIsDir))
	if err != nil {
		panic(err)
	}

	return eventName, uint32(op), isDir
}

func EstablishFileConnection(connectionEstablished chan<- net.Conn) {
	for {
		connection, err := net.Dial("tcp", "localhost:27001")
		if err != nil {
			//panic(err)
		} else {
			connectionEstablished <- connection
		}
	}
}

func ExtractRemotePort(conn net.Conn) string {
	slc := strings.SplitAfter(conn.RemoteAddr().String(), ":")
	return slc[len(slc)-1]
}
