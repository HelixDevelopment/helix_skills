# DebateOrchestrator

- **GitHub URL**: <https://github.com/HelixDevelopment/DebateOrchestrator>
- **Description**: Multi-agent debate orchestration library for structured argumentation and consensus-building across LLM agents
- **Category**: AI / Cognitive Engine
- **Status**: Active

## Overview

DebateOrchestrator implements structured debate protocols where multiple LLM agents argue competing positions, challenge each other's reasoning, and converge on well-justified conclusions. It provides the adversarial reasoning layer for the Helix cognitive stack -- where HelixAgent provides single-agent cognition, DebateOrchestrator provides multi-agent consensus. Currently in Phase 1 reconstruction with 7 real and 5 stub sub-packages.

## Tech Stack

- Language: Go
- Concurrency: Go channels and goroutines for agent coordination
- Architecture: Library with 12 sub-packages (7 implemented, 5 stubs)
- Key patterns: Pipeline-based debate flow, round-robin and adversarial modes

## Key Features

- Multi-agent debate protocol orchestration with configurable round-robin and adversarial modes
- Structured argumentation frameworks -- agents present claims, evidence, and rebuttals
- Consensus-building algorithms that weigh agent confidence and evidence quality
- Argument scoring and ranking with weighted evaluation criteria
- Transcript capture for full auditability of the debate process

## Related Repos

- [HelixAgent](../HelixAgent/README.md) -- uses DebateOrchestrator as a reasoning-pluggable cognitive module
- [HelixSpecifier](../HelixSpecifier/README.md) -- consumes debate for spec-driven decision-making via adversarial review
- [LLMProvider](../LLMProvider/README.md) -- provides multi-model agent instantiation within debates
- [HelixMemory](../HelixMemory/README.md) -- stores debate transcripts for long-term decision recall
- [Agentic](../../vasic-digital/Agentic/README.md) -- graph-based workflow orchestration that can embed debate nodes

---
*Part of the [HelixDevelopment catalogue](../README.md)*
