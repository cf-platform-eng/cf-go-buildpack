package detect

import (
	"io"
	"io/ioutil"
	"strings"
)

func Detect(writer io.Writer, buildDir string) int {
	files, _ := ioutil.ReadDir(buildDir)
	for _, file := range files {
		if strings.Contains(file.Name(), ".go") {
			writer.Write([]byte("go"))
			return 0
		}
	}
	writer.Write([]byte(""))
	return 1
}
