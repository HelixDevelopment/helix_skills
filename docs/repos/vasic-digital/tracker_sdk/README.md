# tracker_sdk

- **GitHub URL**: <https://github.com/vasic-digital/tracker_sdk>
- **Description**: Generic, tracker-agnostic SDK primitives: mirror manager, plugin registry, testing harness. Used by the Lava project (https://github.com/milos85vasic/Lava) to support multiple torrent trackers.
- **Category**: Messaging + Observability + Storage
- **Status**: Active

## Overview

A generic SDK providing tracker-agnostic primitives for torrent tracker integrations. Includes a mirror manager for failover, a plugin registry for adding new tracker backends, and a testing harness for validation. Used by the Lava project to support multiple torrent tracker protocols.

## Tech Stack

- Language: Go
- Module path: digital.vasic.tracker_sdk
- Used by: [Lava](https://github.com/milos85vasic/Lava)

## Key Features

- Mirror manager for tracker endpoint failover
- Plugin registry for adding tracker backends
- Testing harness for tracker integration validation

## Related Repos

- [Plugins](../Plugins/README.md) — general plugin system
- [ShareConnect](../ShareConnect/README.md) — URL sharing to processing endpoints
- [TransmissionConnect](../TransmissionConnect/README.md) — Transmission client integration

---
*Part of the [vasic-digital catalogue](../README.md)*
