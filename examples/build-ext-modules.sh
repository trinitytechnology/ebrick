#!/bin/bash

# Set the Go compiler
GO="go"

# Enable cgo
# export CGO_ENABLED=1

# Set the root directory for plugins
PLUGIN_DIR="./modules"

# Output directory for compiled plugins
BUILD_DIR="./cmd/mono/modules"

# Build flags
BUILD_FLAGS="-buildmode=plugin"

# Function to build all plugins
build_plugins() {
    # Ensure the build directory exists
    mkdir -p "$BUILD_DIR"

    # Find all subdirectories in the PLUGIN_DIR
    for PLUGIN in "$PLUGIN_DIR"/*; do
        if [ -d "$PLUGIN" ]; then
            PLUGIN_NAME=$(basename "$PLUGIN")
            echo "Building modules: $PLUGIN_NAME"
            # Compile the plugin
            ($GO build -x $BUILD_FLAGS -o "$BUILD_DIR/${PLUGIN_NAME}.so" "./$PLUGIN")
        fi
    done
}

# Function to clean up the build directory
clean_plugins() {
    echo "Cleaning up build directory..."
    rm -rf "$BUILD_DIR"
}

# Check command line arguments
if [ "$1" == "clean" ]; then
    clean_plugins
else
    build_plugins
fi
