# Progress

## Go standard library

`sgt` is capable of compiling the following packages of the Go standard library to LLVM IR.

- [ ] [archive/tar](https://golang.org/pkg/archive/tar)
- [ ] [archive/zip](https://golang.org/pkg/archive/zip)
- [ ] [bufio](https://golang.org/pkg/bufio)
- [ ] [builtin](https://golang.org/pkg/builtin)
- [ ] [bytes](https://golang.org/pkg/bytes)
- [ ] [cmd/addr2line](https://golang.org/cmd/addr2line)
- [ ] [cmd/api](https://golang.org/cmd/api)
- [ ] [cmd/asm](https://golang.org/cmd/asm)
- [ ] [cmd/asm/internal/arch](https://golang.org/cmd/asm/internal/arch)
- [ ] [cmd/asm/internal/asm](https://golang.org/cmd/asm/internal/asm)
- [ ] [cmd/asm/internal/flags](https://golang.org/cmd/asm/internal/flags)
- [ ] [cmd/asm/internal/lex](https://golang.org/cmd/asm/internal/lex)
- [ ] [cmd/buildid](https://golang.org/cmd/buildid)
- [ ] [cmd/cgo](https://golang.org/cmd/cgo)
- [ ] [cmd/compile](https://golang.org/cmd/compile)
- [ ] [cmd/compile/internal/amd64](https://golang.org/cmd/compile/internal/amd64)
- [ ] [cmd/compile/internal/arm](https://golang.org/cmd/compile/internal/arm)
- [ ] [cmd/compile/internal/arm64](https://golang.org/cmd/compile/internal/arm64)
- [ ] [cmd/compile/internal/gc](https://golang.org/cmd/compile/internal/gc)
- [ ] [cmd/compile/internal/gc/builtin](https://golang.org/cmd/compile/internal/gc/builtin)
- [ ] [cmd/compile/internal/logopt](https://golang.org/cmd/compile/internal/logopt)
- [ ] [cmd/compile/internal/mips](https://golang.org/cmd/compile/internal/mips)
- [ ] [cmd/compile/internal/mips64](https://golang.org/cmd/compile/internal/mips64)
- [ ] [cmd/compile/internal/ppc64](https://golang.org/cmd/compile/internal/ppc64)
- [ ] [cmd/compile/internal/s390x](https://golang.org/cmd/compile/internal/s390x)
- [ ] [cmd/compile/internal/ssa](https://golang.org/cmd/compile/internal/ssa)
- [ ] [cmd/compile/internal/ssa/gen](https://golang.org/cmd/compile/internal/ssa/gen)
- [ ] [cmd/compile/internal/syntax](https://golang.org/cmd/compile/internal/syntax)
- [ ] [cmd/compile/internal/test](https://golang.org/cmd/compile/internal/test)
- [ ] [cmd/compile/internal/types](https://golang.org/cmd/compile/internal/types)
- [ ] [cmd/compile/internal/wasm](https://golang.org/cmd/compile/internal/wasm)
- [ ] [cmd/compile/internal/x86](https://golang.org/cmd/compile/internal/x86)
- [ ] [cmd/cover](https://golang.org/cmd/cover)
- [ ] [cmd/dist](https://golang.org/cmd/dist)
- [ ] [cmd/doc](https://golang.org/cmd/doc)
- [ ] [cmd/fix](https://golang.org/cmd/fix)
- [ ] [cmd/go](https://golang.org/cmd/go)
- [ ] [cmd/gofmt](https://golang.org/cmd/gofmt)
- [ ] [cmd/go/internal/auth](https://golang.org/cmd/go/internal/auth)
- [ ] [cmd/go/internal/base](https://golang.org/cmd/go/internal/base)
- [ ] [cmd/go/internal/bug](https://golang.org/cmd/go/internal/bug)
- [ ] [cmd/go/internal/cache](https://golang.org/cmd/go/internal/cache)
- [ ] [cmd/go/internal/cfg](https://golang.org/cmd/go/internal/cfg)
- [ ] [cmd/go/internal/clean](https://golang.org/cmd/go/internal/clean)
- [ ] [cmd/go/internal/cmdflag](https://golang.org/cmd/go/internal/cmdflag)
- [ ] [cmd/go/internal/doc](https://golang.org/cmd/go/internal/doc)
- [ ] [cmd/go/internal/envcmd](https://golang.org/cmd/go/internal/envcmd)
- [ ] [cmd/go/internal/fix](https://golang.org/cmd/go/internal/fix)
- [ ] [cmd/go/internal/fmtcmd](https://golang.org/cmd/go/internal/fmtcmd)
- [ ] [cmd/go/internal/generate](https://golang.org/cmd/go/internal/generate)
- [ ] [cmd/go/internal/get](https://golang.org/cmd/go/internal/get)
- [ ] [cmd/go/internal/help](https://golang.org/cmd/go/internal/help)
- [ ] [cmd/go/internal/imports](https://golang.org/cmd/go/internal/imports)
- [ ] [cmd/go/internal/list](https://golang.org/cmd/go/internal/list)
- [ ] [cmd/go/internal/load](https://golang.org/cmd/go/internal/load)
- [ ] [cmd/go/internal/lockedfile](https://golang.org/cmd/go/internal/lockedfile)
- [ ] [cmd/go/internal/lockedfile/internal/filelock](https://golang.org/cmd/go/internal/lockedfile/internal/filelock)
- [ ] [cmd/go/internal/modcmd](https://golang.org/cmd/go/internal/modcmd)
- [ ] [cmd/go/internal/modconv](https://golang.org/cmd/go/internal/modconv)
- [ ] [cmd/go/internal/modfetch](https://golang.org/cmd/go/internal/modfetch)
- [ ] [cmd/go/internal/modfetch/codehost](https://golang.org/cmd/go/internal/modfetch/codehost)
- [ ] [cmd/go/internal/modfetch/zip_sum_test](https://golang.org/cmd/go/internal/modfetch/zip_sum_test)
- [ ] [cmd/go/internal/modget](https://golang.org/cmd/go/internal/modget)
- [ ] [cmd/go/internal/modinfo](https://golang.org/cmd/go/internal/modinfo)
- [ ] [cmd/go/internal/modload](https://golang.org/cmd/go/internal/modload)
- [ ] [cmd/go/internal/mvs](https://golang.org/cmd/go/internal/mvs)
- [ ] [cmd/go/internal/par](https://golang.org/cmd/go/internal/par)
- [ ] [cmd/go/internal/renameio](https://golang.org/cmd/go/internal/renameio)
- [ ] [cmd/go/internal/robustio](https://golang.org/cmd/go/internal/robustio)
- [ ] [cmd/go/internal/run](https://golang.org/cmd/go/internal/run)
- [ ] [cmd/go/internal/search](https://golang.org/cmd/go/internal/search)
- [ ] [cmd/go/internal/str](https://golang.org/cmd/go/internal/str)
- [ ] [cmd/go/internal/test](https://golang.org/cmd/go/internal/test)
- [ ] [cmd/go/internal/tool](https://golang.org/cmd/go/internal/tool)
- [ ] [cmd/go/internal/txtar](https://golang.org/cmd/go/internal/txtar)
- [ ] [cmd/go/internal/version](https://golang.org/cmd/go/internal/version)
- [ ] [cmd/go/internal/vet](https://golang.org/cmd/go/internal/vet)
- [ ] [cmd/go/internal/web](https://golang.org/cmd/go/internal/web)
- [ ] [cmd/go/internal/work](https://golang.org/cmd/go/internal/work)
- [ ] [cmd/internal/bio](https://golang.org/cmd/internal/bio)
- [ ] [cmd/internal/browser](https://golang.org/cmd/internal/browser)
- [ ] [cmd/internal/buildid](https://golang.org/cmd/internal/buildid)
- [ ] [cmd/internal/diff](https://golang.org/cmd/internal/diff)
- [ ] [cmd/internal/dwarf](https://golang.org/cmd/internal/dwarf)
- [ ] [cmd/internal/edit](https://golang.org/cmd/internal/edit)
- [ ] [cmd/internal/gcprog](https://golang.org/cmd/internal/gcprog)
- [ ] [cmd/internal/goobj](https://golang.org/cmd/internal/goobj)
- [ ] [cmd/internal/goobj2](https://golang.org/cmd/internal/goobj2)
- [ ] [cmd/internal/obj](https://golang.org/cmd/internal/obj)
- [ ] [cmd/internal/objabi](https://golang.org/cmd/internal/objabi)
- [ ] [cmd/internal/objfile](https://golang.org/cmd/internal/objfile)
- [ ] [cmd/internal/obj/arm](https://golang.org/cmd/internal/obj/arm)
- [ ] [cmd/internal/obj/arm64](https://golang.org/cmd/internal/obj/arm64)
- [ ] [cmd/internal/obj/mips](https://golang.org/cmd/internal/obj/mips)
- [ ] [cmd/internal/obj/ppc64](https://golang.org/cmd/internal/obj/ppc64)
- [ ] [cmd/internal/obj/riscv](https://golang.org/cmd/internal/obj/riscv)
- [ ] [cmd/internal/obj/s390x](https://golang.org/cmd/internal/obj/s390x)
- [ ] [cmd/internal/obj/wasm](https://golang.org/cmd/internal/obj/wasm)
- [ ] [cmd/internal/obj/x86](https://golang.org/cmd/internal/obj/x86)
- [ ] [cmd/internal/src](https://golang.org/cmd/internal/src)
- [ ] [cmd/internal/sys](https://golang.org/cmd/internal/sys)
- [ ] [cmd/internal/test2json](https://golang.org/cmd/internal/test2json)
- [ ] [cmd/link](https://golang.org/cmd/link)
- [ ] [cmd/link/internal/amd64](https://golang.org/cmd/link/internal/amd64)
- [ ] [cmd/link/internal/arm](https://golang.org/cmd/link/internal/arm)
- [ ] [cmd/link/internal/arm64](https://golang.org/cmd/link/internal/arm64)
- [ ] [cmd/link/internal/ld](https://golang.org/cmd/link/internal/ld)
- [ ] [cmd/link/internal/loadelf](https://golang.org/cmd/link/internal/loadelf)
- [ ] [cmd/link/internal/loader](https://golang.org/cmd/link/internal/loader)
- [ ] [cmd/link/internal/loadmacho](https://golang.org/cmd/link/internal/loadmacho)
- [ ] [cmd/link/internal/loadpe](https://golang.org/cmd/link/internal/loadpe)
- [ ] [cmd/link/internal/loadxcoff](https://golang.org/cmd/link/internal/loadxcoff)
- [ ] [cmd/link/internal/mips](https://golang.org/cmd/link/internal/mips)
- [ ] [cmd/link/internal/mips64](https://golang.org/cmd/link/internal/mips64)
- [ ] [cmd/link/internal/objfile](https://golang.org/cmd/link/internal/objfile)
- [ ] [cmd/link/internal/ppc64](https://golang.org/cmd/link/internal/ppc64)
- [ ] [cmd/link/internal/riscv64](https://golang.org/cmd/link/internal/riscv64)
- [ ] [cmd/link/internal/s390x](https://golang.org/cmd/link/internal/s390x)
- [ ] [cmd/link/internal/sym](https://golang.org/cmd/link/internal/sym)
- [ ] [cmd/link/internal/wasm](https://golang.org/cmd/link/internal/wasm)
- [ ] [cmd/link/internal/x86](https://golang.org/cmd/link/internal/x86)
- [ ] [cmd/nm](https://golang.org/cmd/nm)
- [ ] [cmd/objdump](https://golang.org/cmd/objdump)
- [ ] [cmd/pack](https://golang.org/cmd/pack)
- [ ] [cmd/pprof](https://golang.org/cmd/pprof)
- [ ] [cmd/test2json](https://golang.org/cmd/test2json)
- [ ] [cmd/trace](https://golang.org/cmd/trace)
- [ ] [cmd/vet](https://golang.org/cmd/vet)
- [ ] [compress/bzip2](https://golang.org/pkg/compress/bzip2)
- [ ] [compress/flate](https://golang.org/pkg/compress/flate)
- [ ] [compress/gzip](https://golang.org/pkg/compress/gzip)
- [ ] [compress/lzw](https://golang.org/pkg/compress/lzw)
- [ ] [compress/zlib](https://golang.org/pkg/compress/zlib)
- [ ] [container/heap](https://golang.org/pkg/container/heap)
- [ ] [container/list](https://golang.org/pkg/container/list)
- [ ] [container/ring](https://golang.org/pkg/container/ring)
- [ ] [context](https://golang.org/pkg/context)
- [ ] [crypto](https://golang.org/pkg/crypto)
- [ ] [crypto/aes](https://golang.org/pkg/crypto/aes)
- [ ] [crypto/cipher](https://golang.org/pkg/crypto/cipher)
- [ ] [crypto/des](https://golang.org/pkg/crypto/des)
- [ ] [crypto/dsa](https://golang.org/pkg/crypto/dsa)
- [ ] [crypto/ecdsa](https://golang.org/pkg/crypto/ecdsa)
- [ ] [crypto/ed25519](https://golang.org/pkg/crypto/ed25519)
- [ ] [crypto/ed25519/internal/edwards25519](https://golang.org/pkg/crypto/ed25519/internal/edwards25519)
- [ ] [crypto/elliptic](https://golang.org/pkg/crypto/elliptic)
- [ ] [crypto/hmac](https://golang.org/pkg/crypto/hmac)
- [ ] [crypto/internal/randutil](https://golang.org/pkg/crypto/internal/randutil)
- [ ] [crypto/internal/subtle](https://golang.org/pkg/crypto/internal/subtle)
- [ ] [crypto/md5](https://golang.org/pkg/crypto/md5)
- [ ] [crypto/rand](https://golang.org/pkg/crypto/rand)
- [ ] [crypto/rc4](https://golang.org/pkg/crypto/rc4)
- [ ] [crypto/rsa](https://golang.org/pkg/crypto/rsa)
- [ ] [crypto/sha1](https://golang.org/pkg/crypto/sha1)
- [ ] [crypto/sha256](https://golang.org/pkg/crypto/sha256)
- [ ] [crypto/sha512](https://golang.org/pkg/crypto/sha512)
- [ ] [crypto/subtle](https://golang.org/pkg/crypto/subtle)
- [ ] [crypto/tls](https://golang.org/pkg/crypto/tls)
- [ ] [crypto/x509](https://golang.org/pkg/crypto/x509)
- [ ] [crypto/x509/pkix](https://golang.org/pkg/crypto/x509/pkix)
- [ ] [database/sql](https://golang.org/pkg/database/sql)
- [ ] [database/sql/driver](https://golang.org/pkg/database/sql/driver)
- [ ] [debug/dwarf](https://golang.org/pkg/debug/dwarf)
- [ ] [debug/elf](https://golang.org/pkg/debug/elf)
- [ ] [debug/gosym](https://golang.org/pkg/debug/gosym)
- [ ] [debug/macho](https://golang.org/pkg/debug/macho)
- [ ] [debug/pe](https://golang.org/pkg/debug/pe)
- [ ] [debug/plan9obj](https://golang.org/pkg/debug/plan9obj)
- [ ] [encoding](https://golang.org/pkg/encoding)
- [ ] [encoding/ascii85](https://golang.org/pkg/encoding/ascii85)
- [ ] [encoding/asn1](https://golang.org/pkg/encoding/asn1)
- [ ] [encoding/base32](https://golang.org/pkg/encoding/base32)
- [ ] [encoding/base64](https://golang.org/pkg/encoding/base64)
- [ ] [encoding/binary](https://golang.org/pkg/encoding/binary)
- [ ] [encoding/csv](https://golang.org/pkg/encoding/csv)
- [ ] [encoding/gob](https://golang.org/pkg/encoding/gob)
- [ ] [encoding/hex](https://golang.org/pkg/encoding/hex)
- [ ] [encoding/json](https://golang.org/pkg/encoding/json)
- [ ] [encoding/pem](https://golang.org/pkg/encoding/pem)
- [ ] [encoding/xml](https://golang.org/pkg/encoding/xml)
- [ ] [errors](https://golang.org/pkg/errors)
- [ ] [expvar](https://golang.org/pkg/expvar)
- [ ] [flag](https://golang.org/pkg/flag)
- [ ] [fmt](https://golang.org/pkg/fmt)
- [ ] [go/ast](https://golang.org/pkg/go/ast)
- [ ] [go/build](https://golang.org/pkg/go/build)
- [ ] [go/constant](https://golang.org/pkg/go/constant)
- [ ] [go/doc](https://golang.org/pkg/go/doc)
- [ ] [go/format](https://golang.org/pkg/go/format)
- [ ] [go/importer](https://golang.org/pkg/go/importer)
- [ ] [go/internal/gccgoimporter](https://golang.org/pkg/go/internal/gccgoimporter)
- [ ] [go/internal/gcimporter](https://golang.org/pkg/go/internal/gcimporter)
- [ ] [go/internal/srcimporter](https://golang.org/pkg/go/internal/srcimporter)
- [ ] [go/parser](https://golang.org/pkg/go/parser)
- [ ] [go/printer](https://golang.org/pkg/go/printer)
- [ ] [go/scanner](https://golang.org/pkg/go/scanner)
- [ ] [go/token](https://golang.org/pkg/go/token)
- [ ] [go/types](https://golang.org/pkg/go/types)
- [ ] [hash](https://golang.org/pkg/hash)
- [ ] [hash/adler32](https://golang.org/pkg/hash/adler32)
- [ ] [hash/crc32](https://golang.org/pkg/hash/crc32)
- [ ] [hash/crc64](https://golang.org/pkg/hash/crc64)
- [ ] [hash/fnv](https://golang.org/pkg/hash/fnv)
- [ ] [hash/maphash](https://golang.org/pkg/hash/maphash)
- [ ] [html](https://golang.org/pkg/html)
- [ ] [html/template](https://golang.org/pkg/html/template)
- [ ] [image](https://golang.org/pkg/image)
- [ ] [image/color](https://golang.org/pkg/image/color)
- [ ] [image/color/palette](https://golang.org/pkg/image/color/palette)
- [ ] [image/draw](https://golang.org/pkg/image/draw)
- [ ] [image/gif](https://golang.org/pkg/image/gif)
- [ ] [image/internal/imageutil](https://golang.org/pkg/image/internal/imageutil)
- [ ] [image/jpeg](https://golang.org/pkg/image/jpeg)
- [ ] [image/png](https://golang.org/pkg/image/png)
- [ ] [index/suffixarray](https://golang.org/pkg/index/suffixarray)
- [ ] [internal/bytealg](https://golang.org/pkg/internal/bytealg)
- [ ] [internal/cfg](https://golang.org/pkg/internal/cfg)
- [x] [internal/cpu](https://golang.org/pkg/internal/cpu) (as of 2020-01-04, rev [da6d870](https://github.com/mewmew/skumgummitomte/commit/da6d8704f1462c383211a9a56b61ed4bf55c07c2))
- [ ] [internal/fmtsort](https://golang.org/pkg/internal/fmtsort)
- [ ] [internal/goroot](https://golang.org/pkg/internal/goroot)
- [ ] [internal/goversion](https://golang.org/pkg/internal/goversion)
- [ ] [internal/lazyregexp](https://golang.org/pkg/internal/lazyregexp)
- [ ] [internal/lazytemplate](https://golang.org/pkg/internal/lazytemplate)
- [ ] [internal/nettrace](https://golang.org/pkg/internal/nettrace)
- [ ] [internal/oserror](https://golang.org/pkg/internal/oserror)
- [ ] [internal/poll](https://golang.org/pkg/internal/poll)
- [ ] [internal/race](https://golang.org/pkg/internal/race)
- [ ] [internal/reflectlite](https://golang.org/pkg/internal/reflectlite)
- [ ] [internal/singleflight](https://golang.org/pkg/internal/singleflight)
- [ ] [internal/syscall/unix](https://golang.org/pkg/internal/syscall/unix)
- [ ] [internal/syscall/windows](https://golang.org/pkg/internal/syscall/windows)
- [ ] [internal/syscall/windows/registry](https://golang.org/pkg/internal/syscall/windows/registry)
- [ ] [internal/syscall/windows/sysdll](https://golang.org/pkg/internal/syscall/windows/sysdll)
- [ ] [internal/testenv](https://golang.org/pkg/internal/testenv)
- [ ] [internal/testlog](https://golang.org/pkg/internal/testlog)
- [ ] [internal/trace](https://golang.org/pkg/internal/trace)
- [ ] [internal/xcoff](https://golang.org/pkg/internal/xcoff)
- [ ] [io](https://golang.org/pkg/io)
- [ ] [io/ioutil](https://golang.org/pkg/io/ioutil)
- [ ] [log](https://golang.org/pkg/log)
- [ ] [log/syslog](https://golang.org/pkg/log/syslog)
- [x] [math](https://golang.org/pkg/math) (as of 2020-01-03, rev [6edefbc](https://github.com/mewmew/skumgummitomte/commit/6edefbc00cbb33451b22acbed8702012e17c7913))
- [ ] [math/big](https://golang.org/pkg/math/big)
- [ ] [math/bits](https://golang.org/pkg/math/bits)
- [ ] [math/cmplx](https://golang.org/pkg/math/cmplx)
- [ ] [math/rand](https://golang.org/pkg/math/rand)
- [ ] [mime](https://golang.org/pkg/mime)
- [ ] [mime/multipart](https://golang.org/pkg/mime/multipart)
- [ ] [mime/quotedprintable](https://golang.org/pkg/mime/quotedprintable)
- [ ] [net](https://golang.org/pkg/net)
- [ ] [net/http](https://golang.org/pkg/net/http)
- [ ] [net/http/cgi](https://golang.org/pkg/net/http/cgi)
- [ ] [net/http/cookiejar](https://golang.org/pkg/net/http/cookiejar)
- [ ] [net/http/fcgi](https://golang.org/pkg/net/http/fcgi)
- [ ] [net/http/httptest](https://golang.org/pkg/net/http/httptest)
- [ ] [net/http/httptrace](https://golang.org/pkg/net/http/httptrace)
- [ ] [net/http/httputil](https://golang.org/pkg/net/http/httputil)
- [ ] [net/http/internal](https://golang.org/pkg/net/http/internal)
- [ ] [net/http/pprof](https://golang.org/pkg/net/http/pprof)
- [ ] [net/internal/socktest](https://golang.org/pkg/net/internal/socktest)
- [ ] [net/mail](https://golang.org/pkg/net/mail)
- [ ] [net/rpc](https://golang.org/pkg/net/rpc)
- [ ] [net/rpc/jsonrpc](https://golang.org/pkg/net/rpc/jsonrpc)
- [ ] [net/smtp](https://golang.org/pkg/net/smtp)
- [ ] [net/textproto](https://golang.org/pkg/net/textproto)
- [ ] [net/url](https://golang.org/pkg/net/url)
- [ ] [os](https://golang.org/pkg/os)
- [ ] [os/exec](https://golang.org/pkg/os/exec)
- [ ] [os/signal](https://golang.org/pkg/os/signal)
- [ ] [os/signal/internal/pty](https://golang.org/pkg/os/signal/internal/pty)
- [ ] [os/user](https://golang.org/pkg/os/user)
- [ ] [path](https://golang.org/pkg/path)
- [ ] [path/filepath](https://golang.org/pkg/path/filepath)
- [ ] [plugin](https://golang.org/pkg/plugin)
- [ ] [reflect](https://golang.org/pkg/reflect)
- [ ] [regexp](https://golang.org/pkg/regexp)
- [ ] [regexp/syntax](https://golang.org/pkg/regexp/syntax)
- [ ] [runtime](https://golang.org/pkg/runtime)
- [ ] [runtime/cgo](https://golang.org/pkg/runtime/cgo)
- [ ] [runtime/debug](https://golang.org/pkg/runtime/debug)
- [ ] [runtime/internal/atomic](https://golang.org/pkg/runtime/internal/atomic)
- [ ] [runtime/internal/math](https://golang.org/pkg/runtime/internal/math)
- [ ] [runtime/internal/sys](https://golang.org/pkg/runtime/internal/sys)
- [ ] [runtime/msan](https://golang.org/pkg/runtime/msan)
- [ ] [runtime/pprof](https://golang.org/pkg/runtime/pprof)
- [ ] [runtime/pprof/internal/profile](https://golang.org/pkg/runtime/pprof/internal/profile)
- [ ] [runtime/race](https://golang.org/pkg/runtime/race)
- [ ] [runtime/trace](https://golang.org/pkg/runtime/trace)
- [ ] [sort](https://golang.org/pkg/sort)
- [ ] [strconv](https://golang.org/pkg/strconv)
- [ ] [strings](https://golang.org/pkg/strings)
- [ ] [sync](https://golang.org/pkg/sync)
- [ ] [sync/atomic](https://golang.org/pkg/sync/atomic)
- [ ] [syscall](https://golang.org/pkg/syscall)
- [ ] [syscall/js](https://golang.org/pkg/syscall/js)
- [ ] [testing](https://golang.org/pkg/testing)
- [ ] [testing/internal/testdeps](https://golang.org/pkg/testing/internal/testdeps)
- [ ] [testing/iotest](https://golang.org/pkg/testing/iotest)
- [ ] [testing/quick](https://golang.org/pkg/testing/quick)
- [ ] [text/scanner](https://golang.org/pkg/text/scanner)
- [ ] [text/tabwriter](https://golang.org/pkg/text/tabwriter)
- [ ] [text/template](https://golang.org/pkg/text/template)
- [ ] [text/template/parse](https://golang.org/pkg/text/template/parse)
- [ ] [time](https://golang.org/pkg/time)
- [ ] [unicode](https://golang.org/pkg/unicode)
- [ ] [unicode/utf8](https://golang.org/pkg/unicode/utf8)
- [ ] [unicode/utf16](https://golang.org/pkg/unicode/utf16)
- [ ] [unsafe](https://golang.org/pkg/unsafe)
