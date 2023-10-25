package reader

import (
	"crypto/rand"
	"io"
)

type RandomReader struct {
	Size  uint64
	Index uint64
}

func (r *RandomReader) Read(b []byte) (n int, err error) {
	if r.Index >= r.Size {
		return 0, io.EOF
	}

	n, err = rand.Read(b)
	r.Index += uint64(n)

	return
}
