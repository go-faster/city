package cityhash

// U128 is uint128.
type U128 struct {
	Low, High uint64 // much fater than uint64[2]
}
