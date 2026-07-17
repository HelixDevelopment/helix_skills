# LLMOrchestrator

- **GitHub URL**: <https://github.com/HelixDevelopment/LLMOrchestrator>
- **Description**: Headless CLI agent management for LLM orchestration. Provides the lifecycle management layer for LLM-powered agents -- spawning, monitoring, resuming, and controlling headless agent sessions with crash resilience, rate-limit handling, and multi-alias support.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Headless agent session lifecycle management (spawn, monitor, resume, terminate)
- Multi-alias support with native-first priority and provider fallback chains
- Rate-limit detection and automatic rebind to healthy aliases
- Crash-resilient session recovery with durable state persistence
- Per-subscription auth management with token isolation
- Real-time session monitoring via JSONL event stream consumption
- Background agent pool management with configurable concurrency limits
- Session transcript capture for debugging and audit trails
- Watchdog system for detecting and recovering stalled agents

## Technology

- **Language**: Go
- **Frameworks**: Go concurrency (goroutines, channels), process management
- **Architecture**: Orchestrator pattern with pluggable agent backends
- **Key patterns**: Supervisor tree, circuit breaker for rate limits, temp-then-rename state

## Integration

- Manages HelixAgent instances as headless workers in multi-track development
- Uses LLMProvider for alias resolution and model routing
- Integrates with HelixConstitution multi-track orchestration rules (§11.4.187)
- Provides the conductor layer for parallel-development methodology
- Connects to HelixGitpx for agent git operations during autonomous workflows
- Session state feeds into the CONTINUATION.md resumption system
- Rate-limit tracking integrates with the native-alias-first priority system

## Status

Active development. Core session lifecycle management stable. Multi-alias support with rate-limit handling operational. Crash recovery and durable state persistence production-ready. Ongoing work on enhanced monitoring and pool optimization.
