package main

import (
	"net"
	"os"
	"strconv"
	"strings"

	"./fileoperations"
	"./network2"
	"github.com/fsnotify/fsnotify"
)

type fileSystemChange struct {
	Event     fsnotify.Event
	FileLines []string
	IsDir     bool
}

const (
	FILE_PORT  = "27001"
	EVENT_PORT = "27002"
)

func receiveFile(connectionEstablished <-chan net.Conn, dirPath string) {
	for {
		// select {
		// case connection := <-connectionEstablished:
		// 	//fmt.Println("conn est")
		// 	network2.CreateFileFromSocket(connection, dirPath)
		// }
		connection := <-connectionEstablished
		if network2.ExtractRemotePort(connection) == FILE_PORT {
			network2.CreateFileFromSocket(connection, dirPath)
		}
	}
}

func receiveEvent(dirPath string) {
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
		eventSplit := strings.SplitAfter(event, "/")

		eventName := strings.TrimSuffix(eventSplit[0], "/")

		sOp := strings.TrimSuffix(eventSplit[1], "/")
		op, err := strconv.Atoi(sOp)
		if err != nil {
			panic(err)
		}

		sIsNewDir := eventSplit[2]
		isNewDir, err := strconv.ParseBool(sIsNewDir[0:4])
		if err != nil {
			if sIsNewDir[0:4] != "fals" {
				panic(err)
			}
		}

		if isNewDir {
			fileoperations.CreateDir(dirPath, eventName)
		}
		if op == 4 {
			fileoperations.Remove(dirPath, eventName)
		}
	}
}

// main
func main() {
	//dirPath := os.Args[1]
	dirPath := "/tmp/dropbox/server"

	done := make(chan bool)

	connectionEstablished := make(chan net.Conn)
	go network2.EstablishFileConnection(connectionEstablished)
	go receiveFile(connectionEstablished, dirPath)
	go receiveEvent(dirPath)

	// go func() {
	// 	//connection <- connectionEstabllished
	// 	for {
	// 		select {
	// 		case change := <-receive:
	// 			fmt.Print("server event: ")
	// 			fmt.Println(change)
	// 			//filename := etc.ExtractFileName(change.Event.Name)
	// 			//changeName := etc.ExtractChangeName(change.Event.Name, dirPath)
	// 			fmt.Println(" ")
	// 			switch change.Event.Op {

	// 			case fsnotify.Create:
	// 				//fmt.Println("create event")
	// 				if change.IsDir {
	// 					fileoperations.CreateDir(dirPath, change.Event.Name)
	// 				} else {
	// 					fileoperations.CreateFile(dirPath + "/" + change.Event.Name)
	// 				}

	// 			case fsnotify.Write:
	// 				//fmt.Println("case: write")
	// 				//fmt.Println(change.FileLines)
	// 				fileoperations.WriteSliceToFile(dirPath, change.Event.Name, change.FileLines)

	// 			case fsnotify.Remove:
	// 				fileoperations.Remove(dirPath, change.Event.Name)

	// 			case fsnotify.Rename:

	// 			case fsnotify.Chmod:
	// 				fmt.Println("Chmod")
	// 			default:
	// 				fmt.Println("default")
	// 			}

	// 		}
	// 	}
	// }()

	<-done
}
