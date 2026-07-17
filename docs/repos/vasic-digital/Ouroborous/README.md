# Ouroborous

- **GitHub URL**: <https://github.com/vasic-digital/Ouroborous>
- **Description**: Recursive/self-referential safety patterns
- **Category**: Auth + Security + Middleware
- **Status**: SCAFFOLD / WIP

## Overview

Ouroborous explores and implements recursive and self-referential safety patterns for AI systems. It addresses edge cases where AI agents interact with their own outputs, create feedback loops, or attempt self-modification. The project provides guardrails and detection mechanisms for these unusual but critical safety scenarios.

## Tech Stack

- Language: Multiple
- Framework: Custom safety pattern engine

## Key Features

- Detection of recursive self-referential loops in AI agent behavior
- Safety guardrails for self-modifying or self-referencing agent outputs
- Pattern library for known recursive safety failure modes
- Integration hooks for middleware-level safety enforcement

## Related Repos

- [Claritas](../Claritas/README.md) — prompt extraction detection that complements recursive safety checks
- [Veritas](../Veritas/README.md) — truth verification that can flag self-referential inconsistencies

---
*Part of the [vasic-digital catalogue](../README.md)*
