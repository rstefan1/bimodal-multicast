BINDIR ?= $(CURDIR)/bin

# Run tests
test: generate
	@$(BINDIR)/ginkgo version
	$(BINDIR)/ginkgo \
		--randomize-all --randomize-suites --fail-on-pending \
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

# Run golangci-lint
lint:
	@$(BINDIR)/golangci-lint version
	$(BINDIR)/golangci-lint run ./pkg/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@v2.1.4
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- -b $(BINDIR) v1.48.0
