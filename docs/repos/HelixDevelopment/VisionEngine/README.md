# VisionEngine

- **GitHub URL**: <https://github.com/HelixDevelopment/VisionEngine>
- **Description**: Computer vision and LLM Vision engine for UI analysis -- OCR, template matching, visual regression testing, and pixel-oracle pattern for non-introspectable UIs
- **Category**: AI / Testing
- **Status**: Active

## Overview

VisionEngine combines traditional computer vision techniques (template matching, OCR, edge detection) with LLM-powered vision capabilities for automated UI interaction, screenshot analysis, and visual regression testing. It provides the pixel-oracle fallback for non-introspectable UIs (games, canvas, GL surfaces) per HelixConstitution §11.4.117, enabling testing where accessibility trees are empty or unreliable.

## Tech Stack

- Language: Go (engine core), Python (ML/CV model integration)
- CV: OpenCV bindings, Tesseract OCR, LLM Vision APIs (GPT-4V, Claude Vision)
- Architecture: Pipeline-based with pluggable CV processors and LLM vision backends
- Key patterns: ROI-based processing, confidence-threshold filtering, golden-image diffing

## Key Features

- OCR with per-word confidence scoring and region-of-interest targeting
- Template matching for UI element detection and click-target identification
- LLM Vision integration for natural-language visual understanding and scene description
- Visual regression testing with golden-image comparison and perceptual hashing
- Pixel-oracle pattern for non-introspectable UIs (games, canvas, GL surfaces)

## Related Repos

- [helixqa](../helixqa/README.md) -- visual regression testing and automated screenshot validation
- [HelixAgent](../HelixAgent/README.md) -- visual reasoning and UI understanding capabilities
- [HelixCode](../HelixCode/README.md) -- screenshot-to-code and UI-driven development
- [HelixPlay](../HelixPlay/README.md) -- automated game screen testing
- [DocProcessor](../DocProcessor/README.md) -- visual document analysis and OCR extraction

---
*Part of the [HelixDevelopment catalogue](../README.md)*
