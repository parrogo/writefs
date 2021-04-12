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
	assert := assert.New(t)
	require := require.New(t)

	t.Run("WriteFile", func(t *testing.T) {
		t.Run("open a files for write, writes buf, closes the file", func(t *testing.T) {
			testfs := mockfs.FS{}

			writer := mockfs.FileWriter{}
			writer.On("Write", data).Return(len(data), nil)
			writer.On("Close").Return(nil)
			testfs.On("OpenFile", "dir1/file2", mock.Anything, mock.Anything).Return(&writer, nil)

			n, err := writefs.WriteFile(&testfs, "dir1/file2", data)
			assert.Equal(len(data), n)
			assert.NoError(err)

			writer.AssertExpectations(t)
			testfs.AssertExpectations(t)
		})

		t.Run("return PathError for unvalid path", func(t *testing.T) {

			n, err := writefs.WriteFile(fixtureFS, "/", nil)
			assert.Zero(n)
			require.Error(err)
			assert.ErrorIs(err, fs.ErrInvalid)
			perr, isPathErr := err.(*fs.PathError)
			require.True(isPathErr)
			assert.Equal("WriteFile", perr.Op)
			assert.Equal("/", perr.Path)
			assert.Equal("WriteFile /: invalid argument name: not a valid path", perr.Error())
		})

		t.Run("wraps OpenFile errors if any", func(t *testing.T) {

			n, err := writefs.WriteFile(fixtureFS, "not-existent", nil)
			assert.Zero(n)
			require.Error(err)

			assert.ErrorIs(err, fs.ErrInvalid)
			perr, isPathErr := err.(*fs.PathError)
			require.True(isPathErr)
			assert.Equal("WriteFile: OpenFile", perr.Op)
			assert.Equal("not-existent", perr.Path)
			assert.Equal("WriteFile: OpenFile not-existent: invalid argument fsys: does not implement WriteFS", perr.Error())
		})
	})

	t.Run("OpenFile", func(t *testing.T) {
		t.Run("Call fsys.OpenFile when fsys implements writefs.WriteFS", func(t *testing.T) {
			t.Run("For Write", func(t *testing.T) {
				testfs := mockfs.FS{}
				testfs.On("OpenFile", "dir1/file2", mock.Anything, mock.Anything).Return(nil, nil)

				f, err := writefs.OpenFile(&testfs, "dir1/file2", int(writefs.WriteOnly), fs.FileMode(0644))
				assert.NoError(err)
				assert.Nil(f)

				testfs.AssertExpectations(t)
			})

			t.Run("For Read", func(t *testing.T) {
				testfs := mockfs.FS{}
				testfs.On("OpenFile", "dir1/file2", mock.Anything, mock.Anything).Return(nil, nil)

				f, err := writefs.OpenFile(&testfs, "dir1/file2", int(writefs.ReadOnly), fs.FileMode(0))
				assert.NoError(err)
				assert.Nil(f)

				testfs.AssertExpectations(t)
			})
		})

		t.Run("Call fsys.Open for os.O_RDONLY when fsys not implements writefs.WriteFS", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "dir1/file2", os.O_RDONLY, fs.FileMode(0))
			assert.NoError(err)
			require.NotNil(f)

			defer f.Close()

			buf := make([]byte, len(data))
			n, err := f.Read(buf)
			assert.NoError(err)
			assert.Equal(n, len(data))

			_, err = f.Write(buf)
			assert.Error(err)
			assert.True(errors.Is(err, fs.ErrInvalid))

			assert.Equal(data, buf)
		})

		t.Run("return PathError for unvalid path", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "/", os.O_RDONLY, fs.FileMode(0))
			assert.Nil(f)
			require.Error(err)
			assert.ErrorIs(err, fs.ErrInvalid)
			perr, isPathErr := err.(*fs.PathError)
			require.True(isPathErr)
			assert.Equal("OpenFile", perr.Op)
			assert.Equal("/", perr.Path)
			assert.Equal("OpenFile /: invalid argument name: not a valid path", perr.Error())
		})

		t.Run("Return fsys.OpenFile error if any", func(t *testing.T) {
			testfs := mockfs.FS{}

			testfs.On("OpenFile", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("expected"))
			f, err := writefs.WriteFile(&testfs, "notexists", data)
			assert.Equal("WriteFile: OpenFile notexists: expected", err.Error())
			assert.Zero(f)
			testfs.AssertExpectations(t)
		})

		t.Run("return original error for RO open", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "notexists", os.O_RDONLY, fs.FileMode(0))
			assert.True(errors.Is(err, fs.ErrNotExist))
			assert.Nil(f)
		})

		t.Run("return invalid for RO open for write", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "notexists", os.O_WRONLY, fs.FileMode(0))
			assert.True(errors.Is(err, fs.ErrInvalid))
			assert.Nil(f)
		})

		t.Run("return PathError for unvalid paths", func(t *testing.T) {
			f, err := writefs.OpenFile(fixtureFS, "/", os.O_RDONLY, fs.FileMode(0))
			_, ok := err.(*fs.PathError)
			assert.True(ok)
			assert.Nil(f)
		})
	})

}
