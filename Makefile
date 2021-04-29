SRC = src/main/main.go #Add more src files -> scr/YOURPACKAGE/FILENAME.go

# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=src/bin/

#DONTCOMPILE = ls -1 src/*/*.go | grep -v _test.go

all: test build

build:
	$(GOBUILD) ./...

run:
	go run $(SRC)
	
test: 
	go test  -v ./.../shapeitup
	
clean:
	go clean -testcache
	rm -f $(BINARY_NAME)