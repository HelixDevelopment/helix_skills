# Linux OS

> **Repo:** [HelixDevelopment/helix_skills — constitution/skills/linux-os](https://github.com/HelixDevelopment/helix_skills)
> **Type:** constitution skill · **Domain:** os · **Status:** Draft

## Overview

linux.os is a core reference skill for the Linux operating system. It
covers the kernel architecture (process scheduler, memory management,
VFS, device model), major distribution families (Debian/Ubuntu,
RHEL/Fedora, Arch, Alpine), filesystem hierarchy (FHS), process model
(fork/exec, signals, namespaces, cgroups), and essential tooling (systemd,
shell, coreutils, package managers). This skill serves as foundational
knowledge for any project targeting Linux hosts, containers, or embedded
systems.

## Key capabilities

- Kernel architecture reference (scheduler, mm, VFS, net, device model,
  loadable modules)
- Process model (fork/exec, signals, IPC, namespaces, cgroups,
  seccomp)
- Filesystem hierarchy (FHS layout, mount model, tmpfs, procfs,
  sysfs)
- Distribution families and package management (apt, dnf, pacman,
  apk)
- systemd unit model (service, timer, socket, mount, path units)
- Shell and coreutils reference (POSIX sh, bash, common toolchain)
- Container primitives (namespaces, cgroups, overlayfs, rootless
  runtimes)

## Architecture

linux.os is structured as an atomic reference skill:

1. **Kernel model** — subsystem overviews (scheduler, mm, VFS, net,
   char/block device model) and loadable module management
2. **Process and IPC** — lifecycle, signals, pipes, sockets, shared
   memory, and systemd's process supervision
3. **Filesystem and storage** — FHS layout, mount semantics, LVM,
   RAID, and filesystem types (ext4, xfs, btrfs, tmpfs)
4. **Systemd reference** — unit types, dependency ordering, journal,
   and cgroup delegation
5. **Container layer** — namespace/cgroup isolation, overlayfs
   stacking, and rootless runtime model (podman, crun)

## Integration points

- **Android/AOSP** — Android's kernel is Linux; platform-level skills
  reference this for kernel and driver concepts
- **Container runtimes** — rootless container mandates (constitution
  §11.4.161) build on Linux namespace/cgroup primitives
- **Embedded and IoT** — Linux-based embedded targets reference this
  skill for kernel configuration and device-tree concepts
- **Server and cloud** — deployment, systemd service management, and
  network configuration all ground in this skill

## Configuration

This is a reference skill — it provides knowledge rather than executable
tooling. Consumers load it to ground agents in Linux OS concepts, kernel
internals, and system administration conventions.

## Status

**Draft.** The skill definition and content are under development.
Referenced from the Helix Skills catalog with 0 requirements declared
and 2 upstream dependencies.
