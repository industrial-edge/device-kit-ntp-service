package utils

import (
	"io/fs"
	"os"
)

// FileSystemOperations is an interface to bring functionality of necessary file system operations for `NtpLastSyncFileMigration`
// Note: missing functionalities can be added and implemented on demand, functions must be one-to-one mapping from `os` package
type FileSystemOperations interface {
	// Move Migration requires moving the file(s) to newPath from oldPath
	Move(oldPath string, newPath string) error

	// MkdirAll If new target folder doesn't exist migration needs to create required folders
	MkdirAll(path string, perm fs.FileMode) error

	// Stat Retrieves file stat to check if file exist
	Stat(name string) (fs.FileInfo, error)
}

// OsFileSystemOperations FileSystemOperations implementation, wrapper for functionalities from `os` package
type OsFileSystemOperations struct{}

func (file *OsFileSystemOperations) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (file *OsFileSystemOperations) Move(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (file *OsFileSystemOperations) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
