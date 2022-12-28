// Package id provides encoding functions for id generation, nonces, and such.
package id

import (
	"bytes"
	crnd "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sempernow/kit/convert"
	"github.com/sempernow/kit/timestamp"

	"github.com/pkg/errors"

	"github.com/gofrs/uuid"
	//uuip "github.com/pborman/uuid"

	"github.com/oklog/ulid"
	"lukechampine.com/blake3"
)

// Base32 Alphabets : https://en.wikipedia.org/wiki/Base32
const (
	Base32Hex = "0123456789ABCDEFGHIJKLMNOPQRSTUV" // All upper case; indistinguishables (0-O, I-L-1).
	RFC4648   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567" // All upper case; more distinquishable; sans 0(O), 1(I), 8(B).
	WordSafe  = "23456789CFGHJMPQRVWXcfghjmpqrvwx" // No (bad) words.
	Zbase32   = "ybndrfg8ejkmcpqxot1uwisza345h769" // Easy to read.
)

// Base32 returns a base32-encoded UUID v4 sans padding;
// a unique string of 26 alphanumeric characters based on alphabet.
//
//	Alphabets: Base32Hex, RFC4648 (default), WordSafe, Zbase32.
func Base32(alphabet ...string) string {
	var a string
	if len(alphabet) == 0 {
		a = RFC4648
	} else {
		a = alphabet[0]
	}
	var encoding = base32.NewEncoding(a)
	var bb bytes.Buffer
	encoder := base32.NewEncoder(encoding, &bb)
	x := uuid.Must(uuid.NewV4())
	encoder.Write(x.Bytes())
	encoder.Close()
	bb.Truncate(26)

	return bb.String()
}

// UUIDv5 returns the static ID per namespace and name.
// Namespace may be any UUID string, or named reference to one predifined;
// "DNS", "OID", "URL", OR "X500".  Both parameters are case-insenstive.
func UUIDv5(space, name string) (string, error) {

	if name == "" {
		return "", errors.New("missing NAME")
	}
	var ns = uuid.UUID{}

	switch strings.ToUpper(space) {
	case "DNS":
		ns = uuid.NamespaceDNS
	case "OID":
		ns = uuid.NamespaceOID
	case "URL":
		ns = uuid.NamespaceURL
	case "X500":
		ns = uuid.NamespaceX500
	default:
		var err error
		ns, err = uuid.FromString(space)
		if err != nil {
			return "", errors.New(
				"invalid NAMESPACE ['dns', 'oid', 'url', 'x500', (any UUID)]",
			)
		}
	}
	return uuid.NewV5(ns, strings.ToLower(name)).String(), nil
}

// SumBlake3_256 of  string, int, ... to hex string
func SumBlake3_256(i interface{}) string {
	xx := [32]byte{}
	switch x := i.(type) {
	case []byte:
		xx = blake3.Sum256(x)
	case string:
		xx = blake3.Sum256([]byte(x))
	case int:
		xx = blake3.Sum256([]byte(strconv.Itoa(x)))
	case int64:
		xx = blake3.Sum256(convert.Int64ToBytes(x))
	case uint64:
		xx = blake3.Sum256([]byte(convert.Uint64ToString(x)))
	default:
		return convert.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// SumBlake3_512 of  string, int, ... to hex string
func SumBlake3_512(i interface{}) string {
	xx := [64]byte{}
	switch x := i.(type) {
	case []byte:
		xx = blake3.Sum512(x)
	case string:
		xx = blake3.Sum512([]byte(x))
	case int:
		xx = blake3.Sum512([]byte(strconv.Itoa(x)))
	case int64:
		xx = blake3.Sum512(convert.Int64ToBytes(x))
	case uint64:
		xx = blake3.Sum512([]byte(convert.Uint64ToString(x)))
	default:
		return convert.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// SumSHA1 of *os.File, string, int, ... to hex string
func SumSHA1(i interface{}) string {
	xx := [20]byte{}
	switch x := i.(type) {
	case []byte:
		xx = sha1.Sum(x)
	case string:
		xx = sha1.Sum([]byte(x))
	case int:
		xx = sha1.Sum([]byte(strconv.Itoa(x)))
	case int64:
		xx = sha1.Sum(convert.Int64ToBytes(x))
	case uint64:
		xx = sha1.Sum([]byte(convert.Uint64ToString(x)))
	case *os.File:
		h := sha1.New()
		if _, err := io.Copy(h, x); err != nil {
			return convert.BytesToString(xx[:])
		}
		return fmt.Sprintf("%x", h.Sum(nil))
	default:
		return convert.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// SumSHA256 of *os.File, string, int, ... to hex string
func SumSHA256(i interface{}) string {
	xx := [32]byte{}
	switch x := i.(type) {
	case []byte:
		xx = sha256.Sum256(x)
	case string:
		xx = sha256.Sum256([]byte(x))
	case int:
		xx = sha256.Sum256([]byte(strconv.Itoa(x)))
	case int64:
		xx = sha256.Sum256(convert.Int64ToBytes(x))
	case uint64:
		xx = sha256.Sum256([]byte(convert.Uint64ToString(x)))
	case *os.File:
		h := sha256.New()
		if _, err := io.Copy(h, x); err != nil {
			return convert.BytesToString(xx[:])
		}
		return fmt.Sprintf("%x", h.Sum(nil))
	default:
		return convert.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// // Etag : NNN-NNN.part1.part2.part3....
// func Etag(cfg Config, parts ...interface{}) string {
// 	etag := cfg.ver // SERVER_VER + "-" + APP_VER
// 	for _, part := range parts {
// 		etag += fmt.Sprintf(".%v", part)
// 	}
// 	return etag
// }

// UUID is 36 runes; 32+ 4 hyphens
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Format
// UUID v4 a.k.a. "Random"

// NowULID generates a ULID having a timestamp of `tt[0]` else `time.Now()`.
func NowULID(tt ...time.Time) ulid.ULID { // 26 runes
	//t := time.Unix(0, 0) // 000000000006AFVGQT5ZYC0GEK
	var t time.Time
	if len(tt) == 0 {
		t = time.Now()
	} else {
		t = tt[0]
	} // Nanosecond precision
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

// DoULID ...
func DoULID() {
	id := NowULID()
	debug(id.String())
	parse, _ := ulid.Parse("01EC5ZTYF6C7CMX079XS81D1G9")
	debug(parse)
	debug(string(id.Entropy()))
	debug(id.Time()) // Msec
	debug(timestamp.NowEpochMsec())
	debug(ulid.Time(id.Time()))
	debug(time.Now())
	debug(NowULID(time.Now()))
}
func debug(msgs interface{}) {
	fmt.Println(msgs)
}

// Nonce of size bytes.
func Nonce(size int) ([]byte, error) {
	bb := make([]byte, size)
	_, err := io.ReadFull(crnd.Reader, bb)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Nonce : source of randomness unavailable")
	}
	return bb, nil
}

// XOR ...
func XOR(input, key string) (output string) {
	var str strings.Builder

	for i := 0; i < len(input); i++ {
		str.WriteString(string(input[i] ^ key[i%len(key)]))
	}
	return str.String()
}

// EncodeSecret per XOR with nonce bisected for use as both key and padding.
func EncodeSecret(s, nonce string) (string, error) {
	n := nonce[0:(len(nonce)/2 - len(nonce)%2)]
	p := nonce[len(n):]
	//delta := (len(n) - len(s) - 1)
	delta := (len(n) - len(s))
	xx := (len(n) - (delta / 2) - (delta % 2))
	pad := p[0:(delta / 2)]
	//bb, err := XORbytes([]byte(s+"|"+pad), []byte(n[0:xx]))
	bb, err := XORbytes([]byte(s+pad), []byte(n[0:xx]))
	if err != nil {
		return "", err
	}
	return convert.BytesToString(bb), nil
}

// DecodeSecret per XOR with its nonce key/pad.
func DecodeSecret(secret, nonce string) (string, error) {
	bb, err := XORbytes([]byte(secret), []byte(nonce[0:len(secret)]))
	if err != nil {
		return "", err
	}
	//return strings.Split(BytesToString(bb), "|")[0], nil
	return convert.BytesToString(bb), nil
}

// XORbytes performs XOR (encrypt/decrypt) of two equal-sized byte slices.
// https://sourcegraph.com/github.com/hashicorp/vault/-/blob/helper/xor/xor.go#L8
func XORbytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("Slice LENGTHs differ : [%d]a vs. [%d]b", len(a), len(b))
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf, nil
}

// XORstrings performs XOR (encrypt/decrypt) on input against a nonce.
// https://kylewbanks.com/blog/xor-encryption-using-go
func XORstrings(input, nonce string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ nonce[i%len(nonce)])
	}
	return output
}
