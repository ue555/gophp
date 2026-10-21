package php

/*
#cgo darwin CFLAGS: -I/opt/homebrew/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -lpcre2-8
#cgo linux LDFLAGS: -lpcre2-8
#define PCRE2_CODE_UNIT_WIDTH 8
#include <pcre2.h>
#include <stdlib.h>
*/
import "C"

import (
	"strings"
	"unsafe"
)

// pcre2MatchAll provides the PCRE constructs that Go's regexp package cannot
// express, without changing the fast path for patterns already supported there.
func pcre2MatchAll(pattern, subject string) ([][]int, []string, error) {
	pattern = strings.TrimSpace(pattern)
	if len(pattern) < 2 {
		return nil, nil, &runtimeError{0, "Error", "gophp: unsupported regular expression delimiter"}
	}
	delim := pattern[0]
	closing := delim
	switch delim {
	case '{':
		closing = '}'
	case '[':
		closing = ']'
	case '(':
		closing = ')'
	case '<':
		closing = '>'
	}
	end := strings.LastIndexByte(pattern, closing)
	if end < 1 {
		return nil, nil, &runtimeError{0, "Error", "gophp: unsupported regular expression delimiter"}
	}
	body, flags := pattern[1:end], pattern[end+1:]
	options := C.uint32_t(0)
	for _, flag := range flags {
		switch flag {
		case 'i':
			options |= C.PCRE2_CASELESS
		case 'm':
			options |= C.PCRE2_MULTILINE
		case 's':
			options |= C.PCRE2_DOTALL
		case 'x':
			options |= C.PCRE2_EXTENDED
		case 'u':
			options |= C.PCRE2_UTF
		case 'D':
			options |= C.PCRE2_DOLLAR_ENDONLY
		case 'A':
			options |= C.PCRE2_ANCHORED
		case 'J':
			options |= C.PCRE2_DUPNAMES
		default:
			return nil, nil, &runtimeError{0, "Error", "gophp: unsupported regular expression modifier"}
		}
	}
	patternBytes := C.CBytes([]byte(body))
	defer C.free(patternBytes)
	var codeErr C.int
	var errOffset C.PCRE2_SIZE
	code := C.pcre2_compile_8((*C.PCRE2_UCHAR8)(patternBytes), C.PCRE2_SIZE(len(body)), options, &codeErr, &errOffset, nil)
	if code == nil {
		return nil, nil, &runtimeError{0, "Error", "gophp: invalid PCRE pattern"}
	}
	defer C.pcre2_code_free_8(code)
	var captureCount C.uint32_t
	C.pcre2_pattern_info_8(code, C.PCRE2_INFO_CAPTURECOUNT, unsafe.Pointer(&captureCount))
	names := make([]string, int(captureCount)+1)
	var nameCount C.uint32_t
	C.pcre2_pattern_info_8(code, C.PCRE2_INFO_NAMECOUNT, unsafe.Pointer(&nameCount))
	if nameCount > 0 {
		var size C.uint32_t
		var table *C.PCRE2_UCHAR8
		C.pcre2_pattern_info_8(code, C.PCRE2_INFO_NAMEENTRYSIZE, unsafe.Pointer(&size))
		C.pcre2_pattern_info_8(code, C.PCRE2_INFO_NAMETABLE, unsafe.Pointer(&table))
		for i := 0; i < int(nameCount); i++ {
			entry := unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(table), uintptr(i)*uintptr(size))), int(size))
			index := int(entry[0])<<8 | int(entry[1])
			if index < len(names) {
				end := len(entry)
				if n := strings.IndexByte(string(entry[2:]), 0); n >= 0 {
					end = n + 2
				}
				names[index] = string(entry[2:end])
			}
		}
	}
	data := C.pcre2_match_data_create_from_pattern_8(code, nil)
	if data == nil {
		return nil, nil, &runtimeError{0, "Error", "gophp: PCRE match allocation failed"}
	}
	defer C.pcre2_match_data_free_8(data)
	input := C.CBytes([]byte(subject))
	defer C.free(input)
	var matches [][]int
	for offset := 0; offset <= len(subject); {
		rc := C.pcre2_match_8(code, (*C.PCRE2_UCHAR8)(input), C.PCRE2_SIZE(len(subject)), C.PCRE2_SIZE(offset), 0, data, nil)
		if rc == C.PCRE2_ERROR_NOMATCH {
			break
		}
		if rc < 0 {
			return nil, nil, &runtimeError{0, "Error", "gophp: PCRE match failed"}
		}
		vector := unsafe.Slice(C.pcre2_get_ovector_pointer_8(data), int(captureCount+1)*2)
		row := make([]int, len(vector))
		for i, v := range vector {
			if uint64(v) == ^uint64(0) {
				row[i] = -1
			} else {
				row[i] = int(v)
			}
		}
		matches = append(matches, row)
		if row[1] > offset {
			offset = row[1]
		} else {
			offset++
		}
	}
	return matches, names, nil
}
