#!/bin/bash
# =============================================================================
# HelixKnowledge Skill Graph System - Start Script
# =============================================================================
# Usage: ./start.sh [--compose | --systemd | --foreground]
# Defaults to systemd if available, falls back to docker compose.
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_DIR="$(dirname "$SCRIPT_DIR")"
SERVICE_NAME="skill-system"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Detect container runtime and compose
detect_runtime() {
    if command -v docker &> /dev/null; then
        CONTAINER_RUNTIME="docker"
        if docker compose version &> /dev/null; then
            COMPOSE_CMD="docker compose"
        elif command -v docker-compose &> /dev/null; then
            COMPOSE_CMD="docker-compose"
        fi
    elif command -v podman &> /dev/null; then
        CONTAINER_RUNTIME="podman"
        if command -v podman-compose &> /dev/null; then
            COMPOSE_CMD="podman-compose"
        fi
    else
        log_error "No container runtime found"
        exit 1
    fi
}

# Check if systemd user service is available
has_systemd() {
    systemctl --user status &> /dev/null && \
        [ -f "$HOME/.config/systemd/user/${SERVICE_NAME}.service" ]
}

# Start via systemd
start_systemd() {
    log_info "Starting via systemd user service..."
    systemctl --user start "$SERVICE_NAME"
    sleep 3
    
    if systemctl --user is-active "$SERVICE_NAME" &> /dev/null; then
        log_success "Service started successfully"
        systemctl --user status "$SERVICE_NAME" --no-pager
    else
        log_error "Failed to start service"
        journalctl --user -u "$SERVICE_NAME" --no-pager -n 20
        exit 1
    fi
}

# Start via docker compose
start_compose() {
    log_info "Starting via $COMPOSE_CMD..."
    cd "$INSTALL_DIR"
    
    # Pull latest images
    log_info "Pulling images..."
    $COMPOSE_CMD pull
    
    # Start services
    $COMPOSE_CMD up -d --remove-orphans
    
    log_success "Stack started"
    $COMPOSE_CMD ps
}

# Start in foreground (for debugging)
start_foreground() {
    log_info "Starting in foreground mode..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD up --remove-orphans
}

# Wait for health check
wait_for_health() {
    log_info "Waiting for API health check..."
    local retries=30
    
    while [ $retries -gt 0 ]; do
        if curl -sf http://localhost:8080/health &> /dev/null; then
            log_success "API is healthy"
            return 0
        fi
        retries=$((retries - 1))
        echo -n "."
        sleep 2
    done
    
    log_warn "API not responding yet. Check logs with: ./scripts/status.sh"
    return 1
}

# Main
main() {
    detect_runtime
    
    local mode=""
    
    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --systemd)
                mode="systemd"
                ;;
            --compose)
                mode="compose"
                ;;
            --foreground)
                mode="foreground"
                ;;
            --help|-h)
                echo "Usage: $0 [--compose | --systemd | --foreground]"
                echo ""
                echo "Options:"
                echo "  --compose     Start via docker/podman compose (default fallback)"
                echo "  --systemd     Start via systemd user service"
                echo "  --foreground  Start in foreground for debugging"
                echo "  --help        Show this help"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
        shift
    done
    
    # Auto-detect mode
    if [ -z "$mode" ]; then
        if has_systemd; then
            mode="systemd"
        else
            mode="compose"
        fi
    fi
    
    log_info "Starting Skill Graph System (mode: $mode)..."
    
    case "$mode" in
        systemd)
            start_systemd
            ;;
        compose)
            start_compose
            wait_for_health
            ;;
        foreground)
            start_foreground
            ;;
    esac
    
    echo ""
    log_success "Skill Graph System is running!"
    echo "  API:     http://localhost:8080"
    echo "  Health:  curl http://localhost:8080/health"
}

main "$@"
