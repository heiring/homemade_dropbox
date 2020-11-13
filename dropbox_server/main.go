package main

import (
	"net"

	"./network2"
	"github.com/fsnotify/fsnotify"
)

type fileSystemChange struct {
	Event     fsnotify.Event
	FileLines []string
	IsDir     bool
}

func receiveFile(connectionEstablished <-chan net.Conn, dirPath string) {
	for {
		select {
		case connection := <-connectionEstablished:
			//fmt.Println("conn est")
			network2.CreateFileFromSocket(connection, dirPath)
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
