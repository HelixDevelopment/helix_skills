# VisualRegression

- **GitHub URL**: <https://github.com/vasic-digital/VisualRegression>
- **Description**: Reusable VisualRegression module for visual testing and automation
- **Category**: Testing + QA + Benchmarking
- **Status**: Active

## Overview

VisualRegression provides a complete visual regression testing framework that compares current UI renders against baseline images to detect unintended visual changes. It orchestrates screenshot capture, comparison, threshold management, and reporting — producing structured verdicts that integrate into CI pipelines. Designed as the top-level harness in the visual testing pipeline that composes ScreenDiff, ReplayBuffer, and TrainingCollector.

## Tech Stack

- Language: Multiple
- Framework: Visual regression testing harness

## Key Features

- Automated baseline capture and comparison workflow
- Configurable pixel and perceptual diff thresholds
- Structured pass/fail verdicts with diff image artifacts
- CI pipeline integration for continuous visual regression detection
- Cross-platform support for web, desktop, and mobile targets

## Related Repos

- [ScreenDiff](../ScreenDiff/README.md) — core image comparison engine used by VisualRegression
- [ReplayBuffer](../ReplayBuffer/README.md) — interaction replay for deterministic UI state reproduction
- [VisionEngine](../VisionEngine/README.md) — computer vision engine providing advanced visual analysis

---
*Part of the [vasic-digital catalogue](../README.md)*
