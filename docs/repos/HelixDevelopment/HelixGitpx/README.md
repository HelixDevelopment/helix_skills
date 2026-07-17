# HelixGitpx

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixGitpx>
- **Description**: Helix Git Proxy eXtended -- an advanced git operations layer providing intelligent git workflow automation, multi-upstream management, conflict resolution assistance, and safe force-push alternatives. Extends standard git with AI-assisted operations and safety guards.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Multi-upstream push management -- fan out commits to all configured remotes atomically
- Intelligent merge conflict resolution with AI-assisted union-preserving strategies
- Safe branch operations with pre-operation backup and post-operation verification
- Stale lock detection and auto-reaping for dead-holder git locks
- Commit wrapper with staged-file auditing and forbidden-pattern detection
- Branch consistency enforcement across main repo and owned submodules
- Fetch-before-edit guards preventing stale-state operations
- Git history analysis tools (pickaxe search, blame enrichment, change-impact analysis)

## Technology

- **Language**: Go (core library), Bash (hooks and wrappers)
- **Frameworks**: Go git libraries (go-git), shell script integration
- **Architecture**: Library + CLI tool with hook-based integration points
- **Key patterns**: Advisory locking, temp-then-rename atomicity, pre/post-operation guards

## Integration

- Used by HelixCode for git-aware code operations and branch management
- Integrates with HelixConstitution enforcement hooks (guard-branch-consistency, guard-work-track-binding)
- Provides the commit/push wrappers used by the multi-track orchestration system
- Consumed by LLMOrchestrator for agent git operations during autonomous workflows
- Connects to the workable-items system for branch-to-logic-group binding enforcement
- Foundation for the parallel-development methodology's merge discipline

## Status

Active development. Core multi-upstream push and merge operations are stable. AI-assisted conflict resolution is operational. Stale lock reaping and branch consistency enforcement are production-ready. Ongoing work on enhanced history analysis tools.
