package city

// U128 is uint128.
type U128 struct {
	Low, High uint64 // much faster than uint64[2]
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
		c = b*k1 + hashLen0to16(s, length)

		tmp := c
		if length >= 8 {
			tmp = fetch64(s)
		}
		d = shiftMix(a + tmp)
	} else { // length > 16
		c = hashLen16(fetch64(s[length-8:])+k1, a)
		d = hashLen16(b+uint64(length), c+fetch64(s[length-16:]))
		a += d
		for {
			a ^= shiftMix(fetch64(s)*k1) * k1
			a *= k1
			b ^= a
			c ^= shiftMix(fetch64(s[8:])*k1) * k1
			c *= k1
			d ^= c
			s = s[16:]
			l -= 16
			if l <= 0 {
				break
			}
		}
	}
	a = hashLen16(a, c)
	b = hashLen16(d, b)
	return U128{a ^ b, hashLen16(b, a)}
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

	v.Low = rotate(y^k1, 49)*k1 + fetch64(s)
	v.High = rotate(v.Low, 42)*k1 + fetch64(s[8:])
	w.Low = rotate(y+z, 35)*k1 + x
	w.High = rotate(x+fetch64(s[88:]), 53) * k1

	// This is the same inner loop as Hash64(), manually unrolled.
	for {
		x = rotate(x+y+v.Low+fetch64(s[8:]), 37) * k1
		y = rotate(y+v.High+fetch64(s[48:]), 42) * k1
		x ^= w.High
		y += v.Low + fetch64(s[40:])
		z = rotate(z+w.Low, 33) * k1
		v = weakHashLen32WithSeedsByte(s, v.High*k1, x+w.Low)
		w = weakHashLen32WithSeedsByte(s[32:], z+w.High, y+fetch64(s[16:]))
		z, x = x, z
		s = s[64:]
		x = rotate(x+y+v.Low+fetch64(s[8:]), 37) * k1
		y = rotate(y+v.High+fetch64(s[48:]), 42) * k1
		x ^= w.High
		y += v.Low + fetch64(s[40:])
		z = rotate(z+w.Low, 33) * k1
		v = weakHashLen32WithSeedsByte(s, v.High*k1, x+w.Low)
		w = weakHashLen32WithSeedsByte(s[32:], z+w.High, y+fetch64(s[16:]))
		z, x = x, z
		s = s[64:]
		length -= 128
		if length < 128 {
			break
		}
	}
	x += rotate(v.Low+z, 49) * k0
	y = y*k0 + rotate(w.High, 37)
	z = z*k0 + rotate(w.Low, 27)
	w.Low *= 9
	v.Low *= k0
	// If 0 < length < 128, hash up to 4 chunks of 32 bytes each form the end
	// of s.
	for tailDone := 0; tailDone < length; {
		tailDone += 32
		y = rotate(x+y, 42)*k0 + v.High
		w.Low += fetch64(savedS[savedLength-tailDone+16:])
		x = x*k0 + w.Low
		z += w.High + fetch64(savedS[savedLength-tailDone:])
		w.High += v.Low
		v = weakHashLen32WithSeedsByte(savedS[savedLength-tailDone:], v.Low+z, v.High)
		v.Low *= k0
	}

	// At this point our 56 bytes of state should contain more than
	// enough information for a strong 128-bit hash. We use two different
	// 56-byte-to-8-byte hashes to get a 16-byte final result.
	x = hashLen16(x, v.Low)
	y = hashLen16(y+z, w.Low)
	return U128{hashLen16(x+v.High, w.High) + y, hashLen16(x+w.High, y+v.High)}
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
			fetch64(s),
			fetch64(s[8:]) + k0})
	}
	return Hash128Seed(s, U128{k0, k1})
}
