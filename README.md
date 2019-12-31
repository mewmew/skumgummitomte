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