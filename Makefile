VERSION := $(shell git describe --tags)
COMMIT  := $(shell git log -1 --format='%H')

all: install

LD_FLAGS = -X github.com/stafihub/rtoken-relay-core/relay/cmd.Version=$(VERSION) \
	-X github.com/stafihub/rtoken-relay-core/relay/cmd.Commit=$(COMMIT) \

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	cd relay && go mod tidy && go mod download
	cd common && go mod tidy && go mod download

build:
	@echo " > \033[32mBuilding relay...\033[0m "
	cd relay && go build -mod readonly $(BUILD_FLAGS) -o ../build/relay

install:
	@echo " > \033[32mInstalling relay...\033[0m "
	cd relay && go install -mod readonly $(BUILD_FLAGS)

build-linux:
	@GOOS=linux GOARCH=amd64 cd relay && go build --mod readonly $(BUILD_FLAGS) -o ../build/relay

clean:
	@echo " > \033[32mCleanning build files ...\033[0m "
	rm -rf build
fmt :
	@echo " > \033[32mFormatting go files ...\033[0m "
	cd common && go fmt ./...
	cd relay && go fmt ./...
lint:
	cd relay && golangci-lint run ./... --skip-files ".+_test.go"

.PHONY: all lint test race msan tools clean build
