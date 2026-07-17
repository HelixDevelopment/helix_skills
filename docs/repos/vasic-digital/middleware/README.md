# middleware

- **GitHub URL**: <https://github.com/vasic-digital/middleware>
- **Description**: digital.vasic.middleware - Reusable Go module
- **Category**: Auth + Security + Middleware
- **Status**: Active

## Overview

middleware is a reusable Go module published as `digital.vasic.middleware` that provides HTTP middleware components for Go web services. It includes request logging, error handling, CORS, rate limiting, and other cross-cutting concerns that are composed into HTTP handler chains. The module is designed for clean integration with standard Go HTTP servers.

## Tech Stack

- Language: Go
- Framework: Go standard library, net/http

## Key Features

- Composable HTTP middleware chain for Go services
- Request logging, error handling, and recovery middleware
- CORS and request validation middleware
- Clean integration with Go net/http handler patterns

## Related Repos

- [auth](../auth/README.md) — authentication middleware that plugs into the middleware chain
- [security](../security/README.md) — security utilities used by middleware components

---
*Part of the [vasic-digital catalogue](../README.md)*
