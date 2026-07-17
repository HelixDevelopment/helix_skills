# media-validator

> **Path:** `constitution/skills/media-validator/` · `constitution/mcp/media-validator-mcp.json`
> **Type:** Skill + MCP Server · **Status:** Active

## What it provides

Validate media files (MP4, PNG, TXT) via OCR/metadata/pattern matching for
PASS/FAIL with evidence. Implements §11.4.163 Universal Media Validation.

## How consumed

- **As skill:** installed via `register.sh` into the agent's skill set
- **As MCP server:** defined in `media-validator-mcp.json`, executed as bash script

## Source paths

- Skill: `constitution/skills/media-validator/`
- MCP config: `constitution/mcp/media-validator-mcp.json`
- Script: `constitution/skills/media-validator/media-validator.sh`

## Dependencies

None.

## Constitution references

§11.4.163
