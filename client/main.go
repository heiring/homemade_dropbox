package main

import (
	"os"
	"path/filepath"
	"strings"

	"./serversynch"
	"github.com/fsnotify/fsnotify"
)

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

	serversynch.InitialSynch(clientPath, clientPath)

	for {
		select {
		case err := <-watcher.Errors:
			panic(err)

		case event := <-watcher.Events:
			fileName := extractFileName(event.Name, clientPath)

			switch event.Op {

			case fsnotify.Write:
				serversynch.TransmitFile(clientPath, fileName)

			case fsnotify.Remove:
				serversynch.Remove(fileName)

			case fsnotify.Create:
				fInfo, err := os.Stat(event.Name)
				if err != nil {
					panic(err)
				}

				if fInfo.IsDir() {
					watcher.Add(event.Name)
					serversynch.NewDir(fileName)
				} else {
					serversynch.TransmitFile(clientPath, fileName)
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
