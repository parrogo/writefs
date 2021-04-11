package writefs_test

// # TODO: manage how to write an example using maybe a mock fs?

/*
import (
	"embed"
	"io/fs"

	"github.com/parrogo/writefs"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

// This example show how to use writefs.OpenFile()
func ExampleOpenFile(fsys writefs.WriteFS) {
	f, err := writefs.OpenFile(fsys, "filename", 0, fs.FileMode(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte("ciao"))
	// Output: 42
}
*/
