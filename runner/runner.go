package runner

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/mohibeyki/storage-bench/reader"
	"github.com/mohibeyki/storage-bench/writer"
	"github.com/schollz/progressbar/v3"
)

type Runner struct {
	Path       string
	Files      uint64
	Threads    uint64
	Size       uint64
	Writer     writer.Writer
	ReaderType reflect.Type
	Bar        *progressbar.ProgressBar
}

func formatByteSize(bytes uint64) string {
	const unit = uint64(1024)
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (r *Runner) Run() error {
	fmt.Printf("Running storage-bench with [%d] files of size [%s] using [%s] by [%d] threads in [%s]\n", r.Files, formatByteSize(r.Size), r.ReaderType, r.Threads, r.Path)

	r.Bar = progressbar.Default(
		int64(r.Files),
	)

	// start timing
	startTime := time.Now()

	jobs := make(chan string, r.Files)

	var wg sync.WaitGroup
	for i := uint64(0); i < r.Threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := r.StartWorker(jobs); err != nil {
				panic(err)
			}
		}()
	}

	for i := uint64(0); i < r.Files; i++ {
		jobs <- fmt.Sprintf("%s/%06d.tmp", r.Path, i)
	}
	close(jobs)

	wg.Wait()

	duration := time.Since(startTime)
	durationMS := duration.Milliseconds()
	if durationMS == 0 {
		durationMS = 1
	}

	fmt.Printf("files: [%d]\ttotal: [%s]\tduration: [%dms]\taverage: [%s/s]\n", r.Files, formatByteSize(r.Files*r.Size), durationMS, formatByteSize(r.Files*r.Size/uint64(durationMS)*1000))

	return nil
}

func (r *Runner) StartWorker(jobs <-chan string) error {
	var err error

	for fileName := range jobs {
		var inputReader io.Reader

		switch r.ReaderType {
		case reflect.TypeOf(reader.ZeroReader{}):
			inputReader = &reader.ZeroReader{Reader: reader.Reader{Size: r.Size}}
		case reflect.TypeOf(reader.RandomReader{}):
			inputReader = &reader.RandomReader{Reader: reader.Reader{Size: r.Size}}
		default:
			panic("No suitable reader type found")
		}

		if err = r.Writer.WriteFile(fileName, inputReader); err != nil {
			return err
		}

		if err = r.Bar.Add(1); err != nil {
			return err
		}
	}

	return nil
}
