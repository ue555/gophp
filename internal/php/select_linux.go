package php

import "syscall"

func selectFDs(n int, read, write, except *syscall.FdSet, timeout *syscall.Timeval) error {
	_, err := syscall.Select(n, read, write, except, timeout)
	return err
}
