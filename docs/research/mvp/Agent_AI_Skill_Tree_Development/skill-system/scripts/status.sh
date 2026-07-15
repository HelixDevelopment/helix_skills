#!/bin/bash
# =============================================================================
# HelixKnowledge Skill Graph System - Status Script
# =============================================================================
# Usage: ./status.sh [--watch | --json]
#   --watch  Continuous monitoring mode (refresh every 5s)
#   --json   Output in JSON format
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
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Detect compose command
detect_compose() {
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
        CONTAINER_RUNTIME="docker"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
        CONTAINER_RUNTIME="docker"
    elif command -v podman-compose &> /dev/null; then
        COMPOSE_CMD="podman-compose"
        CONTAINER_RUNTIME="podman"
    else
        echo "No compose command found"
        exit 1
    fi
}

# Show container status
show_containers() {
    echo -e "\n${BOLD}Container Status${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    
    cd "$INSTALL_DIR"
    
    local containers=$("$COMPOSE_CMD" ps --format json 2>/dev/null || "$COMPOSE_CMD" ps 2>/dev/null)
    
    if [ -z "$containers" ]; then
        log_warn "No containers running"
        return
    fi
    
    # Try JSON format first
    if echo "$containers" | head -1 | grep -q '{'; then
        echo "$containers" | while read -r line; do
            local name=$(echo "$line" | jq -r '.Name // .name // "unknown"' 2>/dev/null || echo "unknown")
            local state=$(echo "$line" | jq -r '.State // .state // "unknown"' 2>/dev/null || echo "unknown")
            local health=$(echo "$line" | jq -r '.Health // .health // "N/A"' 2>/dev/null || echo "N/A")
            local status=$(echo "$line" | jq -r '.Status // .status // ""' 2>/dev/null || echo "")
            
            local state_color="$RED"
            [ "$state" = "running" ] && state_color="$GREEN"
            
            printf "  %-20s ${state_color}%-10s${NC} Health: %-10s %s\n" \
                "$name" "$state" "$health" "$status"
        done
    else
        # Fallback to table format
        "$COMPOSE_CMD" ps
    fi
}

# Show systemd status
show_systemd() {
    if [ -f "$HOME/.config/systemd/user/${SERVICE_NAME}.service" ]; then
        echo -e "\n${BOLD}systemd Service${NC}"
        echo "═══════════════════════════════════════════════════════════════"
        
        local status=$(systemctl --user is-active "$SERVICE_NAME" 2>/dev/null || echo "inactive")
        local enabled=$(systemctl --user is-enabled "$SERVICE_NAME" 2>/dev/null || echo "disabled")
        
        local status_color="$RED"
        [ "$status" = "active" ] && status_color="$GREEN"
        
        printf "  Service:  ${status_color}%s${NC}\n" "$status"
        printf "  Enabled:  %s\n" "$enabled"
        
        if [ "$status" = "active" ]; then
            local uptime=$(systemctl --user show "$SERVICE_NAME" --property=ActiveEnterTimestamp --value 2>/dev/null)
            printf "  Uptime:   %s\n" "$uptime"
        fi
    fi
}

# Check API health
check_api_health() {
    echo -e "\n${BOLD}API Health${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    
    if command -v curl &> /dev/null; then
        local response
        response=$(curl -sf http://localhost:8080/health 2>/dev/null || echo "")
        
        if [ -n "$response" ]; then
            log_success "API responding"
            
            if command -v jq &> /dev/null; then
                local status=$(echo "$response" | jq -r '.status // "unknown"' 2>/dev/null)
                local version=$(echo "$response" | jq -r '.version // "unknown"' 2>/dev/null)
                local db=$(echo "$response" | jq -r '.database // "unknown"' 2>/dev/null)
                
                printf "  Status:   %s\n" "$status"
                printf "  Version:  %s\n" "$version"
                printf "  Database: %s\n" "$db"
            else
                echo "  Response: $response"
            fi
        else
            log_error "API not responding at http://localhost:8080/health"
        fi
    else
        log_warn "curl not available, skipping API health check"
    fi
}

# Check database
check_database() {
    echo -e "\n${BOLD}Database${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    
    cd "$INSTALL_DIR"
    
    if "$COMPOSE_CMD" exec -T db pg_isready -U skilluser -d skilldb &> /dev/null; then
        log_success "PostgreSQL is ready"
        
        # Get database stats
        local stats
        stats=$("$COMPOSE_CMD" exec -T db psql -U skilluser -d skilldb -t -c "
            SELECT 
                (SELECT COUNT(*) FROM skills) as skills,
                (SELECT COUNT(*) FROM skill_relationships) as relationships,
                (SELECT COUNT(*) FROM evidence) as evidence,
                pg_size_pretty(pg_database_size('skilldb')) as db_size;
        " 2>/dev/null || echo "")
        
        if [ -n "$stats" ]; then
            echo "  Stats:"
            echo "$stats" | while read -r line; do
                echo "    $line"
            done
        fi
    else
        log_error "PostgreSQL not ready"
    fi
}

# Show recent logs
show_logs() {
    echo -e "\n${BOLD}Recent Logs (last 20 lines)${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    
    cd "$INSTALL_DIR"
    "$COMPOSE_CMD" logs --tail=20 --no-log-prefix 2>/dev/null || true
}

# Show resource usage
show_resources() {
    echo -e "\n${BOLD}Resource Usage${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    
    cd "$INSTALL_DIR"
    local containers
    containers=$("$COMPOSE_CMD" ps -q 2>/dev/null || true)
    
    if [ -n "$containers" ]; then
        for container in $containers; do
            local name=$($CONTAINER_RUNTIME inspect --format='{{.Name}}' "$container" 2>/dev/null | sed 's/\///')
            local mem_usage=$($CONTAINER_RUNTIME stats --no-stream --format='{{.MemUsage}}' "$container" 2>/dev/null || echo "N/A")
            local cpu_usage=$($CONTAINER_RUNTIME stats --no-stream --format='{{.CPUPerc}}' "$container" 2>/dev/null || echo "N/A")
            
            printf "  %-20s CPU: %-8s MEM: %s\n" "$name" "$cpu_usage" "$mem_usage"
        done
    fi
}

# JSON output
json_output() {
    local api_response=""
    local api_status="unhealthy"
    
    if command -v curl &> /dev/null; then
        api_response=$(curl -sf http://localhost:8080/health 2>/dev/null || echo "")
        [ -n "$api_response" ] && api_status="healthy"
    fi
    
    cd "$INSTALL_DIR"
    local container_json
    container_json=$("$COMPOSE_CMD" ps --format json 2>/dev/null || echo "[]")
    
    cat << EOF
{
    "service": "$SERVICE_NAME",
    "timestamp": "$(date -Iseconds)",
    "api": {
        "status": "$api_status",
        "response": $(echo "$api_response" | jq . 2>/dev/null || echo "null")
    },
    "containers": $container_json
}
EOF
}

# Watch mode
watch_mode() {
    while true; do
        clear
        echo -e "${CYAN}$(date '+%Y-%m-%d %H:%M:%S') - Skill Graph System Monitor${NC}"
        echo "Press Ctrl+C to exit"
        
        show_systemd
        show_containers
        check_api_health
        check_database
        show_resources
        
        sleep 5
    done
}

# Main
main() {
    detect_compose
    
    local mode=""
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --watch)
                mode="watch"
                ;;
            --json)
                mode="json"
                ;;
            --help|-h)
                echo "Usage: $0 [--watch | --json]"
                echo ""
                echo "Options:"
                echo "  --watch  Continuous monitoring mode"
                echo "  --json   JSON output"
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
    
    case "$mode" in
        watch)
            watch_mode
            ;;
        json)
            json_output
            ;;
        *)
            # Default status display
            echo -e "${BOLD}"
            echo "╔═══════════════════════════════════════════════════════════════╗"
            echo "║         Skill Graph System Status                             ║"
            echo "╚═══════════════════════════════════════════════════════════════╝"
            echo -e "${NC}"
            
            show_systemd
            show_containers
            check_api_health
            check_database
            show_resources
            
            echo -e "\n${CYAN}Run with --watch for continuous monitoring${NC}"
            ;;
    esac
}

main "$@"
