# Auth + Security + Middleware

Back to [vasic-digital index](./README.md) | [Main index](../README.md)

This group covers authentication, authorization, security hardening, and HTTP middleware for both application-layer and LLM-layer defence. The active Go modules (auth, middleware, security) provide the standard request-pipeline primitives used by every backend service, while the scaffold repos explore LLM-specific attack surfaces -- prompt-leak detection, system-prompt extraction defence, and recursive safety patterns for agentic systems.

| Repo | Description | Status |
|---|---|---|
| [Claritas](./Claritas/README.md) | System-prompt extraction detection library. Monitors LLM outputs for patterns that indicate the system prompt is being leaked or extracted through adversarial prompting, and triggers configurable countermeasures (refusal, redaction, alerting). | SCAFFOLD / WIP |
| [GandalfSolutions](./GandalfSolutions/README.md) | Read-only solutions archive for prompt-leak-defense testing. Contains known attack vectors and their expected outcomes, used as a reference corpus for validating that guardrail pipelines correctly block prompt-leak attempts. | SCAFFOLD / WIP |
| [LeakHub](./LeakHub/README.md) | Prompt-leak corpus and defensive-use fixtures for red-team training. Curates real-world prompt-injection and prompt-leak examples as YAML fixtures that feed into adversarial testing harnesses like RedTeam. | SCAFFOLD / WIP |
| [Ouroborous](./Ouroborous/README.md) | Recursive and self-referential safety patterns for agentic systems. Explores defence mechanisms against agents that attempt to modify their own guardrails, escalate privileges, or bypass safety constraints through self-referential reasoning loops. | SCAFFOLD / WIP |
| [Veritas](./Veritas/README.md) | Truth and verification auxiliary library. Provides fact-checking and claim-verification primitives that can be integrated into LLM pipelines to cross-reference generated outputs against authoritative knowledge sources. | SCAFFOLD / WIP |
| [auth](./auth/README.md) | Generic reusable Go module (`digital.vasic.auth`) providing authentication and authorization primitives -- JWT validation, OAuth2 flows, API key management, role-based access control, and session token lifecycle. Used as the standard auth layer by all backend services. | Active |
| [middleware](./middleware/README.md) | Reusable Go module (`digital.vasic.middleware`) providing HTTP middleware for request logging, rate limiting, CORS, request-ID propagation, panic recovery, and metrics collection. Composes into standard `net/http` handler chains. | Active |
| [security](./security/README.md) | Generic reusable Go module (`digital.vasic.security`) providing security utilities -- encryption helpers, secret management, input sanitization, CSRF protection, and secure header configuration. Used alongside auth and middleware for defense-in-depth. | Active |

**Related skills:** [normalize](../skills/normalize.md)
