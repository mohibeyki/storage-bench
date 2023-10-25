package runner

import (
	"crypto/rand"
	"fmt"
	"github.com/mohibeyki/storage-bench/writer"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Runner struct {
	Path    string
	Files   uint64
	Size    uint64
	Threads uint64
	Data    []byte
	Total   *atomic.Uint64
	Writer  writer.Writer
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
	fmt.Printf("Running storage-bench with [%d] files of size [%s] using [%d] threads in [%s]\n", r.Files, formatByteSize(r.Size), r.Threads, r.Path)

	data := make([]byte, 1024)
	_, err := rand.Read(data)
	if err != nil {
		return err
	}

	// start timing
	startTime := time.Now()

	var wg sync.WaitGroup
	for i := uint64(0); i < r.Threads; i++ {
		wg.Add(1)

		go func(threadName string) {
			if err := r.RunThread(threadName); err != nil {
				panic(err)
			}

			wg.Done()
		}(strconv.FormatUint(i, 10))
	}

	quitWatcher := make(chan bool)
	go func() {
		for {
			select {
			case <-quitWatcher:
				return
			default:
				fmt.Printf("[%d]/[%d]\n", r.Total.Load(), r.Files)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
	quitWatcher <- true

	duration := time.Since(startTime)
	durationMS := duration.Milliseconds()
	if durationMS == 0 {
		durationMS = 1
	}

	// input size is in 1KiB
	byteSize := r.Size * 1024
	fmt.Printf("files: [%d]\ttotal: [%s]\tduration: [%dms]\taverage: [%s/s]\n", r.Files, formatByteSize(r.Files*byteSize), durationMS, formatByteSize(r.Files*byteSize/uint64(durationMS)*1000))

	return nil
}

func (r *Runner) RunThread(name string) error {
	var err error
	files := r.Files / r.Threads

	for i := uint64(0); i < files; i++ {
		var fileName strings.Builder
		fileName.WriteString(r.Path)
		fileName.WriteString("/")
		fileName.WriteString(name)
		fileName.WriteString("-")
		fileName.WriteString(strconv.FormatUint(i, 10))
		fileName.WriteString(".tmp")

		if err = r.Writer.WriteFile(fileName.String(), r.Size, r.Data); err != nil {
			return err
		}

		r.Total.Add(1)
	}

	return nil
}
