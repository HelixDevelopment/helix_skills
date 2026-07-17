# Storage-KMP

- **GitHub URL**: <https://github.com/vasic-digital/Storage-KMP>
- **Description**: digital.vasic.storage - KMP network storage service interfaces and abstractions
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

Storage-KMP provides Kotlin Multiplatform interfaces and abstractions for network-based storage services. Published as `digital.vasic.storage`, it defines the contracts for storage operations -- upload, download, listing, and metadata management -- that each platform target implements against its native storage backend. This is the foundational storage abstraction layer for KMP applications.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, kotlinx.coroutines

## Key Features

- Platform-agnostic storage service interfaces
- Upload, download, and metadata management abstractions
- Network storage service contracts for cross-platform use
- Foundation layer for higher-level database and document modules

## Related Repos

- [Database-KMP](../Database-KMP/README.md) — database interfaces built on top of Storage-KMP abstractions
- [Config-KMP](../Config-KMP/README.md) — configuration types for storage connection and protocol settings
- [Concurrency-KMP](../Concurrency-KMP/README.md) — async utilities used by storage operations

---
*Part of the [vasic-digital catalogue](../README.md)*
