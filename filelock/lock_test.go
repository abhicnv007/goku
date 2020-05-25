package filelock

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestAquire(t *testing.T) {

	datafile := ".goku"

	// Test 1. Lockfile does not already exist
	if err := Acquire(datafile); err != nil {
		t.Error("Acquire, got err:", err)
	}

	// Test 2. Lockfile already exists, must return error
	if err := Acquire(datafile); err == nil {
		t.Error("Acquire, able to acquire lock on the same file again")
	}

	// Test 3. Another PID in lockfile, must return error
	lock, err := os.OpenFile(lockFileName(datafile), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		t.Error("Acquire, could not create dummy lockfile:", err)
	}

	// write the parent's pid
	writeLockData(lock, os.Getppid(), uuid.New())
	lock.Close()

	if err := Acquire(datafile); err == nil {
		t.Error("Acquire, able to acquire other process's lock:", err)
	}

	os.Remove(lockFileName(datafile))

}

func TestRelease(t *testing.T) {
	datafile := ".goku"

	if err := Acquire(datafile); err != nil {
		t.Error("Release, could not acquire:", err)
	}

	// Test 1. Lockfile exists
	if err := Release(datafile); err != nil {
		t.Error("Release, could not release after acquire:", err)
	}

	// Test 2. Lockfile does exist, must return error
	if err := Release(datafile); err == nil {
		t.Error("Release, able to release lock on the same file again")
	}

}
