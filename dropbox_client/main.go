package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func initialSync(rootPath string, scanPath string, dirEventTransmissionQueue chan<- dirEvent, fileTransmissionQueue chan<- fileInfo) {
	files, err := ioutil.ReadDir(scanPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			drEvnt := fsnotify.Event{Name: strings.SplitAfter(scanPath+"/"+file.Name(), rootPath)[1], Op: fsnotify.Create}
			dirEventTransmissionQueue <- dirEvent{Event: drEvnt, IsNewDir: true}
			initialSync(rootPath, scanPath+"/"+file.Name(), dirEventTransmissionQueue, fileTransmissionQueue)
		} else {
			fileTransmissionQueue <- fileInfo{fileName: strings.SplitAfter(scanPath+"/"+file.Name(), rootPath)[1], filePath: rootPath}
		}

	}
}

func main() {
	clientRoot := os.Args[1]

	fileTransmissionQueue := make(chan fileInfo, 10)
	go fileTransmission(fileTransmissionQueue)

	dirEventTransmissionQueue := make(chan dirEvent, 10)
	go eventTransmission(dirEventTransmissionQueue)

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(clientRoot, watchDir); err != nil {
		panic(err)
	}

	initialSync(clientRoot, clientRoot, dirEventTransmissionQueue, fileTransmissionQueue)

	for {
		select {
		case err := <-watcher.Errors:
			panic(err)

		case event := <-watcher.Events:
			fmt.Printf("client: EVENT! %#v\n", event)
			fileName := extractFileName(event.Name, clientRoot)

			switch event.Op {

			case fsnotify.Write:
				fileTransmissionQueue <- fileInfo{fileName: fileName, filePath: clientRoot}

			case fsnotify.Remove:
				dirEventTransmissionQueue <- dirEvent{Event: fsnotify.Event{Name: fileName, Op: event.Op}, IsNewDir: false}

			case fsnotify.Create:

				fInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Fatal(err)
				}

				if fInfo.IsDir() {
					watcher.Add(event.Name)
					dirEventTransmissionQueue <- dirEvent{Event: fsnotify.Event{Name: fileName, Op: event.Op}, IsNewDir: true}
				} else {
					fileTransmissionQueue <- fileInfo{fileName: fileName, filePath: clientRoot}
				}
			}
		}
	}
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func extractFileName(filepath string, clientRoot string) string {
	slc := strings.SplitAfter(filepath, clientRoot+"/")
	filename := slc[len(slc)-1]
	return filename
}
