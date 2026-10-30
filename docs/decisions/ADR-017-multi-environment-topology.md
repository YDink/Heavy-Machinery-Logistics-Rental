# ADR-017: Multi-Environment Architecture & Flavor Strategy (Dev, Staging, Prod)

- **Status:** Accepted
- **Date:** 2026-10-30
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 09-non-functional, 10-data-model, 12-tech-stack
- **Related Requirements:** INF-01-OPS, ENV-01-OPS, SEC-05-BE

---

## 1. Context & Problem Statement

Dink Machinery operates across three client and server tiers:
1. **Go 1.27.3 Backend Modular Monolith**
2. **Next.js 16 Administrative Console**
3. **Flutter 3.47.9 Mobile Client (Renter, Owner, Operator)**

The platform handles real-money escrow transactions, automated double-entry ledger mutations, live equipment GPS breadcrumbs, and national citizen identity verifications under the Ethiopian Personal Data Protection Proclamation No. 1321/2024. 

Operating without strictly isolated multi-tier environments creates critical risks:
- Developers running local experiments could accidentally trigger live SMS OTP dispatches, charge real bank accounts via Telebirr/CBE, or corrupt production database tables.
- QA and automated E2E test suites would compete for resources with live customer traffic or expose unreleased features to end users.
- Testers and operations staff could not run internal staging mobile builds alongside the production consumer app on the same physical device without overwriting app data.

A formal, isolated, three-tier environment architecture (`Development`, `Staging`, `Production`) is required across all backend services, administrative web portals, and mobile application packages.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated 2-tier (Dev/Prod) vs. 3-tier (Dev/Staging/Prod) vs. dynamic ephemeral branch environments.
- [x] **Surprise Test:** Establishes distinct mobile application IDs (`.dev`, `.staging`, prod) allowing side-by-side device installations, and mandates mock OTP/payment modes in development while enforcing real sandbox rails in staging.
- [x] **Constraint Test:** Driven by legal data sovereignty (Proclamation No. 1321/2024) mandating that live citizen PII must never be mirrored into unencrypted or unhardened development workstations.
- [x] **Irreversibility Test:** Infrastructure configuration, database migration pipelines, mobile bundle identifiers, app store release tracks, and secrets provisioning are tightly bound to this environment topology.
- [x] **Scope Test:** Cross-cuts the entire ecosystem: Go backend configuration, Next.js environment injection, Flutter mobile build flavors, Docker compose stacks, and CI/CD promotion pipelines.

---

## 3. Considered Options

### Option 1: Formal 3-Tier Multi-Environment Architecture with Mobile Flavors (Chosen)
- **Summary:** Enforce three distinct tiers across all subsystems:
  - **Development (`dev`):** Local Docker Compose infrastructure (`localhost`), `.env.dev` configuration, mock OTP/sandbox payment adapters, relaxed CORS, Flutter flavor `dev` (`com.dinkmachinery.dink_mobile.dev`), Next.js `localhost:3000`.
  - **Staging (`staging`):** Dedicated staging server hosted in Addis Ababa, staging PostgreSQL 18.6 with synthetic/anonymized data, sandbox telco SMS rails, staging payment gateway sandbox, Flutter flavor `staging` (`com.dinkmachinery.dink_mobile.staging`) distributed via TestFlight / Google Play Internal App Sharing, Next.js `staging-admin.dinkmachinery.com`.
  - **Production (`prod`):** Colocated Tier III facility in Addis Ababa, clustered active-active PostgreSQL with streaming WAL replication, live telecom SMS, live bank APIs, strict ASVS L3 zero-PII logging, Flutter production build (`com.dinkmachinery.dink_mobile`), Next.js `admin.dinkmachinery.com` with mandatory MFA and strict CSP.
- **Pros:**
  - Absolute data safety: zero risk of test actions mutating production ledgers or live payment rails.
  - Side-by-side mobile testing: QA engineers can install Dev, Staging, and Production builds simultaneously on a single smartphone.
  - Complete compliance with national data privacy regulations (PII strictly quarantined to production).
- **Cons:**
  - Requires maintaining configuration manifests, environment files, and build pipelines across three environments.

### Option 2: 2-Tier Architecture (Local Development & Direct Production)
- **Summary:** Eliminate the intermediate staging environment; test locally and deploy directly to production.
- **Pros:**
  - Minimal infrastructure cost and fewer deployment pipelines.
- **Cons:**
  - High probability of customer-impacting regressions during releases.
  - Inability to run full-scale E2E regression drills against real telecom SMS and payment sandbox gateways prior to production launch.

### Option 3: Dynamic Ephemeral Cloud Environments per Pull Request
- **Summary:** Spin up on-demand cloud environments for each pull request using Kubernetes / AWS ECS.
- **Pros:**
  - Isolated test harness per PR.
- **Cons:**
  - Prohibited by sovereign data residency mandates (Proclamation No. 1321/2024) and central bank foreign currency restrictions for cloud hosting.
  - Excessive DevOps operational complexity for a modular monolith platform.

---

## 4. Chosen Solution & Architectural Rationale

The platform adopts **Option 1: Formal 3-Tier Multi-Environment Architecture with Mobile Flavors**.

### Subsystem Environment Matrix

| Dimension | Development (`dev`) | Staging (`staging`) | Production (`prod`) |
|---|---|---|---|
| **Go Backend URL** | `http://localhost:8080/api/v1` | `https://staging-api.dinkmachinery.com/api/v1` | `https://api.dinkmachinery.com/api/v1` |
| **Go Log Level** | `DEBUG` (human-readable / console) | `INFO` (structured JSON) | `INFO` (structured JSON + PII masking) |
| **Go DB Ingress** | Local Docker Compose (`localhost:5432`) | Staging PostgreSQL (NVMe, Dedicated Host) | Clustered Primary-Replica WAL (Tier III DC) |
| **Go OTP / SMS** | Mock Mode (fixed `123456` or console log) | Telco Sandbox APIs (Ethio Tel / Safaricom) | Live Telco Production Aggregator |
| **Go Payment Rails** | Mock Payment Webhook Simulator | Chapa / Telebirr Sandbox Gateway | Live Telebirr / CBE / M-Pesa Rails |
| **Admin Web URL** | `http://localhost:3000` | `https://staging-admin.dinkmachinery.com` | `https://admin.dinkmachinery.com` |
| **Admin Env File** | `admin/.env.development` | `admin/.env.staging` | `admin/.env.production` (Secrets Injection) |
| **Admin Features** | Mock auth bypass, React DevTools | "STAGING" banner watermark, E2E harness | Strict CSP, Mandatory MFA, Audit Log |
| **Mobile Flavor** | `dev` | `staging` | `prod` |
| **Mobile App ID** | `com.dinkmachinery.dink_mobile.dev` | `com.dinkmachinery.dink_mobile.staging` | `com.dinkmachinery.dink_mobile` |
| **Mobile App Name** | `Dink (Dev)` | `Dink (Staging)` | `Dink Machinery` |
| **Mobile Entrypoint**| `mobile/lib/main_dev.dart` | `mobile/lib/main_staging.dart` | `mobile/lib/main_prod.dart` |
| **Mobile Storage** | Embedded Isar `isar_dev` | Embedded Isar `isar_staging` | Embedded Isar `isar_prod` |
| **Mobile SSL Pinning**| Disabled (allows Charles/Proxyman) | Disabled (allows QA proxy inspection) | Enforced (strict public key pinning) |

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (2-Tier Dev/Prod):** Rejected because staging acceptance drills (G3 quality gate) are legally and operationally mandatory for financial systems under NBE directives.
- **Option 3 (Ephemeral Cloud Environments):** Rejected due to strict sovereign colocation mandates in Ethiopia prohibiting foreign hyperscaler persistence for financial transactions and identity verifications.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Strict isolation: zero accidental leakage of test events into production ledgers or real banking rails.
  - Side-by-side mobile coexistence: engineers, QA, and executives can run all three builds on a single physical device without cache corruption.
  - Uncompromised compliance with Ethiopian Personal Data Protection Proclamation No. 1321/2024.
- **Negative / Technical Debt:**
  - Mobile developers must pass `--flavor <env> -t lib/main_<env>.dart` during Flutter CLI operations.
  - Infrastructure team must manage staging environment patching and database schema migrations alongside production.
- **Operational / Observability Impacts:**
  - Traces and metrics are tagged with `service.environment` (`development`, `staging`, `production`) across OpenTelemetry and Prometheus dashboards.

---

## 7. Verification Strategy & Invariant Tests

- **Test Suite:** Configuration Validation & Build Smoke Tests
- **Assertions:**
  - `fvm flutter build apk --flavor dev -t lib/main_dev.dart` generates APK with package `com.dinkmachinery.dink_mobile.dev`.
  - `fvm flutter build apk --flavor staging -t lib/main_staging.dart` generates APK with package `com.dinkmachinery.dink_mobile.staging`.
  - `fvm flutter build apk --flavor prod -t lib/main_prod.dart` generates APK with package `com.dinkmachinery.dink_mobile`.
  - Backend integration test verifies that `ENV=production` rejects startup if mock OTP or unmasked logger cores are configured.

---

## 8. Revisit Conditions

- Revisit if sovereign cloud providers in East Africa (e.g., local Ethio Telecom / Raxio cloud zones) achieve ISO 27001 and local residency certification permitting automated ephemeral staging environments.
