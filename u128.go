package city

// U128 is uint128.
type U128 struct {
	Low, High uint64 // much faster than uint64[2]
}
