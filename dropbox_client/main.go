package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"./servercommunication"
	"github.com/fsnotify/fsnotify"
)

func initialSync(clientPath string, dirReadPath string) {
	files, err := ioutil.ReadDir(dirReadPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		pathFromClientDir := strings.SplitAfter(dirReadPath+"/"+file.Name(), clientPath)[1]

		if file.IsDir() {
			servercommunication.TransmitEvent(pathFromClientDir, fsnotify.Create, true)
			initialSync(clientPath, dirReadPath+"/"+file.Name())

		} else {
			servercommunication.TransmitFile(clientPath, pathFromClientDir)
		}

	}
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		return watcher.Add(path)
	}

	return nil
}

var watcher *fsnotify.Watcher

func main() {
	clientPath := os.Args[1]

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(clientPath, watchDir)
	if err != nil {
		panic(err)
	}

	initialSync(clientPath, clientPath)

	for {
		select {
		case err := <-watcher.Errors:
			panic(err)

		case event := <-watcher.Events:
			fileName := extractFileName(event.Name, clientPath)

			switch event.Op {

			case fsnotify.Write:
				servercommunication.TransmitFile(clientPath, fileName)

			case fsnotify.Remove:
				servercommunication.TransmitEvent(fileName, event.Op, false)

			case fsnotify.Create:
				fInfo, err := os.Stat(event.Name)
				if err != nil {
					panic(err)
				}

				if fInfo.IsDir() {
					watcher.Add(event.Name)
					servercommunication.TransmitEvent(fileName, event.Op, true)
				} else {
					servercommunication.TransmitFile(clientPath, fileName)
				}
			}
		}
	}
}

func extractFileName(filepath string, clientPath string) string {
	slc := strings.SplitAfter(filepath, clientPath+"/")
	filename := slc[len(slc)-1]
	return filename
}
