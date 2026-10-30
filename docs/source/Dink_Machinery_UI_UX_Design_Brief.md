# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## System Architecture, Application Flow Diagrams & UI/UX Design System Specification

**Document Identifier:** `ARCH-UIUX-DINK-CORE-V7.0`  
**Classification:** Canonical Architecture, Sequence Choreography & UX Baseline (Single Source of Truth)  
**Governing Specifications:** `SRS-DINK-CORE-V7.0` (Parts 1–5), `contracts/openapi.yaml`, `ADR-013` (OTA Typography)  
**Target Platforms:** Flutter 3.47.8 Mobile (iOS & Android), Next.js 16 Admin Web Console (Desktop & Tablet)  
**Backend & Data Runtimes:** Go 1.27.3 Modular Monolith, PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3-alpine, SeaweedFS S3, ClamAV Daemon  
**Target Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, W3C WCAG 2.2 AA / AAA (Outdoor Mode), RFC 9457 (Problem Details), ISO 4217, ISO 3166-1 alpha-2  
**Status:** Approved Unified Baseline

---

## Table of Contents

- [Part I: System Architecture & Flow Choreography](#part-i-system-architecture--flow-choreography)
  - 1. [High-Level System Context & Multi-Persona Ingress](#1-high-level-system-context--multi-persona-ingress)
  - 2. [Identity, Registration & Pluggable Verification Flows](#2-identity-registration--pluggable-verification-flows)
    - 2.1 [Passwordless Phone OTP Authentication & Refresh Rotation](#21-passwordless-phone-otp-authentication--refresh-rotation)
    - 2.2 [Pluggable Identity Verification (Fayda / National ID Adapter)](#22-pluggable-identity-verification-fayda--national-id-adapter)
    - 2.3 [Multi-Factor Authentication (Passkey WebAuthn & TOTP)](#23-multi-factor-authentication-passkey-webauthn--totp)
  - 3. [Fleet Registration, Dynamic EAV Taxonomy & Availability Scheduling](#3-fleet-registration-dynamic-eav-taxonomy--availability-scheduling)
  - 4. [Geospatial Radius Discovery, Corridor Search & Filtering (PostGIS)](#4-geospatial-radius-discovery-corridor-search--filtering-postgis)
  - 5. [End-to-End Booking Reservation, 60-Minute Valkey Hold & Dual Settlement](#5-end-to-end-booking-reservation-60-minute-valkey-hold--dual-settlement)
    - 5.1 [Manual Direct Transfer Flow (`MANUAL_TRANSFER` / Direct Transfer Model)](#51-manual-direct-transfer-flow-manual_transfer--direct-transfer-model)
    - 5.2 [Online Payment Gateway Flow (`GATEWAY` / Pluggable Adapter)](#52-online-payment-gateway-flow-gateway--pluggable-adapter)
  - 6. [The Canonical 20-State Booking Machine Choreography](#6-the-canonical-20-state-booking-machine-choreography)
  - 7. [Handover Check-In & Check-Out Inspection Workflow (Monotonic Meters & ClamAV)](#7-handover-check-in--check-out-inspection-workflow-monotonic-meters--clamav)
  - 8. [Damage Dispute Arbitration, SLA Escalation & Escrow Settlement](#8-damage-dispute-arbitration-sla-escalation--escrow-settlement)
  - 9. [Logistics Transport Sub-Booking & Bidirectional Cascade Cancellation](#9-logistics-transport-sub-booking--bidirectional-cascade-cancellation)
  - 10. [Double-Entry Financial Ledger & Wallet-Debit Commission Accounting](#10-double-entry-financial-ledger--wallet-debit-commission-accounting)
  - 11. [Operator Dispatch, 30-Second GPS Telematics & Geofence Breach Engine](#11-operator-dispatch-30-second-gps-telematics--geofence-breach-engine)
  - 12. [Server-Driven Dynamic Admin UI Engine (`/api/v1/admin/ui-config`)](#12-server-driven-dynamic-admin-ui-engine-apiv1adminui-config)
  - 13. [Dynamic Multilingual Typography & OTA Font Pipeline (`ADR-013`)](#13-dynamic-multilingual-typography--ota-font-pipeline-adr-013)
  - 14. [Offline-First Resilience, Isar Local Store & Background Sync Engine](#14-offline-first-resilience-isar-local-store--background-sync-engine)
- [Part II: Industrial Ergonomics & UI/UX Design System](#part-ii-industrial-ergonomics--uiux-design-system)
  - 15. [Executive Summary & Product Persona Context](#15-executive-summary--product-persona-context)
  - 16. [Industrial Ergonomics & Environmental Constraints](#16-industrial-ergonomics--environmental-constraints)
    - 16.1 [Direct Sunlight & High-Glare Visibility](#161-direct-sunlight--high-glare-visibility)
    - 16.2 [Gloved Hands, Vibration & Tactile Feedback](#162-gloved-hands-vibration--tactile-feedback)
    - 16.3 [One-Handed Field Operation & Thumb Zones](#163-one-handed-field-operation--thumb-zones)
  - 17. [Design System Foundations & Tokens](#17-design-system-foundations--tokens)
    - 17.1 [Color System (Standard & Outdoor High-Contrast Modes)](#171-color-system-standard--outdoor-high-contrast-modes)
    - 17.2 [Typography Engine & Multi-Script Baseline Alignment (`ADR-013`)](#172-typography-engine--multi-script-baseline-alignment-adr-013)
    - 17.3 [8dp Spatial Grid, Layout Boundaries & Elevation](#173-8dp-spatial-grid-layout-boundaries--elevation)
    - 17.4 [Component Touch Targets ($48 \times 48\text{ dp}$)](#174-component-touch-targets-48-times-48text-dp)
  - 18. [Cross-Platform Persona Journeys & Wireflows](#18-cross-platform-persona-journeys--wireflows)
    - 18.1 [Equipment Renter Experience (Discovery to Handover)](#181-equipment-renter-experience-discovery-to-handover)
    - 18.2 [Fleet Owner Experience (Fleet Registry, Availability & Wallet)](#182-fleet-owner-experience-fleet-registry-availability--wallet)
    - 18.3 [Field Operator Experience (Inspection Stepper & GPS Telematics)](#183-field-operator-experience-inspection-stepper--gps-telematics)
    - 18.4 [Operations & Compliance Admin Console (Next.js 16 Web)](#184-operations--compliance-admin-console-nextjs-16-web)
  - 19. [Core Component Library & Interaction Specifications](#19-core-component-library--interaction-specifications)
    - 19.1 [Dynamic EAV Form Accordion Engine](#191-dynamic-eav-form-accordion-engine)
    - 19.2 [Camera-Enforced Telematics Capture Overlay](#192-camera-enforced-telematics-capture-overlay)
    - 19.3 [Countdown Timer & Bank Transfer Card](#193-countdown-timer--bank-transfer-card)
    - 19.4 [Interactive Handover Checklist & Glass Signature Canvas](#194-interactive-handover-checklist--glass-signature-canvas)
    - 19.5 [Split-Screen Dispute Arbitration Lightbox](#195-split-screen-dispute-arbitration-lightbox)
  - 20. [Motion, Micro-Interactions & Audio-Haptic Feedback](#20-motion-micro-interactions--audio-haptic-feedback)
  - 21. [Offline Cues, Empty States & RFC 9457 Error Presentation](#21-offline-cues-empty-states--rfc-9457-error-presentation)
    - 21.1 [Persistent Offline Status Banner](#211-persistent-offline-status-banner)
    - 21.2 [RFC 9457 Problem Details Presentation](#212-rfc-9457-problem-details-presentation)

---

# Part I: System Architecture & Flow Choreography

## 1. High-Level System Context & Multi-Persona Ingress

The Dink Machinery platform coordinates three mobile client personas (Renters, Fleet Owners, Field Operators) and desktop administrative workstations through Traefik 3.x edge ingress, routing into a single compiled Go 1.27.3 modular monolith backend.

```mermaid
flowchart TD
    subgraph EdgeClients["Edge Clients & Personas"]
        Renter["Equipment Renter<br>(Flutter Mobile Client)"]
        Owner["Fleet Owner<br>(Flutter Mobile Client)"]
        Operator["Field Operator<br>(Flutter Mobile Client)"]
        Admin["Compliance & Ops Admin<br>(Next.js 16 Web Console)"]
    end

    subgraph EdgeIngress["Edge Gateway & Hardening Layer"]
        Traefik["Traefik 3.x API Gateway<br>(TLS 1.3 Termination, HSTS, Rate Limiter)"]
        WAF["OWASP CRS / DDoS Shield"]
    end

    subgraph CoreBackend["Go 1.27.3 Modular Monolith (Clean Architecture)"]
        AuthMod["Auth & Identity Module<br>(ES256 JWT, Valkey Sessions)"]
        FleetMod["Machinery & Catalog Module<br>(Dynamic EAV Taxonomy)"]
        GeoMod["Spatial Search Engine<br>(PostGIS 3.6.4 Bounded Queries)"]
        BookingMod["Booking State Engine<br>(20-State Machine, OCC Versioning)"]
        PaymentMod["Escrow Ledger & Wallet Engine<br>(Double-Entry, BK-17-BE)"]
        InspectMod["Inspection & Media Ingress<br>(ClamAV Stream Scan)"]
        AdminUIMod["Server-Driven UI Engine<br>(/api/v1/admin/ui-config)"]
    end

    subgraph PersistenceInfra["Storage & State Infrastructure"]
        Postgres[("PostgreSQL 18.6 + PostGIS 3.6.4<br>(Spatial Points, JSONB Specs, pgxpool)")]
        Valkey[("Valkey 9.1.3 Cache<br>(60-Min Mutex Holds, Rate-Limit Counters)")]
        SeaweedFS[("SeaweedFS S3 Object Storage<br>(Photos, Licenses, Encrypted PDFs)")]
        ClamAV["ClamAV Daemon<br>(tcp://clamav:3310 zINSTREAM)"]
    end

    Renter -->|HTTPS / TLS 1.3| Traefik
    Owner -->|HTTPS / TLS 1.3| Traefik
    Operator -->|WSS / HTTPS| Traefik
    Admin -->|HTTPS / TLS 1.3| Traefik

    Traefik --> WAF
    WAF --> CoreBackend

    CoreBackend -->|Binary Protocol pgx/v5| Postgres
    CoreBackend -->|TCP RESP Protocol| Valkey
    CoreBackend -->|S3 REST API| SeaweedFS
    InspectMod -->|TCP Stream| ClamAV
```

---

## 2. Identity, Registration & Pluggable Verification Flows

### 2.1 Passwordless Phone OTP Authentication & Refresh Rotation

Authentication follows OWASP ASVS 4.0.3 Level 3 with single-use rotating refresh tokens stored in Valkey and cryptographic ES256 session JWTs.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Flutter Client (Renter/Owner/Operator)
    participant GW as Traefik API Gateway
    participant Auth as Go Auth Module (net/http)
    participant VK as Valkey 9.1.3 Cache
    participant DB as PostgreSQL 18.6 (users)
    participant SMS as SMS Gateway (Ethio/Safaricom)

    Client->>GW: POST /api/v1/auth/otp/challenge { phone: "+251911223344" }
    GW->>Auth: Validate E.164 (+251 / +254)
    Auth->>VK: Check Rate Limit (max 5 OTPs/hr per IP & phone)
    alt Rate Limit Exceeded
        VK-->>Client: HTTP 429 Too Many Requests (RFC 9457)
    else Under Rate Limit
        Auth->>Auth: Generate 6-digit crypto OTP (crypto/rand)
        Auth->>VK: SETEX otp:challenge:+251911223344 120s { code_hash, attempts: 0 }
        Auth->>SMS: Dispatch SMS Challenge
        SMS-->>Client: SMS with 6-digit code received
        Auth-->>Client: HTTP 200 OK { challengeId, expiresIn: 120 }
    end

    Client->>GW: POST /api/v1/auth/otp/verify { phone: "+251911223344", code: "849201" }
    GW->>Auth: Forward verification payload
    Auth->>VK: GET otp:challenge:+251911223344
    alt Invalid Code / Attempts > 5
        Auth->>VK: INCR attempts
        Auth-->>Client: HTTP 400 Bad Request / 429 Locked (15 min lockout)
    else Valid Code
        Auth->>VK: DEL otp:challenge:+251911223344
        Auth->>DB: Upsert User Record & Resolve RBAC Roles
        Auth->>Auth: Sign Access JWT (ES256, 15-min TTL)
        Auth->>VK: SETEX refresh:{token_id} 30d { user_id, family_id }
        Auth-->>Client: HTTP 200 OK { accessToken, refreshToken, userProfile }
        Client->>Client: Store tokens in flutter_secure_storage (Keychain/Keystore)
    end
```

### 2.2 Pluggable Identity Verification (Fayda / National ID Adapter)

The backend provides a pluggable adapter framework (`REG-IDP-01-BE`). When active, incomplete verification blocks mutating transactional endpoints (`AUTH-05-BE`).

```mermaid
sequenceDiagram
    autonumber
    actor User as Authenticated User
    participant App as Mobile Client
    participant API as Auth Identity Interactor
    participant Policy as Policy Engine (REG-IDP-02-BE)
    participant Adapter as IdentityVerificationProvider (Fayda)
    participant Registry as National Fayda API
    participant DB as PostgreSQL 18.6 (identity_verifications)

    User->>App: Submits Identity ID (e.g. Fayda FIN) + Consent Checkbox
    App->>API: POST /api/v1/user/identity/verify { providerKey: "fayda", identifier: "...", consentAccepted: true }
    API->>DB: Persist Consent Audit Record (timestamp, policy_version, request_ref)
    API->>Adapter: Dispatch Verification Challenge
    Adapter->>Registry: HTTPS REST /verify { identifier, consentRef }

    alt Identity Match
        Registry-->>Adapter: 200 OK { verified: true, name: "...", dob: "1992-05-12", photoUrl: "..." }
        Adapter-->>API: ProviderResult { status: VERIFIED, verifiedAttributes }
        API->>API: Age Gate Verification (DOB >= 18 years)
        API->>DB: INSERT INTO identity_verifications (user_id, provider_key, status='VERIFIED', verified_attributes)
        API->>Policy: Re-evaluate Policy (ALL / ANY / N_OF_M)
        Policy-->>API: { satisfied: true, identity_status: VERIFIED }
        API-->>App: HTTP 200 OK { identityStatus: "VERIFIED" }
        App->>App: Dismiss Identity Lock Overlay (AUTH-05-MOB)
    else Mismatch / Not Found
        Registry-->>Adapter: 404 / 422 Identity Mismatch
        Adapter-->>API: ProviderResult { status: UNVERIFIED, errorCode: IDENTITY_MISMATCH }
        API->>DB: Record Failed Attempt (rate-limit 3 per 24h)
        API-->>App: HTTP 422 Unprocessable Entity { code: "IDENTITY_MISMATCH", remainingAttempts: 2 }
    else Registry 5xx / Timeout (Failover)
        Adapter-->>API: ProviderResult { status: PENDING_REVIEW, error: "TIMEOUT" }
        API->>DB: INSERT INTO identity_verifications (status='PENDING_REVIEW')
        API->>DB: Enqueue in Admin Manual Verification Queue (AD-13-BE)
        API-->>App: HTTP 202 Accepted { identityStatus: "PENDING_REVIEW", message: "Under manual review" }
    end
```

### 2.3 Multi-Factor Authentication (Passkey WebAuthn & TOTP)

Enforces mandatory TOTP setup for administrative staff (`AUTH-04-WEB`) while offering progressive Passkey and Authenticator App options for standard personas.

```mermaid
flowchart TD
    Start([User Initiates Authentication]) --> PasswordCheck{Password / OTP Valid?}
    PasswordCheck -- No --> LockoutCheck{Failed Attempts >= 10?}
    LockoutCheck -- Yes --> Lockout[HTTP 429 Account Locked 30 Min]
    LockoutCheck -- No --> Retry[Return Invalid Credentials]

    PasswordCheck -- Yes --> RoleCheck{Is Admin / Super Admin?}
    RoleCheck -- Yes --> ForceMFA[Mandatory TOTP Verification Required]
    RoleCheck -- No --> UserMFAConfig{User mfa_enabled == true?}

    UserMFAConfig -- No --> IssueSession[Issue Full ES256 Session JWT]
    UserMFAConfig -- Yes --> ResolvePriority[Resolve Factor Priority List: Passkey > TOTP > SMS]

    ForceMFA --> TOTPChallenge[Display TOTP Code Prompt]
    ResolvePriority --> PresentChallenge[Present Primary Registered Factor Challenge]

    PresentChallenge --> SwitchMethod{User Taps 'Try another method'?}
    SwitchMethod -- Yes --> FactorSheet[Display Alternate Registered Factors]
    FactorSheet --> PresentChallenge

    SwitchMethod -- No --> VerifyFactor{Verify Cryptographic Assertion / Code}
    TOTPChallenge --> VerifyFactor

    VerifyFactor -- Success --> RegressionCheck{Passkey Counter Regression Check}
    RegressionCheck -- Cloned Detected --> RejectAuth[HTTP 401 Cloned Authenticator Blocked]
    RegressionCheck -- Valid --> IssueSession
    VerifyFactor -- Fail --> CountMFA[Increment Failed MFA Counter]
```

---

## 3. Fleet Registration, Dynamic EAV Taxonomy & Availability Scheduling

Fleet asset registration validates technical attributes dynamically using JSON Schema definitions stored in `spec_field_definitions` for the chosen category.

```mermaid
sequenceDiagram
    autonumber
    actor Owner as Fleet Owner
    participant App as Mobile Client
    participant CatAPI as Category & Spec Module
    participant MediaAPI as Media & Ingress Controller
    participant ClamAV as ClamAV Antivirus Daemon
    participant S3 as SeaweedFS Object Storage
    participant DB as PostgreSQL 18.6 (fleets)
    participant Admin as Next.js Admin Queue

    Owner->>App: Selects Machinery Category (e.g. Excavation > Crawler Excavator)
    App->>CatAPI: GET /api/v1/categories/{id}/spec-templates
    CatAPI-->>App: SpecDefinitionJSON { fields: [operating_weight, bucket_capacity, boom_length] }
    App->>App: Render Dynamic Accordion Form Inputs (FL-07-MOB)

    Owner->>App: Captures 5 Required Photos + Hour Meter Photo
    App->>MediaAPI: POST /api/v1/media/upload (Multipart image stream)
    MediaAPI->>ClamAV: TCP Stream zINSTREAM (tcp://clamav:3310)
    alt Virus Detected
        ClamAV-->>MediaAPI: VIRUS DETECTED (EICAR / Malware)
        MediaAPI-->>App: HTTP 422 Unprocessable Entity { code: "INSPECT_MEDIA_SCAN_FAILED" }
    else Clean File
        ClamAV-->>MediaAPI: OK
        MediaAPI->>S3: PUT /fleets/media/{uuid}.webp (Compressed WebP, max 1600px)
        MediaAPI-->>App: HTTP 201 Created { mediaUrl: "https://s3.../fleet_01.webp" }
    end

    Owner->>App: Enters Pricing (Minor Units), Availability Dates, Operator Assigned
    App->>CatAPI: POST /api/v1/fleets (Base attributes + Dynamic EAV specs JSON)
    CatAPI->>CatAPI: Validate Integer Minor-Unit Daily Rate (daily_rate_amount_minor > 0)
    CatAPI->>CatAPI: Validate Mandatory Category Specs (FL-09-BE)
    CatAPI->>DB: INSERT INTO fleets (status='PENDING_APPROVAL', version=1, location_point=ST_SetSRID(...))
    CatAPI->>DB: INSERT INTO fleet_spec_indexed (Extract logistics weight & dimensions)
    CatAPI->>DB: INSERT INTO fleet_availability (Blocked calendar intervals)
    CatAPI-->>App: HTTP 201 Created { fleetId, status: "PENDING_APPROVAL" }
    CatAPI->>Admin: Emits WebSocket Event 'fleet.submitted_for_approval'
```

---

## 4. Geospatial Radius Discovery, Corridor Search & Filtering (PostGIS)

Dink Machinery executes high-throughput spheroidal distance queries over PostGIS GiST indexes, ensuring queries complete within $\le 80\text{ ms}$ at P95 across 1,000,000 indexed assets.

```mermaid
flowchart TD
    Client["Renter Mobile Client<br>(GPS Coordinates: Lat, Lng)"] -->|GET /api/v1/geo/fleets/radius-search| Gateway[Traefik 3.x Edge Gateway]
    Gateway --> GeoEngine[Geospatial Query Interactor]

    subgraph QueryParser["Parameter Parsing & Bounds Enforcement"]
        GeoEngine --> BoundingCheck{"Radius Bounds Check<br>(1km <= radiusKm <= 500km)"}
        BoundingCheck -- Exceeded --> Err400[HTTP 400 Bad Request]
        BoundingCheck -- Valid --> CountryContext["Resolve Country Context<br>(country_configs: ET or KE)"]
    end

    subgraph PostGISExecution["PostgreSQL 18.6 + PostGIS 3.6.4 Engine"]
        CountryContext --> SQLQuery["SELECT * FROM fleets f<br>WHERE f.status = 'ACTIVE'<br>AND f.deleted_at IS NULL<br>AND ST_DWithin(f.location_point::geography,<br>    ST_SetSRID(ST_MakePoint(lng, lat), 4326)::geography,<br>    radius_meters)"]
        SQLQuery --> SpatialIndex["idx_fleets_location_point (GiST Index Scan)"]
        SpatialIndex --> SpecFilter["Filter JSONB Dynamic Specs (GIN Index Scan)"]
        SpecFilter --> AvailabilitySubquery["Filter Overlapping Bookings<br>(NOT EXISTS in active booking calendar)"]
    end

    AvailabilitySubquery --> RankingEngine[Relevance & Sponsored Rank Engine]

    subgraph RankingLogic["Weighted Ordering Pipeline"]
        RankingEngine --> BoostCheck{"IsBoosted == true?<br>(Active Sponsored Placement)"}
        BoostCheck -- Yes --> FloatTop["Float Tier to Top (Gold > Silver > Bronze)"]
        BoostCheck -- No --> OrganicScore["Compute Organic Score:<br>0.35*Distance + 0.25*Rating + 0.20*Recency + 0.20*Completeness"]
        FloatTop --> SortOutput[Order By: Tier DESC, Score DESC]
        OrganicScore --> SortOutput
    end

    SortOutput --> KeysetPagination["Keyset Cursor Pagination (id, created_at)"]
    KeysetPagination --> PayloadEnvelope["Wrap in SuccessEnvelope JSON"]
    PayloadEnvelope --> Client
```

---

## 5. End-to-End Booking Reservation, 60-Minute Valkey Hold & Dual Settlement

### 5.1 Manual Direct Transfer Flow (`MANUAL_TRANSFER` / Direct Transfer Model)

The primary MVP settlement rail where the renter pays the fleet owner directly via commercial bank transfer. The platform holds zero customer escrow funds, collecting marketplace commission via owner wallet debit upon confirmation (`BK-17-BE`).

```mermaid
sequenceDiagram
    autonumber
    actor Renter as Equipment Renter
    actor Owner as Fleet Owner
    participant App as Mobile Client
    participant BookingAPI as Booking State Interactor
    participant VK as Valkey 9.1.3 Cache
    participant DB as PostgreSQL 18.6 (pgx.Tx)
    participant S3 as SeaweedFS Object Storage

    Renter->>App: Taps "Book Now" (Dates, Attachments, Delivery Site)
    App->>BookingAPI: POST /api/v1/bookings (Idempotency-Key: UUIDv4)
    BookingAPI->>VK: SET lock:hold:fleet:{id} {renterId} NX EX 3600
    alt Lock Acquisition Failed (Already Held)
        VK-->>BookingAPI: Lock Failure (Key exists)
        BookingAPI-->>App: HTTP 409 Conflict (FLEET_ALREADY_RESERVED)
    else Lock Acquired
        BookingAPI->>DB: INSERT INTO bookings (status='HOLD_PENDING', hold_expires_at=NOW()+1h, version=1)
        BookingAPI-->>App: HTTP 201 Created { bookingId, status: "HOLD_PENDING", remainingSeconds: 3600 }
    end

    App->>App: Render 60-Minute Circular Countdown Timer (BK-06-MOB)
    Renter->>Renter: Executes Bank Transfer (CBE / Awash / Telebirr) to Owner Account
    Renter->>App: Uploads Transfer Receipt Screenshot
    App->>BookingAPI: POST /api/v1/bookings/{id}/payment-proof (Multipart Image)
    BookingAPI->>S3: Stream through ClamAV & Persist to S3
    BookingAPI->>DB: UPDATE bookings SET status='PAYMENT_EVIDENCE_UPLOADED', version=version+1
    BookingAPI-->>App: HTTP 200 OK (Status: PAYMENT_EVIDENCE_UPLOADED)
    BookingAPI->>Owner: Push Notification 'Payment slip uploaded for review'

    Owner->>App: Inspects Receipt Image & Transaction Reference
    Owner->>BookingAPI: PUT /api/v1/bookings/{id}/confirm-payment (Idempotency-Key)
    BookingAPI->>DB: BEGIN TRANSACTION (pgx.Tx)
    Note over BookingAPI,DB: Atomic Ledger & Balance Mutation (BK-17-BE)
    BookingAPI->>DB: SELECT * FROM owner_wallets WHERE owner_id = $1 FOR UPDATE
    BookingAPI->>DB: Compute Commission: (RentalFee * CommissionBps) / 10000
    BookingAPI->>DB: INSERT INTO owner_wallet_transactions (type='COMMISSION_DEBIT', amount=-feeMinor)
    BookingAPI->>DB: UPDATE owner_wallets SET balance_minor = balance_minor - feeMinor
    BookingAPI->>DB: UPDATE bookings SET status='CONFIRMED', version=version+1 WHERE id=$id AND version=$v
    BookingAPI->>DB: COMMIT TRANSACTION
    BookingAPI->>VK: DEL lock:hold:fleet:{id} (Hold successfully converted)
    BookingAPI->>Renter: Push Notification 'Booking Confirmed! Machine in Mobilization'
    BookingAPI-->>Owner: HTTP 200 OK { status: "CONFIRMED" }
```

### 5.2 Online Payment Gateway Flow (`GATEWAY` / Pluggable Adapter)

When pluggable online gateways (Telebirr, MPESA) are activated, state machine transitions occur automatically via cryptographic webhook callbacks.

```mermaid
sequenceDiagram
    autonumber
    actor Renter as Equipment Renter
    participant App as Mobile Client
    participant PaymentAPI as Payment Engine
    participant Adapter as PaymentGatewayAdapter (Telebirr/MPESA)
    participant Gateway as Telco Payment Gateway
    participant DB as PostgreSQL 18.6 (pgx.Tx)

    Renter->>App: Selects Payment Method: "Telebirr" / "M-Pesa"
    App->>PaymentAPI: POST /api/v1/payments/gateway/initiate { bookingId, gatewayKey: "telebirr" }
    PaymentAPI->>Adapter: initiate(paymentLeg, amountMinor, currency)
    Adapter->>Gateway: POST /checkout/session (Merchant Credentials, CallbackURL)
    Gateway-->>Adapter: { checkoutUrl, gatewayTxnRef }
    Adapter-->>PaymentAPI: GatewayInitResponse
    PaymentAPI->>DB: INSERT INTO gateway_transactions (status='INITIATED', gateway_txn_ref=...)
    PaymentAPI-->>App: HTTP 200 OK { checkoutUrl, ussdPayload }

    App->>Gateway: Launches USSD Push / In-App Telco Deeplink
    Renter->>Gateway: Authorizes Payment via Mobile Wallet PIN
    Gateway->>PaymentAPI: POST /api/v1/payments/callback/telebirr (HMAC-SHA256 Signed Payload)
    
    PaymentAPI->>Adapter: verifyCallback(rawPayload, signatureHeader)
    alt Signature Invalid / Stale Timestamp (>300s)
        PaymentAPI-->>Gateway: HTTP 401 Unauthorized
    else Signature Verified & Amount Matches
        PaymentAPI->>DB: BEGIN TRANSACTION (pgx.Tx)
        PaymentAPI->>DB: UPDATE gateway_transactions SET status='SUCCEEDED' WHERE gateway_txn_ref=...
        PaymentAPI->>DB: INSERT INTO financial_ledger_entries (Double-Entry Legs Balanced to 0)
        PaymentAPI->>DB: UPDATE bookings SET status='CONFIRMED', version=version+1
        PaymentAPI->>DB: COMMIT TRANSACTION
        PaymentAPI-->>Gateway: HTTP 200 OK { received: true }
        PaymentAPI->>App: WebSocket Broadcast 'booking.confirmed'
        App->>App: Navigate to "Booking Confirmed" Screen
    end
```

---

## 6. The Canonical 20-State Booking Machine Choreography

The rental booking machine enforces 20 deterministic states, strictly governed by optimistic concurrency control (`WHERE id = $1 AND version = $2`). Unmapped transitions abort with RFC 9457 error `BOOKING_INVALID_STATE_TRANSITION`.

```mermaid
stateDiagram-v2
    [*] --> DRAFT: Renter selects parameters
    DRAFT --> HOLD_PENDING: Request Hold (Idempotency Key)
    
    HOLD_PENDING --> HOLD_EXPIRED: 60-Min Timer Elapses
    HOLD_PENDING --> CANCELLED_BY_RENTER: Renter Aborts (Free Cancel)
    HOLD_PENDING --> PAYMENT_EVIDENCE_UPLOADED: Bank Slip Uploaded (Manual)
    HOLD_PENDING --> PAYMENT_PROCESSING: Online Gateway Initiated
    
    PAYMENT_PROCESSING --> CONFIRMED: Gateway Webhook Verified
    PAYMENT_PROCESSING --> PAYMENT_FAILED: Telco Rail Declined
    PAYMENT_FAILED --> HOLD_PENDING: Retry Payment Rail
    PAYMENT_FAILED --> HOLD_EXPIRED: Window Closes
    
    PAYMENT_EVIDENCE_UPLOADED --> CONFIRMED: Owner Approves Slip (Commission Debited)
    PAYMENT_EVIDENCE_UPLOADED --> PAYMENT_DECLINED: Owner Declines Slip
    
    PAYMENT_DECLINED --> PAYMENT_EVIDENCE_UPLOADED: Renter Re-uploads Slip
    PAYMENT_DECLINED --> CANCELLED: Renter Cancels (Free Cancel)
    PAYMENT_DECLINED --> DISPUTE_OPENED: Renter Disputes Decline
    
    DISPUTE_OPENED --> CONFIRMED: Admin Rules Slip Valid
    DISPUTE_OPENED --> CANCELLED: Admin Rules Payment Invalid
    
    CONFIRMED --> AWAITING_FLEET_READINESS: Owner Reports Mobilization Delay
    CONFIRMED --> READY_FOR_PICKUP: Equipment Staged at Site/Depot
    CONFIRMED --> CANCELLED: Cancel with Timeline Penalty
    
    AWAITING_FLEET_READINESS --> READY_FOR_PICKUP: Owner Marks Ready
    AWAITING_FLEET_READINESS --> PENDING_ADDITIONAL_PAYMENT: Extra Mobilization Fee Asked
    AWAITING_FLEET_READINESS --> CANCELLED: ETA Exceeded (Free Renter Refund)
    
    PENDING_ADDITIONAL_PAYMENT --> READY_FOR_PICKUP: Renter Pays & Owner Confirms
    PENDING_ADDITIONAL_PAYMENT --> CANCELLED: Renter Declines (Owner Penalty)
    
    READY_FOR_PICKUP --> ACTIVE: Joint Pre-Inspection & Sign-off
    READY_FOR_PICKUP --> DISPUTED_HOUR_METER: Pre-Rental Meter Conflict
    READY_FOR_PICKUP --> CANCELLED: 48h Handover Timeout
    
    DISPUTED_HOUR_METER --> READY_FOR_PICKUP: Admin Resolves Pre-Rental Reading
    
    ACTIVE --> EXTENSION_REQUESTED: Renter Requests Extra Days
    EXTENSION_REQUESTED --> ACTIVE: Owner Accepts & Debits Extra Comm
    EXTENSION_REQUESTED --> RETURN_INITIATED: Extension Denied / Expired
    
    ACTIVE --> RETURN_INITIATED: Scheduled End Date Reached / Return Triggered
    
    RETURN_INITIATED --> DISPUTED_HOUR_METER: Return Meter Conflict Flagged
    RETURN_INITIATED --> DAMAGE_PARTIAL_RESOLVED: Return Damage Claim Flagged
    RETURN_INITIATED --> PENDING_FINAL_PAYMENT: Excess Hours / Accepted Damage > 0
    RETURN_INITIATED --> COMPLETED: Zero Extra Charges (Inspection Clean)
    
    DISPUTED_HOUR_METER --> PENDING_FINAL_PAYMENT: Admin Resolves Reading (Due > 0)
    DISPUTED_HOUR_METER --> COMPLETED: Admin Resolves Reading (Due = 0)
    
    DAMAGE_PARTIAL_RESOLVED --> PENDING_FINAL_PAYMENT: Damage Sum Added to Invoice
    
    PENDING_FINAL_PAYMENT --> COMPLETED: Final Payment Verified
    
    COMPLETED --> RATED: Both Parties Submit Reviews (or 7-Day Cron)
    RATED --> CLOSED: Final Ledger Audit Reconciled
    
    HOLD_EXPIRED --> [*]
    CANCELLED_BY_RENTER --> [*]
    CANCELLED --> [*]
    CLOSED --> [*]
```

---

## 7. Handover Check-In & Check-Out Inspection Workflow (Monotonic Meters & ClamAV)

Physical handovers require joint digital inspection enforcing monotonic hour meter progression and antivirus scanning.

```mermaid
sequenceDiagram
    autonumber
    actor Renter as Equipment Renter
    actor Operator as Field Operator / Owner
    participant App as Mobile Inspection Stepper
    participant API as Inspection Module
    participant ClamAV as ClamAV Daemon (zINSTREAM)
    participant S3 as SeaweedFS S3
    participant DB as PostgreSQL 18.6 (inspections)

    Note over Renter,Operator: Equipment Arrives at Construction Job Site
    Operator->>App: Launches Pre-Rental Inspection Stepper
    Operator->>App: Takes Photo of Physical Hour Meter (Camera-Enforced)
    Operator->>App: Inputs Numeric Reading: e.g. 1,450.50 Hours
    Operator->>App: Selects Fuel Gauge Level: 85%
    Operator->>App: Captures 4 Mandatory Angles (Front, Rear, Left Track, Right Track)

    App->>API: POST /api/v1/inspections/media/stream (Chunked Upload)
    API->>ClamAV: Scan Raw Stream
    ClamAV-->>API: Stream Clean (OK)
    API->>S3: Persist WebP Media Artifacts
    API-->>App: { mediaUrls: [...] }

    Operator->>App: Signs Handover Attestation on Glass (Touch Canvas)
    Renter->>App: Reviews Recorded Meter, Fuel, and Photos on Own Device
    Renter->>App: Counter-signs Attestation / Validates via Phone OTP
    App->>API: POST /api/v1/bookings/{id}/handover/attest { meterReading: 1450.5, fuelPercent: 85, signatureUrl: "..." }
    
    API->>DB: INSERT INTO booking_inspections (type='PRE_RENTAL', meter_reading=1450.5, fuel=85)
    API->>DB: UPDATE bookings SET status='ACTIVE', started_at=NOW(), version=version+1
    API-->>App: HTTP 200 OK (Status: ACTIVE)

    Note over Renter,Operator: Machinery Active on Job Site (Rental Concludes)

    Operator->>App: Launches Post-Rental Return Inspection Stepper
    Operator->>App: Enters Return Hour Meter: e.g. 1,498.50 Hours
    App->>API: POST /api/v1/bookings/{id}/inspection/post { meterReading: 1498.5 }
    
    API->>DB: SELECT meter_reading FROM booking_inspections WHERE booking_id=$id AND type='PRE_RENTAL'
    alt Monotonic Meter Regression (Return < Pre-Rental)
        API-->>App: HTTP 422 Unprocessable Entity (RFC 9457: INSPECT_METER_REGRESSION)
    else Reading Valid (1498.50 >= 1450.50)
        API->>DB: Compute Billable Hours: (1498.5 - 1450.5) = 48.0 Hours
        API->>DB: Compare against Agreed Limit: (6 Days * 8h/Day = 48h)
        Note over API,DB: Overtime Hours = 0. Clean Return!
        API->>DB: UPDATE bookings SET status='COMPLETED', version=version+1
        API-->>App: HTTP 200 OK (Status: COMPLETED)
    end
```

---

## 8. Damage Dispute Arbitration, SLA Escalation & Escrow Settlement

Disputed claims (damage, hour meter discrepancies) are isolated into a first-class `Dispute` domain with strict 3-day SLA timers and dual administrative communication channels.

```mermaid
sequenceDiagram
    autonumber
    actor Owner as Fleet Owner
    actor Renter as Equipment Renter
    actor Admin as Dispute Arbitration Officer
    participant API as Dispute Module
    participant DB as PostgreSQL 18.6 (pgx.Tx)
    participant Cron as SLA Escalation Daemon

    Owner->>API: POST /api/v1/bookings/{id}/damage { amountMinor: 4500000, description: "Cracked boom hydraulic cylinder" }
    API->>DB: INSERT INTO damage_claims (amount_minor=4500000, status='CLAIMED')
    API->>Renter: Push Alert 'Damage Claim Filed: 45,000.00 ETB'

    Renter->>API: PUT /api/v1/bookings/{id}/damage/respond { action: "DISPUTE", rebuttalNotes: "Damage existed prior" }
    API->>DB: INSERT INTO disputes (booking_id, origin='return_damage', status='OPEN', dispute_deadline=NOW()+3days)
    API->>DB: UPDATE bookings SET status='DISPUTED', version=version+1
    
    Cron->>DB: Nightly Scan: SELECT * FROM disputes WHERE status='OPEN' AND dispute_deadline < NOW()
    opt Overdue SLA Breach
        Cron->>Admin: High-Priority PagerDuty / Web Alert 'SLA Overdue: Booking bk-1029'
    end

    Admin->>API: GET /api/v1/admin/disputes/{id}
    API-->>Admin: Consolidates Pre/Post Photos, Damage Lightbox Diffs, Tri-Party Chat
    Admin->>API: POST /api/v1/disputes/{id}/messages { scope: "owner-only", text: "Provide certified repair quote" }
    
    Admin->>API: POST /api/v1/admin/disputes/{id}/rule { ruling: "PARTIAL_SPLIT", resolvedAmountMinor: 2500000 }
    API->>DB: BEGIN TRANSACTION (pgx.Tx)
    API->>DB: UPDATE disputes SET status='RESOLVED_PARTIAL_SPLIT', resolved_amount_minor=2500000
    API->>DB: UPDATE bookings SET status='DAMAGE_PARTIAL_RESOLVED', version=version+1
    API->>DB: Append 25,000.00 ETB to Booking Final Settlement Invoice
    API->>DB: UPDATE bookings SET status='PENDING_FINAL_PAYMENT'
    API->>DB: COMMIT TRANSACTION
    API->>Renter: Push Alert 'Dispute Ruled: 25,000.00 ETB Added to Final Balance'
    API->>Owner: Push Alert 'Dispute Ruled: 25,000.00 ETB Approved for Damage'
```

---

## 9. Logistics Transport Sub-Booking & Bidirectional Cascade Cancellation

Logistics transport vehicles (lowboys, flatbeds) are booked as formal sub-bookings linked via `parent_booking_id`. Cancellation cascades preserve operational balance across carriers.

```mermaid
sequenceDiagram
    autonumber
    actor Renter as General Contractor
    actor MachineOwner as Machinery Owner
    actor CarrierOwner as Logistics Carrier Owner
    participant API as Logistics & Booking Engine
    participant DB as PostgreSQL 18.6 (bookings)

    Renter->>API: Books 30-Ton Crawler Excavator (Confirmed Booking #101)
    MachineOwner->>API: GET /api/v1/bookings/101/logistics-suggestions
    API->>DB: Query fleets WHERE category='LOGISTICS' AND payload_capacity >= 30000kg
    API-->>MachineOwner: Returns Compatible Carriers (Lowboy 40-Ton #202)
    
    MachineOwner->>API: POST /api/v1/bookings (Create Transport Sub-Booking)
    API->>DB: INSERT INTO bookings (id=505, parent_booking_id=101, fleet_id=202, status='HOLD_PENDING')
    MachineOwner->>API: Uploads Transport Prepayment Slip
    CarrierOwner->>API: Confirms Transport Payment -> Sub-Booking #505 Transitions to 'CONFIRMED'

    alt Main Machinery Booking Cancelled by Renter
        Renter->>API: Cancels Machinery Booking #101
        API->>DB: UPDATE bookings SET status='CANCELLED' WHERE id=101
        Note over API,DB: Automatic Downward Cascade Cancellation
        API->>DB: UPDATE bookings SET status='CANCELLED' WHERE parent_booking_id=101
        API->>CarrierOwner: Alert 'Transport Sub-Booking Cancelled due to Main Lease Cancellation'
    else Carrier Cancels Transport Sub-Booking #505
        CarrierOwner->>API: Cancels Sub-Booking #505
        API->>DB: UPDATE bookings SET status='CANCELLED' WHERE id=505
        Note over API,DB: Upward Viability Alert on Main Lease
        API->>DB: UPDATE bookings SET logistics_failed=TRUE WHERE id=101
        API->>MachineOwner: High-Priority Alert 'Carrier Cancelled Transport! Re-book or cancel without penalty'
        opt Machine Owner Fails to Resolve Within min(24h, StartTime)
            API->>DB: Auto-cancel Booking #101 with ZERO Penalty to Machine Owner
            API->>Renter: Full Prepayment Refund Obligation Issued
        end
    end
```

---

## 10. Double-Entry Financial Ledger & Wallet-Debit Commission Accounting

Under requirement `BK-17-BE`, every mutating financial action writes immutable, balanced credit and debit rows into `financial_ledger_entries`. The sum of all legs in a transaction must equal exactly zero minor units.

```mermaid
flowchart LR
    subgraph TriggerEvent["Trigger Event: Booking Confirmed (BK-17-BE)"]
        Confirm[Booking Confirmed Event<br>Rental Base: 350,000.00 ETB<br>Commission: 10% = 35,000.00 ETB]
    end

    subgraph AtomicTransaction["Single Atomic DB Transaction (pgx.Tx)"]
        direction TB
        LockWallet["SELECT * FROM owner_wallets<br>WHERE owner_id = $1 FOR UPDATE"]
        
        subgraph DoubleEntryLedger["financial_ledger_entries (Immutable Table)"]
            Leg1["DEBIT Leg (Account 1010: Cash Clearing)<br>+35,000,000 Santims"]
            Leg2["CREDIT Leg (Account 2010: Owner Payable)<br>-31,500,000 Santims"]
            Leg3["CREDIT Leg (Account 3010: Platform Commission)<br>-3,500,000 Santims"]
        end
        
        BalanceCheck{"Zero-Balance Trigger<br>SUM(amount_minor) == 0?"}
        DebitWallet["UPDATE owner_wallets<br>balance_minor = balance_minor - 3,500,000<br>version = version + 1"]
        UpdateBooking["UPDATE bookings<br>status = 'CONFIRMED'<br>version = version + 1"]
    end

    Confirm --> LockWallet
    LockWallet --> Leg1 & Leg2 & Leg3
    Leg1 & Leg2 & Leg3 --> BalanceCheck
    BalanceCheck -- True (Zero Net Sum) --> DebitWallet
    DebitWallet --> UpdateBooking
    UpdateBooking --> Commit([COMMIT TRANSACTION])
    BalanceCheck -- False --> Rollback([ABORT TRANSACTION & RAISE EXCEPTION])
```

---

## 11. Operator Dispatch, 30-Second GPS Telematics & Geofence Breach Engine

Operating under requirements `GPS-01` to `GPS-04`, the mobile client streams high-precision GPS telematics via WebSocket every 30 seconds when in `ACTIVE` or `IN_TRANSIT` states.

```mermaid
sequenceDiagram
    autonumber
    actor Operator as Field Machine Operator
    participant Mobile as Operator Mobile App
    participant WS as WebSocket Telematics Ingress
    participant VK as Valkey 9.1.3 Cache
    participant Geo as PostGIS Geofence Evaluator
    participant DB as PostgreSQL 18.6 (fleet_location_histories)
    participant Owner as Fleet Owner App

    Operator->>Mobile: Launches App on Assigned Active Booking
    Mobile->>Mobile: Prompts Explicit Location Tracking Consent (GPS-04)
    Operator->>Mobile: Taps "Agree & Enable Tracking"
    Mobile->>WS: Connects WSS /ws/v1/bookings/{id}/tracking (Bearer JWT)
    WS->>DB: Verify Active Consent Record in operator_consent_logs

    loop Every 30 Seconds During ACTIVE State
        Mobile->>Mobile: Query GPS Sensor (High Accuracy, distanceFilter=25m)
        Mobile->>WS: Send JSON Packet { lat: 9.032, lng: 38.742, speedKnots: 1.2, heading: 180 }
        WS->>VK: SETEX telem:booking:{id} 60s { coords, speed }
        WS->>Geo: ST_DWithin(current_point, job_site_center, approved_radius_meters)
        
        alt Outside Approved Geofence (> 1,000m)
            Geo-->>WS: GEOFENCE BREACH DETECTED
            WS->>DB: INSERT INTO geofence_breach_events (booking_id, coords, timestamp)
            WS->>Owner: High-Priority FCM Push 'ALERT: Machine moved outside approved site boundary!'
        else Inside Geofence
            Geo-->>WS: Geofence Normal (OK)
        end
        
        WS->>DB: INSERT INTO fleet_location_histories (fleet_id, location_point, speed, recorded_at)
        WS->>Owner: WebSocket Broadcast Location Marker Update
    end

    Note over Mobile,DB: Rental Ends -> Status Transitions to COMPLETED
    WS->>WS: Immediate Connection Termination (GPS-03-BE)
    Mobile->>Mobile: Background Tracking Service Automatically Shuts Down
    DB->>DB: Scheduled Cron: Prune GPS Coordinate Trails Older than 30 Days (ASVS L3)
```

---

## 12. Server-Driven Dynamic Admin UI Engine (`/api/v1/admin/ui-config`)

The Next.js 16 back-office console renders menus, table columns, search filters, and RBAC-gated action buttons dynamically from `/api/v1/admin/ui-config`.

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Operational Staff / Super Admin
    participant Browser as Next.js 16 Web Client
    participant AuthContext as React Server Component
    participant UIEngine as Admin UI Config Interactor
    participant DB as PostgreSQL 18.6 (admin_table_configs)

    Admin->>Browser: Opens Admin Console URL (/admin/fleets)
    Browser->>AuthContext: Validates ES256 Session Cookie
    AuthContext->>UIEngine: GET /api/v1/admin/ui-config (Role: "finance_admin")
    
    UIEngine->>DB: SELECT * FROM admin_sidebar_items WHERE required_permission IN (user_permissions)
    UIEngine->>DB: SELECT * FROM admin_table_configs WHERE page_route = '/admin/fleets'
    UIEngine->>DB: SELECT * FROM admin_actions WHERE page_route = '/admin/fleets' AND required_permission IN (user_permissions)
    
    UIEngine-->>Browser: UIConfigPayload {<br>  sidebar: [...],<br>  columns: [{ key: "registration_number", label: "Plate", sortable: true }, ...],<br>  filters: [{ key: "category", type: "dropdown" }, ...],<br>  actions: [{ label: "Approve Listing", actionId: "ACT_FLEET_APPROVE", variant: "primary" }]<br>}

    Browser->>Browser: DynamicLayout Component Mounts
    Browser->>Browser: Sidebar Renders Permitted Links Only
    Browser->>Browser: TanStack Table Dynamically Compiles Grid from Column Metadata
    Browser->>Browser: Action Buttons Render Filtered by User Permissions
    Admin->>Browser: Taps Permitted Action Button "Approve Listing"
    Browser->>UIEngine: Executes Action API Endpoint (/api/v1/admin/fleets/{id}/approve)
```

---

## 13. Dynamic Multilingual Typography & OTA Font Pipeline (`ADR-013`)

To render Ethiopic (`am`, `om`, `ti`) and Arabic scripts cleanly without inflating the base application binary, non-Latin typography is delivered over-the-air (OTA) with SHA-256 integrity validation.

```mermaid
flowchart TD
    AppLaunch[Mobile App / Web Bootstrap] --> ReadPref[Read Preferred Language: preferred_lang]
    ReadPref --> LatinCheck{Is Language Latin-Script?<br>e.g. English}
    
    LatinCheck -- Yes --> LoadLocalLatin[Load Pre-bundled Asset Font: Inter / Roboto]
    LatinCheck -- No --> CheckStorage{Font Binary Cached in<br>Local App Sandbox Directory?}

    CheckStorage -- Exists --> HashVerify{Compute SHA-256 Checksum<br>Matches Baked Config Hash?}
    HashVerify -- Valid --> LoadMemory[Add to Flutter FontLoader via Memory]
    HashVerify -- Corrupt / Mismatch --> EvictCorrupt[Delete Corrupted Blob from Storage]
    
    CheckStorage -- Missing --> FetchOTA[Fetch Font Binary from S3 CDN Endpoint]
    EvictCorrupt --> FetchOTA
    
    FetchOTA --> DownloadBytes[Download Font Subset: NotoSansEthiopic.ttf]
    DownloadBytes --> ValidateDownloaded{Verify Download SHA-256 Hash}
    
    ValidateDownloaded -- Failure --> AlertFallback[Log Error & Fallback to System Standard Font]
    ValidateDownloaded -- Success --> PersistSandbox[Save Binary to Secure App Directory]
    
    PersistSandbox --> LoadMemory
    LoadMemory --> ApplyTheme[Register in FontLoader & Trigger Global Text Theme Rebuild]
    ApplyTheme --> CleanRender[Render Ge'ez / Arabic Glyphs without Clipped Ascenders]
```

---

## 14. Offline-First Resilience, Isar Local Store & Background Sync Engine

Field inspectors in remote locations record condition reports offline. Data is cached in local Isar NoSQL storage and synchronized idempotently upon connection recovery.

```mermaid
flowchart TD
    subgraph OfflineCapture["Offline Field Capture Phase (No Connectivity)"]
        FieldOperator[Field Operator on Remote Site] --> SubmitInsp[Completes Handover Inspection Form]
        SubmitInsp --> NetCheck{Network Available?}
        NetCheck -- Offline --> LocalStore[Write Inspection Draft to Isar DB]
        LocalStore --> GenerateDraftId[Assign UUIDv4 localDraftId]
        GenerateDraftId --> StatusQueued[Set syncStatus = SyncStatus.queued]
        StatusQueued --> CacheMedia[Persist Full-Res Photos in App File Cache]
    end

    subgraph SyncLoop["Connectivity Restoration & Background Sync Daemon"]
        NetListener[connectivity_plus Listener] --> ReconnectEvent[Network Interface Restored: 3G / LTE]
        ReconnectEvent --> SyncEngine[Trigger SyncEngine.processQueue]
        
        SyncEngine --> ReadQueue[Read All Inspections WHERE syncStatus == queued]
        ReadQueue --> CompressIsolate["Background Worker Isolate (compute)<br>Resize Photos to WebP 1600px, Quality 82%"]
        
        CompressIsolate --> StreamMedia["POST /media/upload<br>(Stream Bytes through ClamAV Daemon)"]
        StreamMedia --> S3Urls[Receive S3 Object URLs]
        
        S3Urls --> MutatingRPC["POST /api/v1/inspections/sync<br>(Header: Idempotency-Key == localDraftId)"]
        
        MutatingRPC --> BackendResponse{Backend Response Status}
        BackendResponse -- 200 / 201 Success --> MarkSynced[Update Isar: syncStatus = SyncStatus.synced]
        BackendResponse -- 409 Conflict --> MarkConflict[Update Isar: syncStatus = SyncStatus.failed]
        
        MarkSynced --> PurgeLocal[Purge Temporary File Cache]
        MarkConflict --> PromptOperator[Surface Manual Concurrency Resolution Prompt]
    end
```

---

# Part II: Industrial Ergonomics & UI/UX Design System

## 15. Executive Summary & Product Persona Context

Dink Machinery bridges heavy civil engineering contractors, industrial equipment owners, and field operators across East Africa. The user experience must accommodate high-value financial commitments while functioning reliably on dusty, sunlit construction corridors with intermittent cellular connectivity.

The platform balances four distinct personas:
1. **The Renter (Civil Contractor / Project Manager):** Needs fast spatial discovery, reliable machine dimensions, verified owner credibility, and transparent, binding rental agreements.
2. **The Fleet Owner (Asset Investor / Equipment Yard Manager):** Manages multi-million Birr machinery, demands payment verification, tracks operator assignments, and requires visibility into wallet commission balances.
3. **The Operator (Field Technician / Heavy Equipment Driver):** Works on active job sites, interacts with physical hour meters and fuel gauges, wears protective gloves, and uses entry-level or mid-tier smartphones.
4. **The Platform Administrator (Compliance, Finance & Dispute Officer):** Operates on desktop browsers, requires high-density server-driven data tables, reviews photographic inspection diffs, and executes legally binding dispute rulings.

---

## 16. Industrial Ergonomics & Environmental Constraints

### 16.1 Direct Sunlight & High-Glare Visibility
Construction corridors in East Africa frequently exceed 80,000 to 100,000 lux ambient sunlight. Pastel palettes become illegible in direct glare.

* **Outdoor High-Contrast Mode:** Automatically triggered when device lux sensors detect ambient light $> 25,000\text{ lux}$, or manually enabled via a toggle in the top app bar.
* **Contrast Ratios:** Standard UI targets WCAG 2.2 AA ($4.5:1$ for body text, $3:1$ for large text). Outdoor High-Contrast Mode enforces WCAG 2.2 AAA ($\ge 7.0:1$ across all interactive surfaces and $\ge 10:1$ for engine telemetry).
* **High-Contrast Badges:** Informational states use solid-fill pills with distinct geometric icons (e.g., solid triangle for warnings, filled octagonal stop sign for credit lockouts).

### 16.2 Gloved Hands, Vibration & Tactile Feedback
Operators face heavy mechanical vibrations and wear protective leather or polyurethane gloves.

* **Touch Boundaries:** All primary touch controls meet or exceed $48 \times 48\text{ dp}$ (target $56 \times 56\text{ dp}$ for camera triggers and inspection checks).
* **Touch Separation:** Minimum physical separation of $8\text{ dp}$ between adjacent targets.
* **Tactile Haptic Confirmation:** Mutating operations trigger a haptic pulse (`HapticFeedback.heavyImpact` on iOS and Android).
* **Audio Cues:** Inspection steps emit an optional high-frequency chime ($85\text{ dB}$, $1200\text{ Hz}$) confirming photo uploads and valid meter inputs.

### 16.3 One-Handed Field Operation & Thumb Zones
Operators frequently operate their devices with one hand while climbing access ladders or inspecting machine linkages.

```
+-----------------------------------+
|      Hard-to-Reach Zone           |  <- Read-only metadata, Status pills
|   (Status badges, Title, Back)    |
+-----------------------------------+
|                                   |
|      Natural Reach Zone           |  <- Content cards, Inspection photo grids,
|     (Form fields, Sliders)        |     Radio selectors
|                                   |
+-----------------------------------+
|      Optimal Thumb Action Zone    |  <- Primary CTAs: "Confirm Handover",
|   (Sticky Bottom Bar, Shutter)    |     "Take Photo", "Upload Proof"
+-----------------------------------+
```

---

## 17. Design System Foundations & Tokens

### 17.1 Color System (Standard & Outdoor High-Contrast Modes)

The palette balances heavy industrial heritage (safety amber, komatsu blue) with strict functional hierarchy.

#### Semantic Color Tokens

| Token Name | Standard Light Mode | Standard Dark Mode | Outdoor High-Contrast Mode | Functional Role |
|---|:---:|:---:|:---:|---|
| `--color-brand-primary` | `#F59E0B` (Amber 500) | `#FBBF24` (Amber 400) | `#FFE600` (Safety Yellow) | Primary brand accents, key milestones |
| `--color-brand-secondary` | `#0F172A` (Slate 900) | `#F8FAFC` (Slate 50) | `#000000` (Pure Black) | Core surface contrast, primary text |
| `--color-surface-bg` | `#F8FAFC` (Slate 50) | `#0B0F17` (Dark Pitch) | `#FFFFFF` (Pure White) | Application canvas background |
| `--color-surface-card` | `#FFFFFF` (Pure White) | `#1E293B` (Slate 800) | `#FFFFFF` (Bordered 2px Black) | Elevated content containers |
| `--color-action-cta` | `#D97706` (Amber 600) | `#F59E0B` (Amber 500) | `#000000` (White on Black) | Primary action buttons |
| `--color-status-success` | `#10B981` (Emerald 500) | `#34D399` (Emerald 400) | `#056B38` (Deep Forest) | Active rental, confirmed payment |
| `--color-status-warning` | `#F59E0B` (Amber 500) | `#FBBF24` (Amber 400) | `#8F5B00` (Dark Ochre) | Hold pending, SLA warning |
| `--color-status-danger` | `#EF4444` (Red 500) | `#F87171` (Red 400) | `#B91C1C` (Deep Crimson) | Disputed state, cancellation penalty |
| `--color-telematics-meter`| `#2563EB` (Blue 600) | `#60A5FA` (Blue 400) | `#002D9C` (Navy High-Viz) | Engine hours, fuel gauges, GPS markers |

### 17.2 Typography Engine & Multi-Script Baseline Alignment (`ADR-013`)

The application supports simultaneous Latin and Ge'ez (Ethiopic) scripts. Ge'ez characters include tall ascenders and deep descenders; inadequate line-height settings cause glyph clipping.

```
Ge'ez Ascender Line    -----------------------------  (ማ, ፕ, ፭)
Latin Cap Height       =============================  (H, M, K)
Baseline               _____________________________
Ge'ez Descender Line   -----------------------------  (ሟ, ጴ, ጇ)
```

#### Multi-Script Typography Scale

| Style Token | Font Family | Weight | Size | Line Height | Tracking | Script Constraints |
|---|---|---|---|---|---|---|
| `display-large` | Roboto / Noto Sans Ethiopic | Bold (700) | 32sp | 44sp | -0.5px | Min line-height multiplier $1.38\times$ for Ge'ez |
| `headline-medium` | Roboto / Noto Sans Ethiopic | SemiBold (600) | 24sp | 34sp | 0px | Prevents upper diacritic clipping |
| `title-medium` | Roboto / Noto Sans Ethiopic | Medium (500) | 18sp | 26sp | +0.15px | Field labels, dialog titles |
| `body-large` | Roboto / Noto Sans Ethiopic | Regular (400) | 16sp | 24sp | +0.5px | Primary body, spec descriptions |
| `body-medium` | Roboto / Noto Sans Ethiopic | Regular (400) | 14sp | 22sp | +0.25px | Secondary metadata, review comments |
| `telematics-mono` | JetBrains Mono / Space Mono | Bold (700) | 20sp | 28sp | +1.0px | Engine hour counters, GPS coordinates |

* **OTA Font Pipeline (`ADR-013`):** Latin characters use pre-bundled Inter/Roboto. Ethiopic and Arabic glyphs are delivered over-the-air from SeaweedFS S3 upon locale selection, verified by SHA-256 integrity checks, and loaded dynamically into memory via Flutter's `FontLoader`.

### 17.3 8dp Spatial Grid, Layout Boundaries & Elevation
Layouts strictly adhere to an 8dp spatial grid system ($8, 16, 24, 32, 48, 64\text{ dp}$). Half-grid increments ($4\text{ dp}$) are used exclusively for tight badge padding and icon-to-label spacing.

* **Card Corner Radii:** Uniform $12\text{ dp}$ for interactive cards; $8\text{ dp}$ for internal input containers; $24\text{ dp}$ for pill badges and button caps.
* **Elevation & Layering:**
  - `Level 0` ($0\text{ dp}$): Background canvas (`#F8FAFC`).
  - `Level 1` ($2\text{ dp}$ shadow / $1\text{ px}$ border): Listing cards, spec accordions.
  - `Level 2` ($4\text{ dp}$ shadow): Floating filter chips, scrolling app bars.
  - `Level 3` ($8\text{ dp}$ shadow): Sticky bottom action sheets, date-picker drawers.
  - `Level 4` ($16\text{ dp}$ shadow): Blocking modals, handover signature canvas overlay.

### 17.4 Component Touch Targets ($48 \times 48\text{ dp}$)

Every interactive control enforces a minimum touch target envelope:

```
+---------------------------------------+
|  Invisible Touch Target: 48 x 48 dp   |
|      +-------------------------+      |
|      | Visible Icon: 24 x 24 dp|      |
|      +-------------------------+      |
+---------------------------------------+
```

---

## 18. Cross-Platform Persona Journeys & Wireflows

### 18.1 Equipment Renter Experience (Discovery to Handover)

```mermaid
journey
    title Equipment Renter Discovery & Booking Journey
    section Spatial Discovery
      Launch App & Check Location: 5: Renter
      Apply Dynamic Spec Filters (Tonnage, Boom): 4: Renter
      Inspect PostGIS Distance & Base Rates: 5: Renter
    section Commitment & Hold
      Select Date Interval (Min Rental Duration Check): 4: Renter
      Review Cost Card (Subtotal, 15% VAT, Informational Deposit): 5: Renter
      Request 60-Min Reservation Hold: 5: Renter
    section Bank Settlement
      View Owner Bank Coordinates & Copy Account: 4: Renter
      Transfer via Mobile Banking (CBE / Telebirr): 3: Renter
      Upload Cropped Payment Slip Image: 4: Renter
    section Handover
      Receive Machine at Construction Site: 5: Renter
      Inspect Baseline Hour Meter & 4-Angle Photos: 5: Renter
      Counter-Sign Handover on Glass Canvas: 5: Renter
```

#### Key Renter Screens
1. **Catalog Home & Geospatial Facets:** Search header with PostGIS radius slider ($5\text{–}200\text{ km}$), category icons, and keyset-paginated cards showing price/day, distance in kilometers, operator badge, and verified owner status.
2. **Dynamic Fleet Detail View:** Swipeable 5-to-20 photo carousel, interactive specs accordion grouped by system (Engine, Hydraulics, Dimensions), public insurance badge, and sticky bottom booking bar.
3. **Pre-Booking Summary Modal (`BK-06B-MOB`):** Non-editable confirmation sheet itemizing rental subtotal, 15% VAT, total due, and separate informational deposit notice.
4. **Active 60-Minute Payment Hold Screen (`BK-06-MOB`):** Animated circular countdown timer, copyable owner bank coordinates, transfer instructions, and direct image upload button.

### 18.2 Fleet Owner Experience (Fleet Registry, Availability & Wallet)

```mermaid
journey
    title Fleet Owner Operational Journey
    section Fleet Onboarding
      Select Subcategory (e.g. Midi Excavator): 5: Owner
      Fill Dynamic EAV Specs (Engine Power, Weight): 4: Owner
      Capture 5 Mandatory Angles + Meter Photo: 4: Owner
      Set Base Rates & Cancellation Penalty Rules: 4: Owner
    section Operational Control
      Block Maintenance Dates on Calendar: 5: Owner
      Assign Certified Team Operator: 5: Owner
      Review Payment Proof & Confirm Reservation: 5: Owner
    section Financial Oversight
      Inspect Pre-Funded OwnerWallet Balance: 4: Owner
      Upload Wallet Top-Up Receipt: 4: Owner
      Monitor Automated Commission Debits: 5: Owner
```

#### Key Owner Screens
1. **Fleet Listing Wizard:** Cascading category pickers triggering dynamic accordion fields, camera upload slots for required angles, and pricing inputs with minor-unit currency formatting.
2. **Availability Calendar Board (`FL-35-MOB`):** Interactive monthly calendar where owners drag and tap date intervals to toggle active availability, pending reservations, or manual maintenance blocks.
3. **Payment Evidence Review Drawer (`BK-09-MOB`):** Zoomable viewer for bank slips with duplicate detection flags (`BK-07B`), transaction metadata, and "Confirm Payment" / "Decline" action buttons.
4. **OwnerWallet Dashboard (`CW-02-MOB`):** Balance display, credit limit bar, top-up entry point, and infinite-scroll ledger table showing commission debits linked to booking IDs.

### 18.3 Field Operator Experience (Inspection Stepper & GPS Telematics)

```mermaid
journey
    title Field Operator Machine Inspection Flow
    section Job Check-In
      Receive Job Assignment & Location Alert: 5: Operator
      Review Job Site Coordinates & Accept: 5: Operator
      Launch Pre-Rental Inspection Stepper: 5: Operator
    section Machine Inspection
      Capture Photo of Physical Engine Hour Meter: 4: Operator
      Input Numeric Meter (Auto-compared to photo): 4: Operator
      Adjust Fuel Level Gauge Slider (0-100%): 5: Operator
      Capture 4 Mandatory Angle Photos: 4: Operator
    section Attestation
      Sign Handover Attestation on Touch Pad: 5: Operator
      Await Counter-Signature from Renter: 4: Operator
      Toggle 30-Second GPS Telematics Stream: 5: Operator
```

#### Key Operator Screens
1. **Operator Active Task Dashboard:** High-contrast card showing assigned machine make/model, job site pin, delivery deadline, and large "Start Field Inspection" CTA.
2. **Camera-Enforced Meter Capture (`AR-01-MOB`):** Custom camera overlay featuring a rectangular reticle for the physical hour meter gauge, with numerical confirmation entry immediately below.
3. **Digital Fuel & Condition Check:** Chunky horizontal slider for fuel percentage ($0\text{–}100\%$ in $5\%$ increments) and toggleable chips for equipment defects.
4. **Glass Signature Pad (`AR-08-MOB`):** Full-screen orientation-locked canvas capturing finger signatures, paired with recipient name and role fields.

### 18.4 Operations & Compliance Admin Console (Next.js 16 Web)

* **Server-Driven Dynamic Data Grids (`AD-06-WEB`):** The console does not hardcode table columns. Grids parse `/api/v1/admin/ui-config`, rendering dynamic columns, sorting rules, and permission-gated action buttons.
* **Split-Screen Dispute Lightbox (`AD-04-WEB`):** Renders pre-rental check-in photos directly beside post-rental check-out photos with synchronized pan-and-zoom inspection.
* **Audit Diff Viewer (`AD-09-WEB`):** Expandable rows displaying visual JSON diffs comparing previous and updated system states.

---

## 19. Core Component Library & Interaction Specifications

### 19.1 Dynamic EAV Form Accordion Engine
Renders dynamic specification fields based on category schema JSON.

```
+-------------------------------------------------------------------+
|  [v] ENGINE & POWERTRAIN SPECIFICATIONS                           |
+-------------------------------------------------------------------+
|  Engine Manufacturer *                                            |
|  [ Isuzu 4HK1X                              ]                     |
|                                                                   |
|  Gross Power (hp) *                                               |
|  [ 162                                      ] hp                  |
|                                                                   |
|  Fuel Type *                                                      |
|  (o) Diesel    ( ) Electric    ( ) Hybrid                         |
+-------------------------------------------------------------------+
|  [>] HYDRAULIC & WORKING RANGES (3 Incomplete)                    |
+-------------------------------------------------------------------+
|  [>] PHYSICAL DIMENSIONS & WEIGHT (Complete)                      |
+-------------------------------------------------------------------+
```

* **Section Completion Pills:** Closed accordions display a badge showing incomplete required fields (e.g., `"3 Required Missing"` in amber).
* **Auto-Scroll to Error:** On submit, if validation fails, the accordion collapses unaffected sections, opens the first invalid section, and scrolls smoothly to the invalid input.

### 19.2 Camera-Enforced Telematics Capture Overlay
Restricts photo acquisition to the active device camera (banning gallery uploads for field inspections).

```
+-------------------------------------------------------------------+
|  <- Back                   Pre-Rental Check-In             [ ? ]  |
+-------------------------------------------------------------------+
|                                                                   |
|         +-----------------------------------------------+         |
|         |  ALIGN ENGINE HOUR METER GAUGE INSIDE RETICLE |         |
|         |                                               |         |
|         |            [   1 4 5 0 . 5 0   ]              |         |
|         |                                               |         |
|         +-----------------------------------------------+         |
|                                                                   |
|  Enter Observed Gauge Reading *                                   |
|  [ 1450.50                                  ] Hours               |
|                                                                   |
|  Fuel Level: 85%                                                  |
|  [===|=========================================o=======] 85%       |
|                                                                   |
|                         [ (O) SHUTTER ]                           |
+-------------------------------------------------------------------+
```

* **Reticle Guide:** Semi-transparent viewfinder ($4:3$ aspect ratio) indicating gauge placement.
* **Client-Side Image Optimization:** Captured raw images resize on a background Dart isolate to a maximum dimension of $1600\text{ px}$ at $82\%$ WebP quality, reducing uploads from $\approx 12\text{MB}$ to $< 800\text{KB}$ before virus scanning.

### 19.3 Countdown Timer & Bank Transfer Card
Maintains urgency and clarity during the 60-minute payment hold window.

```
+-------------------------------------------------------------------+
|  RESERVATION HOLD EXPIRES IN:                                      |
|                            [ 48:22 ]                              |
|  Complete bank transfer and upload receipt before timer expires.  |
+-------------------------------------------------------------------+
|  OWNER DIRECT BANK COORDINATES                                    |
|                                                                   |
|  Bank: Commercial Bank of Ethiopia (CBE)                          |
|  Account Name: Mesfin Industrial Equipment PLC                    |
|  Account Number: 1000123456789                      [ COPY NO ]   |
|                                                                   |
|  Required Prepayment: 35,000.00 ETB                               |
+-------------------------------------------------------------------+
|  [ UPLOAD TRANSFER RECEIPT SCREENSHOT ]                           |
+-------------------------------------------------------------------+
```

* **Color Interpolation:** Timer circle stays Emerald from 60 to 20 minutes, transitions to Amber between 20 and 5 minutes, and pulses Crimson under 5 minutes.
* **Clipboard Interaction:** Tapping `[ COPY NO ]` updates the system clipboard, triggers haptic feedback, and changes label to `"COPIED!"` for 2.5 seconds.

### 19.4 Interactive Handover Checklist & Glass Signature Canvas
Ensures binding dual-party sign-off at equipment dispatch and return.

```
+-------------------------------------------------------------------+
|  EQUIPMENT CONDITION ATTESTATION                                   |
|                                                                   |
|  [x] Engine oil & hydraulic fluid verified at operational levels   |
|  [x] All safety decals & ROPS/FOPS structural points intact       |
|  [x] Demountable bucket attachment securely pinned & locked       |
|                                                                   |
|  Receiving Party Signature (Sign inside box)                      |
|  +-------------------------------------------------------------+  |
|  |                                                             |  |
|  |                /s/ Dawit Wolde                              |  |
|  |                                                             |  |
|  +-------------------------------------------------------------+  |
|  Signee Name: Dawit Wolde             Signee Role: Site PM        |
|                                                                   |
|  [ CONFIRM & ACTIVATE LEASE ]                                     |
+-------------------------------------------------------------------+
```

* **Minimum Stroke Requirement:** The signature canvas requires a minimum of 20 vector points before enabling submission, preventing accidental taps from submitting empty signatures.
* **Fallback Authentication:** If the physical touch glass is damaged or wet, the user can switch to `"Verify via Phone OTP"`.

### 19.5 Split-Screen Dispute Arbitration Lightbox
Administrative view comparing check-in vs check-out machine condition.

```
+-------------------------------------------------------------------+
|  DISPUTE #DSP-9042: CRAWLER EXCAVATOR (BOOM CYLINDER CRACK)       |
+---------------------------------+---------------------------------+
|  CHECK-IN (2026-10-12 09:14 UTC)|  CHECK-OUT (2026-10-18 16:30 UTC|
|                                 |                                 |
|       [ Image Viewport A ]      |       [ Image Viewport B ]      |
|    (Hydraulic Cylinder Clean)   |    (Hydraulic Cylinder Crack)   |
|                                 |                                 |
|  Recorded Meter: 1,450.50 Hrs   |  Recorded Meter: 1,498.50 Hrs   |
+---------------------------------+---------------------------------+
|  Synchronized Zoom: [ 100% | 200% | 400% ]     [ Synchronize Pan ]|
|                                                                   |
|  Claimed Damage: 45,000.00 ETB       Suggested Split: 50% / 50%   |
|  [ RULE FOR RENTER ]    [ RULE PARTIAL SPLIT ]    [ RULE FOR OWNER]|
+-------------------------------------------------------------------+
```

---

## 20. Motion, Micro-Interactions & Audio-Haptic Feedback

1. **Page Transitions:** Mobile views use horizontal slide transitions (`CupertinoPageRoute` on iOS, shared axis material motion on Android) with a duration of $250\text{ ms}$.
2. **Pull-to-Refresh:** Custom pull-to-refresh animation showing a rotating excavator gear; triggers light haptic pulse upon crossing release threshold.
3. **Card Expansion Animation:** Expanding dynamic spec accordions uses a cubic ease-in-out curve ($300\text{ ms}$) without layout shifts.
4. **Shutter Micro-Interaction:** Tapping the inspection camera shutter triggers a $50\text{ ms}$ screen flash and haptic tap, followed by an immediate progress indicator.
5. **Reduced Motion Mode:** Respects system accessibility toggles (`prefers-reduced-motion`). When enabled, all non-essential spatial transitions switch to instant opacity cross-fades ($100\text{ ms}$).

---

## 21. Offline Cues, Empty States & RFC 9457 Error Presentation

### 21.1 Persistent Offline Status Banner
When network connectivity is lost, Flutter's `connectivity_plus` listener mounts a sticky, high-contrast banner below the app bar.

```
+-------------------------------------------------------------------+
|  [!] OFFLINE MODE: Showing cached fleet data (Updated 2h ago).    |
|      Field inspections will automatically sync once reconnected.  |
+-------------------------------------------------------------------+
```

* **Color:** High-visibility amber background (`#D97706`) with pure white text and warning triangle icon.
* **Behavior:** Persistent across all mobile views; does not dismiss until reliable socket connectivity is restored.

### 21.2 RFC 9457 Problem Details Presentation
Errors from the backend map into user-friendly localized cards while preserving tracking IDs for support assistance.

```
+-------------------------------------------------------------------+
|  (X) UNABLE TO COMPLETE RESERVATION                               |
|                                                                   |
|  This equipment was just reserved by another contractor for the   |
|  selected dates.                                                  |
|                                                                   |
|  Reason: Conflicting Reservation Active                           |
|  Reference Code: FLEET_ALREADY_RESERVED                           |
|  Support Tracking ID: req-8842-1982-9912                          |
|                                                                   |
|  [ CHOOSE ALTERNATIVE DATES ]         [ VIEW SIMILAR EXCAVATORS ] |
+-------------------------------------------------------------------+
```

* **Inline Form Errors:** Validation errors targeting specific fields render red text (`#EF4444`) with an error outline around the input container.
* **Copyable Support ID:** Tapping the Support Tracking ID copies the UUID to clipboard, facilitating rapid support ticket resolution.

---

*This document defines the canonical unified architecture, sequence choreography, and UI/UX design system specification for Dink Machinery Version 7.0.*