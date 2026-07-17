# Helix Skills

**Revision:** 2
**Last modified:** 2026-07-18T00:30:00Z

Skills system for CLI AI Agents. Inherits the
[Helix Constitution](https://github.com/HelixDevelopment/HelixConstitution)
as the `constitution/` submodule — every universal rule from
`constitution/CLAUDE.md` and `constitution/Constitution.md` applies
unconditionally.

---

## Package Inventory

### Direct Git Submodules

| Repo | URL | Path |
|---|---|---|
| **HelixConstitution** | `git@github.com:HelixDevelopment/HelixConstitution.git` | `constitution/` |

### Constitution Skills (7 installed via `register.sh`)

| Skill | Directory | Complexity |
|---|---|---|
| action-prefix-system | `constitution/skills/action-prefix-system/` | intermediate |
| media-validator | `constitution/skills/media-validator/` | intermediate |
| multitrack | `constitution/skills/multitrack/` | advanced |
| reporting-workable-items | `constitution/skills/reporting-workable-items/` | intermediate |
| scheduled-work-queue | `constitution/skills/scheduled-work-queue/` | intermediate |
| session-sync | `constitution/skills/session-sync/` | advanced |
| workable-item-lifecycle | `constitution/skills/workable-item-lifecycle/` | intermediate |

### Draft Skills (in INDEX.md, pending activation)

| Skill | Domain |
|---|---|
| android.overview | android |
| java.language | language |
| kotlin.language | language |
| linux.os | os |

### MCP Tool Servers (2)

| Server | Definition |
|---|---|
| media-validator | `constitution/mcp/media-validator-mcp.json` |
| scheduled-work | `constitution/mcp/scheduled-work-mcp.json` |

### Claude Code Plugins (2)

| Plugin | Source |
|---|---|
| helix | `constitution/plugins/helix/` |
| scheduled-work | `constitution/plugins/scheduled-work/` |

### Depth-1 Reusable Engines (4)

| Engine | Path | Status |
|---|---|---|
| continuum | `constitution/submodules/continuum/` | implemented |
| session_orchestrator | `constitution/submodules/session_orchestrator/` | design |
| token_optimizer | `constitution/submodules/token_optimizer/` | design |
| clickup_sync | `constitution/submodules/clickup_sync/` | design (Phase 0) |

### Declared Dependencies (6)

| Repo | Org | Required by |
|---|---|---|
| TOON | vasic-digital | token_optimizer |
| Embeddings | vasic-digital | token_optimizer |
| VectorDB | vasic-digital | token_optimizer |
| Normalize | vasic-digital | token_optimizer |
| LLMProvider | HelixDevelopment | token_optimizer |
| conversation | vasic-digital | token_optimizer |

---

## Documentation

- [Skills Catalog](docs/repos/skills.md) — full repo/tool/skill inventory
- [Skills Index](docs/skills/INDEX.md) — auto-generated skill graph index
- [Skills-Bearing Repos](docs/repos/skills/README.md) — per-repo detail pages
- [Gaps & Risks Register](docs/research/mvp/Agent_AI_Skill_Tree_Development/GAPS_AND_RISKS_REGISTER.md) — 136 tracked findings (95 open, 40 fixed, 1 N/A)
- [Constitution Integration Guide](docs/guides/HELIX_SKILLS_CONSTITUTION.md)

### Tracked-Items + Status Documents

| Document | Last modified | Revision | Markdown | HTML | PDF |
|---|---|---|---|---|---|
| Gaps & Risks Register | 2026-07-17T22:00:00Z | 9 | [md](docs/research/mvp/Agent_AI_Skill_Tree_Development/GAPS_AND_RISKS_REGISTER.md) | — | — |
| CONTINUATION | 2026-07-18T00:30:00Z | 4 | [md](CONTINUATION.md) | — | — |
| Skills Catalog | — | — | [md](docs/skills/README.md) | — | — |
| Repo Detail Index | — | — | [md](docs/repos/skills/README.md) | — | — |

---

## Upstream Mirrors

| Host | URL |
|---|---|
| GitHub (HelixDevelopment) | `git@github.com:HelixDevelopment/helix_skills.git` |
| GitLab (helixdevelopment1) | `git@gitlab.com:helixdevelopment1/helix_skills.git` |
| GitFlic | `git@gitflic.ru:helixdevelopment/helix_skills.git` |
| GitVerse | `git@gitverse.ru:helixdevelopment/helix_skills.git` |

---

## Quick Start

```bash
# Clone with submodules
git clone --recurse-submodules git@github.com:HelixDevelopment/helix_skills.git
cd helix_skills

# Install upstreams (§11.4.36)
install_upstreams

# Register constitution skills
for s in constitution/skills/*/register.sh; do bash "$s"; done
```
