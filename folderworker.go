package folderworker

import (
	"flag"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var (
	root string
)

// Start dsfs
func Start() {
	flag.StringVar(&root, "root", "", "Absolute path for root directory")
	flag.Parse()

	if len(root) < 1 {
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		root = workingDir
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(root)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
