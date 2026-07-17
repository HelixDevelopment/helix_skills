# Enterprise Scalability Documentation — Skill Graph System

**Date:** 2026-07-17
**Status:** DESIGN
**Scope:** Scalability planning for enterprise deployment of the skill graph system

---

## 1. Overview

This document addresses scalability concerns for deploying the HelixKnowledge Skill Graph System at enterprise scale, covering:

1. **Database scalability** — PostgreSQL + pgvector at scale
2. **API scalability** — Horizontal scaling of the REST/MCP surface
3. **Worker scalability** — Background job processing at scale
4. **Storage scalability** — Skill content, evidence, and embeddings
5. **Multi-tenant isolation** — Per-team/per-org skill namespaces

---

## 2. Current Architecture Constraints

### 2.1 Single-Process Bottlenecks

| Component | Constraint | Impact |
|-----------|------------|--------|
| PostgreSQL | Single instance | Read/write contention at high concurrency |
| Worker | Single goroutine pool | Job processing limited by single machine |
| Embedding | Synchronous API calls | Throughput limited by provider rate limits |
| MCP Server | stdio transport | Single-client connection |

### 2.2 Resource Limits

| Resource | Current Config | Enterprise Need |
|----------|----------------|-----------------|
| DB connections | `max_connections = 25` | 100-500 |
| Embedding batch | Single skill at a time | 100+ skills/batch |
| Worker concurrency | 1 goroutine per job type | 10-50 concurrent jobs |
| API rate limiting | Token bucket (G22) | Per-tenant quotas |

---

## 3. Scalability Strategies

### 3.1 Database Scalability

#### 3.1.1 Read Replicas
- Deploy PostgreSQL read replicas for search/query workloads
- Route read operations to replicas, writes to primary
- Implementation: `internal/db/pool.go` — add read/write pool separation

#### 3.1.2 Connection Pooling
- Use PgBouncer for connection pooling
- Increase `max_connections` to 100-500
- Monitor connection usage via Prometheus metrics

#### 3.1.3 Partitioning
- Partition `audit_log` by timestamp (monthly)
- Partition `evidences` by `source_project`
- Implementation: PostgreSQL native partitioning

#### 3.1.4 Indexing Strategy
- HNSW indexes for vector search (already implemented)
- GIN indexes for JSONB metadata queries (already implemented)
- BRIN indexes for timestamp columns (audit_log)

### 3.2 API Scalability

#### 3.2.1 Horizontal Scaling
```
┌─────────────────────────────────────────────────────────┐
│                    Load Balancer                         │
│                    (nginx/HAProxy)                       │
└───────────┬───────────────┬───────────────┬─────────────┘
            │               │               │
            ▼               ▼               ▼
    ┌───────────┐   ┌───────────┐   ┌───────────┐
    │ API       │   │ API       │   │ API       │
    │ Instance 1│   │ Instance 2│   │ Instance 3│
    └─────┬─────┘   └─────┬─────┘   └─────┬─────┘
          │               │               │
          └───────────────┼───────────────┘
                          │
                          ▼
                  ┌───────────────┐
                  │  PostgreSQL   │
                  │  (Primary +   │
                  │   Replicas)   │
                  └───────────────┘
```

- Stateless API instances (no session state in-memory)
- Shared PostgreSQL backend
- Redis for rate limiting and caching (optional)

#### 3.2.2 Caching Strategy
- **Skill content**: Cache in Redis with TTL (5 min)
- **Search results**: Cache popular queries
- **Embeddings**: Cache in PostgreSQL (already stored)
- **Dependency trees**: Cache computed trees with invalidation on write

#### 3.2.3 Rate Limiting
- Per-tenant rate limits (configurable)
- Per-endpoint rate limits (search vs CRUD)
- Implementation: Redis-backed token bucket

### 3.3 Worker Scalability

#### 3.3.1 Job Queue
- Replace in-memory job queue with durable queue (Redis/PostgreSQL)
- Support distributed workers across multiple machines
- Implementation: `internal/worker/runner.go` — add queue backend

#### 3.3.2 Worker Pools
```
┌─────────────────────────────────────────────────────────┐
│                    Job Queue                             │
│                    (Redis/PostgreSQL)                    │
└───────────┬───────────────┬───────────────┬─────────────┘
            │               │               │
            ▼               ▼               ▼
    ┌───────────┐   ┌───────────┐   ┌───────────┐
    │ Worker    │   │ Worker    │   │ Worker    │
    │ Node 1    │   │ Node 2    │   │ Node 3    │
    │           │   │           │   │           │
    │ - expand  │   │ - expand  │   │ - validate│
    │ - validate│   │ - validate│   │ - analyze │
    │ - analyze │   │ - analyze │   │ - sync    │
    └───────────┘   └───────────┘   └───────────┘
```

- Each worker node runs all job types
- Horizontal scaling by adding worker nodes
- Job distribution via queue (round-robin or priority)

#### 3.3.3 Embedding Batching
- Batch embedding requests (100+ skills per API call)
- Async embedding pipeline with progress tracking
- Implementation: `internal/db/embedding.go` — add batch API

### 3.4 Storage Scalability

#### 3.4.1 Content Storage
- Store large skill content in object storage (S3/MinIO)
- Keep metadata and embeddings in PostgreSQL
- Implementation: `content_url` column instead of `content TEXT`

#### 3.4.2 Evidence Storage
- Store code snippets in object storage
- Keep metadata and embeddings in PostgreSQL
- Implementation: `code_snippet_url` column

#### 3.4.3 Backup Strategy
- Daily PostgreSQL backups (pg_dump)
- Weekly full backups
- Point-in-time recovery (WAL archiving)

---

## 4. Multi-Tenant Isolation

### 4.1 Namespace Strategy

```sql
-- Add tenant_id to all tables
ALTER TABLE skills ADD COLUMN tenant_id UUID NOT NULL;
ALTER TABLE skill_dependencies ADD COLUMN tenant_id UUID NOT NULL;
ALTER TABLE evidences ADD COLUMN tenant_id UUID NOT NULL;
-- ... etc

-- Row-level security
ALTER TABLE skills ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON skills
    USING (tenant_id = current_setting('app.tenant_id')::uuid);
```

### 4.2 Tenant Management

| Feature | Implementation |
|---------|----------------|
| Tenant creation | `POST /api/v1/tenants` |
| Tenant deletion | `DELETE /api/v1/tenants/:id` (cascade) |
| Skill isolation | Row-level security per tenant |
| Embedding isolation | Separate HNSW indexes per tenant |
| Worker isolation | Tenant-aware job routing |

### 4.3 Shared vs Isolated Resources

| Resource | Strategy |
|----------|----------|
| PostgreSQL instance | Shared (row-level security) |
| Embedding API keys | Per-tenant (config) |
| Worker nodes | Shared (tenant-aware routing) |
| Object storage | Per-tenant buckets |

---

## 5. Performance Targets

### 5.1 Latency Targets

| Operation | Target | Current |
|-----------|--------|---------|
| `GET /skills/:name` | < 50ms | ~10ms |
| `POST /search` | < 200ms | ~100ms |
| `GET /skills/:name/tree` | < 500ms | ~200ms |
| `POST /skills` (create) | < 1s | ~500ms |

### 5.2 Throughput Targets

| Operation | Target | Current |
|-----------|--------|---------|
| Search queries | 1000 req/s | ~100 req/s |
| Skill creates | 100 req/s | ~10 req/s |
| Embedding generation | 1000 skills/min | ~10 skills/min |
| Worker jobs | 50 jobs/min | ~5 jobs/min |

### 5.3 Scale Targets

| Metric | Target | Current |
|--------|--------|---------|
| Total skills | 100K+ | ~100 |
| Total evidence | 1M+ | ~1000 |
| Concurrent users | 1000+ | ~10 |
| Tenants | 100+ | 1 |

---

## 6. Monitoring & Observability

### 6.1 Metrics (Prometheus)

| Metric | Description |
|--------|-------------|
| `skill_api_requests_total` | Total API requests by endpoint |
| `skill_api_latency_seconds` | Request latency histogram |
| `skill_search_latency_seconds` | Search latency histogram |
| `skill_worker_jobs_total` | Worker jobs by type and status |
| `skill_embedding_latency_seconds` | Embedding generation latency |
| `skill_db_connections_active` | Active DB connections |
| `skill_cache_hits_total` | Cache hit/miss counters |

### 6.2 Logging

- Structured JSON logging (already implemented)
- Request tracing with correlation IDs
- Error tracking with stack traces

### 6.3 Alerting

| Alert | Condition |
|-------|-----------|
| High latency | p99 > 1s for 5 min |
| Error rate | > 1% for 5 min |
| DB connections | > 80% of max for 5 min |
| Worker backlog | > 1000 jobs for 10 min |
| Embedding failures | > 10% for 5 min |

---

## 7. Implementation Phases

### Phase 1: Single-Node Optimization (P0)
1. Connection pooling (PgBouncer)
2. Read replicas for search
3. Caching layer (Redis)
4. Batch embedding

### Phase 2: Horizontal Scaling (P1)
1. Stateless API instances
2. Durable job queue (Redis/PostgreSQL)
3. Load balancer configuration
4. Multi-worker deployment

### Phase 3: Multi-Tenancy (P2)
1. Tenant management API
2. Row-level security
3. Per-tenant configuration
4. Tenant-aware worker routing

### Phase 4: Enterprise Features (P3)
1. SSO/SAML integration
2. Audit logging compliance
3. Data retention policies
4. Disaster recovery

---

## 8. Honest Gaps

1. **Multi-tenancy**: Row-level security adds complexity and performance overhead. Needs benchmarking at scale.
2. **Object storage**: Moving content to S3/MinIO requires significant refactoring of read paths.
3. **Distributed workers**: Job queue durability and exactly-once processing need careful design.
4. **Real-time updates**: WebSocket support for live skill tree updates requires server-sent events or WebSocket proxy.
