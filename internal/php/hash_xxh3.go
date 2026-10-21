package php

import (
	binenc "encoding/binary"
	"fmt"
	"math/bits"
	"strings"
)

// The fixed API Platform dependency uses the 64-bit xxh3 digest as a cache
// key. This implements the unseeded input paths from xxHash's published
// algorithm. The 64-bit and 128-bit long paths share their low accumulator.
// Secret and constants: https://github.com/Cyan4973/xxHash/blob/dev/xxhash.h
var xxh3Secret = [...]byte{
	0xb8, 0xfe, 0x6c, 0x39, 0x23, 0xa4, 0x4b, 0xbe, 0x7c, 0x01, 0x81, 0x2c, 0xf7, 0x21, 0xad, 0x1c,
	0xde, 0xd4, 0x6d, 0xe9, 0x83, 0x90, 0x97, 0xdb, 0x72, 0x40, 0xa4, 0xa4, 0xb7, 0xb3, 0x67, 0x1f,
	0xcb, 0x79, 0xe6, 0x4e, 0xcc, 0xc0, 0xe5, 0x78, 0x82, 0x5a, 0xd0, 0x7d, 0xcc, 0xff, 0x72, 0x21,
	0xb8, 0x08, 0x46, 0x74, 0xf7, 0x43, 0x24, 0x8e, 0xe0, 0x35, 0x90, 0xe6, 0x81, 0x3a, 0x26, 0x4c,
	0x3c, 0x28, 0x52, 0xbb, 0x91, 0xc3, 0x00, 0xcb, 0x88, 0xd0, 0x65, 0x8b, 0x1b, 0x53, 0x2e, 0xa3,
	0x71, 0x64, 0x48, 0x97, 0xa2, 0x0d, 0xf9, 0x4e, 0x38, 0x19, 0xef, 0x46, 0xa9, 0xde, 0xac, 0xd8,
	0xa8, 0xfa, 0x76, 0x3f, 0xe3, 0x9c, 0x34, 0x3f, 0xf9, 0xdc, 0xbb, 0xc7, 0xc7, 0x0b, 0x4f, 0x1d,
	0x8a, 0x51, 0xe0, 0x4b, 0xcd, 0xb4, 0x59, 0x31, 0xc8, 0x9f, 0x7e, 0xc9, 0xd9, 0x78, 0x73, 0x64,
	0xea, 0xc5, 0xac, 0x83, 0x34, 0xd3, 0xeb, 0xc3, 0xc5, 0x81, 0xa0, 0xff, 0xfa, 0x13, 0x63, 0xeb,
	0x17, 0x0d, 0xdd, 0x51, 0xb7, 0xf0, 0xda, 0x49, 0xd3, 0x16, 0x55, 0x26, 0x29, 0xd4, 0x68, 0x9e,
	0x2b, 0x16, 0xbe, 0x58, 0x7d, 0x47, 0xa1, 0xfc, 0x8f, 0xf8, 0xb8, 0xd1, 0x7a, 0xd0, 0x31, 0xce,
	0x45, 0xcb, 0x3a, 0x8f, 0x95, 0x16, 0x04, 0x28, 0xaf, 0xd7, 0xfb, 0xca, 0xbb, 0x4b, 0x40, 0x7e,
}

const xxh3Prime1 uint64 = 11400714785074694791
const xxh3Prime2 uint64 = 14029467366897019727
const xxh3Prime3 uint64 = 1609587929392839161

func xxh3Avalanche(x uint64) uint64 {
	x ^= x >> 37
	x *= 0x165667919e3779f9
	return x ^ x>>32
}

func xxh3SmallAvalanche(x uint64) uint64 {
	x ^= x >> 33
	x *= xxh3Prime2
	x ^= x >> 29
	x *= xxh3Prime3
	return x ^ x>>32
}

func xxh128SmallAvalanche(x uint64) uint64 {
	x *= xxh3Prime2
	x ^= x >> 29
	x *= xxh3Prime3
	return x ^ x>>32
}

func xxh3Mix(input []byte, inputAt, secretAt int) uint64 {
	a := binenc.LittleEndian.Uint64(input[inputAt:]) ^ binenc.LittleEndian.Uint64(xxh3Secret[secretAt:])
	b := binenc.LittleEndian.Uint64(input[inputAt+8:]) ^ binenc.LittleEndian.Uint64(xxh3Secret[secretAt+8:])
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func xxh3Digest(input []byte) uint64 {
	n := len(input)
	if n > 240 {
		// XXH3_finalizeLong_128b.low64 is XXH3_finalizeLong_64b.
		_, low := xxh128(input)
		return low
	}
	s := xxh3Secret[:]
	if n <= 16 {
		if n > 8 {
			lo := binenc.LittleEndian.Uint64(input) ^ (binenc.LittleEndian.Uint64(s[24:]) ^ binenc.LittleEndian.Uint64(s[32:]))
			hi := binenc.LittleEndian.Uint64(input[n-8:]) ^ (binenc.LittleEndian.Uint64(s[40:]) ^ binenc.LittleEndian.Uint64(s[48:]))
			mulHi, mulLo := bits.Mul64(lo, hi)
			return xxh3Avalanche(uint64(n) + bits.ReverseBytes64(lo) + hi + (mulHi ^ mulLo))
		}
		if n >= 4 {
			v := uint64(binenc.LittleEndian.Uint32(input[n-4:])) | uint64(binenc.LittleEndian.Uint32(input))<<32
			v ^= binenc.LittleEndian.Uint64(s[8:]) ^ binenc.LittleEndian.Uint64(s[16:])
			v ^= bits.RotateLeft64(v, 49) ^ bits.RotateLeft64(v, 24)
			v *= 0x9fb21c651e98df25
			v ^= v>>35 + uint64(n)
			v *= 0x9fb21c651e98df25
			return v ^ v>>28
		}
		if n == 0 {
			return 0x2d06800538d394c2
		}
		first, middle, last := uint32(input[0]), uint32(input[n>>1]), uint32(input[n-1])
		packed := uint64(first<<16 | middle<<24 | last | uint32(n)<<8)
		return xxh3SmallAvalanche(packed ^ uint64(binenc.LittleEndian.Uint32(s[0:])^binenc.LittleEndian.Uint32(s[4:])))
	}
	acc := uint64(n) * xxh3Prime1
	if n <= 128 {
		acc += xxh3Mix(input, 0, 0) + xxh3Mix(input, n-16, 16)
		if n > 32 {
			acc += xxh3Mix(input, 16, 32) + xxh3Mix(input, n-32, 48)
		}
		if n > 64 {
			acc += xxh3Mix(input, 32, 64) + xxh3Mix(input, n-48, 80)
		}
		if n > 96 {
			acc += xxh3Mix(input, 48, 96) + xxh3Mix(input, n-64, 112)
		}
		return xxh3Avalanche(acc)
	}
	for i := 0; i < 8; i++ {
		acc += xxh3Mix(input, i*16, i*16)
	}
	acc = xxh3Avalanche(acc)
	for at := 128; at < n&^15; at += 16 {
		acc += xxh3Mix(input, at, at-125)
	}
	acc += xxh3Mix(input, n-16, 119)
	return xxh3Avalanche(acc)
}

// xxh128 computes the unseeded XXH3 128-bit digest. The two lanes are kept
// separate because PHP's hash() renders their canonical big-endian order.
// Algorithm constants and rounds follow xxHash (BSD-2-Clause):
// https://github.com/Cyan4973/xxHash/blob/dev/xxhash.h
func xxh128(input []byte) (high, low uint64) {
	n := len(input)
	s := xxh3Secret[:]
	read := func(at int) uint64 { return binenc.LittleEndian.Uint64(input[at:]) }
	secret := func(at int) uint64 { return binenc.LittleEndian.Uint64(s[at:]) }
	fold := func(a, b uint64) uint64 { hi, lo := bits.Mul64(a, b); return hi ^ lo }
	if n == 0 {
		return 0x99aa06d3014798d8, 0x6001c324468d497f
	}
	if n <= 3 {
		first, middle, last := uint32(input[0]), uint32(input[n>>1]), uint32(input[n-1])
		lo := uint64(first<<16 | middle<<24 | last | uint32(n)<<8)
		hi := uint64(bits.RotateLeft32(bits.ReverseBytes32(uint32(lo)), 13))
		return xxh128SmallAvalanche(hi ^ uint64(binenc.LittleEndian.Uint32(s[8:])^binenc.LittleEndian.Uint32(s[12:]))), xxh128SmallAvalanche(lo ^ uint64(binenc.LittleEndian.Uint32(s)^binenc.LittleEndian.Uint32(s[4:])))
	}
	if n <= 8 {
		v := uint64(binenc.LittleEndian.Uint32(input)) | uint64(binenc.LittleEndian.Uint32(input[n-4:]))<<32
		keyed := v ^ (secret(16) ^ secret(24))
		hi, lo := bits.Mul64(keyed, xxh3Prime1+(uint64(n)<<2))
		hi += lo << 1
		lo ^= hi >> 3
		lo ^= lo >> 35
		lo *= 0x9fb21c651e98df25
		lo ^= lo >> 28
		return xxh3Avalanche(hi), lo
	}
	if n <= 16 {
		a, b := read(0), read(n-8)
		hi, lo := bits.Mul64(a^b^secret(32)^secret(40), xxh3Prime1)
		lo += uint64(n-1) << 54
		b ^= secret(48) ^ secret(56)
		hi += b + uint64(uint32(b))*(2246822519-1)
		lo ^= bits.ReverseBytes64(hi)
		mulHi, mulLo := bits.Mul64(lo, xxh3Prime2)
		return xxh3Avalanche(mulHi + hi*xxh3Prime2), xxh3Avalanche(mulLo)
	}
	lo := uint64(n) * xxh3Prime1
	var hi uint64
	mix := func(at, from, keyAt int) {
		a, b, c, d := read(at), read(at+8), read(from), read(from+8)
		hi += fold(c^secret(keyAt+16), d^secret(keyAt+24))
		hi ^= a + b
		lo += fold(a^secret(keyAt), b^secret(keyAt+8))
		lo ^= c + d
	}
	finish := func() (uint64, uint64) {
		hi, lo = lo*xxh3Prime1+hi*9650029242287828579+uint64(n)*xxh3Prime2, hi+lo
		return -xxh3Avalanche(hi), xxh3Avalanche(lo)
	}
	if n <= 128 {
		if n > 96 {
			mix(48, n-64, 96)
		}
		if n > 64 {
			mix(32, n-48, 64)
		}
		if n > 32 {
			mix(16, n-32, 32)
		}
		mix(0, n-16, 0)
		return finish()
	}
	if n <= 240 {
		for at := 0; at < 128; at += 32 {
			mix(at, at+16, at)
		}
		hi, lo = xxh3Avalanche(hi), xxh3Avalanche(lo)
		for at := 128; at < n&^31; at += 32 {
			// The sliding secret starts at offset 3 for the first trailing group.
			keyAt := at - 125
			a, b, c, d := read(at), read(at+8), read(at+16), read(at+24)
			hi += fold(c^secret(keyAt+16), d^secret(keyAt+24))
			hi ^= a + b
			lo += fold(a^secret(keyAt), b^secret(keyAt+8))
			lo ^= c + d
		}
		a, b, c, d := read(n-32), read(n-24), read(n-16), read(n-8)
		hi += fold(a^secret(119), b^secret(127))
		hi ^= c + d
		lo += fold(c^secret(103), d^secret(111))
		lo ^= a + b
		return finish()
	}
	// Long messages use 64-byte stripes, 1 KiB blocks, and accumulator
	// scrambling. This scalar path is used on every host architecture.
	acc := [8]uint64{3266489917, xxh3Prime1, xxh3Prime2, xxh3Prime3, 9650029242287828579, 2246822519, 2870177450012600261, 2654435761}
	stripe := func(at, keyAt int) {
		for i := 0; i < 8; i++ {
			v := read(at + i*8)
			k := v ^ secret(keyAt+i*8)
			acc[i^1] += v
			acc[i] += uint64(uint32(k)) * (k >> 32)
		}
	}
	scramble := func() {
		for i := 0; i < 8; i++ {
			acc[i] ^= acc[i] >> 47
			acc[i] ^= secret(128 + i*8)
			acc[i] *= 2654435761
		}
	}
	at := 0
	remaining := n
	for remaining > 1024 {
		for j := 0; j < 16; j++ {
			stripe(at+j*64, j*8)
		}
		scramble()
		at += 1024
		remaining -= 1024
	}
	if remaining > 0 {
		count := (remaining - 1) / 64
		for j := 0; j < count; j++ {
			stripe(at+j*64, j*8)
		}
		stripe(n-64, 121)
	}
	lo += fold(acc[0]^secret(11), acc[1]^secret(19)) + fold(acc[2]^secret(27), acc[3]^secret(35)) + fold(acc[4]^secret(43), acc[5]^secret(51)) + fold(acc[6]^secret(59), acc[7]^secret(67))
	hi = ^(uint64(n) * xxh3Prime2)
	hi += fold(acc[0]^secret(117), acc[1]^secret(125)) + fold(acc[2]^secret(133), acc[3]^secret(141)) + fold(acc[4]^secret(149), acc[5]^secret(157)) + fold(acc[6]^secret(165), acc[7]^secret(173))
	return xxh3Avalanche(hi), xxh3Avalanche(lo)
}

func (e *evaluator) callHash(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "algo"}, {name: "data"}, {name: "binary"}, {name: "options"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "hash() expects between 2 and 4 arguments"}
	}
	algo, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "hash", 1, "algo", "string", bound[0])
	}
	data, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "hash", 2, "data", "string", bound[1])
	}
	if !strings.EqualFold(algo.s, "xxh3") && !strings.EqualFold(algo.s, "xxh128") {
		return value{}, &runtimeError{node.line, "ValueError", "hash(): Argument #1 ($algo) must be a valid hashing algorithm"}
	}
	if present[3] && bound[3].kind == arrayKind && len(bound[3].a.entries) != 0 {
		return value{}, &runtimeError{node.line, "Error", "gophp: xxh3 options are not supported"}
	}
	if strings.EqualFold(algo.s, "xxh128") {
		hi, lo := xxh128([]byte(data.s))
		if present[2] && bound[2].truthy() {
			var raw [16]byte
			binenc.BigEndian.PutUint64(raw[:8], hi)
			binenc.BigEndian.PutUint64(raw[8:], lo)
			return value{kind: stringKind, s: string(raw[:])}, nil
		}
		return value{kind: stringKind, s: fmt.Sprintf("%016x%016x", hi, lo)}, nil
	}
	digest := xxh3Digest([]byte(data.s))
	if present[2] && bound[2].truthy() {
		var raw [8]byte
		binenc.BigEndian.PutUint64(raw[:], digest)
		return value{kind: stringKind, s: string(raw[:])}, nil
	}
	return value{kind: stringKind, s: fmt.Sprintf("%016x", digest)}, nil
}
