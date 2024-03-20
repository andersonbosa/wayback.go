.PHONY: all clean build-linux build-macos build-windows generate-checksum

MODULE_NAME=github.com/andersonbosa/security-goat@latest
BINARY_NAME=wayback.go
RELEASE_PATH=../build
CURRENT_VERSION := $(shell cat ./version)

release-build: clean build-linux build-macos build-windows generate-checksum

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(RELEASE_PATH)/$(BINARY_NAME)-linux-amd64 main.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(RELEASE_PATH)/$(BINARY_NAME)-macos-amd64 main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(RELEASE_PATH)/$(BINARY_NAME)-windows-amd64.exe main.go

generate-checksum:
	cd $(RELEASE_PATH) && \
	shasum -a 256 $(BINARY_NAME)-linux-amd64 > $(BINARY_NAME)-linux-amd64.sha256 && \
	shasum -a 256 $(BINARY_NAME)-macos-amd64 > $(BINARY_NAME)-macos-amd64.sha256 && \
	shasum -a 256 $(BINARY_NAME)-windows-amd64.exe > $(BINARY_NAME)-windows-amd64.sha256

clean:
	rm -f $(RELEASE_PATH)/$(BINARY_NAME)-* $(RELEASE_PATH)/*.sha256

# Make the module available by running the go list command to prompt Go to update its index of modules with information about the module you’re publishing.
go-publish:
	GOPROXY=proxy.golang.org go list -m $(MODULE_NAME)
	go install $(MODULE_NAME) 

go-upload-pkg:
	go clean -modcache
	go mod tidy
	go mod download
	go mod verify
	go list -m all

tag-version:
	git tag -a v$(CURRENT_VERSION) -m "Release version: v$(CURRENT_VERSION)"
	git push origin v$(CURRENT_VERSION)

