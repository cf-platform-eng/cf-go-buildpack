package release

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Release Suite")
}
