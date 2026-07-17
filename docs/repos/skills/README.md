# Skills-Bearing Repositories

> Per-repo detail pages for every repository that contributes skills, MCP
> tool servers, plugins, or reusable engines to the Helix Skills ecosystem.
>
> Master catalog: [skills.md](../skills.md)

## Index

### Constitution Submodule

| Repo | Type | Path |
|---|---|---|
| [HelixConstitution](helixconstitution.md) | Git submodule | `constitution/` |

### Constitution Skills

| Skill | Type | Path |
|---|---|---|
| [action-prefix-system](action-prefix-system.md) | Skill | `constitution/skills/action-prefix-system/` |
| [media-validator](media-validator.md) | Skill + MCP server | `constitution/skills/media-validator/` |
| [multitrack](multitrack.md) | Skill | `constitution/skills/multitrack/` |
| [reporting-workable-items](reporting-workable-items.md) | Skill | `constitution/skills/reporting-workable-items/` |
| [scheduled-work-queue](scheduled-work-queue.md) | Skill | `constitution/skills/scheduled-work-queue/` |
| [session-sync](session-sync.md) | Skill | `constitution/skills/session-sync/` |
| [workable-item-lifecycle](workable-item-lifecycle.md) | Skill | `constitution/skills/workable-item-lifecycle/` |

### MCP Tool Servers

| Server | Type | Path |
|---|---|---|
| [media-validator-mcp](media-validator.md) | MCP server | `constitution/mcp/media-validator-mcp.json` |
| [scheduled-work-mcp](scheduled-work.md) | MCP server | `constitution/mcp/scheduled-work-mcp.json` |

### Plugins

| Plugin | Type | Path |
|---|---|---|
| [helix](helix-plugin.md) | Plugin | `constitution/plugins/helix/` |
| [scheduled-work](scheduled-work.md) | Plugin | `constitution/plugins/scheduled-work/` |

### Reusable Engines (Submodules)

| Engine | Type | Path |
|---|---|---|
| [continuum](continuum.md) | Go submodule | `constitution/submodules/continuum/` |
| [session_orchestrator](session-orchestrator.md) | Go submodule | `constitution/submodules/session_orchestrator/` |
| [token_optimizer](token-optimizer.md) | Go submodule | `constitution/submodules/token_optimizer/` |

### Declared Dependencies

| Repo | Type | Required by |
|---|---|---|
| [TOON](toon.md) | Go module | token_optimizer |
| [Embeddings](embeddings.md) | Go module | token_optimizer |
| [VectorDB](vectordb.md) | Go module | token_optimizer |
| [Normalize](normalize.md) | Go module | token_optimizer |
| [LLMProvider](llmprovider.md) | Go module | token_optimizer |
| [conversation](conversation.md) | Go module | token_optimizer |

### HelixDevelopment Repos

| Repo | Type |
|---|---|
| [SkillRegistry](skillregistry.md) | vasic-digital repo |
| [ToolSchema](toolschema.md) | vasic-digital repo |
| [Agentic](agentic.md) | vasic-digital repo |
| [AgentWrapper](agentwrapper.md) | vasic-digital repo |
| [LLMOrchestrator](llmorchestrator.md) | vasic-digital repo |
| [MCP_Module](mcp-module.md) | vasic-digital repo |
