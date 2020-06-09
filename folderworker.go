package folderworker

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	err := readDir()
	if err != nil {
		log.Fatal(err)
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

func readDir() error {
	files, err := ioutil.ReadDir(root)
	fmt.Println("AAAA")
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println("BBB")
		name := f.Name()
		if f.IsDir() || filepath.Ext(name) != ".url" {
			continue
		}
		content, err := ioutil.ReadFile(name)
		if err != nil {
			continue
		}

		fmt.Println(string(content))
	}

	return nil
}
