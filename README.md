# city

Google CityHash in Go. 

CityHash provides hash functions for strings. 

Fork of [tenfyzhong/cityhash](https://github.com/tenfyzhong/cityhash).

[cityhash homepage](https://github.com/google/cityhash)


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