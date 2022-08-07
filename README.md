# city [![](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/go-faster/city#section-documentation) [![](https://img.shields.io/codecov/c/github/go-faster/city?label=cover)](https://codecov.io/gh/go-faster/city) [![stable](https://img.shields.io/badge/-stable-brightgreen)](https://go-faster.org/docs/projects/status#stable)
[CityHash](https://github.com/google/cityhash) in Go. Fork of [tenfyzhong/cityhash](https://github.com/tenfyzhong/cityhash).

Note: **prefer [xxhash](https://github.com/cespare/xxhash) as non-cryptographic hash algorithm**, this package is intended 
for places where CityHash is already used.

CityHash **is not compatible** to [FarmHash](https://github.com/google/farmhash), use [go-farm](https://github.com/dgryski/go-farm).

```console
go get github.com/go-faster/city
```

```go
city.Hash128([]byte("hello"))
```

* Faster
* Supports ClickHouse hash

```
name            old time/op    new time/op    delta
CityHash64-32      333ns ± 2%     108ns ± 3%   -67.57%  (p=0.000 n=10+10)
CityHash128-32     347ns ± 2%     112ns ± 2%   -67.74%  (p=0.000 n=9+10)

name            old speed      new speed      delta
CityHash64-32   3.08GB/s ± 2%  9.49GB/s ± 3%  +208.40%  (p=0.000 n=10+10)
CityHash128-32  2.95GB/s ± 2%  9.14GB/s ± 2%  +209.98%  (p=0.000 n=9+10)
```

## Examples

Let's take 64-bit hash from `Moscow` string.

```sql
:) SELECT cityHash64('Moscow')
12507901496292878638
```

```go
s := []byte("Moscow")
fmt.Print("ClickHouse: ")
fmt.Println(city.CH64(s))
fmt.Print("CityHash:   ")
fmt.Println(city.Hash64(s))
// Output:
// ClickHouse: 12507901496292878638
// CityHash:   5992710078453357409
```

You can use [test data corpus](https://github.com/go-faster/city/blob/main/_testdata/data.json) to check your implementation of ClickHouse CityHash variant if needed.

```json
{
  "Seed": {
    "Low": 5577006791947779410,
    "High": 8674665223082153551
  },
  "entries": [
    {
      "Input": "Moscow",
      "City32": 431367057,
      "City64": 5992710078453357409,
      "City128": {
        "Low": 10019773792274861915,
        "High": 12276543986707912152
      },
      "City128Seed": {
        "Low": 13396466470330251720,
        "High": 5508504338941663328
      },
      "ClickHouse64": 12507901496292878638,
      "ClickHouse128": {
        "Low": 3377444358654451565,
        "High": 2499202049363713365
      },
      "ClickHouse128Seed": {
        "Low": 568168482305327346,
        "High": 1719721512326527886
      }
    }
  ]
}
```

## Benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/go-faster/city
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkClickHouse128/16     2213.98 MB/s
BenchmarkClickHouse128/64     4712.24 MB/s
BenchmarkClickHouse128/256    7561.58 MB/s
BenchmarkClickHouse128/1024  10158.98 MB/s
BenchmarkClickHouse64        10379.89 MB/s
BenchmarkCityHash32           3140.54 MB/s
BenchmarkCityHash64           9508.45 MB/s
BenchmarkCityHash128          9304.27 MB/s
BenchmarkCityHash64Small      2700.84 MB/s
BenchmarkCityHash128Small     1175.65 MB/s
```
