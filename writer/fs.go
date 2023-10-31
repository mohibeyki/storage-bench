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

	_, err = io.Copy(f, reader)
	return
}
