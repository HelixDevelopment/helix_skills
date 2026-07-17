# ReplayBuffer

- **GitHub URL**: <https://github.com/vasic-digital/ReplayBuffer>
- **Description**: Reusable ReplayBuffer module for visual testing and automation
- **Category**: Testing + QA + Benchmarking
- **Status**: Active

## Overview

ReplayBuffer provides a storage and replay mechanism for UI interaction sequences captured during automated test runs. It records user interactions (clicks, inputs, navigation) into structured buffers that can be deterministically replayed for regression testing, debugging, and training data generation. Designed as the persistence layer in the visual testing pipeline between capture and replay.

## Tech Stack

- Language: Multiple
- Framework: Visual testing pipeline, interaction recording

## Key Features

- Structured recording of UI interaction sequences
- Deterministic replay of captured interaction buffers
- Buffer management with configurable retention and compression
- Integration with visual regression and training data pipelines
- Support for cross-platform interaction format normalisation

## Related Repos

- [TrainingCollector](../TrainingCollector/README.md) — captures interaction data that ReplayBuffer stores
- [VisualRegression](../VisualRegression/README.md) — regression testing that replays buffers against current UI state
- [ScreenDiff](../ScreenDiff/README.md) — visual comparison applied to replayed interaction screenshots

---
*Part of the [vasic-digital catalogue](../README.md)*
