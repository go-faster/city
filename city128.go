package cityhash

// A subroutine for CityHash128(). Returns a decent 128-bit hash for strings
// of any length representable in signed long. Based on City and Mumur.
func cityMurmur(s []byte, length int, seed Uint128) Uint128 {
	a := seed.Low64()
	b := seed.High64()
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
	return Uint128{a ^ b, hashLen16(b, a)}
}

// CityHash128WithSeed return a 128-bit hash with a seed.
func CityHash128WithSeed(s []byte, seed Uint128) Uint128 {
	length := len(s)
	if length < 128 {
		return cityMurmur(s, length, seed)
	}

	savedLength := length
	savedS := s

	// We expect len >= 128 to be the common case. Keep 56 bytes of state:
	// v, w, x, y and z.
	var v, w Uint128
	x := seed.Low64()
	y := seed.High64()
	z := uint64(length) * k1

	v[0] = rotate(y^k1, 49)*k1 + fetch64(s)
	v[1] = rotate(v[0], 42)*k1 + fetch64(s[8:])
	w[0] = rotate(y+z, 35)*k1 + x
	w[1] = rotate(x+fetch64(s[88:]), 53) * k1

	// This is the same inner loop as CityHash64(), manually unrolled.
	for {
		x = rotate(x+y+v[0]+fetch64(s[8:]), 37) * k1
		y = rotate(y+v[1]+fetch64(s[48:]), 42) * k1
		x ^= w[1]
		y += v[0] + fetch64(s[40:])
		z = rotate(z+w[0], 33) * k1
		v = weakHashLen32WithSeedsByte(s, v[1]*k1, x+w[0])
		w = weakHashLen32WithSeedsByte(s[32:], z+w[1], y+fetch64(s[16:]))
		z, x = x, z
		s = s[64:]
		x = rotate(x+y+v[0]+fetch64(s[8:]), 37) * k1
		y = rotate(y+v[1]+fetch64(s[48:]), 42) * k1
		x ^= w[1]
		y += v[0] + fetch64(s[40:])
		z = rotate(z+w[0], 33) * k1
		v = weakHashLen32WithSeedsByte(s, v[1]*k1, x+w[0])
		w = weakHashLen32WithSeedsByte(s[32:], z+w[1], y+fetch64(s[16:]))
		z, x = x, z
		s = s[64:]
		length -= 128
		if length < 128 {
			break
		}
	}
	x += rotate(v[0]+z, 49) * k0
	y = y*k0 + rotate(w[1], 37)
	z = z*k0 + rotate(w[0], 27)
	w[0] *= 9
	v[0] *= k0
	// If 0 < length < 128, hash up to 4 chunks of 32 bytes each form the end
	// of s.
	for tailDone := 0; tailDone < length; {
		tailDone += 32
		y = rotate(x+y, 42)*k0 + v.High64()
		w[0] += fetch64(savedS[savedLength-tailDone+16:])
		x = x*k0 + w.Low64()
		z += w.High64() + fetch64(savedS[savedLength-tailDone:])
		w[1] += v.Low64()
		v = weakHashLen32WithSeedsByte(savedS[savedLength-tailDone:], v.Low64()+z, v.High64())
		v[0] *= k0
	}

	// At this point our 56 bytes of state should contain more than
	// enough information for a strong 128-bit hash. We use two different
	// 56-byte-to-8-byte hashes to get a 16-byte final result.
	x = hashLen16(x, v.Low64())
	y = hashLen16(y+z, w.Low64())
	return Uint128{hashLen16(x+v[1], w[1]) + y, hashLen16(x+w[1], y+v[1])}
}

// CityHash128 return a 128-bit hash and are tuned for strings of at least
// a few hundred bytes.  Depending on your compiler and hardware,
// it's likely faster than CityHash64() on sufficiently long strings.
// It's slower than necessary on shorter strings, but we expect
// that case to be relatively unimportant.
func CityHash128(s []byte) Uint128 {
	length := len(s)
	if length >= 16 {
		return CityHash128WithSeed(s[16:length], Uint128{
			fetch64(s),
			fetch64(s[8:]) + k0})
	}
	return CityHash128WithSeed(s, Uint128{k0, k1})
}
