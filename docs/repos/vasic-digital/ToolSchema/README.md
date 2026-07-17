# ToolSchema

- **GitHub URL**: <https://github.com/vasic-digital/ToolSchema>
- **Description**: Generic tool schema definition, validation, and execution for AI agent tool systems
- **Category**: AI / LLM Provider + Agent
- **Status**: Active

## Overview

ToolSchema provides a unified schema format for defining, validating, and executing tools that AI agents can invoke. It standardises how tool inputs and outputs are described, validates arguments against declared schemas before execution, and provides a common execution harness. This ensures agent tool integrations are type-safe and self-documenting, reducing the surface area for runtime errors in tool-calling loops.

## Tech Stack

- Language: Go
- Module: digital.vasic.toolschema

## Key Features

- Declarative tool schema definition with JSON Schema compatibility
- Input argument validation before tool execution
- Common execution harness with error handling and result formatting
- Self-documenting tool descriptions for LLM function-calling integration
- Composable schema registry for multi-tool agent environments

## Related Repos

- [SkillRegistry](../SkillRegistry/README.md) — registers skills whose tool invocations are validated by ToolSchema
- [LLMProvider](../LLMProvider/README.md) — provider adapters that use tool schemas for function-calling APIs
- [conversation](../conversation/README.md) — conversation context that tracks tool call results validated by this module

---
*Part of the [vasic-digital catalogue](../README.md)*
