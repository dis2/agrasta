# Experimental instances of Agrasta

https://eprint.iacr.org/2018/181

* 5 rounds
* 255 bit block size
* 254 bit claimed security level
* Affine matrices with proper rank generated as product of LU factors

This is still naive (Go compiler), but bitsliced. That gets us:

```
BenchmarkNewMatrix          3000            394505 ns/op
BenchmarkCrypt               500           2460989 ns/op
```

SIMD is expected to bring this down 4-8x more. Majority of the time is spent
generating the matrices. The cost of doing the encryption as such is neglible.

Matrix generation complexity is O(BlockSize^2), the operation being population
count of row products.

For 127bit security, things are considerably faster, but such parameterization
doesn't seem to have a large security margin anymore.

```
BenchmarkNewMatrix         20000             87211 ns/op
BenchmarkCrypt              3000            442646 ns/op
```

In [ZKBoo based schemes](https://github.com/IAIK/Picnic), ballpark figure for Rasta is cutting communicated view sizes in half as a result of lower MC. The higher computational complexity
is generally worth the resulting smaller signatures.

