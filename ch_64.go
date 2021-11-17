package city

// Ref:
// https://github.com/xzkostyan/python-cityhash/commit/f4091154ff2c6c0de11d5d6673b5007fdd6355ad

const k3 uint64 = 0xc949d7c7509e6557

func ch16(u, v uint64) uint64 {
	return hash128to64(U128{u, v})
}

// Return an 8-byte hash for 33 to 64 bytes.
func ch33to64(s []byte, length int) uint64 {
	z := fetch64(s[24:])
	a := fetch64(s) + (uint64(length)+fetch64(s[length-16:]))*k0
	b := rotate(a+z, 52)
	c := rotate(a, 37)

	a += fetch64(s[8:])
	c += rotate(a, 7)
	a += fetch64(s[16:])

	vf := a + z
	vs := b + rotate(a, 31) + c

	a = fetch64(s[16:]) + fetch64(s[length-32:])
	z = fetch64(s[length-8:])
	b = rotate(a+z, 52)
	c = rotate(a, 37)
	a += fetch64(s[length-24:])
	c += rotate(a, 7)
	a += fetch64(s[length-16:])

	wf := a + z
	ws := b + rotate(a, 31) + c
	r := shiftMix((vf+ws)*k2 + (wf+vs)*k0)
	return shiftMix(r*k0+vs) * k2
}

func ch17to32(s []byte, length int) uint64 {
	a := fetch64(s) * k1
	b := fetch64(s[8:])
	c := fetch64(s[length-8:]) * k2
	d := fetch64(s[length-16:]) * k0
	return hashLen16(
		rotate(a-b, 43)+rotate(c, 30)+d,
		a+rotate(b^k3, 20)-c+uint64(length),
	)
}

func ch0to16(s []byte, length int) uint64 {
	if length > 8 {
		a := fetch64(s)
		b := fetch64(s[length-8:])
		return ch16(a, rotatePositive(b+uint64(length), length)) ^ b
	}
	if length >= 4 {
		a := uint64(fetch32(s))
		return ch16(uint64(length)+(a<<3), uint64(fetch32(s[length-4:])))
	}
	if length > 0 {
		a := s[0]
		b := s[length>>1]
		c := s[length-1]
		y := uint32(a) + (uint32(b) << 8)
		z := uint32(length) + (uint32(c) << 2)
		return shiftMix(uint64(y)*k2^uint64(z)*k3) * k2
	}
	return k2
}

// CH64 returns ClickHouse version of Hash64.
func CH64(s []byte) uint64 {
	length := len(s)
	if length <= 16 {
		return ch0to16(s, length)
	}
	if length <= 32 {
		return ch17to32(s, length)
	}
	if length <= 64 {
		return ch33to64(s, length)
	}

	x := fetch64(s)
	y := fetch64(s[length-16:]) ^ k1
	z := fetch64(s[length-56:]) ^ k0

	v := weakHashLen32WithSeedsByte(s[length-64:], uint64(length), y)
	w := weakHashLen32WithSeedsByte(s[length-32:], uint64(length)*k1, k0)
	z += shiftMix(v.High) * k1
	x = rotate(z+x, 39) * k1
	y = rotate(y, 33) * k1

	// Decrease length to the nearest multiple of 64, and operate on 64-byte chunks.
	tmpLength := uint32(length)
	tmpLength = (tmpLength - 1) & ^uint32(63)
	for {
		x = rotate(x+y+v.Low+fetch64(s[16:]), 37) * k1
		y = rotate(y+v.High+fetch64(s[48:]), 42) * k1

		x ^= w.High
		y ^= v.Low

		z = rotate(z^w.Low, 33)
		v = weakHashLen32WithSeedsByte(s, v.High*k1, x+w.Low)
		w = weakHashLen32WithSeedsByte(s[32:], z+w.High, y)
		z, x = x, z
		s = s[64:]
		tmpLength -= 64
		if tmpLength == 0 {
			break
		}
	}

	return ch16(
		ch16(v.Low, w.Low)+shiftMix(y)*k1+z,
		ch16(v.High, w.High)+x,
	)
}
