// Package filelock implements a lock for mutually exclusive access of files
package filelock

import (
	"encoding/binary"
	"errors"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/google/uuid"
)

var packageUUID uuid.UUID

func init() {
	packageUUID = uuid.New()
}

// Acquire creates a lock for the file
func Acquire(f string) error {

	// open file with O_CREATE and O_EXCL, which gives an error if the file
	// was not created by this call
	lock, err := os.OpenFile(lockFileName(f), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)

	// lock already exists
	if err != nil {

		// open the lock file and inspect it's contents
		lock, err = os.OpenFile(lockFileName(f), os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		// TODO get the uid from lockfile and compare it to our uid
		// if not match, another instance of goku is trying to open the same file
		// throw an error
		lockpid, uid := readLockData(lock)
		// lock not owned by current process and the other process is running
		if lockpid != os.Getpid() && pidExists(lockpid) {
			return errors.New("File locked by another process")
		}

		// we are trying to lock this file again, return an error
		if uid == packageUUID {
			return errors.New("File already locked by this process")
		}

	}

	// write lock data now
	writeLockData(lock, os.Getpid(), packageUUID)

	lock.Close()

	return nil

}

// Release deletes the lock for the file
func Release(f string) error {

	// make sure we have acquired this lock
	// open the lock file and inspect it's contents
	lock, err := os.OpenFile(lockFileName(f), os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	lockpid, uid := readLockData(lock)
	if lockpid != os.Getpid() || uid != packageUUID {
		return errors.New("File not locked by current process")
	}

	lock.Close()
	return os.Remove(lockFileName(f))
}

func pidExists(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	if err.Error() == "os: process already finished" {
		return false
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	switch errno {
	case syscall.ESRCH:
		return false
	case syscall.EPERM:
		return true
	}
	return false
}

func lockFileName(datafile string) string {
	// create a lockfile of the same name, but with a .lock suffix
	// eg: .db and .db.lock
	filename := filepath.Base(datafile)
	dir := filepath.Dir(datafile)
	return filepath.Join(dir, filename+".lock")
}

func writeLockData(f *os.File, pid int, uid uuid.UUID) {

	// truncate the file
	f.Truncate(0)
	f.Seek(0, 0)

	// write pid of process
	p := make([]byte, 8)
	binary.LittleEndian.PutUint64(p, uint64(pid))
	if _, err := f.Write(p); err != nil {
		log.Fatalln("Failed to write PID:", err)
	}

	uidbin, err := uid.MarshalBinary()
	if err != nil {
		log.Fatalln("Failed to marshal UUID to binary:", err)
	}

	if _, err := f.Write(uidbin); err != nil {
		log.Fatalln("Failed to write uuid:", err)
	}
}

func readLockData(f *os.File) (int, uuid.UUID) {
	pid := make([]byte, 8)
	if _, err := f.Read(pid); err != nil {
		log.Fatalln("Failed to read PID:", err)
	}
	p := int(binary.LittleEndian.Uint64(pid))

	uidbin := make([]byte, 16)
	if _, err := f.Read(uidbin); err != nil {
		log.Fatalln("Failed to read uuid:", err)
	}

	// create a new uuid and then read from bytes
	uid := uuid.New()
	uid.UnmarshalBinary(uidbin)

	return p, uid
}
