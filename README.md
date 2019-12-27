# sgt

A winter times adventure into compilers.

## Installation

```bash
go get github.com/mewmew/skumgummitomte/cmd/sgt
```

## Example

```bash
$ sgt -o hello.ll examples/hello/hello.go
$ llvm-link -S -o main.ll hello.ll std/builtin.ll
$ lli main.ll
# Output:
#
# hello world
```
