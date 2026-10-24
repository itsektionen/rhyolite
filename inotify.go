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

	err = watchAllRocksInFolder(config.RockDir)
	if err != nil {
		log.Fatalf("Failed to watch rocks: %s", err.Error())
	}

	go watcherHandle(rockConfigWatcher, rockContentsChanged)
	go watcherHandle(rockRootConfigWatcher, rockRootChanged)
}

func watcherHandle(w *fsnotify.Watcher, callback func(fsnotify.Event)) {
	for {
		select {
		case event, ok := <-rockConfigWatcher.Events:
			if !ok {
				return
			}
			if event.Op.Has(fsnotify.Create) || event.Op.Has(fsnotify.Write) || event.Op.Has(fsnotify.Remove) {
				callback(event)
			}
		case err, ok := <-rockConfigWatcher.Errors:
			if !ok {
				return
			}
			log.Println("Rock Watcher error:", err.Error())
		}
	}
}

func rockRootChanged(event fsnotify.Event) {
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
}

func watchAllRocksInFolder(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("Cannot read directory:", directory)
	}
	errors := ""
	for _, entry := range entries {
		if entry.IsDir() {
			fullPath := config.joinRockDir(entry.Name())
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

	return nil
}

func watchRock(directory string) error {
	err := refreshPath(directory)
	if err != nil {
		return err
	}
	rockConfigWatcher.Add(directory)
	return nil
}
func refreshPath(directory string) error {
	delta, err := filepath.Rel(config.RockDir, directory)
	if err != nil {
		return fmt.Errorf("Failed to parse relative path: %w", err)
	}
	err = reloadRock(delta)
	if err != nil {
		log.Println("Rock config error:", err.Error())
	}
	return nil
}

func rockContentsChanged(event fsnotify.Event) {
	directory, file := filepath.Split(event.Name)
	if file != ".rhyolite.yaml" {
		return
	}
	refreshPath(directory)
}

func closeWatchForFileChanges() {
	log.Println("Closing the rock watchers")
	rockConfigWatcher.Close()
	rockRootConfigWatcher.Close()
}
