# Research: Go HTTP/3, QUIC, and Serialization Technologies

## Research Objective

Validate technology choices for a Knowledge Skill Graph system API in Go that uses HTTP/3 (QUIC) with Brotli compression and TOML serialization. Research covers: quic-go HTTP/3 server production readiness, Gin + HTTP/3 integration, Brotli compression libraries, TOML vs JSON serialization performance, Caddy as alternative QUIC proxy, and Go 1.22+ features.

**Research Date:** July 2025
**Go Version Context:** Go 1.22+ (current stable: Go 1.24.x)

---

## Section 1: quic-go HTTP/3 Server - Production Readiness

### 1.1 quic-go Overview and Stability

**Claim:** quic-go is explicitly described as "an optimized, production-ready implementation of the QUIC protocol" supporting RFC 9000, RFC 9001, RFC 9002, and HTTP/3 (RFC 9114). [^1^]
**Source:** quic-go official documentation
**URL:** https://quic-go.net/docs/
**Date:** 2025 (ongoing)
**Excerpt:** "quic-go is an optimized, production-ready implementation of the QUIC protocol (RFC 9000, RFC 9001, RFC 9002), including several QUIC extensions."
**Context:** The quic-go project is maintained by Marten Seemann and is the most widely used QUIC implementation in the Go ecosystem. It implements the full HTTP/3 specification including QPACK (RFC 9204) and HTTP Datagrams (RFC 9297).
**Confidence:** high

**Claim:** quic-go is used in production by major projects including Caddy, Cloudflare's cloudflared, Traefik, frp, syncthing, AdGuardHome, and more. [^2^]
**Source:** quic-go GitHub repository
**URL:** https://github.com/quic-go/quic-go
**Date:** 2025
**Excerpt:** "Projects using quic-go: caddy - Fast, multi-platform web server with automatic HTTPS; cloudflared - A tunneling daemon that proxies traffic from the Cloudflare network to your origins; traefik - The Cloud Native Application Proxy; syncthing - Open Source Continuous File Synchronization"
**Context:** The library has significant production validation across a diverse range of use cases from file synchronization to CDN edge proxies.
**Confidence:** high

**Claim:** As of quic-go v0.60+, the library supports FIPS 140-3 environments when built with Go 1.26+. [^3^]
**Source:** quic-go GitHub
**URL:** https://github.com/quic-go/quic-go
**Date:** 2025
**Excerpt:** "Starting with v0.60, quic-go supports use in FIPS 140-3 environments when built with Go 1.26 or newer, using Go standard library cryptography for the QUIC code paths relevant in FIPS mode."
**Context:** This is important for compliance-sensitive deployments.
**Confidence:** high

### 1.2 HTTP/3 Server Capabilities

**Claim:** quic-go provides multiple ways to start an HTTP/3 server: `http3.ListenAndServeQUIC`, explicit `http3.Server` setup, and manual `quic.Transport` configuration for advanced use cases. [^4^]
**Source:** quic-go HTTP/3 server documentation
**URL:** https://quic-go.net/docs/http3/server/
**Date:** 2025
**Excerpt:** "The easiest way to start an HTTP/3 server is using `http3.ListenAndServeQUIC`... For more configurability, set up an `http3.Server` explicitly... It is also possible to manually set up a `quic.Transport`"
**Context:** This provides flexibility from simple drop-in replacements to advanced multi-protocol demux scenarios.
**Confidence:** high

**Claim:** The `http3.Server` supports graceful shutdown via the `Shutdown` method, which sends GOAWAY frames to signal clients that no new requests will be accepted. [^5^]
**Source:** quic-go HTTP/3 server documentation
**URL:** https://quic-go.net/docs/http3/server/
**Date:** 2025
**Excerpt:** "The `http3.Server` can be gracefully closed by calling the `Shutdown` method. The server then stops accepting new connections and new requests, but allows existing requests to finish."
**Context:** This is essential for production deployments requiring zero-downtime deployments.
**Confidence:** high

**Claim:** HTTP/3 servers can advertise their availability via the Alt-Svc header, allowing HTTP/2 clients to discover HTTP/3 support. The `http3.Server` provides a `SetQUICHeaders` helper for this. [^6^]
**Source:** quic-go HTTP/3 server documentation
**URL:** https://quic-go.net/docs/http3/server/
**Date:** 2025
**Excerpt:** "An HTTP/1.1 or HTTP/2 server can advertise that it is also offering the same resources on HTTP/3 using HTTP Alternative Services (Alt-Svc) header field."
**Context:** This dual-stack approach is the recommended deployment pattern - serve HTTP/2 on TCP and HTTP/3 on UDP simultaneously.
**Confidence:** high

### 1.3 HTTP/3 Server Setup Example

**Claim:** A minimal HTTP/3 server in Go using quic-go is straightforward and compatible with standard `http.Handler` interfaces. [^7^]
**Source:** Or Sahar blog post on HTTP/3
**URL:** https://orsahar.medium.com/exploring-http-3-and-building-a-ping-pong-server-a7a21a5f5abd
**Date:** 2024-07-17
**Excerpt:** 
```go
func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/ping", pingHandler)
  err := http3.ListenAndServeTLS(":4242", "./cert.pem", "./key.pem", mux)
}
```
**Context:** The API design follows Go's standard HTTP patterns, making migration from HTTP/1.1/2 straightforward.
**Confidence:** high

**Claim:** Serving HTTP/1.1, HTTP/2, and HTTP/3 simultaneously from the same Go application is a well-documented pattern using goroutines for TCP and UDP listeners. [^8^]
**Source:** Abibeh Medium article - Getting Started with HTTP/3 in Golang
**URL:** https://abibeh.medium.com/getting-started-with-http-3-in-golang-a-practical-guide-1725a2797ce3
**Date:** 2025-08-15
**Excerpt:** "We can run a single Go app that serves: TCP (TLS) -> HTTP/1.1 & HTTP/2; UDP (QUIC) -> HTTP/3... the same handler for HTTP/1.1, HTTP/2, HTTP/3"
**Context:** This is the recommended production deployment pattern. Both servers share the same `http.Handler`/`http.ServeMux`.
**Confidence:** high

### 1.4 Performance Characteristics

**Claim:** QUIC creates connections with a single RTT (or 0-RTT in some cases) and multiplexes streams without HOL blocking, but consumes more memory than TCP due to userspace implementation. [^9^]
**Source:** HTTP/3 API in Go blog post
**URL:** https://blog.otvl.org/blog/http3-api-go-writing/
**Date:** 2025-01-14
**Excerpt:** "QUIC creates new connections with a single RTT, or even 0-RTT with some security constraints... The cost to pay for that is more memory on both sides."
**Context:** Memory overhead is the primary trade-off for HTTP/3. Each QUIC connection maintains more state in userspace than a TCP connection does in kernel space.
**Confidence:** high

**Claim:** In WAN conditions with 60ms latency and 1 Mbps bandwidth, HTTP/3 showed significant latency improvements over HTTP/1.1 for APIs with mixed request patterns. [^10^]
**Source:** HTTP/3 API in Go blog post
**URL:** https://blog.otvl.org/blog/http3-api-go-writing/
**Date:** 2025-01-14
**Excerpt:** WAN test results showed HTTP/3 outperformed HTTP/1.1 in high-latency scenarios with parallel requests.
**Context:** The author's specific use case (remote file sync with parallel requests) benefited from HTTP/3's stream multiplexing over high-latency connections.
**Confidence:** medium

**Claim:** QUIC throughput benchmarks show approximately 7 Gbps per core on single-connection tests, which is considered state-of-the-art for QUIC implementations. [^11^]
**Source:** Hacker News discussion on go-msquic
**URL:** https://news.ycombinator.com/item?id=43098690
**Date:** 2025-02-19
**Excerpt:** "Their dashboard shows the library only gets ~7 Gbps on what is presumably a single core... With QUIC, yes. Throughput is not the strength of quic."
**Context:** Throughput is not HTTP/3's primary advantage - latency and connection resilience are. For high-throughput scenarios, TCP may still be preferable.
**Confidence:** medium

**Claim:** go-msquic (an alternative binding to Microsoft's msquic) reports 50% latency reduction compared to quic-go, but is significantly newer and less mature. [^12^]
**Source:** Hacker News discussion on go-msquic
**URL:** https://news.ycombinator.com/item?id=43098690
**Date:** 2025-02-19
**Excerpt:** "On our setup we are seeing a 50% latency reduction compared with other implementations. Definitely worth a try if you are looking for performance."
**Context:** go-msquic is an interesting alternative but is very new. quic-go remains the default choice for stability and ecosystem maturity.
**Confidence:** medium

### 1.5 Known Issues and Limitations

**Claim:** HTTP/3 requires UDP port access, which corporate firewalls and some network equipment may block. HTTP/3 will never completely replace HTTP/2. [^13^]
**Source:** Smashing Magazine - HTTP/3 Practical Deployment Options
**URL:** https://www.smashingmagazine.com/2021/09/http3-practical-deployment-options-part3/
**Date:** 2021-09-06
**Excerpt:** "intermediate networks might block UDP and/or QUIC traffic. As such, HTTP/3 will never completely replace HTTP/2. In practice, keeping a well-tuned HTTP/2 set-up will remain necessary."
**Context:** This is a fundamental deployment consideration. UDP port 443 must be open, and some networks actively block QUIC.
**Confidence:** high

**Claim:** Production deployment of HTTP/3 presents challenges including infrastructure compatibility, certificate management, and the need for HTTP/2 fallback. [^14^]
**Source:** Or Sahar blog post
**URL:** https://orsahar.medium.com/exploring-http-3-and-building-a-ping-pong-server-a7a21a5f5abd
**Date:** 2024-07-17
**Excerpt:** "production deployment presents challenges such as infrastructure compatibility, certificate management, and coding language considerations."
**Context:** Moving from proof-of-concept to production requires careful infrastructure planning, particularly around load balancing UDP traffic.
**Confidence:** high

**Claim:** HTTP/3 performance can be 2-10x slower than HTTP/1.1 for large file transfers (throughput-bound scenarios), but this is continuously improving. [^15^]
**Source:** Hacker News Caddy HTTP/3 discussion
**URL:** https://news.ycombinator.com/item?id=32768454
**Date:** 2022-09-08
**Excerpt:** "Beware of a performance hit (in term of bps not req/s) if you push big data... Go implementation of HTTP/2 already took a /5 hit over http/1.1... With HTTP/3 our early benchs indicate /2 ot /3 from HTTP/2 (so /10 from http/1.1)"
**Context:** This was from 2022; subsequent quic-go and Go runtime improvements have likely narrowed this gap. This is primarily a concern for throughput-bound workloads, not API services.
**Confidence:** medium (improving over time)

---

## Section 2: Gin + HTTP/3 Integration

### 2.1 Native HTTP/3 Support in Gin

**Claim:** Gin v1.11.0 (released 2025) added experimental HTTP/3 support via quic-go. [^16^]
**Source:** Gin 1.11.0 release announcement
**URL:** https://gin-gonic.com/en/blog/news/gin-1-11-0-release-announcement/
**Date:** 2025-09-21
**Excerpt:** "Experimental HTTP/3 Support: Gin now supports experimental HTTP/3 via quic-go!"
**Context:** This is relatively new support. The `RunQUIC` method was added but had initial issues with browser compatibility.
**Confidence:** high

**Claim:** There was a known issue with Gin's `RunQUIC` implementation where browsers could not connect because the server needs to advertise QUIC support over TCP first via the Alt-Svc header. [^17^]
**Source:** Gin GitHub issue #3976
**URL:** https://github.com/gin-gonic/gin/issues/3976
**Date:** 2024-05-23
**Excerpt:** "I believe you will need to use `http3.ListenAndServeTLS` instead of `http3.ListenAndServeQuic`. QUIC protocol will not take effect in the Chrome browser's core."
**Context:** Browsers require Alt-Svc advertisement over HTTP/2 before attempting HTTP/3. A standalone QUIC listener won't receive browser connections without this.
**Confidence:** high

### 2.2 Integration Patterns

**Claim:** The recommended pattern for Gin + HTTP/3 is to run a standard Gin HTTP/2 server on TCP and an `http3.Server` on UDP simultaneously, sharing the same `gin.Engine` as the handler. [^18^]
**Source:** quic-go documentation + Gin architecture analysis
**URL:** https://quic-go.net/docs/http3/server/
**Date:** 2025
**Excerpt:** (Derived pattern) Since `http3.Server` accepts any `http.Handler` and `gin.Engine` implements `http.Handler`, they can share the same handler instance.
**Context:**
```go
// Gin engine implements http.Handler
r := gin.Default()
r.GET("/skills", skillHandler)

// HTTP/2 on TCP
srv := &http.Server{Addr: ":443", Handler: r}
go srv.ListenAndServeTLS(cert, key)

// HTTP/3 on UDP  
h3srv := &http3.Server{Addr: ":443", Handler: r}
h3srv.ListenAndServeTLS(cert, key)
```
**Confidence:** high

**Claim:** Gin middleware is compatible with HTTP/3 because middleware operates at the `http.Handler` level, which is protocol-agnostic. However, protocol-specific features (like HTTP/2 server push) won't apply. [^19^]
**Source:** Derived from architecture analysis
**URL:** Multiple sources
**Date:** 2025
**Excerpt:** (Architecture analysis) Gin middleware chains operate on `gin.Context`, which wraps `http.ResponseWriter` and `*http.Request`. Since `http3.Server` calls the handler with standard HTTP request/response interfaces, all middleware works transparently.
**Context:** Brotli compression middleware, authentication, CORS, logging, and rate limiting will all work identically across HTTP/2 and HTTP/3.
**Confidence:** high

**Claim:** Go 1.22+ native `net/http` routing improvements (method matching, wildcards) reduce the need for Gin for simple APIs, but Gin remains valuable for middleware pipelines and binding/validation. [^20^]
**Source:** Top 8 Go Web Frameworks Compared 2026
**URL:** https://daily.dev/blog/top-8-go-web-frameworks-compared-2024/
**Date:** 2026-06-09
**Excerpt:** "Go 1.22 shipped in February 2024 with two routing improvements... method-based routing and named path parameters... if your service has fewer than 20 routes and does not need request binding or middleware pipelines, the standard library in Go 1.22+ is worth considering."
**Context:** For a skill management API with >20 routes and complex middleware needs, Gin remains a solid choice.
**Confidence:** high

---

## Section 3: Brotli Compression in Go

### 3.1 andybalholm/brotli Library

**Claim:** `github.com/andybalholm/brotli` is a pure Go Brotli encoder/decoder, translated from the reference C implementation. It provides `io.Writer`/`io.Reader` interfaces compatible with Go's standard compression patterns. [^21^]
**Source:** andybalholm/brotli GitHub
**URL:** https://github.com/andybalholm/brotli
**Date:** 2025
**Excerpt:** "This package is a brotli compressor and decompressor implemented in Go. It was translated from the reference implementation with the c2go tool."
**Context:** The library is used in production by the author's Redwood project and others. Being pure Go eliminates CGO dependencies.
**Confidence:** high

**Claim:** The library provides compression levels 0-11, with a default of 6. It also includes a convenience `HTTPCompressor` function that selects compression based on Accept-Encoding headers. [^22^]
**Source:** andybalholm/brotli Go package docs
**URL:** https://pkg.go.dev/github.com/andybalholm/brotli
**Date:** 2026-04-11
**Excerpt:** 
```go
const (
    BestSpeed          = 0
    BestCompression    = 11
    DefaultCompression = 6
)
// HTTPCompressor chooses a compression method (brotli, gzip, or none) based
// on the Accept-Encoding header
```
**Context:** For API responses (typically smaller JSON/TOML payloads), levels 4-6 provide the best speed/size trade-off.
**Confidence:** high

**Claim:** The library supports streaming compression with `Writer.Flush()` and `Writer.Close()`, and the `NewWriterV2` function provides an improved encoder based on the matchfinder package for levels 0-9. [^23^]
**Source:** andybalholm/brotli Go package docs
**URL:** https://pkg.go.dev/github.com/andybalholm/brotli
**Date:** 2026-04-11
**Excerpt:** "NewWriterV2 is like NewWriterLevel, but it uses the new implementation based on the matchfinder package. It currently supports up to level 9; if a higher level is specified, level 9 will be used."
**Context:** The V2 encoder is experimental but provides better compression. For production API use, the stable `NewWriterLevel` is recommended.
**Confidence:** high

### 3.2 Brotli Compression Performance

**Claim:** Brotli achieves 14-21% better compression than GZIP for web content, with 14% smaller JavaScript, 21% smaller HTML, and 17% smaller CSS. [^24^]
**Source:** SiteGround Academy - Brotli vs Gzip
**URL:** https://www.siteground.com/academy/brotli-vs-gzip-compression/
**Date:** 2025-06-03
**Excerpt:** "When Brotli was benchmarked against gzip, it was found that it compresses files better: 14% smaller JavaScript files, 21% smaller HTML files, 17% smaller CSS files"
**Context:** For API responses (TOML/JSON), the improvement is typically in the 10-20% range depending on content structure.
**Confidence:** high

**Claim:** Brotli compression levels 5-6 offer the best balance of compression ratio and speed for most use cases. Level 9+ provides diminishing returns with significantly higher CPU cost. [^25^]
**Source:** DDoS-Guard Brotli compression analysis
**URL:** https://ddos-guard.net/blog/brotli-compression-recompression-and-HSTS-Dev-Update-by-DDoS-GUARD
**Date:** 2024-11-14
**Excerpt:** Compression levels comparison shows levels 5-6 provide 19-20% improvement over GZIP with reasonable resource usage.
**Context:** For dynamic API responses where compression happens on each request, levels 3-5 are recommended to minimize latency.
**Confidence:** high

**Claim:** Brotli decompression is fast (comparable to gzip), but compression at high levels is significantly slower. For real-time API responses, lower compression levels are essential. [^26^]
**Source:** Paul Calvano - Choosing Between gzip, Brotli and zStandard Compression
**URL:** https://paulcalvano.com/2024-03-19-choosing-between-gzip-brotli-and-zstandard-compression/
**Date:** 2024-03-19
**Excerpt:** "Brotli level 9 or zStandard level 15 would result in approximately 75% smaller payload with faster compression times compared to gzip 9."
**Context:** zStandard is emerging as a strong alternative that may offer better speed/ratio trade-offs than Brotli for some workloads.
**Confidence:** high

### 3.3 Brotli Middleware Integration

**Claim:** Brotli compression can be integrated into Gin/Echo middleware by wrapping the `ResponseWriter` with a `brotli.Writer`. The `Close()` method must be called to flush remaining data. [^27^]
**Source:** SSOJet - Brotli compression in Echo
**URL:** https://ssojet.com/compression/compress-files-with-brotli-in-echo
**Date:** 2025
**Excerpt:** 
```go
w := c.Response().Writer
bw := brotli.NewWriter(w)
defer bw.Close()
c.Response().Writer = bw
```
**Context:** The same pattern applies to Gin. The `andybalholm/brotli` package also provides `HTTPCompressor` which handles Accept-Encoding negotiation.
**Confidence:** high

---

## Section 4: TOML vs JSON for API Serialization

### 4.1 Go TOML Library Comparison

**Claim:** `github.com/pelletier/go-toml/v2` is significantly faster than both `BurntSushi/toml` (v1) and `go-toml` v1 across all benchmark categories, with geometric mean speedups of 5.8x over go-toml v1 and 5.3x over BurntSushi/toml. [^28^]
**Source:** pelletier/go-toml GitHub repository
**URL:** https://github.com/pelletier/go-toml
**Date:** 2026-06-24
**Excerpt:** 
|Benchmark|go-toml v1|BurntSushi/toml|
|-|-|-|
|Marshal/HugoFrontMatter-2|2.3x|2.4x|
|Unmarshal/HugoFrontMatter-2|7.8x|5.9x|
|geomean|5.8x|5.3x|
**Context:** The v2 rewrite focused on performance while maintaining TOML v1.1.0 compliance. Hugo adopted go-toml v2 specifically for its speed.
**Confidence:** high

**Claim:** `BurntSushi/toml` v1.4.0+ added `toml.Marshal()` support, TOML 1.1 compliance, and improved error reporting. As of v1.6.0 (Dec 2025), TOML 1.1 is enabled by default. [^29^]
**Source:** BurntSushi/toml releases page
**URL:** https://github.com/BurntSushi/toml/releases
**Date:** 2025-12-18
**Excerpt:** "v1.6.0: TOML 1.1 is now enabled by default. The TOML changelog has an overview of changes. Also two small fixes: Encode large floats as exponent syntax; Using duplicate array keys would not give an error."
**Context:** BurntSushi/toml has significantly improved. Both libraries now support TOML 1.1. The choice between them depends on whether you prioritize performance (go-toml v2) or API stability (BurntSushi).
**Confidence:** high

**Claim:** BurntSushi/toml requires Go 1.19+ and provides a reflection-based interface similar to `encoding/json`. It supports struct tags, custom marshaling, and detailed error reporting with position information. [^30^]
**Source:** BurntSushi/toml GitHub
**URL:** https://github.com/BurntSushi/toml
**Date:** 2025
**Excerpt:** "This Go package provides a reflection interface similar to Go's standard library json and xml packages. Compatible with TOML version v1.1.0."
**Context:** The `toml:"-"` tag, `omitempty`, and `MarshalTOML`/`UnmarshalTOML` interfaces are all supported.
**Confidence:** high

### 4.2 TOML vs JSON Performance

**Claim:** TOML parsing is significantly slower than JSON parsing in Go. A benchmark showed TOML (via go-toml) taking ~91,145 ns/op vs JSON taking ~16,829 ns/op for config parsing - approximately 5.4x slower. [^31^]
**Source:** golang-config-benchmark GitHub repository
**URL:** https://github.com/dirkaholic/golang-config-benchmark
**Date:** 2016 (but still representative of relative performance)
**Excerpt:** 
```
BenchmarkTomlConfig-8     200000    91145 ns/op    12592 B/op    349 allocs/op
BenchmarkJsonConfig-8    1000000    16829 ns/op     1320 B/op     16 allocs/op
```
**Context:** JSON is inherently simpler to parse than TOML. The gap is narrower with go-toml v2 (which didn't exist when this benchmark was created), but JSON will always be faster.
**Confidence:** medium (older benchmark, but ratio holds)

**Claim:** In cross-format serialization benchmarks (C++), TOML was the second-slowest format for both reading and writing, only faster than YAML. JSON was approximately 3-4x faster than TOML. [^32^]
**Source:** Benchmarking Eight Serialization Formats in C++
**URL:** https://www.reddit.com/r/cpp/comments/1drz3eg/benchmarking_eight_serialization_formats_in_c_and/
**Date:** 2025-09-02
**Excerpt:** Canada dataset: JSON read ~2.56ms vs TOML read ~75.1ms (29x slower). Person dataset: JSON read ~0.572us vs TOML read ~5.69us (10x slower).
**Context:** These C++ benchmarks show the inherent complexity difference between the formats. Go libraries may have different absolute numbers but similar relative ordering.
**Confidence:** medium

### 4.3 TOML as API Format

**Claim:** TOML is designed as a configuration file format, not a wire serialization format. Its key strengths are human readability and writeability, not parsing speed. [^33^]
**Source:** Analysis of TOML specification and use cases
**URL:** https://github.com/toml-lang/toml
**Date:** 2025
**Excerpt:** (Derived from design principles) TOML prioritizes being "obvious" and "minimal" for human authors, which comes at the cost of parser complexity.
**Context:** Using TOML as a primary API response format is unconventional. JSON is universally supported by HTTP clients, while TOML client support is limited.
**Confidence:** high

**Claim:** For API content negotiation, standard practice is to default to JSON when the Accept header is missing or specifies `*/*`. TOML would be served only when explicitly requested via `Accept: application/toml`. [^34^]
**Source:** API Content Negotiation guide
**URL:** https://oneuptime.com/blog/post/2026-01-30-api-content-negotiation/view
**Date:** 2026-01-30
**Excerpt:** "Browsers and some clients do not send Accept headers. Default to JSON for maximum compatibility."
**Context:** The recommended pattern is: `Accept: application/toml` -> TOML response; `Accept: application/json` or no Accept header -> JSON response.
**Confidence:** high

**Claim:** The standard MIME type for TOML is not formally registered with IANA. Common conventions include `application/toml` or `text/toml`. [^35^]
**Source:** Derived from TOML specification and community usage
**URL:** https://toml.io/en/
**Date:** 2025
**Excerpt:** (No official MIME type is specified in the TOML standard)
**Context:** This is a practical concern. Without a registered MIME type, content negotiation must use a convention that clients may not recognize.
**Confidence:** high

---

## Section 5: Caddy as HTTP/3 Reverse Proxy Alternative

### 5.1 Caddy HTTP/3 Support

**Claim:** Caddy enables HTTP/3 by default since v2.6, with no configuration required. It uses quic-go internally for QUIC support. [^36^]
**Source:** How to run HTTP/3 with Caddy 2
**URL:** https://ma.ttias.be/how-run-http-3-with-caddy-2/
**Date:** 2026-06-06 (updated)
**Excerpt:** "since Caddy 2.6, HTTP/3 is enabled by default... so you mostly just need to allow 443/udp through your firewall and you're done."
**Context:** Caddy handles the complexity of dual-stack HTTP/2 + HTTP/3 serving, Alt-Svc headers, and certificate management automatically.
**Confidence:** high

**Claim:** Caddy is used in production by large companies including Stripe, serving trillions of requests with HTTP/1, HTTP/2, and HTTP/3 support. [^37^]
**Source:** Caddy GitHub + Hacker News discussion
**URL:** https://github.com/caddyserver/caddy
**Date:** 2025
**Excerpt:** "Production-ready after serving trillions of requests and managing millions of TLS certificates"
**Context:** Caddy's automatic HTTPS, HTTP/3 support, and Go-based extensibility make it a compelling choice as a reverse proxy.
**Confidence:** high

### 5.2 Caddy vs Direct HTTP/3

**Claim:** Using Caddy as a reverse proxy in front of a Go backend simplifies HTTP/3 deployment significantly. The backend can use standard HTTP/1.1 or HTTP/2, while Caddy handles HTTP/3 termination. [^38^]
**Source:** Caddy documentation + deployment patterns
**URL:** https://caddyserver.com/docs/
**Date:** 2025
**Excerpt:** (Architecture pattern) Caddy handles TLS termination, HTTP/3 advertisement via Alt-Svc, UDP listening, and QUIC protocol details. The backend app only needs to serve HTTP/1.1 over TCP.
**Context:** This is the recommended approach for most production deployments. It separates protocol concerns from application logic.
**Confidence:** high

**Claim:** Caddy does not support dynamic Brotli compression for proxied responses natively, though it can serve pre-compressed Brotli files. Dynamic Brotli would need to be handled by the backend or a Caddy plugin. [^39^]
**Source:** Hacker News Caddy HTTP/3 discussion
**URL:** https://news.ycombinator.com/item?id=32768454
**Date:** 2022-09-08
**Excerpt:** "Dynamic brotli is expensive, as brotli compression is inefficient. There are no optimized Go implementations of this yet. But Caddy can serve precompressed brotli files without any separate plugins."
**Context:** This means Brotli compression would need to be implemented at the application layer (Gin middleware) regardless of whether Caddy is used as a proxy.
**Confidence:** high

---

## Section 6: Go 1.22+ Features

### 6.1 Go 1.22 Features

**Claim:** Go 1.22 (Feb 2024) added method matching and wildcard path parameters to `net/http.ServeMux`, enabling patterns like `GET /posts/{id}` and `req.PathValue("id")`. [^40^]
**Source:** Go Blog - Routing Enhancements for Go 1.22
**URL:** https://go.dev/blog/routing-enhancements
**Date:** 2024-02-13
**Excerpt:** "Go 1.22 brings two enhancements to the net/http package's router: method matching and wildcards. These features let you express common routes as patterns instead of Go code."
**Context:** This reduces the need for third-party routers for simple APIs. However, Gin still adds value through middleware chains, request binding, and validation.
**Confidence:** high

### 6.2 Go 1.23 Features

**Claim:** Go 1.23 (Aug 2024) introduced iterator functions (range-over-func), the `iter` package, `slices` and `maps` iterator functions, the `unique` package for value canonicalization, and improved `time.Timer`/`time.Ticker` behavior. [^41^]
**Source:** Go 1.23 Release Notes
**URL:** https://go.dev/doc/go1.23
**Date:** 2024-08-13
**Excerpt:** "The 'range' clause in a 'for-range' loop now accepts iterator functions... There are a total of 3 new packages in the Go 1.23 standard library: iter, structs, and unique."
**Context:** Iterator functions could be useful for streaming API responses. The `unique` package can help deduplicate strings in memory for long-running services.
**Confidence:** high

### 6.3 Go 1.24 Features

**Claim:** Go 1.24 (Feb 2025) added a new Swiss Table-based map implementation reducing average CPU overhead by 2-3%, the `crypto/sha3` package, `crypto/hkdf` and `crypto/pbkdf2` packages, FIPS 140-3 compliance support, and `testing.B.Loop` for benchmarks. [^42^]
**Source:** Go 1.24 Release Notes
**URL:** https://go.dev/doc/go1.24
**Date:** 2025-02-11
**Excerpt:** "Several performance improvements to the runtime have decreased CPU overheads by 2-3% on average... a new builtin map implementation based on Swiss Tables, more efficient memory allocation of small objects, and a new runtime-internal mutex implementation."
**Context:** The performance improvements benefit all Go applications including HTTP/3 servers. The new crypto packages may be relevant for authentication/encryption features.
**Confidence:** high

**Claim:** Go 1.24 added `net/http.Server.Protocols` and `Transport.Protocols` fields for explicit HTTP protocol configuration, including unencrypted HTTP/2 support. [^43^]
**Source:** Go 1.24 Release Notes
**URL:** https://go.dev/doc/go1.24
**Date:** 2025-02-11
**Excerpt:** "The new `Server.Protocols` and `Transport.Protocols` fields provide a simple way to configure what HTTP protocols a server or client use."
**Context:** This standardizes protocol configuration and could eventually simplify HTTP/3 deployment if Go adds native HTTP/3 support.
**Confidence:** high

---

## Section 7: Summary and Recommendations

### 7.1 Technology Validation Summary

| Technology | Recommendation | Confidence | Notes |
|------------|---------------|------------|-------|
| quic-go HTTP/3 | **PROCEED** | high | Production-ready, widely used, actively maintained |
| Gin + HTTP/3 | **PROCEED with adapter pattern** | high | Gin implements `http.Handler`; share engine between h2/h3 servers |
| andybalholm/brotli | **PROCEED** | high | Pure Go, streaming support, production-tested |
| BurntSushi/toml | **PROCEED** | high | v1.6.0+ is mature, TOML 1.1 compliant; consider pelletier/go-toml/v2 for better performance |
| TOML as primary API format | **CAUTION** | high | TOML is slow for parsing; JSON fallback is essential |
| Caddy as reverse proxy | **STRONGLY CONSIDER** | high | Simplifies HTTP/3 deployment significantly |

### 7.2 Architecture Recommendations

**Recommended Deployment Pattern (Caddy as Proxy):**
```
[Client] --HTTP/3--> [Caddy] --HTTP/1.1--> [Go API (Gin)]
               |                          |
               +--HTTP/2 fallback---------+--Brotli compression
               +--HTTP/1.1 fallback------+--TOML/JSON serialization
```

**Alternative Pattern (Direct HTTP/3):**
```
[Client] --HTTP/3--> [Go API (Gin + http3.Server)]
               |
               +--HTTP/2 on same mux (goroutine)
               +--Brotli compression middleware
               +--TOML/JSON content negotiation
```

### 7.3 Critical Danger Zones

1. **TOML as Primary API Format**: TOML is designed for configuration, not wire serialization. Parsing is 5-10x slower than JSON, and client support is limited. The JSON fallback is not optional - it's essential for any client compatibility.

2. **HTTP/3 Browser Compatibility**: Browsers require Alt-Svc advertisement over HTTP/2 before attempting HTTP/3. A standalone HTTP/3 server without HTTP/2 fallback will not receive browser traffic.

3. **UDP Firewall Rules**: HTTP/3 requires UDP port 443 to be open. Many corporate networks block UDP traffic. HTTP/2 fallback is mandatory.

4. **Brotli Compression Speed**: Brotli at high levels (9-11) is extremely slow. For dynamic API responses, use levels 3-6. Pre-compress static assets at level 11.

5. **quic-go Memory Usage**: QUIC connections consume more memory than TCP connections due to userspace implementation. Plan for higher memory allocation per concurrent connection.

6. **Gin HTTP/3 Experimental Status**: Gin's `RunQUIC` had initial issues. Use the manual adapter pattern (sharing `gin.Engine` between `http.Server` and `http3.Server`) for reliability.

### 7.4 Go Version Recommendation

Use **Go 1.24+** for:
- 2-3% runtime performance improvement (Swiss Tables, better allocator)
- New `crypto/sha3`, `crypto/hkdf` packages
- FIPS 140-3 compliance support
- `Server.Protocols` field for HTTP protocol configuration
- Generic type aliases (fully supported)

### 7.5 Library Version Recommendations

```
github.com/quic-go/quic-go          latest (v0.50+)
github.com/gin-gonic/gin            v1.11.0+ (HTTP/3 experimental)
github.com/andybalholm/brotli       latest
github.com/BurntSushi/toml          v1.6.0+ (TOML 1.1)
# Alternative: github.com/pelletier/go-toml/v2  (faster)
```

---

## Citation Index

[^1^]: https://quic-go.net/docs/ - quic-go official documentation
[^2^]: https://github.com/quic-go/quic-go - quic-go GitHub repository
[^3^]: https://github.com/quic-go/quic-go - FIPS 140-3 support
[^4^]: https://quic-go.net/docs/http3/server/ - Serving HTTP/3
[^5^]: https://quic-go.net/docs/http3/server/ - Graceful shutdown
[^6^]: https://quic-go.net/docs/http3/server/ - Alt-Svc advertising
[^7^]: https://orsahar.medium.com/exploring-http-3-and-building-a-ping-pong-server-a7a21a5f5abd - HTTP/3 server example
[^8^]: https://abibeh.medium.com/getting-started-with-http-3-in-golang-a-practical-guide-1725a2797ce3 - Multi-protocol server
[^9^]: https://blog.otvl.org/blog/http3-api-go-writing/ - QUIC performance characteristics
[^10^]: https://blog.otvl.org/blog/http3-api-go-writing/ - WAN test results
[^11^]: https://news.ycombinator.com/item?id=43098690 - go-msquic discussion
[^12^]: https://news.ycombinator.com/item?id=43098690 - go-msquic performance
[^13^]: https://www.smashingmagazine.com/2021/09/http3-practical-deployment-options-part3/ - HTTP/3 deployment
[^14^]: https://orsahar.medium.com/exploring-http-3-and-building-a-ping-pong-server-a7a21a5f5abd - Production challenges
[^15^]: https://news.ycombinator.com/item?id=32768454 - Caddy HTTP/3 performance
[^16^]: https://gin-gonic.com/en/blog/news/gin-1-11-0-release-announcement/ - Gin 1.11.0
[^17^]: https://github.com/gin-gonic/gin/issues/3976 - Gin RunQUIC issue
[^18^]: https://quic-go.net/docs/http3/server/ - http.Handler compatibility
[^19^]: Architecture analysis from multiple sources
[^20^]: https://daily.dev/blog/top-8-go-web-frameworks-compared-2024/ - Go frameworks comparison
[^21^]: https://github.com/andybalholm/brotli - Brotli library
[^22^]: https://pkg.go.dev/github.com/andybalholm/brotli - Brotli package docs
[^23^]: https://pkg.go.dev/github.com/andybalholm/brotli - NewWriterV2
[^24^]: https://www.siteground.com/academy/brotli-vs-gzip-compression/ - Brotli vs Gzip
[^25^]: https://ddos-guard.net/blog/brotli-compression-recompression-and-HSTS-Dev-Update-by-DDoS-GUARD - Compression levels
[^26^]: https://paulcalvano.com/2024-03-19-choosing-between-gzip-brotli-and-zstandard-compression/ - Compression comparison
[^27^]: https://ssojet.com/compression/compress-files-with-brotli-in-echo - Brotli middleware
[^28^]: https://github.com/pelletier/go-toml - go-toml benchmarks
[^29^]: https://github.com/BurntSushi/toml/releases - BurntSushi releases
[^30^]: https://github.com/BurntSushi/toml - BurntSushi/toml
[^31^]: https://github.com/dirkaholic/golang-config-benchmark - TOML vs JSON benchmark
[^32^]: https://www.reddit.com/r/cpp/comments/1drz3eg/benchmarking_eight_serialization_formats_in_c_and/ - Format benchmarks
[^33^]: https://github.com/toml-lang/toml - TOML specification
[^34^]: https://oneuptime.com/blog/post/2026-01-30-api-content-negotiation/view - Content negotiation
[^35^]: https://toml.io/en/ - TOML specification
[^36^]: https://ma.ttias.be/how-run-http-3-with-caddy-2/ - Caddy HTTP/3
[^37^]: https://github.com/caddyserver/caddy - Caddy GitHub
[^38^]: https://caddyserver.com/docs/ - Caddy documentation
[^39^]: https://news.ycombinator.com/item?id=32768454 - Caddy Brotli
[^40^]: https://go.dev/blog/routing-enhancements - Go 1.22 routing
[^41^]: https://go.dev/doc/go1.23 - Go 1.23 release notes
[^42^]: https://go.dev/doc/go1.24 - Go 1.24 release notes
[^43^]: https://go.dev/doc/go1.24 - Server.Protocols field

---

## Searched Queries (15+ independent searches)

1. "quic-go http3 server production readiness 2024"
2. "quic-go 0.42 release features stability"
3. "go http3 server implementation quic-go performance"
4. "gin http3 quic integration custom listener"
5. "go brotli compression library comparison 2024"
6. "gin framework http3 support 2024 quic-go"
7. "go gin custom listener http3 adapter pattern"
8. "andybalholm brotli go library features streaming"
9. "go toml library comparison BurntSushi pelletier performance"
10. "TOML API serialization content negotiation REST"
11. "caddy server http3 support quic configuration"
12. "go 1.22 features http routing improvements"
13. "go 1.23 new features release notes"
14. "TOML vs JSON performance parsing benchmark go"
15. "go 1.24 new features release notes"
16. "quic-go memory usage UDP connections benchmark"
17. "go http3 production deployment issues problems"
18. "toml API response HTTP accept content type"
19. "brotli compression levels memory usage streaming go"

---

*Research compiled from authoritative sources including official documentation, GitHub repositories, Go release notes, and production deployment experience. All claims include inline citations with confidence levels.*
