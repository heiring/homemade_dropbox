package main

import (
	"os"
	"path/filepath"
	"strings"

	"./serversync"
	"github.com/fsnotify/fsnotify"
)

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		// watch for any file system changes in this directory
		return watcher.Add(path)
	}

	return nil
}

var watcher *fsnotify.Watcher

func main() {
	clientPath := os.Args[1]

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// watch for any file system changes in the client directory and in possible subdirectories
	err := filepath.Walk(clientPath, watchDir)
	if err != nil {
		panic(err)
	}

	serversync.InitialSynch(clientPath, clientPath)

	for {
		select {
		case err := <-watcher.Errors:
			panic(err)

		case event := <-watcher.Events:
			fileName := extractFileName(event.Name, clientPath)

			switch event.Op {

			case fsnotify.Write:
				serversync.TransmitFile(clientPath, fileName)

			case fsnotify.Remove:
				serversync.Remove(fileName)

			case fsnotify.Create:
				fInfo, err := os.Stat(event.Name)
				if err != nil {
					panic(err)
				}

				if fInfo.IsDir() {
					// watch for any file system changes in the new directory
					watcher.Add(event.Name)

					serversync.NewDir(fileName)
				} else {
					serversync.TransmitFile(clientPath, fileName)
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
