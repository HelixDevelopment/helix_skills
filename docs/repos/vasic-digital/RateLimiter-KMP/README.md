# RateLimiter-KMP

- **GitHub URL**: <https://github.com/vasic-digital/RateLimiter-KMP>
- **Description**: Kotlin Multiplatform rate limiting: semaphore, token bucket, adaptive, throttler
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

RateLimiter-KMP provides multiple rate limiting strategies as a Kotlin Multiplatform module. It implements semaphore-based limiting, token bucket algorithms, adaptive rate adjustment, and request throttling -- all accessible through a shared API across Android, iOS, desktop, and server targets. This is essential for managing API call rates and resource contention in multi-platform apps.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, Kotlin Coroutines

## Key Features

- Semaphore-based concurrency limiting
- Token bucket algorithm for smooth rate control
- Adaptive rate limiting with dynamic adjustment
- Request throttling with configurable windows

## Related Repos

- [Concurrency-KMP](../Concurrency-KMP/README.md) — concurrency primitives that rate limiting builds upon
- [Storage-KMP](../Storage-KMP/README.md) — storage services that benefit from rate-limited access

---
*Part of the [vasic-digital catalogue](../README.md)*
