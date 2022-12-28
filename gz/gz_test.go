package gz_test

import (
	"math/rand"
	"testing"

	"gd9/prj3/kit/gz"
	"gd9/prj3/kit/testkit"
)

func TestGzipper(t *testing.T) {

	t.Log("@ Write() / Read() : NON-ZERO byte slice")
	{
		bb := []byte("This may be running the klauspost/compress/gzip package ..." + RandStringBytes(44000))
		n := 55
		t.Log("\t@ Orig:   ", len(bb), string(bb)[:n])

		zw, err := gz.Write(bb)
		var c float64
		c = float64(len(bb) - len(zw))
		c = c / float64(len(bb))
		c = float64((c * 100))
		testkit.Log(t, "Comp", err)
		t.Logf("\t@ Comp:    %d (-%.0f%s) %s", len(zw), c, "%", string(zw)[:n])

		zr, err := gz.Read(zw)
		testkit.Log(t, "Decomp", err)
		t.Log("\t@ Decomp: ", len(zr), string(zr)[:n])
	}
	t.Log("@ Write() / Read() : ZERO byte slice")
	{
		bb := []byte{}
		t.Log("\t@ Orig:   ", len(bb), string(bb))

		zw, err := gz.Write(bb)
		testkit.Log(t, "Comp", err)
		t.Log("\t@ Comp:   ", len(zw), string(zw))

		zr, err := gz.Read(zw)
		testkit.Log(t, "Decomp", err)
		t.Log("\t@ Decomp: ", len(zr), string(zr))
	}
}

// @ klauspost
// BenchmarkGzipper-4   	    4520	    251327 ns/op	 1401881 B/op	      15 allocs/op
// @ stdlib
// BenchmarkGzipper-4   	    3428	    298142 ns/op	  822566 B/op	      23 allocs/op
func BenchmarkGzipper(b *testing.B) {
	bb := []byte(RandStringBytes(4096))
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		gz.Write(bb)
	}
}

func RandStringBytes(n int) string {
	const dictionary = " 1234567890=+-_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(b)
}
