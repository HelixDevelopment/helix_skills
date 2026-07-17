# Auth-KMP

- **GitHub URL**: <https://github.com/vasic-digital/Auth-KMP>
- **Description**: Kotlin Multiplatform OAuth2 authentication: flows, token management, secure storage interface
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

Auth-KMP provides OAuth2 authentication as a Kotlin Multiplatform module, enabling shared auth logic across Android, iOS, desktop, and server targets. It implements standard OAuth2 flows, handles token lifecycle management including refresh, and defines a platform-agnostic interface for secure token storage that each target implements natively.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, OAuth2

## Key Features

- OAuth2 authorization code and client credentials flows
- Token lifecycle management with automatic refresh
- Platform-agnostic secure storage interface
- Shared authentication logic across all KMP targets

## Related Repos

- [Security-KMP](../Security-KMP/README.md) — provides the AES encryption and Keychain/KeyStore integration for secure token storage
- [Config-KMP](../Config-KMP/README.md) — configuration types for network protocol settings used by auth flows
- [auth](../auth/README.md) — Go-native equivalent for server-side authentication

---
*Part of the [vasic-digital catalogue](../README.md)*
