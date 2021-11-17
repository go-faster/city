//go:build go1.18
// +build go1.18

package city

import (
	"bytes"
	"testing"
)

func FuzzHash128(f *testing.F) {
	for _, s := range [][]byte{
		nil,
		{},
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6},
		[]byte("hello"),
		[]byte("hello world"),
		bytes.Repeat([]byte("hello"), 100),
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		Hash128(data)
	})
}
