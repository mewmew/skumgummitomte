# sgt

A winter times adventure into compilers.

## Installation

```bash
go get github.com/mewmew/skumgummitomte/cmd/sgt
```

## Example


Print "hello world":
```bash
$ sgt -o hello.ll examples/hello/hello.go
$ llvm-link -S -o main.ll hello.ll std/builtin.ll
$ lli main.ll
# Output:
#
# hello world
```

Use of local variables:
```bash
$ sgt -o locals.ll examples/locals/locals.go
$ lli locals.ll ; echo $?
# Output:
#
# 42
```
