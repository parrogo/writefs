// Package writefs provides interface WriteFS that extend fs.FS to support
// write operations.
//
// The main type of the package is WriteFS interface. It inheriths fs.FS
// and add as single OpenFile method that allows to open files for write,
// and it works similarly to os.OpenFile.
//
// There are also two subpackages:
// * mock contains a WriteFS implementation based on github.com/stretchr/testify/mock
//   that simplify testing of your code.
// * test is similar to fs/fstest but allows you to check your writefs.WriteFS
//   implementation.
package writefs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Flag is a type that represents the
// values accepted by OpenFile function
// flag argument.
type Flag int

const (
	// ReadOnly flag opens the file read-only.
	ReadOnly = Flag(os.O_RDONLY)
	// WriteOnly flag opens the file write-only.
	WriteOnly = Flag(os.O_WRONLY)
	// ReadWrite flag opens the file read-write.
	ReadWrite = Flag(os.O_RDWR)
	// Append flag appends data to the file when writing.
	Append = Flag(os.O_APPEND)
	// Create flag creates a new file if none exists.
	Create = Flag(os.O_CREATE)
	// Exclusive flag, when used with Create flag, requires that the file must not exist.
	Exclusive = Flag(os.O_EXCL)
	// Synchronous flag opens for synchronous I/O.
	Synchronous = Flag(os.O_SYNC)
	// Truncate flag truncates regular writable file when opened.
	Truncate = Flag(os.O_TRUNC)
)

// WriteFS extends fs.FS interface to provide write operations
// on file systems.
// OpenFile method could be used to open files for write
// but also to create directories and delete files or directories.
type WriteFS interface {
	fs.FS
	// OpenFile is the generalized open call; It opens the named file with
	// specified flags (O_RDONLY etc.).
	// If the file does not exist, and the O_CREATE flag is passed, it is
	// created with mode perm. If successful, methods on the
	// returned File can be used for I/O. If there is an error, it will
	// be of type *fs.PathError.
	// When flag Create is set and perm is fs.ModeDir, a directory is created
	// with path `name`, creating parent directories as needed when missing.
	// When flag Truncate is set, but not WriteOnly nor ReadWrite, file or
	// directory at path `name` is deleted. If the directory is not empty,
	// any content will be deleted recursively.
	// On both these circumstances, the function returns a nil FileWriter
	// even in case of success.
	// If this default semantic of directory creation and deletion is not
	// sufficient or if your filesystem implementation support optimized
	// algorithm, you can implements writefs.RemoveFS or writefs.MkDirFS
	// that provides more control on the operations.
	OpenFile(name string, flag int, perm fs.FileMode) (FileWriter, error)
}

// FileWriter ...
type FileWriter interface {
	fs.File
	io.Writer
}

// ReadOnlyWriteFile ...
type ReadOnlyWriteFile struct {
	fs.File
}

func (f ReadOnlyWriteFile) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("file does not support write: %w", fs.ErrInvalid)
}

func openFileReadOnly(fsInst fs.FS, name string) (FileWriter, error) {
	file, err := fsInst.Open(name)
	if err != nil {
		return nil, err
	}
	return ReadOnlyWriteFile{file}, nil
}

// OpenFile ...
func OpenFile(fsInst fs.FS, name string, flag int, perm fs.FileMode) (w FileWriter, err error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	if fs, ok := fsInst.(WriteFS); ok {
		return fs.OpenFile(name, flag, perm)
	}

	if flag == os.O_RDONLY {
		return openFileReadOnly(fsInst, name)
	}

	return nil, fmt.Errorf("file system does not support write: %w", fs.ErrInvalid)
}

// WriteFile ...
func WriteFile(fsys fs.FS, name string, buf []byte) (n int, err error) {
	file, err := OpenFile(fsys, name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.FileMode(0644))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err = file.Write(buf)
	return
}
