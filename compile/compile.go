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

	os.MkdirAll(path.Join(cacheDir, version.Version), 0777)
	out, _ := os.Create(path.Join(cacheDir, version.Version, goFile(version)))
	defer out.Close()

	resp, _ := http.Get(goUrl(version))
	defer resp.Body.Close()

	// n, _ := io.Copy(out, resp.Body)

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

func goUrl(version *GoVersion) string {
	return "http://go.googlecode.com/files/" + goFile(version)
}

func goFile(version *GoVersion) string {
	arch := version.Architecture
	if version.Platform == "darwin" && !strings.HasPrefix(version.Version, "go1.0") && !strings.HasPrefix(version.Version, "go1.1") {
		arch = arch + "-osx10.8"
	}

	return version.Version + "." + version.Platform + "-" + arch + ".tar.gz"
}
