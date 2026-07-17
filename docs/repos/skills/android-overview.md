# Android Overview

> **Repo:** [HelixDevelopment/helix_skills — constitution/skills/android-overview](https://github.com/HelixDevelopment/helix_skills)
> **Type:** constitution skill · **Domain:** android · **Status:** Draft

## Overview

android.overview is a general reference skill for the Android operating
system and application framework. It covers core Android components
(Activity, Service, BroadcastReceiver, ContentProvider), Jetpack Compose
declarative UI, the Gradle/AGP build pipeline, and application
distribution via APK/AAB. It requires java.language and kotlin.language
as prerequisites, since Android development is grounded in the JVM
ecosystem and Kotlin is the preferred language since 2019.

## Key capabilities

- Android application component lifecycle reference (Activity, Service,
  BroadcastReceiver, ContentProvider)
- Jetpack Compose declarative UI model and state management
- Gradle/AGP build pipeline configuration and dependency management
- APK/AAB packaging, signing, and distribution workflows
- Android manifest configuration and permissions model
- Integration with AOSP platform-level concepts (framework, HAL, system
  services)

## Architecture

android.overview is structured as an atomic skill bundle:

1. **Component reference** — lifecycle, intents, and inter-component
   communication patterns
2. **Compose guide** — declarative UI primitives, theming, navigation,
   and state hoisting
3. **Build pipeline** — Gradle/AGP tasks, variant management, and
   dependency resolution
4. **Distribution** — signing, bundletool, Play Store requirements,
   and sideloading

## Integration points

- **java.language** — required prerequisite; Android's original language
  and runtime foundation (Dalvik/ART bytecode)
- **kotlin.language** — required prerequisite; preferred Android language
  since Google I/O 2019
- **android.studio** — IDE integration for project templates, emulator,
  and profiling tools
- **AOSP platform skills** — platform-level reference skills build on
  the app-layer concepts this skill covers

## Configuration

This is a reference skill — it provides knowledge rather than executable
tooling. Consumers load it to ground agents in Android development
concepts. No runtime configuration is required beyond declaring the
skill in the project's skill manifest.

## Status

**Draft.** The skill definition and content are under development.
Referenced from the Helix Skills catalog with 2 requirements declared
and 2 upstream dependencies (java.language, kotlin.language).
