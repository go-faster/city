package city

import (
	"encoding/binary"
	"fmt"
)

// U128 is uint128.
type U128 struct {
	Low  uint64 // first 32 bits
	High uint64 // last 32 bits

	// much faster than uint64[2]
}

// Arr returns byte array that represents uint128 value of U128 in big endian.
func (u U128) Arr() (v [16]byte) {
	binary.BigEndian.PutUint64(v[:8], u.Low)
	binary.BigEndian.PutUint64(v[8:], u.High)
	return v
}

// Append appends uint128 big value to buf in big endian.
//
// Append is much faster that appending Arr result.
func (u U128) Append(buf []byte) []byte {
	s := len(buf)
	buf = append(buf, make([]byte, 16)...)
	binary.BigEndian.PutUint64(buf[s:], u.Low)
	binary.BigEndian.PutUint64(buf[s+8:], u.High)
	return buf
}

func (u U128) String() string {
	return fmt.Sprintf("%x", u.Arr())
}

// A subroutine for Hash128(). Returns a decent 128-bit hash for strings
// of any length representable in signed long. Based on City and Mumur.
func cityMurmur(s []byte, length int, seed U128) U128 {
	a := seed.Low
	b := seed.High
	c := uint64(0)
	d := uint64(0)
	l := length - 16
	if l <= 0 { // length <= 16
		a = shiftMix(a*k1) * k1
		c = b*k1 + hash0to16(s, length)

		tmp := c
		if length >= 8 {
			tmp = u64(s)
		}
		d = shiftMix(a + tmp)
	} else { // length > 16
		c = hash16(u64(s[length-8:])+k1, a)
		d = hash16(b+uint64(length), c+u64(s[length-16:]))
		a += d
		for {
			a ^= shiftMix(u64(s)*k1) * k1
			a *= k1
			b ^= a
			c ^= shiftMix(u64(s[8:])*k1) * k1
			c *= k1
			d ^= c
			s = s[16:]
			l -= 16
			if l <= 0 {
				break
			}
		}
	}
	a = hash16(a, c)
	b = hash16(d, b)
	return U128{a ^ b, hash16(b, a)}
}

// Hash128Seed return a 128-bit hash with a seed.
func Hash128Seed(s []byte, seed U128) U128 {
	length := len(s)
	if length < 128 {
		return cityMurmur(s, length, seed)
	}

	savedLength := length
	savedS := s

	// We expect len >= 128 to be the common case. Keep 56 bytes of state:
	// v, w, x, y and z.
	var v, w U128
	x := seed.Low
	y := seed.High
	z := uint64(length) * k1

	v.Low = rot64(y^k1, 49)*k1 + u64(s)
	v.High = rot64(v.Low, 42)*k1 + u64(s[8:])
	w.Low = rot64(y+z, 35)*k1 + x
	w.High = rot64(x+u64(s[88:]), 53) * k1

	// This is the same inner loop as Hash64(), manually unrolled.
	for {
		x = rot64(x+y+v.Low+u64(s[8:]), 37) * k1
		y = rot64(y+v.High+u64(s[48:]), 42) * k1
		x ^= w.High
		y += v.Low + u64(s[40:])
		z = rot64(z+w.Low, 33) * k1
		v = weakHash32SeedsByte(s, v.High*k1, x+w.Low)
		w = weakHash32SeedsByte(s[32:], z+w.High, y+u64(s[16:]))
		z, x = x, z
		s = s[64:]
		x = rot64(x+y+v.Low+u64(s[8:]), 37) * k1
		y = rot64(y+v.High+u64(s[48:]), 42) * k1
		x ^= w.High
		y += v.Low + u64(s[40:])
		z = rot64(z+w.Low, 33) * k1
		v = weakHash32SeedsByte(s, v.High*k1, x+w.Low)
		w = weakHash32SeedsByte(s[32:], z+w.High, y+u64(s[16:]))
		z, x = x, z
		s = s[64:]
		length -= 128
		if length < 128 {
			break
		}
	}
	x += rot64(v.Low+z, 49) * k0
	y = y*k0 + rot64(w.High, 37)
	z = z*k0 + rot64(w.Low, 27)
	w.Low *= 9
	v.Low *= k0
	// If 0 < length < 128, hash up to 4 chunks of 32 bytes each form the end
	// of s.
	for tailDone := 0; tailDone < length; {
		tailDone += 32
		y = rot64(x+y, 42)*k0 + v.High
		w.Low += u64(savedS[savedLength-tailDone+16:])
		x = x*k0 + w.Low
		z += w.High + u64(savedS[savedLength-tailDone:])
		w.High += v.Low
		v = weakHash32SeedsByte(savedS[savedLength-tailDone:], v.Low+z, v.High)
		v.Low *= k0
	}

	// At this point our 56 bytes of state should contain more than
	// enough information for a strong 128-bit hash. We use two different
	// 56-byte-to-8-byte hashes to get a 16-byte final result.
	x = hash16(x, v.Low)
	y = hash16(y+z, w.Low)
	return U128{hash16(x+v.High, w.High) + y, hash16(x+w.High, y+v.High)}
}

// Hash128 returns a 128-bit hash and are tuned for strings of at least
// a few hundred bytes.  Depending on your compiler and hardware,
// it's likely faster than Hash64() on sufficiently long strings.
// It's slower than necessary on shorter strings, but we expect
// that case to be relatively unimportant.
func Hash128(s []byte) U128 {
	length := len(s)
	if length >= 16 {
		return Hash128Seed(s[16:length], U128{
			u64(s),
			u64(s[8:]) + k0})
	}
	return Hash128Seed(s, U128{k0, k1})
}
