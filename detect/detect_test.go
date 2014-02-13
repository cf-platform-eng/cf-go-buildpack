package detect

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
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
		buildDir, _ = ioutil.TempDir("/tmp", "cf_go_bp")
		fakeWriter = new(FakeWriter)
	})

	AfterEach(func() {
		os.RemoveAll(buildDir)
	})

	Describe("Detecting a Go application", func() {
		Context("There is a .go file present in the root", func() {
			var tmpGoFile *os.File

			JustBeforeEach(func() {
				tmpGoFile, _ = os.Create(buildDir + string(os.PathSeparator) + "hello.go")
				defer tmpGoFile.Close()
			})

			It("should output 'go' to stdout", func() {
				Detect(fakeWriter, buildDir)

				Ω(tmpGoFile.Name()).To(Equal(buildDir + string(os.PathSeparator) + "hello.go"))
				Ω(fakeWriter.buffer).To(Equal([]byte("go")))
			})

			It("should exit '0'", func() {
				exit := Detect(fakeWriter, buildDir)
				Ω(exit).To(Equal(0))
			})
		})

		Context("There is not a .go file present in the root", func() {
			It("should output nothing to stdout", func() {
				Detect(fakeWriter, buildDir)

				Ω(fakeWriter.buffer).To(Equal([]byte("")))
			})

			It("should exit '1'", func() {
				exit := Detect(fakeWriter, buildDir)
				Ω(exit).To(Equal(1))
			})
		})
	})

})
