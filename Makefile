cmd := $(shell ls cmd/)

all: $(cmd)

Vendor= github.com/dearcode/netpi/vendor/
Project = $(Vendor)github.com/dearcode/doodle/service/debug.Project
GitHash = $(Vendor)github.com/dearcode/doodle/service/debug.GitHash
GitTime = $(Vendor)github.com/dearcode/doodle/service/debug.GitTime
GitMessage = $(Vendor)github.com/dearcode/doodle/service/debug.GitMessage


LDFLAGS += -X "$(Project)=github.com/dearcode/netpi"
LDFLAGS += -X "$(GitHash)=$(shell git log --pretty=format:'%H' -1)"
LDFLAGS += -X "$(GitTime)=$(shell git log --pretty=format:'%ct' -1)"
LDFLAGS += -X "$(GitMessage)=$(shell git log --pretty=format:'%cn %s %b' -1)"

source := $(shell ls -ld */|awk '$$NF !~ /bin\/|logs\/|config\/|_vendor\/|vendor\/|web\/|Godeps\/|docs\// {printf $$NF" "}')

golint:
	go get golang.org/x/lint/golint

staticcheck:
	go get honnef.co/go/tools/cmd/staticcheck

lint: golint staticcheck
	for path in $(source); do golint "$$path..."; done;
	for path in $(source); do gofmt -s -l -w $$path;  done;
	go vet ./...
	staticcheck ./...


clean:
	@rm -rf bin

.PHONY: $(cmd)

$(cmd):
	go build -o bin/$@ -ldflags '$(LDFLAGS)' cmd/$@/*


test:
	@for path in $(source); do echo "go test ./$$path"; go test "./"$$path; done;

