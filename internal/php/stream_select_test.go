package php

import (
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestSelectFDLayout(t *testing.T) {
	var set syscall.FdSet
	capacity := int(unsafe.Sizeof(set)) * 8
	for _, fd := range []int{0, 31, 32, 63, 64, capacity - 1} {
		set = syscall.FdSet{}
		if !setFD(&set, fd) || !fdReady(&set, fd) {
			t.Fatalf("descriptor %d was not set", fd)
		}
		// Inspect the ABI bytes independently of fdReady's indexing logic.
		data := unsafe.Slice((*byte)(unsafe.Pointer(&set)), int(unsafe.Sizeof(set)))
		if data[fd/8] != byte(1<<uint(fd%8)) {
			t.Fatalf("descriptor %d has incorrect ABI layout", fd)
		}
	}
	if setFD(&set, -1) || setFD(&set, capacity) || fdReady(&set, capacity) {
		t.Fatal("out-of-range descriptor accepted")
	}
}

func TestSelectPipeReadiness(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	fd := int(reader.Fd())
	for _, ready := range []bool{false, true} {
		if ready {
			if _, err := writer.Write([]byte("x")); err != nil {
				t.Fatal(err)
			}
		}
		var read syscall.FdSet
		if !setFD(&read, fd) {
			t.Fatalf("pipe descriptor %d exceeds select capacity", fd)
		}
		timeout := syscall.Timeval{}
		if err := selectFDs(fd+1, &read, nil, nil, &timeout); err != nil {
			t.Fatal(err)
		}
		if fdReady(&read, fd) != ready {
			t.Fatalf("pipe readiness: want %v", ready)
		}
	}
}
