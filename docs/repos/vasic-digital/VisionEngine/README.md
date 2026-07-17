# VisionEngine

- **GitHub URL**: <https://github.com/vasic-digital/VisionEngine>
- **Description**: Computer vision and LLM Vision for UI analysis and navigation
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

VisionEngine combines traditional computer vision techniques with LLM Vision capabilities to analyse, understand, and navigate user interfaces. It provides screen capture analysis, UI element detection, OCR, and visual state comparison — enabling AI agents to interact with applications that lack programmatic APIs or accessibility trees. Designed as the pixel-oracle fallback for non-introspectable UIs.

## Tech Stack

- Language: Multiple
- Framework: Computer vision, LLM Vision APIs

## Key Features

- Screen capture and visual state analysis for UI understanding
- UI element detection via template matching and OCR
- Visual diff and state comparison for change detection
- LLM Vision integration for high-level UI comprehension
- Pixel-level interaction driving for non-accessible applications

## Related Repos

- [ScreenDiff](../ScreenDiff/README.md) — visual regression detection used by VisionEngine for state comparison
- [VisualRegression](../VisualRegression/README.md) — regression testing harness that consumes VisionEngine's visual analysis
- [Panoptic](../Panoptic/README.md) — automated testing tool that uses VisionEngine for cross-platform UI interaction

---
*Part of the [vasic-digital catalogue](../README.md)*
