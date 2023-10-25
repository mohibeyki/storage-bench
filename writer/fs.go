package writer

import (
	"io"
	"os"
)

type FSWriter struct {
}

func (w *FSWriter) WriteFile(path string, reader io.Reader) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	buf := make([]byte, 64*1024)
	_, err = io.CopyBuffer(f, reader, buf)
	return
}
