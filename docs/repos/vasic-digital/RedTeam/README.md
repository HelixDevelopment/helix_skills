# RedTeam

- **GitHub URL**: <https://github.com/vasic-digital/RedTeam>
- **Description**: YAML-driven adversarial prompt fixture harness for defensive LLM guardrail regression testing
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

RedTeam provides a YAML-driven test harness for adversarial prompt regression testing against LLM guardrails. It defines adversarial fixtures in declarative YAML, runs them against guardrail implementations, and reports pass/fail results for continuous security validation. Designed as the red-team testing layer in the defensive LLM pipeline, ensuring guardrails remain effective as models and attack techniques evolve.

## Tech Stack

- Language: Multiple
- Config: YAML-driven fixtures

## Key Features

- Declarative YAML adversarial prompt fixture definitions
- Guardrail regression test execution with structured reporting
- CI/CD integration for continuous security validation
- Fixture versioning and categorisation by attack class
- Extensible fixture library for new adversarial patterns

## Related Repos

- [Normalize](../Normalize/README.md) — input canonicalisation that RedTeam fixtures exercise and validate
- [Claritas](../Claritas/README.md) — system-prompt extraction detection tested by adversarial fixtures
- [LeakHub](../LeakHub/README.md) — prompt-leak corpus providing source material for RedTeam fixture generation

---
*Part of the [vasic-digital catalogue](../README.md)*
