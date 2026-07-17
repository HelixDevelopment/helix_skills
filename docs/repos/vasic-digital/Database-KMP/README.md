# Database-KMP

- **GitHub URL**: <https://github.com/vasic-digital/Database-KMP>
- **Description**: digital.vasic.database - KMP network storage database interfaces and entity types
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

Database-KMP provides Kotlin Multiplatform interfaces and entity types for network-based database access. It defines platform-agnostic abstractions for database operations, entity mapping, and query patterns that each target implements against its native database layer. Published as `digital.vasic.database`.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, kotlinx.serialization

## Key Features

- Cross-platform database access interfaces
- Entity type definitions with serialization support
- Platform-agnostic query and mutation abstractions
- Shared database model across all KMP targets

## Related Repos

- [Storage-KMP](../Storage-KMP/README.md) — higher-level storage service interfaces that Database-KMP builds on
- [Config-KMP](../Config-KMP/README.md) — configuration types for database connection settings
- [Concurrency-KMP](../Concurrency-KMP/README.md) — async utilities used by database operations

---
*Part of the [vasic-digital catalogue](../README.md)*
