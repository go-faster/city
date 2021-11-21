package city_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/go-faster/city"
)

// Generated with:
// 	go run ./internal/citygen > _testdata/data.json
// Not worth it to use go:generate here.

//go:embed _testdata/data.json
var testData []byte

func TestData(t *testing.T) {
	var data struct {
		Seed    []byte `json:"seed"`
		Entries []struct {
			Input             string `json:"input"`
			City32            uint32 `json:"city_32"`
			City64            uint64 `json:"city_64,string"`
			City128           []byte `json:"city_128"`
			City128Seed       []byte `json:"city_128_seed"`
			ClickHouse64      uint64 `json:"clickhouse_64,string"`
			ClickHouse128     []byte `json:"clickhouse_128"`
			ClickHouse128Seed []byte `json:"clickhouse_128_seed"`
		}
	}
	if err := json.Unmarshal(testData, &data); err != nil {
		t.Fatal(err)
	}
	var seed city.U128
	seed.Set(data.Seed)

	for _, e := range data.Entries {
		var exp128, exp128Seed, expCH128, expCH128Seed city.U128
		exp128.Set(e.City128)
		exp128Seed.Set(e.City128Seed)
		expCH128.Set(e.ClickHouse128)
		expCH128Seed.Set(e.ClickHouse128Seed)
		input := []byte(e.Input)

		if v := city.Hash32(input); v != e.City32 {
			t.Errorf("Hash32(%q) %d (got) != %d (expected)", input, v, e.City32)
		}
		if v := city.Hash64(input); v != e.City64 {
			t.Errorf("Hash64(%q) %d (got) != %d (expected)", input, v, e.City64)
		}
		if v := city.Hash128(input); v != exp128 {
			t.Errorf("Hash128(%q) %s (got) != %s (expected)", input, v, exp128)
		}
		if v := city.Hash128Seed(input, seed); v != exp128Seed {
			t.Errorf("Hash128Seed(%q, %s) %s (got) != %s (expected)", input, seed, v, exp128Seed)
		}
		if v := city.CH64(input); v != e.ClickHouse64 {
			t.Errorf("CH64(%q) %d (got) != %d (expected)", input, v, e.City64)
		}
		if v := city.CH128(input); v != expCH128 {
			t.Errorf("CH128(%q) %s (got) != %s (expected)", input, v, expCH128)
		}
		if v := city.CH128Seed(input, seed); v != expCH128Seed {
			t.Errorf("CH128Seed(%q, %s) %s (got) != %s (expected)", input, seed, v, expCH128Seed)
		}
	}
}
