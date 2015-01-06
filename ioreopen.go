package goreopen

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

// Reopener interface defines something that can be reopened
type Reopener interface {
	Reopen() error
}

// ReopenWriter is a writer that also can be reopened
type Writer interface {
	Reopener
	io.Writer
}

// ReopenWriteCloser is a io.WriteCloser that can also be reopened
type WriteCloser interface {
	Reopener
	io.Writer
	io.Closer
}

// FileWriter that can also be reopened
type FileWriter struct {
	mu   sync.Mutex // ensures close / reopen / write are not called at the same time, protects f
	f    *os.File
	Name string
}

// Close calls the underlyding File.Close()
func (f *FileWriter) Close() error {
	f.mu.Lock()
	err := f.f.Close()
	f.mu.Unlock()
	return err
}

// mutex free version
func (f *FileWriter) reopen() error {
	if f.f != nil {
		f.f.Close()
		f.f = nil
	}
	newf, err := os.OpenFile(f.Name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		f.f = nil
		return err
	}
	f.f = newf

	return nil
}

// Reopen the file
func (f *FileWriter) Reopen() error {
	f.mu.Lock()
	err := f.reopen()
	f.mu.Unlock()
	return err
}

// Write implements the stander io.Writer interface
func (f *FileWriter) Write(p []byte) (int, error) {
	f.mu.Lock()
	n, err := f.f.Write(p)
	f.mu.Unlock()
	return n, err
}

// NewFileWriter opens a file for appending and writing and can be reopened.
// it is a ReopenWriteCloser...
func NewFileWriter(name string) (*FileWriter, error) {
	writer := FileWriter{
		f:    nil,
		Name: name,
	}
	err := writer.reopen()
	if err != nil {
		return nil, err
	}
	return &writer, nil
}

// BufferedWriter is buffer writer than can be reopned
type BufferedFileWriter struct {
	mu         sync.Mutex
	OrigWriter *FileWriter
	BufWriter  *bufio.Writer
}

// Reopen implement Reopener
func (bw *BufferedFileWriter) Reopen() error {
	bw.mu.Lock()
	bw.BufWriter.Flush()

	// use non-mutex version since we are using this one
	err := bw.OrigWriter.reopen()

	bw.BufWriter.Reset(io.Writer(bw.OrigWriter))
	bw.mu.Unlock()

	return err
}

func (bw *BufferedFileWriter) Close() error {
	bw.mu.Lock()
	bw.BufWriter.Flush()
	bw.OrigWriter.f.Close()
	bw.mu.Unlock()
	return nil
}

// Write implements io.Writer (and ReopenWriter)
func (bw *BufferedFileWriter) Write(p []byte) (int, error) {
	bw.mu.Lock()
	n, err := bw.BufWriter.Write(p)

	// Special Case... if the used space in the buffer is LESS than
	// the input, then we did a flush in the middle of the line
	// and the full log line was not sent on its way.
	if bw.BufWriter.Buffered() < len(p) {
		bw.BufWriter.Flush()
	}

	bw.mu.Unlock()
	return n, err
}

// flushDaemon periodically flushes the log file buffers.
func (bw *BufferedFileWriter) flushDaemon() {
	for range time.NewTicker(flushInterval).C {
		bw.mu.Lock()
		bw.BufWriter.Flush()
		bw.OrigWriter.f.Sync()
		bw.mu.Unlock()
	}
}

const bufferSize = 256 * 1024
const flushInterval = 30 * time.Second

func NewBufferedFileWriter(w *FileWriter) *BufferedFileWriter {
	bw := BufferedFileWriter{
		OrigWriter: w,
		BufWriter:  bufio.NewWriterSize(w, bufferSize),
	}
	go bw.flushDaemon()
	return &bw
}

type multiReopenWriter struct {
	writers []Writer
}

func (t *multiReopenWriter) Reopen() error {
	for _, w := range t.writers {
		err := w.Reopen()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *multiReopenWriter) Write(p []byte) (int, error) {
	for _, w := range t.writers {
		n, err := w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

// MultiWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
//  Also allow reopen
func MultiWriter(writers ...Writer) Writer {
	w := make([]Writer, len(writers))
	copy(w, writers)
	return &multiReopenWriter{w}
}

type nopReopenWriter struct {
	io.Writer
}

func (nopReopenWriter) Reopen() error {
	return nil
}

// NopReopenerWriter turns a normal writer into a ReopenWriter
//  by doing a NOP on Reopen
func NopWriter(w io.Writer) Writer {
	return nopReopenWriter{w}
}

// Reopenable versions of os.Stdout and os.Stderr (reopen does nothing)
var (
	Stdout = NopWriter(os.Stdin)
	Stderr = NopWriter(os.Stdin)
)
