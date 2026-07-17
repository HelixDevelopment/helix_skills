# HyperTune

- **GitHub URL**: <https://github.com/vasic-digital/HyperTune>
- **Description**: Hyperparameter tuning orchestration
- **Category**: Container + Lifecycle
- **Status**: SCAFFOLD / WIP

## Overview

HyperTune provides orchestration for hyperparameter tuning across LLM-based and ML workflows. It manages the search space definition, trial execution, and result tracking needed to find optimal configurations. The system is designed to work alongside other orchestration components for end-to-end tuning pipelines.

## Tech Stack

- Language: Multiple
- Framework: Custom tuning orchestration engine

## Key Features

- Configurable hyperparameter search space definition
- Parallel trial execution for efficient exploration
- Result tracking and comparison across tuning runs
- Integration with LLM orchestration and benchmarking tools

## Related Repos

- [AutoTemp](../AutoTemp/README.md) — temperature-specific tuning that HyperTune generalizes beyond
- [LLMOrchestrator](../LLMOrchestrator/README.md) — manages the agent runtimes that HyperTune tunes parameters for

---
*Part of the [vasic-digital catalogue](../README.md)*
