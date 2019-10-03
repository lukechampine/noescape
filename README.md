noescape
--------

[![GoDoc](https://godoc.org/lukechampine.com/noescape?status.svg)](https://godoc.org/lukechampine.com/noescape)
[![Go Report Card](http://goreportcard.com/badge/lukechampine.com/noescape)](https://goreportcard.com/report/lukechampine.com/noescape)

```
go get lukechampine.com/noescape
```

`noescape` provides `Read` and `Write` functions that do not heap-allocate their argument.

Normally, when you pass a `[]byte` to an `io.Reader` or `io.Writer`, the
compiler must heap-allocate the slice data. This is because, at compile-time,
there is no way to know which concrete type is satisfying the interface, and
therefore the compiler cannot prove that the slice data will not be retained.

This is sad, because the vast majority of `Read` and `Write` methods do not
retain their argument, but still incur the performance penalty of
heap-allocation. The `noescape` package allows you to promise to the compiler
that your `Read` or `Write` method is perfectly safe, thank you very much, thus
allowing you to avoid the allocation.

This can be illustrated via benchmark:

```go
type yesReader struct{}

func (yesReader) Read(p []byte) (int, error) {
    return copy(p, "yes"), nil
}

func BenchmarkConcrete(b *testing.B) {
    r := yesReader{}
    for i := 0; i < b.N; i++ {
        buf := make([]byte, 100)
        r.Read(buf)
    }
}

func BenchmarkInterface(b *testing.B) {
    var r io.Reader = yesReader{}
    for i := 0; i < b.N; i++ {
        buf := make([]byte, 100)
        r.Read(buf)
    }
}

func BenchmarkNoEscape(b *testing.B) {
    var r io.Reader = yesReader{}
    for i := 0; i < b.N; i++ {
        buf := make([]byte, 100)
        noescape.Read(r, buf)
    }
}
```

```
BenchmarkConcrete-4     1000000000        0.372 ns/op       0 B/op      0 allocs/op
BenchmarkInterface-4      28114684       44.8 ns/op       112 B/op      1 allocs/op
BenchmarkNoEscape-4       88440447       11.8 ns/op        0 B/op       0 allocs/op
```

## How?

The `gc` compiler recognizes a `//go:noescape` pragma that promises to the
compiler that a function's arguments do not escape. So we just need to stick
this pragma on top of our `Read` and `Write` functions:

```go
//go:noescape
func Read(r io.Reader, b []byte) (int, error) { return r.Read(b) }
//go:noescape
func Write(w io.Writer, b []byte) (int, error) { return w.Write(b) }
```

There's a complication, though: the `go:noescape` pragma can only be applied to
externally-defined functions. So in order to use it, we need to implement `Read`
and `Write` in assembly! Not a big deal, it's just one line of Go, after all;
but it turned out to be trickier than I thought. If you're interested in the
details, check out the comments in `noescape_amd64.s`.

On that note: only `amd64` is supported for now. If you want to contribute
implementations for other architectures, I'll gladly merge them. Also, for
maximum compatibility, I should add non-assembly implementations; this would
defeat the whole point, but it would also allow this package to be used in
cross-platform code as an architecture-dependent optimization.