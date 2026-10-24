package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var rockRootConfigWatcher *fsnotify.Watcher
var rockConfigWatcher *fsnotify.Watcher

func watchForFileChanges() {
	var err error
	rockRootConfigWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not create watcher: %s", err.Error())
	}
	rockConfigWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not create watcher: %s", err.Error())
	}

	err = watchAllRocksInFolder(RhyoliteDir)
	if err != nil {
		log.Fatalf("Failed to watch rocks: %s", err.Error())
	}

	go rockRootWatcherHandle()
	go rockWatcherHandle()
}

func rockRootWatcherHandle() {
	for {
		select {
		case event, ok := <-rockRootConfigWatcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				if err := watchRock(event.Name); err != nil {
					log.Println("Error watching rock: ", err.Error())
				}
			}

			if event.Has(fsnotify.Remove) {
				if err := refreshPath(event.Name); err != nil {
					log.Println("Error on rock removal: ", err.Error())
				}
			}
		case err, ok := <-rockRootConfigWatcher.Errors:
			if !ok {
				return
			}
			log.Println("Rock Root Watcher error:", err.Error())
		}
	}
}

func watchAllRocksInFolder(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("Cannot read directory:", directory)
	}
	errors := ""
	for _, entry := range entries {
		if entry.IsDir() {
			fullPath := filepath.Join(RhyoliteDir, entry.Name())
			if err := watchRock(fullPath); err != nil {
				errors += err.Error()
			}
		}
	}
	if errors != "" {
		return fmt.Errorf("Watching rocks gave errors: %w", errors)
	}

	err = rockRootConfigWatcher.Add(directory)
	if err != nil {
		return fmt.Errorf("Failed to watch root dir: %w", err)
	}
	println("Watching: " + directory)

	return nil
}

func watchRock(directory string) error {
	println("Watching rock:", directory)
	err := refreshPath(directory)
	if err != nil {
		return err
	}
	rockConfigWatcher.Add(directory)
	return nil
}
func refreshPath(directory string) error {
	delta, err := filepath.Rel(RhyoliteDir, directory)
	if err != nil {
		return fmt.Errorf("Failed to parse relative path: %w", err)
	}
	err = reloadPath(delta)
	if err != nil {
		fmt.Println("Rock config error: %w", err)
	}
	return nil
}

func rockWatcherHandle() {
	for {
		select {
		case event, ok := <-rockConfigWatcher.Events:
			if !ok {
				return
			}
			if event.Op.Has(fsnotify.Create) || event.Op.Has(fsnotify.Write) || event.Op.Has(fsnotify.Remove) {
				rockContentsChanged(event.Name)
			}
		case err, ok := <-rockConfigWatcher.Errors:
			if !ok {
				return
			}
			log.Println("Rock Root Watcher error:", err.Error())
		}
	}
}

func rockContentsChanged(fullPath string) {
	directory, file := filepath.Split(fullPath)
	if file != ".rhyolite.yaml" {
		return
	}
	refreshPath(directory)
}

func closeWatchForFileChanges() {
	println("Closing the config watcher")
	rockConfigWatcher.Close()
}
