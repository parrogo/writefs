package mock

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestFileWriter(t *testing.T) {
	data := []byte{0xca, 0xfe, 0xba, 0xbe}

	t.Run("Close", func(t *testing.T) {
		w := &FileWriter{}

		w.On("Close").Return(nil)

		err := w.Close()
		assert.NoError(t, err)

		w.AssertExpectations(t)
	})

	t.Run("Write", func(t *testing.T) {

		w := &FileWriter{}

		w.On("Write", data).Return(42, nil)

		n, err := w.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, 42, n)

		w.AssertExpectations(t)

	})

	t.Run("Read", func(t *testing.T) {
		w := &FileWriter{}
		w.On("Read", data).Return(42, nil)
		n, err := w.Read(data)
		assert.NoError(t, err)
		assert.Equal(t, 42, n)
		w.AssertExpectations(t)
	})

	t.Run("Stat", func(t *testing.T) {
		w := &FileWriter{}
		w.On("Stat").Return(nil, nil)
		info, err := w.Stat()
		assert.NoError(t, err)
		assert.Nil(t, info)

		w.AssertExpectations(t)
	})

}

func newMemDirInfo(name string) fs.FileInfo {
	tmp := fstest.MapFS{}
	tmp[name] = &fstest.MapFile{}
	file, err := tmp.Open(name)
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return info
}
