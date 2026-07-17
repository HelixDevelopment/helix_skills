# Helix-Game-1

- **GitHub URL**: <https://github.com/HelixDevelopment/Helix-Game-1>
- **Description**: First game project implementing a Stake-based gaming experience -- reference implementation for the Helix gaming platform
- **Category**: Gaming
- **Status**: Active

## Overview

Helix-Game-1 is the flagship game project from HelixDevelopment, demonstrating the full Helix gaming platform stack in action. It implements a complete game loop with player interaction systems, Stake integration for in-game economy, and AI-driven features. Serves as the primary reference implementation for game loop architecture, asset pipelines, and platform integration patterns.

## Tech Stack

- Language: Go / TypeScript (depending on platform target)
- Engine: HelixPlay game engine with ECS architecture
- Rendering: Multi-platform (WebGL, OpenGL, Vulkan abstraction)
- Key patterns: Game loop, state machine, event-driven input handling

## Key Features

- Complete game loop with update/render cycle and frame-rate management
- Player input handling with keyboard, mouse, and gamepad support
- Stake integration for in-game economy and reward systems
- Scene management with transitions and state persistence
- Asset loading pipeline for sprites, audio, and level data

## Related Repos

- [HelixPlay](../HelixPlay/README.md) -- gaming platform infrastructure powering Helix-Game-1
- [Stake-Tetris](../Stake-Tetris/README.md) -- shares game utilities and Stake integration patterns
- [Helix-Game-Demo](../Helix-Game-Demo/README.md) -- lightweight demo showcasing platform fundamentals
- [HelixAgent](../HelixAgent/README.md) -- AI-driven NPC behavior and adaptive difficulty
- [VisionEngine](../VisionEngine/README.md) -- visual testing and screenshot-based validation

---
*Part of the [HelixDevelopment catalogue](../README.md)*
