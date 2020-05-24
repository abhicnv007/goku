package main

import (
	pb "github.com/abhicnv007/goku/entry"
	"testing"
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
		t.Errorf("Replay Entries INSERT expected 2, got %d", l)
	}
}
