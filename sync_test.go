package main

import (
	pb "goku/entry"
	"os"
	"testing"
)

func TestWriteEntry(t *testing.T) {
	fname := "sync_test_file"

	f, _ := os.Create(fname)
	f.Close()

	e := pb.Entry{
		Key:   "Hello",
		Value: "World",
		Type:  pb.Entry_INSERT,
	}
	writeEntry(&e, fname)

	e = pb.Entry{
		Key:   "Hi",
		Value: "World",
		Type:  pb.Entry_INSERT,
	}
	writeEntry(&e, fname)

	//clean up
	os.Remove(fname)
}

func TestReadEntry(t *testing.T) {
	fname := "sync_test_file"

	f, _ := os.Create(fname)
	f.Close()

	e := pb.Entry{
		Key:   "Hello",
		Value: "World",
		Type:  pb.Entry_INSERT,
	}
	writeEntry(&e, fname)

	e = pb.Entry{
		Key:   "Hi",
		Value: "World",
		Type:  pb.Entry_INSERT,
	}
	writeEntry(&e, fname)

	entries := readEntry(fname)
	if len(entries) != 2 {
		t.Errorf("Read entry: expect 2, got %d", len(entries))
	}

	//clean up
	os.Remove(fname)
}
