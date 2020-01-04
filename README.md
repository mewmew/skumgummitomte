# sgt

A winter times adventure into compilers.

## Installation

```bash
go get github.com/mewmew/skumgummitomte/cmd/sgt
```

## Example

### "hello world"

Compile and run [examples/hello/hello.go](examples/hello/hello.go).
```bash
$ sgt -o hello.ll examples/hello/hello.go
$ llvm-link -S -o main.ll hello.ll std/builtin.ll
$ lli main.ll
# Output:
#
# hello world
```

### Local variables

Compile and run [examples/locals/main.go](examples/locals/main.go).
```bash
$ sgt -o locals.ll examples/locals/main.go
$ lli locals.ll ; echo $?
# Output:
#
# 42
```

### Closures

Compile and run [examples/closures/closures.go](examples/closures/closures.go).
```bash
$ sgt -o closures.ll examples/closures/closures.go
$ llvm-link -S -o main.ll closures.ll std/builtin.ll
$ lli main.ll
# Output:
#
# lol
```

### Type definitions

Compile [examples/types/types.go](examples/types/types.go).
```bash
$ sgt -o types.ll examples/types/types.go
```

### Methods

Compile and run [examples/methods/methods.go](examples/methods/methods.go).
```bash
$ sgt -o methods.ll examples/methods/methods.go
$ llvm-link -S -o main.ll methods.ll std/builtin.ll
$ lli main.ll
# Output:
#
# T.M1
# T.M2
```

### Package imports

Compile and run `main` program [examples/imports/cmd/foo](examples/imports/cmd/foo/main.go) importing Go package [examples/imports/p](examples/imports/p/p.go).
```bash
$ sgt -o foo.ll ./examples/imports/cmd/foo
$ sgt -o p.ll ./examples/imports/p
$ llvm-link -S -o main.ll foo.ll p.ll std/builtin.ll
$ lli main.ll
# Output:
#
# p.Foo
```

### Named constants

Compile and run `main` program [examples/consts/cmd/foo](examples/consts/cmd/foo/main.go) importing Go package [examples/consts/p](examples/consts/p/p.go).
```bash
$ sgt -o foo.ll ./examples/consts/cmd/foo
$ sgt -o p.ll ./examples/consts/p
$ llvm-link -S -o main.ll foo.ll p.ll std/builtin.ll
$ lli main.ll
# Output:
#
# test
```

### Multiple return values

Compile and run `main` program [examples/multiple_results](examples/multiple_results/main.go).
```bash
$ sgt -o multiple_results.ll ./examples/multiple_results
$ llvm-link -S -o main.ll multiple_results.ll std/builtin.ll
$ lli main.ll
# Output:
#
# bar
# foo
```

### Slices

Compile and run `main` program [examples/slices](examples/slices/main.go).
```bash
$ sgt -o slices.ll ./examples/slices
$ llvm-link -S -o main.ll slices.ll std/builtin.ll
$ lli main.ll
# Output:
#
# qux
# baz
# bar
# foo
```

### Synthesized `len` function

Compile and run [examples/length/main.go](examples/length/main.go).
```bash
$ sgt -o length.ll examples/length/main.go
$ lli length.ll ; echo $?
# Output:
#
# 42
```

### String comparison

Compile and run [examples/string_compare/main.go](examples/string_compare/main.go).
```bash
$ sgt -o string_compare.ll examples/string_compare/main.go
$ llvm-link -S -o main.ll string_compare.ll std/builtin.ll
$ lli main.ll
# Output:
#
# x:
# foo
# y:
# foo
# x less than or equal to y
# x greater than or equal to y
# x equal to y
#
# x:
# abc
# y:
# def
# x less than y
# x less than or equal to y
# x not equal to y
#
# x:
# abc
# y:
# abcd
# x less than y
# x less than or equal to y
# x not equal to y
#
# x:
# abcd
# y:
# abc
# x greater than y
# x greater than or equal to y
# x not equal to y
#
# x:
# abx
# y:
# abc
# x greater than y
# x greater than or equal to y
# x not equal to y
```
