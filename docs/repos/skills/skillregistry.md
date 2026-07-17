# SkillRegistry

> **Repo:** [vasic-digital/SkillRegistry](https://github.com/vasic-digital/SkillRegistry)
> **Type:** vasic-digital repo · **Status:** Active

## Overview

A registration and discovery system for CLI agent skills. Enables
AI agents to register, discover, and invoke reusable skills at
runtime, forming the backbone of composable agent architectures.

## Key capabilities

- Skill registration and metadata management
- Runtime skill discovery and invocation
- Versioning and dependency resolution for skills
- Skill lifecycle management (install, update, uninstall)

## Architecture

SkillRegistry is structured as a service registry:

1. **Registry store** — persistent storage of skill definitions,
   metadata, and version history
2. **Discovery engine** — query interface for agents to find
   available skills by capability, name, or version
3. **Dependency resolver** — resolves skill dependency graphs and
   validates compatibility
4. **Invocation bridge** — dispatches skill calls to the appropriate
   runtime implementation

## Integration points

- **ToolSchema** — tool schema definitions used for registered skill
  interfaces
- **Agentic** — graph workflows can reference registered skills as
  node implementations
- **Planning** — planning algorithms selecting skills for execution
  based on task requirements
- **Helix Skills** — the Helix Skills project itself uses
  SkillRegistry patterns for skill management

## Configuration

Registry storage backend, skill search paths, and version policies
are configurable. Check the repo for CLI and API documentation.

## Status

**Active.** Standalone repository. Referenced in the Helix Skills
ecosystem as a core skill management primitive.
