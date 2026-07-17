# ToolSchema

> **Repo:** [vasic-digital/ToolSchema](https://github.com/vasic-digital/ToolSchema)
> **Type:** vasic-digital repo · **Status:** Active

## Overview

A library for defining, validating, and executing tool schemas in
AI agent systems. Provides a standardised format for declaring
agent-callable tools with input/output validation, enabling reliable
tool use across different agent frameworks.

## Key capabilities

- Declarative tool schema definition format
- Input and output validation against schemas
- Tool execution orchestration and error handling
- Schema versioning and compatibility checks

## Architecture

ToolSchema is structured as a schema-driven tool framework:

1. **Schema definition** — declarative format for describing tool
   inputs, outputs, and capabilities
2. **Validator** — runtime validation of tool calls against
   declared schemas
3. **Execution engine** — orchestrates tool invocation with
   error handling and result validation
4. **Schema registry** — stores and versions tool schemas for
   discovery and compatibility

## Integration points

- **SkillRegistry** — skill registration using tool schemas for
  interface definitions
- **MCP_Module** — MCP tool registration using ToolSchema
  definitions
- **Agentic** — graph workflow nodes reference tool schemas for
  task definitions
- **JVM-Toolkit** — JVM-level abstractions for tool implementations
- **Planning** — planning algorithms use tool schemas to reason
  about available capabilities

## Configuration

Schema format options, validation strictness, and execution
policies are configurable. Check the repo for schema format
documentation and API reference.

## Status

**Active.** Standalone repository. Referenced in the Helix Skills
ecosystem as a core tool definition primitive.
