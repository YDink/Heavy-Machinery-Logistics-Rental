# Partition Decisions: 07-06 — Phase 3 Direct CAN-Bus Telematics & IoT Hardware Integration Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `07-06`.

- **Governing System ADRs:** `ADR-007`, `ADR-009`
- **Subsystem Scope:** backend
- **Requirement Targets:** IOT-01-BE to IOT-12-BE

---

## D-07-06-001: High-Throughput MQTT / J1939 CAN-Bus Telematics Ingestion Pipeline

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [IOT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
In Phase 3, hardware OBD/CAN-bus telematics gateways installed directly on machinery stream real-time engine metrics: true J1939 engine hours, fuel consumption rate, hydraulic temperature, fault codes (DTCs), and GPS. Ingesting tens of thousands of high-frequency telemetry streams must not overload the core HTTP API.

### 2. Considered Alternatives
1. **Option A:** Dedicated lightweight MQTT / TCP broker gateway (e.g. EMQX / Go MQTT ingestor) writing to time-series optimized hypertable in PostgreSQL (TimescaleDB / partitioning)
   - *Pros:* Low-latency binary protocol; handles intermittent cellular reconnections cleanly; decouples telematics ingestion from REST API.
   - *Cons:* Requires running an MQTT ingestion gateway container.
2. **Option B:** Sending telematics data via standard HTTP POST calls to the main Chi monolith API
   - *Pros:* Reuses existing HTTP routes.
   - *Cons:* HTTP overhead and TLS handshakes exhaust server threads under thousands of concurrent IoT devices.
3. **Option C:** Third-party proprietary telematics cloud (e.g. Caterpillar VisionLink API only)
   - *Pros:* Vendor-managed ingestion.
   - *Cons:* Locked to single manufacturer; cannot ingest multi-brand fleets (Komatsu, SANY, CAT, XCMG) in emerging markets.

### 3. Chosen Solution & Technical Rationale
Option A: Dedicated MQTT ingestion pipeline. IoT hardware gateways transmit compact J1939 telemetry frames over MQTT with TLS 1.3. An in-process MQTT worker parses frames, validates device cryptographic credentials, and writes downsampled metrics to partitioned storage.

### 4. Rejected Alternatives
- Option B: Rejected because HTTP/REST protocol overhead is unsuited for continuous streaming IoT devices.
- Option C: Rejected because African equipment fleets are highly heterogeneous and multi-brand.

### 5. Consequences & Invariants
- **Invariants Maintained:** Telematics engine hour telemetry automatically reconciles and updates the official verified fleet engine hours.
- **Operational Trade-offs:** Requires installing and wiring physical telematics hardware into machine CAN-bus harnesses.

### 6. Verification
- **Test File:** `backend/internal/modules/telematics/infrastructure/mqtt_ingest_test.go`
- **Assertion:** TestMQTT_ParsesJ1939Frame asserts valid CAN-bus frame extracts engine hours and fuel level accurately.

### 7. Revisit Trigger
- Revisit if satellite IoT transceivers (Iridium/Starlink Direct-to-Cell) become standard on heavy equipment.

---

## D-07-06-002: Real-Time Geofence Violation and Tamper Alerting with Remote Engine Immobilization Capability

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [IOT-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Machinery theft or unauthorized night operation outside contracted hours is a severe threat for fleet owners. Owners demand immediate alerts and the ability to remotely immobilize equipment (prevent engine restart) during verified breach or non-payment disputes.

### 2. Considered Alternatives
1. **Option A:** Hardware relay immobilization: backend sends authenticated cryptographic immobilize command to telematics gateway, preventing starter motor actuation after next engine shutdown
   - *Pros:* Prevents high-speed vehicle shutdown hazards (safely immobilizes only when machine is stationary); stops machine theft cold.
   - *Cons:* Requires strict multi-signature administrative authorization to prevent accidental immobilization.
2. **Option B:** Immediate mid-operation engine kill-switch cutting fuel solenoid instantly
   - *Pros:* Stops machine immediately.
   - *Cons:* Extreme safety hazard: sudden engine loss on a moving crane or excavator on a steep slope can cause fatal rollovers.
3. **Option C:** Passive GPS alerting only without any remote immobilization hardware
   - *Pros:* Zero equipment intervention.
   - *Cons:* Useless against professional equipment thieves who disable GPS trackers once the machine is moved.

### 3. Chosen Solution & Technical Rationale
Option A: Safe starter-relay immobilization. The platform allows verified owners or super admins to issue an IMMOBILIZE command. The telematics firmware verifies the machine is stationary (speed == 0, hydraulic pressure == 0) and disables starter relay on next engine stop.

### 4. Rejected Alternatives
- Option B: Rejected because sudden inflight engine cutoffs create unacceptable job site fatality risks.
- Option C: Rejected because passive tracking does not prevent stolen equipment from crossing national borders.

### 5. Consequences & Invariants
- **Invariants Maintained:** Remote immobilization can only engage when verified machine ground speed is zero and requires two-factor authorization.
- **Operational Trade-offs:** Immobilization hardware relays must be professionally installed and tamper-monitored.

### 6. Verification
- **Test File:** `backend/internal/modules/telematics/application/immobilization_test.go`
- **Assertion:** TestImmobilize_BlockedWhileMachineMoving asserts immobilize command is rejected if engine RPM > 0 and speed > 0.

### 7. Revisit Trigger
- Revisit if national transport regulations mandate specific remote shutdown certification.

---

