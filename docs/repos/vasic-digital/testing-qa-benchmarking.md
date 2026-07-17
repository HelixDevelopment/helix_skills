# Testing + QA + Benchmarking

Back to [vasic-digital index](./README.md) | [Main index](../README.md)

This group provides the visual testing, regression detection, and QA automation infrastructure that underpins the Helix anti-bluff validation pipeline. It includes cross-platform UI recording and screenshot capture, pixel-level visual diffing, replay buffer management for test reproducibility, training data collection for vision models, and the Challenges framework that enforces anti-bluff quality gates. These modules are the evidence-capture layer -- every user-visible feature's PASS verdict depends on captured artifacts produced by this group.

| Repo | Description | Status |
|---|---|---|
| [Panoptic](./Panoptic/README.md) | Comprehensive cross-platform automated testing tool for UI recording and screenshot capture across web, desktop, and mobile applications. Drives application UIs via Playwright (web), platform-native automation (desktop), and ADB/UIAutomator (Android), capturing frames and interaction sequences as structured test evidence. | Active |
| [ReplayBuffer](./ReplayBuffer/README.md) | Reusable Go module providing a replay buffer for visual testing and automation. Stores captured UI frames, interaction events, and timing data in a circular buffer, enabling deterministic replay of test sequences for debugging flaky tests and reproducing visual regressions. | Active |
| [ScreenDiff](./ScreenDiff/README.md) | Reusable Go module for pixel-level screen diffing in visual testing pipelines. Compares captured screenshots against golden references using SSIM, perceptual hashing, and configurable tolerance thresholds, producing diff images and structured PASS/FAIL reports for visual regression detection. | Active |
| [TrainingCollector](./TrainingCollector/README.md) | Reusable Go module for collecting visual training data during test runs. Captures annotated screenshots, UI element bounding boxes, and interaction labels that feed into vision model training pipelines for OCR, UI element detection, and layout analysis. | Active |
| [VisualRegression](./VisualRegression/README.md) | Reusable Go module orchestrating visual regression test suites. Manages golden image baselines, runs screenshot comparisons across builds/commits, generates regression reports with diff visualizations, and integrates with CI pipelines to gate deployments on visual correctness. | Active |
| [challenges](./challenges/README.md) | Anti-bluff challenge framework providing a YAML-driven test bank where each challenge defines acceptance criteria, evidence requirements, and scoring rules. Challenges are executed by HelixQA sessions and produce structured PASS/FAIL verdicts backed by captured physical evidence -- the mechanism that prevents "tests pass but features broken" outcomes. | Active |

**Related skills:** [media-validator](../skills/media-validator.md)
