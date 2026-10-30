# ADR-012: Phase 0 Contract-First Governance, Sovereign Colocation and SeaweedFS S3

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 01-introduction, 05-06, 05-13, 09-non-functional-requirements, 12-tech-stack
- **Related Requirements:** MED-01-BE, SEC-05-BE, INF-01-OPS, ETH-DATA-1321-2024

---

## 1. Context & Problem Statement

Dink Machinery processes sensitive national citizen records, commercial equipment ownership titles, trade licenses, and high-resolution inspection media. Deploying to public cloud hyperscalers (AWS Ireland, Azure South Africa) introduces severe foreign exchange currency (FX) payment constraints under National Bank of Ethiopia regulations, causes high cross-border network latency over cellular networks, and violates data sovereignty mandates under Ethiopian Personal Data Protection Proclamation No. 1321/2024.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated on-premise Tier III datacenter colocation with SeaweedFS S3 vs. public cloud hyperscalers (AWS S3) vs. basic office server hosting.
- [x] **Surprise Test:** Chooses self-hosted SeaweedFS S3 object storage and bare-metal colocation in Addis Ababa over convenient global hyperscalers.
- [x] **Constraint Test:** Driven by mandatory legal data residency compliance (Proclamation No. 1321/2024) and central bank foreign exchange limits.
- [x] **Irreversibility Test:** Fundamental infrastructure and storage architecture; changing physical hosting location requires complex data migration drills.
- [x] **Scope Test:** Governs all media storage, database persistence, network ingress, antivirus pipelines, and backup topologies.

---

## 3. Considered Options

### Option 1: Sovereign Tier III Datacenter Colocation + Self-Hosted SeaweedFS S3 + ClamAV
- **Summary:** Lease bare-metal server infrastructure in certified Tier III datacenters (Raxio / Wingu DC in Addis Ababa paired with an edge PoP in Nairobi). Store inspection photographs, trade licenses, and PDF contracts in a self-hosted SeaweedFS S3 cluster. Stream all uploads through synchronous ClamAV `zINSTREAM` sockets (:3310) prior to persistence.
- **Pros:**
  - 100% statutory compliance with Ethiopian and Kenyan data residency laws.
  - Zero foreign currency (USD) monthly cloud billing; operational expenses paid in local currency (`ETB`, `KES`).
  - Sub-10ms domestic latency across local telecom exchanges (EthIO-IX).
  - Open-source, high-throughput media storage with full AWS S3 API compatibility.
- **Cons:**
  - Requires engineering responsibility for hardware capacity management and disk replacement drills.

### Option 2: Public Cloud Hyperscalers (AWS EU-West / AWS Cape Town)
- **Summary:** Host databases on Amazon RDS and media on Amazon S3.
- **Pros:**
  - Fully managed infrastructure and automated scalability.
- **Cons:**
  - Explicit violation of Ethiopian citizen data residency laws; severe foreign exchange currency shortages prevent paying monthly AWS USD invoices; 120ms+ network latency over undersea cables.

### Option 3: MinIO S3 on Local Office Hardware
- **Summary:** Run MinIO on office rack servers.
- **Pros:**
  - Low capital cost.
- **Cons:**
  - MinIO AGPLv3 licensing liabilities; office power and fiber line single points of failure.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Sovereign Tier III Datacenter Colocation + Self-Hosted SeaweedFS S3 + ClamAV**.

1. **Physical Hosting Topology:**
   - **Primary Core:** Tier III certified colocation facility in Addis Ababa (dual utility power feeds, N+1 generators, 24/7 physical biometric security).
   - **Regional Edge PoP:** Colocation rack in Nairobi terminating Kenyan user traffic and routing over an encrypted private tunnel.
2. **Object Storage Engine (SeaweedFS S3):**
   - High-performance Go-based distributed object store optimized for billions of small and medium inspection photos ($< 2\text{MB}$).
   - Replicated locally across physical NVMe drive arrays (profile `001`).
   - S3-compatible REST API allows transparent use of standard AWS S3 SDKs in Go and Dart.
3. **Synchronous Malware Inspection (ClamAV):**
   - Upload streams pass through ClamAV `zINSTREAM` over local Unix sockets.
   - Any malicious or infected file is aborted mid-flight before writing bytes to SeaweedFS, returning an RFC 9457 HTTP 422 error (`INSPECT_MEDIA_SCAN_FAILED`).

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (AWS Hyperscaler):** Rejected due to legal illegality under national data residency proclamations, crippling foreign currency allocation bottlenecks, and high latency for field mobile apps.
- **Option 3 (MinIO / Office Hosting):** Rejected due to commercial enterprise unreliability and restrictive AGPLv3 licensing encumbrances.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Total operational sovereignty and full compliance with national security and privacy laws.
  - Immunity from foreign exchange payment freezes or international cloud price hikes.
  - Blazing fast image upload and retrieval speeds on regional cellular networks.
- **Negative / Technical Debt:**
  - DevOps team must manage disk monitoring, offsite encrypted backups, and physical hardware procurement.
- **Operational / Observability Impacts:**
  - SeaweedFS volume disk usage and ClamAV scan latencies are monitored via Prometheus alerts.

---

## 7. Verification Strategy & Invariant Tests

- **ClamAV Stream Scan Test:** `backend/internal/platform/storage/clamav_test.go` streams the EICAR standard antivirus test string to the upload endpoint; asserts immediate rejection with HTTP 422 Problem Details.
- **S3 API Compatibility Test:** Verifies that the Go backend and Dart mobile client read and write objects using standard AWS S3 SDK drivers without custom patches.

---

## 8. Revisit Conditions

- Local cloud providers in Ethiopia receive official Tier III certification, achieve central bank compliance certification, and accept local currency payments at competitive pricing.
