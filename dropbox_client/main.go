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

type dirEvent struct {
	Event    fsnotify.Event
	IsNewDir bool
}

type fileInfo struct {
	fileName string
	filePath string
}

var watcher *fsnotify.Watcher

func fileTransmission(fileTransmissionQueue <-chan fileInfo) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C
		file := <-fileTransmissionQueue
		network2.TransmitFile(file.filePath, file.fileName)
	}
}

func eventTransmission(dirEventTransmissionQueue <-chan dirEvent) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C
		dirEvent := <-dirEventTransmissionQueue
		network2.TransmitEvent(dirEvent.Event.Name, dirEvent.Event.Op, dirEvent.IsNewDir)
	}
}

func main() {
	//arg := os.Args[1]
	dirPath := "/tmp/dropbox/client"

	fileTransmissionQueue := make(chan fileInfo, 10)
	go fileTransmission(fileTransmissionQueue)

	dirEventTransmissionQueue := make(chan dirEvent, 10)
	go eventTransmission(dirEventTransmissionQueue)

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

				switch event.Op {

				case fsnotify.Create:

					fInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Fatal(err)
					}

					if fInfo.IsDir() {
						watcher.Add(event.Name)
						dirEventTransmissionQueue <- dirEvent{Event: fsnotify.Event{Name: fileName, Op: event.Op}, IsNewDir: true}
					} else {

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
					fileTransmissionQueue <- fileInfo{fileName: fileName, filePath: dirPath}

				case fsnotify.Remove:
					dirEventTransmissionQueue <- dirEvent{Event: fsnotify.Event{Name: fileName, Op: event.Op}, IsNewDir: false}
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
