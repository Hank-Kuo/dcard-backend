# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

ENTRY_POINT=./cmd/main.go
BINARY_NAME=./build/server
BINARY_WINDOW=$(BINARY_NAME)_window.exe
BINARY_DARWIN=$(BINARY_NAME)_darwin
BINARY_LINUX=$(BINARY_NAME)_linux

ifeq (test, $(firstword $(MAKECMDGOALS)))
  runargs := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(eval $(runargs):;@true)
endif

all: test build
run:
	$(GORUN) $(ENTRY_POINT)
build:
	$(GOBUILD) -o $(BINARY_NAME) $(ENTRY_POINT) -v
test:
	$(GOTEST) -v ./$(runargs)/...
clean:
	$(GOCLEAN)
	-rm -f $(BINARY_NAME)
	-rm -f $(BINARY_LINUX)
	-rm -f $(BINARY_DARWIN)
	-rm -f $(BINARY_WINDOW)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) $(ENTRY_POINT)
build-window:
	CGO_ENABLED=0 GOOS=window GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) $(ENTRY_POINT)
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) $(ENTRY_POINT)

docker-build-image:
	docker build -t dcard-backend . --no-cache 
docker-run:
	docker run -d -p 8080:8080 dcard-backend
docker-compose:
	docker-compose up -d