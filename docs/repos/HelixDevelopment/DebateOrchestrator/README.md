# DebateOrchestrator

- **GitHub URL**: <https://github.com/HelixDevelopment/DebateOrchestrator>
- **Description**: Multi-agent debate orchestration library for structured argumentation and consensus-building across LLM agents. Implements structured debate protocols where multiple agents argue competing positions, challenge each other's reasoning, and converge on a well-justified conclusion. Currently in Phase 1 reconstruction with 7 real and 5 stub sub-packages.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Multi-agent debate protocol orchestration with configurable round-robin and adversarial modes
- Structured argumentation frameworks -- agents present claims, evidence, and rebuttals
- Consensus-building algorithms that weigh agent confidence and evidence quality
- Configurable debate depth (number of rounds) and participant count
- Argument scoring and ranking with weighted evaluation criteria
- Transcript capture for full auditability of the debate process
- Integration with LLMProvider for agent model diversity within a single debate
- Pluggable evaluation strategies for determining debate winners

## Technology

- **Language**: Go
- **Frameworks**: Go standard library, structured concurrency patterns
- **Architecture**: Library with 12 sub-packages (7 implemented, 5 stubs pending reconstruction)
- **Key patterns**: Pipeline-based debate flow, channel-based agent coordination

## Integration

- Consumes LLMProvider for multi-model agent instantiation within debates
- Used by HelixSpecifier for spec-driven decision-making via adversarial review
- Integrates with HelixAgent as a reasoning-pluggable cognitive module
- Feeds debate transcripts to HelixMemory for long-term decision recall
- Part of the Helix cognitive stack -- provides structured reasoning where LLMOrchestrator provides execution

## Status

Phase 1 reconstruction underway. Core orchestration and 7 sub-packages are implemented. 5 remaining sub-packages are stubbed with implementation tracked in RECONSTRUCTION_ROADMAP.md. Active development focus is on completing the reconstruction and adding evaluation strategy plugins.
