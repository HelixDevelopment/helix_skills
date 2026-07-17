# GandalfSolutions

- **GitHub URL**: <https://github.com/vasic-digital/GandalfSolutions>
- **Description**: Read-only solutions archive for prompt-leak-defense testing
- **Category**: Auth + Security + Middleware
- **Status**: SCAFFOLD / WIP

## Overview

GandalfSolutions maintains a read-only archive of known solutions and defensive patterns for prompt-leak testing. It serves as a reference corpus for validating that LLM applications resist known prompt-extraction techniques. The archive is intentionally read-only to preserve the integrity of test fixtures.

## Tech Stack

- Language: Multiple
- Framework: Static archive / test fixture repository

## Key Features

- Read-only archive of prompt-leak defense solutions
- Curated test fixtures for validating LLM application security
- Reference patterns for known extraction techniques
- Versioned archive for regression testing

## Related Repos

- [LeakHub](../LeakHub/README.md) — prompt-leak corpus that GandalfSolutions draws defensive patterns from
- [Claritas](../Claritas/README.md) — detection engine that validates against GandalfSolutions fixtures
- [Ouroborous](../Ouroborous/README.md) — recursive safety patterns that complement static defense archives

---
*Part of the [vasic-digital catalogue](../README.md)*
