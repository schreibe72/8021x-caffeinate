all: 8021x-caffeinate.x86_64 8021x-caffeinate.arm64

clean:
	find . -type f -a \( -name 8021x-caffeinate.arm64 -o -name 8021x-caffeinate.x86_64 \) -delete	

8021x-caffeinate.x86_64: $(wildcard *.go)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 SDKROOT=$(xcrun --sdk macosx --show-sdk-path) go build -ldflags "-X main.version=`git describe --tags HEAD`" -o build/8021x-caffeinate.app/Contents/Resources/8021x-caffeinate.x86_64

8021x-caffeinate.arm64: $(wildcard *.go)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 SDKROOT=$(xcrun --sdk macosx --show-sdk-path) go build -ldflags "-X main.version=`git describe --tags HEAD`" -o build/8021x-caffeinate.app/Contents/Resources/8021x-caffeinate.arm64

.PHONY: all clean
