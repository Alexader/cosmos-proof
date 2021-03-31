APP_NAME = proof
GO  = GO111MODULE=on go

# build with verison infos
VERSION_DIR = github.com/meshplus/${APP_NAME}
BUILD_DATE = $(shell date +%FT%T)
GIT_COMMIT = $(shell git log --pretty=format:'%h' -n 1)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

GOLDFLAGS += -X "${VERSION_DIR}.BuildDate=${BUILD_DATE}"
GOLDFLAGS += -X "${VERSION_DIR}.CurrentCommit=${GIT_COMMIT}"
GOLDFLAGS += -X "${VERSION_DIR}.CurrentBranch=${GIT_BRANCH}"
GOLDFLAGS += -X "${VERSION_DIR}.CurrentVersion=${APP_VERSION}"

RED=\033[0;31m
GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m

.PHONY: test

help: Makefile
	@echo "Choose a command run:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'

## make install: Go install the project
install:
	$(GO) install -ldflags '${GOLDFLAGS}' ./cmd/${APP_NAME}
	@printf "${GREEN}Build cosmos proof tool successfully${NC}\n"

build:
	@mkdir -p bin
	$(GO) build -ldflags '${GOLDFLAGS}' ./cmd/${APP_NAME}
	@mv ./proof bin
	@printf "${GREEN}Build cosmos proof tool successfully!${NC}\n"