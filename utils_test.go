package goku

import (
	"os"
	"testing"

	pb "github.com/abhicnv007/goku/entry"
	"github.com/abhicnv007/goku/filelock"
)

func TestReplayEntries(t *testing.T) {
	entries := []*pb.Entry{
		{
			Key:   "1",
			Value: "19",
			Type:  pb.Entry_INSERT,
		},
		{
			Key:   "2",
			Value: "20",
			Type:  pb.Entry_INSERT,
		}, {
			Key:   "3",
			Value: "1234",
			Type:  pb.Entry_INSERT,
		},
	}

	items := replayEntries(entries)
	if l := len(items); l != 3 {
		t.Error("Replay Entries INSERT expected 2, got", l)
	}
}

func TestCreateDataFile(t *testing.T) {
	datafile := ".goku"
	f, err := createDataFile(datafile)
	if err != nil {
		t.Error("Create Data file, did not expect any errors, got", err)
	}
	f.Close()

	//cleanup
	deleteDataFile(datafile)
}

func TestOpenDataFile(t *testing.T) {
	datafile := ".goku"

	// create an empty datafile
	f, err := os.Create(datafile)
	if err != nil {
		t.Error("Open data file, cannot create sample datafile, got", err)
	}
	f.Close()

	f, err = openDataFile(datafile)
	if err != nil {
		t.Error("Open data file, could not open data file, got", err)
	}

	// cleanup
	f.Close()
	os.Remove(datafile)
	filelock.Release(datafile)

}
