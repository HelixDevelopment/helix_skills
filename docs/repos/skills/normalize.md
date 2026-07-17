# Normalize

> **Repo:** [vasic-digital/Normalize](https://github.com/vasic-digital/Normalize.git)
> **Type:** Declared dependency · **Status:** Active · **Org:** vasic-digital

## Overview

A canonicalisation library that normalises adversarial inputs before
they reach LLM guardrail pipelines. Detects and defuses common prompt
injection techniques, encoding tricks, and obfuscation patterns used
to bypass safety filters.

## Key capabilities

- Adversarial input detection and canonicalisation
- Prompt injection defanging
- Encoding and obfuscation pattern normalisation
- Request canonicalization for consistent input processing

## Architecture

Normalize operates as an input preprocessing pipeline:

1. **Detection layer** — identifies adversarial patterns, encoding
   tricks, and injection attempts
2. **Canonicalisation engine** — normalises detected patterns into
   safe, canonical forms
3. **Pattern registry** — extensible set of known adversarial
   patterns and their canonical replacements
4. **Pass-through mode** — benign inputs pass through unchanged
   with minimal overhead

## Integration points

- **token_optimizer** — direct dependency (request canonicalization);
  normalizes inputs before token optimization
- **RedTeam** — adversarial test fixtures that exercise normalisation
  patterns
- **conversation** — conversation pipeline consuming normalised
  inputs
- **LLMProvider** — provider adapters receiving canonicalised
  requests

## Configuration

Pattern sets, severity thresholds, and pass-through policies are
configurable. Check the repo for Go module API documentation.

## Status

**Active.** Go module dependency consumed via `helix-deps.yaml` by
the token_optimizer submodule.
