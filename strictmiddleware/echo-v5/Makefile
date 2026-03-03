MINIMUM_GO_MINOR := 25
CURRENT_GO_MINOR := $(shell go version | awk '{split($$3,a,"."); print a[2]+0}')
GO_VERSION_OK := $(shell [ $(CURRENT_GO_MINOR) -ge $(MINIMUM_GO_MINOR) ] && echo yes || echo no)

ifeq ($(GO_VERSION_OK),yes)

lint:
	$(GOBIN)/golangci-lint run ./...

lint-ci:
	$(GOBIN)/golangci-lint run ./... --output.text.path=stdout --timeout=5m

generate:
	go generate ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

tidy-ci:
	tidied -verbose

else

SKIP_MSG := @echo "Skipping echo-v5: requires Go 1.$(MINIMUM_GO_MINOR)+ (have $$(go version | awk '{print $$3}'))"

lint lint-ci generate test tidy tidy-ci:
	$(SKIP_MSG)

endif
