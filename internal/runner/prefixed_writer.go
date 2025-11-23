package runner

import (
	"bufio"
	"io"
	"sync"

	"github.com/azuyamat/pace/internal/logger"
)

type PrefixedWriter struct {
	taskName string
	isStdout bool
	scanner  *bufio.Scanner
	mu       sync.Mutex
	pr       *io.PipeReader
	pw       *io.PipeWriter
	done     chan struct{}
}

func NewPrefixedWriter(taskName string, isStdout bool) *PrefixedWriter {
	pr, pw := io.Pipe()
	writer := &PrefixedWriter{
		taskName: taskName,
		isStdout: isStdout,
		scanner:  bufio.NewScanner(pr),
		pr:       pr,
		pw:       pw,
		done:     make(chan struct{}),
	}
	go writer.processLines()
	return writer
}

func (w *PrefixedWriter) Write(p []byte) (n int, err error) {
	return w.pw.Write(p)
}

func (w *PrefixedWriter) Close() error {
	err := w.pw.Close()
	<-w.done
	return err
}

func (w *PrefixedWriter) processLines() {
	defer close(w.done)
	defer w.pr.Close()

	for w.scanner.Scan() {
		line := w.scanner.Text()
		w.mu.Lock()
		if w.isStdout {
			logger.Default.TaskOutput(w.taskName, "%s", line)
		} else {
			logger.Default.TaskError(w.taskName, "%s", line)
		}
		w.mu.Unlock()
	}
}
