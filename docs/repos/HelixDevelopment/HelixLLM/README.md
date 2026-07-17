# HelixLLM

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixLLM>
- **Description**: Local-running LLM system providing on-premises inference with GPU acceleration, model quantization, and OpenAI-compatible API endpoint
- **Category**: AI / Infrastructure
- **Status**: Active

## Overview

HelixLLM enables running large language models locally for privacy-sensitive workloads, offline operation, and low-latency inference without external API dependencies. It provides an OpenAI-compatible API endpoint for drop-in replacement of cloud LLM services, with GPU acceleration, model quantization, and multi-model serving capabilities. Acts as the local fallback provider in the multi-alias native-first priority system.

## Tech Stack

- Language: Go (server), Python (model integration), C++ (inference kernels)
- Inference: llama.cpp integration, ONNX Runtime, custom inference engine
- Architecture: Model server with REST API and WebSocket streaming
- Key patterns: Model pooling, request batching, KV-cache management

## Key Features

- Local LLM inference with support for multiple model architectures (LLaMA, Mistral, CodeLlama, etc.)
- OpenAI-compatible API endpoint for drop-in replacement of cloud LLM services
- GPU acceleration with automatic device detection (CUDA, ROCm, Metal)
- Model quantization support (GGUF, GPTQ, AWQ) for memory-efficient inference
- Multi-model serving -- run multiple models concurrently with automatic routing

## Related Repos

- [LLMProvider](../LLMProvider/README.md) -- uses HelixLLM as a local provider adapter alongside cloud providers
- [HelixAgent](../HelixAgent/README.md) -- consumes HelixLLM for privacy-sensitive agent operations
- [LLMOrchestrator](../LLMOrchestrator/README.md) -- integrates HelixLLM as a local model provider in the provider chain
- [HelixCode](../HelixCode/README.md) -- uses HelixLLM for offline code generation and review
- [HelixMemory](../HelixMemory/README.md) -- local embedding generation

---
*Part of the [HelixDevelopment catalogue](../README.md)*
