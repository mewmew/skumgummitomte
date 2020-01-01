# sgt

A winter times adventure into compilers.

## Installation

```bash
go get github.com/mewmew/skumgummitomte/cmd/sgt
```

## Example


Compile and run [examples/hello/hello.go](examples/hello/hello.go).
```bash
$ sgt -o hello.ll examples/hello/hello.go
$ llvm-link -S -o main.ll hello.ll std/builtin.ll
$ lli main.ll
# Output:
#
# hello world
```

Compile and run [examples/locals/locals.go](examples/locals/locals.go).
```bash
$ sgt -o locals.ll examples/locals/locals.go
$ lli locals.ll ; echo $?
# Output:
#
# 42
```

Compile and run [examples/closures/closures.go](examples/closures/closures.go).
```bash
$ sgt -o closures.ll examples/closures/closures.go
$ llvm-link -S -o main.ll closures.ll std/builtin.ll
$ lli main.ll
# Output:
#
# lol
```

Compile [examples/types/types.go](examples/types/types.go).
```bash
$ sgt -o types.ll examples/types/types.go
```
