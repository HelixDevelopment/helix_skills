# AutoTemp

- **GitHub URL**: <https://github.com/vasic-digital/AutoTemp>
- **Description**: Benchmark-driven temperature auto-tuning orchestration
- **Category**: Container + Lifecycle
- **Status**: SCAFFOLD / WIP

## Overview

AutoTemp automates the process of finding optimal temperature settings for LLM-based workflows. It runs benchmark suites across a range of temperature values, measures output quality metrics, and selects the best-performing configuration. This removes manual guesswork from tuning creative vs. deterministic output balance.

## Tech Stack

- Language: Multiple
- Framework: Custom benchmark orchestration

## Key Features

- Automated temperature sweep across configurable ranges
- Benchmark-driven evaluation of LLM output quality
- Orchestrated parallel runs for efficient tuning
- Result aggregation and best-parameter selection

## Related Repos

- [HyperTune](../HyperTune/README.md) — broader hyperparameter tuning that complements temperature-specific tuning
- [LLMOrchestrator](../LLMOrchestrator/README.md) — manages the LLM agents whose temperature settings AutoTemp optimizes

---
*Part of the [vasic-digital catalogue](../README.md)*
