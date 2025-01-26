# Build variables
BINARY_NAME=Space-invader
WINDOWS_BINARY=$(BINARY_NAME).exe
LINUX_BINARY=$(BINARY_NAME)
OUTPUT_DIR=executables

# Build flags
BUILD_FLAGS=-ldflags="-s -w"

# Windows build settings
WINDOWS_ENV=CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64

# Ensure output directory exists
$(shell mkdir -p $(OUTPUT_DIR))

.PHONY: all clean windows linux

# Build for all platforms
all: windows linux

# Build for Windows
windows:
	@echo "Building for Windows..."
	$(WINDOWS_ENV) go build $(BUILD_FLAGS) -o $(OUTPUT_DIR)/$(WINDOWS_BINARY)
	@echo "Windows build complete: $(OUTPUT_DIR)/$(WINDOWS_BINARY)"

# Build for Linux
linux:
	@echo "Building for Linux..."
	go build $(BUILD_FLAGS) -o $(OUTPUT_DIR)/$(LINUX_BINARY)
	@echo "Linux build complete: $(OUTPUT_DIR)/$(LINUX_BINARY)"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR)
	@echo "Cleanup complete"

# Help command
help:
	@echo "Available commands:"
	@echo "  make all      - Build for all platforms"
	@echo "  make windows  - Build for Windows (requires mingw-w64)"
	@echo "  make linux    - Build for Linux"
	@echo "  make clean    - Remove build artifacts"
