package compile

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/cf-platform-eng/cf-go-buildpack/util"
)

func Compile(buildDir string, cacheDir string, version *GoVersion) int {
	if version.Version == "arm" {
		util.PrintlnErr(" !     ARM is not supported")
		return 1
	}

	sep := string(os.PathSeparator)

	if e, _ := exists(cacheDir + sep + version.Version + sep + "go"); e {
		fmt.Println("-----> Using", version.Version)
	}

	return 0
}

type GoVersion struct {
	Version      string // go version (e.g. 1.2)
	Platform     string // os platform (e.g. linux, darwin, windows, freebsd)
	Architecture string // os architecture (e.g. amd64, i386)
}

func DetectVersion(buildDir string) *GoVersion {
	version := new(GoVersion)
	version.Platform = runtime.GOOS
	version.Architecture = runtime.GOARCH
	version.Version = "go1.2"

	// TODO: Detect Version From Godeps

	return version
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func goUrl(version *GoVersion) string {
	arch := version.Architecture
	if version.Platform == "darwin" && !strings.HasPrefix(version.Version, "go1.0") && !strings.HasPrefix(version.Version, "go1.1") {
		arch = arch + "-osx10.8"
	}

	return "http://go.googlecode.com/files/" + version.Version + "." + version.Platform + "-" + arch + ".tar.gz"
}
