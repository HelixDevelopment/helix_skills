# HelixCode

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixCode>
- **Description**: AI coding agent that autonomously writes, reviews, tests, and refactors code across multiple programming languages
- **Category**: AI / Development Tools
- **Status**: Active

## Overview

HelixCode is an autonomous coding agent that provides intelligent code generation with context-aware completions, automated code review with security and quality analysis, and self-healing test generation that adapts to project conventions. Built on HelixAgent as its cognitive substrate, it understands git context, project structure, and coding conventions to produce high-quality code across multiple languages.

## Tech Stack

- Language: Go (core engine), TypeScript (editor integrations)
- Protocols: LSP for editor integration, git for version control awareness
- Architecture: Agent-based with tool-use for file operations, terminal commands, and git
- Key patterns: RAG-based code context, tree-sitter AST analysis

## Key Features

- Autonomous code generation from natural language specifications and task descriptions
- Context-aware code completions that understand project structure and conventions
- Automated code review with security vulnerability detection and quality scoring
- Self-healing test generation -- writes tests that catch real regressions, not just syntax
- Multi-language support across Go, Python, TypeScript, Java, Rust, and more

## Related Repos

- [HelixAgent](../HelixAgent/README.md) -- cognitive substrate powering HelixCode
- [LLMProvider](../LLMProvider/README.md) -- multi-model code generation (different models for different tasks)
- [HelixMemory](../HelixMemory/README.md) -- project context persistence across sessions
- [HelixGitpx](../HelixGitpx/README.md) -- advanced git operations and branch management
- [VisionEngine](../VisionEngine/README.md) -- screenshot-to-code and UI-driven development
- [HelixBuilder](../HelixBuilder/README.md) -- consumes HelixCode for AI-assisted build pipeline code generation

---
*Part of the [HelixDevelopment catalogue](../README.md)*
