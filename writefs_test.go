package writefs_test

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/parrogo/writefs"
	mockfs "github.com/parrogo/writefs/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed mock/fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "mock/fixtures")

func TestWriteFS(t *testing.T) {
	data := []byte("ciao")

	t.Run("OpenFile", func(t *testing.T) {
		t.Run("Call fsys.OpenFile when fsys implements writefs.WriteFS", func(t *testing.T) {
			testfs := mockfs.FS{}
			writer := mockfs.FileWriter{}
			writer.On("Write", data).Return(len(data), nil)
			writer.On("Close").Return(nil)
			testfs.On("OpenFile", "dir1/file2", mock.Anything, mock.Anything).Return(&writer, nil)

			n, err := writefs.WriteFile(&testfs, "dir1/file2", data)
			assert.NoError(t, err)
			assert.Equal(t, n, len(data))
			testfs.AssertExpectations(t)
		})

		t.Run("Return fsys.OpenFile error if any", func(t *testing.T) {
			testfs := mockfs.FS{}

			testfs.On("OpenFile", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("expected"))
			f, err := writefs.WriteFile(&testfs, "notexists", data)
			assert.Equal(t, "expected", err.Error())
			assert.Zero(t, f)
			testfs.AssertExpectations(t)
		})

		t.Run("defaults to Open for os.O_RDONLY", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "dir1/file2", os.O_RDONLY, fs.FileMode(0))
			require.NoError(t, err)

			defer f.Close()

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
			f, err := writefs.OpenFile(fixtureFS, "notexists", os.O_RDONLY, fs.FileMode(0))
			assert.True(t, errors.Is(err, fs.ErrNotExist))
			assert.Nil(t, f)
		})

		t.Run("return invalid for RO open for write", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "notexists", os.O_WRONLY, fs.FileMode(0))
			assert.True(t, errors.Is(err, fs.ErrInvalid))
			assert.Nil(t, f)
		})

		t.Run("return PathError for unvalid paths", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "/", os.O_RDONLY, fs.FileMode(0))
			_, ok := err.(*fs.PathError)
			assert.True(t, ok)
			assert.Nil(t, f)
		})
	})

}
