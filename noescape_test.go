package noescape

import (
	"bytes"
	"strings"
	"testing"
)

type readImpl struct{}

func (readImpl) Read(b []byte) (int, error) {
	n := copy(b, "Hello, world!")
	return n, nil
}

func TestRead(t *testing.T) {
	r := strings.NewReader("Hello, world!")
	b := make([]byte, 5)
	n, err := Read(r, b)
	if n != 5 || err != nil || string(b) != "Hello" {
		t.Fatal("bad Read:", n, err, b)
	}
}

func TestReadAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(10, func() {
		s := readImpl{}
		b := make([]byte, 5)
		Read(s, b)
	})
	if allocs > 0 {
		t.Fatal("Read should not allocate", allocs)
	}
}

func TestReaderAllocs(t *testing.T) {
	r := Reader{readImpl{}}
	allocs := testing.AllocsPerRun(10, func() {
		b := make([]byte, 5)
		r.Read(b)
	})
	if allocs > 0 {
		t.Fatal("Read should not allocate", allocs)
	}
}

type writeImpl struct{}

func (writeImpl) Write(b []byte) (int, error) {
	return len(b), nil
}

func TestWrite(t *testing.T) {
	w := new(bytes.Buffer)
	b := []byte("Hello, world!")
	n, err := Write(w, b)
	if n != len(b) || err != nil {
		t.Fatal("bad Write:", n, err)
	}
}

func TestWriteAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(10, func() {
		b := make([]byte, 5)
		Write(writeImpl{}, b)
	})
	if allocs > 0 {
		t.Fatal("Write should not allocate", allocs)
	}
}

func TestWriterAllocs(t *testing.T) {
	w := Writer{writeImpl{}}
	allocs := testing.AllocsPerRun(10, func() {
		b := make([]byte, 5)
		w.Write(b)
	})
	if allocs > 0 {
		t.Fatal("Write should not allocate", allocs)
	}
}
