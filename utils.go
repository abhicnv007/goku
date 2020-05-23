package main

import pb "goku/entry"

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
