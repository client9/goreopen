package ioreopen

import (
	"bufio"
	"io"
	"os"
)

// Reopener interface defines something that can be reopened
type Reopener interface {
	Reopen() error
}

// ReopenWriter is a writer that also can be reopened
type ReopenWriter interface {
	Reopener
	io.Writer
}

// ReopenWriteCloser is a io.WriteCloser that can also be reopened
type ReopenWriteCloser interface {
	Reopener
	io.Writer
	io.Closer
}

// File that can also be reopened
type File struct {
	f    *os.File
	Name string
}

// Close calls the underlyding File.Close()
func (f *File) Close() error {
	return f.f.Close()
}

// Reopen the file
func (f *File) Reopen() error {
	// f.f.Sync?
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

func (f *File) Write(p []byte) (int, error) {
	return f.f.Write(p)
}

// NewFile opens a file for appending and writing and can be reopened.
// it is a ReopenWriteCloser...
func NewFile(name string) (*File, error) {
	writer := File{
		f:    nil,
		Name: name,
	}
	err := writer.Reopen()
	if err != nil {
		return nil, err
	}
	return &writer, nil
}

// BufferedWriter is buffer writer than can be reopned
type BufferedWriter struct {
	OrigWriter ReopenWriter
	BufWriter  *bufio.Writer
}

// Reopen implement Reopener
func (bw *BufferedWriter) Reopen() error {
	bw.BufWriter.Flush()
	bw.OrigWriter.Reopen()
	bw.BufWriter.Reset(bw.OrigWriter)
	return nil
}

// Write implements io.Writer (and ReopenWriter)
func (bw *BufferedWriter) Write(p []byte) (int, error) {
	return bw.BufWriter.Write(p)
}

func NewBufferedWriter(w ReopenWriter) *BufferedWriter {
	return &BufferedWriter{
		OrigWriter: w,
		BufWriter:  bufio.NewWriter(w),
	}
}

type multiReopenWriter struct {
	writers []ReopenWriter
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
func MultiWriter(writers ...ReopenWriter) ReopenWriter {
	w := make([]ReopenWriter, len(writers))
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
func NopReopenerWriter(w io.Writer) ReopenWriter {
	return nopReopenWriter{w}
}

// Reopenable versions of os.Stdout and os.Stderr (reopen does nothing)
var (
	Stdout = NopReopenerWriter(os.Stdin)
	Stderr = NopReopenerWriter(os.Stdin)
)
