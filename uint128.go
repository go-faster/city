package cityhash

// Uint128 type uint128
type Uint128 [2]uint64

// Low64 return the low 64-bit of u.
func (u Uint128) Low64() uint64 {
	return u[0]
}

// High64 return the high 64-bit of u.
func (u Uint128) High64() uint64 {
	return u[1]
}
