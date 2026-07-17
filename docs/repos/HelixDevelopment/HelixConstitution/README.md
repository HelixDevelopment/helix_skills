# HelixConstitution

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixConstitution>
- **Description**: Universal agent rules, constraints, and governance framework providing the constitutional foundation for all Helix AI agents. Defines CLAUDE.md, AGENTS.md, and equivalent governance files that enforce anti-bluff policies, testing mandates, safety constraints, and operational discipline across every agent and project in the ecosystem.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Universal governance rules applicable to any AI agent platform (Claude, Cursor, Aider, etc.)
- Anti-bluff covenant enforcement -- every agent output must carry positive evidence
- Testing mandate framework with four-layer coverage requirements
- Host-session safety constraints preventing destructive operations
- Documentation sync mandates ensuring all docs stay current
- Workable-item tracking with SQLite-backed single source of truth
- Auto-propagation hooks for distributing rule changes across consuming projects
- Multi-track orchestration rules for parallel development workflows

## Technology

- **Language**: Markdown (governance files), Bash (enforcement scripts), Go (tooling)
- **Frameworks**: Git submodule distribution, hook-based enforcement
- **Architecture**: Constitutional hierarchy with universal rules in the submodule and project-specific extensions in consuming repos
- **Key patterns**: Submodule inheritance, paired mutation testing, gate-based enforcement

## Integration

- Used as a git submodule by every Helix project and consuming repository
- Provides governance rules consumed by HelixAgent, HelixCode, and all other agents
- Enforcement scripts integrate with git hooks (pre-commit, pre-push, PreToolUse)
- Pairs with helixqa for testing mandate enforcement and validation
- Feeds the workable-items system via reporting directives (ISSUE/BUG/TASK/FEATURE)
- Auto-propagation via post_update_hook.sh on constitution pulls

## Status

Active development. Over 213 numbered constitutional anchors covering anti-bluff, testing, safety, documentation, and operational discipline. Continuous expansion driven by operator mandates and forensic incident remediation. Revision 55+ of the core governance files.
