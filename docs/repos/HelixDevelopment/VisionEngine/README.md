# VisionEngine

- **GitHub URL**: <https://github.com/HelixDevelopment/VisionEngine>
- **Description**: Computer vision and LLM Vision engine for UI analysis, navigation, and visual understanding. Combines traditional CV techniques (template matching, OCR, edge detection) with LLM-powered vision capabilities for automated UI interaction, screenshot analysis, and visual regression testing.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Screenshot capture and analysis with multi-monitor and multi-window support
- OCR (Optical Character Recognition) with per-word confidence scoring and region-of-interest targeting
- Template matching for UI element detection and click-target identification
- LLM Vision integration for natural-language visual understanding and scene description
- Visual regression testing with golden-image comparison and perceptual hashing
- UI automation -- drive applications by visual recognition when accessibility trees are empty
- Screen recording with frame extraction and analysis pipelines
- Image preprocessing (cropping, scaling, contrast adjustment) for improved recognition
- Pixel-oracle pattern for non-introspectable UIs (games, canvas, GL surfaces)

## Technology

- **Language**: Go (engine core), Python (ML/CV model integration)
- **Frameworks**: OpenCV bindings, Tesseract OCR, LLM Vision APIs (GPT-4V, Claude Vision)
- **Architecture**: Pipeline-based with pluggable CV processors and LLM vision backends
- **Key patterns**: ROI-based processing, confidence-threshold filtering, golden-image diffing

## Integration

- Used by helixqa for visual regression testing and automated screenshot validation
- Integrates with HelixAgent for visual reasoning and UI understanding capabilities
- Consumed by HelixCode for screenshot-to-code generation and UI-driven development
- Provides the pixel-oracle fallback for non-introspectable UIs per HelixConstitution §11.4.117
- Connects to HelixPlay and game projects for automated game screen testing
- Feeds visual evidence into the anti-bluff captured-evidence chain
- Used by DocProcessor for visual document analysis and OCR extraction

## Status

Active development. Core screenshot capture and OCR operational. Template matching and golden-image comparison stable. LLM Vision integration functional with multiple providers. UI automation via visual recognition is a key capability under continuous improvement. Per-word OCR confidence scoring and ROI targeting production-ready.
