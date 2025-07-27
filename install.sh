#!/bin/bash

# hab installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/grovesjosephn/hab/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Configuration
REPO_URL="https://github.com/grovesjosephn/hab.git"
BINARY_NAME="hab"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Check dependencies
check_dependencies() {
    local missing_deps=()
    
    if ! command -v git >/dev/null 2>&1; then
        missing_deps+=("git")
    fi
    
    if ! command -v go >/dev/null 2>&1; then
        missing_deps+=("go")
    fi
    
    if ! command -v make >/dev/null 2>&1; then
        missing_deps+=("make")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        error "Missing required dependencies: ${missing_deps[*]}"
    fi
}

# Clone and build hab
install_hab() {
    local temp_dir=$(mktemp -d)
    
    info "Cloning hab repository..."
    if ! git clone "$REPO_URL" "$temp_dir" --depth 1 --quiet; then
        error "Failed to clone repository"
    fi
    
    info "Building hab from source..."
    cd "$temp_dir"
    
    if ! make build >/dev/null 2>&1; then
        error "Failed to build hab"
    fi
    
    info "Installing hab to $INSTALL_DIR..."
    
    # Check if we need sudo
    if [ ! -w "$INSTALL_DIR" ]; then
        if command -v sudo >/dev/null 2>&1; then
            warning "Installing to $INSTALL_DIR requires sudo privileges"
            sudo cp .bin/hab "$INSTALL_DIR/$BINARY_NAME"
        else
            error "Cannot write to $INSTALL_DIR and sudo is not available"
        fi
    else
        cp .bin/hab "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Clean up
    cd /
    rm -rf "$temp_dir"
    
    success "hab installed to $INSTALL_DIR/$BINARY_NAME"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local installed_version=$($BINARY_NAME --version 2>/dev/null || echo "hab")
        success "Installation verified! $installed_version"
        info "Run '$BINARY_NAME --help' to get started"
        info "Run '$BINARY_NAME' to launch the interactive TUI"
    else
        warning "Binary installed but not found in PATH"
        info "You may need to add $INSTALL_DIR to your PATH"
        info "Or run: export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

# Main installation flow
main() {
    echo "🎯 hab - Terminal Habit Tracker Installer"
    echo "========================================="
    echo
    
    # Check dependencies
    check_dependencies
    
    # Check if already installed
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version=$($BINARY_NAME --version 2>/dev/null || echo "hab")
        warning "$BINARY_NAME is already installed: $current_version"
        echo -n "Do you want to reinstall? [y/N]: "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            info "Installation cancelled"
            exit 0
        fi
    fi
    
    # Install hab
    install_hab
    
    # Verify installation
    verify_installation
    
    echo
    success "Installation complete! 🎉"
    echo
    info "Get started with:"
    echo "  $BINARY_NAME new exercise    # Create your first habit"
    echo "  $BINARY_NAME exercise        # Log an activity"
    echo "  $BINARY_NAME                 # View your progress"
}

# Run the installer
main "$@"