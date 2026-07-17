# HelixCode

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixCode>
- **Description**: AI coding agent that autonomously writes, reviews, tests, and refactors code across multiple programming languages. Provides intelligent code generation with context-aware completions, automated code review with security and quality analysis, and self-healing test generation that adapts to project conventions.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Autonomous code generation from natural language specifications and task descriptions
- Context-aware code completions that understand project structure and conventions
- Automated code review with security vulnerability detection and quality scoring
- Self-healing test generation -- writes tests that catch real regressions, not just syntax
- Multi-language support across Go, Python, TypeScript, Java, Rust, and more
- Refactoring engine with blast-radius analysis and safe transformation guarantees
- Git-aware operations -- understands branches, diffs, and merge contexts
- Interactive debugging assistance with error diagnosis and fix suggestions

## Technology

- **Language**: Go (core engine), TypeScript (editor integrations)
- **Frameworks**: Go concurrency, LSP protocol integration
- **Architecture**: Agent-based with tool-use for file operations, terminal commands, and git
- **Key patterns**: RAG-based code context, tree-sitter AST analysis

## Integration

- Built on HelixAgent as its cognitive substrate
- Consumes LLMProvider for multi-model code generation (different models for different tasks)
- Uses HelixMemory for project context persistence across sessions
- Integrates with HelixGitpx for advanced git operations and branch management
- Connects to VisionEngine for screenshot-to-code and UI-driven development
- Consumed by HelixBuilder for AI-assisted build pipeline code generation

## Status

Active development. Core code generation and review capabilities are stable. Multi-language support expanding. Integration with editor environments (VS Code, JetBrains) via LSP protocol. Self-healing test generation is a key differentiator under continuous improvement.
