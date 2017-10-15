# YouFoodz
Awesome parser, v2



## Run all Go-tests
```
go test -run ''
```


## Fuzz testing for parser:
```
go-fuzz-build github.com/kaanixir/YouFoodz
mkdir -p fuzzcorpus
go-fuzz -bin=kaanparser-fuzz.zip -workdir=fuzzcorpus
```
