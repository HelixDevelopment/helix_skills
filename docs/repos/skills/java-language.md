# Java Language

> **Repo:** [HelixDevelopment/helix_skills — constitution/skills/java-language](https://github.com/HelixDevelopment/helix_skills)
> **Type:** constitution skill · **Domain:** language · **Status:** Draft

## Overview

java.language is a reference skill for the Java programming language.
It covers the JVM and bytecode execution model, language evolution from
Java 1.0 through modern LTS releases (17, 21, 25), garbage collector
strategies (G1, ZGC, Shenandoah), and the role of Java as the
foundation for Android development and the Kotlin ecosystem. This skill
serves as a prerequisite for both android.overview and kotlin.language.

## Key capabilities

- JVM architecture reference (class loading, bytecode verification,
  JIT compilation, memory model)
- Language feature evolution across LTS releases (records, sealed
  classes, pattern matching, virtual threads)
- Garbage collector comparison (G1, ZGC, Shenandoah) and tuning
  fundamentals
- Build tooling reference (Maven, Gradle) and dependency management
- Java module system (JPMS) and classpath model
- Concurrency primitives (java.util.concurrent, Project Loom virtual
  threads)

## Architecture

java.language is structured as an atomic reference skill:

1. **JVM runtime model** — class loading, bytecode format, JIT
   compilation tiers, and the memory model (heap, stack, metaspace)
2. **Language reference** — syntax, type system, generics, lambdas,
   and modern features per LTS release
3. **GC and performance** — collector strategies, tuning flags, and
   profiling approaches
4. **Build ecosystem** — Maven/Gradle conventions, dependency scopes,
   and plugin architecture

## Integration points

- **kotlin.language** — Kotlin compiles to JVM bytecode and interops
  seamlessly with Java; this skill is its prerequisite
- **android.overview** — Android's original runtime (Dalvik, now ART)
  is derived from the JVM; Java remains supported alongside Kotlin
- **Gradle ecosystem** — Java is a primary Gradle language; build
  pipeline skills reference this for project configuration
- **Maven ecosystem** — Maven POM conventions and lifecycle phases
  are part of the Java build reference

## Configuration

This is a reference skill — it provides knowledge rather than executable
tooling. Consumers load it to ground agents in Java language concepts,
JVM internals, and build tooling conventions.

## Status

**Draft.** The skill definition and content are under development.
Referenced from the Helix Skills catalog with 0 requirements declared
and 2 downstream dependents (android.overview, kotlin.language).
