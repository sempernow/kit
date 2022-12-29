package str_test

import (
	"fmt"
	"testing"

	"github.com/sempernow/kit/types/str"
)

func BenchmarkRandAlphaNum(b *testing.B) {
	len := 32                          // 92ns @ 16, 160ns @ 32
	fmt.Println(str.RandAlphaNum(len)) // gy2RCNho42n21Wh7
	for i := 0; i < b.N; i++ {
		str.RandAlphaNum(len)
	}
} // BenchmarkRandAlphaNum-4         12554822                91.4 ns/op            16 B/op          1 allocs/op
