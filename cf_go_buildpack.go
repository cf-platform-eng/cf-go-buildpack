package main

import (
	"os"

	"github.com/cf-platform-eng/cf-go-buildpack/compile"
	"github.com/cf-platform-eng/cf-go-buildpack/detect"
	"github.com/cf-platform-eng/cf-go-buildpack/release"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "cf_go_buildpack"
	app.Usage = "Run your Go applications on Cloud Foundry!"
	app.Commands = []cli.Command{
		{
			Name:      "detect",
			ShortName: "d",
			Usage:     "is this a Go application?",
			Action: func(c *cli.Context) {
				exit := detect.Detect(os.Stdout, c.Args()[0])
				os.Exit(exit)
			},
		},
		{
			Name:      "compile",
			ShortName: "c",
			Usage:     "prepare Go runtime and compile application",
			Action: func(c *cli.Context) {
				compile.Compile()
			},
		},
		{
			Name:      "release",
			ShortName: "r",
			Usage:     "generate application metadata",
			Action: func(c *cli.Context) {
				release.Release()
			},
		},
	}

	app.Run(os.Args)
}
