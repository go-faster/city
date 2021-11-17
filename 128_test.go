package city

import (
	"bytes"
	"testing"
)

func TestU128_Append(t *testing.T) {
	v := U128{0x18062081bf558df, 0x63416eb68f104a36}

	a := v.Arr()
	if !bytes.Equal(v.Append(nil), a[:]) {
		t.Error("mismatch")
	}
}

func BenchmarkU128_Append(b *testing.B) {
	b.ReportAllocs()
	v := U128{0x18062081bf558df, 0x63416eb68f104a36}

	var buf []byte

	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = v.Append(buf)
	}
}

func BenchmarkU128_Arr(b *testing.B) {
	b.ReportAllocs()
	v := U128{0x18062081bf558df, 0x63416eb68f104a36}

	for i := 0; i < b.N; i++ {
		a := v.Arr()
		_ = a[15]
	}
}
