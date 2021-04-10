package writefs_test

import (
	"embed"
	"io/fs"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

// This example show how to use writefs.Func()
func ExampleOpenFile() {
	//fmt.Println(writefs.Func())
	// Output: 42
}
