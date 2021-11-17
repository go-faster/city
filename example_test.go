package city_test

import (
	"fmt"

	"github.com/go-faster/city"
)

func ExampleHash32() {
	s := []byte("hello")
	hash32 := city.Hash32(s)
	fmt.Printf("the 32-bit hash of 'hello' is: 0x%x\n", hash32)

	// Output:
	// the 32-bit hash of 'hello' is: 0x79969366
}

func ExampleHash64() {
	s := []byte("hello")
	hash64 := city.Hash64(s)
	fmt.Printf("the 64-bit hash of 'hello' is: 0x%x\n", hash64)

	// Output:
	// the 64-bit hash of 'hello' is: 0xb48be5a931380ce8
}

func ExampleHash128() {
	fmt.Println(city.Hash128([]byte("hello")))

	// Output: 6f72e4abb491a74a65148f580b45f347
}
