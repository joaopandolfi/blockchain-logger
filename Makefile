
GOPATH ?= ~/go
# https://medium.com/mabar/today-i-learned-fix-go-get-private-repository-return-error-reading-sum-golang-org-lookup-93058a058dd8

GO ?= go
GOGET ?= $(GO) get
GOBUILD ?= $(GO) build
GOMOD ?= $(GO) mod
GOTEST ?= $(GO) test
GOTOOL ?= $(GO) tool
GORUN ?= $(GO) run

SRCPATH ?= src

GOIMPORTS ?= $(GORUN) golang.org/x/tools/cmd/goimports
GOLANGCI_LINT := golangci-lint --timeout 10m
GOLINT := $(GORUN) golang.org/x/lint/golint

GIT ?= git
GITDIFF ?= $(GIT) diff


################################################################################
## Go Tools
################################################################################

setup ::
	@echo "==> Setup: installing tools"
	go install github.com/rakyll/gotest@latest
	
################################################################################
## Go targets
################################################################################

.PHONY: build test cover

build:
	cd $(SRCPATH)
	rm -rf dist
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux go build -o dist/ ./...

test:
	cd $(SRCPATH) && go test -race -coverpkg= -coverprofile=coverage.out ./...

cover: cover/text

cover/html:
	cd $(SRCPATH) && $(GOTOOL) cover -html=coverage.out

cover/text:
	cd $(SRCPATH) && $(GOTOOL) cover -func=coverage.out

.PHONY: run
run:
	$(GORUN) ./app/service


################################################################################
## Linters and formatters
################################################################################

.PHONY: goimports lint git/diff

lint:
	cd src && $(GOLANGCI_LINT) run -c ./.golangci.yml ./...

# go get installation aren't guaranteed to work. We recommend using binary installation.
# more on https://golangci-lint.run/usage/install/#ci-installation
lint/CI:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2
	golangci-lint run -c ./.golangci.yml ./...

git/diff:
	@if ! $(GITDIFF) --quiet; then \
		printf 'Found changes on local workspace. Please run this target and commit the changes\n' ; \
		exit 1; \
	fi

SOURCES := $(shell \
	find . -name '*.go' | \
	grep -Ev './(proto|protogen|third_party|vendor|dal|.history)/' | \
	xargs)

goimports:
	$(GOIMPORTS) -w $(SOURCES)

################################################################################
## Metrics
################################################################################

.PHONY: metrics/up
metrics/up:
	docker-compose -f ./.prometheus/docker-compose.yml up -d

.PHONY: metrics/down
metrics/down:
	docker-compose -f ./.prometheus/docker-compose.yml down