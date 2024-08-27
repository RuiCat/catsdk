//go:build clibs
// +build clibs

package sereal

import (
	"boltdb/util/zstd"
)

func zstdEncode(buf []byte, level int) ([]byte, error) {
	dst, err := zstd.CompressLevel(nil, buf, level)
	return dst, err
}

func zstdDecode(d, buf []byte) ([]byte, error) {
	dst, err := zstd.Decompress(d, buf)
	return dst, err
}
