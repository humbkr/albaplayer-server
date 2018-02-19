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

# Copy config file for each os.
cp ../alba.yml linux
cp ../alba.yml macos
cp ../alba.yml windows