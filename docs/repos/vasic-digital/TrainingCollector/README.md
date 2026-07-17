# TrainingCollector

- **GitHub URL**: <https://github.com/vasic-digital/TrainingCollector>
- **Description**: Reusable TrainingCollector module for visual testing and automation
- **Category**: Testing + QA + Benchmarking
- **Status**: Active

## Overview

TrainingCollector gathers and structures visual training data from automated test runs and UI interactions. It captures screenshots, interaction sequences, and UI state transitions into organised datasets suitable for training visual regression models and improving automated test generation. Works as a data-collection layer in the visual testing pipeline.

## Tech Stack

- Language: Multiple
- Framework: Visual testing pipeline

## Key Features

- Automated capture of UI screenshots and interaction sequences
- Structured dataset generation from test run artifacts
- State transition tracking for UI workflow documentation
- Integration with visual regression and replay testing pipelines
- Configurable capture frequency and output formats

## Related Repos

- [ReplayBuffer](../ReplayBuffer/README.md) — replay storage for captured interaction sequences
- [ScreenDiff](../ScreenDiff/README.md) — visual diff analysis applied to collected training data
- [VisualRegression](../VisualRegression/README.md) — regression detection that consumes training datasets

---
*Part of the [vasic-digital catalogue](../README.md)*
