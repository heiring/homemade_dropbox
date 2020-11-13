package singlefile

import (
	"fmt"
	"os"

	"./network/bcast"
	"github.com/fsnotify/fsnotify"
)

// main
func singlefile() {
	arg := os.Args[1]

	send := make(chan fsnotify.Event)
	go bcast.Transmitter(32000, send)

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)
				send <- event

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(arg); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
