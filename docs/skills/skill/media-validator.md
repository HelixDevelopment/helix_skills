# media-validator

> **GENERATED FILE — DO NOT HAND-EDIT.** Regenerated from the live skill
> graph by the `skills-catalog` generator. Edit the skill via CLI/REST/MCP
> (see `docs/scripts/` / `docs/API.md`) — this file will be overwritten.

<!-- skills-catalog:section=header -->
## Header

- **Name:** media-validator
- **Title:** Universal Media Validation and Verification
- **Version:** 0.1.0
- **Kind:** atomic
- **Status:** active
- **Domain:** validation
- **Complexity:** intermediate
- **Tags:** media, validation, OCR, vision, liveness, evidence

<!-- /skills-catalog:section=header -->
<!-- skills-catalog:section=description -->
## Description

Universal media validation & verification (§11.4.163, §11.4.107). Validates
media files (MP4, PNG, TXT) via OCR, metadata extraction, and pattern
matching for PASS/FAIL verdicts with captured evidence. Implements liveness
detection, freeze-detection oracles, and self-validated analyzer fixtures
(golden-good/golden-bad) per §11.4.107(10). Also exposed as an MCP server.

<!-- /skills-catalog:section=description -->
<!-- skills-catalog:section=dependencies -->
## Dependencies

### Requires

_(none — standalone validation tool)_

### Optional

- `continuum` — for persisting validation evidence across sessions

<!-- /skills-catalog:section=dependencies -->
<!-- skills-catalog:section=resources -->
## Resources

| Resource | Type | Description |
|---|---|---|
| `constitution/skills/media-validator/` | Directory | Skill source |
| `constitution/mcp/media-validator-mcp.json` | File | MCP server definition |
| `constitution/skills/media-validator/media-validator.sh` | Script | Validation script |

<!-- /skills-catalog:section=resources -->
<!-- skills-catalog:section=metadata -->
## Metadata

- **Created:** 2026-06-21
- **Updated:** 2026-07-16
- **Author:** HelixDevelopment
- **Source:** `constitution/skills/media-validator/`

<!-- /skills-catalog:section=metadata -->
