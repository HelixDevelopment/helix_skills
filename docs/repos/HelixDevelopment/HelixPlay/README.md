# HelixPlay

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixPlay>
- **Description**: Gaming platform providing engine foundations, rendering pipelines, input management, audio systems, and Stake integration for building immersive games
- **Category**: Gaming
- **Status**: Active

## Overview

HelixPlay is the core gaming infrastructure for the Helix ecosystem. It delivers a configurable game engine with multi-platform rendering (WebGL, OpenGL, Vulkan abstraction), input management, spatial audio, scene graph, and Stake integration. All Helix games (Helix-Game-1, Stake-Tetris, Helix-Game-Demo) are built on HelixPlay as their shared platform layer.

## Tech Stack

- Language: Go (engine core), TypeScript (web rendering), C++ (performance-critical paths)
- Rendering: WebGL/WebGPU for browser, native rendering for desktop
- Architecture: Layered engine with pluggable subsystems (render, audio, input, physics)
- Key patterns: Entity Component System, command pattern for input, observer pattern for events

## Key Features

- Game engine core with configurable update/render loop and frame pacing
- Multi-platform rendering pipeline (WebGL, OpenGL, Vulkan abstraction layer)
- Input management system supporting keyboard, mouse, touch, and gamepad
- Audio engine with spatial sound, mixing, and streaming playback
- Stake integration layer for in-game economy, rewards, and transactions

## Related Repos

- [Helix-Game-1](../Helix-Game-1/README.md) -- flagship game built on HelixPlay
- [Stake-Tetris](../Stake-Tetris/README.md) -- puzzle game using HelixPlay infrastructure
- [Helix-Game-Demo](../Helix-Game-Demo/README.md) -- demo showcasing HelixPlay fundamentals
- [VisionEngine](../VisionEngine/README.md) -- visual testing and automated screenshot validation
- [HelixAgent](../HelixAgent/README.md) -- AI-driven gameplay features (NPCs, adaptive difficulty)

---
*Part of the [HelixDevelopment catalogue](../README.md)*
