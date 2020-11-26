package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func isDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return f.IsDir()
}

func WatchFolder(path string) (chan fsnotify.Event, func() error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		log.Fatal("walk", err)
	}

	changes := make(chan fsnotify.Event, 10)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Remove:
					if isDir(event.Name) {
						watcher.Remove(event.Name)
					}
				case fsnotify.Create:
					if isDir(event.Name) {
						watcher.Add(event.Name)
					}
				}
				changes <- event
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	return changes, watcher.Close
}
