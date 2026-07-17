# ScreenDiff

- **GitHub URL**: <https://github.com/vasic-digital/ScreenDiff>
- **Description**: Reusable ScreenDiff module for visual testing and automation
- **Category**: Testing + QA + Benchmarking
- **Status**: Active

## Overview

ScreenDiff provides pixel-level and perceptual image comparison for detecting visual regressions in UI applications. It compares captured screenshots against baseline images, highlights differences, and produces structured verdicts with configurable sensitivity thresholds. Designed as the core comparison engine for visual regression testing pipelines across web, desktop, and mobile targets.

## Tech Stack

- Language: Multiple
- Framework: Image comparison, visual testing pipeline

## Key Features

- Pixel-level and perceptual hash-based image comparison
- Configurable sensitivity thresholds for diff detection
- Structured diff output with highlighted change regions
- Baseline management for regression tracking over time
- Integration with visual testing and CI pipelines

## Related Repos

- [VisualRegression](../VisualRegression/README.md) — regression testing harness that uses ScreenDiff for change detection
- [VisionEngine](../VisionEngine/README.md) — computer vision engine that uses ScreenDiff for UI state comparison
- [TrainingCollector](../TrainingCollector/README.md) — captures screenshots that ScreenDiff compares against baselines

---
*Part of the [vasic-digital catalogue](../README.md)*
