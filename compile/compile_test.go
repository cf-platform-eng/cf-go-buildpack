package compile_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	. "github.com/cf-platform-eng/cf-go-buildpack/compile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compile", func() {
	var buildDir string
	var cacheDir string

	BeforeEach(func() {
		buildDir, _ = ioutil.TempDir("", "cf-go-buildpack-builddir")
		cacheDir, _ = ioutil.TempDir("", "cf-go-buildpack-cachedir")
	})

	AfterEach(func() {
		os.RemoveAll(buildDir)
		os.RemoveAll(cacheDir)
	})

	Describe("Detecting a Go application", func() {
		Context("There is a .go file present in the root", func() {
			It("Should download and extract go", func() {
				Compile(buildDir, cacheDir, DetectVersion(buildDir))
				_, err := os.Stat(filepath.Join(cacheDir, DetectVersion(buildDir).Version))
				Ω(err).Should(BeNil())
			})
		})
	})
})
