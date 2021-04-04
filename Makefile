GO111MODULE=on
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: test build

build: openrms-windows-amd64.exe openrms-linux-arm openrms-linux-amd64

test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)

deps:
	$(GOCMD) mod download

openrms: test
	$(GOBUILD) -o openrms github.com/qvistgaard/openrms/cmd/openrms

openrms-windows-amd64.exe: test
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $@ -v github.com/qvistgaard/openrms/cmd/openrms

openrms-linux-arm: test
	GOOS=linux GOARCH=arm $(GOBUILD) -o $@ -v github.com/qvistgaard/openrms/cmd/openrms

openrms-linux-amd64: test
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $@ -v github.com/qvistgaard/openrms/cmd/openrms


# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
