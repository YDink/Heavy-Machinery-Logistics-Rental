# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Architecture Specification — 10. Deployment & Infrastructure Architecture (`docs/shared/architecture/10-deployment-architecture.md`)

**Document Identifier:** `ARCH-DEPLOY-DINK-CORE-V7.0`  
**Classification:** Canonical Architecture Specification & Engineering Reference (Single Source of Truth)  
**Governing Documents:** `SRS-DINK-CORE-V7.0` (Parts 1–5), `ARCH-INTRO-DINK-CORE-V7.0`, `ARCH-PRIN-DINK-CORE-V7.0`, `ARCH-OVERVIEW-DINK-CORE-V7.0`, `ARCH-BACKEND-GO-DINK-CORE-V7.0`, `ARCH-DB-DINK-CORE-V7.0`, `ARCH-COMM-DINK-CORE-V7.0`, `ARCH-SEC-DINK-CORE-V7.0`, `contracts/openapi.yaml`, `ADR-001` through `ADR-014`  
**Target Runtimes:** Linux AMD64 / ARM64, Docker Engine 27+, Traefik 3.x, Go 1.27.3 (`distroless/static-debian12`), Next.js 16 Standalone (Node.js 22 LTS), PostgreSQL 18.4 + PostGIS 3.4, Valkey 8.0-alpine, SeaweedFS S3, ClamAV Daemon, OpenBao KMS  
**Target Hosting Environments:** Sovereign On-Premise / Colocation Datacenters (Tier III Facilities, Addis Ababa) & Regional Edge (Nairobi)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, ISO/IEC 27001:2022, Ethiopian Personal Data Protection Proclamation No. 1321/2024 (Local Data Residency), NBE Directives  
**Status:** Approved Engineering Baseline

---

## Table of Contents

- 1. [Executive Summary & Infrastructure Tenets](#1-executive-summary--infrastructure-tenets)
- 2. [Multi-Environment Architecture Strategy](#2-multi-environment-architecture-strategy)
  - 2.1 [Environment Matrix & Isolation Boundaries](#21-environment-matrix--isolation-boundaries)
  - 2.2 [Macro Infrastructure Topology](#22-macro-infrastructure-topology)
- 3. [Production Network Architecture & Edge Ingress](#3-production-network-architecture--edge-ingress)
  - 3.1 [VPC Network Segmentation & Firewall Zones](#31-vpc-network-segmentation--firewall-zones)
  - 3.2 [Traefik 3.x Edge Gateway Configuration](#32-traefik-3x-edge-gateway-configuration)
  - 3.3 [Ingress Port Mapping Matrix](#33-ingress-port-mapping-matrix)
- 4. [Containerization & Workload Specifications](#4-containerization--workload-specifications)
  - 4.1 [Go Modular Monolith Container (`distroless`)](#41-go-modular-monolith-container-distroless)
  - 4.2 [Next.js 16 Admin Web Console Container](#42-nextjs-16-admin-web-console-container)
  - 4.3 [Production Compose Orchestration Baseline](#43-production-compose-orchestration-baseline)
- 5. [High Availability, Clustering & Storage Topography](#5-high-availability-clustering--storage-topography)
  - 5.1 [PostgreSQL 18.4 Streaming WAL Replication](#51-postgresql-184-streaming-wal-replication)
  - 5.2 [Valkey 8.0 Sentinel Failover Topology](#52-valkey-80-sentinel-failover-topology)
  - 5.3 [SeaweedFS Distributed S3 Storage Replication](#53-seaweedfs-distributed-s3-storage-replication)
- 6. [Local Data Residency & Regional Hosting Topology](#6-local-data-residency--regional-hosting-topology)
  - 6.1 [Statutory Data Residency Governance (Proclamation No. 1321/2024)](#61-statutory-data-residency-governance-proclamation-no-13212024)
  - 6.2 [Cross-Border Edge Peering (Addis Ababa Core & Nairobi PoP)](#62-cross-border-edge-peering-addis-ababa-core--nairobi-pop)
- 7. [CI/CD Pipeline & Automated Release Engineering](#7-cicd-pipeline--automated-release-engineering)
  - 7.1 [Automated Quality & Invariant Gates](#71-automated-quality--invariant-gates)
  - 7.2 [Zero-Downtime Blue-Green Deployment Lifecycle](#72-zero-downtime-blue-green-deployment-lifecycle)
- 8. [Disaster Recovery, Backup & Point-in-Time Recovery (PITR)](#8-disaster-recovery-backup--point-in-time-recovery-pitr)
  - 8.1 [Recovery Objectives (RPO & RTO)](#81-recovery-objectives-rpo--rto)
  - 8.2 [Automated pgBackRest Cadence](#82-automated-pgbackrest-cadence)
  - 8.3 [Offsite Cold Storage Encryption](#83-offsite-cold-storage-encryption)
- 9. [Observability, Telemetry & SRE Stack](#9-observability-telemetry--sre-stack)
  - 9.1 [Prometheus Metrics Harvesting](#91-prometheus-metrics-harvesting)
  - 9.2 [Distributed Tracing with OpenTelemetry](#92-distributed-tracing-with-opentelemetry)
  - 9.3 [Zero-PII Log Aggregation with Promtail & Loki](#93-zero-pii-log-aggregation-with-promtail--loki)
  - 9.4 [Health Check Probes & Readiness Probes](#94-health-check-probes--readiness-probes)
- 10. [Secrets Provisioning with OpenBao KMS](#10-secrets-provisioning-with-openbao-kms)
- 11. [Operational Runbooks & Incident Playbooks](#11-operational-runbooks--incident-playbooks)

---

## 1. Executive Summary & Infrastructure Tenets

The Dink Machinery infrastructure layer provides an enterprise operational environment for heavy machinery discovery, logistics coordination, direct bank transfer verification, and live telematics across East Africa.

```mermaid
flowchart TD
    subgraph Tenets [Infrastructure Architecture Tenets]
        T1["1. Sovereign Data Residency - Local Tier III Facility"]
        T2["2. Minimal Operational Complexity - Single Binary Monolith"]
        T3["3. Self-Hosted S3 Storage - On-Premise SeaweedFS Cluster"]
        T4["4. Zero-Downtime Releases - Blue-Green Container Cutover"]
    end

    subgraph HostPlatform [Production Hosting Topology]
        PhysicalTier["Tier III Colocation Facility - Addis Ababa"]
        ContainerTier["Stateless Docker Engine Workloads"]
        EncryptedStorage["NVMe Storage Arrays with LUKS2 Encryption"]
    end

    Tenets --> HostPlatform
```

### The Four Infrastructure Tenets

1. **National Data Sovereignty & Local Residency:** Primary transactional databases, national identity verifications, citizen phone numbers, and raw GPS tracks reside strictly within Ethiopian sovereign territory (Tier III certified facility in Addis Ababa). Foreign cloud dependencies for transactional persistence are strictly prohibited under Ethiopian Personal Data Protection Proclamation No. 1321/2024.
2. **Minimal Operational Complexity:** Microservice sprawl on Kubernetes is rejected in favor of a Go 1.27.3 modular monolith deployed as static, minimal Linux containers via orchestrated Docker Engine / Swarm. This delivers low operational cognitive overhead while sustaining 10,000+ requests per second on modest bare-metal hardware.
3. **Self-Hosted, Air-Gapped S3 Storage:** Inspection photographs, damage condition records, and payment bank slips are persisted to an on-premise, horizontally scalable **SeaweedFS S3** cluster, eliminating external cloud egress bandwidth charges and ensuring local internet exchange (EthIO-IX) latency.
4. **Deterministic Zero-Downtime Releases:** Deployments, database schema migrations via Goose, and runtime dynamic UI reconfigurations execute without dropping in-flight transactions or active telematics streaming sessions.

---

## 2. Multi-Environment Architecture Strategy

### 2.1 Environment Matrix & Isolation Boundaries

| Attribute | Development (`dev`) | Staging (`stage`) | Production (`prod`) |
|---|---|---|---|
| **Primary Goal** | Local feature iteration | End-to-end integration & load benchmarks | Live customer operations |
| **Hosting Target** | Developer Workstation | Dedicated Staging Host (Addis Ababa) | Colocated Tier III DC (Addis Ababa Core + Nairobi PoP) |
| **Compute Profile** | 4 Cores, 16GB RAM | 16 Cores, 64GB RAM | Clustered Active-Active (64 Cores, 256GB RAM) |
| **Data Ingestion** | Synthetic Seed Fixtures | Anonymized Production Mirror | Live Citizen Data & Production Telco Rails |
| **Database Mode** | Single PostgreSQL 18.4 Container | Replicated PostgreSQL Container | Physical Primary + Standby Streaming WAL Cluster |
| **Object Store** | Local SeaweedFS Standalone | Replicated SeaweedFS Master/Volume | Multi-Node Replicated SeaweedFS Cluster |

### 2.2 Macro Infrastructure Topology

```mermaid
flowchart TD
    subgraph Clients [Client Workstations and Edge Personas]
        MobileApp["Flutter Mobile App - iOS and Android"]
        WebAdmin["Next.js 16 Admin Console - Browsers"]
    end

    subgraph EdgeDMZ [Edge Ingress DMZ Network]
        TraefikGateway["Traefik 3.x Ingress Proxy<br/>TLS 1.3 Termination and WAF"]
    end

    subgraph AppNetwork [Application Workload Network]
        MonolithNode1["Go 1.27.3 Monolith Node 01"]
        MonolithNode2["Go 1.27.3 Monolith Node 02"]
        AdminWebNode1["Next.js 16 Web Console Node 01"]
        AdminWebNode2["Next.js 16 Web Console Node 02"]
        ClamAVDaemon["ClamAV Antivirus Daemon"]
    end

    subgraph DataNetwork [Isolated Persistence Network]
        PGPrimary[("PostgreSQL 18.4 Primary<br/>Read/Write on NVMe")]
        PGStandby[("PostgreSQL 18.4 Standby<br/>Read Replica")]
        ValkeyCluster[("Valkey 8.0 Cache Cluster<br/>Sentinel Failover")]
        SeaweedCluster[("SeaweedFS S3 Storage<br/>Replicated Volumes")]
        OpenBaoKMS["OpenBao Key Management System"]
    end

    MobileApp -->|HTTPS REST and WSS| TraefikGateway
    WebAdmin -->|HTTPS REST| TraefikGateway

    TraefikGateway --> MonolithNode1
    TraefikGateway --> MonolithNode2
    TraefikGateway --> AdminWebNode1
    TraefikGateway --> AdminWebNode2

    MonolithNode1 --> PGPrimary
    MonolithNode2 --> PGPrimary
    MonolithNode1 -.->|Read Splitting| PGStandby
    MonolithNode2 -.->|Read Splitting| PGStandby

    MonolithNode1 --> ValkeyCluster
    MonolithNode2 --> ValkeyCluster
    MonolithNode1 --> SeaweedCluster
    MonolithNode2 --> SeaweedCluster
    MonolithNode1 --> ClamAVDaemon
    MonolithNode2 --> ClamAVDaemon
    MonolithNode1 --> OpenBaoKMS
    MonolithNode2 --> OpenBaoKMS

    PGPrimary -->|Streaming WAL Replication| PGStandby
```

---

## 3. Production Network Architecture & Edge Ingress

### 3.1 VPC Network Segmentation & Firewall Zones

1. **Zone 1: DMZ Ingress (`172.28.10.0/24`)**
   - Houses the Traefik 3.x edge proxies.
   - Publicly exposes only TCP ports `80` (HTTP to HTTPS redirection) and `443` (HTTPS / WSS).
2. **Zone 2: Application Workloads (`172.28.20.0/24`)**
   - Houses the Go monolith instances and Next.js admin frontend instances.
   - Isolated from public routing; receives traffic exclusively from Zone 1.
3. **Zone 3: Secure Persistence & Storage (`172.28.30.0/24`)**
   - Houses PostgreSQL, Valkey, SeaweedFS, ClamAV, and OpenBao.
   - No direct internet connectivity. Ingress permitted exclusively from Zone 2 on authorized application ports.
4. **Zone 4: Management & Telemetry (`172.28.40.0/24`)**
   - Houses Prometheus, Grafana, Promtail/Loki, and OpenTelemetry collectors.
   - Accessible only via authenticated internal WireGuard VPN tunnels.

### 3.2 Traefik 3.x Edge Gateway Configuration

```yaml
# /etc/traefik/traefik.yaml - Production Static Configuration
global:
  checkNewVersion: false
  sendAnonymousUsage: false

entryPoints:
  web:
    address: ":80"
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https
          permanent: true

  websecure:
    address: ":443"
    http:
      tls:
        options: strictTLS@file
      middlewares:
        - securityHeaders@file
        - globalRateLimit@file

providers:
  file:
    directory: "/etc/traefik/dynamic"
    watch: true

api:
  dashboard: false
```

```yaml
# /etc/traefik/dynamic/routing.yaml - Production Dynamic Routing Rules
http:
  routers:
    core-api:
      rule: "Host(`api.dinkmachinery.com`)"
      entryPoints: ["websecure"]
      service: core-api-service
      middlewares:
        - corsHeaders
        - compression

    admin-console:
      rule: "Host(`admin.dinkmachinery.com`)"
      entryPoints: ["websecure"]
      service: admin-console-service
      middlewares:
        - adminSecurityHeaders

    storage-s3:
      rule: "Host(`s3.dinkmachinery.com`)"
      entryPoints: ["websecure"]
      service: storage-s3-service
      middlewares:
        - s3CorsHeaders

  services:
    core-api-service:
      loadBalancer:
        healthCheck:
          path: /healthz
          interval: "10s"
          timeout: "3s"
        servers:
          - url: "[http://172.28.20.10:8080](http://172.28.20.10:8080)"
          - url: "[http://172.28.20.11:8080](http://172.28.20.11:8080)"

    admin-console-service:
      loadBalancer:
        healthCheck:
          path: /api/health
          interval: "10s"
          timeout: "3s"
        servers:
          - url: "[http://172.28.20.20:3000](http://172.28.20.20:3000)"
          - url: "[http://172.28.20.21:3000](http://172.28.20.21:3000)"

    storage-s3-service:
      loadBalancer:
        servers:
          - url: "[http://172.28.30.50:8333](http://172.28.30.50:8333)"

  middlewares:
    securityHeaders:
      headers:
        sslRedirect: true
        stsSeconds: 63072000
        stsIncludeSubdomains: true
        stsPreload: true
        contentTypeNosniff: true
        browserXssFilter: true
        customFrameOptionsValue: "DENY"
        customResponseHeaders:
          X-Powered-By: "Dink-Core-V7"

    globalRateLimit:
      rateLimit:
        average: 1000
        burst: 200

    corsHeaders:
      headers:
        accessControlAllowMethods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
        accessControlAllowHeaders: ["Authorization", "Content-Type", "Idempotency-Key", "X-Country-Code", "X-Correlation-ID", "Accept-Language"]
        accessControlAllowOriginList: ["[https://dinkmachinery.com](https://dinkmachinery.com)", "[https://admin.dinkmachinery.com](https://admin.dinkmachinery.com)"]
        accessControlMaxAge: 3600

tls:
  options:
    strictTLS:
      minVersion: "VersionTLS13"
      cipherSuites:
        - "TLS_AES_128_GCM_SHA256"
        - "TLS_AES_256_GCM_SHA384"
        - "TLS_CHACHA20_POLY1305_SHA256"
      curvePreferences:
        - "X25519"
        - "SecP256r1"
      sniStrict: true
```

### 3.3 Ingress Port Mapping Matrix

| Segment | Service | Port | Protocol | Scope |
|---|---|---|---|---|
| **Public Edge** | Traefik Reverse Proxy | `80`, `443` | HTTP/2, TLS 1.3 | Public Internet (`0.0.0.0/0`) |
| **Zone 1 -> 2** | Go Modular Monolith | `8080` | HTTP REST / WSS | App Network (`172.28.20.10:8080`) |
| **Zone 1 -> 2** | Next.js Admin Console | `3000` | HTTP/1.1 | App Network (`172.28.20.20:3000`) |
| **Zone 2 -> 3** | PostgreSQL Primary | `5432` | Extended Binary Protocol | Persistence Network (`172.28.30.10:5432`) |
| **Zone 2 -> 3** | Valkey Cache Cluster | `6379` | RESP3 TCP Protocol | Persistence Network (`172.28.30.20:6379`) |
| **Zone 2 -> 3** | ClamAV Antivirus | `3310` | TCP (`zINSTREAM`) | Persistence Network (`172.28.30.40:3310`) |
| **Zone 2 -> 3** | SeaweedFS S3 API | `8333` | HTTP S3 REST | Persistence Network (`172.28.30.50:8333`) |
| **Zone 2 -> 3** | OpenBao KMS Server | `8200` | HTTPS (mTLS) | Persistence Network (`172.28.30.60:8200`) |
| **Zone 4 (VPN)** | Prometheus Telemetry | `9090` | HTTP Web UI | Management Subnet (`10.99.0.5:9090`) |
| **Zone 4 (VPN)** | Grafana Dashboards | `3001` | HTTPS Web UI | Management Subnet (`10.99.0.6:3001`) |

---

## 4. Containerization & Workload Specifications

### 4.1 Go Modular Monolith Container (`distroless`)

The backend builds as a single static binary (`CGO_ENABLED=0`) deployed in Google's `distroless/static-debian12` container without shell access or auxiliary utilities.

```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.27.3-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.version=v7.0.0 -X main.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    -o /bin/dink-server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

LABEL maintainer="Dink Machinery Engineering <infra@dinkmachinery.com>"
LABEL security.asvs="level-3"

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=nonroot:nonroot /bin/dink-server /app/dink-server

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/dink-server", "-healthcheck"]

ENTRYPOINT ["/app/dink-server"]
```

### 4.2 Next.js 16 Admin Web Console Container

```dockerfile
# syntax=docker/dockerfile:1.7
FROM node:22-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM node:22-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

RUN npm run build

FROM node:22-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

RUN addgroup --system --gid 10001 nodejs && \
    adduser --system --uid 10001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

CMD ["node", "server.js"]
```

### 4.3 Production Compose Orchestration Baseline

```yaml
# /etc/dink/docker-compose.prod.yaml - Production Stack Baseline
version: "3.9"

networks:
  dmz_net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.28.10.0/24
  app_net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.28.20.0/24
  data_net:
    driver: bridge
    internal: true
    ipam:
      config:
        - subnet: 172.28.30.0/24

services:
  traefik:
    image: traefik:v3.1
    container_name: dink-traefik-edge
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /etc/traefik:/etc/traefik:ro
      - /etc/ssl/certs:/etc/ssl/certs:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    networks:
      dmz_net:
        ipv4_address: 172.28.10.2
      app_net:
        ipv4_address: 172.28.20.2
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "5"

  monolith-primary:
    image: dink-registry.local/dink/core-server:v7.0.0
    container_name: dink-monolith-01
    restart: always
    environment:
      - APP_PORT=8080
      - APP_ENV=production
      - DATABASE_URL=postgres://dink_app_user:${DB_PASSWORD}@172.28.30.10:5432/dink_core_prod?sslmode=verify-full&sslrootcert=/etc/ssl/certs/pg-root.crt
      - VALKEY_ADDR=172.28.30.20:6379
      - VALKEY_PASSWORD=${VALKEY_PASSWORD}
      - CLAMAV_ADDR=172.28.30.40:3310
      - S3_ENDPOINT=[http://172.28.30.50:8333](http://172.28.30.50:8333)
      - S3_BUCKET=dink-production-media
      - OPENBAO_ADDR=[https://172.28.30.60:8200](https://172.28.30.60:8200)
    networks:
      app_net:
        ipv4_address: 172.28.20.10
      data_net:
        ipv4_address: 172.28.30.100
    depends_on:
      - postgres-primary
      - valkey-primary
      - clamav-daemon
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp:rw,noexec,nosuid,size=64m

  admin-console:
    image: dink-registry.local/dink/admin-console:v7.0.0
    container_name: dink-admin-01
    restart: always
    environment:
      - NODE_ENV=production
      - BACKEND_INTERNAL_URL=[http://172.28.20.10:8080](http://172.28.20.10:8080)
    networks:
      app_net:
        ipv4_address: 172.28.20.20
    depends_on:
      - monolith-primary

  postgres-primary:
    image: postgis/postgis:18-3.4-alpine
    container_name: dink-postgres-primary
    restart: always
    environment:
      - POSTGRES_DB=dink_core_prod
      - POSTGRES_USER=dink_app_user
      - POSTGRES_PASSWORD_FILE=/run/secrets/pg_pass
    volumes:
      - /mnt/nvme/pgdata:/var/lib/postgresql/data
      - /etc/postgresql/postgresql.conf:/etc/postgresql/postgresql.conf:ro
    networks:
      data_net:
        ipv4_address: 172.28.30.10
    shm_size: 4g

  valkey-primary:
    image: valkey/valkey:8.0-alpine
    container_name: dink-valkey-primary
    restart: always
    command: ["valkey-server", "/etc/valkey/valkey.conf"]
    volumes:
      - /mnt/nvme/valkeydata:/data
      - /etc/valkey/valkey.conf:/etc/valkey/valkey.conf:ro
    networks:
      data_net:
        ipv4_address: 172.28.30.20

  clamav-daemon:
    image: clamav/clamav:latest
    container_name: dink-clamav-worker
    restart: always
    networks:
      data_net:
        ipv4_address: 172.28.30.40

  seaweedfs:
    image: chrislusf/seaweedfs:latest
    container_name: dink-seaweedfs-s3
    restart: always
    command: ["server", "-dir=/data", "-s3", "-s3.port=8333"]
    volumes:
      - /mnt/nvme/seaweeddata:/data
    networks:
      data_net:
        ipv4_address: 172.28.30.50
```

---

## 5. High Availability, Clustering & Storage Topography

### 5.1 PostgreSQL 18.4 Streaming WAL Replication

```mermaid
flowchart LR
    AppPrimary["Go Monolith Primary Nodes"] -->|Write Traffic / Port 5432| MasterDB[("PostgreSQL 18.4 Primary<br/>Addis Ababa DC (Zone 3)")]
    AppPrimary -.->|Read Splitting / Port 5432| StandbyDB[("PostgreSQL 18.4 Standby<br/>Addis Ababa DC (Zone 3)")]
    MasterDB -->|Physical Streaming WAL Replication| StandbyDB
    MasterDB -->|Continuous WAL Archiving| WALStorage[("pgBackRest Repository<br/>Local NVMe + Remote Cold Storage")]
```

* **Failover Orchestration:** Monitored via automated health agents and Patroni consensus. If the primary instance misses 3 consecutive 2-second heartbeats, the standby node promotes to primary and updates internal DNS / VIP routes.
* **Synchronous Commits:** For financial transactions, `synchronous_commit = on` guarantees that transaction commits are acknowledged on disk on the primary before the client request returns.

### 5.2 Valkey 8.0 Sentinel Failover Topology

Distributed mutual exclusion locks (the 60-minute reservation hold `lock:hold:fleet:{id}`) and rate-limiting token buckets require high cache availability:
* **Topology:** 1 Primary node, 2 Replica nodes, and 3 Sentinel instances distributed across independent host machines.
* **Failover SLA:** Downed primary detection within 1,500ms; full failover election and promotion completed in $< 3,000\text{ ms}$.

### 5.3 SeaweedFS Distributed S3 Storage Replication

Photographs, digital signatures, and PDF tax invoices are stored in an on-premise SeaweedFS cluster:
* **Replication Profile `001`:** Each volume is replicated across two separate physical disks within the datacenter rack.
* **Capacity Strategy:** 1 Master Node handling topology and metadata lookups; 3 Volume Nodes offering a combined raw capacity of $24\text{TB}$ NVMe storage.

---

## 6. Local Data Residency & Regional Hosting Topology

### 6.1 Statutory Data Residency Governance (Proclamation No. 1321/2024)

Under Ethiopian law, personal identity data, citizen tax identifiers, and location tracking information cannot leave sovereign borders without formal regulatory dispensation from the Ethiopian Data Protection Authority.

* **Data Residency Rule:** All operational databases (`PostgreSQL`), in-memory caches (`Valkey`), and media attachments (`SeaweedFS`) reside physically inside the Addis Ababa Tier III Datacenter facility.
* **Cross-Border Edge Routing:** Kenya (`KE`) traffic enters via a dedicated Point of Presence (PoP) in Nairobi. Cross-border mobile requests are reverse-proxied over dedicated private peering fiber to the Addis Ababa primary core; financial records are persisted directly to the sovereign Ethiopian data layer.

### 6.2 Cross-Border Edge Peering (Addis Ababa Core & Nairobi PoP)

```mermaid
flowchart TD
    subgraph AddisAbabaDC [Primary Sovereign Datacenter - Tier III Addis Ababa]
        direction TB
        AddisEdge["Traefik 3.x Primary Gateway"]
        AddisMonolith["Go Monolith Replica Pool"]
        AddisDB[("PostgreSQL 18.4 + PostGIS Cluster")]
        AddisS3[("SeaweedFS S3 Cluster")]
        AddisKMS["OpenBao Sovereign KMS"]

        AddisEdge --> AddisMonolith
        AddisMonolith --> AddisDB
        AddisMonolith --> AddisS3
        AddisMonolith --> AddisKMS
    end

    subgraph NairobiPoP [Cross-Border Edge PoP - Nairobi Kenya]
        direction TB
        NairobiEdge["Traefik Edge Reverse Proxy<br/>Local SSL Termination & Caching"]
        NairobiUsers["Kenyan Contractors & Fleet Owners"]
        NairobiUsers --> NairobiEdge
    end

    subgraph EthiopianUsers [Ethiopian Mobile & Web Clients]
        EthUsers["Addis Ababa & Regional Contractors"]
    end

    EthUsers -->|Ethio Telecom / Safaricom ET LTE| AddisEdge
    NairobiEdge -->|Encrypted Cross-Border Private Peering| AddisEdge
```

---

## 7. CI/CD Pipeline & Automated Release Engineering

### 7.1 Automated Quality & Invariant Gates

Every pull request initiates a pipeline enforcing the architectural invariants established across earlier documentation sections:

```mermaid
flowchart TD
    Commit["Git Push / Pull Request Trigger"] --> LintGate{"Gate 1: Static Analysis<br/>golangci-lint, biome, spectral"}
    LintGate -->|Fail| AbortPR["Abort Build & Notify Author"]
    
    LintGate -->|Pass| ContractGate{"Gate 2: Contract Supremacy<br/>kin-openapi contract linting"}
    ContractGate -->|Fail| AbortPR
    
    ContractGate -->|Pass| ForbiddenAPIGate{"Gate 3: Invariant Enforcement<br/>Banned Floats & deleted_at check"}
    ForbiddenAPIGate -->|Fail| AbortPR
    
    ForbiddenAPIGate -->|Pass| TestGate{"Gate 4: Automated Testing<br/>Unit, BLoC, Testcontainers Integration"}
    TestGate -->|Fail| AbortPR
    
    TestGate -->|Pass| SecurityGate{"Gate 5: DevSecOps Scan<br/>Semgrep SAST & Trivy Container Scan"}
    SecurityGate -->|Fail| AbortPR
    
    SecurityGate -->|Pass| ImageBuild["Build Static Distroless Images<br/>Sign with Cosign & Push to Registry"]
```

#### The Five Invariant Checks Enforced in CI
1. **Contract Supremacy:** Validates that `contracts/openapi.yaml` matches backend models with zero schema drift.
2. **Forbidden Float Check:** Custom AST grep rules verify that `float32` and `float64` types are not used for monetary calculations.
3. **Database Migration Checks:** Verifies that every new table migration defines `version INT NOT NULL DEFAULT 1` (OCC) and `deleted_at TIMESTAMPTZ DEFAULT NULL` (Soft-delete).
4. **Zero-Sum Ledger Test:** Integration tests verify that mutating financial operations post ledger journal entries summing algebraically to zero minor units.
5. **ASVS Level 3 Static Scan:** Semgrep scans codebases to verify PII scrubbing, secure crypto usage (AES-256-GCM), and ES256 signature verification.

### 7.2 Zero-Downtime Blue-Green Deployment Lifecycle

Deployments use an automated blue-green container substitution strategy:
1. **Database Migration Step:** Goose executes backward-compatible SQL schema updates (`goose up`).
2. **Green Instance Spin-up:** The new container version starts on an alternate internal port.
3. **Health Check Assertion:** Traefik probes `GET /healthz` on the green instance until 5 consecutive successful responses are recorded.
4. **Traffic Cutover:** Traefik shifts the load-balancer weights to 100% Green.
5. **Connection Draining:** Blue instances receive a `SIGTERM` signal, gracefully completing existing HTTP requests and WebSocket telematics loops within a 30-second drain window before terminating.

---

## 8. Disaster Recovery, Backup & Point-in-Time Recovery (PITR)

### 8.1 Recovery Objectives (RPO & RTO)

* **Recovery Point Objective (RPO):** $\le 1\text{ hour}$ for catastrophic physical failure; target $0\text{ seconds}$ for database crashes via synchronous WAL replication.
* **Recovery Time Objective (RTO):** $\le 4\text{ hours}$ to restore full operational capability from cold disaster backups.

### 8.2 Automated pgBackRest Cadence

* **Continuous WAL Archiving:** PostgreSQL pushes WAL segments to an encrypted, dedicated backup volume immediately upon file rotation (typically every 16MB or 5 minutes).
* **Full Snapshots:** Executed every Sunday at `01:00 EAT`.
* **Differential Snapshots:** Executed daily at `01:00 EAT` (Monday through Saturday).
* **Point-in-Time Recovery (PITR):** Enables reconstructing the database to any specific microsecond within the rolling 30-day retention window.

```ini
# /etc/pgbackrest/pgbackrest.conf - Production Archival Configuration
[global]
repo1-path=/mnt/backups/pgbackrest
repo1-retention-full=4
repo1-retention-diff=14
repo1-cipher-type=aes-256-cbc
repo1-cipher-pass=ENCRYPTED_BACKUP_PASSPHRASE_FROM_OPENBAO
process-max=4
log-level-console=info
log-level-file=detail
start-fast=y
compress-type=zst
compress-level=6

[dink_core_prod]
pg1-path=/var/lib/postgresql/data
pg1-user=postgres
pg1-port=5432
```

### 8.3 Offsite Cold Storage Encryption

Daily differential backups and weekly full database archives are synced offsite to a physically isolated, air-gapped facility outside Addis Ababa over an encrypted point-to-point tunnel. Archives are encrypted with AES-256-CBC using keys managed outside the production database host.

---

## 9. Observability, Telemetry & SRE Stack

```mermaid
flowchart TD
    subgraph DataProducers [Monitored Services]
        GoCore["Go Monolith Core Engine"]
        AdminWebCore["Next.js 16 Web Console"]
        HostNode["Linux Host OS (Node Exporter)"]
        TraefikEdgeCore["Traefik Edge Gateway"]
    end

    subgraph TelemetryPipeline [Telemetry Ingestion Pipeline]
        PrometheusServer["Prometheus Time-Series Scraper"]
        OTelCollector["OpenTelemetry Trace Collector"]
        PromtailAgent["Promtail PII-Sanitizing Agent"]
    end

    subgraph VisualizationLayer [Visualization & Alerting]
        GrafanaUI["Grafana Dashboards (Zone 4)"]
        LokiEngine["Grafana Loki Log Storage"]
        AlertManager["Alertmanager (PagerDuty / SMS Webhook)"]
    end

    GoCore -->|/metrics (Port 8080)| PrometheusServer
    TraefikEdgeCore -->|/metrics (Port 8082)| PrometheusServer
    HostNode -->|/metrics (Port 9100)| PrometheusServer

    GoCore -->|gRPC Spans (Port 4317)| OTelCollector
    GoCore -->|JSON Log Files| PromtailAgent

    PromtailAgent --> LokiEngine
    PrometheusServer --> GrafanaUI
    PrometheusServer --> AlertManager
    OTelCollector --> GrafanaUI
    LokiEngine --> GrafanaUI
```

### 9.1 Prometheus Metrics Harvesting

The Go monolith exposes operational metrics over `GET /metrics` formatted for Prometheus:
* **HTTP Metrics:** Request rate, error rate, and duration buckets (`http_requests_total`, `http_request_duration_seconds_bucket`).
* **Financial Metrics:** Active reservation hold count (`booking_holds_active_count`), double-entry balance assertion errors (`financial_ledger_imbalance_errors_total`).
* **Database Pool Metrics:** Open connections, idle connections, wait durations (`pgxpool_active_conns`, `pgxpool_acquire_duration_seconds`).

### 9.2 Distributed Tracing with OpenTelemetry

Transactions propagate W3C trace contexts (`traceparent` header). OpenTelemetry instrumentation records database queries, cache hits, and external telco API requests, enabling end-to-end distributed latency analysis in Grafana Tempo.

### 9.3 Zero-PII Log Aggregation with Promtail & Loki

Promtail instances harvest structured JSON logs emitted by container runtimes. A mandatory sanitization stage ensures that any residual phone numbers, tokens, or national ID identifiers are scrubbed before log ingestion into Grafana Loki.

### 9.4 Health Check Probes & Readiness Probes

The Go monolith exposes dual-tier health probing:
1. **Liveness Probe (`GET /healthz`):** Verifies that the internal HTTP event loop is active. Returns `HTTP 200 OK`.
2. **Readiness Probe (`GET /readyz`):** Verifies downstream dependencies:
   - PostgreSQL connection pool ping (`pgxpool.Ping()`).
   - Valkey cache reachability (`valkey.Ping()`).
   - ClamAV socket responsiveness.
   - If any core dependency fails, returns `HTTP 503 Service Unavailable`, prompting Traefik to take the instance out of rotation.

---

## 10. Secrets Provisioning with OpenBao KMS

Production secrets (PostgreSQL credentials, Valkey authentication tokens, ES256 private keys, and AES-256-GCM data encryption keys) are managed through **OpenBao KMS** (an open-source fork of HashiCorp Vault):

```mermaid
flowchart TD
    ContainerBoot["Container Bootstrap Engine"] --> VaultAuth["Authenticate using AppRole / mTLS"]
    VaultAuth --> TokenGrant["Issue Ephemeral Lease Token (1-Hour TTL)"]
    TokenGrant --> SecretsRequest["Request Database Passwords & AES Master Keys"]
    SecretsRequest --> IngestMemory["Mount directly into Process Memory"]
    IngestMemory --> RunDaemon["Run Monolith Core (Zero Disk Residue)"]
```

* **No Plaintext Configuration Files:** Secrets are never committed to Git or stored in Docker environment variables on disk.
* **Lease Expiration & Dynamic Revocation:** Database credentials use short-lived dynamic leases rotated automatically every 24 hours.

---

## 11. Operational Runbooks & Incident Playbooks

```text
+------------------------------------------------------------------------------------+
|                         PRODUCTION INCIDENT PLAYBOOKS                              |
+------------------------------------------------------------------------------------+
|  Playbook 1: Database Primary Failover                                             |
|  - Trigger: PostgreSQL primary fails to respond for > 10 seconds.                  |
|  - Action: Patroni / SRE promotes standby replica to primary.                       |
|  - Verify: Assert pgxpool connections switch to new primary without errors.        |
|  - Recovery: Re-provision former primary as new streaming replica.                 |
|                                                                                    |
|  Playbook 2: Valkey Cache Flush & Reservation Hold Re-Sync                        |
|  - Trigger: Valkey cluster corruption or cold reboot.                             |
|  - Action: Flush cache. Execute background recovery task in Go backend:            |
|    SELECT fleet_id, renter_id, hold_expires_at FROM bookings WHERE status='HOLD_PENDING'|
|  - Recovery: Re-populate Valkey mutex keys (SET fleet_id renter_id EX remaining).  |
|                                                                                    |
|  Playbook 3: ClamAV Daemon Failure                                                 |
|  - Trigger: ClamAV TCP socket connection refused.                                  |
|  - Action: Upload controller temporarily transitions to HTTP 503 for new media      |
|    uploads; existing text and booking operations proceed normally.                 |
|  - Recovery: Restart ClamAV daemon container and verify virus database reload.     |
+------------------------------------------------------------------------------------+
```

---

*This document defines the canonical deployment and infrastructure architecture for Dink Machinery Version 7.0. All downstream infrastructure provisioning, container deployments, and disaster recovery procedures must strictly adhere to these standards.*