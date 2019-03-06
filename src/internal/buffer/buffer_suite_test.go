package buffer

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuffers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buffer Suite Test")
}
