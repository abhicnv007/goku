// Package goku is a library that implements an in memory but persistant datastore
package goku

import (
	"log"
	"os"

	pb "github.com/abhicnv007/goku/entry"
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
		f, err := createFile(dbPath)
		if err != nil {
			log.Fatalln("Failed to create data file", dbPath, "Got error", err)
		}
		return Goku{items: make(map[string]string), dbPath: dbPath, logPtr: f}
	}

	entries := readEntry(dbPath)

	f, err := openFile(dbPath)
	if err != nil {
		log.Fatalln("Failed to open file data file", dbPath, "Got error", err)
	}
	return Goku{items: replayEntries(entries), dbPath: dbPath, logPtr: f}

}

// Add a key value pair to the Goku instance and also persists the operation to disk.
func (g *Goku) Add(key string, value string) {

	if g.logPtr == nil {
		f, err := createFile(g.dbPath)
		if err != nil {
			log.Fatalln("Failed to create data file", g.dbPath, "Got error", err)
		}
		g.logPtr = f
	}

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

// Close all file handlers and free up memory
//
// Example:
//	g := goku.New(".db")
//	defer g.Close()
func (g *Goku) Close() {
	g.logPtr.Close()
}

// Clear deletes all elements and removes the log from disk.
//
// Example:
//	g := goku.New(".db")
//	defer g.Clear()
func (g *Goku) Clear() {
	// close the file pointer
	g.Close()
	g.logPtr = nil

	// remove the contents of the file
	if err := os.Remove(g.dbPath); err != nil {
		log.Fatal("Could not remove file log file", g.dbPath, "; got error:", err)
	}

	g.items = make(map[string]string)
}
