// Package gz provides gzip read/write (de)compression of byte slices.
package gz

// REF: gzip pkg @ stdlib: https://golang.org/pkg/compress/gzip
import (
	"bytes"
	"compress/gzip"

	"io"
	"time"
	//"github.com/klauspost/compress/gzip" // WARNING: Default compression is NO COMPRESSION
)

// An example use case for `Write(..)` is at a static-file server
// that utilizes an embedded cache to store its requested resources.
// `Read(..)` would be useful there as a unit-testing helper.

// Write compresses a byte slice per gzip.
func Write(bb []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	zw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	// zw.Name = "@APP"
	// zw.Comment = "kit"
	zw.ModTime = time.Now().UTC()

	_, err = zw.Write(bb)
	if err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Read decompresses a byte slice per gzip.
func Read(bb []byte) ([]byte, error) {
	br := bytes.NewReader(bb)
	zr, err := gzip.NewReader(br)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, zr); err != nil {
		return nil, err
	}
	if err := zr.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// ============================================================================
// 	PRIOR EFFORTs
// ============================================================================

// GzipCompress bb into an io.Reader
func GzipCompress(bb []byte) (io.Reader, error) {
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	defer func() {
		gz.Flush()
		gz.Close()
	}()
	_, err := gz.Write(bb)
	return buf, err
}

// GzipDecompress from io.Reader to bb
func GzipDecompress(gzd io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(gzd)
	chk(err)
	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, gz)
	return buf.Bytes(), err
}

// GzipCompressToBytes ...
func GzipCompressToBytes(bb []byte) ([]byte, error) {
	var err error
	compressed, err := GzipCompress(bb)
	chk(err)

	b := new(bytes.Buffer)
	_, err = b.ReadFrom(compressed)
	return b.Bytes(), err
}

// GzipDecompressFromBytes ...
func GzipDecompressFromBytes(gzb []byte) ([]byte, error) {

	gzbuf := new(bytes.Buffer)
	gz, err := gzip.NewReader(gzbuf)
	chk(err)
	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, gz) // decompress and read into buffer

	return buf.Bytes(), err

	// gz := func() io.Reader {
	// 	gzbuf := new(bytes.Buffer)
	// 	_, err := gzbuf.Read(gzb)
	// 	chk(err)
	// 	gz, err := gzip.NewReader(gzbuf)
	// 	chk(err)
	// 	defer func() {
	// 		gz.Close()
	// 	}()
	// 	return gz
	// }()

	// // buf := new(bytes.Buffer)
	// // _, err := buf.ReadFrom(gz)
	// return ioutil.ReadAll(gz)
}

// ToBytes FAILs ... nil pointer
// Must do all at once (compress and to bytes)
func ToBytes(gzr io.Reader) ([]byte, error) {
	// Fail 1
	//return ioutil.ReadAll(gzr)

	// Fail 2
	// buf := new(bytes.Buffer)
	// _, err := io.Copy(buf, gzr)
	// return buf.Bytes(), err

	// Fail 3
	b := new(bytes.Buffer)
	_, err := b.ReadFrom(gzr)
	//_, err := b.ReadFrom(bufio.NewReader(gzr))
	return b.Bytes(), err
}

func chk(e error) error {
	if e != nil {
		return e
	}
	return nil
}
