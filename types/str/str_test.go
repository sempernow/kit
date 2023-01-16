package str_test

import (
	"fmt"
	"testing"

	"github.com/sempernow/kit/testkit"
	"github.com/sempernow/kit/types/str"
)

// ☩ go test -v -count=1 -run=TestIsAlphaNum ./kit/types/str/...

func TestIsAlphaNum(t *testing.T) {
	t.Log("@ IsAlphaNum(s) ...")
	ss := []string{"abc1234", "abc 123", "abc/123"}
	exp := []bool{true, false, false}
	for i := range ss {
		testkit.LogDiff(t, fmt.Sprintf("'%s' : %v", ss[i], exp[i]),
			str.IsAlphaNum(ss[i]), exp[i],
		)
	}
}
func BenchmarkRandAlphaNum(b *testing.B) {
	len := 32                          // 92ns @ 16, 160ns @ 32
	fmt.Println(str.RandAlphaNum(len)) // gy2RCNho42n21Wh7
	for i := 0; i < b.N; i++ {
		str.RandAlphaNum(len)
	}
} // BenchmarkRandAlphaNum-4         12554822                91.4 ns/op            16 B/op          1 allocs/op
