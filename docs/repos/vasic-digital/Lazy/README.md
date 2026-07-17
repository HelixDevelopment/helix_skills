# Lazy

- **GitHub URL**: <https://github.com/vasic-digital/Lazy>
- **Description**: Generic reusable Go module: digital.vasic.lazy - Lazy loading with sync.Once generics
- **Category**: Messaging + Observability + Storage
- **Status**: Active

## Overview

A Go module providing generic lazy-loading primitives built on sync.Once. Enables deferred initialization of expensive resources with thread-safe guarantees, ensuring values are computed only once on first access.

## Tech Stack

- Language: Go
- Module path: digital.vasic.lazy
- Built on: sync.Once

## Key Features

- Generic lazy-loaded value containers
- Thread-safe initialization via sync.Once
- Deferred resource computation on first access

## Related Repos

- [Assets](../Assets/README.md) — lazy asset loading with strategy-based resolution
- [cache](../cache/README.md) — caching layer for computed values
- [concurrency](../concurrency/README.md) — concurrent primitives

---
*Part of the [vasic-digital catalogue](../README.md)*
