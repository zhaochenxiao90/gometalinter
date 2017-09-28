package regressiontests

import (
	"fmt"
	"testing"

	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
)

func TestGoType(t *testing.T) {
	t.Parallel()

	dir := fs.NewDir(t, "test-gotype",
		fs.WithFile("file.go", fileContent("root")),
		fs.WithDir("sub",
			fs.WithFile("file.go", fileContent("sub"))),
		fs.WithDir("excluded",
			fs.WithFile("file.go", fileContent("excluded"))))
	defer dir.Remove()

	expected := Issues{
		{Linter: "gotype", Severity: "error", Path: "file.go", Line: 4, Col: 6, Message: "foo declared but not used"},
		{Linter: "gotype", Severity: "error", Path: "sub/file.go", Line: 4, Col: 6, Message: "foo declared but not used"},
	}
	actual := RunLinter(t, "gotype", dir.Path(), "--skip=excluded")
	assert.Equal(t, expected, actual)
}

func fileContent(pkg string) string {
	return fmt.Sprintf(`package %s

func badFunction() {
	var foo string
}
	`, pkg)
}

func TestGoTypeWithMultiPackageDirectoryTest(t *testing.T) {
	t.Parallel()

	dir := fs.NewDir(t, "test-gotype",
		fs.WithFile("file.go", fileContent("root")),
		fs.WithFile("file_test.go", fileContent("root_test")))
	defer dir.Remove()

	// Expect only one issue because the other file is in an external package and
	// requires `gotype -x`
	expected := Issues{
		{Linter: "gotype", Severity: "error", Path: "file.go", Line: 4, Col: 6, Message: "foo declared but not used"},
	}
	actual := RunLinter(t, "gotype", dir.Path())
	assert.Equal(t, expected, actual)
}
