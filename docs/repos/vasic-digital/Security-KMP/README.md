# Security-KMP

- **GitHub URL**: <https://github.com/vasic-digital/Security-KMP>
- **Description**: Kotlin Multiplatform secure storage: AES encryption, platform Keychain/KeyStore integration
- **Category**: Cross-platform (KMP) Modules
- **Status**: Active

## Overview

Security-KMP provides cross-platform secure storage for Kotlin Multiplatform applications. It implements AES encryption for data at rest and defines interfaces for platform-native secure storage -- iOS Keychain on Apple targets and Android KeyStore on Android. This enables sensitive data like tokens and credentials to be stored securely on every platform.

## Tech Stack

- Language: Kotlin (KMP)
- Framework: Kotlin Multiplatform, AES cryptography, platform Keychain/KeyStore APIs

## Key Features

- AES encryption for data at rest across all platforms
- iOS Keychain integration for Apple targets
- Android KeyStore integration for Android targets
- Platform-agnostic secure storage interface

## Related Repos

- [Auth-KMP](../Auth-KMP/README.md) — OAuth2 module that uses Security-KMP for secure token storage
- [security](../security/README.md) — Go-native equivalent for server-side security utilities

---
*Part of the [vasic-digital catalogue](../README.md)*
