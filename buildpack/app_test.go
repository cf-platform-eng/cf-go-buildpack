package buildpack_test

import (
	. "github.com/cf-platform-eng/cf-go-buildpack/buildpack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	Describe("Cli should be created correctly", func() {
		Context("Without arguments", func() {
			It("Should contain at least one mapped variable", func() {
				app := CreateApp()
				Ω(app).ShouldNot(BeNil())
			})
		})
	})
})
