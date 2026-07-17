# SkillRegistry

- **GitHub URL**: <https://github.com/vasic-digital/SkillRegistry>
- **Description**: CLI agent skill registration and management for AI agent systems
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

SkillRegistry provides a centralised registry for discovering, loading, and managing skills available to CLI-based AI coding agents. It handles skill metadata, version tracking, dependency resolution, and runtime activation so agent orchestrators can dynamically compose capabilities without hardcoding skill references. The registry is designed to be consumed by any agent framework that needs pluggable skill discovery.

## Tech Stack

- Language: Go
- Module: digital.vasic.skillregistry

## Key Features

- Centralised skill registration with metadata and version tracking
- Dynamic skill discovery and runtime activation
- Dependency resolution between skills
- Pluggable skill loader supporting multiple source formats
- Integration hooks for agent orchestrators and workflow engines

## Related Repos

- [ToolSchema](../ToolSchema/README.md) — defines the schema and validation for tool invocations that skills expose
- [Agentic](../Agentic/README.md) — graph-based orchestration that composes registered skills into workflows
- [LLMOrchestrator](../LLMOrchestrator/README.md) — headless agent management that loads skills from the registry

---
*Part of the [vasic-digital catalogue](../README.md)*
