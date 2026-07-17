# HelixConstitution

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixConstitution>
- **Description**: Universal agent rules, constraints, and governance framework providing the constitutional foundation for all Helix AI agents
- **Category**: Governance / Framework
- **Status**: Active

## Overview

HelixConstitution defines the CLAUDE.md, AGENTS.md, and equivalent governance files that enforce anti-bluff policies, testing mandates, safety constraints, and operational discipline across every agent and project in the Helix ecosystem. With over 213 numbered constitutional anchors, it is the single source of truth for how agents must behave -- ensuring every output carries positive evidence, every test produces captured proof, and every operation respects host and data safety.

## Tech Stack

- Language: Markdown (governance files), Bash (enforcement scripts), Go (tooling)
- Distribution: Git submodule with auto-propagation hooks
- Architecture: Constitutional hierarchy with universal rules in the submodule and project-specific extensions in consuming repos
- Key patterns: Submodule inheritance, paired mutation testing, gate-based enforcement

## Key Features

- Universal governance rules applicable to any AI agent platform (Claude, Cursor, Aider, etc.)
- Anti-bluff covenant enforcement -- every agent output must carry positive evidence
- Testing mandate framework with four-layer coverage requirements
- Host-session safety constraints preventing destructive operations
- Workable-item tracking with SQLite-backed single source of truth

## Related Repos

- [HelixAgent](../HelixAgent/README.md) -- governed by HelixConstitution rules
- [helixqa](../helixqa/README.md) -- enforces testing mandates and validates compliance
- [LLMOrchestrator](../LLMOrchestrator/README.md) -- multi-track orchestration rules (§11.4.187)
- [HelixGitpx](../HelixGitpx/README.md) -- enforcement hooks for branch consistency and commit discipline
- [HelixSpecifier](../HelixSpecifier/README.md) -- specification-driven governance compliance

---
*Part of the [HelixDevelopment catalogue](../README.md)*
