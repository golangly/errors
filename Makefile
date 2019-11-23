PKGS := github.com/arikkfir/go-errors
SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))
GO := go

check: test vet gofmt unconvert ineffassign errcheck

test:
	$(GO) test $(PKGS)

vet:
	$(GO) vet $(PKGS)

gofmt:
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"

unconvert:
	$(GO) get github.com/mdempsky/unconvert
	unconvert -v $(PKGS)

ineffassign:
	$(GO) get github.com/gordonklaus/ineffassign
	find $(SRCDIRS) -name '*.go' | xargs ineffassign

errcheck:
	$(GO) get github.com/kisielk/errcheck
	errcheck $(PKGS)
