package files

import (
	"io"
	"io/fs"
	"os"
)

// FileSystemOperations is an interface to bring functionality of required file system operations
// Note: missing functionalities can be added and implemented on demand, functions must be one-to-one mapping from `os` package
type FileSystemOperations interface {
	OpenFile(name string, flag int, perm fs.FileMode) (FileIO, error)

	// Open os.Open
	Open(path string) (FileIO, error)

	// Create os.Create
	Create(name string) (FileIO, error)

	//Remove os.Remove
	Remove(path string) error

	// Move os.Rename moving the file to newPath from an oldPath
	Move(oldPath string, newPath string) error

	// MkdirAll os.MkdirAll
	MkdirAll(path string, perm fs.FileMode) error

	// Stat os.Stat file stats
	Stat(name string) (fs.FileInfo, error)
}

// FileIO is an interface to bring functionality of required file IO operations
// Note: missing functionalities can be added and implemented on demand
type FileIO interface {
	io.Reader
	io.Writer
	io.Closer
}

// OsFileSystemOperations FileSystemOperations implementation, wrapper for functionalities from `os` package
type OsFileSystemOperations struct{}

func (fileSystem *OsFileSystemOperations) Open(path string) (FileIO, error) {
	if file, err := os.Open(path); err == nil {
		return &OsFile{file}, nil
	} else {
		return nil, err
	}
}

func (fileSystem *OsFileSystemOperations) Remove(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	} else {
		return nil
	}
}

func (fileSystem *OsFileSystemOperations) OpenFile(path string, flag int, perm fs.FileMode) (FileIO, error) {
	if file, err := os.OpenFile(path, flag, perm); err == nil {
		return &OsFile{file}, nil
	} else {
		return nil, err
	}
}

func (fileSystem *OsFileSystemOperations) Create(name string) (FileIO, error) {
	if file, err := os.Create(name); err == nil {
		return &OsFile{file}, nil
	} else {
		return nil, err
	}
}

func (fileSystem *OsFileSystemOperations) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fileSystem *OsFileSystemOperations) Move(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (fileSystem *OsFileSystemOperations) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}


// OsFile wrapper type for os.File
type OsFile struct {
	file *os.File
}

func (osFile *OsFile) Read(p []byte) (n int, err error) {
	return osFile.file.Read(p)
}

func (osFile *OsFile) Write(p []byte) (n int, err error) {
	return osFile.file.Write(p)
}

func (osFile *OsFile) Close() error {
	return osFile.file.Close()
}
