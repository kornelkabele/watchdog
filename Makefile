# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFLAGS=-ldflags="-s -w"
MAIN=./cmd/watchdog/main.go
CFG=./config.yml
EXT=
ifeq (${GOOS},windows)
    EXT=.exe
endif
ifeq ($(OS),Windows_NT)
    EXT=.exe
endif
BINPATH=./bin
BINARY_NAME=$(BINPATH)/watchdog$(EXT)
BINARY_PI=$(BINPATH)/watchdog-pi
IMAGE_DIR=c:/temp/watchdog/images
LOG_DIR=c:/temp/watchdog/log

all: test build
build: 
		$(GOBUILD) $(GOFLAGS) -o $(BINARY_NAME) -v $(MAIN)
		@cp -u $(CFG) $(BINPATH)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		@rm -f $(BINARY_PI)
		@rm -rf $(BINPATH)/log
		@rm -rf $(BINPATH)/images
run:
		$(GOBUILD) $(GOFLAGS) -o $(BINARY_NAME) -v $(MAIN)
		@cp $(CFG) $(BINPATH)
		export $$(cat .secrets | tr -d '\r' | xargs) && $(BINARY_NAME)
docker-build:
		docker build -t watchdog .
		docker image prune --filter label=stage=builder -f
docker-run:
		docker run -it --name=watchdog --rm --env-file .secrets --mount type=bind,source=$(IMAGE_DIR),target=/images --mount type=bind,source=$(LOG_DIR),target=/log watchdog
docker-stop:
		docker stop --time=20 watchdog
docker-killall:
		docker kill $$(docker ps -a -q)
deps:
		$(GOGET) github.com/secsy/goftp@latest
		$(GOGET) github.com/jordan-wright/email@latest
		$(GOGET) github.com/disintegration/imaging@latest
		$(GOGET) gopkg.in/yaml.v2
vendor:
		$(GOGET) mod vendor

# Cross compilation
pi-build:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) $(GOFLAGS) -o $(BINARY_PI) -v $(MAIN)
		@cp -u $(CFG) $(BINPATH)
