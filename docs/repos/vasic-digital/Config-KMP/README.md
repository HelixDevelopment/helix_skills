# Config-KMP

- **GitHub URL**: <https://github.com/vasic-digital/Config-KMP>
- **Description**: Kotlin Multiplatform storage configuration types for network protocols
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

Config-KMP defines Kotlin Multiplatform configuration types for network storage protocols. It provides shared data classes and enums for connection settings, protocol parameters, and storage configuration that are used by other KMP modules to configure network-based storage and database services across platforms.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, kotlinx.serialization

## Key Features

- Cross-platform configuration data classes for network protocols
- Serializable storage and connection settings types
- Shared configuration model for network-based services
- Platform-agnostic protocol parameter definitions

## Related Repos

- [Storage-KMP](../Storage-KMP/README.md) — storage service interfaces that consume Config-KMP types
- [Database-KMP](../Database-KMP/README.md) — database interfaces that use Config-KMP for connection configuration
- [Auth-KMP](../Auth-KMP/README.md) — authentication flows that reference Config-KMP for network settings

---
*Part of the [vasic-digital catalogue](../README.md)*
