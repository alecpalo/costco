package filesystem

import (
	"io"
	"os"
)

// file is an interface for mocking the outputs from fs
type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
	Write(b []byte) (int, error)
}

// fs is an interface for mocking file system functions. This does not include
// close as I don't care about it.
type fs interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (file, error)
}

// filesystem is the implementation of fs with basic os functions.
type filesystem struct {
}

func (f *filesystem) OpenFile(name string, flag int, perm os.FileMode) (file, error) {
	return os.OpenFile(name, flag, perm)
}

// Open takes in a filepath and attempts to open it, returning
// a file and nil upon success. If the file cannot be found or there is an
// error while opening the file it will return an error.
func (f *filesystem) Open(name string) (file, error) {
	return os.Open(name)
}

// Stat takes in a filepath and attempts to get the file info from that path.
// If there is an error finding the filepath, it will return an error.
func (f *filesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Mkdir takes in a file path and attempts to create a directory with the given
// permissions. If the directory exists or cannot be made, an error is returned.
func (f *filesystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}
