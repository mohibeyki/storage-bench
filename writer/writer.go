package writer

import "io"

type Writer interface {
	WriteFile(path string, reader io.Reader) error
}
