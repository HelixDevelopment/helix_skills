# http3

- **GitHub URL**: <https://github.com/vasic-digital/http3>
- **Description**: Generic Go module wrapping quic-go/http3 for net/http.Handler servers -- drop-in HTTP/3 support
- **Category**: Messaging + Observability + Storage
- **Status**: Active

## Overview

A Go module that wraps quic-go/http3 to provide drop-in HTTP/3 support for standard net/http.Handler servers. Enables existing HTTP servers to upgrade to HTTP/3 (QUIC) with minimal code changes.

## Tech Stack

- Language: Go
- Module path: digital.vasic.http3
- Built on: quic-go/http3
- Protocol: HTTP/3 (QUIC)

## Key Features

- Drop-in HTTP/3 support for net/http.Handler
- QUIC protocol transport layer
- Backward-compatible with HTTP/1.1 and HTTP/2

## Related Repos

- [Proxy](../Proxy/README.md) — proxy server with protocol support
- [ratelimiter](../ratelimiter/README.md) — HTTP request rate limiting
- [observability](../observability/README.md) — HTTP request monitoring

---
*Part of the [vasic-digital catalogue](../README.md)*
