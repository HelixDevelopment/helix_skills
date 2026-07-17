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

## Usage

### Command-line invocation

```bash
# Validate a screenshot against expected patterns
bash constitution/skills/media-validator/media-validator.sh \
  /path/to/screenshot.png "Settings" "Connected" "API 35"

# Validate a video recording
bash constitution/skills/media-validator/media-validator.sh \
  /path/to/recording.mp4 "PASS" "Playback started"

# Validate text output
bash constitution/skills/media-validator/media-validator.sh \
  /path/to/output.txt "Build successful" "0 failures"
```

### Exit codes

| Code | Meaning |
|---|---|
| 0 | PASS — all expected patterns matched |
| 1 | FAIL — one or more patterns unmatched; evidence written to `qa-results/media-validator/` |
| 2 | SKIP — file type not supported or missing dependencies |

### MCP server usage

The media-validator is also available as an MCP server (`constitution/mcp/media-validator-mcp.json`). When configured, it exposes validation as a tool call rather than a shell invocation.

## Constitution References

| Reference | Meaning |
|---|---|
| **§11.4.163** | Universal Media Validation & Verification Mandate. Every recorded artifact (MP4, PNG, TXT, asciinema) MUST pass through a self-validated media validation pipeline before being accepted as evidence. Defines the seven requirements: content extraction, pattern matching, self-validated analyzer, structured verdict, exact pinpoint on FAIL, dual trigger mode, and paired §1.1 meta-test. |
| **§11.4.107** | Anti-bluff AV/test-validation techniques. The pipeline must be mutation-tested with golden-good/golden-bad fixture pairs. A pipeline that passes its golden-bad fixture is a bluff gate. |
| **§11.4.117** | Computer-vision / OCR pixel-oracle fallback. OCR extraction with a confidence floor for non-introspectable UIs. |
| **§11.4.137** | Subtitle/caption content-correctness oracle. Validates subtitle rendering against expected text. |
| **§11.4.69** | Universal sink-side positive-evidence taxonomy. The structured verdict is a taxonomy-class observable. |
| **§11.4.5** | Captured-evidence quality. The verdict's evidence path and pinpoint data meet the captured-evidence standard. |
| **§11.4.102** | Systematic-debugging. On FAIL, the pinpoint data is the root-cause entry point for investigation before re-recording. |
| **§11.4.83** | docs/qa/ end-user evidence mandate. Verdicts are committed under `docs/qa/<run-id>/`. |

## Cross-links

- **Parent domain:** [`validation`](../by-domain/validation.md)
- **Related skills:** `helixqa` (HelixQA bank consumes media-validator verdicts), `systematic-debugging` (FAIL pinpoint triggers §11.4.102 investigation)
- **Constitution source:** [`constitution/skills/media-validator/`](../../../constitution/skills/media-validator/)
- **MCP config:** [`constitution/mcp/media-validator-mcp.json`](../../../constitution/mcp/media-validator-mcp.json)

## Integration

| Surface | How it hooks in |
|---|---|
| **CLI script** | `media-validator.sh` — standalone bash script. Extracts content (OCR via tesseract, frame extraction via ffmpeg, text parsing via `file` + grep), matches against expected patterns, writes structured verdict. |
| **MCP server** | `media-validator-mcp.json` — exposes validation as an MCP tool. Agents call it as a tool rather than shelling out. |
| **Dual trigger mode** | Post-recording (validate after full recording completes) and real-time monitoring (analyze frames/output during recording per §11.4.159(M)). |
| **Evidence pipeline** | Verdicts committed alongside artifacts under `docs/qa/<run-id>/`. FAIL verdicts include exact pinpoint data (file + frame/line/timestamp + OCR region + expected vs actual). |
| **Self-validation** | Golden-good fixture MUST produce PASS, golden-bad fixture MUST produce FAIL. The pipeline is inadmissible until its self-validation is restored. |
