package id_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"testing"

	"gd9/prj3/kit/id"
	"gd9/prj3/kit/str"

	"github.com/gofrs/uuid" // FORKED from github.com/satori
	"lukechampine.com/blake3"
)

func TestOBA(t *testing.T) {
	//id.XOR("", "")
	//str.Reverse("")
	t.Skip()
	var u uuid.UUID
	u, _ = uuid.NewV4()
	_na := sha512.Sum512(u.Bytes())
	u, _ = uuid.NewV4()
	_nb := sha512.Sum512(u.Bytes())

	na := hex.EncodeToString(_na[:])
	nb := hex.EncodeToString(_nb[:])

	// t.Log(na) // 128 ch
	// n, _ := web.Nonce(100)
	// t.Log(n)
	// t.Log(id.Base32UUID())
	// return

	email := "foo@bar.com"
	pass := "1234"

	email = id.XOR(email, str.Reverse(na))
	pass = id.XOR(pass, id.XOR(na, str.Reverse(nb)))

	t.Log("email", hex.EncodeToString([]byte(email)))
	t.Log("pass", hex.EncodeToString([]byte(pass)))

	email = id.XOR(email, str.Reverse(na))
	pass = id.XOR(pass, id.XOR(na, str.Reverse(nb)))

	t.Log(email, pass)

	//t.Fail()
}

// ******************************************************
// go test -count=1 -benchmem -bench=EncodeSecret ./kit
// ******************************************************

func BenchmarkEncodeSecret(b *testing.B) {
	// Run it b.N times
	secret := id.SumSHA256("a")
	nonce := id.SumSHA256("b")
	//x, _ := id.DecodeSecret(secret[:], nonce)

	bb, _ := id.XORbytes([]byte(secret), []byte(nonce[0:len(secret)]))
	fmt.Println("hex:", hex.EncodeToString(bb))

	for n := 0; n < b.N; n++ {
		bb, _ := id.XORbytes([]byte(secret), []byte(nonce[0:len(secret)]))
		hex.EncodeToString(bb)
	}
} // BenchmarkEncodeSecret-4   	 5881406	       194 ns/op	     192 B/op	       3 allocs/op

func BenchmarkOBA(b *testing.B) {
	var u uuid.UUID
	u, _ = uuid.NewV4()
	_na := sha512.Sum512(u.Bytes())
	u, _ = uuid.NewV4()
	_nb := sha512.Sum512(u.Bytes())

	na := hex.EncodeToString(_na[:])
	nb := hex.EncodeToString(_nb[:])

	email := "foo@bar.com"
	pass := "1234"

	email = id.XOR(email, str.Reverse(na))
	pass = id.XOR(pass, id.XOR(na, str.Reverse(nb)))

	for n := 0; n < b.N; n++ {
		email = id.XOR(email, str.Reverse(na))
		pass = id.XOR(pass, id.XOR(na, str.Reverse(nb)))
	}

} // BenchmarkOBA-4   	   83332	     14496 ns/op	    1592 B/op	      12 allocs/op

func BenchmarkXORbytes(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	u, _ = uuid.NewV4()
	y := u.Bytes()
	for n := 0; n < b.N; n++ {
		bb, _ := id.XORbytes(x, y)
		hex.EncodeToString(bb)
	}
} //BenchmarkXORbytes-4   	18176034	        63.4 ns/op	      48 B/op	       1 allocs/op
// BenchmarkXORbytes-4   	 5940584	       219 ns/op	     208 B/op	       3 allocs/op

func BenchmarkXORbytesOnStrings(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	u, _ = uuid.NewV4()
	y := u.Bytes()
	bb, _ := id.XORbytes(x, y)
	fmt.Println(hex.EncodeToString(bb))

	for n := 0; n < b.N; n++ {
		//id.XORbytes([]byte(x), []byte(y))
		bb, _ := id.XORbytes([]byte(x), []byte(y))
		hex.EncodeToString(bb)
	}
} // BenchmarkXORbytesOnStrings-4   	 3680936	       313 ns/op	     304 B/op	       5 allocs/op

func BenchmarkXOR(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.String()
	u, _ = uuid.NewV4()
	y := u.String()
	for n := 0; n < b.N; n++ {
		id.XOR(x, y)
	}
} // go test -bench=XOR -count=1 ./kit
// BenchmarkXOR-4   	 1300089	       919 ns/op	     120 B/op	       4 allocs/op

func BenchmarkXORstrings(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.String()
	u, _ = uuid.NewV4()
	y := u.String()
	for n := 0; n < b.N; n++ {
		id.XORstrings(x, y)
	}
} //BenchmarkXORstrings-4   	  461542	      2383 ns/op	     864 B/op	      35 allocs/op

func TestSumSHA1(t *testing.T) {
	t.Skip()
	t.Log(id.SumSHA1([]byte("1594562637"))) //c2db288bfaeaf238a793d7a7f9c51ab50b2b9196
	t.Log(id.SumSHA1("1594562637"))         //c2db288bfaeaf238a793d7a7f9c51ab50b2b9196
	t.Log(id.SumSHA1(1594562637))           //c2db288bfaeaf238a793d7a7f9c51ab50b2b9196
	t.Log(id.SumSHA1("☩"))                  //ec7b5357c265c5832ad34aaa6eb4305b3bd83bc5
	t.Log(id.SumSHA1('☩'))                  //0000000000000000000000000000000000000000

	f, err := os.Open("kit_test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	t.Log(id.SumSHA1(f)) //0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33
}

func TestULID(t *testing.T) {
	t.Skip()
	id.DoULID()
	//t.Fatal()
}

func BenchmarkULID(b *testing.B) {
	// Run it b.N times
	//fmt.Println(id.NowULID().String())
	for n := 0; n < b.N; n++ {
		id.NowULID() //.String()
	}
} // BenchmarkULID-4   	   63837	     18485 ns/op	    9696 B/op	       6 allocs/p

func BenchmarkBase32(b *testing.B) {
	// Run it b.N times
	//fmt.Println(id.Base32())
	//fmt.Println(id.Base32(id.Zbase32))
	//fmt.Println(id.Base32(id.WordSafe))
	//fmt.Println(id.Base32(id.Base32Hex))
	for n := 0; n < b.N; n++ {
		id.Base32() //.String()
	}
} // BenchmarkBase32-4   	  857222	      1199 ns/op	    1648 B/op	       7 allocs/op

func BenchmarkSHA256(b *testing.B) {
	// Run it b.N times
	//s := sha256.Sum256([]byte("foobar"))
	for n := 0; n < b.N; n++ {
		id.SumSHA256([]byte("foobar"))
		// if s == sha256.Sum256([]byte("foobar")) {
		// 	//true
		// } else {
		// 	//false
		// }
	}
} // BenchmarkSHA256-4   	 2368640	       500 ns/op	     168 B/op	       4 allocs/op
//   BenchmarkSHA256-4       4286598           271 ns/op           0 B/op          0 allocs/op @ sha256 + test

func BenchmarkBlake3256(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	for n := 0; n < b.N; n++ {
		blake3.Sum256(x)
	}
} // BenchmarkBlake3256-4     7118822        162 ns/op               0 B/op          0 allocs/op
func BenchmarkSumBlake3_256(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.String()
	for n := 0; n < b.N; n++ {
		//id.SumBlake3_256(bb)
		id.SumBlake3_256(x)
	}
} // BenchmarkSumBlake3_256-4    3558349	       337 ns/op	     160 B/op	       3 allocs/op  @ bb
// BenchmarkSumBlake3_256-4      2787312	       387 ns/op	     208 B/op	       4 allocs/op  @ s

func BenchmarkSHA512(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	for n := 0; n < b.N; n++ {
		sha512.Sum512(x)
	}
} // BenchmarkSHA512-4   	 3251862	       364 ns/op	       0 B/op	       0 allocs/op

// https://pkg.go.dev/lukechampine.com/blake3
func BenchmarkBlake3512(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	for n := 0; n < b.N; n++ {
		//blake3.Sum512(bb)
		b := blake3.Sum512(x)
		hex.EncodeToString(b[:])
		//id.SumBlake3_512(bb)
	}
} // BenchmarkBlake3512-4   7205895          164 ns/op               0 B/op          0 allocs/op
//   BenchmarkBlake3512-4   2819072          417 ns/op             320 B/op          3 allocs/op @ string

func BenchmarkNewUUID(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		uuid.NewV4()
	}
} // BenchmarkNewUUID-4   	 3045514	       393 ns/op	      64 B/op	       2 allocs/o

func BenchmarkHMAC(b *testing.B) {
	// Run it b.N times

	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	u, _ = uuid.NewV4()
	y := u.Bytes()

	rid := fmt.Sprintf("%x", sha256.Sum256(x))
	rKey := fmt.Sprintf("%x", sha256.Sum256(y))
	for n := 0; n < b.N; n++ {
		hmac.New(sha256.New, []byte(rKey+rid))
	}
} // BenchmarkHMAC-4   	  399993	      2914 ns/op	    1152 B/op	      10 allocs/op
