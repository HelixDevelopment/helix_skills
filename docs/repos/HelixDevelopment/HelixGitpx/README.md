# HelixGitpx

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixGitpx>
- **Description**: Helix Git Proxy eXtended -- advanced git operations layer with intelligent workflow automation, multi-upstream management, and safe force-push alternatives
- **Category**: DevOps / Git Tooling
- **Status**: Active

## Overview

HelixGitpx extends standard git with AI-assisted operations and safety guards. It provides multi-upstream push management, intelligent merge conflict resolution, stale lock detection, and branch consistency enforcement across main repos and owned submodules. The library is the git operations backbone for the multi-track parallel-development methodology, ensuring safe, atomic, and auditable git workflows.

## Tech Stack

- Language: Go (core library), Bash (hooks and wrappers)
- Git integration: go-git library, shell script hooks
- Architecture: Library + CLI tool with hook-based integration points
- Key patterns: Advisory locking, temp-then-rename atomicity, pre/post-operation guards

## Key Features

- Multi-upstream push management -- fan out commits to all configured remotes atomically
- Intelligent merge conflict resolution with AI-assisted union-preserving strategies
- Stale lock detection and auto-reaping for dead-holder git locks
- Commit wrapper with staged-file auditing and forbidden-pattern detection
- Branch consistency enforcement across main repo and owned submodules

## Related Repos

- [HelixCode](../HelixCode/README.md) -- uses HelixGitpx for git-aware code operations
- [HelixConstitution](../HelixConstitution/README.md) -- enforcement hooks for branch consistency
- [LLMOrchestrator](../LLMOrchestrator/README.md) -- agent git operations during autonomous workflows
- [HelixAgent](../HelixAgent/README.md) -- git-aware agent operations

---
*Part of the [HelixDevelopment catalogue](../README.md)*
