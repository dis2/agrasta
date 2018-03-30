# Experimental instances of Agrasta

https://eprint.iacr.org/2018/181

* 5 rounds
* 255 bit block size
* 254 bit claimed security level
* Matrices generated as LU decomposition (probably the fastest way)

This is still naive (Go compiler), but bitsliced. That gets us:

```
BenchmarkNewMatrix          3000            394505 ns/op
BenchmarkCrypt               500           2460989 ns/op
```

SIMD is expected to bring this down 4-8x more. Majority of the time is spent
in generating the matrices. The cost of doing the encryption as such is neglible.

Matrix generation cost is O(BlockSize^2), spent counting parities of row products.

For 127bit security, things are considerably faster, but such parameterization
doesn't seem to have a large security margin anymore.

```
BenchmarkNewMatrix         20000             87211 ns/op
BenchmarkCrypt              3000            442646 ns/op
```

For GWE setting, the ballpark figure seems that Rasta cuts proof sizes in half
compared to LowMC, but has a higher startup cost - making the usual PQ PK schemes
in there about 2x slower. Still, this is very favorable tradeoff.



