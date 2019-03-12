package gossip

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHTTPServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gossip Server Suite")
}
