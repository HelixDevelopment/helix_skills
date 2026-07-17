# Claritas

- **GitHub URL**: <https://github.com/vasic-digital/Claritas>
- **Description**: System-prompt extraction detection
- **Category**: Auth + Security + Middleware
- **Status**: SCAFFOLD / WIP

## Overview

Claritas detects attempts to extract system prompts from LLM-based applications. It monitors input patterns and conversation flows to identify prompt-leak attack vectors, providing early warning when adversaries attempt to reverse-engineer protected instructions. The tool is designed for integration into LLM middleware pipelines.

## Tech Stack

- Language: Multiple
- Framework: Custom detection engine

## Key Features

- Detection of system-prompt extraction attempts in LLM conversations
- Pattern analysis for known prompt-leak attack vectors
- Configurable sensitivity and alerting thresholds
- Middleware integration for real-time monitoring

## Related Repos

- [GandalfSolutions](../GandalfSolutions/README.md) — solutions archive for testing prompt-leak defenses that Claritas detects
- [LeakHub](../LeakHub/README.md) — prompt-leak corpus providing training data for detection models
- [Veritas](../Veritas/README.md) — truth verification auxiliary that complements extraction detection

---
*Part of the [vasic-digital catalogue](../README.md)*
