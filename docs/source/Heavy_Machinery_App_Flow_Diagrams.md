# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## System Application Flow Diagrams & Architecture Choreography (`docs/source/Heavy_Machinery_App_Flow_Diagrams.md`)

**Document Identifier:** `FLOW-DINK-CORE-V7.0`  
**Classification:** Canonical Architecture & Development Reference (Single Source of Truth)  
**Governing Documents:** `SRS-DINK-CORE-V7.0` (Parts 1–5), `contracts/openapi.yaml`, `docs/components/mobile/DESIGN.md`, `docs/components/mobile/AGENTS.md`  
**Target Runtimes:** Go 1.27.3 Modular Monolith, Flutter 3.44 Cross-Platform Mobile Client, Next.js 16 Admin Web Console, PostgreSQL 18.4 + PostGIS 3.4, Valkey 8.0-alpine, SeaweedFS S3  
**Target Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), W3C WCAG 2.1/2.2 AA, ISO/IEC/IEEE 29148, ISO 4217, ISO 3166-1 alpha-2  
**Status:** Approved Production Architecture

---

## Table of Contents

1. [High-Level System Context & Multi-Persona Ingress](#1-high-level-system-context--multi-persona-ingress)
2. [Identity, Registration & Pluggable Verification Flows](#2-identity-registration--pluggable-verification-flows)
   - 2.1 [Passwordless Phone OTP Authentication & Refresh Rotation](#21-passwordless-phone-otp-authentication--refresh-rotation)
   - 2.2 [Pluggable Identity Verification (Fayda / National ID Adapter)](#22-pluggable-identity-verification-fayda--national-id-adapter)
   - 2.3 [Multi-Factor Authentication (Passkey WebAuthn & TOTP)](#23-multi-factor-authentication-passkey-webauthn--totp)
3. [Fleet Registration, Dynamic EAV Taxonomy & Availability Scheduling](#3-fleet-registration-dynamic-eav-taxonomy--availability-scheduling)
4. [Geospatial Radius Discovery, Corridor Search & Filtering (PostGIS)](#4-geospatial-radius-discovery-corridor-search--filtering-postgis)
5. [End-to-End Booking Reservation, 60-Minute Valkey Hold & Dual Settlement](#5-end-to-end-booking-reservation-60-minute-valkey-hold--dual-settlement)
   - 5.1 [Manual Direct Transfer Flow (`MANUAL_TRANSFER` / Direct Transfer Model)](#51-manual-direct-transfer-flow-manual_transfer--direct-transfer-model)
   - 5.2 [Online Payment Gateway Flow (`GATEWAY` / Pluggable Adapter)](#52-online-payment-gateway-flow-gateway--pluggable-adapter)
6. [The Canonical 20-State Booking Machine Choreography](#6-the-canonical-20-state-booking-machine-choreography)
7. [Handover Check-In & Check-Out Inspection Workflow (Monotonic Meters & ClamAV)](#7-handover-check-in--check-out-inspection-workflow-monotonic-meters--clamav)
8. [Damage Dispute Arbitration, SLA Escalation & Escrow Settlement](#8-damage-dispute-arbitration-sla-escalation--escrow-settlement)
9. [Logistics Transport Sub-Booking & Bidirectional Cascade Cancellation](#9-logistics-transport-sub-booking--bidirectional-cascade-cancellation)
10. [Double-Entry Financial Ledger & Wallet-Debit Commission Accounting](#10-double-entry-financial-ledger--wallet-debit-commission-accounting)
11. [Operator Dispatch, 30-Second GPS Telematics & Geofence Breach Engine](#11-operator-dispatch-30-second-gps-telematics--geofence-breach-engine)
12. [Server-Driven Dynamic Admin UI Engine (`/api/v1/admin/ui-config`)](#12-server-driven-dynamic-admin-ui-engine-apiv1adminui-config)
13. [Dynamic Multilingual Typography & OTA Font Pipeline (`ADR-013`)](#13-dynamic-multilingual-typography--ota-font-pipeline-adr-013)
14. [Offline-First Resilience, Isar Local Store & Background Sync Engine](#14-offline-first-resilience-isar-local-store--background-sync-engine)

---

## 1. High-Level System Context & Multi-Persona Ingress

The Dink Machinery ecosystem orchestrates interactions across three primary edge personas (Renters, Fleet Owners, Field Operators) on mobile, plus back-office administrative operations on the web console. All traffic terminates at the Traefik 3.x edge ingress, routing to a single compiled Go 1.27.3 modular monolith backend.

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
        GeoMod["Spatial Search Engine<br>(PostGIS 3.4 Bounded Queries)"]
        BookingMod["Booking State Engine<br>(20-State Machine, OCC Versioning)"]
        PaymentMod["Escrow Ledger & Wallet Engine<br>(Double-Entry, BK-17-BE)"]
        InspectMod["Inspection & Media Ingress<br>(ClamAV Stream Scan)"]
        AdminUIMod["Server-Driven UI Engine<br>(/api/v1/admin/ui-config)"]
    end

    subgraph PersistenceInfra["Storage & State Infrastructure"]
        Postgres[("PostgreSQL 18.4 + PostGIS 3.4<br>(Spatial Points, JSONB Specs, pgxpool)")]
        Valkey[("Valkey 8.0 Cache<br>(60-Min Mutex Holds, Rate-Limit Counters)")]
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

Authentication adheres to OWASP ASVS 4.0.3 Level 3 with single-use rotating refresh tokens stored in Valkey and cryptographic ES256 session JWTs.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Flutter Client (Renter/Owner/Operator)
    participant GW as Traefik API Gateway
    participant Auth as Go Auth Module (net/http)
    participant VK as Valkey 8.0 Cache
    participant DB as PostgreSQL 18.4 (users)
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

Verification uses a provider-agnostic framework (`REG-IDP-01-BE`). When active, missing required identity proofs block mutating transactional APIs (`AUTH-05-BE`).

```mermaid
sequenceDiagram
    autonumber
    actor User as Authenticated User
    participant App as Mobile Client
    participant API as Auth Identity Interactor
    participant Policy as Policy Engine (REG-IDP-02-BE)
    participant Adapter as IdentityVerificationProvider (Fayda)
    participant Registry as National Fayda API
    participant DB as PostgreSQL 18.4 (identity_verifications)

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

Fleet asset registration dynamically validates technical specifications using dynamic JSON Schema attributes stored in `spec_field_definitions` based on the chosen category.

```mermaid
sequenceDiagram
    autonumber
    actor Owner as Fleet Owner
    participant App as Mobile Client
    participant CatAPI as Category & Spec Module
    participant MediaAPI as Media & Ingress Controller
    participant ClamAV as ClamAV Antivirus Daemon
    participant S3 as SeaweedFS Object Storage
    participant DB as PostgreSQL 18.4 (fleets)
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

Dink Machinery executes high-throughput, low-latency spheroidal distance queries over PostGIS GiST indexes, ensuring queries complete within $\le 80\text{ ms}$ at P95 for up to 1,000,000 assets.

```mermaid
flowchart TD
    Client["Renter Mobile Client<br>(GPS Coordinates: Lat, Lng)"] -->|GET /api/v1/geo/fleets/radius-search| Gateway[Traefik 3.x Edge Gateway]
    Gateway --> GeoEngine[Geospatial Query Interactor]

    subgraph QueryParser["Parameter Parsing & Bounds Enforcement"]
        GeoEngine --> BoundingCheck{"Radius Bounds Check<br>(1km <= radiusKm <= 500km)"}
        BoundingCheck -- Exceeded --> Err400[HTTP 400 Bad Request]
        BoundingCheck -- Valid --> CountryContext["Resolve Country Context<br>(country_configs: ET or KE)"]
    end

    subgraph PostGISExecution["PostgreSQL 18.4 + PostGIS 3.4 Engine"]
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
    participant VK as Valkey 8.0 Cache
    participant DB as PostgreSQL 18.4 (pgx.Tx)
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
    participant DB as PostgreSQL 18.4 (pgx.Tx)

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

The rental booking machine enforces 20 deterministic states, strictly governed by optimistic concurrency control (`WHERE id = $1 AND version = $2`). Arbitrary or unmapped transitions abort immediately with RFC 9457 error `BOOKING_INVALID_STATE_TRANSITION`.

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

Physical equipment handovers require a joint digital inspection enforcing monotonic hour meter progression and real-time antivirus scanning.

```mermaid
sequenceDiagram
    autonumber
    actor Renter as Equipment Renter
    actor Operator as Field Operator / Owner
    participant App as Mobile Inspection Stepper
    participant API as Inspection Module
    participant ClamAV as ClamAV Daemon (zINSTREAM)
    participant S3 as SeaweedFS S3
    participant DB as PostgreSQL 18.4 (inspections)

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

Disputed claims (unmatched damage, hour-meter discrepancies) are isolated into a first-class `Dispute` domain with strict 3-day SLA timers and dual administrative communication channels.

```mermaid
sequenceDiagram
    autonumber
    actor Owner as Fleet Owner
    actor Renter as Equipment Renter
    actor Admin as Dispute Arbitration Officer
    participant API as Dispute Module
    participant DB as PostgreSQL 18.4 (pgx.Tx)
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

Logistics transport vehicles (lowboys, flatbeds) are booked as formal sub-bookings linked via `parent_booking_id`. Cancellation cascades preserve operational fairness.

```mermaid
sequenceDiagram
    autonumber
    actor Renter as General Contractor
    actor MachineOwner as Machinery Owner
    actor CarrierOwner as Logistics Carrier Owner
    participant API as Logistics & Booking Engine
    participant DB as PostgreSQL 18.4 (bookings)

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
    participant VK as Valkey 8.0 Cache
    participant Geo as PostGIS Geofence Evaluator
    participant DB as PostgreSQL 18.4 (fleet_location_histories)
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

The Next.js 16 back-office console eliminates hardcoded administrative pages. Menus, data table columns, search filters, and RBAC-gated action buttons are delivered dynamically from `/api/v1/admin/ui-config`.

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Operational Staff / Super Admin
    participant Browser as Next.js 16 Web Client
    participant AuthContext as React Server Component
    participant UIEngine as Admin UI Config Interactor
    participant DB as PostgreSQL 18.4 (admin_table_configs)

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

To guarantee flawless rendering of Ethiopic (`am`, `om`, `ti`) and Arabic scripts without bloating the initial base application binary, non-Latin typography is delivered over-the-air (OTA) with SHA-256 integrity validation.

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

Field inspectors in remote quarries or road construction corridors can record full condition reports offline. Data is isolated in local Isar NoSQL storage and synchronized idempotently upon connection recovery.

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

*This document represents the canonical flow diagram and architectural choreography specification for Dink Machinery Version 7.0. All downstream code generation, sequence verification, and component wiring must strictly mirror the visual and state transitions specified herein.*