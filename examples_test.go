package writefs_test

import (
	"embed"
	"fmt"
	"io/fs"
	
	"github.com/parrogo/writefs"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

// This example show how to use writefs.Func()
func ExampleFunc() {
	fmt.Println(writefs.Func())
	// Output: 42
}
