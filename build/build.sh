#!/usr/bin/env bash

# Clean previously generated files.
rm -rf linux macos windows
mkdir linux macos windows

# Build for Linux.
echo "Start build for Linux..."
env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o linux/alba ../main.go
echo "Finished."

# Build for MacOs.
echo "Start build for MacOs..."
env CC=o64-clang GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o macos/alba ../main.go
echo "Finished."

# Build for Windows.
echo "Start build for Windows..."
env CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o windows/alba.exe ../main.go
echo "Finished."

echo "Generate archives..."

# Copy config file for each os.
cp prod.alba.yml linux/alba.yml
cp prod.alba.yml macos/alba.yml
cp prod.alba.yml windows/alba.yml

# Copy web directory for each os.
if [ ! -f ../web/index.html ]; then
    echo "The frontend app seems to not be present in the /web directory, did you forget to build it?"
    exit 1
else
    cp -R ../web linux
    cp -R ../web macos
    cp -R ../web windows
fi

# TODO zip / tar files.

echo "Application archives generated."
exit 0
