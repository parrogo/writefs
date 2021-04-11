package writefs_test

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/parrogo/writefs"
	mockfs "github.com/parrogo/writefs/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed mock/fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

func TestWriteFS(t *testing.T) {
	data := []byte{0xca, 0xfe, 0xba, 0xbe}
	roFS := fstest.MapFS{
		"adir2/afile2": &fstest.MapFile{Data: data},
	}
	t.Run("WriteFile", func(t *testing.T) {
		t.Run("Call openfile when fsys implements writefs.WriteFS", func(t *testing.T) {
			testfs := mockfs.FS{}
			writer := mockfs.FileWriter{}
			writer.On("Write", data).Return(len(data), nil)
			writer.On("Close").Return(nil)
			testfs.On("OpenFile", "adir2/afile2", mock.Anything, mock.Anything).Return(&writer, nil)

			n, err := writefs.WriteFile(&testfs, "adir2/afile2", data)
			assert.NoError(t, err)
			assert.Equal(t, n, len(data))
			testfs.AssertExpectations(t)
		})
		t.Run("return original error for", func(t *testing.T) {
			testfs := mockfs.FS{}

			testfs.On("OpenFile", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("expected"))
			f, err := writefs.WriteFile(&testfs, "notexists", data)
			assert.Equal(t, "expected", err.Error())
			assert.Zero(t, f)
			testfs.AssertExpectations(t)
		})
	})
	t.Run("OpenFile", func(t *testing.T) {
		t.Run("defaults to Open for os.O_RDONLY", func(t *testing.T) {
			f, err := writefs.OpenFile(roFS, "adir2/afile2", os.O_RDONLY, fs.FileMode(0))
			defer f.Close()
			assert.NoError(t, err)

			buf := make([]byte, len(data))
			n, err := f.Read(buf)
			assert.NoError(t, err)
			assert.Equal(t, n, len(data))

			_, err = f.Write(buf)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, fs.ErrInvalid))

			assert.Equal(t, data, buf)
		})
		t.Run("return original error for RO open", func(t *testing.T) {
			f, err := writefs.OpenFile(roFS, "notexists", os.O_RDONLY, fs.FileMode(0))
			assert.True(t, errors.Is(err, fs.ErrNotExist))
			assert.Nil(t, f)
		})

		t.Run("return invalid for RO open for write", func(t *testing.T) {
			f, err := writefs.OpenFile(roFS, "notexists", os.O_WRONLY, fs.FileMode(0))
			assert.True(t, errors.Is(err, fs.ErrInvalid))
			assert.Nil(t, f)
		})

		t.Run("return PathError for unvalid paths", func(t *testing.T) {
			f, err := writefs.OpenFile(roFS, "/", os.O_RDONLY, fs.FileMode(0))
			_, ok := err.(*fs.PathError)
			assert.True(t, ok)
			assert.Nil(t, f)
		})
	})

}
