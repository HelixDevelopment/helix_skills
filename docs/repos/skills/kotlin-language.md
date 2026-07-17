# Kotlin Language

> **Repo:** [HelixDevelopment/helix_skills — constitution/skills/kotlin-language](https://github.com/HelixDevelopment/helix_skills)
> **Type:** constitution skill · **Domain:** language · **Status:** Draft

## Overview

kotlin.language is a reference skill for the Kotlin programming language.
It covers null safety, coroutine-based structured concurrency, seamless
Java interoperability, and Kotlin's role as Android's preferred language
since Google I/O 2019. This skill requires java.language as a
prerequisite, since Kotlin targets the JVM and interoperates directly
with Java libraries and frameworks.

## Key capabilities

- Null safety model (nullable vs non-nullable types, safe calls,
  Elvis operator, lateinit, and platform types)
- Coroutines and structured concurrency (suspend functions, flows,
  coroutine scopes, exception handling)
- Java interoperability patterns (annotations, SAM conversions,
  platform type mapping)
- Kotlin-specific language features (data classes, sealed classes,
  extension functions, delegated properties, inline functions)
- Kotlin Multiplatform (KMP) reference for shared code across JVM,
  JS, Native, and WASM targets
- Build tooling (Gradle Kotlin DSL, kapt/KSP annotation processing)

## Architecture

kotlin.language is structured as an atomic reference skill:

1. **Language core** — type system, null safety, control flow, and
   functional programming patterns
2. **Coroutines model** — suspend/async primitives, dispatchers,
   flows, and structured concurrency hierarchy
3. **Interop layer** — Java/Kotlin boundary patterns, annotation
   mapping, and mixed-project conventions
4. **Multiplatform** — expect/actual mechanism, target configuration,
   and shared module patterns

## Integration points

- **java.language** — required prerequisite; Kotlin compiles to JVM
  bytecode and calls Java code directly
- **android.overview** — Kotlin is the recommended language for Android
  development; Compose UI is Kotlin-first
- **Gradle ecosystem** — Kotlin DSL is the preferred Gradle scripting
  model; KSP replaces kapt for annotation processing
- **Jetpack Compose** — Compose is a Kotlin-specific declarative UI
  framework; this skill provides the language grounding Compose requires

## Configuration

This is a reference skill — it provides knowledge rather than executable
tooling. Consumers load it to ground agents in Kotlin language concepts,
coroutine patterns, and JVM interop conventions.

## Status

**Draft.** The skill definition and content are under development.
Referenced from the Helix Skills catalog with 1 requirement declared
(java.language) and 2 upstream dependencies.
