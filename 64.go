package city

import (
	"encoding/binary"
)

func bswap64(x uint64) uint64 {
	return ((x & 0xff00000000000000) >> 56) |
		((x & 0x00ff000000000000) >> 40) |
		((x & 0x0000ff0000000000) >> 24) |
		((x & 0x000000ff00000000) >> 8) |
		((x & 0x00000000ff000000) << 8) |
		((x & 0x0000000000ff0000) << 24) |
		((x & 0x000000000000ff00) << 40) |
		((x & 0x00000000000000ff) << 56)
}

func fetch64(p []byte) uint64 {
	r := binary.LittleEndian.Uint64(p)
	return r
}

// Bitwise right rotate
func rotate(val uint64, shift uint) uint64 {
	// Avoid shifting by 64: doing so yields an undefined result.
	if shift == 0 {
		return val
	}
	return (val >> shift) | val<<(64-shift)
}

func shiftMix(val uint64) uint64 {
	return val ^ (val >> 47)
}

func hash128to64(x U128) uint64 {
	const mul = uint64(0x9ddfea08eb382d69)
	a := (x.Low ^ x.High) * mul
	a ^= a >> 47
	b := (x.High ^ a) * mul
	b ^= b >> 47
	b *= mul
	return b
}

func hashLen16(u, v uint64) uint64 {
	return hash128to64(U128{u, v})
}

func hashLen16mul(u, v, mul uint64) uint64 {
	// Murmur-inspired hashing.
	a := (u ^ v) * mul
	a ^= a >> 47
	b := (v ^ a) * mul
	b ^= b >> 47
	b *= mul
	return b
}

func hashLen0to16(s []byte, length int) uint64 {
	if length >= 8 {
		mul := k2 + uint64(length)*2
		a := fetch64(s) + k2
		b := fetch64(s[length-8:])
		c := rotate(b, 37)*mul + a
		d := (rotate(a, 25) + b) * mul
		return hashLen16mul(c, d, mul)
	}
	if length >= 4 {
		mul := k2 + uint64(length)*2
		a := uint64(fetch32(s))
		first := uint64(length) + (a << 3)
		second := uint64(fetch32(s[length-4:]))
		result := hashLen16mul(
			first,
			second,
			mul)
		return result
	}
	if length > 0 {
		a := s[0]
		b := s[length>>1]
		c := s[length-1]
		y := uint32(a) + (uint32(b) << 8)
		z := uint32(length) + (uint32(c) << 2)
		return shiftMix(uint64(y)*k2^uint64(z)*k0) * k2
	}
	return k2
}

// This probably works well for 16-byte strings as well, but is may be overkill
// in that case
func hashLen17to32(s []byte, length int) uint64 {
	mul := k2 + uint64(length)*2
	a := fetch64(s) * k1
	b := fetch64(s[8:])
	c := fetch64(s[length-8:]) * mul
	d := fetch64(s[length-16:]) * k2
	return hashLen16mul(
		rotate(a+b, 43)+rotate(c, 30)+d,
		a+rotate(b+k2, 18)+c,
		mul)
}

// Return a 16-byte hash for 48 bytes. Quick and dirty.
// callers do best to use "random-looking" values for a and b.
func weakHashLen32WithSeeds(w, x, y, z, a, b uint64) U128 {
	a += w
	b = rotate(b+a+z, 21)
	c := a
	a += x
	a += y
	b += rotate(a, 44)
	return U128{a + z, b + c}
}

// Return a 16-byte hash for s[0] ... s[31], a, and b. Quick and dirty.
func weakHashLen32WithSeedsByte(s []byte, a, b uint64) U128 {
	return weakHashLen32WithSeeds(
		fetch64(s),
		fetch64(s[8:]),
		fetch64(s[16:]),
		fetch64(s[24:]),
		a,
		b)
}

// Return an 8-byte hash for 33 to 64 bytes.
func hashLen33to64(s []byte, length int) uint64 {
	mul := k2 + uint64(length)*2
	a := fetch64(s) * k2
	b := fetch64(s[8:])
	c := fetch64(s[length-24:])
	d := fetch64(s[length-32:])
	e := fetch64(s[16:]) * k2
	f := fetch64(s[24:]) * 9
	g := fetch64(s[length-8:])
	h := fetch64(s[length-16:]) * mul
	u := rotate(a+g, 43) + (rotate(b, 30)+c)*9
	v := ((a + g) ^ d) + f + 1
	w := bswap64((u+v)*mul) + h
	x := rotate(e+f, 42) + c
	y := (bswap64((v+w)*mul) + g) * mul
	z := e + f + c
	a = bswap64((x+z)*mul+y) + b
	b = shiftMix((z+a)*mul+d+h) * mul
	return b + x
}

// Hash64 return a 64-bit hash.
func Hash64(s []byte) uint64 {
	length := len(s)
	if length <= 32 {
		if length <= 16 {
			return hashLen0to16(s, length)
		}
		return hashLen17to32(s, length)
	} else if length <= 64 {
		return hashLen33to64(s, length)
	}

	// For string over 64 bytes we hash the end first, and then as we
	// loop we keep 56 bytes of state: v, w, x, y and z.
	x := fetch64(s[length-40:])
	y := fetch64(s[length-16:]) + fetch64(s[length-56:])
	z := hashLen16(fetch64(s[length-48:])+uint64(length), fetch64(s[length-24:]))
	v := weakHashLen32WithSeedsByte(s[length-64:], uint64(length), z)
	w := weakHashLen32WithSeedsByte(s[length-32:], y+k1, x)
	x = x*k1 + fetch64(s)

	// Decrease len to the nearest multiple of 64, and operate on 64-byte chunks.
	tmpLength := uint32(length)
	tmpLength = (tmpLength - 1) & ^uint32(63)
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
		tmpLength -= 64
		if tmpLength == 0 {
			break
		}
	}

	return hashLen16(
		hashLen16(v.Low, w.Low)+shiftMix(y)*k1+z,
		hashLen16(v.High, w.High)+x)
}

// Hash64WithSeed return a 64-bit hash with a seed.
func Hash64WithSeed(s []byte, seed uint64) uint64 {
	return Hash64WithSeeds(s, k2, seed)
}

// Hash64WithSeeds return a 64-bit hash with two seeds.
func Hash64WithSeeds(s []byte, seed0, seed1 uint64) uint64 {
	return hashLen16(Hash64(s)-seed0, seed1)
}
