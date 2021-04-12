package mock

import (
	"io/fs"

	"github.com/stretchr/testify/mock"

	"github.com/parrogo/writefs"
)

// # TODO: add an example file

// FileWriter provides a struct based on testify/mock.Mock
// that implements the fs.Writer.
//
// The struct is used to test the writefs package itself,
// and could be generally used to test methods that
// expects a writefs.WriteFS instance, by providing a
// mocked FileWriter type that you can return from you mocked
// writefs.WriteFS objects.
type FileWriter struct {
	mock.Mock
}

var _ writefs.FileWriter = &FileWriter{}

// Close implements fs.Close
func (w *FileWriter) Close() error {
	args := w.Called()
	return args.Error(0)
}

// Write implements writefs.Write
func (w *FileWriter) Write(buf []byte) (int, error) {
	args := w.Called(buf)
	return args.Int(0), args.Error(1)

}

// Read implements fs.Read
func (w *FileWriter) Read(buf []byte) (int, error) {
	args := w.Called(buf)
	return args.Int(0), args.Error(1)
}

// Stat implements fs.File
func (w *FileWriter) Stat() (fs.FileInfo, error) {
	args := w.Called()
	res := args.Get(0)
	res2, _ := res.(fs.FileInfo)
	return res2, args.Error(1)
}
