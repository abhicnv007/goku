package main

import (
	pb "goku/entry"
	"log"
	"os"
)

// Goku is the core data strcture of the library
type Goku struct {
	items   map[string]string
	logName string
	logPtr  *os.File
}

// New creates a new instance of Goku
func New() Goku {
	logName := "goku_data"
	_, err := os.Stat(logName)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(logName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln("Failed to create data file", logName, "Got error", err)
		}
		return Goku{items: make(map[string]string), logName: logName, logPtr: f}
	}

	entries := readEntry(logName)

	f, err := os.OpenFile(logName, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open file data file", logName, "Got error", err)
	}
	return Goku{items: replayEntries(entries), logName: logName, logPtr: f}

}

// Add a key value pair to the Goku instance
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

// Close all file handlers and free up memory
func (g *Goku) Close() {
	g.logPtr.Close()
}

// Clear deletes all elements and removes its log from disk
func (g *Goku) Clear() {
	g.logPtr.Close()
	if err := os.Remove(g.logName); err != nil {
		log.Fatal("Could not remove file log file", g.logName, "; got error:", err)
	}
	f, err := os.OpenFile(g.logName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Failed to create data file", g.logName, "Got error", err)
	}

	g.logPtr = f
	g.items = make(map[string]string)
}
