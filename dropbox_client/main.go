package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func initialSync(rootPath string, scanPath string) {
	files, err := ioutil.ReadDir(scanPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			network2.TransmitEvent(strings.SplitAfter(scanPath+"/"+file.Name(), rootPath)[1], fsnotify.Create, true)
			initialSync(rootPath, scanPath+"/"+file.Name())
		} else {
			network2.TransmitFile2(rootPath, strings.SplitAfter(scanPath+"/"+file.Name(), rootPath)[1])
		}

	}
}

func main() {
	clientRoot := os.Args[1]

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(clientRoot, watchDir); err != nil {
		panic(err)
	}

	initialSync(clientRoot, clientRoot)

	for {
		select {
		case err := <-watcher.Errors:
			panic(err)

		case event := <-watcher.Events:
			fmt.Printf("client: EVENT! %#v\n", event)
			fileName := extractFileName(event.Name, clientRoot)

			switch event.Op {

			case fsnotify.Write:
				network2.TransmitFile2(clientRoot, fileName)

			case fsnotify.Remove:
				network2.TransmitEvent(fileName, event.Op, false)

			case fsnotify.Create:
				fInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Fatal(err)
				}

				if fInfo.IsDir() {
					watcher.Add(event.Name)
					network2.TransmitEvent(fileName, event.Op, true)
				} else {
					network2.TransmitFile2(clientRoot, fileName)
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
