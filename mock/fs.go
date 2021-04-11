package mock

import (
	"io/fs"

	"github.com/parrogo/writefs"
	"github.com/stretchr/testify/mock"
)

// FS mocks
type FS struct {
	mock.Mock
}

// DoSomething is a method on MyMockedObject that implements some interface
// and just records the activity, and returns what the Mock object tells it to.
//
// In the real object, this method would do something useful, but since this
// is a mocked object - we're just going to stub it out.
//
// NOTE: This method is not being tested here, code that uses this object is.
func (fsys *FS) DoSomething(number int) (bool, error) {

	args := fsys.Called(number)
	return args.Bool(0), args.Error(1)

}

var (
	_ fs.StatFS     = &FS{}
	_ fs.ReadFileFS = &FS{}
	_ fs.SubFS      = &FS{}
	_ fs.ReadDirFS  = &FS{}
	_ fs.GlobFS     = &FS{}

	_ writefs.WriteFS = &FS{}
	//_ writefs.RemoveFS = &FS{}
	//_ writefs.MkDirFS  = &FS{}
)

// MkDir implements writefs.MkDirFS
func (fsys *FS) MkDir(name string, perm fs.FileMode) error {
	args := fsys.Called(name, perm)
	return args.Error(0)
}

// Remove implements writefs.RemoveFS
func (fsys *FS) Remove(name string) error {
	args := fsys.Called(name)
	return args.Error(0)
}

// OpenFile implements writefs.WriteFS
func (fsys *FS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	args := fsys.Called(name, flag, perm)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(writefs.FileWriter), args.Error(1)
}

// Stat implements fs.StatFS
func (fsys *FS) Stat(name string) (fs.FileInfo, error) {
	args := fsys.Called(name)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(fs.FileInfo), args.Error(1)
}

// ReadFile implements fs.ReadFileFS
func (fsys *FS) ReadFile(name string) ([]byte, error) {
	args := fsys.Called(name)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.([]byte), args.Error(1)
}

// Sub implements fs.SubFS
func (fsys *FS) Sub(dir string) (fs.FS, error) {
	args := fsys.Called(dir)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(fs.FS), args.Error(1)
}

// Open implements fs.FS
func (fsys *FS) Open(name string) (fs.File, error) {
	args := fsys.Called(name)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(fs.File), args.Error(1)
}

// ReadDir implements fs.ReadDirFS
func (fsys *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	args := fsys.Called(name)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.([]fs.DirEntry), args.Error(1)
}

// Glob implements fs.GlobFS
func (fsys *FS) Glob(pattern string) ([]string, error) {
	args := fsys.Called(pattern)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.([]string), args.Error(1)
}
