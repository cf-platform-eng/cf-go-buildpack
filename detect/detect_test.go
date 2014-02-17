package detect_test

import (
	"io/ioutil"
	"os"
	"path"
	. "github.com/cf-platform-eng/cf-go-buildpack/detect"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeWriter struct {
	buffer []byte
}

func (this *FakeWriter) Write(p []byte) (n int, err error) {
	this.buffer = p
	return len(p), nil
}

var _ = Describe("Detect", func() {
	var buildDir string
	var fakeWriter *FakeWriter

	BeforeEach(func() {
		buildDir, _ = ioutil.TempDir("", "cf_go_bp")
		fakeWriter = new(FakeWriter)
	})

	AfterEach(func() {
		os.RemoveAll(buildDir)
	})

	Describe("Detecting a Go application", func() {
		Context("There is a .go file present in the root", func() {
			var tmpGoFile *os.File

			JustBeforeEach(func() {
				tmpGoFile, _ = os.Create(path.Join(buildDir, "hello.go"))
				defer tmpGoFile.Close()
			})

			It("should output 'Go' to writer", func() {
				Detect(fakeWriter, buildDir)
				Ω(string(fakeWriter.buffer[:])).To(Equal("Go"))
			})

			It("should exit '0'", func() {
				exit := Detect(fakeWriter, buildDir)
				Ω(exit).To(Equal(0))
			})
		})

		Context("There is a Godeps file present in the root", func() {
			var tmpGoFile *os.File

			JustBeforeEach(func() {
				tmpGoFile, _ = os.Create(path.Join(buildDir, "Godeps"))
				defer tmpGoFile.Close()
			})

			It("should output 'Go' to writer", func() {
				Detect(fakeWriter, buildDir)
				Ω(string(fakeWriter.buffer[:])).To(Equal("Go"))
			})

			It("should exit '0'", func() {
				exit := Detect(fakeWriter, buildDir)
				Ω(exit).To(Equal(0))
			})
		})

		Context("There is a Godeps folder present in the root", func() {
			JustBeforeEach(func() {
				if err := os.MkdirAll(path.Join(buildDir, "Godeps"), 0777); err != nil {
					Fail(err.Error())
				}
			})

			It("should output 'Go' to writer", func() {
				Detect(fakeWriter, buildDir)
				Ω(string(fakeWriter.buffer[:])).To(Equal("Go"))
			})

			It("should exit '0'", func() {
				exit := Detect(fakeWriter, buildDir)
				Ω(exit).To(Equal(0))
			})
		})

		Context("There is not a '.go' or 'Godeps' file, or a 'Godeps' folder present in the root", func() {
			It("should output 'No Go' to Writer", func() {
				Detect(fakeWriter, buildDir)
				Ω(string(fakeWriter.buffer[:])).To(Equal("No Go"))
			})

			It("should exit '1'", func() {
				exit := Detect(fakeWriter, buildDir)
				Ω(exit).To(Equal(1))
			})
		})
	})

})
