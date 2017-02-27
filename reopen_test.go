package reopen

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// TestReopenAppend -- make sure we always append to an existing file
//
// 1. Create a sample file using normal means
// 2. Open a ioreopen.File
//    write line 1
// 3. call Reopen
//    write line 2
// 4. close file
// 5. read file, make sure it contains line0,line1,line2
//
func TestReopenAppend(t *testing.T) {
	// TODO fix
	var fname = "/tmp/foo"

	// Step 1 -- Create a sample file using normal means
	forig, err := os.Create(fname)
	if err != nil {
		t.Fatalf("Unable to create initial file %s: %s", fname, err)
	}
	_, err = forig.Write([]byte("line0\n"))
	if err != nil {
		t.Fatalf("Unable to write initial line %s: %s", fname, err)
	}
	err = forig.Close()
	if err != nil {
		t.Fatalf("Unable to close initial file: %s", err)
	}

	// Test that making a new File appends
	f, err := NewFileWriter(fname)
	if err != nil {
		t.Fatalf("Unable to create %s", fname)
	}
	_, err = f.Write([]byte("line1\n"))
	if err != nil {
		t.Errorf("Got write error1: %s", err)
	}

	// Test that reopen always appends
	err = f.Reopen()
	if err != nil {
		t.Errorf("Got reopen error %s: %s", fname, err)
	}
	_, err = f.Write([]byte("line2\n"))
	if err != nil {
		t.Errorf("Got write error2 on %s: %s", fname, err)
	}
	err = f.Close()
	if err != nil {
		t.Errorf("Got closing error for %s: %s", fname, err)
	}

	out, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("Unable read in final file %s: %s", fname, err)
	}

	outstr := string(out)
	if outstr != "line0\nline1\nline2\n" {
		t.Errorf("Result was %s", outstr)
	}
}

// Test that reopen works when Inode is swapped out
// 1. Create a sample file using normal means
// 2. Open a ioreopen.File
//    write line 1
// 3. call Reopen
//    write line 2
// 4. close file
// 5. read file, make sure it contains line0,line1,line2
//
func TestChangeInode(t *testing.T) {
	// TODO fix
	var fname = "/tmp/foo"

	// Step 1 -- Create a empty sample file
	forig, err := os.Create(fname)
	if err != nil {
		t.Fatalf("Unable to create initial file %s: %s", fname, err)
	}
	err = forig.Close()
	if err != nil {
		t.Fatalf("Unable to close initial file: %s", err)
	}

	// Test that making a new File appends
	f, err := NewFileWriter(fname)
	if err != nil {
		t.Fatalf("Unable to create %s", fname)
	}
	_, err = f.Write([]byte("line1\n"))
	if err != nil {
		t.Errorf("Got write error1: %s", err)
	}

	// Now move file
	err = os.Rename(fname, fname+".orig")
	if err != nil {
		t.Errorf("Renaming error: %s", err)
	}
	f.Write([]byte("after1\n"))

	// Test that reopen always appends
	err = f.Reopen()
	if err != nil {
		t.Errorf("Got reopen error %s: %s", fname, err)
	}
	_, err = f.Write([]byte("line2\n"))
	if err != nil {
		t.Errorf("Got write error2 on %s: %s", fname, err)
	}
	err = f.Close()
	if err != nil {
		t.Errorf("Got closing error for %s: %s", fname, err)
	}

	out, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("Unable read in final file %s: %s", fname, err)
	}
	outstr := string(out)
	if outstr != "line2\n" {
		t.Errorf("Result was %s", outstr)
	}
}

func TestBufferedWriter(t *testing.T) {
	const fname = "/tmp/foo"
	fw, err := NewFileWriter(fname)
	if err != nil {
		t.Fatalf("Unable to create test file %q: %s", fname, err)
	}
	logger := NewBufferedFileWriterSize(fw, 0, time.Millisecond*100)
	// make sure 3-4 flush events happen
	time.Sleep(time.Millisecond * 4)

	// close should shutdown the flush
	logger.Close()

	// if flush is still happening we should get one or two here
	time.Sleep(time.Millisecond * 200)
}
