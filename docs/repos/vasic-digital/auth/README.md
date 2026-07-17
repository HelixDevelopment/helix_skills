# auth

- **GitHub URL**: <https://github.com/vasic-digital/auth>
- **Description**: Generic reusable Go module: digital.vasic.auth
- **Category**: Auth + Security + Middleware
- **Status**: Active

## Overview

auth is a reusable Go module published as `digital.vasic.auth` that provides authentication primitives for Go applications. It offers token validation, session management, and authentication middleware that can be imported into any Go service. The module follows Go module conventions for clean dependency management.

## Tech Stack

- Language: Go
- Framework: Go standard library, net/http

## Key Features

- Token validation and verification for HTTP requests
- Session management with configurable storage backends
- HTTP authentication middleware for Go services
- Reusable Go module with clean import path

## Related Repos

- [security](../security/README.md) — complementary Go module for security utilities beyond authentication
- [middleware](../middleware/README.md) — HTTP middleware module that auth integrates with
- [Auth-KMP](../Auth-KMP/README.md) — Kotlin Multiplatform equivalent for cross-platform auth

---
*Part of the [vasic-digital catalogue](../README.md)*
