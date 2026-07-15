#!/bin/bash
# =============================================================================
# HelixKnowledge Skill Graph System - Stop Script
# =============================================================================
# Usage: ./stop.sh [--systemd | --compose | --kill]
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

# Detect compose command
detect_compose() {
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif command -v podman-compose &> /dev/null; then
        COMPOSE_CMD="podman-compose"
    fi
}

# Stop via systemd
stop_systemd() {
    log_info "Stopping via systemd..."
    systemctl --user stop "$SERVICE_NAME"
    log_success "Service stopped"
}

# Stop via compose
stop_compose() {
    log_info "Stopping via compose..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD down --timeout 30
    log_success "Stack stopped"
}

# Kill all containers (emergency)
stop_kill() {
    log_warn "Forcefully stopping all containers..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD kill
    $COMPOSE_CMD down
    log_success "Containers killed"
}

# Main
main() {
    detect_compose
    
    local mode=""
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --systemd)
                mode="systemd"
                ;;
            --compose)
                mode="compose"
                ;;
            --kill)
                mode="kill"
                ;;
            --help|-h)
                echo "Usage: $0 [--compose | --systemd | --kill]"
                echo ""
                echo "Options:"
                echo "  --compose  Stop via docker/podman compose"
                echo "  --systemd  Stop via systemd user service"
                echo "  --kill     Force kill all containers"
                echo "  --help     Show this help"
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
        if systemctl --user is-active "$SERVICE_NAME" &> /dev/null; then
            mode="systemd"
        else
            mode="compose"
        fi
    fi
    
    log_info "Stopping Skill Graph System (mode: $mode)..."
    
    case "$mode" in
        systemd)
            stop_systemd
            ;;
        compose)
            stop_compose
            ;;
        kill)
            stop_kill
            ;;
    esac
    
    log_success "Skill Graph System stopped"
}

main "$@"
