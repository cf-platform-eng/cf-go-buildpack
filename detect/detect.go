package detect

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Detect(writer io.Writer, buildDir string) int {
	godeps := false
	gofile := false
	visit := func(path string, f os.FileInfo, err error) error {
		if strings.ToLower(f.Name()) == "godeps" {
			godeps = true
		}

		if strings.HasSuffix(f.Name(), ".go") {
			gofile = true
		}
		return nil
	}

	filepath.Walk(buildDir, visit)

	if gofile || godeps {
		io.WriteString(writer, "Go")
		return 0
	}

	io.WriteString(writer, "No Go")
	return 1
}
