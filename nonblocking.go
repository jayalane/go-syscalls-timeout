// -*- tab-width:2 -*-

package timeouts

import (
	"errors"
	"os"
	"time"

	count "github.com/jayalane/go-counter"
)

const (
	two = 2
)

var suffix = "timeouts"

// this module is non-blocking versions of various os. calls that block or never
// return in the NFS case

// ReadDir is a ReadDir with a timeout in case you are calling it on a
// big NAS that might never reply Default is 60 seconds (1 minute);
// use ReadDirTimeout to tune this

// ReadDir is a wrapper to os.ReadDir with a 1 minute timeout.
func ReadDir(name string) ([]os.DirEntry, error) {
	var r []os.DirEntry

	var e error

	count.TimeFuncRunSuffix("readdir", func() {
		r, e = ReadDirTimeout(name, time.Minute)
	}, suffix)

	return r, e
}

// ReadDirTimeout is a ReadDir with a timeout in case you are calling
// it on a big NAS that might never reply.
func ReadDirTimeout(name string, t time.Duration) ([]os.DirEntry, error) {
	type res struct {
		de  []os.DirEntry
		err error
	}
	// start the ReadDir
	resCh := make(chan res, two)
	go func(name string) {
		de, err := os.ReadDir(name)
		resCh <- res{de, err}
	}(name)
	// now wait for it
	select {
	case result := <-resCh:
		count.IncrSuffix("readdir-ok", suffix)

		return result.de, result.err
	case <-time.After(t):
		count.IncrSuffix("readdir-timeout", suffix)

		return nil, errors.New("timeout on ReadDir") //nolint:err113
	}
}

// Open is an os.Open with a timeout in case you are calling it on a
// big NAS that might never reply Default is 60 seconds (1 minute);
// use OpenTimeout to tune the timeout.
func Open(name string) (*os.File, error) {
	var f *os.File

	var e error

	count.TimeFuncRunSuffix("open", func() {
		f, e = OpenTimeout(name, time.Minute)
	}, suffix)

	return f, e
}

// OpenTimeout is an fs.Open  with a timeout in case you are calling
// it on a big NAS that might never reply.
func OpenTimeout(name string, t time.Duration) (*os.File, error) {
	type res struct {
		f   *os.File
		err error
	}
	// start the Open
	resCh := make(chan res, two)
	go func(name string) {
		f, err := os.Open(name)
		resCh <- res{f, err}
	}(name)
	// now wait for it
	select {
	case result := <-resCh:
		count.IncrSuffix("open-ok", suffix)

		return result.f, result.err
	case <-time.After(t):
		count.IncrSuffix("open-timeout", suffix)

		return nil, errors.New("timeout on Open") //nolint:err113
	}
}

// Lstat is an os.Lstat with a timeout. Default is 60 seconds (1 minute);
// use LstatTimeout to tune the timeout.
func Lstat(name string) (os.FileInfo, error) {
	var fi os.FileInfo

	var e error

	count.TimeFuncRunSuffix("lstat", func() {
		fi, e = LstatTimeout(name, time.Minute)
	}, suffix)

	return fi, e
}

// LstatTimeout is an os.Lstat with a timeout parameter. Tmed out
// calls will leave a go routine and OS thread around.
func LstatTimeout(name string, t time.Duration) (os.FileInfo, error) {
	type res struct {
		fi  os.FileInfo
		err error
	}
	// start the Lstat
	resCh := make(chan res, two)
	go func(name string) {
		fi, err := os.Lstat(name)
		resCh <- res{fi, err}
	}(name)
	select {
	case result := <-resCh:
		count.IncrSuffix("lstat-ok", suffix)

		return result.fi, result.err
	case <-time.After(t):
		count.IncrSuffix("lstat-timeout", suffix)

		return nil, errors.New("timeout on Lstat") //nolint:err113
	}
}
