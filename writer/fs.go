package writer

import "os"

type FSWriter struct {
}

func (fsWriter *FSWriter) WriteFile(path string, size uint64, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	for ; size > 0; size-- {
		if _, err = f.Write(data); err != nil {
			return err
		}
	}

	return nil
}
