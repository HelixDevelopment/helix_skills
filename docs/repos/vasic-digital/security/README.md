# security

- **GitHub URL**: <https://github.com/vasic-digital/security>
- **Description**: Generic reusable Go module: digital.vasic.security
- **Category**: Auth + Security + Middleware
- **Status**: Active

## Overview

security is a reusable Go module published as `digital.vasic.security` that provides security utilities for Go applications. It covers encryption helpers, secure random generation, input sanitization, and other security primitives that Go services commonly need. The module is designed to be imported alongside auth and middleware for a complete security stack.

## Tech Stack

- Language: Go
- Framework: Go standard library, crypto packages

## Key Features

- Encryption and hashing utility functions
- Secure random generation for tokens and secrets
- Input sanitization and validation helpers
- Reusable Go module with clean import path

## Related Repos

- [auth](../auth/README.md) — authentication module that uses security primitives
- [middleware](../middleware/README.md) — HTTP middleware that applies security utilities
- [Security-KMP](../Security-KMP/README.md) — Kotlin Multiplatform equivalent for cross-platform security

---
*Part of the [vasic-digital catalogue](../README.md)*
