# scheduled-work

> **Path:** `constitution/mcp/scheduled-work-mcp.json` · `constitution/plugins/scheduled-work/`
> **Type:** MCP Server + Plugin · **Status:** Active

## What it provides

Track scheduled work/reminders (background-queue entries): create, list,
status, overdue, needs-verification, mark-done. Decoupled Go MCP server
backing the REMINDER/BACKGROUND action-prefix (§11.4.140).

## How consumed

- **As MCP server:** compiled Go binary at `scripts/scheduled-work-engine/bin/scheduled-work`
- **As plugin:** wired via `constitution/plugins/scheduled-work/.mcp.json`

## Source paths

- MCP config: `constitution/mcp/scheduled-work-mcp.json`
- Plugin: `constitution/plugins/scheduled-work/`
- Go module: `constitution/scripts/scheduled-work-engine/`
- Binary: `scripts/scheduled-work-engine/bin/scheduled-work`

## Dependencies

gin-gonic, quic-go, brotli, uuid, yaml (Go dependencies)

## Constitution references

§11.4.140, §11.4.6, §11.4.108
