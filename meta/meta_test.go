package meta_test

import (
	"strconv"
	"strings"
	"testing"

	"gd9/prj3/kit/convert"
	// FORKED from github.com/satori
)

func BenchmarkStringContains(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		strings.Contains("2q43o46u34pi 58303 4foo86bar l35325l4j", "foo86bar")
	}
} // BenchmarkStringContains-4   	100000000	        10.7 ns/op	       0 B/op	       0 allocs/op

// ================================================
// ToString functions

var i uint64 = 999999999999

func BenchmarkFormatUint(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		strconv.FormatUint(i, 10)
	}
}

func BenchmarkUint64ToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		convert.Uint64ToString(i)
	}
}

func BenchmarkToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		convert.ToString(i)
	} //... 3x SLOWER than other 3 "ToString" funcs.
}

func BenchmarkIntToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		convert.IntToString(int(999999999999))
	}
}

func BenchmarkInt64ToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		convert.Int64ToString(int64(999999999999))
	}
}
