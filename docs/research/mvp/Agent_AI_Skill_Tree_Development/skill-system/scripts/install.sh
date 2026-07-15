#!/bin/bash
# =============================================================================
# HelixKnowledge Skill Graph System - One-Command Installer
# =============================================================================
# Usage: ./install.sh [install-dir]
# Default install directory: /opt/skill-system
# 
# This script:
#   1. Checks dependencies (docker/podman, compose)
#   2. Creates install directory
#   3. Copies deployment files
#   4. Pulls container images
#   5. Runs database migrations
#   6. Creates systemd user service
#   7. Starts the stack
#   8. Prints status and next steps
# =============================================================================

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="${1:-/opt/skill-system}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SERVICE_NAME="skill-system"

# Logging functions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║         HelixKnowledge Skill Graph System Installer           ║"
    echo "║                                                               ║"
    echo "║  AI-powered skill tracking with auto-growth and validation   ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Check if command exists
check_command() {
    if command -v "$1" &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check for container runtime (docker or podman)
    if check_command docker; then
        CONTAINER_RUNTIME="docker"
        log_success "Found Docker: $(docker --version)"
    elif check_command podman; then
        CONTAINER_RUNTIME="podman"
        log_success "Found Podman: $(podman --version)"
    else
        log_error "Neither Docker nor Podman found. Please install one of them first."
        echo ""
        echo "  Docker: https://docs.docker.com/get-docker/"
        echo "  Podman: https://podman.io/getting-started/installation"
        exit 1
    fi
    
    # Check for compose plugin
    if $CONTAINER_RUNTIME compose version &> /dev/null; then
        COMPOSE_CMD="$CONTAINER_RUNTIME compose"
        log_success "Found Compose plugin"
    elif check_command docker-compose; then
        COMPOSE_CMD="docker-compose"
        log_success "Found docker-compose"
    elif check_command podman-compose; then
        COMPOSE_CMD="podman-compose"
        log_success "Found podman-compose"
    else
        log_error "No compose plugin found. Please install docker-compose or podman-compose."
        exit 1
    fi
    
    # Check for curl
    if ! check_command curl; then
        log_warn "curl not found. Health checks will be limited."
    fi
    
    # Check for jq
    if ! check_command jq; then
        log_warn "jq not found. JSON parsing will be limited."
    fi
    
    # Check systemd availability for user services
    if ! systemctl --user status &> /dev/null; then
        log_warn "systemd user services not available. Will use direct compose commands."
        USE_SYSTEMD=false
    else
        USE_SYSTEMD=true
        log_success "systemd user services available"
    fi
    
    # Check available resources
    if check_command free; then
        AVAILABLE_MEM=$(free -m | awk '/^Mem:/{print $7}')
        if [ "$AVAILABLE_MEM" -lt 2048 ]; then
            log_warn "Available memory is ${AVAILABLE_MEM}MB. Recommended: 4096MB+"
        else
            log_success "Available memory: ${AVAILABLE_MEM}MB"
        fi
    fi
    
    log_success "All required dependencies found"
}

# Create install directory
setup_install_dir() {
    log_info "Setting up installation directory: $INSTALL_DIR"
    
    if [ -d "$INSTALL_DIR" ]; then
        log_warn "Directory $INSTALL_DIR already exists"
        read -p "Overwrite existing installation? [y/N]: " confirm
        if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
            log_info "Installation cancelled"
            exit 0
        fi
        BACKUP_DIR="${INSTALL_DIR}.backup.$(date +%Y%m%d%H%M%S)"
        log_info "Backing up existing installation to $BACKUP_DIR"
        cp -r "$INSTALL_DIR" "$BACKUP_DIR"
    fi
    
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"/{scripts,config,docs,data/evidence,data/backups}
    
    # Copy files
    cp "$PROJECT_DIR/docker-compose.yml" "$INSTALL_DIR/"
    cp "$PROJECT_DIR/Dockerfile" "$INSTALL_DIR/"
    cp "$PROJECT_DIR/Makefile" "$INSTALL_DIR/" 2>/dev/null || true
    cp -r "$PROJECT_DIR/scripts/"* "$INSTALL_DIR/scripts/"
    cp -r "$PROJECT_DIR/config/"* "$INSTALL_DIR/config/" 2>/dev/null || true
    cp -r "$PROJECT_DIR/migrations" "$INSTALL_DIR/" 2>/dev/null || true
    
    # Create .env from example if it doesn't exist
    if [ ! -f "$INSTALL_DIR/.env" ]; then
        if [ -f "$PROJECT_DIR/.env.example" ]; then
            cp "$PROJECT_DIR/.env.example" "$INSTALL_DIR/.env"
            log_info "Created .env from template. Please review and customize it."
        else
            # Create minimal .env
            cat > "$INSTALL_DIR/.env" << 'EOF'
VERSION=1.0.0
DB_NAME=skilldb
DB_USER=skilluser
DB_PASSWORD=skillpassword
LOG_LEVEL=info
HTTP_PORT=8080
HTTP3_PORT=8443
ENABLE_HTTP3=true
EMBEDDING_PROVIDER=local
EMBEDDING_DIMENSION=768
AUTO_EXPAND_ENABLED=true
EOF
        fi
    fi
    
    # Make scripts executable
    chmod +x "$INSTALL_DIR/scripts/"*.sh
    
    log_success "Installation directory prepared"
}

# Pull container images
pull_images() {
    log_info "Pulling container images..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD pull
    log_success "Images pulled"
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."
    cd "$INSTALL_DIR"
    
    # Start database first
    $COMPOSE_CMD up -d db
    
    # Wait for database to be ready
    log_info "Waiting for database to be ready..."
    sleep 5
    
    local retries=30
    while [ $retries -gt 0 ]; do
        if $COMPOSE_CMD exec -T db pg_isready -U skilluser -d skilldb &> /dev/null; then
            log_success "Database is ready"
            break
        fi
        retries=$((retries - 1))
        sleep 2
    done
    
    if [ $retries -eq 0 ]; then
        log_error "Database failed to start"
        exit 1
    fi
    
    # Run migrations using skillctl or direct SQL
    if [ -f "$INSTALL_DIR/migrations/001_initial.up.sql" ]; then
        log_info "Applying initial schema..."
        $COMPOSE_CMD exec -T db psql -U skilluser -d skilldb < \
            "$INSTALL_DIR/migrations/001_initial.up.sql"
        log_success "Initial schema applied"
    fi
    
    log_success "Database migrations complete"
}

# Create systemd user service
create_systemd_service() {
    if [ "$USE_SYSTEMD" != true ]; then
        log_warn "Skipping systemd service creation"
        return
    fi
    
    log_info "Creating systemd user service..."
    
    # Determine compose command for service file
    if [ "$CONTAINER_RUNTIME" = "podman" ]; then
        COMPOSE_BIN="podman-compose"
    else
        COMPOSE_BIN="docker compose"
    fi
    
    # Create systemd user directory
    mkdir -p "$HOME/.config/systemd/user"
    
    # Create service file
    cat > "$HOME/.config/systemd/user/${SERVICE_NAME}.service" << EOF
[Unit]
Description=HelixKnowledge Skill Graph System
Documentation=https://github.com/helixdevelopment/skill-system
Requires=${CONTAINER_RUNTIME}.socket
After=${CONTAINER_RUNTIME}.socket network-online.target

[Service]
Type=simple
WorkingDirectory=${INSTALL_DIR}
Environment="COMPOSE_CMD=${COMPOSE_CMD}"
Environment="PATH=/usr/local/bin:/usr/bin:/bin"
Environment="CONTAINER_RUNTIME=${CONTAINER_RUNTIME}"

# Load environment from .env file
EnvironmentFile=${INSTALL_DIR}/.env

# Start command
ExecStartPre=-${COMPOSE_BIN} -f ${INSTALL_DIR}/docker-compose.yml pull
ExecStart=${COMPOSE_BIN} -f ${INSTALL_DIR}/docker-compose.yml up --remove-orphans

# Stop command
ExecStop=${COMPOSE_BIN} -f ${INSTALL_DIR}/docker-compose.yml down --timeout 30

# Restart policy
Restart=on-failure
RestartSec=10
StartLimitInterval=60
StartLimitBurst=3

# Graceful shutdown
TimeoutStopSec=60
KillSignal=SIGTERM

[Install]
WantedBy=default.target
EOF
    
    # Reload systemd
    systemctl --user daemon-reload
    
    # Enable service
    systemctl --user enable "$SERVICE_NAME"
    
    log_success "systemd user service created"
    log_info "Service file: $HOME/.config/systemd/user/${SERVICE_NAME}.service"
}

# Start the stack
start_stack() {
    log_info "Starting the stack..."
    cd "$INSTALL_DIR"
    
    if [ "$USE_SYSTEMD" = true ]; then
        systemctl --user start "$SERVICE_NAME"
        sleep 5
        
        if systemctl --user is-active "$SERVICE_NAME" &> /dev/null; then
            log_success "Stack started via systemd"
        else
            log_error "Failed to start stack via systemd"
            systemctl --user status "$SERVICE_NAME" --no-pager
            exit 1
        fi
    else
        $COMPOSE_CMD up -d
        sleep 5
        log_success "Stack started via compose"
    fi
}

# Wait for services to be healthy
wait_for_healthy() {
    log_info "Waiting for services to be healthy..."
    cd "$INSTALL_DIR"
    
    local retries=30
    while [ $retries -gt 0 ]; do
        local all_healthy=true
        
        # Check API health
        if curl -sf http://localhost:8080/health &> /dev/null; then
            log_success "API is healthy"
            break
        else
            all_healthy=false
        fi
        
        if [ "$all_healthy" = false ]; then
            retries=$((retries - 1))
            echo -n "."
            sleep 2
        fi
    done
    
    if [ $retries -eq 0 ]; then
        log_warn "Services may not be fully healthy yet. Check logs with: ./scripts/status.sh"
    else
        log_success "All services are healthy"
    fi
}

# Print status and next steps
print_next_steps() {
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║              Installation Complete!                           ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Installation Directory:${NC} $INSTALL_DIR"
    echo -e "${BLUE}Container Runtime:${NC} $CONTAINER_RUNTIME"
    echo -e "${BLUE}Compose Command:${NC} $COMPOSE_CMD"
    echo ""
    echo -e "${BLUE}Services:${NC}"
    echo "  API Server:     http://localhost:8080"
    echo "  HTTP/3 API:     https://localhost:8443 (UDP)"
    echo "  PostgreSQL:     localhost:5432"
    echo ""
    echo -e "${BLUE}Useful Commands:${NC}"
    if [ "$USE_SYSTEMD" = true ]; then
        echo "  Start:          systemctl --user start $SERVICE_NAME"
        echo "  Stop:           systemctl --user stop $SERVICE_NAME"
        echo "  Status:         systemctl --user status $SERVICE_NAME"
        echo "  Logs:           journalctl --user -u $SERVICE_NAME -f"
    fi
    echo "  Compose:        cd $INSTALL_DIR && $COMPOSE_CMD ps"
    echo "  Health:         curl http://localhost:8080/health"
    echo "  API Docs:       curl http://localhost:8080/api/v1/docs"
    echo ""
    echo -e "${BLUE}Management Scripts:${NC}"
    echo "  $INSTALL_DIR/scripts/start.sh   - Start the stack"
    echo "  $INSTALL_DIR/scripts/stop.sh    - Stop the stack"
    echo "  $INSTALL_DIR/scripts/status.sh  - Show stack status"
    echo "  $INSTALL_DIR/scripts/backup.sh  - Create backup"
    echo "  $INSTALL_DIR/scripts/migrate.sh - Run migrations"
    echo ""
    echo -e "${YELLOW}Next Steps:${NC}"
    echo "  1. Review configuration in: $INSTALL_DIR/.env"
    echo "  2. Access the API at: http://localhost:8080"
    echo "  3. Read the documentation: $INSTALL_DIR/docs/"
    echo "  4. Set up MCP integration: docs/MCP_INTEGRATION.md"
    echo ""
    echo -e "${YELLOW}Security Note:${NC}"
    echo "  - Change default database password in .env"
    echo "  - Set JWT_SECRET and API_KEY for production"
    echo "  - Configure firewall rules for ports 8080, 8443, 5432"
    echo ""
}

# Main installation flow
main() {
    print_banner
    
    log_info "Starting installation..."
    log_info "Install directory: $INSTALL_DIR"
    
    check_dependencies
    setup_install_dir
    pull_images
    run_migrations
    create_systemd_service
    start_stack
    wait_for_healthy
    print_next_steps
    
    log_success "Installation complete!"
}

# Handle interrupts
trap 'log_error "Installation interrupted"; exit 130' INT TERM

# Run main function
main "$@"
