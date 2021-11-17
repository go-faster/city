package city

import "fmt"

func ExampleHash32() {
	s := []byte("hello")
	hash32 := Hash32(s)
	fmt.Printf("the 32-bit hash of 'hello' is: 0x%x\n", hash32)

	// Output:
	// the 32-bit hash of 'hello' is: 0x79969366
}

func ExampleHash64() {
	s := []byte("hello")
	hash64 := Hash64(s)
	fmt.Printf("the 64-bit hash of 'hello' is: 0x%x\n", hash64)

	// Output:
	// the 64-bit hash of 'hello' is: 0xb48be5a931380ce8
}

func ExampleHash128() {
	s := []byte("hello")
	hash128 := Hash128(s)
	fmt.Printf("the 128-bit hash of 'hello' is: 0x%x%x\n", hash128.High, hash128.Low)

	// Output:
	// the 128-bit hash of 'hello' is: 0x65148f580b45f3476f72e4abb491a74a
}
