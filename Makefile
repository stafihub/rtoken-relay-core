VERSION := $(shell git describe --tags)
COMMIT  := $(shell git log -1 --format='%H')

all: install

LD_FLAGS = -X github.com/stafihub/rtoken-relay-core/relay/cmd.Version=$(VERSION) \
	-X github.com/stafihub/rtoken-relay-core/relay/cmd.Commit=$(COMMIT) \

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	cd common && go mod tidy && go mod download
fmt :
	@echo " > \033[32mFormatting go files ...\033[0m "
	cd common && go fmt ./...

.PHONY: all lint test race msan tools clean build
