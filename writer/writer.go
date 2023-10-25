package writer

type Writer interface {
	WriteFile(path string, size uint64, data []byte) error
}
