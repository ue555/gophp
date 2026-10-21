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
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// PCRE2 is the source of both match results and PHP preg error state.
func pcre2MatchAll(pattern, subject string) ([][]int, []string, error) {
	return pcre2Matches(pattern, subject, 0, -1)
}

func pcre2Matches(pattern, subject string, start, limit int) ([][]int, []string, error) {
	pattern = strings.TrimSpace(pattern)
	if len(pattern) < 2 {
		return nil, nil, &pregError{1, "No ending delimiter found"}
	}
	delim := pattern[0]
	if delim == 0 || delim == '\\' || delim >= 'a' && delim <= 'z' || delim >= 'A' && delim <= 'Z' || delim >= '0' && delim <= '9' {
		return nil, nil, &pregError{1, "Delimiter must not be alphanumeric, backslash, or NUL"}
	}
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
		return nil, nil, &pregError{1, "No ending delimiter found"}
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
			options |= C.PCRE2_UTF | C.PCRE2_UCP
		case 'D':
			options |= C.PCRE2_DOLLAR_ENDONLY
		case 'A':
			options |= C.PCRE2_ANCHORED
		case 'U':
			options |= C.PCRE2_UNGREEDY
		case 'n':
			options |= C.PCRE2_NO_AUTO_CAPTURE
		case 'S': // PCRE2 optimizes without the legacy study hint.
		case 'J':
			options |= C.PCRE2_DUPNAMES
		default:
			return nil, nil, &pregError{1, fmt.Sprintf("Unknown modifier '%c'", flag)}
		}
	}
	patternBytes := C.CBytes([]byte(body))
	defer C.free(patternBytes)
	var codeErr C.int
	var errOffset C.PCRE2_SIZE
	code := C.pcre2_compile_8((*C.PCRE2_UCHAR8)(patternBytes), C.PCRE2_SIZE(len(body)), options, &codeErr, &errOffset, nil)
	if code == nil {
		buf := make([]byte, 256)
		C.pcre2_get_error_message_8(codeErr, (*C.PCRE2_UCHAR8)(unsafe.Pointer(&buf[0])), C.PCRE2_SIZE(len(buf)))
		return nil, nil, &pregError{1, fmt.Sprintf("Compilation failed: %s at offset %d", strings.TrimRight(string(buf), "\x00"), uint64(errOffset))}
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
		return nil, nil, &pregError{1, ""}
	}
	defer C.pcre2_match_data_free_8(data)
	input := C.CBytes([]byte(subject))
	defer C.free(input)
	var matches [][]int
	matchOptions := C.uint32_t(0)
	for offset := start; offset <= len(subject) && (limit < 0 || len(matches) < limit); {
		rc := C.pcre2_match_8(code, (*C.PCRE2_UCHAR8)(input), C.PCRE2_SIZE(len(subject)), C.PCRE2_SIZE(offset), matchOptions, data, nil)
		if rc == C.PCRE2_ERROR_NOMATCH {
			if matchOptions == 0 {
				break
			}
			matchOptions = 0
			if offset == len(subject) {
				break
			}
			if options&C.PCRE2_UTF != 0 {
				_, width := utf8.DecodeRuneInString(subject[offset:])
				offset += width
			} else {
				offset++
			}
			continue
		}
		if rc < 0 {
			errorCode := int64(1)
			switch {
			case rc == C.PCRE2_ERROR_MATCHLIMIT:
				errorCode = 2
			case rc == C.PCRE2_ERROR_DEPTHLIMIT:
				errorCode = 3
			case rc <= C.PCRE2_ERROR_UTF8_ERR1 && rc >= C.PCRE2_ERROR_UTF8_ERR21:
				errorCode = 4
			case rc == C.PCRE2_ERROR_BADUTFOFFSET:
				errorCode = 5
			case rc == C.PCRE2_ERROR_JIT_STACKLIMIT:
				errorCode = 6
			}
			return nil, nil, &pregError{errorCode, ""}
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
		offset = row[1]
		matchOptions = 0
		if row[0] == row[1] {
			matchOptions = C.PCRE2_NOTEMPTY_ATSTART | C.PCRE2_ANCHORED
		}
	}

	return matches, names, nil
}

type pregError struct {
	code    int64
	warning string
}

func (p *pregError) Error() string {
	if p.warning != "" {
		return p.warning
	}
	return pregErrorMessage(p.code)
}
func pregErrorMessage(code int64) string {
	switch code {
	case 0:
		return "No error"
	case 2:
		return "Backtrack limit exhausted"
	case 3:
		return "Recursion limit exhausted"
	case 4:
		return "Malformed UTF-8 characters, possibly incorrectly encoded"
	case 5:
		return "The offset did not correspond to the beginning of a valid UTF-8 code point"
	case 6:
		return "JIT stack limit exhausted"
	}
	return "Internal error"
}
func (e *evaluator) pregFailure(name string, line int, err error) (value, error) {
	p, ok := err.(*pregError)
	if !ok {
		return value{}, err
	}
	e.pregLastError = p.code
	if p.warning != "" {
		if err := e.warn(line, name+"(): "+p.warning); err != nil {
			return value{}, err
		}
	}
	if name == "preg_replace" || name == "preg_replace_callback" {
		return value{}, nil
	}
	return boolean(false), nil
}
func pcreLibraryInfo() (string, int64, int64, bool) {
	buf := make([]byte, 256)
	C.pcre2_config_8(C.PCRE2_CONFIG_VERSION, unsafe.Pointer(&buf[0]))
	s := strings.TrimRight(string(buf), "\x00")
	pieces := strings.SplitN(s, ".", 2)
	major, _ := strconv.ParseInt(pieces[0], 10, 64)
	minor := int64(0)
	if len(pieces) > 1 {
		minor, _ = strconv.ParseInt(strings.Fields(pieces[1])[0], 10, 64)
	}
	var jit C.uint32_t
	C.pcre2_config_8(C.PCRE2_CONFIG_JIT, unsafe.Pointer(&jit))
	return s, major, minor, jit != 0
}
