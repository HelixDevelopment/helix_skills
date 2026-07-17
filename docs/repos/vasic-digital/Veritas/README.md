# Veritas

- **GitHub URL**: <https://github.com/vasic-digital/Veritas>
- **Description**: Truth/verification auxiliary
- **Category**: Auth + Security + Middleware
- **Status**: SCAFFOLD / WIP

## Overview

Veritas provides truth verification and fact-checking capabilities as an auxiliary service for AI systems. It cross-references LLM outputs against trusted sources and flags potential hallucinations or factual inaccuracies. The module is designed to sit alongside LLM pipelines as a post-generation verification step.

## Tech Stack

- Language: Multiple
- Framework: Custom verification engine

## Key Features

- Post-generation truth verification for LLM outputs
- Cross-referencing against trusted knowledge sources
- Hallucination and factual inaccuracy flagging
- Auxiliary service integration for pipeline-level verification

## Related Repos

- [Claritas](../Claritas/README.md) — prompt extraction detection that shares the security monitoring surface
- [Ouroborous](../Ouroborous/README.md) — recursive safety patterns that Veritas can flag as inconsistencies

---
*Part of the [vasic-digital catalogue](../README.md)*
