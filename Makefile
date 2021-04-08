SRC = src/main/main.go #Add more src files -> scr/YOURPACKAGE/FILENAME.go
TEST = scr/main/main_test.go #Add more test files -> scr/YOURPACKAGE/FILENAME_test.go

# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=src/bin/main
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(SRC)

run:
	$(GORUN) -v $(SRC)
	
test: 
	$(GOTEST) -v ./...
	
clean:
	rm -f $(BINARY_NAME)