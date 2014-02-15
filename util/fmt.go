package util

import (
	"fmt"
	"os"
)

func PrintlnErr(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}
