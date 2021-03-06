REPO=grafana-provisioner
CONTAINER=quay.io/mad01/grafana-provisioner
VERSION ?= $(shell ./hacks/git-version)
LD_FLAGS="-X main.Version=$(VERSION) -w -s -extldflags \"-static\" "

$( shell mkdir -p _bin )
$( shell mkdir -p _release )


default: format format-verify build-dev

clean:
	@rm -r _bin _release

test: format-verify
	@echo "----- running tests -----"
	@go test -v -i $(shell go list ./... | grep -v '/vendor/')
	@go test -v $(shell go list ./... | grep -v '/vendor/')

install:
	@GOBIN=$(GOPATH)/bin && go install -v -ldflags $(LD_FLAGS) 

build: build-release
build-dev:
	@echo "----- running dev build-----"
	@export GOBIN=$(PWD)/_bin && go install -v -ldflags $(LD_FLAGS) 

build-release:
	@echo "----- running release build -----"
	@go build -v -o _release/$(REPO) -ldflags $(LD_FLAGS) 

container: container-build
container-build:
	@docker build -t $(CONTAINER):$(VERSION) --file Dockerfile .

container-push:
	@docker push $(CONTAINER):$(VERSION)

dep-ensure:
	@dep ensure -v -vendor-only

setup-deps:
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/tools/cmd/goimports

format:
	@echo "----- running gofmt -----"
	@gofmt -w -s *.go
	@echo "----- running goimports -----"
	@goimports -w *.go

format-verify:
	@echo "----- running gofmt verify -----"
	@hacks/verify-gofmt
	@echo "----- running goimports verify -----"
	@hacks/verify-goimports

lines-diff:
	@echo "----- business logic lindes of code -----"
	@git ls |grep .go |grep -v test.go|xargs cat |wc -l
	@echo "----- test lines of code -----"
	@git ls |grep .go |grep -e test.go|xargs cat |wc -l
