package reader

import (
	"crypto/rand"
	"io"
)

type RandomReader struct {
	Reader
}

func (r *RandomReader) Read(b []byte) (n int, err error) {
	if r.Index >= r.Size {
		return 0, io.EOF
	}

	// Limiting buffer size to 32KiB
	n = int(min(uint64(len(b)), uint64(32*1024), r.Size-r.Index))

	values := make([]byte, n)
	if _, err = rand.Read(values); err != nil {
		return
	}

	n = copy(b, values)
	r.Index += uint64(n)

	return
}
