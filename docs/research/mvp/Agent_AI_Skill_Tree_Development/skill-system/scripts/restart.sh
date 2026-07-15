#!/bin/bash
# =============================================================================
# HelixKnowledge Skill Graph System - Restart Script
# =============================================================================
# Usage: ./restart.sh [--quick | --full]
#   --quick  Rolling restart (default)
#   --full   Full stop, rebuild, and start
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

# Quick rolling restart
restart_quick() {
    log_info "Performing quick rolling restart..."
    cd "$INSTALL_DIR"
    
    # Restart worker first (no external ports)
    log_info "Restarting worker..."
    $COMPOSE_CMD restart worker
    sleep 2
    
    # Restart API
    log_info "Restarting API..."
    $COMPOSE_CMD restart api
    sleep 2
    
    # Verify API is up
    local retries=30
    while [ $retries -gt 0 ]; do
        if curl -sf http://localhost:8080/health &> /dev/null; then
            log_success "API is healthy after restart"
            return 0
        fi
        retries=$((retries - 1))
        echo -n "."
        sleep 2
    done
    
    log_warn "API not responding after restart. Check logs."
    return 1
}

# Full restart with rebuild
restart_full() {
    log_info "Performing full restart with rebuild..."
    
    # Stop
    "$SCRIPT_DIR/stop.sh"
    
    # Rebuild images
    log_info "Rebuilding images..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD build --no-cache
    
    # Start
    "$SCRIPT_DIR/start.sh"
    
    log_success "Full restart complete"
}

# Main
main() {
    detect_compose
    
    local mode="quick"
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --quick)
                mode="quick"
                ;;
            --full)
                mode="full"
                ;;
            --help|-h)
                echo "Usage: $0 [--quick | --full]"
                echo ""
                echo "Options:"
                echo "  --quick  Rolling restart (default)"
                echo "  --full   Full stop, rebuild, and start"
                echo "  --help   Show this help"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
        shift
    done
    
    log_info "Restarting Skill Graph System (mode: $mode)..."
    
    case "$mode" in
        quick)
            restart_quick
            ;;
        full)
            restart_full
            ;;
    esac
    
    echo ""
    log_success "Restart complete!"
    echo "  API: http://localhost:8080"
    echo "  Run ./scripts/status.sh for details"
}

main "$@"
