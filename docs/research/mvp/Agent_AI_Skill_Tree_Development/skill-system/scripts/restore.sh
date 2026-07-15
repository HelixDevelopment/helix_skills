#!/bin/bash
# =============================================================================
# HelixKnowledge Skill Graph System - Restore Script
# =============================================================================
# Usage: ./restore.sh <backup-file> [--force] [--db-only]
# Restores from a backup archive created by backup.sh
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

# Defaults
FORCE=false
DB_ONLY=false

# Load environment
load_env() {
    if [ -f "$INSTALL_DIR/.env" ]; then
        export $(grep -v '^#' "$INSTALL_DIR/.env" | xargs 2>/dev/null || true)
    fi
    
    DB_HOST="${DB_HOST:-db}"
    DB_PORT="${DB_PORT:-5432}"
    DB_NAME="${DB_NAME:-skilldb}"
    DB_USER="${DB_USER:-skilluser}"
    DB_PASSWORD="${DB_PASSWORD:-skillpassword}"
}

# Detect compose command
detect_compose() {
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif command -v podman-compose &> /dev/null; then
        COMPOSE_CMD="podman-compose"
    else
        log_error "No compose command found"
        exit 1
    fi
}

# Show backup info
show_backup_info() {
    local archive="$1"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    tar -xzf "$archive" -C "$temp_dir" 2>/dev/null || {
        log_error "Invalid backup archive"
        rm -rf "$temp_dir"
        exit 1
    }
    
    local metadata_file
    metadata_file=$(find "$temp_dir" -name "backup.json" 2>/dev/null | head -1)
    
    if [ -n "$metadata_file" ] && [ -f "$metadata_file" ]; then
        echo -e "\n${BLUE}Backup Information${NC}"
        echo "═══════════════════════════════════════════════════════════════"
        
        if command -v jq &> /dev/null; then
            echo "  Name:        $(jq -r '.name' "$metadata_file")"
            echo "  Type:        $(jq -r '.type' "$metadata_file")"
            echo "  Created:     $(jq -r '.created_at' "$metadata_file")"
            echo "  Version:     $(jq -r '.version' "$metadata_file")"
            echo "  Hostname:    $(jq -r '.hostname' "$metadata_file")"
            echo ""
            echo "  Contents:"
            echo "    Database:      $(jq -r '.contents.database' "$metadata_file")"
            echo "    Config:        $(jq -r '.contents.configuration' "$metadata_file")"
            echo "    Evidence:      $(jq -r '.contents.evidence' "$metadata_file")"
        else
            cat "$metadata_file"
        fi
    else
        log_warn "No metadata found in backup"
    fi
    
    rm -rf "$temp_dir"
}

# Restore database
restore_database() {
    local archive="$1"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    log_info "Extracting database backup..."
    tar -xzf "$archive" -C "$temp_dir" 2>/dev/null
    
    local db_dump
    db_dump=$(find "$temp_dir" -name "skilldb.sql.gz" -o -name "skilldb.sql" 2>/dev/null | head -1)
    
    if [ -z "$db_dump" ]; then
        log_error "Database dump not found in backup"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    log_info "Database dump: $db_dump"
    
    # Decompress if needed
    local sql_file="$db_dump"
    if [[ "$db_dump" == *.gz ]]; then
        sql_file="${db_dump%.gz}"
        gunzip -c "$db_dump" > "$sql_file"
    fi
    
    # Check if database exists and prompt
    local db_exists
    db_exists=$(cd "$INSTALL_DIR" && $COMPOSE_CMD exec -T db psql -U "$DB_USER" -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME';" 2>/dev/null | tr -d '[:space:]')
    
    if [ "$db_exists" = "1" ] && [ "$FORCE" = false ]; then
        echo ""
        log_warn "Database '$DB_NAME' already exists"
        read -p "Drop and recreate? [y/N]: " confirm
        if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
            log_info "Restore cancelled"
            rm -rf "$temp_dir"
            exit 0
        fi
    fi
    
    # Drop and recreate database
    log_info "Recreating database..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD exec -T db psql -U "$DB_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || true
    $COMPOSE_CMD exec -T db psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;" 2>/dev/null
    
    # Restore from dump
    log_info "Restoring database (this may take a while)..."
    if ! $COMPOSE_CMD exec -T db psql -U "$DB_USER" -d "$DB_NAME" < "$sql_file" 2>/dev/null; then
        log_error "Database restore failed"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    log_success "Database restored"
    rm -rf "$temp_dir"
}

# Restore configuration
restore_config() {
    local archive="$1"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    log_info "Restoring configuration..."
    tar -xzf "$archive" -C "$temp_dir" 2>/dev/null
    
    local config_dir
    config_dir=$(find "$temp_dir" -type d -name "config" 2>/dev/null | head -1)
    
    if [ -n "$config_dir" ] && [ -d "$config_dir" ]; then
        # Backup current config
        if [ -f "$INSTALL_DIR/.env" ]; then
            cp "$INSTALL_DIR/.env" "$INSTALL_DIR/.env.restore-backup.$(date +%s)"
        fi
        
        # Restore files
        [ -f "$config_dir/.env" ] && cp "$config_dir/.env" "$INSTALL_DIR/"
        [ -f "$config_dir/config.toml" ] && cp "$config_dir/config.toml" "$INSTALL_DIR/config/"
        
        log_success "Configuration restored"
    else
        log_warn "No configuration found in backup"
    fi
    
    rm -rf "$temp_dir"
}

# Restore evidence
restore_evidence() {
    local archive="$1"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    log_info "Restoring evidence data..."
    tar -xzf "$archive" -C "$temp_dir" 2>/dev/null
    
    local evidence_archive
    evidence_archive=$(find "$temp_dir" -name "evidence.tar.gz" 2>/dev/null | head -1)
    
    if [ -n "$evidence_archive" ] && [ -f "$evidence_archive" ]; then
        mkdir -p "$INSTALL_DIR/data"
        tar -xzf "$evidence_archive" -C "$INSTALL_DIR/data" 2>/dev/null || true
        log_success "Evidence data restored"
    else
        log_warn "No evidence data found in backup"
    fi
    
    rm -rf "$temp_dir"
}

# Main
main() {
    load_env
    detect_compose
    
    local backup_file=""
    
    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --force)
                FORCE=true
                shift
                ;;
            --db-only)
                DB_ONLY=true
                shift
                ;;
            --help|-h)
                echo "Usage: $0 <backup-file> [options]"
                echo ""
                echo "Arguments:"
                echo "  <backup-file>  Path to backup archive (.tar.gz)"
                echo ""
                echo "Options:"
                echo "  --force     Skip confirmation prompts"
                echo "  --db-only   Restore only database"
                echo "  --help      Show this help"
                echo ""
                echo "Examples:"
                echo "  $0 data/backups/skill-system_v1.0_20240115.tar.gz"
                echo "  $0 backup.tar.gz --force"
                echo "  $0 backup.tar.gz --db-only"
                exit 0
                ;;
            -*)
                log_error "Unknown option: $1"
                exit 1
                ;;
            *)
                if [ -z "$backup_file" ]; then
                    backup_file="$1"
                fi
                shift
                ;;
        esac
    done
    
    if [ -z "$backup_file" ]; then
        log_error "Backup file required"
        echo "Run '$0 --help' for usage"
        exit 1
    fi
    
    if [ ! -f "$backup_file" ]; then
        log_error "Backup file not found: $backup_file"
        exit 1
    fi
    
    # Show backup info
    show_backup_info "$backup_file"
    
    # Confirm restore
    if [ "$FORCE" = false ]; then
        echo ""
        log_warn "This will overwrite existing data!"
        read -p "Continue with restore? [y/N]: " confirm
        if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
            log_info "Restore cancelled"
            exit 0
        fi
    fi
    
    # Stop services before restore
    log_info "Stopping services..."
    "$SCRIPT_DIR/stop.sh" --compose 2>/dev/null || true
    
    # Start just the database
    log_info "Starting database..."
    cd "$INSTALL_DIR"
    $COMPOSE_CMD up -d db
    
    # Wait for database
    log_info "Waiting for database..."
    sleep 5
    
    local retries=30
    while [ $retries -gt 0 ]; do
        if $COMPOSE_CMD exec -T db pg_isready -U "$DB_USER" -d "$DB_NAME" &> /dev/null; then
            break
        fi
        retries=$((retries - 1))
        sleep 2
    done
    
    # Restore database
    restore_database "$backup_file"
    
    # Restore configuration and evidence (unless db-only)
    if [ "$DB_ONLY" = false ]; then
        restore_config "$backup_file"
        restore_evidence "$backup_file"
    fi
    
    # Restart services
    log_info "Restarting services..."
    "$SCRIPT_DIR/start.sh"
    
    echo ""
    log_success "Restore complete!"
    log_info "Verify with: curl http://localhost:8080/health"
}

main "$@"
