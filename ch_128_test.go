package city

import (
	"bufio"
	"bytes"
	_ "embed"
	"strconv"
	"strings"
	"testing"
)

//go:embed _testdata/ch128.csv
var ch128data []byte

func TestCH128(t *testing.T) {
	r := bufio.NewScanner(bytes.NewReader(ch128data))
	for r.Scan() {
		elems := strings.Split(r.Text(), ",")

		s := []byte(elems[0])

		lo, _ := strconv.ParseUint(elems[1], 16, 64)
		hi, _ := strconv.ParseUint(elems[2], 16, 64)

		v := CH128(s)
		if lo != v.Low || hi != v.High {
			t.Errorf("mismatch %d", len(s))
		}
	}
}

func BenchmarkCH128(b *testing.B) {
	setup()
	b.ResetTimer()

	b.ReportAllocs()
	b.SetBytes(1024)

	for i := 0; i < b.N; i++ {
		CH128(data[:1024])
	}
}
