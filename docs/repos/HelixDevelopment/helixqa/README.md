# helixqa

- **GitHub URL**: <https://github.com/HelixDevelopment/helixqa>
- **Description**: AI-driven QA orchestration for multi-platform testing with anti-bluff evidence enforcement
- **Category**: QA / Testing
- **Status**: Active

## Overview

helixqa is the QA orchestration backbone of the Helix ecosystem. It provides autonomous test planning, execution, and reporting across unit, integration, E2E, stress, chaos, and visual regression test types. Its defining characteristic is anti-bluff enforcement -- every PASS must carry captured evidence (audio, video, screenshots, logs), ensuring that test success genuinely means the feature works for the end user.

## Tech Stack

- Language: Go (orchestration engine), Bash (test harnesses), Python (ML-based analysis)
- Streaming: JSONL event channels for real-time test execution monitoring
- Architecture: Orchestrator with pluggable test runners and evidence collectors
- Key patterns: Risk-descending execution, paired mutation testing, evidence-chain validation

## Key Features

- Multi-platform test orchestration across desktop, mobile, web, and embedded targets
- Autonomous test planning with risk-ordered execution priority
- Anti-bluff enforcement -- every PASS must carry captured evidence
- Challenge system -- structured validation challenges with scoring and evidence requirements
- Visual regression testing with golden-image comparison and OCR-based content verification

## Related Repos

- [DocProcessor](../DocProcessor/README.md) -- specification-driven test planning
- [HelixAgent](../HelixAgent/README.md) -- AI-assisted test generation and failure diagnosis
- [VisionEngine](../VisionEngine/README.md) -- screenshot-based visual regression and OCR validation
- [LLMOrchestrator](../LLMOrchestrator/README.md) -- headless test agent management
- [HelixConstitution](../HelixConstitution/README.md) -- testing mandate enforcement
- [HelixBuilder](../HelixBuilder/README.md) -- post-build validation in the build pipeline
- [HelixQA](../../vasic-digital/HelixQA/README.md) -- related QA framework in the vasic-digital ecosystem

---
*Part of the [HelixDevelopment catalogue](../README.md)*
