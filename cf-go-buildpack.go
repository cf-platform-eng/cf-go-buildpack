package main

import (
	"os"

	"github.com/cf-platform-eng/cf-go-buildpack/buildpack"
)

func main() {
	app := buildpack.CreateApp()
	app.Run(os.Args)
}
