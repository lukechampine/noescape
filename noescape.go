package noescape // import "lukechampine.com/noescape"

import "io"

//go:noescape
func callRead(r io.Reader, b []byte) (int, error)

// Read calls r.Read(b) while promising that b does not escape.
//
// If you call Read with an io.Reader that actually does retain b, you are
// likely to be eaten by a grue.
func Read(r io.Reader, b []byte) (int, error) { return callRead(r, b) }

// Reader adds a "non-escaping" promise to an existing io.Reader.
type Reader struct {
	R io.Reader
}

// Read implements io.Reader.
func (r Reader) Read(b []byte) (int, error) { return Read(r.R, b) }

//go:noescape
func callWrite(w io.Writer, b []byte) (int, error)

// Write calls r.Write(b) while promising that b does not escape.
//
// If you call Write with an io.Writer that actually does retain b, you are
// likely to be eaten by a grue.
func Write(w io.Writer, b []byte) (int, error) { return callWrite(w, b) }

// Writer adds a "non-escaping" promise to an existing io.Writer.
type Writer struct {
	W io.Writer
}

// Write implements io.Writer.
func (w Writer) Write(b []byte) (int, error) { return Write(w.W, b) }
