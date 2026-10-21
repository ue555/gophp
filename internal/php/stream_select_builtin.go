package php

import (
	"syscall"
	"time"
	"unsafe"
)

func setFD(set *syscall.FdSet, fd int) bool {
	bits := int(unsafe.Sizeof(set.Bits[0])) * 8
	if fd < 0 || fd/bits >= len(set.Bits) {
		return false
	}
	set.Bits[fd/bits] |= 1 << uint(fd%bits)
	return true
}

func fdReady(set *syscall.FdSet, fd int) bool {
	bits := int(unsafe.Sizeof(set.Bits[0])) * 8
	return fd >= 0 && fd/bits < len(set.Bits) && set.Bits[fd/bits]&(1<<uint(fd%bits)) != 0
}

func (e *evaluator) callStreamSelect(node functionCall) (value, error) {
	if len(node.args) < 4 || len(node.args) > 5 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "stream_select() expects 4 to 5 arguments"}
	}
	var places [3]referencePlace
	var arrays [3]value
	for i := 0; i < 3; i++ {
		place, err := e.referencePlace(node.args[i], false)
		if err != nil {
			return value{}, err
		}
		places[i] = place
		arrays[i] = dereference(place.read())
		if arrays[i].kind != arrayKind && arrays[i].kind != nullKind {
			return value{}, builtinTypeError(node.line, "stream_select", i+1, [...]string{"read", "write", "except"}[i], "?array", arrays[i])
		}
	}
	seconds, err := e.eval(node.args[3])
	if err != nil {
		return value{}, err
	}
	var micros value
	if len(node.args) == 5 {
		micros, err = e.eval(node.args[4])
		if err != nil {
			return value{}, err
		}
	}
	var readFDs, writeFDs, exceptFDs syscall.FdSet
	sets := []*syscall.FdSet{&readFDs, &writeFDs, &exceptFDs}
	maxFD, count := -1, 0
	for i, array := range arrays {
		if array.kind != arrayKind || array.a == nil {
			continue
		}
		for _, entry := range array.a.entries {
			v := dereference(entry.value)
			if v.kind != resourceKind || v.stream == nil || v.stream.file == nil || v.stream.closed {
				continue
			}
			fd := int(v.stream.file.Fd())
			if !setFD(sets[i], fd) {
				return boolean(false), nil
			}
			if fd > maxFD {
				maxFD = fd
			}
			count++
		}
	}
	if count == 0 {
		return value{}, &runtimeError{node.line, "ValueError", "stream_select(): No stream arrays were passed"}
	}
	var deadline time.Time
	if seconds.kind != nullKind {
		if seconds.kind != intKind || seconds.i < 0 {
			return boolean(false), nil
		}
		duration := time.Duration(seconds.i)*time.Second + time.Duration(micros.i)*time.Microsecond
		deadline = time.Now().Add(duration)
	}
	originalRead, originalWrite, originalExcept := readFDs, writeFDs, exceptFDs
	for {
		readFDs, writeFDs, exceptFDs = originalRead, originalWrite, originalExcept
		var timeout *syscall.Timeval
		if !deadline.IsZero() {
			remaining := time.Until(deadline)
			if remaining < 0 {
				remaining = 0
			}
			v := syscall.NsecToTimeval(remaining.Nanoseconds())
			timeout = &v
		}
		err := selectFDs(maxFD+1, &readFDs, &writeFDs, &exceptFDs, timeout)
		if err == syscall.EINTR {
			continue
		}
		if err != nil {
			if warning := e.warn(node.line, "stream_select(): "+err.Error()); warning != nil {
				return value{}, warning
			}
			return boolean(false), nil
		}
		break
	}
	ready := 0
	for i, array := range arrays {
		if array.kind != arrayKind || array.a == nil {
			continue
		}
		filtered := arrayValue()
		for _, entry := range array.a.entries {
			v := dereference(entry.value)
			if v.kind == resourceKind && v.stream != nil && v.stream.file != nil && fdReady(sets[i], int(v.stream.file.Fd())) {
				filtered.a = filtered.a.with(entry.key, v)
				ready++
			}
		}
		if err := places[i].assign(filtered); err != nil {
			return value{}, err
		}
	}
	return integer(int64(ready)), nil
}
