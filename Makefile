# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINPATH=./bin
BINARY_NAME=$(BINPATH)/watchdog.exe
BINARY_PI=$(BINPATH)/watchdog-pi

all: test build
build: 
		$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) -v
		cp ./config.yml $(BINPATH)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_PI)
		rm -f $(BINPATH)/config.yml
		rm -f $(BINPATH)/*.log
		rm -rf $(BINPATH)/images
run:
		$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) -v ./...
		cp ./config.yml $(BINPATH)
		export $$(cat .secrets | tr -d '\r' | xargs) && ./$(BINARY_NAME)
deps:
		$(GOGET) github.com/secsy/goftp
		$(GOGET) github.com/jordan-wright/email
		$(GOGET) github.com/disintegration/imaging

# Cross compilation
build-pi:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -ldflags="-s -w" -o $(BINARY_PI) -v
