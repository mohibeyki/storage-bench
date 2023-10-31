package reader

import (
	"io"
)

type ZeroReader struct {
	Reader
}

func (r *ZeroReader) Read(b []byte) (n int, err error) {
	if r.Index >= r.Size {
		return 0, io.EOF
	}

	// Limiting our buffer size to 32KiB
	n = int(min(uint64(len(b)), uint64(32*1024), r.Size-r.Index))

	for i := 0; i < n; i++ {
		b[i] = 0
	}

	r.Index += uint64(n)

	return
}
