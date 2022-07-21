# HAL vector

HTTP Accept-Language vector parser. Based on [accept-language-parser](https://github.com/opentable/accept-language-parser)
library, but based on vector parser instead of regexp due to performance reasons and reduce pointers policy.

## Usage

```go
src := "fr-CA,fr;q=0.2,en-US;q=0.6,en;q=0.4,*;q=0.5"
vec := halvector.Acquire()
defer halvector.Release(vec)
_ = vec.ParseStr(src)
vec.Sort().Root().Each(func(idx int, node *vector.Node) {
print(idx, ":")
println("code:", node.GetString("code"))
println("script:", node.GetString("script"))
println("region:", node.GetString("region"))
q, _ := node.Get("quality").Float()
println("quality:", q)
})
```

Output:
```
0 :
code: fr
script: 
region: CA
quality: +1.000000e+000
1 :
code: en
script: 
region: US
quality: +6.000000e-001
2 :
code: *
script: 
region: 
quality: +5.000000e-001
3 :
code: en
script: 
region: 
quality: +4.000000e-001
4 :
code: fr
script: 
region: 
quality: +2.000000e-001

```
