# HelixLLM

- **GitHub URL**: <https://github.com/HelixDevelopment/HelixLLM>
- **Description**: Helix LLM -- a local-running super model system that provides on-premises LLM inference capabilities. Enables running large language models locally for privacy-sensitive workloads, offline operation, and low-latency inference without external API dependencies.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Local LLM inference with support for multiple model architectures (LLaMA, Mistral, CodeLlama, etc.)
- Model management -- download, quantize, and serve models from a unified interface
- OpenAI-compatible API endpoint for drop-in replacement of cloud LLM services
- GPU acceleration with automatic device detection (CUDA, ROCm, Metal)
- Model quantization support (GGUF, GPTQ, AWQ) for memory-efficient inference
- Streaming response generation with token-by-token output
- Multi-model serving -- run multiple models concurrently with automatic routing
- Context caching and prompt optimization for reduced inference latency

## Technology

- **Language**: Go (server), Python (model integration), C++ (inference kernels)
- **Frameworks**: llama.cpp integration, ONNX Runtime, custom inference engine
- **Architecture**: Model server with REST API and WebSocket streaming
- **Key patterns**: Model pooling, request batching, KV-cache management

## Integration

- Provides local inference backend for LLMProvider when cloud APIs are unavailable or undesired
- Consumed by HelixAgent for privacy-sensitive agent operations
- Integrates with LLMOrchestrator as a local model provider in the provider chain
- Used by HelixCode for offline code generation and review
- Fallback provider in the multi-alias native-first priority system
- Connects to HelixMemory for local embedding generation

## Status

Active development. Core inference engine supports LLaMA-family models. GPU acceleration operational for CUDA devices. OpenAI-compatible API endpoint stable. Ongoing work on model quantization improvements and additional architecture support.
