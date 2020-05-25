package goku

import (
	"encoding/binary"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/proto"

	pb "github.com/abhicnv007/goku/entry"
)

func writeEntry(entry *pb.Entry, f *os.File) {
	// marshal the entry to bytes
	out, err := proto.Marshal(entry)
	if err != nil {
		log.Fatalln("Failed to encode entry:", err)
	}

	// get the size of the entry
	sz := make([]byte, 4)
	binary.LittleEndian.PutUint32(sz, uint32(len(out)))

	// write the size first
	if _, err := f.Write(sz); err != nil {
		log.Fatalln("Failed to write sz:", err)
	}

	// write the entry
	if _, err := f.Write(out); err != nil {
		log.Fatalln("Failed to write entry:", err)
	}
}

func readEntry(f *os.File) []*pb.Entry {

	sizeBuf := make([]byte, 4)
	var entries []*pb.Entry
	for {
		if _, err := io.ReadFull(f, sizeBuf); err != nil {
			break
		}

		size := binary.LittleEndian.Uint32(sizeBuf)
		msg := make([]byte, size)
		if _, err := io.ReadFull(f, msg); err != nil {
			log.Fatalln("Failed to read entry:", err)
		}

		entry := &pb.Entry{}
		if err := proto.Unmarshal(msg, entry); err != nil {
			log.Fatalln("Failed to parse entry:", err)
		}
		entries = append(entries, entry)
	}
	return entries
}
