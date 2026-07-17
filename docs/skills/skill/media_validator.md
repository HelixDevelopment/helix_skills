# media-validator

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** media-validator
- **Title:** Universal Media Validation
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** validation
- **Complexity:** intermediate
- **Tags:** media, validation, ocr, metadata, pattern-matching

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Validate media files (MP4, PNG, TXT) via OCR/metadata/pattern matching for
PASS/FAIL with evidence. Implements the §11.4.163 Universal Media Validation
mandate. Available both as a constitution skill and as an MCP server
(`media-validator-mcp.json`).

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

_None._

### Optional

_None._

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/media-validator/` | Directory | Skill source |
| `constitution/mcp/media-validator-mcp.json` | MCP config | MCP server definition |
| `constitution/skills/media-validator/media-validator.sh` | Script | Validation script |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Source:** `constitution/skills/media-validator/`
- **Consumed via:** Constitution skill + MCP server
- **Constitution references:** §11.4.163
- **Created:** 2026-07-15
- **Last updated:** 2026-07-17

<!-- /skills-catalog:section=metadata -->
