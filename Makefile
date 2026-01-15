BINARY_NAME=video-publisher

.PHONY: all build run clean

# Default target
all: build

# Build the binary
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) .

# Run the application
run:
	go run .

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)

