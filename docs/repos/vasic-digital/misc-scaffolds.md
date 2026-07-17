# Misc / Scaffolds

Back to [vasic-digital index](./README.md) | [Main index](../README.md)

This group holds repositories that do not fit neatly into the other capability domains -- event notification fan-out, container orchestration infrastructure, and the organization's GitHub Pages site. These are supporting infrastructure and early-stage scaffolds that serve cross-cutting needs across the ecosystem.

| Repo | Description | Status |
|---|---|---|
| [Herald](./Herald/README.md) | Multi-channel event notification system written in Go. Ingests system events (build completions, test results, CI alerts, agent status changes) and reliably fans them out to multiple notification channels (Telegram, Slack, Discord, email, webhook) with per-channel routing rules, retry policies, and delivery confirmation. Ensures every alert reaches the right destination without duplication or loss. | Active |
| [containers](./containers/README.md) | Container orchestration infrastructure providing Podman-based rootless container definitions, compose files, and health-check tooling for ecosystem services. Serves as the shared container runtime layer that other repos reference for building and deploying their containerized workloads. | Active |
| [vasic-digital.github.io](./vasic-digital.github.io/README.md) | GitHub Pages site for the vasic-digital organization. Serves as the public-facing web presence with project documentation, repository indexes, and developer guides built from static site generators. | Active |

**Related skills:** [agentwrapper](../skills/agentwrapper.md)
