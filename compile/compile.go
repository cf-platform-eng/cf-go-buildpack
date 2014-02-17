package compile

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/cf-platform-eng/cf-go-buildpack/util"
)

func Compile(buildDir string, cacheDir string, version *GoVersion) int {
	if version.Version == "arm" {
		util.PrintlnErr(" !     ARM is not supported")
		return 1
	}

	if e, _ := util.Exists(path.Join(cacheDir, version.Version, "go")); e {
		fmt.Println("-----> Using", version.Version)
	}

	os.MkdirAll(path.Join(cacheDir, version.Version), 0755)
	out, _ := os.Create(path.Join(cacheDir, version.Version, version.Name()))
	defer out.Close()

	resp, _ := http.Get(version.Url())
	defer resp.Body.Close()

	util.ExtractTarGz(path.Join(cacheDir, version.Version), resp.Body, true)

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

func (version *GoVersion) Url() string {
	return "http://go.googlecode.com/files/" + version.Name()
}

func (version *GoVersion) Name() string {
	arch := version.Architecture
	if version.Platform == "darwin" && !strings.HasPrefix(version.Version, "go1.0") && !strings.HasPrefix(version.Version, "go1.1") {
		arch = arch + "-osx10.8"
	}

	return version.Version + "." + version.Platform + "-" + arch + ".tar.gz"
}
