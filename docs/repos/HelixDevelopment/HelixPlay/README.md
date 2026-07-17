# HelixPlay

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixPlay>
- **Description**: Ultimate gaming experience platform providing the core gaming infrastructure for the Helix ecosystem. Delivers game engine foundations, rendering pipelines, input management, audio systems, and Stake integration for building immersive gaming experiences.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Game engine core with configurable update/render loop and frame pacing
- Multi-platform rendering pipeline (WebGL, OpenGL, Vulkan abstraction layer)
- Input management system supporting keyboard, mouse, touch, and gamepad
- Audio engine with spatial sound, mixing, and streaming playback
- Scene graph management with hierarchical transforms and culling
- Asset pipeline with hot-reload support for rapid iteration
- Stake integration layer for in-game economy, rewards, and transactions
- Physics integration with configurable simulation stepping
- UI framework for in-game HUD, menus, and overlays

## Technology

- **Language**: Go (engine core), TypeScript (web rendering), C++ (performance-critical paths)
- **Frameworks**: Custom ECS architecture, WebGL/WebGPU for browser, native rendering for desktop
- **Architecture**: Layered engine with pluggable subsystems (render, audio, input, physics)
- **Key patterns**: Entity Component System, command pattern for input, observer pattern for events

## Integration

- Foundation for all Helix games (Helix-Game-1, Stake-Tetris, Helix-Game-Demo)
- Uses VisionEngine for visual testing and automated screenshot validation
- Integrates with HelixAgent for AI-driven gameplay features (NPCs, adaptive difficulty)
- Connects to Stake systems for economy and reward integration
- Consumed by helixqa for game-specific automated testing
- Shared game utilities library used across the Helix game portfolio

## Status

Active development. Core engine with rendering, input, and audio subsystems operational. WebGL rendering stable for browser targets. Stake integration layer functional. Ongoing work on Vulkan backend and advanced rendering features.
