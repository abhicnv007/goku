// Package goku is a library that implements an in memory but persistant datastore
package goku

import (
	"log"
	"os"

	pb "github.com/abhicnv007/goku/entry"
	"github.com/abhicnv007/goku/filelock"
)

// Goku is the core data strcture of the library.
// It stores all the key value pairs and pointers to the log file
type Goku struct {
	items  map[string]string
	logPtr *os.File
	dbPath string
}

// New creates a new instance of Goku. Each instance has it's own private store.
// It creates a append log at the path given, which contains all operations performed on Goku allowing recovery of
// data in case of program crash.
//
// Example:
//	g := goku.New(".goku_data")
func New(dbPath string) Goku {
	_, err := os.Stat(dbPath)

	if os.IsNotExist(err) {
		f, err := createDataFile(dbPath)
		if err != nil {
			log.Fatalln("Failed to create data file", dbPath, "Got error", err)
		}
		return Goku{items: make(map[string]string), dbPath: dbPath, logPtr: f}
	}

	// open the data file (file is opened in append mode)
	f, err := openDataFile(dbPath)
	if err != nil {
		log.Fatalln("Failed to open file data file", dbPath, "Got error", err)
	}
	entries := readEntry(f)
	return Goku{items: replayEntries(entries), dbPath: dbPath, logPtr: f}

}

// Add a key value pair to the Goku instance and also persists the operation to disk.
func (g *Goku) Add(key string, value string) {

	g.items[key] = value

	writeEntry(&pb.Entry{
		Key:   key,
		Value: value,
		Type:  pb.Entry_INSERT,
	}, g.logPtr)
}

// Get the value saved for a key
func (g *Goku) Get(key string) (val string, ok bool) {
	val, ok = g.items[key]
	return
}

// Count returns the number of items saved
func (g *Goku) Count() int {
	return len(g.items)
}

// Close file handlers and remove lock files
//
// Example:
//	g := goku.New(".db")
//	defer g.Close()
func (g *Goku) Close() {
	g.logPtr.Close()

	// if there are no items, go ahead and delete the datafile
	// else, just delete the lock file
	if len(g.items) == 0 {
		deleteDataFile(g.dbPath)
	} else {
		filelock.Release(g.dbPath)
	}

}

// Clear deletes all elements and truncates the log from disk.
//
// Example:
//	g := goku.New(".db")
//	defer g.Clear()
func (g *Goku) Clear() {

	// truncate the file
	g.logPtr.Truncate(0)
	g.logPtr.Seek(0, 0)

	// allow the previous items to be garbage collected
	g.items = make(map[string]string)
}
