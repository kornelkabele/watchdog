# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MAIN=./cmd/watchdog/main.go
CFG=./internal/cfg/config.yml
BINPATH=./bin
BINARY_NAME=$(BINPATH)/watchdog.exe
BINARY_PI=$(BINPATH)/watchdog-pi

all: test build
build: 
		$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) -v $(MAIN)
		cp -u $(CFG) $(BINPATH)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_PI)
		rm -f $(BINPATH)/*.log
		rm -rf $(BINPATH)/images
run:
		$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) -v $(MAIN)
		cp $(CFG) $(BINPATH)
		export $$(cat .secrets | tr -d '\r' | xargs) && $(BINARY_NAME)
deps:
		$(GOGET) github.com/secsy/goftp@latest
		$(GOGET) github.com/jordan-wright/email@latest
		$(GOGET) github.com/disintegration/imaging@latest
		$(GOGET) gopkg.in/yaml.v2
vendor:
		$(GOGET) mod vendor

# Cross compilation
build-pi:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -ldflags="-s -w" -o $(BINARY_PI) -v $(MAIN)
		cp -u $(CFG) $(BINPATH)
