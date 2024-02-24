BINDIR ?= $(CURDIR)/../../bin

test: generate
#	@$(BINDIR)/ginkgo version
#	$(BINDIR)/ginkgo \
#		--randomize-all --randomize-suites --fail-on-pending \
#		--cover --trace --race -v \
#		./...

e2e-tests:
	docker build -t bmmc-maelstrom -f Dockerfile .
	docker run --rm bmmc-maelstrom

fmt:
	go fmt ./...

vet:
	go vet ./...

generate:
	go generate ./...

lint:
	@$(BINDIR)/golangci-lint version
	$(BINDIR)/golangci-lint run ./...
