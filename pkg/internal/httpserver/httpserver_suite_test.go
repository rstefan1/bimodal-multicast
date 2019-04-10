package httpserver

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHTTPServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Server Suite")
}
