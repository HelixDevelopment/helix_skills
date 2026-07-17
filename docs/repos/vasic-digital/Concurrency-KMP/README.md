# Concurrency-KMP

- **GitHub URL**: <https://github.com/vasic-digital/Concurrency-KMP>
- **Description**: Kotlin Multiplatform concurrency utilities: lazy loading, platform synchronization, flow-based loaders
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

Concurrency-KMP provides cross-platform concurrency utilities for Kotlin Multiplatform projects. It abstracts platform-specific synchronization primitives into a common API, offers lazy-loading constructs that work across targets, and provides flow-based data loaders that integrate naturally with Kotlin coroutines and Flow.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, Kotlin Coroutines, Flow

## Key Features

- Platform-agnostic synchronization primitives
- Lazy loading with coroutine-aware initialization
- Flow-based data loaders for reactive async operations
- Shared concurrency utilities across all KMP targets

## Related Repos

- [RateLimiter-KMP](../RateLimiter-KMP/README.md) — rate limiting primitives that build on concurrency utilities
- [Storage-KMP](../Storage-KMP/README.md) — storage interfaces that use concurrency utilities for async operations

---
*Part of the [vasic-digital catalogue](../README.md)*
