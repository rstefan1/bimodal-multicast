BINDIR ?= $(CURDIR)/bin

# Run tests
test: generate 
	ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race \
		./src/...

# Run go fmt against code
fmt:
	go fmt ./src/...
	
# Run go vet against code
vet:
	go vet ./src/...

# Generate code
generate:
	go generate ./src/...

lint:
	$(BINDIR)/golangci-lint run ./src/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install ./vendor/github.com/onsi/ginkgo/ginkgo
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(BINDIR) v1.10.2

