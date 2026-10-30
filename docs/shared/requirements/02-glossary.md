# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification — 02. Unified Glossary & Domain Terminology (`docs/shared/requirements/02-glossary.md`)

**Document Identifier:** `SRS-DINK-CORE-V7.0-SEC-02`  
**Classification:** Canonical System Requirements Specification & Domain Glossary (Single Source of Truth)  
**Governing Standard:** IEEE Std 830-1998 / IEEE/ISO/IEC 29148:2018 (Systems and Software Engineering — Requirements Engineering)  
**Associated Artifacts:** `SRS-DINK-CORE-V7.0` (Full Suite §1–§14), `FLOW-ENTRY-DINK-CORE-V7.0` through `FLOW-INSURANCE-CLAIM-DINK-CORE-V7.0` (Interaction Flows 01–19), `FLOW-APPENDIX-NOTIF-OFFLINE-DINK-CORE-V7.0`, `FLOW-APPENDIX-AUTH-SEC-DINK-CORE-V7.0`, `ARCH-BACKEND-GO-DINK-CORE-V7.0`, `ARCH-FLUTTER-DINK-CORE-V7.0`, `ARCH-ADMIN-NEXT-DINK-CORE-V7.0`, `ARCH-SEC-DINK-CORE-V7.0`, `DS-DINK-CORE-V7.0`, `DS-PERSONAS-DINK-CORE-V7.0`, `contracts/openapi.yaml`, `ADR-001` through `ADR-010`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, Flutter 3.47.8 / Dart 3.9 (Cross-Platform Mobile), Next.js 16.x (latest CVE security patch) / React 19 (Web Admin & Corporate Portal), PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3, ClamAV Daemon (:3310), SeaweedFS S3, Isar NoSQL 3.1+  
**Target Operating Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), RFC 7519 (JWT), RFC 6238 (TOTP), ISO 9241-110, W3C WCAG 2.2 Level AAA (Field Mobile) & Level AA (Admin Web), National Bank of Ethiopia (NBE) Directives, Kenya Data Protection Act 2019  
**Status:** Approved Engineering Baseline

---

## Table of Contents

- 1. [Document Scope & Architectural Mandate](#1-document-scope--architectural-mandate)
  - 1.1 [Purpose & Semantic Alignment](#11-purpose--semantic-alignment)
  - 1.2 [Taxonomic Boundaries](#12-taxonomic-boundaries)
- 2. [Universal Business & Legal Domain Terminology](#2-universal-business--legal-domain-terminology)
  - 2.1 [Actors, Identity & Corporate Hierarchies](#21-actors-identity--corporate-hierarchies)
  - 2.2 [Contractual, Commercial & Regulatory Terms](#22-contractual-commercial--regulatory-terms)
- 3. [Heavy Machinery, Attachments & Logistics Terminology](#3-heavy-machinery-attachments--logistics-terminology)
  - 3.1 [Plant Asset Classification & Technical Parameters](#31-plant-asset-classification--technical-parameters)
  - 3.2 [Lowboy Transport & Heavy Haulage Logistics](#32-lowboy-transport--heavy-haulage-logistics)
- 4. [Operational, Field & Inspection Terminology](#4-operational-field--inspection-terminology)
  - 4.1 [Handover Reticle, Hours & Shift Telemetry](#41-handover-reticle-hours--shift-telemetry)
  - 4.2 [Damage Assessment, Wear & Insurance Loss Adjustment](#42-damage-assessment-wear--insurance-loss-adjustment)
- 5. [Financial, Settlement & Commission Terminology](#5-financial-settlement--commission-terminology)
  - 5.1 [Non-Custodial Payment Terms & Rails](#51-non-custodial-payment-terms--rails)
  - 5.2 [Platform Revenue, Wallets & Debt Ledgers](#52-platform-revenue-wallets--debt-ledgers)
- 6. [Technical, Architectural & Infrastructure Terminology](#6-technical-architectural--infrastructure-terminology)
  - 6.1 [Concurrency, Locking & Spatial Indexing](#61-concurrency-locking--spatial-indexing)
  - 6.2 [Offline Synchronization, Storage & Protocols](#62-offline-synchronization-storage--protocols)
  - 6.3 [Security, Cryptography & Identity Protocols](#63-security-cryptography--identity-protocols)
- 7. [Acronym & Abbreviation Directory](#7-acronym--abbreviation-directory)

---

## 1. Document Scope & Architectural Mandate

### 1.1 Purpose & Semantic Alignment

This specification establishes the canonical **Universal Glossary and Domain Terminology** for Dink Machinery Version 7.0. Operating across diverse linguistic contexts (Amharic, Afaan Oromoo, Swahili, and English) and complex regulatory frameworks across Ethiopia and Kenya, the platform demands uncompromising semantic consistency.

The definitions codified herein are legally, architecturally, and procedurally binding across:
* **Wire Protocol Schemas:** Exact entity property names, status strings, and enumeration constants in `contracts/openapi.yaml`.
* **Database Relational Models:** PostgreSQL 18.6 table structures, PostGIS functions, and column constraints.
* **Client-Side Domain Models:** Dart classes in Flutter 3.47.8, TypeScript interfaces in Next.js 16, and Isar NoSQL collection schemas.
* **Operational Legal Contracts:** Binding rental terms, damage claims, transport manifests, and arbitration rulings.

### 1.2 Taxonomic Boundaries

Terms in this glossary are classified into five functional quadrants:
1. **Business, Identity & Legal Governance:** Actors, organizational structures, and compliance definitions.
2. **Plant Asset & Heavy Logistics:** Engineering classifications, physical dimensions, and transport mechanics.
3. **Field Operations & Life-Safety:** Reticles, hour meter telemetry, and inspections.
4. **Financial Ledgers & Monetization:** Non-custodial bank settlement models, dynamic commission structures, and debt recovery.
5. **Software Engineering & Cryptography:** Distributed concurrency primitives, offline sync mechanics, and security controls.

---

## 2. Universal Business & Legal Domain Terminology

### 2.1 Actors, Identity & Corporate Hierarchies

#### Civil Contractor (Renter)
The individual proprietor, commercial builder, or general contracting enterprise hiring heavy plant machinery or logistics services to execute civil works. Accountable for site safety, operator deployment limits, equipment physical security during the hire period, and direct payment settlements.

#### Equipment Fleet Owner (Owner)
The natural person or legal commercial entity (e.g., equipment rental house, contractor with surplus assets) holding certified legal title or authorized commercial management rights over registered machinery and transport trailers. Responsible for equipment roadworthiness, mechanical baseline certifications, and team operator assignments.

#### Certified Field Operator (Operator)
A skilled individual holding a valid national or regional heavy plant operator license certified to operate specific machinery classes (e.g., crawler excavators, rough-terrain cranes, motor graders). The operator is physically present in the machinery cab, logs daily shifts, signs digital reticle handovers, and triggers emergency breakdown stoppages.

#### Lowboy Logistics Carrier (Carrier)
A commercial fleet owner or logistics enterprise operating prime movers and multi-axle trailers (lowboys, step-decks, flatbeds) providing heavy plant mobilization and demobilization services.

#### Compliance Officer & Arbitrator
Authorized back-office administrative personnel charged with evaluating KYC verification dossiers, arbitrating contested bank wire declines, adjudicating itemized post-rental damage claims, and auditing manual commission debt clearing receipts.

#### Corporate Organization (`Organization`)
An enterprise entity (e.g., Category 1 General Contractor, infrastructure joint venture) possessing multiple site engineers, centralized procurement hierarchies, cost centers, and approved revolving credit lines.

#### Delegation of Authority (DoA)
An organization's internal administrative configuration specifying spending thresholds, auto-approval boundaries, and multi-signatory approval chains required to execute a corporate equipment requisition.

#### Separation of Duties (SoD)
An immutable platform governance rule asserting that an individual user cannot originate a capital equipment requisition and approve the resulting corporate commitment or credit drawdown.

### 2.2 Contractual, Commercial & Regulatory Terms

#### Booking
A deterministic, legally binding agreement entered into between a Renter and a Fleet Owner for the exclusive operational custody of a designated asset over an inclusive date range $[T_{\text{start}}, T_{\text{end}})$.

#### Parent-Child Booking Hierarchy
The structural relationship linking a primary heavy equipment booking (`ParentBookingID`) to an associated logistics transport sub-booking. Guarantees that logistics cancellations freeze parent machinery readiness.

#### Non-Custodial Marketplace
An operational model complying with National Bank of Ethiopia directives wherein the digital aggregator never holds, pools, or intermediates gross equipment hire funds. All contract settlements flow directly between commercial bank accounts via peer-to-peer wire rails.

#### Force Majeure Logistics Cancellation
A contractual cancellation event triggered when heavy plant transport is rendered physically impossible due to severe mechanical breakdown, road corridor blockades, or regional carrier unavailability, resulting in mutual party exoneration from standard cancellation penalties.

#### Delinquency Lockout
An administrative sanction imposed on a fleet owner whose commission debt balance has remained negative past the 30-day statutory grace period, resulting in the automated un-indexing of all fleet listings from spatial search.

---

## 3. Heavy Machinery, Attachments & Logistics Terminology

### 3.1 Plant Asset Classification & Technical Parameters

#### Fleet Asset (`Fleet`)
Any unique physical capital asset registered on the platform—whether a self-propelled heavy plant machine, a stationary industrial generator, or a non-powered logistics trailer. Identified by an immutable platform UUID, national registration plate, and chassis serial number.

#### Immovable Heavy Plant (`self_move_capable = false`)
Machinery that cannot travel under its own power over public paved roads due to speed restrictions, pavement damage risks, or traffic regulations (e.g., crawler excavators, steel-wheel rollers, tracked bulldozers). Automatically triggers a mandatory lowboy transport integration.

#### Self-Moving Heavy Plant (`self_move_capable = true`)
Pneumatic-tired machinery certified for public roadway travel under its own power between depot yards and job sites (e.g., wheeled backhoe loaders, motor graders, mobile truck cranes).

#### Standalone Attachment (`Attachment`)
A detachable hydraulic or mechanical implement (e.g., hydraulic rock breaker, auger, trenching bucket) rentable independently or bundled with a host prime mover.

#### Dynamic EAV Architecture
The Entity-Attribute-Value schema allowing platform administrators to define technical specifications (`SpecFieldDefinition`) grouped by categories (`SpecFieldGroup`) without mutating physical relational database tables.

#### Indexed Logistics Specifications (`FleetSpecIndexed`)
A denormalized relational table mirroring critical transport parameters (`operating_weight_kg`, `transport_length_mm`, `transport_width_mm`, `transport_height_mm`) utilized by spatial matching algorithms to determine trailer compatibility.

### 3.2 Lowboy Transport & Heavy Haulage Logistics

#### Lowboy Trailer (Heavy Low-Bed)
A specialized semi-trailer featuring a drop-deck design with a low ground clearance deck, engineered specifically to haul high-tonnage, high-profile construction equipment beneath overhead bridge clearances.

#### Step-Deck Trailer (Drop-Deck)
A commercial trailer with an upper deck and a lower deck, utilized for mid-weight machinery (e.g., 10-ton mini-excavators, compact rollers) that do not require an ultra-low lowboy well.

#### Hydraulic Folding Ramps
Heavy articulated ramps attached to the rear of a lowboy trailer, hydraulically lowered to allow tracked excavators to climb onto the trailer deck under their own traction.

#### Abnormal Wide Load
Any machinery transport configuration where overall transit width exceeds $3,000\text{ mm}$ ($3.0\text{ m}$) or height exceeds $4,500\text{ mm}$, legally requiring certified route clearances, day-light transit limits, and escort safety vehicles under Ethiopian Federal Transport Authority and KeNHA regulations.

#### Transshipment (Roadside Cargo Transfer)
The physical roadside transfer of an immovable plant machine from a disabled, broken-down lowboy trailer onto an emergency replacement trailer using specialized recovery winches and ramp shoring.

---

## 4. Operational, Field & Inspection Terminology

### 4.1 Handover Reticle, Hours & Shift Telemetry

#### Baseline Hour Meter
The certified, physically verified mechanical engine hour reading recorded and approved during the initial onboarding of a fleet asset. Represents the immutable lower bound for all future operational records.

#### Handover Inspection (`AR-01`)
The synchronized physical audit conducted at equipment delivery. Requires high-resolution camera reticle captures of the cabin mechanical hour meter, 4-point exterior walkaround photos, fuel fraction documentation, and dual-party confirmation.

#### Dual Challenge PIN Protocol
A cryptographic field verification protocol where the fleet owner’s app generates a time-bounded 6-digit challenge PIN (derived via HMAC-SHA256) which must be physically entered into the contractor’s device alongside an on-screen signature to transition a booking to `ACTIVE`.

#### Standard Daily Shift Allowance
The contractual baseline operating allowance—standardized at **8.0 Billable Engine Hours per Day**. Accumulated operational hours exceeding this quota are billed as excess operating hours.

#### Monotonic Meter Non-Regression
An invariant mathematical assertion enforced across all database triggers:
$$\text{Baseline Meter} \le \text{Handover Check-In Meter} \le \text{Return Check-Out Meter}$$
Any submission violating this condition is rejected as an unprocessable meter rollback attempt (`ERR_METER_READING_REGRESSION`).

#### Accelerometer Steady-State Validation
Client-side computer vision heuristic utilizing mobile device inertial sensors to assert that camera drift is $< 0.05g$ variance during shutter actuation, eliminating blurred or fraudulent photographs of third-party displays.

### 4.2 Damage Assessment, Wear & Insurance Loss Adjustment

#### Ordinary Operational Wear-and-Tear
The natural physical depreciation and surface abrasion resulting from routine equipment operation under standard conditions (e.g., paint chipping, track grouser wear within allowable tolerances, wire rope friction wear). Explicitly excluded from damage claims.

#### Compensable Structural Damage
Physical damage arising from impact, operator negligence, tip-overs, or operational overload (e.g., cracked boom weldments, severed hydraulic cylinders, bent track frames, smashed cab glass). Subject to formal post-rental claims or insurance settlement.

#### 24-Hour Damage Audit Window
The strict statutory window post-equipment return during which an equipment fleet owner must lodge an itemized, photo-evidenced damage claim docket (`RC-04`). Failure to submit within 24 hours results in automatic closure with zero damage liability.

#### Comprehensive Plant All-Risks (CPAR)
An insurance policy underwriting physical machinery destruction, structural rollovers, flood submersion, fire, and theft. Governed by a mandatory policy deductible and statutory 48-hour filing covenants.

#### Independent Loss Adjuster
A licensed, accredited third-party engineering expert appointed by an underwriter to conduct on-site physical failure inspections, evaluate metallurgical stress fractures, determine liability, and issue binding repair authorizations or write-off certificates.

#### Total Constructive Loss (TOL)
An insurance condition declared when the estimated cost of repairing damaged machinery exceeds $75\%$ of its certified insured market value. Results in capital write-off and permanent platform de-registration.

---

## 5. Financial, Settlement & Commission Terminology

### 5.1 Non-Custodial Payment Terms & Rails

#### Direct P2P Bank Settlement
The direct electronic wire transfer of funds between a contractor's commercial bank account and a fleet owner's verified commercial bank account, utilizing domestic banking rails:
* **Ethiopia:** Commercial Bank of Ethiopia (CBE Birr / CBE Mobile), Awash International Bank, Telebirr Enterprise.
* **Kenya:** Safaricom M-PESA Paybill, Kenya Commercial Bank (KCB) Electronic Funds Transfer.

#### Payment Evidence Slip
A digital artifact (image or PDF) uploaded by the paying contractor depicting an official bank transfer confirmation slip, detailing transaction reference code, debit account, beneficiary account, amount, and timestamp.

#### 24-Hour Payment Verification SLA
The time-bounded window during which a fleet owner must audit their corporate bank account and confirm or dispute receipt of a submitted wire prepayment. Inaction automatically triggers escalation warnings.

#### Final Balance Reconciliation Statement
The comprehensive mathematical accounting statement executed at booking termination:
$$\text{Net Balance} = (\text{Base Contract Fee} + \text{Excess Shift Surcharge} + \text{Accepted Damage}) - \text{Prepayment Previously Settled}$$

#### Prepayment Percentage
The mandatory advance wire transfer required to confirm a booking—standardized at **$50\%$ of the gross base hire value**, with the remaining balance settled upon post-rental check-out reconciliation.

### 5.2 Platform Revenue, Wallets & Debt Ledgers

#### Pre-Funded Owner Wallet (`owner_wallets`)
A virtual double-entry balance maintained for each fleet owner in PostgreSQL, tracking accumulated platform commission obligations, cancellation penalties, and manual top-up credits.

#### Approved Overdraft Credit Buffer (`credit_limit_minor`)
A monetary credit allowance extended to verified, high-reputation fleet owners permitting them to confirm bookings and incur commission debits into a negative wallet balance without immediate operational lockout.

#### Pass-Through Freight Commission Discount
A platform monetization policy discounting commission on linked lowboy transport sub-bookings to a flat **$2.5\%$ administrative fee** (compared to standard $10.0\%$ machinery commission) to prevent punitive double-taxation on mandatory mobilization freight.

#### Commission Restitution
The automatic system process crediting previously debited commission fees back to an owner's wallet when a rental is terminated due to verified third-party freight failure or administrative force majeure.

---

## 6. Technical, Architectural & Infrastructure Terminology

### 6.1 Concurrency, Locking & Spatial Indexing

#### Two-Tier Concurrency Control
The hybrid locking architecture combining an in-memory, volatile 60-minute lease in Valkey 9.1.3 (`SET hold:fleet:{id} {renter} EX 3600 NX`) with permanent storage-layer range exclusion constraints in PostgreSQL 18.6 (`EXCLUDE USING gist`).

#### GiST Temporal Exclusion Constraint
A PostgreSQL disk-level constraint utilizing generalized search trees to enforce zero overlapping rental windows (`daterange && daterange`) for any single physical fleet asset:
```sql
CONSTRAINT no_overlapping_fleet_bookings EXCLUDE USING gist (
    fleet_id WITH =,
    hire_period WITH &&
) WHERE (status IN ('HOLD_PENDING', 'PAYMENT_EVIDENCE_UPLOADED', 'CONFIRMED', 'ACTIVE', 'RETURN_INITIATED'));
```

#### Optimistic Concurrency Control (OCC)
A concurrency pattern utilizing monotonically incrementing integer `version` columns and HTTP `If-Match` ETag headers to intercept lost updates caused by parallel administrative or client mutations.

#### Valkey Distributed Lease
A time-bounded, self-evicting distributed lock in memory granting a single contractor an exclusive checkout window to visit a bank branch and transfer funds without competing hold interference.

#### Automated Standby Queue (`ZADD`)
A Valkey sorted set holding secondary contractors competing for a locked asset, sorted chronologically by request timestamp, automatically granting the next user an exclusive lease if the primary hold expires unverified.

### 6.2 Offline Synchronization, Storage & Protocols

#### Isar NoSQL Embedded Storage (`ADR-010`)
A local, ACID-compliant, native C++ NoSQL embedded database compiled into Flutter mobile applications, providing ultra-low-latency ($< 1.5\text{ ms}$) local data persistence without internet connectivity.

#### Client-Side Outbox Mutation Pattern (`SyncOutboxItem`)
An offline architectural pattern where client actions (check-ins, damage photos, meter logs) are written transactionally to a local Isar outbox table and held in a sequential queue until cellular connectivity is restored.

#### Peer-to-Peer Cryptographic QR Handshake
An offline verification protocol wherein two disconnected devices exchange an HMAC-SHA256 signed JSON payload via high-density dynamic QR codes, attesting to the transfer of equipment custodianship without cellular data packets.

#### Two-Way Differential Synchronization
The reconciliation algorithm executed upon mobile reconnection: pulling server-side database deltas since `lastServerSyncVersion` and pushing queued outbox items sequentially with client-generated idempotency keys.

### 6.3 Security, Cryptography & Identity Protocols

#### Refresh Token Rotation (RTR)
An identity security protocol wherein every refresh token exchange issues a completely fresh single-use refresh token and invalidates the previous token. Replay of an already-used token immediately revokes the entire token family.

#### Step-Up Authentication
A dynamic security challenge requiring biometric verification (Secure Enclave / Android BiometricPrompt) or master password re-entry prior to executing high-risk mutations (e.g., altering bank coordinates or drawing down enterprise credit lines).

#### Argon2id Password Hashing
The cryptographic password derivation function parameterizing memory cost ($m=65536$), time iterations ($t=3$), and parallelism ($p=4$) to prevent GPU-based offline dictionary attacks.

#### 1-Tap Account Kill-Switch
A stateless emergency URL endpoint bound to an HMAC-SHA256 token embedded in multi-channel security alerts, allowing a user whose phone was compromised to immediately freeze all account operations and listings without authenticating.

#### Problem Details Standard (RFC 9457)
The machine-readable HTTP error representation standard returning `Content-Type: application/problem+json` containing structured type URIs, titles, status codes, contextual detail, and invalid parameter arrays.

---

## 7. Acronym & Abbreviation Directory

```text
+--------------------------------------------------------------------------------------------------+
| DINK MACHINERY ACRONYM & ABBREVIATION MASTER DIRECTORY                                           |
+--------------------------------------------------------------------------------------------------+
| ACRONYM  | FULL EXPANSION                                    | DOMAIN CONTEXT                    |
+----------+---------------------------------------------------+-----------------------------------+
| ACID     | Atomicity, Consistency, Isolation, Durability     | Relational Database Transactions  |
| ADR      | Architectural Decision Record                     | Technical Governance              |
| API      | Application Programming Interface                 | Software Architecture             |
| ASVS     | Application Security Verification Standard (OWASP)| Security Engineering              |
| CAN      | Controller Area Network (CAN-Bus)                 | Vehicle IoT & Telematics          |
| CBE      | Commercial Bank of Ethiopia                       | Banking Rails (Ethiopia)          |
| CPAR     | Comprehensive Plant All-Risks Insurance           | Machinery Underwriting            |
| CRDT     | Conflict-Free Replicated Data Type                | Distributed Data Merging          |
| DLQ      | Dead-Letter Queue                                 | Message Processing Fault Handling |
| DoA      | Delegation of Authority                           | Enterprise Corporate Governance   |
| EAV      | Entity-Attribute-Value Model                      | Dynamic Schema Architecture       |
| ECA      | Ethiopian Communications Authority                | Telecommunications Regulation     |
| ECDSA    | Elliptic Curve Digital Signature Algorithm        | Cryptography & Token Minting      |
| EIC      | Ethiopian Insurance Corporation                   | Insurance Underwriting            |
| ETB      | Ethiopian Birr (ISO 4217 Currency Code)           | Financial Settlements             |
| EXIF     | Exchangeable Image File Format                    | Image Metadata Sanitation         |
| FCM      | Firebase Cloud Messaging                          | Multi-Channel Push Delivery       |
| FIPS     | Federal Information Processing Standards          | Cryptographic HSM Modules         |
| GC-1     | General Contractor Category 1                     | Ethiopian Construction Licensing  |
| GiST     | Generalized Search Tree Index                     | PostgreSQL Range Exclusions       |
| GPS      | Global Positioning System                         | Telemetry & Geofencing            |
| GVW      | Gross Vehicle Weight                              | Roadway Freight Compliance        |
| HMAC     | Hash-Based Message Authentication Code            | Token Signing & Integrity         |
| HSM      | Hardware Security Module                          | Master Key Cryptographic Storage  |
| HUD      | Heads-Up Display                                  | Real-Time Operational Client UI   |
| IEEE     | Institute of Electrical and Electronics Engineers | Engineering Specifications        |
| IFRS     | International Financial Reporting Standards       | Double-Entry Accounting Invariants |
| IoT      | Internet of Things                                | Plant Telematics Dongles          |
| IRA      | Insurance Regulatory Authority of Kenya           | Insurance Regulation (Kenya)      |
| ISO      | International Organization for Standardization   | Quality & Ergonomic Standards     |
| IVR      | Interactive Voice Response (Automated Telephony)  | Emergency Voice Alerting          |
| JWKS     | JSON Web Key Set                                  | Public Key Verification Protocols |
| JWT      | JSON Web Token (RFC 7519)                         | Stateless Authentication Tokens   |
| KCB      | Kenya Commercial Bank                             | Banking Rails (Kenya)             |
| KeNHA    | Kenya National Highways Authority                 | Axle Load Regulation (Kenya)      |
| KES      | Kenyan Shilling (ISO 4217 Currency Code)          | Financial Settlements             |
| KYC      | Know Your Customer                                | Identity & Business Attestation   |
| LWW      | Last-Write-Wins                                   | Conflict Resolution Heuristic     |
| MFA      | Multi-Factor Authentication                       | Authentication Step-Up            |
| MSISDN   | Mobile Station International Subscriber Directory  | Sovereign Phone Number Identity   |
| MUDCo    | Ministry of Urban Development & Construction (ET)  | Contractor Licensing Verification |
| NBE      | National Bank of Ethiopia                         | Central Banking Regulations       |
| NTP      | Network Time Protocol                             | Distributed Clock Sync            |
| OCC      | Optimistic Concurrency Control                    | Versioned State Collision Defense |
| OS       | Operating System                                  | Android / iOS / Linux Runtimes     |
| OTP      | One-Time Password                                 | Multi-Channel Verification Codes  |
| OWASP    | Open Web Application Security Project             | Cybersecurity Controls            |
| P2P      | Peer-to-Peer                                      | Direct Non-Custodial Transactions |
| PIN      | Personal Identification Number                    | Handover Reticle Verification     |
| PO       | Purchase Order                                    | Enterprise Procurement Accounting |
| PostGIS  | PostgreSQL Spatial Database Extender              | Geographic Vector Indexing        |
| RBAC     | Role-Based Access Control                         | Granular Permissions Engine       |
| RFC      | Request for Comments (IETF Standards)             | Wire & Protocol Compliance        |
| RTR      | Refresh Token Rotation                            | ASVS Identity Security Invariant  |
| S3       | Simple Storage Service (SeaweedFS / AWS)          | Binary Media Storage              |
| SHA      | Secure Hash Algorithm                             | Cryptographic Digest Ingestion    |
| SLA      | Service Level Agreement                           | Time-Bounded Operational Targets  |
| SMPP     | Short Message Peer-to-Peer Protocol               | Direct Carrier SMS Ingress Rails  |
| SMS      | Short Message Service                             | Telephony Alerting Channel        |
| SoD      | Separation of Duties                              | Enterprise Anti-Fraud Governance  |
| SSE      | Server-Sent Events                                | Real-Time Push In-App Streaming   |
| TIN      | Tax Identification Number                         | Fiscal Commercial Attestation     |
| TOL      | Total Constructive Loss                           | Catastrophic Insurance Write-Off  |
| TOTP     | Time-Based One-Time Password (RFC 6238)           | Hardware Authenticator MFA        |
| UI/UX    | User Interface / User Experience                  | Ergonomic Client Design Systems   |
| UUID     | Universally Unique Identifier (RFC 4122 v4)       | Primary Database Key Standard     |
| Valkey   | Valkey In-Memory Key-Value Storage (Redis Fork)    | Distributed Leases & Caching      |
| VAT      | Value Added Tax                                   | Fiscal Invoicing Requirements     |
| WBS      | Work Breakdown Structure                          | Corporate Cost Center Accounting  |
| WCAG     | Web Content Accessibility Guidelines (W3C)        | Accessibility & Contrast Controls |
+--------------------------------------------------------------------------------------------------+
```

---

*This document defines the canonical unified glossary, domain ontology, and architectural taxonomy for Dink Machinery Version 7.0. All downstream requirements modules (`03` through `14`), technical architecture specifications, wire protocol contracts, and client software implementations must strictly adhere to the definitions and invariants established herein.*
```