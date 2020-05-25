package goku

import (
	"log"
	"os"

	pb "github.com/abhicnv007/goku/entry"
	"github.com/abhicnv007/goku/filelock"
)

func replayEntries(entries []*pb.Entry) map[string]string {
	items := make(map[string]string)
	for _, e := range entries {
		switch e.Type {
		case pb.Entry_INSERT:
			items[e.Key] = e.Value
		case pb.Entry_DELETE:
			_, ok := items[e.Key]
			if ok {
				delete(items, e.Key)
			}
		}
	}

	return items
}

// create data file if one is not found
func createDataFile(datafile string) (f *os.File, err error) {

	// return error if file is already present
	_, err = os.Stat(datafile)
	if !os.IsNotExist(err) {
		return nil, err
	}

	if err := filelock.Acquire(datafile); err != nil {
		return nil, err
	}

	// create the datafile
	f, err = os.OpenFile(datafile, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	return
}

// delete data and lock files
func deleteDataFile(datafile string) {
	if err := os.Remove(datafile); err != nil {
		log.Fatalln("Could not delete datafile, got error", err)
	}

	filelock.Release(datafile)
}

// open a datafile and ensure the locks are in order
func openDataFile(datafile string) (*os.File, error) {
	if err := filelock.Acquire(datafile); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(datafile, os.O_APPEND|os.O_RDWR, 0666)
	return f, err
}
