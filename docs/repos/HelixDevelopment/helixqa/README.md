# helixqa

- **GitHub URL**: <https://github.com/HelixDevelopment/helixqa>
- **Description**: AI-driven QA orchestration framework for multi-platform testing. Provides autonomous test planning, execution, and reporting across unit, integration, E2E, stress, chaos, and visual regression test types. Integrates with the Helix Constitution's anti-bluff mandate to ensure every test produces captured evidence.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Multi-platform test orchestration across desktop, mobile, web, and embedded targets
- Autonomous test planning with risk-ordered execution priority
- Test bank management -- register, organize, and execute test suites per feature/domain
- Anti-bluff enforcement -- every PASS must carry captured evidence (audio, video, screenshots, logs)
- Challenge system -- structured validation challenges with scoring and evidence requirements
- Stress and chaos test orchestration with fault injection and recovery validation
- Visual regression testing with golden-image comparison and OCR-based content verification
- Real-time test execution streaming via JSONL event channels
- Test result aggregation with coverage analysis and gap identification

## Technology

- **Language**: Go (orchestration engine), Bash (test harnesses), Python (ML-based analysis)
- **Frameworks**: Go concurrency for parallel test execution, JSONL event streaming
- **Architecture**: Orchestrator with pluggable test runners and evidence collectors
- **Key patterns**: Risk-descending execution, paired mutation testing, evidence-chain validation

## Integration

- Consumes DocProcessor for specification-driven test planning
- Uses HelixAgent for AI-assisted test generation and failure diagnosis
- Integrates with VisionEngine for screenshot-based visual regression and OCR validation
- Connects to LLMOrchestrator for headless test agent management
- Enforces HelixConstitution testing mandates (four-layer coverage, anti-bluff, captured evidence)
- Feeds test results into the workable-items system for defect tracking
- Used by HelixBuilder for post-build validation in the build pipeline

## Status

Active development. Core orchestration engine and challenge system are stable. Multi-platform test execution operational. Anti-bluff evidence chain enforcement is a key differentiator. Visual regression and stress/chaos testing capabilities continuously expanding.
