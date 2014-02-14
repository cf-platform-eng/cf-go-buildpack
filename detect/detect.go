package detect

import (
	"io"
	"io/ioutil"
	"strings"
)

func Detect(writer io.Writer, buildDir string) int {
	files, _ := ioutil.ReadDir(buildDir)
	match := false
	for _, file := range files {
		// Godeps files
		if strings.HasPrefix(strings.ToLower(file.Name()), "godeps") {
			match = true
		}

		// .go files
		if strings.HasSuffix(file.Name(), ".go") {
			match = true
		}
	}

	if match {
		io.WriteString(writer, "Go")
		return 0
	}

	io.WriteString(writer, "No Go")
	return 1
}
