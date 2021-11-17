# city [![](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/go-faster/city#section-documentation) [![](https://img.shields.io/codecov/c/github/go-faster/city?label=cover)](https://codecov.io/gh/go-faster/city) [![experimental](https://img.shields.io/badge/-experimental-blueviolet)](https://go-faster.org/docs/projects/status#experimental)


Google CityHash in Go. Fork of [tenfyzhong/cityhash](https://github.com/tenfyzhong/cityhash).

```console
go get github.com/go-faster/city
```

[CityHash](https://github.com/google/cityhash) provides hash functions for strings. 

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