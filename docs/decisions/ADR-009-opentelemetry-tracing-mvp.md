# ADR-009: OpenTelemetry Tracing, RED Metrics and In-Flight Zero-PII Log Scrubbing

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-10, 05-24, 09-non-functional-requirements, 12-tech-stack
- **Related Requirements:** SEC-08-BE, LOG-04-OPS, ASVS-V8.1, ECA-PRIV-01

---

## 1. Context & Problem Statement

Dink Machinery captures sensitive citizen data, commercial equipment titles, national identity verifications, and financial transactions. Application telemetry (logs, distributed traces, error metrics) must provide deep operational observability without ever leaking Personally Identifiable Information (PII) such as phone numbers, National IDs, TINs, or bank account numbers. Ethiopian Data Protection Proclamation No. 1321/2024 mandates strict data minimization and imposes severe statutory penalties for plaintext personal data leakage.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated in-flight application-layer regex masking vs. post-ingestion SIEM masking vs. manual developer log hygiene.
- [x] **Surprise Test:** Plaintext PII is intercepted and masked in memory before log serialization; raw data never reaches disk or shipping sockets.
- [x] **Constraint Test:** Driven by legal mandate (Proclamation No. 1321/2024) and OWASP ASVS Level 3 zero-PII logging standards.
- [x] **Irreversibility Test:** Central logging and telemetry framework across all backend modules and infrastructure layers.
- [x] **Scope Test:** Governs all logging cores, OpenTelemetry instrumentation, Prometheus metrics, and log shippers.

---

## 3. Considered Options

### Option 1: In-Flight Application-Layer Masking with Uber Zap and OpenTelemetry
- **Summary:** Custom `RedactingZapCore` wraps the Zap logger in memory. Pre-compiled regular expressions match and redact phone numbers (`+25191****344`), TINs, National IDs, and bank account numbers prior to JSON serialization. OpenTelemetry auto-instruments HTTP requests, pgx queries, and Valkey commands using RED metrics (Rate, Errors, Duration).
- **Pros:**
  - Plaintext PII never hits physical log files or network shipping pipelines.
  - Provable statutory compliance under automated static analysis and dynamic penetration testing.
  - Consistent trace context (`trace_id`, `span_id`) correlates logs and traces without leaking sensitive payloads.
- **Cons:**
  - Minor CPU overhead for regex evaluation during log emission.

### Option 2: Post-Ingestion SIEM Masking / Log Collector Redaction
- **Summary:** Raw logs are written to disk and shipped to Grafana Loki / Vector, where regex pipelines redact sensitive data.
- **Pros:**
  - Zero CPU impact on the Go application runtime.
- **Cons:**
  - Raw unredacted citizen PII resides on local disks and in shipping queues, violating ASVS L3 and creating data breach liabilities during disk backup snapshots.

### Option 3: Manual Developer Hygiene & Code Review Guidelines Alone
- **Summary:** Relying on developer discipline and PR reviews to avoid logging sensitive variables.
- **Pros:**
  - Zero performance overhead.
- **Cons:**
  - Completely unreliable; human error inevitably leads to logging entire request structs or error messages containing phone numbers and credentials.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: In-Flight Application-Layer Masking with Uber Zap and OpenTelemetry**.

1. **Two-Tier Masking Architecture:**
   - **Tier 1 (Application Memory):** `RedactingZapCore` executes pre-compiled regex substitutions on all logged fields:
     - Phone Numbers (`^\+?[1-9]\d{1,14}$`): Masked to keep dialing prefix and last 3 digits (`+25191****344`).
     - National IDs & TINs: Masked to `***REDACTED_ID***`.
     - Bearer Tokens & Secrets: Fully stripped.
   - **Tier 2 (Log Shipper):** Promtail / Vector log pipeline enforces secondary regex scrubbing as defense-in-depth prior to ingestion into Grafana Loki.
2. **Standardized Telemetry Schema:**
   All logs output structured JSON including: `timestamp`, `level`, `caller`, `message`, `requestId`, `userId` (UUID only, never phone/name), `traceId`, `spanId`, `module`.
3. **RED Metrics:**
   OpenTelemetry instruments Chi HTTP routes, emitting Prometheus metrics for request rate, error rate, and duration histograms.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Post-Ingestion Masking):** Rejected because writing plaintext PII to local disk violates Ethiopian Data Protection Proclamation No. 1321/2024 and OWASP ASVS L3 requirements.
- **Option 3 (Manual Hygiene):** Rejected because human vigilance fails predictably under production incident pressures.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Zero PII leaks into log files, search indexes, or third-party monitoring services.
  - SREs and developers can safely inspect production logs without security clearances.
  - Full traceability across distributed transactions using correlation IDs.
- **Negative / Technical Debt:**
  - Negligible CPU processing overhead during high-volume log writes ($< 2\%$ increase in CPU utilization).
- **Operational / Observability Impacts:**
  - Prometheus alert triggers if log emission rate or log error count spikes unexpectedly.

---

## 7. Verification Strategy & Invariant Tests

- **Zero-PII Regex Audit Script:** CI executes static and dynamic scans:
  ```bash
  grep -rnEi 'logger\.(Info|Debug|Warn|Error)\(.*(password|pin|tin|national_id|account_number)' backend/
  ```
- **Redaction Unit Test:** `backend/internal/platform/telemetry/redaction_test.go` logs synthetic phone numbers, bank accounts, and bearer tokens; asserts that output JSON contains zero plaintext PII substrings.

---

## 8. Revisit Conditions

- Introduction of machine learning telemetry processors requiring homomorphic or pseudonymized mathematical transformations beyond standard regex masking.
