package goku

import (
	pb "github.com/abhicnv007/goku/entry"
	"os"
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

func createFile(fl string) (f *os.File, err error) {
	f, err = os.OpenFile(fl, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	return
}

func openFile(fl string) (f *os.File, err error) {
	f, err = os.OpenFile(fl, os.O_APPEND|os.O_WRONLY, 0666)
	return
}
