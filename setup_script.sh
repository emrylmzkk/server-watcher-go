#!/bin/bash

# =================================================================
# Server Watcher Go - Setup & Cleanup Script (V2)
# Description: Installs NVM, Node 20, PM2, and Go (arch-based).
# Also manages existing PM2 backend processes.
# =================================================================

# Colors for better logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
GO_VERSION="1.26.1"
NODE_VERSION="20"
PROJECT_NAME="server-watcher-app"
INSTALL_DIR="/usr/local"

log "Starting setup script for $PROJECT_NAME..."

# 1. Architecture Detection
log "Detecting system architecture..."
ARCH=$(uname -m)
GO_BINARY_ARCH=""

case "$ARCH" in
    x86_64)
        GO_BINARY_ARCH="amd64"
        ;;
    aarch64|arm64)
        GO_BINARY_ARCH="arm64"
        ;;
    armv7l|armv6l)
        GO_BINARY_ARCH="armv6l"
        ;;
    *)
        error "Unsupported architecture: $ARCH. Manual installation required."
        exit 1
        ;;
esac

log "System Architecture: $ARCH -> Target Go Arch: $GO_BINARY_ARCH"

# 2. Node.js & NVM Setup
log "Checking for Node.js environment..."

# Check if NVM is already installed
export NVM_DIR="$HOME/.nvm"
if [ -s "$NVM_DIR/nvm.sh" ]; then
    log "NVM already installed. Loading NVM..."
    \. "$NVM_DIR/nvm.sh"
elif command -v nvm &> /dev/null; then
    log "NVM detected in PATH."
else
    log "NVM not found. Installing NVM..."
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
    
    # Load NVM for the current session
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
    
    success "NVM installed and loaded."
fi

# Install Node 20
if ! command -v node &> /dev/null || [[ $(node -v) != v20* ]]; then
    log "Installing Node.js version $NODE_VERSION..."
    nvm install "$NODE_VERSION"
    nvm use "$NODE_VERSION"
    nvm alias default "$NODE_VERSION"
    success "Node.js $(node -v) installed."
else
    log "Node.js $(node -v) is already installed."
fi

# 3. PM2 Setup
if ! command -v pm2 &> /dev/null; then
    log "PM2 not found. Installing PM2 globally via NPM..."
    npm install -g pm2
    success "PM2 installed successfully."
else
    log "PM2 is already installed: $(pm2 -v)"
fi

# 4. Go Installation Check/Install
if ! command -v go &> /dev/null; then
    log "Go not found. Proceeding with installation..."
    
    GO_TAR="go${GO_VERSION}.linux-${GO_BINARY_ARCH}.tar.gz"
    GO_URL="https://golang.org/dl/${GO_TAR}"
    
    log "Downloading Go $GO_VERSION from $GO_URL..."
    if wget -q "$GO_URL" -O "/tmp/$GO_TAR" || curl -L "$GO_URL" -o "/tmp/$GO_TAR"; then
        log "Download successful. Extracting to $INSTALL_DIR..."
        
        # Remove old installation if exists
        sudo rm -rf "${INSTALL_DIR}/go"
        
        # Extract
        if sudo tar -C "$INSTALL_DIR" -xzf "/tmp/$GO_TAR"; then
            success "Go extracted successfully."
            
            # Update PATH for current session
            export PATH=$PATH:${INSTALL_DIR}/go/bin
            
            # Add to profile for future sessions
            if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
                echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                log "Added Go to ~/.bashrc"
            fi
            
            success "Go version $(go version) installed."
        else
            error "Failed to extract Go tarball."
            exit 1
        fi
        rm "/tmp/$GO_TAR"
    else
        error "Failed to download Go. Check connection or version: $GO_VERSION"
        exit 1
    fi
else
    log "Go is already installed: $(go version)"
    export PATH=$PATH:/usr/local/go/bin
fi

# 5. PM2 Backend Cleanup
log "Checking for existing backend process: $PROJECT_NAME"
if pm2 list | grep -q "$PROJECT_NAME"; then
    warn "Existing PM2 process found for '$PROJECT_NAME'. Removing it..."
    pm2 stop "$PROJECT_NAME" &> /dev/null
    pm2 delete "$PROJECT_NAME" &> /dev/null
    success "Successfully removed process '$PROJECT_NAME' from PM2."
else
    log "No existing PM2 process found for '$PROJECT_NAME'."
fi

# Final Summary
echo -e "\n--------------------------------------------"
success "Full Setup Complete!"
log "Node: $(node -v)"
log "NPM:  $(npm -v)"
log "PM2:  $(pm2 -v)"
log "Go:   $(go version)"
log "Environment is ready for $PROJECT_NAME."
warn "IMPORTANT: Run 'source ~/.bashrc' or restart your terminal to apply all PATH changes."
echo -e "--------------------------------------------\n"
