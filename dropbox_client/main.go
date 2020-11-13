package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"./fileoperations"
	"./network2"
	"github.com/fsnotify/fsnotify"
)

type fileSystemChange struct {
	Event     fsnotify.Event
	FileLines []string
	IsDir     bool
}

type fileInfo struct {
	fileName string
	filePath string
}

var watcher *fsnotify.Watcher

func fileTransmission(fileTransmissionQueue <-chan fileInfo) {
	ticker := time.NewTicker(5000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			select {
			case file := <-fileTransmissionQueue:
				fmt.Println("|||||||||||||||||||||||TRANSMITTING FILE||||||||||||||||||")
				network2.TransmitFile(file.filePath, file.fileName)
			}
		}
	}

}

func main() {
	//arg := os.Args[1]
	dirPath := "/tmp/dropbox/client"

	//send := make(chan fileSystemChange)
	//go bcast.Transmitter(32000, send)

	fileTransmissionQueue := make(chan fileInfo, 10)
	go fileTransmission(fileTransmissionQueue) //fileTransmission)

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(dirPath, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Printf("client: EVENT! %#v\n", event)
				fileName := fileoperations.ExtractFileName(event.Name)

				//transmitEvent := fsnotify.Event{Name: fileoperations.ExtractChangeName(event.Name, dirPath), Op: event.Op}
				fmt.Println(" ")
				switch event.Op {

				case fsnotify.Create:
					//fmt.Println("case: create")
					fInfo, err := os.Stat(event.Name)
					if err != nil {
						fmt.Println("ERROR finfo")
						log.Fatal(err)
					}
					//isDir := false
					if fInfo.IsDir() {
						watcher.Add(event.Name)
						//isDir = true
					} else {
						//isDir = false
						//network2.TransmitFile(dirPath, fileName)
						fileTransmissionQueue <- fileInfo{fileName: fileName, filePath: dirPath}
					}

					// } else {
					// 	send <- fileSystemChange{Event: transmitEvent, FileLines: nil, IsDir: false}
					// }
					//network2.TransmitEvent(transmitEvent, isDir)

				case fsnotify.Write:
					//fmt.Println("case: write")
					// fileLines := fileoperations.ReadFileLines(dirPath, filename)
					// send <- fileSystemChange{Event: transmitEvent, FileLines: fileLines, IsDir: false}
					//network2.TransmitFile(dirPath, fileName)

				case fsnotify.Remove:
					//send <- fileSystemChange{Event: transmitEvent, FileLines: nil, IsDir: false}
				}

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
