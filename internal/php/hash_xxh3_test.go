package php

import (
	"fmt"
	"testing"
)

func TestXXH128KnownDigests(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"", "99aa06d3014798d86001c324468d497f"},
		{"a", "a96faf705af16834e6c632b61e964e1f"},
		{"abc", "06b05ab6733a618578af5f94892f3950"},
		{"Hello, World!", "531df2844447dd5077db03842cd75395"},
		{"The quick brown fox jumps over the lazy dog", "ddd650205ca3e7fa24a1cc2e3a8a7651"},
	} {
		hi, lo := xxh128([]byte(tc.input))
		if got := fmt.Sprintf("%016x%016x", hi, lo); got != tc.want {
			t.Errorf("xxh128(%q) = %s, want %s", tc.input, got, tc.want)
		}
	}
}

// The input generator and expected digests come from xxHash's
// cli/xsum_sanity_check.c (unseeded XXH3_64bits vectors).
func TestXXH3LongReferenceVectors(t *testing.T) {
	input := make([]byte, 2367)
	generator := uint64(2654435761)
	for i := range input {
		input[i] = byte(generator >> 56)
		generator *= 11400714785074694797
	}
	for _, tc := range []struct {
		length int
		want   uint64
	}{
		{403, 0xCDEB804D65C6DEA4},
		{512, 0x617E49599013CB6B},
		{2048, 0xDD59E2C3A5F038E0},
		{2099, 0xC6B9D9B3FC9AC765},
		{2240, 0x6E73A90539CF2948},
		{2367, 0xCB37AEB9E5D361ED},
	} {
		if got := xxh3Digest(input[:tc.length]); got != tc.want {
			t.Errorf("xxh3 len=%d got=%016x want=%016x", tc.length, got, tc.want)
		}
	}
}
