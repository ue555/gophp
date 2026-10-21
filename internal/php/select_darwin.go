package php

import "syscall"

func selectFDs(n int, read, write, except *syscall.FdSet, timeout *syscall.Timeval) error {
	return syscall.Select(n, read, write, except, timeout)
}
