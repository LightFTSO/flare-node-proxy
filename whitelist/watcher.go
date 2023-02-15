// sourced from https://github.com/fsnotify/fsnotify/blob/main/cmd/fsnotify/file.go
package whitelist

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func watchFile(files ...string) {
	if len(files) < 1 {
		log.Fatal("Must specify at least one file to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Debug("creating a new watcher: %s", err)
	}
	defer w.Close()

	// Start listening for events.
	go fileLoop(w, files)

	// Add all files from the commandline.
	for _, p := range files {
		st, err := os.Lstat(p)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if st.IsDir() {
			log.Fatalf("%q is a directory, not a file", p)
		}

		// Watch the directory, not the file itself.
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Fatalf("%s", err)
		}
		err = w.Add(filepath.Dir(abs))
		log.Infof("Watching %s", abs)
		if err != nil {
			log.Fatalf("%q: %s", p, err)
		}
	}

	<-make(chan struct{}) // Block forever
}

func fileLoop(w *fsnotify.Watcher, files []string) {
	for {
		select {

		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Fatalf("ERROR: %s", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// Ignore files we're not interested in. Can use a
			// map[string]struct{} if you have a lot of files, but for just a
			// few files simply looping over a slice is faster.
			var found bool
			for _, f := range files {
				if f == e.Name {
					found = true
				}
			}
			if !found {
				continue
			} else {
				// Just print the event
				log.Debugf("File event %s %s", e.Op, e.Name)

				// update the whitelist
				readWhitelistFile(&e.Name)
			}
		}

	}

}
