PACKAGE=payhere
BUILDPATH ?= $(CURDIR)
BASE = $(BUILDPATH)
BIN = $(BASE)/bin
GOBIN = $(shell go env GOPATH)/bin

UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
	GOENV   ?= CGO_ENABLED=0 GOOS=linux
endif
GOBUILD = ${GOENV} go
ARROW  = $(shell printf "\033[34;1m▶\033[0m")
GOV = $(shell go version)
BUILDTAG=-tags 'release'

$(BIN):
	@mkdir -p $@

.PHONY: build
build: vendor 
	$(info $(ARROW) building executable ...)
	@cd $(BASE)/api && $(GOBUILD) build \
		$(BUILDTAG) \
		-o $(BIN)/$(PACKAGE)
	@cd $(BASE) && cp -fp config/.env.json bin/.
	@bin/${PACKAGE} -v

.PHONY: pkg.list
pkg.list:
ifeq (,$(wildcard $(CURDIR)/pkg.list))
	$(info $(ARROW) generate pkg.list file ...)
	@cd $(BASE) && go list -f '{{.Dir}}' ./... 2>&1 | grep -v "^go: " | grep -v "^$(PACKAGE)/vendor/" | grep -v nocompile | grep -v logs | grep -v "\/mock" | grep -v testset | grep -v template > pkg.list
endif 

.PHONY: vendor
vendor:
	$(info $(ARROW) $(GOV) $(GOBIN) retrieving dependencies ...)
ifneq (,$(wildcard $(CURDIR)/go.mod))
	$(info $(ARROW) no needs mod init ...)
else
	$(info $(ARROW) mod init ...)
	@cd $(BASE) && go mod init $(PACKAGE)
	@cd $(BASE) && go mod tidy
	$(MAKE) pkg.list
endif

.PHONY: clean
clean:	
	$(info $(ARROW) cleaning ...)
	@rm -rf bin
	@rm -rf mock
	@rm -rf vendor
	@rm -rf pkg.list
	@rm -rf go.mod go.sum

# $(MAKE) = make command