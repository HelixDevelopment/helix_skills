# Normalize

- **GitHub URL**: <https://github.com/vasic-digital/Normalize>
- **Description**: Adversarial-input canonicalisation library for defensive LLM guardrail pipelines
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

Normalize provides canonicalisation functions that transform adversarial or obfuscated LLM inputs into a standard form before they reach guardrail classifiers. It defends against prompt-injection attacks that use Unicode confusables, homoglyph substitution, invisible characters, whitespace manipulation, and other encoding tricks to bypass content filters. Designed as a preprocessing stage in any defensive LLM pipeline.

## Tech Stack

- Language: Go
- Module: digital.vasic.normalize

## Key Features

- Unicode confusable and homoglyph normalisation to defeat visual spoofing attacks
- Invisible character and zero-width token stripping
- Whitespace and encoding canonicalisation
- Composable normalisation pipeline with pluggable stages
- Designed for low-latency inline preprocessing in LLM request paths

## Related Repos

- [Claritas](../Claritas/README.md) — system-prompt extraction detection that consumes normalised inputs
- [LeakHub](../LeakHub/README.md) — prompt-leak corpus providing adversarial fixtures for testing normalisation coverage
- [RedTeam](../RedTeam/README.md) — adversarial prompt harness that exercises normalisation edge cases

---
*Part of the [vasic-digital catalogue](../README.md)*
