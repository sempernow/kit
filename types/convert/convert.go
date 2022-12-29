// Package convert provides type-conversion functions.
package convert

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"unsafe"
)

// ----------------------------------------------------------------------------
// JSON ...

// TODO: https://github.com/json-iterator/go
// A fast drop-in replacement for encode/json ...
//import jsoniter "github.com/json-iterator/go"
//var json = jsoniter.ConfigCompatibleWithStandardLibrary

// JSONtoStruct ...
func JSONtoStruct(j []byte, ptr interface{}) error {
	return json.Unmarshal(j, ptr)
}

// Stringify any struct to JSON using json.Marxhal(..)
func Stringify(s interface{}) string {
	bb, err := json.Marshal(&s)
	if err != nil {
		return ""
	}
	return BytesToString(bb)
}

// PrettyPrint any struct to JSON using json.MarshalIndent(..)
func PrettyPrint(i interface{}) string {
	bb, _ := json.MarshalIndent(i, "", "\t")
	return BytesToString(bb)
}

// ToJSON converts a struct to a bytes.Buffer
func ToJSON(x interface{}) (*bytes.Buffer, error) {
	j := new(bytes.Buffer) // https://golang.org/pkg/bytes/#Buffer
	err := json.NewEncoder(j).Encode(&x)
	return j, err

	// Returns []byte
	// var j []byte
	// j, err := json.Marshal(x)
	// if err != nil {
	// 	return j, err
	// }
	// return j, nil
}

// ----------------------------------------------------------------------------
// io.Reader TO ...

// ReaderToBytes ...
func ReaderToBytes(r io.Reader) []byte {
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	return b.Bytes()
}

// ReaderToString ...
func ReaderToString(r io.Reader) string {
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	return b.String()
}

// ----------------------------------------------------------------------------
// BYTEs TO ...

// BytesToString @ zero-copy safely https://github.com/golang/go/issues/25484
func BytesToString(bytes []byte) (s string) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = slice.Data
	str.Len = slice.Len
	runtime.KeepAlive(&bytes) // this line is essential.
	return s
}

// BytesToHex ...
func BytesToHex(bb []byte) string {
	return fmt.Sprintf("%x\n", bb)
}

// ----------------------------------------------------------------------------
// NUMERIC TO ...

// IntToString is faster than ToString
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int64ToString ...
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Uint64ToString ...
func Uint64ToString(u uint64) string {
	return strconv.FormatUint(u, 10)
}

// ToString : any
func ToString(x interface{}) string {
	return fmt.Sprintf("%v", x)
} //... 3x SLOWER than others. (See benchmarks @ kit_test.go)

// Int64ToBytes ...
func Int64ToBytes(n int64) []byte {
	cn := make([]byte, 8)
	binary.LittleEndian.PutUint64(cn, uint64(n))
	return cn
}

// ----------------------------------------------------------------------------
// String to ...

// ToInt ...
func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

// ToInt64 ...
func ToInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(0)
	}
	return n
}

// ToUint64 ...
// Else: uid := strconv.FormatUint(u.ID, 10)
// E.g., from BIGINT (db) to uint64 (struct) to string
func ToUint64(s string) uint64 {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return uint64(0)
	}
	return n
}
