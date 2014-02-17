package util_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	//. "github.com/cf-platform-eng/cf-go-buildpack/util"
	"github.com/cf-platform-eng/cf-go-buildpack/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Archive", func() {
	var baseDir string
	var outputDir string
	var tarFile *os.File

	BeforeEach(func() {
		baseDir, _ = ioutil.TempDir("", "cf-go-buildpack-archive")
		outputDir, _ = ioutil.TempDir("", "cf-go-buildpack-output")
		tarFile, _ = ioutil.TempFile("", "cf-go-buildpack-tarfile")
	})

	AfterEach(func() {
		tarFile.Close()
		os.RemoveAll(tarFile.Name())
		os.RemoveAll(baseDir)
		os.RemoveAll(outputDir)
	})

	Describe("Roundtripping an archive", func() {
		var buf *bytes.Buffer

		JustBeforeEach(func() {
			buf = new(bytes.Buffer)
			util.CreateTarGz(buf, baseDir)
		})

		Context("With a single file", func() {
			BeforeEach(func() {
				createFile(filepath.Join(baseDir, "test.txt"), "testing")
			})

			It("Is successful", func() {
				Ω(buf.Len()).Should(BeNumerically(">", 0))
				length, err := io.Copy(tarFile, buf)
				if err != nil {
					fmt.Println(err.Error())
				}
				Ω(err).Should(BeNil())
				Ω(length).Should(BeNumerically(">", 0))
				tarFile.Close()
				createdFile, err := os.Open(tarFile.Name())
				util.ExtractTarGz(outputDir, createdFile, false)
				recursiveCompare(baseDir, filepath.Join(outputDir, filepath.Base(baseDir)))
			})
		})

		Context("With multiple files", func() {
			BeforeEach(func() {
				createFile(filepath.Join(baseDir, "test1.txt"), "testing1")
				createFile(filepath.Join(baseDir, "test2.txt"), "testing2")
				createFile(filepath.Join(filepath.Join(baseDir, "test"), "test3.txt"), "testing3")
			})

			It("Is successful", func() {
				Ω(buf.Len()).Should(BeNumerically(">", 0))
				length, err := io.Copy(tarFile, buf)
				if err != nil {
					fmt.Println(err.Error())
				}
				Ω(err).Should(BeNil())
				Ω(length).Should(BeNumerically(">", 0))
				tarFile.Close()
				createdFile, err := os.Open(tarFile.Name())
				util.ExtractTarGz(outputDir, createdFile, true)
				recursiveCompare(baseDir, outputDir)
			})
		})
	})
})

func recursiveCompare(source string, target string) {
	visit := func(path string, f os.FileInfo, err error) error {
		sourceFilename, _ := filepath.Rel(source, path)
		targetPath := filepath.Join(target, sourceFilename)
		if !f.IsDir() {
			sourceFile, err := ioutil.ReadFile(path)
			targetFile, err := ioutil.ReadFile(targetPath)
			if err != nil {
				Fail(err.Error())
			}

			Ω(targetFile).Should(BeEquivalentTo(sourceFile))
		}

		return nil
	}

	filepath.Walk(source, visit)
}

func createFile(dest string, content string) {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		Fail(err.Error())
	}

	file, err := os.Create(dest)
	if err != nil {
		Fail(err.Error())
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString(content)
	w.Flush()
}
