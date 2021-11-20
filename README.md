# city [![](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/go-faster/city#section-documentation) [![](https://img.shields.io/codecov/c/github/go-faster/city?label=cover)](https://codecov.io/gh/go-faster/city) [![alpha](https://img.shields.io/badge/-alpha-orange)](https://go-faster.org/docs/projects/status#alpha)
[CityHash](https://github.com/google/cityhash) in Go. Fork of [tenfyzhong/cityhash](https://github.com/tenfyzhong/cityhash).

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
CityHash32-32      298ns ± 1%     295ns ± 2%      ~     (p=0.421 n=5+5)
CityHash64-32      336ns ± 0%     145ns ± 2%   -56.92%  (p=0.008 n=5+5)
CityHash128-32     353ns ± 0%     149ns ± 2%   -57.81%  (p=0.008 n=5+5)

name            old speed      new speed      delta
CityHash32-32   3.44GB/s ± 1%  3.47GB/s ± 2%      ~     (p=0.421 n=5+5)
CityHash64-32   3.04GB/s ± 0%  7.07GB/s ± 2%  +132.16%  (p=0.008 n=5+5)
CityHash128-32  2.90GB/s ± 0%  6.87GB/s ± 2%  +137.06%  (p=0.008 n=5+5)
```

## Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/go-faster/city
cpu: AMD Ryzen 9 5950X 16-Core Processor            
BenchmarkClickHouse128-32     137.7 ns/op  7437.19 MB/s  0 B/op  0 allocs/op
BenchmarkClickHouse64-32      131.5 ns/op  7786.84 MB/s  0 B/op  0 allocs/op
BenchmarkCityHash32-32        333.9 ns/op  3066.73 MB/s  0 B/op  0 allocs/op
BenchmarkCityHash64-32        141.9 ns/op  7216.10 MB/s  0 B/op  0 allocs/op
BenchmarkCityHash128-32       148.5 ns/op  6897.44 MB/s  0 B/op  0 allocs/op
BenchmarkCityHash64Small-32   3.659 ns/op  2732.82 MB/s  0 B/op  0 allocs/op
BenchmarkCityHash128Small-32  9.103 ns/op  1098.57 MB/s  0 B/op  0 allocs/op
PASS
ok      github.com/go-faster/city       12.018s
```