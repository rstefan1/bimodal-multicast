BINDIR ?= $(CURDIR)/bin

# Run tests
test: generate
	$(BINDIR)/ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race -v \
		./pkg/...

# Run go fmt against code
fmt:
	go fmt ./pkg/...
	
# Run go vet against code
vet:
	go vet ./pkg/...

# Generate code
generate:
	go generate ./pkg/...

lint:
	$(BINDIR)/golangci-lint run ./pkg/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install ./vendor/github.com/onsi/ginkgo/ginkgo
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(BINDIR) v1.10.2

