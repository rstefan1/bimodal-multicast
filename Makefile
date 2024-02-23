BINDIR ?= $(CURDIR)/bin

GINKGO_VERSION = $(shell go list -f '{{ .Version }}' -m github.com/onsi/ginkgo/v2)
GOLANGCI_LINT_VERSION=v1.52.2

test: generate
	@$(BINDIR)/ginkgo version
	$(BINDIR)/ginkgo \
		--randomize-all --randomize-suites --fail-on-pending \
		--cover --trace --race -v \
		./pkg/...

	make -C _examples/http fmt
	make -C _examples/maelstrom fmt

fmt:
	go fmt ./pkg/...

	make -C _examples/http fmt
	make -C _examples/maelstrom fmt

vet:
	go vet ./pkg/...

	make -C _examples/http vet
	make -C _examples/maelstrom vet

generate:
	go generate ./pkg/...

	make -C _examples/http generate
	make -C _examples/maelstrom generate

lint:
	@$(BINDIR)/golangci-lint version
	$(BINDIR)/golangci-lint run ./pkg/...

	make -C _examples/http lint
	make -C _examples/maelstrom lint

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install github.com/onsi/ginkgo/v2/ginkgo@$(GINKGO_VERSION)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- -b $(BINDIR) $(GOLANCI_LINT_VERSION)
