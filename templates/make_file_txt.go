package templates

const (
	makeFileTxt = `
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build -mod=vendor
GOTEST=$(GOCMD) test -mod=vendor

NAME={{.Name}}

GOPACKAGE=github.com/geekymedic/neon

GITCOMMIT=$(shell git rev-parse HEAD)
#GITTAG=$(shell git rev-list --tags --max-count=1)
#GITTAG:=$(shell git describe --tags $(GITTAG))
GITCOMMITTIME=$(shell git log -1 --format=%cd --date=local)
PRONAME={{.System}}-system-{{.Typ}}-{{.Name}}

ldflags=-X $(GOPACKAGE)/version.GITCOMMIT=$(GITCOMMIT)
#ldflags:=$(ldflags) -X $(GOPACKAGE)/version.GITTAG=$(GITTAG)
ldflags:=$(ldflags) -X '$(GOPACKAGE)/version.GITCOMMITTIME=$(GITCOMMITTIME)'
ldflags:=$(ldflags) -X '$(GOPACKAGE)/version.PRONAME=$(PRONAME)'

all: test build

build: clean
	$(GOBUILD) -o target/$(PRONAME) -ldflags "$(ldflags)" main.go
	cp config/config.yml target/

test:
	$(GOTEST) -v ./...

clean:
	rm -fr target
`
)
