# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification (SRS) — Part 1 of 5

**Document Identifier:** `SRS-DINK-CORE-V7.0-PART-1`  
**Version:** 7.0 – v3 National ID (Fayda / Pluggable Identity), MFA & Advanced Ads Edition  
**Classification:** Canonical Architecture & Development Reference (Single Source of Truth)  
**Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Multi-Country Parameterization (`LOC-01-BE`)  
**Target Runtimes:** Go 1.27.3 Modular Monolith, Flutter 3.47.9 Cross-Platform Mobile Client, Next.js 16.x Admin Web Console (latest CVE security patch), PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3-alpine  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), W3C WCAG 2.1/2.2 AA, ISO/IEC/IEEE 29148, ISO 4217, ISO 3166-1 alpha-2  
**Status:** Approved Final Specification

---

## Table of Contents

- 1. [Introduction](#1-introduction)
  - 1.1 [Purpose & Audience](#11-purpose--audience)
  - 1.2 [Scope & Phases](#12-scope--phases)
  - 1.3 [How to Read This Document](#13-how-to-read-this-document)
- 2. [Glossary & Conventions](#2-glossary--conventions)
- 3. [Product Vision & Strategic Context](#3-product-vision--strategic-context)
  - 3.1 [Vision Statement](#31-vision-statement)
  - 3.2 [Core Value](#32-core-value)
  - 3.3 [Ethiopian Market Fit](#33-ethiopian-market-fit)
- 4. [User Roles & Permissions](#4-user-roles--permissions)
  - 4.1 [Role Definitions](#41-role-definitions)
  - 4.2 [Permission Matrix (MVP & Phase 2)](#42-permission-matrix-mvp--phase-2)
- 5. [Functional Requirements – Phase 1 (MVP)](#5-functional-requirements--phase-1-mvp)
  - 5.1 [Registration & Authentication](#51-registration--authentication)
    - 5.1.1 [Renter Registration](#511-renter-registration)
    - 5.1.2 [Owner Registration](#512-owner-registration)
    - 5.1.3 [Operator Registration](#513-operator-registration)
    - 5.1.4 [Authentication & Session Management](#514-authentication--session-management)
    - 5.1.5 [Combined Role Transitions](#515-combined-role-transitions)
  - 5.2 [Fleet Registration – Machinery & Logistics](#52-fleet-registration--machinery--logistics)
    - 5.2.1 [Basic Information](#521-basic-information)
    - 5.2.2 [Dynamic Specifications](#522-dynamic-specifications)
    - 5.2.3 [Pricing, Booking Rules & Payment Details](#523-pricing-booking-rules--payment-details)
    - 5.2.4 [Operator Assignment](#524-operator-assignment)
    - 5.2.5 [Logistics Configuration](#525-logistics-configuration)
    - 5.2.6 [Attachments Bundling](#526-attachments-bundling)
    - 5.2.7 [Availability & Listing Status](#527-availability--listing-status)
  - 5.3 [Attachment Registration & Inventory](#53-attachment-registration--inventory)
  - 5.4 [Fleet Discovery, Search & Comparison](#54-fleet-discovery-search--comparison)
  - 5.5 [Booking & Reservation Flow](#55-booking--reservation-flow)
    - 5.5.1 [Creating a Booking Request](#551-creating-a-booking-request)
    - 5.5.2 [Payment Evidence Upload & Verification](#552-payment-evidence-upload--verification)
    - 5.5.3 [Post-Confirmation & Readiness](#553-post-confirmation--readiness)
    - 5.5.4 [Cancellation Logic](#554-cancellation-logic)
  - 5.6 [Active Rental Management](#56-active-rental-management)
  - 5.7 [Rental Completion, Return & Damage](#57-rental-completion-return--damage)
  - 5.8 [Logistics Transport Integration](#58-logistics-transport-integration)
  - 5.9 [Commission, Revenue & Owner Payouts](#59-commission-revenue--owner-payouts)
  - 5.10 [Notifications & In-App Center](#510-notifications--in-app-center)
  - 5.11 [Admin Dashboard & Dynamic UI Engine](#511-admin-dashboard--dynamic-ui-engine)
  - 5.12 [Operator GPS Tracking](#512-operator-gps-tracking)
  - 5.13 [Insurance Document Management](#513-insurance-document-management)
  - 5.14 [User Onboarding & Help](#514-user-onboarding--help)
  - 5.15 [Search & Discovery Enhancements](#515-search--discovery-enhancements)
  - 5.16 [Accessibility & Internationalisation](#516-accessibility--internationalisation)
  - 5.17 [Offline & Low-Connectivity Support](#517-offline--low-connectivity-support)
  - 5.18 [Error Handling & User Feedback](#518-error-handling--user-feedback)
  - 5.19 [Account Deletion](#519-account-deletion)
  - 5.20 [Policy & Terms Management](#520-policy--terms-management)
  - 5.21 [Feature Flagging & Remote Configuration](#521-feature-flagging--remote-configuration)
  - 5.22 [Global Soft-Delete Policy](#522-global-soft-delete-policy)
  - 5.23 [Ethiopian Public Holiday Calendar](#523-ethiopian-public-holiday-calendar)
  - 5.24 [Marketplace Integrity & Abuse Controls](#524-marketplace-integrity--abuse-controls)
- 6. [Functional Requirements – Phase 2 (Corporate)](#6-functional-requirements--phase-2-corporate)
- 7. [Functional Requirements – Phase 3 (Advanced)](#7-functional-requirements--phase-3-advanced)
- 8. [Booking State Machine & Business Rules](#8-booking-state-machine--business-rules)
- 9. [Non-Functional Requirements](#9-non-functional-requirements)
  - 9.1 [Performance](#91-performance)
  - 9.2 [Availability](#92-availability)
  - 9.3 [Scalability](#93-scalability)
  - 9.4 [Security](#94-security)
  - 9.5 [Maintainability](#95-maintainability)
  - 9.6 [Data Retention & Privacy](#96-data-retention--privacy)
  - 9.7 [API, Logging & QA Development Standards](#97-api-logging--qa-development-standards)
  - 9.8 [Locale & Calendar Support](#98-locale--calendar-support)
  - 9.9 [Accessibility Standards](#99-accessibility-standards)
- 10. [Data Model Overview](#10-data-model-overview)
- 11. [Dynamic Configuration Architecture](#11-dynamic-configuration-architecture)
  - 11.1 [Dynamic Specification Schema Versioning & Backfill Rules](#111-dynamic-specification-schema-versioning--backfill-rules)
- 12. [Technology Stack & Deployment View](#12-technology-stack--deployment-view)
- 13. [MVP Scope & Implementation Roadmap](#13-mvp-scope--implementation-roadmap)
- 14. [Appendices](#14-appendices)

---

## 1. Introduction

### 1.1 Purpose & Audience
This document constitutes the first partition of the canonical Software Requirements Specification for the Dink Machinery marketplace platform. It defines the foundational system context, architectural invariants, user permission models, and complete functional requirements for Identity, Registration, and Session Security.

| Stakeholder Role | Primary Responsibilities & Usage |
|---|---|
| **Product Managers** | Validate feature completeness, prioritise development backlog, and enforce business rules. |
| **System Architects & Agents** | Ensure contract supremacy, zero architectural drift, and strict invariant compliance. |
| **Backend Engineers (Go 1.27.3)** | Implement Clean Architecture interactors, `sqlc` compiled repositories, and Chi routing. |
| **Mobile Engineers (Flutter 3.47.9)** | Build BLoC state machines, Isar offline persistence, and accessible screens. |
| **Frontend Engineers (Next.js 16)** | Build dynamic server-driven admin data tables (`@tanstack/react-table`) and Server Actions. |
| **QA & Test Automation Engineers** | Derive contract tests (`kin-openapi`), integration suites, and ASVS L3 security audits. |
| **Security & Compliance Officers** | Audit zero-PII masking, identity provider integrations, and immutable audit logs. |

### 1.2 Scope & Phases
- **Phase 1 — MVP Foundation (Individual Marketplace):** Core rental lifecycle, 60-min payment hold locks, manual bank-transfer direct settlement (`MANUAL_TRANSFER`), owner-funded `OwnerWallet` commission debits, dynamic category specs, PostGIS radius search, dual-language dynamic i18n, and complete server-driven admin console.
- **Phase 2 — Corporate & Multi-Tenant Accounts:** Organizational entity accounts, tiered approval hierarchies, purchase order allocations, negotiated corporate discount contracts, and cost-center accounting.
- **Phase 3 — Advanced Market Extensions:** Pluggable payment gateway online splits (`PLATFORM_HELD`), direct CAN-bus IoT telematics ingestion, AI demand-based dynamic pricing recommendations, independent operator marketplace, and white-label public APIs.

### 1.3 System Invariants & How to Read This Document
Every requirement is designated with a unique tracking identifier (e.g., `REG-R01-BE`, `AUTH-04-MOB`) and priority level:
- **P0:** Mandatory for baseline release. Bypassing a P0 blocks release gates.
- **P1:** High-priority functional capability scheduled within the active milestone.
- **P2:** Advanced operational feature or future phase capability.

#### Non-Negotiable System Invariants
1. **Contract Supremacy:** Wire models, endpoints, parameters, and RFC 9457 error envelopes are dictated by `contracts/openapi.yaml`. Manual client/server DTO divergence is an immediate CI failure.
2. **Integer Minor-Unit Currency Math:** Floating-point numbers are banned. All monetary attributes are 64-bit signed integers (`int64` / `BIGINT`) representing minor currency units (Santims for `ETB`, Cents for `KES` and `USD`) paired with an ISO 4217 code (`MoneyAmount`).
3. **In-Transaction Ledger Atomicity (`BK-17-BE`):** State transitions mutating funds and double-entry ledger writes (`financial_ledger_entries`) execute inside a single atomic database transaction (`pgx.Tx`), balancing debits and credits to exactly zero minor units.
4. **Optimistic Concurrency Control (OCC):** Stateful tables (`bookings`, `fleets`, `owner_wallets`, `inspections`) enforce `version INT NOT NULL DEFAULT 1` and update via `WHERE id = $1 AND version = $2`.
5. **Soft-Delete Standard (`SD-01-BE`):** Entities declare `deleted_at TIMESTAMPTZ NULL`. Hard deletes are banned from operational queries.
6. **Multi-Country Parameterization (`LOC-01-BE`):** Regulatory parameters (tax basis points, dialing prefixes, currencies, hold timeouts) resolve dynamically from `country_configs`.
7. **OWASP ASVS Level 3 Zero-PII Security:** Sensitive attributes (phone numbers, National IDs, TINs, bank accounts) are encrypted at rest (AES-256-GCM), masked in structured logs, and excluded from client error envelopes.
8. **Universal Envelopes & RFC 9457:** Responses adhere to `SuccessEnvelope` (`success: true`) or RFC 9457 `ErrorEnvelope` (`success: false`).

---

## 2. Glossary & Conventions

| Term | Canonical Definition |
|---|---|
| **Dink Machinery** | The product/brand name of the platform. "Heavy machinery rental aggregator" is the category descriptor for what Dink Machinery is; "Dink Machinery" is the name it goes to market under. |
| **Platform** | The entire Dink Machinery heavy machinery rental aggregator ecosystem (mobile apps, backend, admin, databases, messaging). |
| **Fleet** | A single rentable unit – can be machinery (excavator) or logistics (lowboy trailer). |
| **Attachment** | Add-on equipment (bucket, hydraulic hammer, auger). Can be rented standalone or bundled with a fleet. |
| **Logistics** | Transport vehicles registered under the category "LOGISTICS". |
| **Booking** | A reservation for a specific date/time range, created by a renter and confirmed by the owner after payment. |
| **Prepayment** | The portion of the total rental fee that the renter must pay upfront. |
| **Payment Evidence** | A digital image (photo/screenshot) of a bank transfer receipt uploaded by the renter as proof of payment. |
| **Booking Timeout** | The window (in hours) after booking creation in which the renter must upload payment evidence; otherwise the booking auto-cancels. Configurable per fleet with a **default of 24h**. |
| **Cancellation Penalty** | Pre-defined fee (percentage or fixed amount) charged to the party that cancels after the booking is confirmed. |
| **Commission** | The platform's revenue, a percentage of the rental fee. |
| **State Machine** | A formal set of booking statuses and the allowed transitions between them, enforced by the backend. |
| **Dynamic Configuration** | Any system element stored in the database rather than hard-coded, editable by admins without code deployment. |
| **i18n** | Internationalisation – all user-facing strings are stored as translation keys. |
| **PostGIS** | PostgreSQL extension for geospatial data; used for location-based searches. |
| **RBAC** | Role-Based Access Control – permissions assigned to roles, which are assigned to users. |
| **Corporate Account** | An organisation that registers on the platform and manages multiple members with internal roles and approvals. |
| **KYC** | Know Your Customer – the process of verifying user identity documents. |
| **Idempotency Key** | A unique token sent with mutating API calls to prevent duplicate operations on retry. |
| **SLA** | Service Level Agreement – minimum performance or availability commitment. |
| **Sponsored** | An owner-purchased paid-ranking placement. Owners buy ranking tiers (e.g. Gold/Silver/Bronze) that float their listing higher in search results. Represented by the `IsBoosted` flag. The canonical term for owner-purchased paid ranking is "Sponsored". |
| **Featured** | An admin-curated placement chosen by platform staff (not purchased by owners). Represented by the `Fleet.IsFeatured` / `Fleet.FeaturedUntil` flags. |
| **Ad Campaign** | Banner advertisements (e.g. home-screen or category banners) managed via the admin ad-campaign tools. Distinct from both Sponsored (paid ranking) and Featured (admin-curated ranking). |
| **Settlement Method** | How a renter payment leg (prepayment, final, additional) is executed and verified. Two methods exist (`CW-00-BE`): **`MANUAL_TRANSFER`** — the Direct Transfer Model (bank transfer + evidence upload + confirmation), always available and the only method enabled at MVP launch; **`GATEWAY`** — an online payment through a configured payment-gateway adapter (`PAY-GW-01-BE`). Which methods/gateways are offered at checkout is admin-configured per country/currency (`PAY-GW-02-BE`); the renter chooses among the enabled options. |
| **Payment Gateway** | A pluggable online-payment integration (e.g. TeleBirr, Chapa, CBE Birr, Stripe, M-Pesa) registered as a row in the `PaymentGateway` registry and implemented as a `PaymentGatewayAdapter` (`PAY-GW-01-BE`). Each gateway declares its own **custody model**: `DIRECT_TO_OWNER` (funds settle straight to the owner; platform never touches them) or `PLATFORM_HELD` (platform receives funds, deducts commission at source, pays out the owner — `PAY-GW-04-BE`). Adding a gateway is config + one adapter; no checkout, state-machine, or client changes. |
| **Direct Transfer Model (MVP)** | The `MANUAL_TRANSFER` settlement method and the canonical MVP fund-flow model (see `CW-00-BE`): the renter pays the owner **directly** by bank transfer; the platform **never** holds, escrows, or disburses renter funds. The only platform-held balance is the owner-funded **OwnerWallet**, from which commission is debited at `CONFIRMED`. Wherever this document says "payout" in the manual/direct context, it means the **owner receiving proceeds directly** (renter → owner), **not** a platform disbursement. Refunds and owner→renter penalties are **owner obligations** that the platform only tracks (owner uploads transfer evidence, renter confirms, overdue escalates to Admin / may debit the owner wallet). Deposit is informational-only in MVP. This model remains first-class permanently — gateways are added **beside** it (per market), never replacing it. |
| **Escrow** | In this document "escrow" refers **only** to the manual **prepayment / payment-evidence** step (renter uploads a bank-transfer receipt, owner verifies). It does **not** denote a held-funds/custodial account: under the Direct Transfer Model (MVP) the platform never holds renter funds. Legacy "escrow" wording in LG-03 and the §9.7.3 E2E flow should be read as "prepayment / payment-evidence". |

---

## 3. Product Vision & Strategic Context

### 3.1 Vision Statement
**"To be the trusted digital marketplace for heavy machinery rental in Ethiopia, enabling equipment owners to maximise utilisation and contractors to access the right machine, at the right time, with full transparency."**

### 3.2 Core Value

| For Renters | For Owners | For the Platform |
|---|---|---|
| Single app to discover & compare machinery | New revenue from idle assets | Asset-light, scalable business |
| Transparent pricing, verified owners | Reach wider customer base | Commission per successful booking |
| Streamlined booking & payment proof | Booking & operator management tools | Data-driven market insights |
| Documented inspection & dispute process | Trust through reviews | Foundation for financial services (Phase 3) |

### 3.3 Ethiopian Market Fit
- **Payment Culture:** Bank transfers are the norm; card penetration is low. Mobile money (TeleBirr, CBE Birr) is rising rapidly.
- **Machinery Mix:** Construction, road infrastructure, and agriculture dominate regional capital investments.
- **Logistics:** Often machinery needs a separate transport vehicle; we model logistics as a first-class fleet category.
- **Trust:** Face-to-face handovers are common; digital verification and photo evidence replicate and formalize this trust.
- **Language:** Amharic and regional languages are essential for mass adoption; all strings are i18n-keyed.
- **Connectivity:** Large parts of Ethiopia have intermittent 3G/2G connectivity; the app must gracefully handle poor or no network via local offline persistence.
- **Digital Literacy:** A significant portion of the target user base may be first-time smartphone app users. Onboarding must be frictionless and intuitive.

---

## 4. User Roles & Permissions

### 4.1 Role Definitions

**MVP (Phase 1)**
- **Renter (Individual):** A person who browses, books, and pays for machinery rental.
- **Owner (Individual):** A person who registers machinery/logistics, sets pricing, manages operators, and confirms payments.
- **Operator:** A skilled professional who operates machinery. Assigned by an owner to a fleet.

**Admin / Staff Roles (MVP)** — the four provisioned staff roles (role keys in parentheses; canonical seed roles per Appendix B):
- **Super Administrator (`super_admin`):** Highest privilege; manages roles, permissions, system configuration, and dynamic UI elements.
- **Administrator (`admin`):** Operational management: verifications/approvals, dispute resolution, user management, analytics.
- **Support Agent (`support_agent`):** Handles disputes and user inquiries only.
- **Finance Administrator (`finance_admin`):** Manages commission tracking, owner payouts, and financial reports.

**Phase 2 Additions (Corporate Roles)**
- **Organization Admin (`org_admin`):** Full control over the organization's settings, members, roles, and approval policies.
- **Corporate Requester (`corp_requester`):** A corporate member who can create booking requests on behalf of the organization.
- **Corporate Approver (`corp_approver`):** A corporate member who can approve or reject internal booking requests.
- **Fleet Manager (`fleet_manager`):** A corporate member who manages specific fleets belonging to the organization.

### 4.2 Permission Matrix (MVP & Phase 2)

> **Authoritative source:** The permission-key matrix in **Appendix B** is the canonical, provisioned permission model for the four staff roles (`super_admin`, `admin`, `support_agent`, `finance_admin`). The table below is a feature-level summary; where it and Appendix B differ, **Appendix B wins**. Corporate role columns use the canonical enum keys `org_admin`, `corp_requester`, `corp_approver` (and `fleet_manager` where fleet-scoped grants apply).

| Permission Area | Renter | Owner | Operator | admin | super_admin | support_agent | finance_admin | org_admin | corp_requester | corp_approver |
|---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| Register / Login | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ (invite) | ✅ | ✅ |
| Edit own profile | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| View fleet list & details | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ |
| Create new fleet | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Edit own fleet | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ (fleet_manager) | ❌ | ❌ |
| Create booking (personal) | ✅ | ❌ | ❌ | ❌ | ❌ | — | — | — | — | — |
| Create booking (corporate) | — | — | — | — | — | — | — | ✅ | ✅ | ❌ |
| Approve corporate booking | — | — | — | — | — | — | — | ✅ | ❌ | ✅ |
| Upload payment evidence | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ (corp) | ✅ | ❌ |
| Confirm / decline payment | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ (fleet_manager) | ❌ | ❌ |
| Rate & review | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Approve owner / fleet | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View disputes | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Resolve disputes | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View commission & financial reports | ❌ | ✅ (own) | ❌ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Confirm commission payouts | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Manage categories & specs | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Manage roles & permissions | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View audit log | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Moderate reviews | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |

---

## 5. Functional Requirements — Phase 1 (MVP)

### 5.1 Registration & Authentication

#### 5.1.1 Renter Registration

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| REG-R01-BE | Register API Validation & OTP | P0 | Backend API to validate signup fields (phone format, email format, username uniqueness) and generate/send a 6-digit OTP code. | 1. Receives signup request payload.<br>2. Validates email and Ethiopian phone format (+251).<br>3. Generates 6-digit OTP and initiates SMS/email dispatch. | API validates inputs, rejects duplicates/invalid formats, and triggers OTP. |
| REG-R01-MOB | Register Form UI | P0 | Mobile application UI allowing new users to sign up as renters with personal details. | 1. User taps "Get Started".<br>2. Selects "Rent Machinery".<br>3. Enters name, phone (+251), email, or username and secure password.<br>4. Taps "Send OTP" to submit form to registration API. | Form fields accept values, validate input formats locally, and navigate to OTP screen on success. |
| REG-R02-BE | OTP Verification Engine | P0 | Validate user-submitted OTP codes, track retry attempts, enforce lockout periods, and set expiration. | 1. Expose verification API.<br>2. Verify OTP code matches cache.<br>3. Enforce 5-minute expiration limit.<br>4. Track failed attempts in Redis (max 3, lock account/phone for 1 minute). | OTP matches cache, expires in 5 mins; rejects after 3 failed attempts and applies 1-min lockout. |
| REG-R02-MOB | OTP Verification Screen | P0 | 6-digit OTP input screen in Flutter mobile app. | 1. Show 6-input verification screen.<br>2. User enters OTP code.<br>3. Call verification API on completion.<br>4. Manage timer for "Resend OTP".<br>5. Show error on failed attempt. | Displays 6 input boxes; handles error alerts; transitions to terms or home screen on valid OTP. |
| REG-R03-BE | Terms Agreement Log | P0 | Database persistence of terms and conditions agreement status. | 1. API registers user acceptance.<br>2. Save terms version, status, and agreement timestamp in the user's database record. | Database saves terms acceptance state with correct timestamp. |
| REG-R03-MOB | Terms Agreement UI | P0 | Checkbox element on signup requiring terms and liability waiver acceptance. | 1. Render checkbox for Terms of Service and Liability Waiver.<br>2. "Complete Registration" button remains disabled until checked.<br>3. Tapping links opens PDF viewer or system browser. | Checkbox is mandatory; links load documents; enables registration button only when checked. |
| REG-R04-BE | Create Renter Profile | P0 | API core database operations for creating a new renter profile record on successful validation. | 1. Create renter user record with role set to `renter` in DB.<br>2. Auto-initialize user profile configurations.<br>3. Return authorization tokens to app client. | Database record successfully created; returns standard access/refresh tokens. |
| REG-R04-MOB | Navigate to Home Screen | P0 | Mobile app post-registration workflow routing user to renter home/discovery. | 1. Receive successful registration response from API.<br>2. Store session keys locally.<br>3. Clear auth navigation stack.<br>4. Redirect user to main category/discovery screen. | User profile data loaded; home/discovery screen rendered with default categories. |
| REG-R05-BE | Duplicate Account Validation | P0 | API checks to ensure credentials are not registered to another account. | 1. API query DB on username, email, and phone index.<br>2. Return HTTP conflict error with descriptive code if a duplicate is found. | Registration endpoint returns clean conflict error; prevents saving duplicate rows. |
| REG-R05-MOB | Duplicate Account Redirection | P0 | Mobile app handling of conflict errors. | 1. Detect duplicate account error code from signup API.<br>2. Prompt user with a dialog: "This account already exists. Log in instead?".<br>3. Navigate user to login screen upon confirmation. | Prompt dialog appears on duplicate; user redirected to login. |
| REG-R06-BE | Anonymous discovery endpoints | P1 | Allow unauthenticated search and discovery API calls. | 1. Configure auth middleware to bypass JWT checks for category listings, fleet searches, and details endpoints. | Endpoints return successful HTTP 200 responses without authorization headers. |
| REG-R06-MOB | Guest browsing layout | P1 | Let unregistered guests browse and search fleets. | 1. Allow guest user to access home page and view fleet details.<br>2. If guest taps "Book Now", intercept navigation and present login/registration prompt. | Guest browsing active; triggers authentication flow before booking. |
| REG-R07-BE | Fayda ID Verification Client (reference identity adapter) | P0 | The Fayda registry integration, implemented as the reference `IdentityVerificationProvider` adapter under the REG-IDP-01-BE framework (not a bespoke one-off client). | 1. Backend connects to Fayda API via the adapter interface.<br>2. Submit user Fayda ID.<br>3. Retrieve official name, birthdate, and photo as `verified_attributes`.<br>4. Fall back to manual input with that provider's `IdentityVerification.status = UNVERIFIED` if Fayda offline.<br>5. Persists the user's informed-consent record (timestamp, accepted policy version, and a reference to the Fayda request/response envelope) captured by the REG-R07-MOB consent checkbox, satisfying the §9.6 Legal Audit Consent Logging mandate. Fayda-specific endpoint/credentials/field-map live in its `IdentityProvider.config_json`, not in code. | Implements the `IdentityVerificationProvider` interface; syncs data automatically when Fayda is active and enabled; a consent record (timestamp + policy version + envelope reference) is persisted for every submission attempt; Fayda connection settings are config-driven. |
| REG-R07-MOB | Identity Verification Forms (provider-driven) | P0 | Renter/any-role app screens to verify identity against a configured provider; provider list and labels are driven by config, not hardcoded to Fayda. | 1. User selects a verification provider (when only Fayda is enabled this is a single "Verify with Fayda" option; the list is rendered from the enabled `IdentityProvider` set).<br>2. Displays an explicit provider-named informed-consent checkbox — e.g. "I authorize Dink Machinery to verify my identity against the Fayda registry" — which must be checked before submit is enabled.<br>3. Input the provider identifier.<br>4. Fetch pre-filled information from the provider's `verified_attributes`.<br>5. If verification is unavailable, fall back to manual text form entry. | Provider options render from enabled-provider config (not a hardcoded Fayda-only screen); submission is blocked until the consent checkbox is checked; pre-fills name/photo/birthdate when the provider returns them; falls back to manual forms. |
| REG-R07-WEB | *(superseded)* Fayda settings — see REG-IDP-03-WEB | P0 | The former standalone "Enforce Fayda Sync" toggle is superseded by the unified Identity Provider & Enforcement Admin Console. | Fayda is now enabled/configured as one provider row in REG-IDP-03-WEB (set `is_enabled`, `enforcement_role`, and `config_json`); there is no separate Fayda-only settings screen. | Retained for traceability; all Fayda admin configuration is performed via REG-IDP-03-WEB. |
| REG-R08-BE | Identity Verification Outcome & Retry Handling | P0 | Provider-agnostic backend handling of identity-provider failure modes — not-found/mismatch and API unavailability — with bounded retries and an admin manual-verification fallback. Shared by every provider adapter (REG-IDP-01-BE), Fayda included. | 1. If a provider returns a not-found or data-mismatch result, respond with error code `IDENTITY_MISMATCH`, keep that provider's `IdentityVerification.status = UNVERIFIED`, and permit at most 3 retry submissions per user per provider per rolling 24h window (further attempts return HTTP 429).<br>2. If the provider API times out or returns 5xx, retry with exponential backoff; after N consecutive failures (configurable, default 3), set the verification to `PENDING_REVIEW` and enqueue the user in the existing admin manual-verification review queue (REG-O02-WEB-style queue).<br>3. On retry exhaustion, notify the user their submission is under manual review.<br>4. Every provider request/response envelope is logged per the §9.6 Legal Audit Consent Logging mandate, including on failure paths. | Mismatch responses never move a provider verification to verified and are rate-limited to 3/24h/provider; transient API failures fail over to `PENDING_REVIEW` + manual queue after N retries; every provider envelope (success or failure) is logged; behaviour is identical across all configured providers. |
| REG-R08-MOB | Identity Verification Error & Retry UI | P0 | Mobile UI surfacing distinct error states for each provider outcome; copy is parameterised by the active provider's display name. | 1. Renders a specific message for `IDENTITY_MISMATCH` ("The details don't match the {provider} records — please check and try again") distinct from a generic API-unavailable message.<br>2. Displays a "X of 3 attempts remaining today" counter that decrements per retry.<br>3. Once retries are exhausted or the verification enters `PENDING_REVIEW`, replaces the retry form with a "Submitted for manual review" state and a manual-fallback contact/support option. | Each outcome shows a distinct, actionable message with the correct provider name; remaining-retry counter is accurate; exhausted retries route to the manual-review state instead of a dead-end error. |
| REG-N01-BE | Role-Agnostic Identity Verification Capture | P0 | Backend endpoint allowing any authenticated user (renter, owner, or operator) to submit identity verification to any configured provider, reusing the pluggable provider framework (REG-IDP-01-BE) regardless of role. | 1. Exposes `POST /user/identity/verify` (accepts a `provider_key`, the provider-specific identifier, and the consent record) for any authenticated user who has not yet satisfied the configured identity policy for that provider.<br>2. Delegates to the selected provider adapter and shared outcome handling (REG-IDP-01-BE / REG-R08-BE).<br>3. Persists an `IdentityVerification` row against the requesting user's record regardless of role (renter/owner/operator); a user may hold verifications from several providers at once.<br>4. When the identity lock policy is active (AUTH-05-BE), this is the single endpoint through which any role satisfies each still-required provider — including owners/operators who have no renter-registration verification step. | Endpoint accepts requests from any role and any configured `provider_key`; each verification result is saved as an `IdentityVerification` against the correct user record; users clear the AUTH-05 lock by satisfying every still-missing required provider exclusively via this endpoint. |
| REG-IDP-01-BE | Pluggable Identity Verification Provider Framework | P0 | A provider-agnostic interface + registry so identity providers (Fayda and any future provider) are plug-and-use via configuration, with no enforcement/capture/client code changes to add one. | 1. Define an `IdentityVerificationProvider` adapter interface: `submit(identifier, consent) -> { status, verified_attributes (name/DOB/photo/etc.), provider_ref, raw_envelope }`, plus outcome/retry handling shared with REG-R08-BE.<br>2. Providers are rows in a new `IdentityProvider` table: `provider_key` (e.g. `fayda`), `display_name`, `adapter_type`, `config_json` (endpoint/credentials/field-map), `is_enabled`, `enforcement_role ∈ {required, optional, disabled}`, `priority`.<br>3. Fayda (REG-R07-BE) is the reference adapter implementing this interface.<br>4. Adding a new provider = implement the adapter + insert/enable an `IdentityProvider` row via REG-IDP-03-WEB — the enforcement middleware (AUTH-05), the capture endpoint (REG-N01), and mobile clients are untouched.<br>5. Every attempt persists an `IdentityVerification` row (`user_id`, `provider_key`, `status`, `provider_ref`, `verified_attributes`, `envelope_ref`, `consent_ref`, `created_at`) so multiple providers can be verified per user. | New provider becomes usable by config alone (adapter + enabled row) with no change to AUTH-05/REG-N01/client code; each attempt stored as an `IdentityVerification`; a user can simultaneously hold verified results from more than one provider. |
| REG-IDP-02-BE | Identity Enforcement Policy Engine | P0 | Evaluates whether a user satisfies the admin-configured identity policy, supporting single-provider, any-of, all-of, and N-of-M multi-provider requirements. | 1. Reads `identity.enforcement_mode ∈ {OFF, ANY, ALL, N_OF_M}` (SystemConfig) and each provider's `enforcement_role` (IdentityProvider).<br>2. Semantics over the set of providers marked `required`: **OFF** = no identity requirement (verification is optional/informational); **ANY** = ≥1 required provider verified; **ALL** = every required provider verified; **N_OF_M** = at least `identity.min_required_providers` of the required set verified.<br>3. Returns `{ satisfied: bool, missing_providers: [...], mode }`.<br>4. Derives the user's overall `identity_status`: `VERIFIED` when the policy is satisfied, `PENDING_REVIEW` if any required provider is mid-review, else `UNVERIFIED`.<br>5. This structured result — not any single hardcoded Fayda flag — is the sole input AUTH-05-BE consumes. | Policy engine returns the correct satisfied/missing set for each mode; requiring two providers (ALL over two required providers) correctly blocks a user who has verified only one; `identity_status` is derived from the policy result, not a single provider. |
| REG-IDP-03-WEB | Identity Provider & Enforcement Admin Console | P0 | Single Super-Admin console (replacing the former separate "Enforce Fayda Sync" and "Enforce National ID Lock" toggles) to configure providers and the enforcement policy dynamically. | 1. List all registered providers; per provider: enable/disable, set `enforcement_role` (required/optional/disabled), edit `config_json` (connection/credentials/field-map), set `priority`.<br>2. Set `identity.enforcement_mode` (OFF/ANY/ALL/N_OF_M) and `identity.min_required_providers`.<br>3. Toggle `identity.lock_enabled` — whether policy-unsatisfied users are blocked from transactional endpoints (AUTH-05) or merely flagged.<br>4. **Reject invalid combinations**: `lock_enabled=true` with `enforcement_mode=OFF`; `ALL`/`N_OF_M`/`ANY` with zero enabled `required` providers; `min_required_providers` greater than the count of required providers; marking a disabled provider as required.<br>5. Guarded by `config:manage` + admin MFA; every change written to `AuditLog`; changes take effect immediately in the middleware. | Admin can add/enable a provider, mark one or more providers required, choose any/all/N-of-M, and toggle the lock — all without code changes; contradictory configurations are rejected with a clear error; changes are audit-logged and apply immediately. |
| REG-R09-BE | Email Ownership Verification API | P0 | Backend flow proving ownership of the email address on file. | 1. On registration (if an email was provided) and on any email change (AUTH-15), generates a signed, single-use verification link/code (24h expiry) and sends it to the address.<br>2. `POST /auth/verify-email` validates the token and sets `email_verified = true` on the user record.<br>3. Email-based login (where supported) and email-based password-reset flows are gated on `email_verified = true`.<br>4. `POST /auth/verify-email/resend` re-sends the link, rate-limited per the §9.4 OTP/notification throttling conventions. | Unverified emails cannot be used for login or password reset; verification link/code expires after 24h; resend is rate-limited. |
| REG-R09-MOB | Email Verification UI | P0 | Mobile UI prompting and tracking email verification. | 1. After signup (or email change) with an unverified email, displays a persistent banner: "Verify your email — check your inbox".<br>2. Handles the verification deep link, calling the confirm endpoint and updating the banner state on success.<br>3. Banner includes a "Resend" action with a visible cooldown timer matching the backend rate limit. | Banner shows until verified; deep link completes verification without manual re-login; resend respects cooldown. |
| REG-R10-BE | Minimum-Age Eligibility Enforcement | P0 | Backend enforcement of a minimum account-holder age derived from verified (or manually-entered) date of birth. | 1. When any identity provider returns a DOB (REG-R07/REG-N01 via the provider framework), computes age; if age is below `identity.min_age` (SystemConfig key, default 18), blocks account creation/verification with error code `UNDER_MINIMUM_AGE`.<br>2. Applies the same DOB-derived gate to the manual-fallback entry path (when a provider is offline), additionally flagging the record for admin review since the DOB is self-reported/unattested.<br>3. Booking creation (`bookings:create`) is blocked for any account that failed or has not passed the age gate, independent of its identity-verification status. | Accounts with computed age below `identity.min_age` cannot complete registration/verification; manually-entered DOB under the threshold is flagged for review rather than silently trusted; booking creation is blocked for age-ineligible accounts. |
| REG-R10-MOB | Minimum-Age Ineligibility UI | P0 | Mobile messaging for age-ineligible registration attempts. | 1. If registration/verification returns `UNDER_MINIMUM_AGE`, displays a clear ineligibility screen explaining the minimum-age requirement rather than a generic error.<br>2. Once DOB has been verified by a provider, the date-of-birth field becomes read-only everywhere in the app (profile, edit-profile) and is no longer user-editable. | Ineligible users see an explanatory (non-generic) message; provider-verified DOB cannot be edited from any app screen. |

---

#### 5.1.2 Owner Registration

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| REG-O01-BE | Register Owner Profile API | P0 | Backend endpoint to process owner profiles and save associated KYC documents. | 1. Accept owner form payload (personal details, company, bank info).<br>2. Store business license and insurance PDF/image links.<br>3. Set status to `PENDING_REVIEW`. | DB record created in `PENDING_REVIEW` state; files correctly stored in file server/S3. |
| REG-O01-MOB | Owner Onboarding Screen | P0 | Flutter form for entering business, personal, and financial details. | 1. User inputs owner details.<br>2. Pick and upload license document (max 10MB).<br>3. Select bank and input routing/account details.<br>4. Upload insurance files.<br>5. Submit to register. | Mobile client validates field lengths, formats file size, and sends files to API. |
| REG-O02-BE | Owner Registration Approval API | P0 | Endpoints for changing owner status based on verification reviews. | 1. Expose admin endpoint to update owner status.<br>2. If approved, set status `ACTIVE`.<br>3. If rejected, set status `REJECTED` and store comment reason.<br>4. Send notification. | Owner state changes correctly; rejection comments stored in DB; email/push alerts triggered. |
| REG-O02-WEB | Admin Owner Review Queue | P0 | Next.js portal interface for administrators to verify new owners. | 1. Admin opens review queue page.<br>2. View owner files (business license, insurance).<br>3. Click "Approve" or "Reject" (requiring text explanation). | Portal displays files inline; action buttons issue status change requests correctly. |
| REG-O02-MOB | Verification Status Notification | P0 | Dashboard and push warning widgets displaying approval state to owner. | 1. Owner opens app dashboard.<br>2. Retrieve verification status.<br>3. Render status banner ("Pending Review", "Active", or "Rejected" with reasons). | Owners receive push update on status changes; rejection explanation renders on failure. |
| REG-O03-BE | Expose Bank Info API | P0 | Supply owner bank details to authorized renter requests during checkout. | 1. Fetch bank information associated with booking's fleet owner.<br>2. Expose bank info in booking detail endpoint when status is `PENDING_PAYMENT`. | Bank details included in response only for authorized renter checkouts. |
| REG-O03-MOB | Display Bank Coordinates | P0 | Render owner banking info on renter's booking checkout details. | 1. Renter navigates to payment screen.<br>2. Retrieve bank details (holder, account, bank name).<br>3. Render info with a "Copy Account" button. | Data is shown correctly; clipboard copy copies the correct account number. |
| REG-O04-BE | Re-submit Application API | P0 | Accept updated documents and reset owner status back to PENDING_REVIEW. | 1. Allow PATCH on owner details if status is `REJECTED`.<br>2. Save new files and set status to `PENDING_REVIEW`.<br>3. Generate validation task for admins. | Application updates database; sets status back to `PENDING_REVIEW` without generating new profiles. |
| REG-O04-MOB | Re-submission Interface | P0 | Re-upload documents for rejected owner accounts. | 1. Owner app detects `REJECTED` status.<br>2. Display "Re-submit" button.<br>3. Pre-fill previous details, open file picker to replace invalid documents, and send update. | UI shows rejection reason, lets user pick new files, and runs update call. |
| REG-O04-WEB | Re-submitted Review Notification | P0 | Admin panel updates for re-submitted owner profiles. | 1. Display "Re-submitted" badge on registration queue items.<br>2. Show history of prior rejection reasons next to new files. | Admin sees updated details and historical rejection notes in approval dashboard. |
| REG-O05-BE | KYC Expiry Checking Job | P1 | Scheduled cron task verifying document dates and applying status freezes. | 1. Nightly cron job calculates time differences between current date and license expiry.<br>2. If 60 days left, dispatch alerts.<br>3. If expired, set status `DOCUMENT_EXPIRED`, pause active listings. | Automated cron updates database flags; blocks new listing creations; triggers warning emails. |
| REG-O05-MOB | KYC Expiry Warning Banners | P1 | Warnings on owner dashboard regarding approaching expiry. | 1. Check owner license expiry threshold.<br>2. If close, show banner warnings on home screen.<br>3. If expired, lock listing management and show document upload modal. | Prominent warning banners show remaining days; user blocked from creating listings if expired. |
| REG-O05-WEB | KYC Expiry Admin Alert Log | P1 | Next.js admin page tracking expired and expiring files. | 1. Show alerts panel of owners with expired licenses or licenses expiring in under 60 days.<br>2. Offer filters for quick admin action or contacting owners. | Admin dashboard lists expiring owners with contact info. |

---

#### 5.1.3 Operator Registration

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| REG-OP01-BE | Operator Registration API | P0 | Endpoint to save operator details, documents, skill taxonomy, and account credentials. | 1. Receive name, contact, profile image, certification uploads, skill tags, and a password meeting the AUTH-07 complexity policy (operators use the same username/password + MFA-priority login model as AUTH-04; there is no separate passwordless login path).<br>2. Run the REG-R05 duplicate phone/email/username check before creating the record; reject with the same conflict error on a match.<br>3. Record terms-of-service/liability-waiver acceptance (version + timestamp) following the REG-R03 pattern.<br>4. Save profile records.<br>5. Set status to `ACTIVE`. | Operator account requires and stores a complexity-valid password usable via AUTH-04 login; duplicate phone/email/username is rejected via the REG-R05 check before any record is created; terms acceptance is logged with version and timestamp identically to REG-R03; DB record created; skills linked in join table; files persisted. |
| REG-OP01-MOB | Operator Signup Form | P0 | Flutter form interface for operators to list skills and upload certificates. | 1. Operator inputs details.<br>2. Selects skill tags from dynamic checklist.<br>3. Uploads certification files.<br>4. Submits form. | User fills name, picks skill checkboxes, uploads files, and receives active screen. |
| REG-OP02-BE | Owner-Operator Association API | P0 | Endpoint to create team invitations and map relationships upon acceptance. | 1. Owner sends invitation with operator phone.<br>2. Dispatch push notification.<br>3. Expose accept/decline endpoint.<br>4. If accepted, create `OwnerOperatorLink` record. | Invite model created in DB; link is written upon acceptance; notifications sent. |
| REG-OP02-MOB | Invitation Acceptance Flow | P0 | Owner dashboard invitation actions and operator push alert dialogs. | 1. Owner inputs phone number in "My Operators".<br>2. Operator gets alert in-app.<br>3. Operator clicks "Accept" or "Decline".<br>4. Link refreshed in owner's team. | Owner sends invite easily; operator responds in app; owner's list updates instantly. |
| REG-OP03-BE | Operator Certification Checker | P1 | Cron job tracking operator certificate validity and flagging expired entries. | 1. Check certification dates nightly.<br>2. Send warning notifications 30 days before expiry.<br>3. If expired, flag skill tag `UNVERIFIED` and reject assignment in booking validation. | Automatic status updates in DB; API rejects expired operator bookings. |
| REG-OP03-MOB | Expiry Warning Displays | P1 | Mobile banners showing expiring certificates to operators and linked owners. | 1. Operator receives warning notification and banner on home.<br>2. Linked owners view warnings on their operator roster screens. | Warning indicators visible on active profiles; operator prompted to upload new file. |

---

#### 5.1.4 Authentication & Session Management

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| AUTH-01-BE | Session Token Generation API | P0 | Generate access JWTs (15 min) and refresh tokens (30 days), and handle refresh requests. | 1. API issues access token and refresh token on auth.<br>2. Verify and swap refresh token at `/auth/refresh` endpoint to generate new access token. | Tokens generated using secure signature; invalid/expired refresh tokens return 401. |
| AUTH-01-MOB | Secure Token Lifecycle Management | P0 | Secure local storage and silent token rotation on mobile client. | 1. Store refresh token in Flutter Secure Storage.<br>2. Detect 401 error in network interceptor.<br>3. Call `/auth/refresh`.<br>4. Retry original query.<br>5. If expired, send to login. | Silent token refreshing requires no user input; clean login redirect if expired. |
| AUTH-01-WEB | Web Token Cookie Middleware | P0 | Secure HttpOnly cookie handling and refresh middleware for Next.js web application. | 1. Web client receives tokens.<br>2. Store refresh token in HttpOnly, Secure, SameSite=Strict cookie.<br>3. Next.js middleware refreshes access token before page loads. | Cookie settings prevent XSS theft; auto-refresh keeps dashboard session active. |
| AUTH-02-BE | Dynamic Concurrent Session Enforcement | P1 | DB-tracked session constraints dynamically resolved from SystemConfig keys. | 1. Track user sessions in `UserActiveSession` table.<br>2. Read SystemConfig keys: `session.allow_multi_session` (default `true` if unconfigured), `session.max_sessions_user` (default 5), and `session.max_sessions_admin` (default 2).<br>3. If `allow_multi_session` is false, revoke all previous sessions on new login.<br>4. If true, verify session count against role-specific limit (`max_sessions_admin` for admins, `max_sessions_user` for standard users) and revoke the oldest session when exceeded.<br>5. Idle-session timeout: if a session's refresh token goes unused for `session.idle_timeout_days` (default 7), invalidate that session even if the user is under the session-count cap.<br>6. Absolute session lifetime cap: regardless of activity, force-invalidate a session once it exceeds `session.absolute_lifetime_days` (default 90) from creation.<br>7. Step-up re-authentication: sensitive operations — bank/payout-detail edits, and MFA/password/phone/email changes (cross-ref AUTH-06, AUTH-12, AUTH-15) — additionally require the user to have re-verified their current password or a valid current MFA factor within the last 15 minutes, even within an otherwise-valid, non-expired session; requests failing this freshness check are rejected and prompted to re-authenticate. | Session limits dynamically enforced from SystemConfig settings (with documented defaults: `max_sessions_user`=5, `max_sessions_admin`=2, `allow_multi_session`=true); database updates revoke sessions instantly; idle sessions unused for `session.idle_timeout_days` (default 7) are invalidated independent of the count cap; sessions older than `session.absolute_lifetime_days` (default 90) are force-invalidated regardless of activity; sensitive operations (bank/payout edits, MFA/password/phone/email changes) are rejected unless the current password or MFA factor was re-verified within the last 15 minutes. |
| AUTH-02-MOB | Session Termination Handling | P1 | Clear local state when active session is revoked. | 1. Mobile app receives token invalidation error (session revoked).<br>2. Clear credentials and local storage.<br>3. Redirect to login.<br>4. Show: "Session expired on another device". | Revoked devices automatically logout and display descriptive notification. |
| AUTH-03-BE | Brute Force Lockout Engine | P0 | Redis-based counting of failed attempts and automated locking. | 1. On failed login, increment Redis counter (30 min TTL).<br>2. If attempts reach 10, lock account for 30 mins.<br>3. Reject logins with HTTP 429.<br>4. Expose unlock API. | Account locked for exactly 30 minutes; returns HTTP 429; admin unlock API functions. |
| AUTH-03-WEB | Admin Account Unlock Dashboard | P0 | Account management features in Next.js admin portal. | 1. Admin views locked user accounts in panel.<br>2. Click "Unlock Account".<br>3. Trigger unlock API, resetting Redis counters. | Locked status visible; unlock actions successfully clear login count blocks. |
| AUTH-03-MOB | Lockout Warning UI | P0 | Display lockout errors and cooldown counts. | 1. Login screen receives HTTP 429.<br>2. Block submit button.<br>3. Display countdown timer showing remaining minutes before lock releases. | Displays user-friendly error stating lockout time; blocks retry input. |
| AUTH-04-BE | Optional MFA Validation API | P0 | Dual-stage auth validation where MFA is optional. Enforces user-configured factor priority and dynamic switching. | 1. Verify username/password.<br>2. If valid, check user profile's `mfa_enabled` flag.<br>3. If false, immediately return full session JWT.<br>4. If true, return temporary pre-auth state with the user's highest priority registered factor (defaulting to Passkey -> Authenticator App -> SMS priority).<br>5. Process `/auth/mfa` verification.<br>6. Expose `/auth/mfa/switch` to return challenge details for alternative registered factor.<br>7. Issue full session JWT upon successful verification of any factor. | Full JWT only issued after password verification + secondary factor verification (if MFA enabled). Default priority and factor switching API are functional. |
| AUTH-04-MOB | MFA Login & Dynamic Switching UI | P0 | Flutter interface for progressive MFA code verification and dynamic switching. | 1. User enters password.<br>2. If MFA is active, app displays input/challenge screen for user's highest priority factor.<br>3. Screen displays a "Try another method" hyperlink.<br>4. Tapping it opens a bottom sheet listing alternate registered factors.<br>5. Selecting one calls the switch API and redraws the verification input screen for that factor. | Displays highest-priority factor screen immediately after password; "Try another method" allows seamless factor swap. |
| AUTH-04-WEB | Web Admin Mandatory MFA Login | P0 | Enforces mandatory MFA for all admin accounts. User settings/preferences are bypassed. | 1. Admin enters credentials.<br>2. On first login, prompt forced setup of TOTP via QR code scan.<br>3. For subsequent logins, require TOTP verification immediately following password. | Admin accounts blocked from accessing dashboard unless TOTP verification passes. |
| AUTH-04B-BE | MFA Preference & Reordering API | P0 | Endpoints to enable/disable MFA and customize factor priority order. | 1. Expose PUT `/user/security/mfa` to update `mfa_enabled` flag.<br>2. When the request sets `mfa_enabled=false`, requires a step-up re-authentication challenge (current password + a valid current MFA factor) submitted with the request; reject with 403 if the challenge is missing or fails.<br>3. If the requesting account's role is `admin` or `super_admin`, reject any `mfa_enabled=false` request outright with 403 regardless of challenge, per §9.4's mandatory-admin-MFA rule.<br>4. On successful disable (non-admin roles only), dispatches a security notification (MFA-disabled alert) to the user's registered phone/email.<br>5. Expose PUT `/user/security/mfa/priority` to save user's custom priority list (ordered array, e.g., `["totp", "passkey", "sms"]`) in `mfa_priority` column. | API updates DB preferences and validates that all priority elements are valid; disabling MFA without a valid current-password + current-factor challenge is rejected (403); Admin/Super Admin `mfa_enabled=false` requests are always rejected (403) at the API layer irrespective of challenge; a security notification is sent whenever a non-admin user disables MFA. |
| AUTH-04B-MOB | Security Settings & MFA Ordering UI | P0 | Mobile security settings screen to manage MFA toggles and reorder factor priority. | 1. User navigates to Settings -> Security Settings.<br>2. Displays toggle for "Enable Multi-Factor Authentication".<br>3. If enabled, displays a re-orderable list (drag-and-drop or rank pickers) showing: Passkey, Authenticator App, SMS.<br>4. On reorder, triggers priority update API. | User can enable/disable MFA and customize factor priority order; changes persist to backend. |
| AUTH-04B-WEB | Web Admin Security Settings | P0 | Web settings for admin account security settings. | 1. Admin navigates to settings.<br>2. Edits factor priority preference.<br>3. Saves to backend. Admin and Super Admin roles have no UI control to disable MFA (`mfa_enabled=false`) — per §9.4's mandatory-admin-MFA rule, only factor priority/ordering may be reconfigured; the disable action is omitted from the admin settings UI, and the underlying AUTH-04B-BE endpoint rejects any admin/super_admin request attempting to set `mfa_enabled=false`. | Supports admin-level MFA factor-priority customization only; Admin and Super Admin accounts cannot set `mfa_enabled=false` — the disable action is unavailable/blocked for these roles at the API layer, not just hidden in UI. |
| AUTH-05-BE | Identity Enforcement Block Middleware | P0 | API enforcement that blocks transactional endpoints when a user does not satisfy the admin-configured identity policy — provider-agnostic, driven entirely by the REG-IDP-02-BE policy engine. | 1. Ask the REG-IDP-02-BE policy engine whether the user is `satisfied` (never read any single hardcoded provider flag).<br>2. If `identity.lock_enabled` is true and the user is not satisfied (`identity_status` = `UNVERIFIED` or `PENDING_REVIEW`), block "transactional" endpoints/permission scopes, enumerated as: `bookings:create`, `payments:*`, `fleets:create`, `wallet:topup`, and bank/payout-detail edits.<br>3. Explicitly permitted while blocked: browse/search listings, view own profile, complete any identity verification itself, and read-only viewing of booking history.<br>4. Applies uniformly to all roles (renter, owner, operator); a user unblocks by satisfying every still-missing required provider via REG-N01-BE.<br>5. The 403 response includes the `missing_providers` set so clients can prompt the exact remaining step(s). | Middleware consumes the policy-engine result (not a Fayda flag); blocks the enumerated transactional scope while unsatisfied; browse/search, own-profile, identity completion, and read-only history stay accessible; multi-provider policies (e.g. two required providers) block a user who has satisfied only one; response lists the missing providers. |
| AUTH-05-MOB | Identity Lock Overlay | P0 | App navigation lock restricting users who don't satisfy the identity policy. | 1. App checks `identity_status` and `missing_providers` post-login.<br>2. If the lock is active and the user is unsatisfied, render a full-screen modal blocking dashboard access, listing each remaining required provider as an actionable verification step. | Navigation is locked; modal cannot be closed until the policy is satisfied; the modal shows every outstanding required provider, not a single hardcoded Fayda step. |
| AUTH-05-WEB | *(superseded)* Identity lock toggle — see REG-IDP-03-WEB | P0 | The former standalone "Enforce National ID Post-Login Lock" toggle is superseded by the unified Identity Provider & Enforcement Admin Console. | The `identity.lock_enabled` toggle, enforcement mode, and per-provider `required/optional/disabled` settings are all managed in REG-IDP-03-WEB. | Retained for traceability; all identity-lock configuration is performed via REG-IDP-03-WEB. |
| AUTH-06-BE | Secure Phone Number Change API | P0 | Endpoints to request and verify phone number changes with user credentials, multi-factor auth, and OTP. | 1. Receives change request at `POST /user/change-phone` with current password, current MFA token, and new phone number.<br>2. Verifies credentials, generates 6-digit OTP (5 min expiry), and dispatches OTP via SMS to the new number.<br>3. Receives OTP validation at `POST /user/change-phone/confirm`.<br>4. Instantly updates user's phone, archives old number in `UserPhoneHistory`, and sends notification SMS to old phone number. | New phone number takes effect instantly; old phone gets notify SMS; history logged. |
| AUTH-06-MOB | Secure Phone Change UI | P0 | Mobile settings interface for requesting and confirming phone number updates. | 1. User taps "Change Phone" in Security settings.<br>2. Inputs password, MFA code (if active), and new phone.<br>3. Receives OTP input page.<br>4. Enters OTP and confirms change. | Secure wizard prompts for password, MFA, and OTP; triggers success alerts. |
| AUTH-06-WEB | Support Revert Phone Change Console | P0 | Admin override endpoint and dashboard view to rollback fraudulent phone number changes. | 1. Endpoint `POST /admin/users/{id}/revert-phone` receives target userId, requiring Admin MFA + 100-character justification.<br>2. Reverts user phone to the previous record in `UserPhoneHistory`.<br>3. Force-logouts all active sessions for target user.<br>4. Logs event in `AuditLog`. | Rollback restores old phone; terminates sessions; logs admin justification immutably. |
| AUTH-07-BE | Password Complexity & Forgot Password API | P0 | Backend validation checks for password complexity, and phone OTP based password reset flow. | 1. Enforces password strength rules during registration and reset (min 8 chars, 1 uppercase, 1 lowercase, 1 digit, 1 special character).<br>2. Endpoint `POST /auth/forgot-password` initiates reset: checks phone and triggers SMS OTP.<br>3. Endpoint `POST /auth/forgot-password/reset` verifies OTP (and MFA if enabled) and updates password. | Complexity constraints fail invalid inputs; reset flow updates password without support agents. |
| AUTH-07-MOB | Forgot Password UI Wizard | P0 | Mobile application wizard interface for user password reset. | 1. Login screen renders "Forgot Password" link.<br>2. Prompts user for registered phone.<br>3. Verifies OTP and MFA input fields.<br>4. Collects new password, validating complexity rules inline. | App UI validates strength criteria; OTP/MFA checks succeed before reset API call. |
| AUTH-07-WEB | Web Forgot Password & Complexity UI | P0 | Next.js admin and user web pages containing password complexity checks and reset wizards. | 1. Web forms display password strength indicator.<br>2. Resolves forgot password request pages. | Enforces complexity and resets password. |
| AUTH-08-BE | MFA Backup Recovery Codes API | P0 | Backend generation and verification of one-time MFA backup recovery codes. | 1. During MFA (TOTP) onboarding, generates 8 hashed alphanumeric backup codes, saving them to `UserMfaBackupCode`.<br>2. Login endpoint `POST /auth/mfa/recovery` verifies backup code.<br>3. If correct, bypasses TOTP, updates session, and marks code as `USED`.<br>4. Exposes `POST /user/security/mfa/recovery/regenerate`, requiring the current password and a valid current MFA factor challenge; on success, revokes (marks `USED`/voided) all prior unused codes and issues a fresh set of 8.<br>5. The profile/security-settings response includes the count of remaining unused backup codes.<br>6. When the remaining unused-code count reaches ≤2, dispatches a low-codes-remaining notification (push/email) prompting the user to regenerate. | User can authenticate using backup codes; used codes are permanently voided; regeneration requires re-auth (current password + valid current MFA factor); old codes are voided on regeneration; a low-count warning notification fires at ≤2 remaining codes. |
| AUTH-08-MOB | MFA Backup Code Onboarding & Login UI | P0 | Mobile UI displays recovery codes during TOTP setup and provides backup login challenges. | 1. Setup screen shows list of 8 backup codes with "Copy" and "Download" buttons.<br>2. MFA screen displays "Use Recovery Code" link.<br>3. User inputs code to log in.<br>4. Security settings screen displays the current count of remaining unused backup codes.<br>5. A "Regenerate codes" action prompts a confirmation dialog requiring password/MFA re-entry, then calls the regenerate endpoint and displays the new 8 codes once (with "Copy"/"Download" options as in onboarding). | Displays codes once during onboarding; fallback code bypass works on MFA screen; remaining-code count is visible in security settings; regeneration requires confirmed password/MFA and shows the new codes exactly once. |
| AUTH-09-BE | Passkey (WebAuthn) Registration & Assertion API | P0 | Backend endpoints enabling passkey enrollment and passwordless/MFA assertion using the WebAuthn standard. | 1. `POST /user/security/passkey/register/begin` returns a WebAuthn creation challenge (rp, user, pubKeyCredParams, excludeCredentials).<br>2. `POST /user/security/passkey/register/finish` verifies the attestation response and stores the credential (`credential_id`, `public_key`, `sign_count`, `aaguid`, `device_label`) in `UserWebAuthnCredential`.<br>3. `POST /auth/passkey/assert/begin` and `/auth/passkey/assert/finish` issue and verify an assertion challenge during login or MFA.<br>4. Every assertion enforces a `sign_count` regression check, rejecting the credential if the returned counter does not exceed the stored value (cloned-authenticator detection). | Passkey enrolls with verified attestation and is stored against the user; assertion succeeds during login/MFA; an assertion from a cloned authenticator (regressed `sign_count`) is rejected. |
| AUTH-09-MOB | Passkey Enrollment & Login UI | P0 | Mobile Settings screens to register, label, list, and remove passkeys, and to use them during login. | 1. Settings -> Security shows "Add Passkey", invoking the platform authenticator (biometric-backed) via the WebAuthn API.<br>2. On success, prompts the user to label the credential (e.g., "My Phone").<br>3. Lists all enrolled passkeys with device label and per-item "Remove" action.<br>4. During login, if passkey is the user's highest-priority MFA factor (AUTH-04B), the app presents the platform authenticator challenge directly. | User can enroll, name, and remove passkeys; passkey login completes successfully without falling back to SMS or TOTP. |
| AUTH-10-BE | User TOTP Enrollment API | P0 | Backend endpoints allowing any end user (not just admins) to self-enroll a TOTP authenticator app. | 1. `POST /user/security/totp/begin` generates a TOTP secret, stores it as pending (unverified), and returns an `otpauth://` provisioning URI plus QR payload.<br>2. `POST /user/security/totp/confirm` requires a valid current TOTP code generated from the pending secret before the secret is activated on the user's account.<br>3. The `mfa_secret` column is encrypted at the column level (AES-256-GCM, per §9.4 Security). | The secret is only activated after one valid confirmation code is verified; unverified pending secrets expire; `mfa_secret` is stored encrypted at rest. |
| AUTH-10-MOB | TOTP Enrollment UI | P0 | Mobile Settings -> Security screen to add an authenticator app. | 1. User taps "Add Authenticator App".<br>2. Screen renders the QR code plus a copyable plain-text secret as a fallback.<br>3. Prompts for a 6-digit confirmation code before activating.<br>4. On successful confirmation, triggers the AUTH-08 backup-code display flow. | QR code renders correctly; enrollment cannot complete without a valid confirmation code; backup codes are shown once immediately after successful enrollment. |
| AUTH-11-BE | Active Session & Device Self-Management API | P0 | Backend endpoints letting a user view and revoke their own active sessions/devices. | 1. `GET /user/security/sessions` returns rows from `UserActiveSession` joined with device metadata (device label, OS, `last_seen`, IP/region, current-session flag).<br>2. `DELETE /user/security/sessions/{id}` revokes the specified session.<br>3. `POST /user/security/sessions/revoke-all-others` revokes every session except the caller's current one.<br>4. Revocation immediately invalidates the target session's refresh token. | User can retrieve a list of all active devices; revoking a session forces that device to logout per AUTH-02; "revoke all others" leaves only the current session active. |
| AUTH-11-MOB | Device & Session Management UI | P0 | Mobile Settings -> Security "Devices" screen. | 1. Lists devices with last-active timestamp and approximate location.<br>2. Each row has a "Log out" action.<br>3. A "Log out of all other devices" button prompts a confirmation dialog before calling revoke-all-others. | Device list renders correctly; individual and bulk revocation succeed; the current device is clearly labelled and excluded from self-revocation. |
| AUTH-12-BE | Change Password & Global Session Invalidation | P0 | Backend endpoint for authenticated password change, plus shared session-invalidation behavior on password reset. | 1. `POST /user/security/change-password` requires the current password (and current MFA verification if `mfa_enabled`).<br>2. Enforces the AUTH-07 complexity rules and rejects reuse of the current password.<br>3. On success — whether via this change-password flow or a forgot-password reset (AUTH-07) — revokes all sessions except the one that initiated the change and rotates refresh tokens.<br>4. Sends a security notification (SMS/email/push) confirming the change. | Wrong current password is rejected; on success, all other active sessions are invalidated; the user receives a security notification. |
| AUTH-12-MOB | Change Password UI | P0 | Mobile Settings -> Security change-password form. | 1. Form collects current password, new password, and confirmation, reusing the `PasswordStrengthIndicator` for a live strength meter.<br>2. Requires an MFA step if `mfa_enabled`.<br>3. On success, displays "You've been logged out of other devices". | Form enforces live complexity feedback; MFA is required when enabled; success message correctly reflects the AUTH-12-BE session invalidation. |
| AUTH-13-BE | Support-Assisted Account Recovery | P1 | Backend admin workflow for recovering accounts where the user has lost both phone access and MFA (including exhausted backup codes). | 1. User re-proves identity via an identity-provider match (reusing REG-N01/REG-R07 through the provider framework) plus a document/liveness check submitted through the existing admin manual-verification queue (REG-R08's queue).<br>2. An admin, authenticated with their own MFA and required to supply a written justification (same pattern as AUTH-06-WEB revert-phone), resets the user's phone and/or MFA configuration.<br>3. All of the user's active sessions are force-logged-out.<br>4. The action is written to `AuditLog`. | Recovery cannot proceed without a passing identity re-proof; admin reset requires admin MFA and justification; all sessions are terminated; the action is fully audit-logged. |
| AUTH-14-BE | New Sign-In Notification | P0 | Backend alerting on login from an unrecognized device or location. | 1. On successful login, checks whether `device_id` (per FR-03 fingerprinting) is unrecognized for the user, or whether the login region/geovelocity signal is unrecognized.<br>2. If either is true, sends a push and SMS/email notification: "New sign-in from {device}, {location}, {time} — Not you? Secure your account", deep-linking into AUTH-11 device/session management. | An unrecognized device or region triggers the alert within 1 minute of login; logins from recognized devices/regions do not trigger it. |
| AUTH-15-BE | Change Email API | P0 | Backend endpoint mirroring the AUTH-06 secure-phone-change pattern for email address changes. | 1. `POST /user/change-email` requires current password, current MFA token (if enabled), and the new email address.<br>2. Sends a verification link to the new address, reusing the REG-R09 verification mechanism.<br>3. On confirm, updates the user's email, sends a notification to the OLD address, and logs the change to `UserEmailHistory`.<br>4. The new email must reach `email_verified = true` (REG-R09) before the change is considered complete. | The old address receives a change notification; the new address must complete verification before the change takes effect; the change is logged to `UserEmailHistory`. |

---

#### 5.1.5 Combined Role Transitions

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| REG-U01-BE | Combined Role Registration API | P1 | Allow existing users (renters) to apply for Owner or Operator roles under a single account. | 1. Receives role application payload (business license, certifications, bank details).<br>2. Verifies existing user session.<br>3. Creates linked `OwnerProfile` or `OperatorProfile` record in DB with status `PENDING_REVIEW`. | API saves new profiles linked to the same parent `User` record; rejects duplicate profile creation. |
| REG-U01-MOB | Combined Role Settings UI | P1 | Mobile profile settings interface allowing renters to apply to list machinery or operate. | 1. Renter navigates to Profile Settings.<br>2. Taps "Become an Owner" or "Register as Operator".<br>3. Displays document upload forms.<br>4. Submits to API. | Form captures required files; calls role API; alerts user that role is pending review. |
| REG-U01-WEB | Combined Role Review Portal | P1 | Next.js admin interface displays additional role applications in standard review queues. | 1. Admin views Owner/Operator verification queues.<br>2. System flags applications from existing users.<br>3. Admin reviews documents and approves/rejects. | Verification table displays combined profile links; approval adds role permission to User. |


---


### 5.2 Fleet Registration – Machinery & Logistics

#### 5.2.1 Basic Information

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-01-BE | Category Lookup API | P0 | Expose category catalog endpoints to populate creation forms. | 1. Retrieve parent-child categories from `Category` database.<br>2. Respond to app client.<br>3. Retrieve spec field models based on selected category. | Category tree returned cleanly; specs associated properly in DB queries. |
| FL-01-MOB | Category Selector Form | P0 | Form input selectors fetching and loading specs dynamically. | 1. Fetch categories on form load.<br>2. User selects parent and child dropdown values.<br>3. Fetch specific specs for choice.<br>4. Redraw spec inputs dynamically. | Dropdowns load sequentially; selecting subcategory triggers dynamic form render. |
| FL-02-BE | Make, Model, Year Validation | P0 | Enforce presence and value constraints on primary fields. | 1. Check fields exist in JSON.<br>2. Verify year is between 1980 and current year + 1.<br>3. Reject invalid registrations. | Registration rejected if fields empty or year is out of bounds. |
| FL-02-MOB | Primary Detail Fields UI | P0 | Text and picker fields on fleet form. | 1. Form displays make and model inputs.<br>2. Year selection is a scrolling picker (1980 to current year + 1). | Inputs require entries; validate inputs before allowing page progression. |
| FL-03-BE | Optional Plate Storage | P0 | Database storage allowing optional alphanumeric plate numbers. | 1. Validate plate field regex (alphanumeric only).<br>2. Save plate string or NULL to the database. | Accepts and stores valid alphanumeric plate numbers; allows blank values. |
| FL-03-MOB | Plate Input field | P0 | Optional text input in registration form. | 1. Owner can enter alphanumeric plate number.<br>2. Form permits submission if left empty. | Accepts alphanumeric values; validation fails on special characters. |
| FL-04-BE | Hour Meter Recording API | P0 | Store initial telemetry and image attachments. | 1. Extract numeric hour reading and image upload URL from payload.<br>2. Save parameters to the fleet's baseline records. | Hour meter value and image URL successfully stored. |
| FL-04-MOB | Telemetry & Photo Capture | P0 | Camera capturing and numeric entry for hour meters. | 1. User inputs meter numeric reading.<br>2. Tap photo button to launch camera.<br>3. Take photo of meter.<br>4. Compress and upload image. | Both inputs are required; user cannot advance without capturing hour meter photo. |
| FL-05-BE | Image Batch Upload Handler | P0 | Accept and save mandatory fleet photos, with count, size, format, and safety bounds. | 1. Verify registration payload has at least 5 photo URLs, with a maximum of 20 photos per fleet.<br>2. Match photo tags to required angles.<br>3. Enforce max file size 10MB per file.<br>4. Reject files whose format is not JPEG, PNG, or WebP.<br>5. Route every uploaded file through the platform's existing malware/virus scan pipeline before it is persisted or made accessible; infected files are quarantined and the upload rejected. | Validates photo counts (min 5, max 20) and per-file size (<=10MB); rejects disallowed formats; returns upload URL paths only for files that pass the malware scan; infected uploads are quarantined and rejected. |
| FL-05-MOB | Multi-Angle Photo Grid | P0 | Image slots mapping required angles on mobile screen. | 1. Show grid with 5 required upload containers (Front, Rear, Left, Right, Interior) plus up to 15 additional optional slots (20 max total).<br>2. Open camera or gallery picker on tap, restricted to JPEG/PNG/WebP.<br>3. Show thumbnails.<br>4. Compress/upload images (client-side warning if a file exceeds 10MB before upload). | Require 5 distinct photos minimum and allow up to 20 total; thumbnail preview updates on upload success; oversized or disallowed-format files are rejected client-side with a clear message. |
| FL-06-BE | Fuel Type Enumeration | P0 | Validate fuel selection against ENUM schemas. | 1. Validate input against ENUM: `Diesel`, `Electric`, `Petrol`, `Hybrid`.<br>2. Save values in database. | API rejects invalid fuel values; saves valid choices successfully. |
| FL-06-MOB | Fuel Type Dropdown | P0 | Selection list inside fleet form. | 1. Form shows a selection widget populated with allowed fuel types.<br>2. Owner chooses options. | Dropdown requires choice; selection correctly maps to API request keys. |
| FL-06B-BE | Notes API Support | P1 | Accept and output descriptive text fields. | 1. Enforce length <= 1000 characters on description.<br>2. Save to fleet table.<br>3. Include in details endpoint. | Saves long text notes; rejects string inputs over 1000 characters. |
| FL-06B-MOB | Fleet Notes Input UI | P1 | Character-restricted text box on forms and details screens. | 1. Owner registration shows text box with character counter (max 1000).<br>2. Details screen displays notes. | Text box restricts input past 1000 characters; displays description correctly. |

---

#### 5.2.2 Dynamic Specifications

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-07-BE | Spec Definition Retrieval API | P0 | Provide custom fields definition list based on category ID. | 1. Fetch `SpecFieldDefinition` rows by category ID.<br>2. Return JSON list including field types and requirements. | Returns correct JSON containing types, keys, and validation properties. |
| FL-07-MOB | Dynamic spec layout builder | P0 | Render custom form widgets from schema JSON. | 1. Query spec schema for selected subcategory.<br>2. Loop over keys to render inputs. | Dynamic fields render correctly in form based on API definitions. |
| FL-08-MOB | Collapsible form accordions | P0 | Group input forms under accordion headers. | 1. Group form widgets based on SpecFieldDefinition `Section` values.<br>2. Render each section in an expandable block. | Accordions collapse/expand; displays number of incomplete required fields. |
| FL-09-BE | Dynamic payload validation | P0 | Enforce required fields on variable spec payloads. | 1. Retrieve required schema properties.<br>2. Validate dynamic JSON data values.<br>3. Reject with missing field list on failure. | Backend blocks creation if any mandatory dynamic field is empty or missing. |
| FL-09-MOB | Dynamic form validation UI | P0 | Client validation for dynamic inputs. | 1. Loop through mandatory spec fields on submit.<br>2. Trigger validation messages.<br>3. Scroll to first missing field. | Shows red validation warnings; blocks API dispatch if incomplete. |
| FL-10-MOB | Spec Field Unit Labels | P0 | Append measurement suffixes in forms and details. | 1. Check unit attribute in field definition.<br>2. Append units (e.g. "m", "kg") in input overlays and specs listings. | Units display correctly in text field placeholders and read-only pages. |
| FL-11-BE | Custom Spec Range Checks | P1 | Validate numeric inputs against range fields. | 1. Compare input with `MinValue` and `MaxValue`.<br>2. If out of range, reject with error details. | Rejects registrations containing values outside defined boundaries. |
| FL-11-MOB | Custom Spec Range Validation UI | P1 | Apply local form validation triggers displaying a message when values exceed the Min/Max range. | 1. Monitor dynamic numeric inputs.<br>2. Validate values locally against min/max ranges.<br>3. Display error under the text field if range is exceeded. | Warning displays instantly on validation failure; prevents submission. |
| FL-12-BE | Enum spec option validation | P1 | Backend validation verifying submitted enum specs. | 1. Ensure value matches predefined enum items in DB schema. | API blocks invalid choices, only saving database-approved enum elements. |
| FL-12-MOB | Enum field selector widget | P1 | Populate selector lists dynamically. | 1. Check if field type is ENUM.<br>2. Load options array.<br>3. Render dropdown option selection. | Dropdown displays correct choices; selected key updates payload structure. |
| FL-13-BE | Renter spec filtering middleware | P0 | Hide restricted fields from renter details. | 1. Query fleet specs.<br>2. Filter out keys where `IsVisibleToRenter` is false.<br>3. Return cleaned list to renter endpoints. | Renters cannot view restricted fields; owners and admins can retrieve full list. |
| FL-13-MOB | Filtered Spec Details layout | P0 | Screen showing allowed details only. | 1. Retrieve specs list for renter page.<br>2. Display specifications inline. | Shows basic specifications; sensitive fields are absent. |
| FL-14-BE | Logistics extraction engine | P0 | Index essential logistics parameters into a separate DB model. | 1. Detect spec records saved with `IsLogisticsCritical=true`.<br>2. Flatten and write fields to `FleetSpecIndexed` columns. | Automatically populates indexing tables on fleet updates and creations. |
| FL-15-BE | Comparative spec API | P1 | Output overlapping specifications for comparative analysis. | 1. Accept list of 2-3 fleet IDs.<br>2. Query specs.<br>3. Filter keys to only include those with `IsComparable=true`. | API returns comparative table with matching fields for active fleets. |
| FL-15-MOB | Comparison comparison screen | P1 | View side-by-side comparison tables. | 1. Select up to 3 fleets.<br>2. Call comparison API.<br>3. Render side-by-side column grid. | Limits choices to 3; layout allows easy scanning of attributes. |
| FL-15B-BE | Booking spec snapshot engine | P1 | Record a historical state of specifications upon booking. | 1. On booking confirmation, read active fleet specs.<br>2. Write copy into `BookingSpecSnapshot`.<br>3. Retrieve this copy for historical reviews. | Future modifications of fleet details do not affect active or historical bookings. |

---

#### 5.2.3 Pricing, Booking Rules & Payment Details

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-16-BE | Base pricing validation | P0 | Enforce required rates and currency fields on creation APIs, including sanity bounds to catch typos/fraud. | 1. Validate daily rate is present and > 0.<br>2. Check other pricing rates (hourly/weekly) are valid numeric values in ETB.<br>3. Validate every rate field (daily/hourly/weekly) against a configurable min/max sanity bound read from `SystemConfig` (e.g. `pricing.daily_rate_min` default 100 ETB, `pricing.daily_rate_max` default 500,000 ETB, with analogous keys for hourly/weekly), analogous to the Year sanity bound already enforced in FL-02-BE.<br>4. Reject with a clear validation error (field name, allowed range) if any rate falls outside its bound. | Registration fails if daily rate is missing or invalid; rates outside the configured min/max sanity bounds are rejected with a clear error identifying the field and allowed range. |
| FL-16-MOB | Rate configuration form | P0 | Numerical entry forms for setting fleet rates. | 1. Display pricing inputs.<br>2. Require daily rate.<br>3. Render warning/error if invalid values. | Owner enters pricing; fields restrict inputs to positive numbers. |
| FL-17-BE | Minimum period checking | P0 | Enforce minimum duration requirements on booking requests. | 1. On booking request, calculate duration.<br>2. Compare against fleet's minimum period limit.<br>3. Block creation if invalid. | Booking fails if duration is shorter than configured minimum. |
| FL-17-MOB | Minimum period selector | P0 | Settings selectors for rental periods. | 1. Display dropdown: 1 hour, 4 hours, 8 hours, 1 day, etc.<br>2. Store selected value in fleet payload. | Choices list renders; selection correctly configures booking rules. |
| FL-18-BE | Overtime calculations | P0 | Calculate costs for exceeding booked rental hours. | 1. Compare actual rental duration with booked window.<br>2. Calculate excess hours.<br>3. Multiply by excess rate to calculate final invoice. | Invoices apply overtime charges correctly when return is late. |
| FL-18-MOB | Overtime rate inputs | P0 | Number settings on owner form. | 1. Render excess hour ETB input.<br>2. Save settings. | Rates are saved and displayed correctly in owner dashboard. |
| FL-19-BE | Deposit data fields | P0 | Save security deposit requirements on fleet profiles. | 1. Save security deposit parameter.<br>2. Display value in details queries. | Field saved correctly and sent to app endpoints. |
| FL-19-MOB | Deposit info displays | P0 | Informational displays of deposits. | 1. Render input on owner settings.<br>2. Display deposit fee to renter during booking checkout. | Renter sees deposit requirement; owner sets it successfully. |
| FL-19B-BE | Security-deposit collection & release lifecycle (future) | P2 | Reserve a full deposit lifecycle for a post-MVP phase; in MVP the deposit is informational-only. | **MVP:** deposit is **informational-only** — never collected, never part of Subtotal/VAT/Total/Prepayment (see BK-04, CW-00-BE).<br>**Future:** define a deposit lifecycle (`collected → held → partially/fully released → refunded/forfeited`) modelled as an **owner obligation** consistent with the Direct Transfer Model (owner-tracked, no platform-held renter funds), reusing the BK-CX-04 obligation/evidence pattern. | MVP treats deposit as informational-only; future collection/release lifecycle documented as owner-obligation (no platform custody). |
| FL-20-BE | Booking cancel scheduler | P0 | Cron checking engine canceling bookings on overdue payment uploads. | 1. Query pending bookings daily/hourly.<br>2. If duration exceeds fleet's timeout (default 24h) and no receipt, update booking to `CANCELLED`. | Unpaid bookings are automatically cancelled on schedule. |
| FL-20-MOB | Timeout parameter config | P0 | Input screen to adjust payment windows. | 1. Input box for hours.<br>2. Renter app displays timer showing payment time remaining. | Owner sets custom hours; renter app updates and runs countdown active timer. |
| FL-21-BE | Prepayment calculator | P0 | Apply rules to determine required deposit sums. | 1. Check prepayment rules (Full, Half, Percentage, Fixed).<br>2. Compute payment request total on checkout. | Computes accurate prepayment required for invoices. |
| FL-21-MOB | Prepayment setup inputs | P0 | UI to configure prepayment types. | 1. Dropdown for prepayment options (Full, Half, %, Fixed).<br>2. Input values.<br>3. Render required amount in renter checkout. | Configuration options are saved; renter sees exact required amount. |
| FL-22-BE | Renter cancellation handler | P0 | Compute the refundable amount (prepayment minus renter forfeiture) and hand it to the owner-obligation refund flow. | 1. Check cancellation time against timeline.<br>2. Deduct penalty percentage from prepayment to compute refundable ETB (penalty basis per §5.5.4 / NE-PAY-05).<br>3. **Do NOT auto-refund** — the platform never disburses renter funds; instead create/hand off a `RefundObligation` for owner-tracked settlement per **BK-CX-04-BE** (owner uploads evidence, renter confirms, overdue → Admin/wallet-debit). | Refundable amount computed correctly per configuration; refund handled as an owner obligation via BK-CX-04 (no platform disbursement). |
| FL-22-MOB | Renter cancellation settings | P0 | Form input to set penalty rates. | 1. Display numeric slider/input for cancellation penalty (%).<br>2. Save. | Settings saved and displayed on owner policies. |
| FL-23-BE | Owner cancellation ledger | P0 | Log cancellation details and compute penalty charges. | 1. On owner cancellation, evaluate penalty rules.<br>2. Apply fees to owner ledger. | Owner account gets correctly debited for cancellations. |
| FL-23-MOB | Owner penalty config UI | P0 | Inputs to select penalty rules. | 1. Option buttons for owner penalty type.<br>2. Text input for value. | Configures owner penalties; fields validated before saving. |
| FL-24-BE | Expose bank settings API | P0 | Deliver owner payment info to checkout API. | 1. Read owner bank data.<br>2. Return coordinates in payment endpoint responses. | Restricts bank details exposure to authorized bookings. |
| FL-24-MOB | Owner bank coordinates displays | P0 | Render bank info on checkout screens. | 1. Show bank, holder name, account details.<br>2. Copy details button. | Renter sees correct coordinates and can copy them. |
| FL-24B-MOB | Pricing details layout | P0 | Detailed pricing grid for fleets. | 1. Render pricing review widget containing daily rates, deposit, and cancellations notes. | Summary display is complete and displays prior to book clicks. |
| FL-24C-BE | Holiday rates engine | P2 | Select rates based on calendars. | 1. Check booking dates against holidays/weekends.<br>2. Adjust daily pricing calculation. | Applies custom rates on weekends and public holiday dates. |
| FL-24C-MOB | Weekend rate options UI | P2 | Toggle and text input for custom rates. | 1. Toggle switch for weekend rates.<br>2. Rate input box in ETB. | Toggle reveals inputs; updates rates correctly on save. |

---

#### 5.2.4 Operator Assignment

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-25-BE | Operator inclusion flag persistence | P0 | Database storage for tracking operator availability settings. | 1. Save boolean setting `IsOperatorIncluded` in fleet details table. | API stores and returns operator availability settings. |
| FL-25-MOB | Operator availability toggle | P0 | Form toggle widget in owner app fleet registration. | 1. Toggle switch on listing forms.<br>2. If checked, reveal selection dropdowns. | Toggle works; visual options collapse/expand correctly. |
| FL-26-BE | Linked operator validation | P0 | Backend checking to ensure only valid team operators are assigned. | 1. Validate assigned operators exist in `OwnerOperatorLink`.<br>2. Link operator IDs to fleet. | API blocks assigning operators who are not linked to owner's team. |
| FL-26-MOB | Operator select selector | P0 | Selector widgets to choose team operators. | 1. Retrieve list of active operators.<br>2. Show checkbox/dropdown selector list.<br>3. Save selections. | Renders active operator lists; user checks multiple options. |
| FL-27-BE | Booking operator update dispatch | P1 | Process operator re-assignments and dispatch notifications. | 1. Update booking records with new operator ID.<br>2. Dispatch push notification to renter. | Operator is updated in database; renter receives alert. |
| FL-27-MOB | Re-assign operator screens | P1 | Interface for updating assigned operator on active bookings. | 1. Owner views booking.<br>2. Tap "Change Operator".<br>3. Select new operator.<br>4. Save. | Action is accessible on active bookings; alerts renter immediately. |
| FL-27B-BE | Secure operator profiles endpoint | P0 | Serve operator details to authorized renter bookings. | 1. Fetch details (name, photo, phone) of assigned operator.<br>2. Return in booking endpoint. | Restricts operator details to active and confirmed renters only. |
| FL-27B-MOB | Operator contact info card | P0 | UI panel on renter booking details containing operator info. | 1. Display operator avatar card.<br>2. Show name and phone number.<br>3. Tap "Call" button to initiate phone dialer. | Card is visible after booking confirmation; dialer triggers successfully. |

---

#### 5.2.5 Logistics Configuration

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-28-BE | Self-move flag database fields | P0 | Save self-move capacity boolean on fleet records. | 1. Write `CanSelfMove` flag to database on fleet registry request. | Flag is stored; default query returns field. |
| FL-28-MOB | Self-move toggle switches | P0 | Toggles for self-move configuration. | 1. Display Yes/No toggle.<br>2. If No, display transport class options. | Toggle functions; toggles input requirements based on choice. |
| FL-29-BE | Transport suggestion calculations | P1 | Algorithm calculating vehicle types from specs. | 1. Read dimensions/weight from `FleetSpecIndexed`.<br>2. Evaluate classification tables.<br>3. Return suggestion ID. | API returns matching transportation type classes. |
| FL-29-MOB | Suggestion overrides UI | P1 | Display suggestions with custom override controls. | 1. Render recommended transport class tag.<br>2. Allow choosing custom category on dropdown override. | Recommended class loads; owner can override choice. |
| FL-30-BE | Force Owner Pay logic | P0 | Enforce owner payment liability rules for MVP logistics. | 1. Force logistics pricing responsibility to `OWNER` in invoice calculations. | System maps logistics billing records to owner accounts. |
| FL-31-BE | Transport sub-booking link | P1 | Core workflows linking transport bookings to primary rentals. | 1. Process transport booking request.<br>2. Create sub-booking.<br>3. Link back using `parent_booking_id`. | Sub-booking registers with correct link reference in DB. |
| FL-31-MOB | Arranging logistics selectors | P1 | Screen to select platform or external transport. | 1. Owner views booking.<br>2. Click "Arrange Transport".<br>3. Shows platform fleets or "Handle Externally". | Submits selection; transitions to payment or sets external flag. |

---

#### 5.2.6 Attachments Bundling

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-32-BE | Bundled attachments links API | P1 | Store relationship updates in `FleetAttachment` databases. | 1. Accept attachment ID list for fleet updates.<br>2. Write rows to `FleetAttachment`. | Mappings write successfully; removes obsolete references. |
| FL-32-MOB | Attachment pick checklists | P1 | Checkbox inventory list in fleet forms. | 1. Fetch available attachments list.<br>2. Render checkbox list.<br>3. Toggle items. | Lists active attachments; captures selected items. |
| FL-33-BE | Bundled rate calculator | P1 | Process custom rates for bundled attachments. | 1. Check custom rate overrides on checkout.<br>2. Apply custom rate if exists; otherwise standalone rate. | Invoice calculations represent correct bundled attachment prices. |
| FL-33-MOB | Dynamic rates fields | P1 | Input fields for setting bundled pricing. | 1. Render ETB rate text inputs beside checked attachments.<br>2. Validate inputs. | Stores custom rates; defaults to standalone rate if empty. |
| FL-34-BE | Overlap verification queries | P1 | Check attachment booking conflicts before confirmation. | 1. Run overlap checks for attachment bookings during booking requests.<br>2. Return conflict errors. | Conflicting bookings blocked; returns specific error code. |
| FL-34-MOB | Attachment unavailable alerts | P1 | Highlight overlapping attachment conflicts during checkouts. | 1. Display message if attachment is unavailable.<br>2. Permit renter to check out without conflicts. | Allows checkout completion by removing unavailable items. |

---

#### 5.2.7 Availability & Listing Status

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| FL-35-BE | Block availability endpoints | P0 | Bulk upsert availability data blocks. | 1. Expose availability updates API.<br>2. Save ranges in `FleetAvailability`. | DB stores date ranges; blocks bookings on these dates. |
| FL-35-MOB | Calendar grids dashboards | P0 | Interactive screens to toggle blocked dates. | 1. Open availability calendar tab.<br>2. Tap dates to toggle blocked (red).<br>3. Save updates. | Grid functions; updates are sent to server on save. |
| FL-36-BE | Status machines validations | P0 | Enforce allowed states transitions on fleets using the canonical fleet-status enum. | 1. Restrict states to the canonical fleet-status enum: `DRAFT`, `PENDING_APPROVAL`, `ACTIVE`, `REJECTED`, `INCOMPLETE_SPEC`, `PAUSED`, `RENTED`, `INACTIVE`. The awaiting-approval state is `PENDING_APPROVAL` (not "Pending"/"PENDING_REVIEW"). Typical flow: `DRAFT -> PENDING_APPROVAL -> ACTIVE -> RENTED -> INACTIVE` (with `REJECTED`, `INCOMPLETE_SPEC`, `PAUSED` as branch states).<br>2. Apply rules. | Rejects forbidden transitions (e.g. `DRAFT -> ACTIVE` without approval). |
| FL-36-MOB | Listing status headers | P0 | Render status indicators on listings management lists. | 1. Show status tags (Draft, Active, Pending) on fleet cards.<br>2. Show corresponding action buttons. | Labels match current API state; actions match permissions. |
| FL-37-BE | Admin approval updates API | P0 | Expose endpoints for updating review states. | 1. Allow admins to submit approval status and reason.<br>2. Update status.<br>3. Send notification. | Fleet status transitions to ACTIVE or REJECTED; notifications sent. |
| FL-37-WEB | Admin Fleet approval tables | P0 | Web dashboards showing pending listings. | 1. Next.js dashboard shows queue of fleets waiting review.<br>2. View details and photos.<br>3. Click Approve/Reject. | Admin can view all specs and files; action calls API correctly. |
| FL-37-MOB | Approval notifications alerts | P0 | Display results to owners. | 1. Display status update alert to owner.<br>2. Show details if rejected. | Prompts review alerts; displays comments for rejections. |
| FL-37B-BE | Modified listing checks | P1 | Categorize changes and trigger approvals when editing. | 1. Identify modified fields.<br>2. If critical (photos, specs), reset status to `PENDING_APPROVAL`. | Non-critical changes save immediately; critical resets state. |
| FL-37B-MOB | Editing warning screens | P1 | Display warnings during fleet updates. | 1. Detect edits on critical fields.<br>2. Show modal warning that re-approval is required. | Warning displays before submission; user confirms action. |
| FL-38-BE | Owner-initiated fleet pause/deactivate | P1 | Allow an owner to temporarily pause or permanently deactivate their own fleet listing. | 1. Expose an endpoint for the owner to toggle a fleet to `PAUSED` (canonical fleet-status enum per FL-36-BE: `DRAFT`, `PENDING_APPROVAL`, `ACTIVE`, `REJECTED`, `INCOMPLETE_SPEC`, `PAUSED`, `RENTED`, `INACTIVE`); a `PAUSED` fleet is immediately hidden from search/discovery results, but any existing confirmed booking on it continues unaffected through to completion.<br>2. Allow the owner to reactivate a `PAUSED` fleet back to `ACTIVE`, making it visible in search again.<br>3. Expose a separate endpoint for the owner to request permanent deactivation (`INACTIVE`); before applying it, check for any booking on the fleet in a non-terminal state (i.e. not `COMPLETED` or `CANCELLED`) and block the request with a clear error if one exists, instructing the owner to wait until the booking reaches a terminal state. | Owner can pause a fleet (removed from search, confirmed bookings unaffected) and reactivate it later; permanent deactivation is blocked with a clear error while any non-terminal booking exists, and succeeds once all bookings are `COMPLETED`/`CANCELLED`. |
| FL-38-MOB | Pause/deactivate listing controls | P1 | Owner-facing controls to pause, reactivate, or deactivate a fleet from the listings management screen. | 1. Show "Pause Listing" / "Reactivate" toggle and a "Deactivate Fleet" action on each fleet card in the owner's listings screen.<br>2. On deactivate, call the API and surface the non-terminal-booking block error inline if returned. | Toggle correctly flips fleet between `ACTIVE` and `PAUSED`; deactivation attempt shows a clear blocking message when a non-terminal booking exists. |

---

### 5.3 Attachment Registration & Inventory

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| AT-01-BE | Add Attachment API | P1 | Process and save new attachment details. | 1. Accept attachment inputs and files.<br>2. Save to attachment inventory.<br>3. Set ACTIVE. | Record created successfully; files stored. |
| AT-01-MOB | Add Attachment screen | P1 | Forms to create new attachment records. | 1. Open "Add Attachment".<br>2. Input details and photos.<br>3. Submit. | Validates fields; sends uploads and submits form. |
| AT-02-BE | Map category indexes | P1 | Store compatibility relations between attachments and categories. | 1. Accept compatibility category IDs.<br>2. Save relationships in DB. | Mappings queryable; filters work. |
| AT-02-MOB | Select category checklists | P1 | Selection UI for compatibility. | 1. Render checkbox lists of categories.<br>2. Store selection payload. | Allows checkbox selections; saves correctly. |
| AT-03-BE | Standalone pricing setup | P1 | Save pricing plans for separate rentals. | 1. Save hourly/daily/weekly rates.<br>2. Compute prices for standalone booking. | Calculations apply pricing rates correctly. |
| AT-03-MOB | Standalone pricing inputs | P1 | Price fields on attachment forms. | 1. Input rates.<br>2. Submit values. | Requires daily price; handles valid numeric entry. |
| AT-04-BE | Attachment schedule tables | P1 | Save unavailable date blocks for attachments. | 1. Expose availability endpoints.<br>2. Save in database. | Availability dates stored; blocks bookings. |
| AT-04-MOB | Block attachment calendars | P1 | Calendars to select blocked dates. | 1. Render calendar grid on attachment.<br>2. Toggle blocked dates.<br>3. Save. | Selection functions; saves inputs. |
| AT-05-BE | Bundle associations API | P1 | Link attachments to fleets in database. | 1. Accept linkage request.<br>2. Map attachment to fleet. | Records written; queries retrieve listings. |
| AT-05-MOB | Bundle select screens | P1 | Roster list in fleet edit for bundling. | 1. Show available attachments.<br>2. Select.<br>3. Link. | Displays matching items; links successfully. |
| AT-06-BE | Attachment overlap checker | P1 | Query database for attachment booking overlaps. | 1. Query bookings table for date collisions.<br>2. Reject booking if overlap exists. | Prevents double booking; returns clear conflict error. |

---

### 5.4 Fleet Discovery, Search & Comparison

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| SR-01-BE | Load Categories API | P0 | Expose active category endpoints. | 1. Query active categories.<br>2. Return tree JSON. | Returns structured parent/child JSON array. |
| SR-01-MOB | Category home page grids | P0 | Navigation category icons on home screen. | 1. Home page renders grid of categories.<br>2. Tapping navigates to subcategory screen. | Category list loads; buttons navigate correctly. |
| SR-02-BE | Search index query API | P0 | Perform text matches against fleets and owners. | 1. Query search keywords against indexed fields.<br>2. Return matches. | Quick search returns matching name, model, make. |
| SR-02-MOB | Debounced search bars | P0 | Text search input with debounce handling. | 1. Input keywords.<br>2. Debounce query 300ms.<br>3. Call API.<br>4. Display result card grid. | Prevents multiple calls; renders correct summaries on cards. |
| SR-03-BE | Dynamic spec filter endpoints | P0 | Returns valid attributes and limits for current category search. | 1. Identify category.<br>2. Retrieve range values and distinct enums for specs.<br>3. Filter fleets by filters. | Returns query metadata; validates filter parameters. |
| SR-03-MOB | Dynamic spec filters panels | P0 | Render matching selectors based on category spec fields. | 1. Query available spec filters.<br>2. Render range sliders or checkbox lists.<br>3. Show active filters as chips. | Panel adapts dynamically; chips close to remove filters. |
| SR-04-BE | Database sorting queries | P0 | Sort query outputs by price, score, or spatial distance. | 1. Parse sorting query param.<br>2. Order records.<br>3. Return sorted results. | Sort matches choices; distance ordering functions. |
| SR-04-MOB | Sort picker dialogs | P0 | Segmented pickers to change search sorting. | 1. Tap Sort.<br>2. Select sorting metric (Price, Rating, Distance).<br>3. Call search API. | Dropdown displays choices; calls API to reorder list. |
| SR-05-BE | Fetch detail summary API | P0 | Return complete dataset for fleet pages. | 1. Query fleet, specs, owner, reviews, and calendar.<br>2. Return combined JSON payload. | Aggregates all necessary details into a single response. |
| SR-05-MOB | Fleet details screen layout | P0 | Renders detail page for selected fleet. | 1. Renders photo gallery carousel, specification table, owner info card, reviews, and booking date picker. | All views load; user can scroll specs and see ratings. |
| SR-06-BE | Date availability API | P0 | Verify fleet availability for target date ranges. | 1. Accept fleet ID and range.<br>2. Check bookings and availability database.<br>3. Return availability state. | Accurate verification; blocks overlap checkouts. |
| SR-06-MOB | Date check selectors | P0 | Calendar overlay to check date status. | 1. Open date picker.<br>2. Query availability endpoint.<br>3. Show "Available" or "Unavailable" reason. | Blocks booked dates; prompts status clearly. |
| SR-07-BE | Get comparison specs | P1 | Return compare metrics for specified fleet lists. | 1. Retrieve specs marked `IsComparable=true` for selected fleet IDs. | Returns JSON matching comparison fields. |
| SR-07-MOB | Fleet compare screens | P1 | Side by side comparison dashboard. | 1. Renders selected fleets side-by-side.<br>2. Lists matching specs in rows. | Renders up to 3 fleets; scrolls horizontally. |
| SR-08-MOB | Local views tracking | P1 | Local log of recently accessed pages. | 1. Save viewed fleet ID to local cache (SQLite/Hive) on page open.<br>2. Display top 10 on home screen. | Persists across app sessions; shows up to 10 items. |
| SR-09-BE | Wishlist persistence endpoints | P1 | Save wishlists items to databases. | 1. Expose add/delete endpoints.<br>2. Save in database.<br>3. Expose list endpoints. | Saves correctly in DB; allows listings retrieval. |
| SR-09-MOB | Wishlist heart controls | P1 | Toggle bookmark states on cards. | 1. Tap heart icon.<br>2. Call API.<br>3. Display in "Saved" page. | Icon changes color; lists populate saved lists tab. |
| SR-10-BE | Geolocation spatial queries | P0 | Query database coordinates using PostGIS filters. | 1. Query records within distance range of lat/lng.<br>2. Order by distance. | Calculates distances accurately; returns sorted lists. |
| SR-10-MOB | GPS locations prompt | P0 | Check device permission and query coordinates. | 1. Request GPS permission.<br>2. Fetch lat/lng.<br>3. Submit values to search endpoint. | Prompt works; handles location-disabled states. |
| SR-11-BE | Paged queries endpoints | P0 | Implement cursor/offset page query handling. | 1. Slice DB records.<br>2. Return data with pagination cursors. | Limits database payload; prevents high database loads. |
| SR-11-MOB | Infinite scrolls listeners | P0 | Scroll controller loading next page of cards. | 1. Monitor scroll offset.<br>2. Trigger API fetch on bottom proximity.<br>3. Append cards. | Loads pages smoothly; prevents duplicate fetches. |
| SR-12-BE | Sponsored listings rank logic | P0 | Weight search algorithms for Sponsored (owner-purchased, `IsBoosted`) listings. | 1. Set rank weights by Sponsored tiers.<br>2. Expose purchase endpoints.<br>3. Schedule auto-expiration tasks. | Gold/Silver/Bronze Sponsored tiers float to top; checks expiration daily. (Sponsored = owner-purchased paid ranking; distinct from admin-curated Featured and from Ad Campaign banners — see §2 glossary.) |
| SR-12-WEB | Setup package tiers pricing | P0 | Next.js admin page managing Sponsored tiers. | 1. Admin views active Sponsored placements.<br>2. Manage prices and tier duration packages. | Admin configures package prices; details update live. |
| SR-12-MOB | Highlight Sponsored badge UI | P0 | Render Sponsored tags on Sponsored search results. | 1. Renders special border styling on Sponsored listings.<br>2. Show "Sponsored" badge label. | Badges display correctly; does not block organic list view. |
| SR-13-BE | Status tag computations | P1 | Compute short string availability tags for listings. | 1. Check bookings schedule.<br>2. Generate summary label ("Available this week"). | Calculations are accurate and return simple strings. |
| SR-13-MOB | Card availability tags | P1 | Display availability strings on listing cards. | 1. Display computed label inline on result cards. | Tags render clean styling on search listings. |
| SR-14-BE | Amharic text translations match | P0 | Support multilingual text searches. | 1. Map keywords to English and Amharic dictionaries.<br>2. Match queries in both strings. | Amharic search strings return matched items. |
| SR-14-MOB | Ge'ez inputs compatibility | P0 | Parse keyboard characters correctly. | 1. Submit text directly from mobile keyboard. | Handles input format without formatting errors. |
| SR-15-BE | Organic search relevance ranking | P1 | Define the default "Recommended" sort as a weighted-score organic ranking, distinct from the paid Sponsored ranking in SR-12-BE. | 1. When sort mode is "Recommended" (default), compute a relevance score per listing as a weighted sum of: (a) recency — listing/spec-update freshness; (b) average rating; (c) distance from the searching user's location (when available, per SR-10-BE); and (d) profile completeness (photo count relative to the 5-20 range in FL-05-BE, and dynamic-spec completeness).<br>2. Read the per-factor weights from `SystemConfig` (e.g. `search.rank_weight_recency`, `search.rank_weight_rating`, `search.rank_weight_distance`, `search.rank_weight_completeness`), admin-configurable without a deploy.<br>3. Sort descending by computed score; break ties by listing `created_at` descending (newest first).<br>4. Sponsored listings (SR-12-BE) are ranked/inserted above organic results per their tier logic and do not participate in this organic score. | "Recommended" sort returns organic listings ordered by the documented weighted score; changing a weight in `SystemConfig` changes ranking without a deploy; equal-score listings are ordered by creation date descending; Sponsored placement logic is unaffected. |


---


### 5.5 Booking & Reservation Flow

#### 5.5.1 Creating a Booking Request

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `BK-01-BE` | Date selection validation | P0 | Backend must validate booking start/end dates. Start date must be in the future and before the end date. | 1. Receives incoming booking creation payload.<br>2. Validates that `start_date` is strictly in the future relative to server UTC time.<br>3. Asserts that `end_date` is chronologically after `start_date`.<br>4. Rejects invalid dates with RFC 9457 400 Bad Request. | Date validation works; rejects past dates or end date before start date with descriptive validation problem details. |
| `BK-01-MOB` | Date selection calendar UI | P0 | Mobile app calendar with time pickers for renters to select start and end dates/times. | 1. User taps "Book Now" on fleet details.<br>2. App presents full-screen calendar modal with start and end time selectors.<br>3. Validates date bounds client-side prior to network submission.<br>4. Disables past dates and dates marked as blocked in local cache. | Calendar displays correctly; past dates are disabled/not selectable; valid date range persists to checkout state. |
| `BK-02-BE` | Check fleet availability API | P0 | Backend must verify fleet availability by checking for overlapping active bookings (excluding CANCELLED, COMPLETED, RATED states) and blocked dates. | 1. Queries database booking records for target `fleet_id`.<br>2. Checks for overlapping intervals where status is not in (`CANCELLED`, `COMPLETED`, `RATED`).<br>3. Soft-locks slots currently in `PENDING_PAYMENT` and `PAYMENT_EVIDENCE_UPLOADED`.<br>4. Rejects conflicting requests with HTTP 409 Conflict. | Overlap check rejects conflicting bookings (including those awaiting payment proof) and returns appropriate error code. |
| `BK-02-MOB` | Display availability conflict message | P0 | Mobile app shows error message to the renter if selected dates overlap with an existing booking. | 1. Intercepts HTTP 409 response from availability or booking initiation API.<br>2. Parses RFC 9457 error payload.<br>3. Displays a persistent modal or alert banner: "Fleet is unavailable for the selected dates."<br>4. Keeps date picker open to allow reselection. | Visual conflict message is clearly displayed on mobile screen when validation fails; prevents accidental resubmissions. |
| `BK-03-BE` | Attachment availability & price API | P1 | Backend returns list of available fleet attachments and validates selected attachments and calculates pricing upon reservation. | 1. Queries database for attachments linked to fleet via `FleetAttachment`.<br>2. Filters out attachments with conflicting booking intervals.<br>3. Calculates updated total price applying bundled rate discounts if defined.<br>4. Returns attachment availability list with calculated minor-unit rates. | Correct attachment availability returned; price updates correctly on checkout; conflicting attachments excluded. |
| `BK-03-MOB` | Attachment selection checkboxes UI | P1 | Mobile UI shows list of attachments with checkboxes. Unavailable items are disabled with an explanation. Price summary updates dynamically. | 1. Renders attachments list inside booking setup sheet.<br>2. Disables checkboxes for unavailable attachments with "Booked for selected dates" badge.<br>3. Toggling checkboxes updates local pricing state and triggers cost breakdown recalculation. | Checkboxes function correctly; disabled state displays reason; pricing summary updates dynamically on toggle. |
| `BK-04-BE` | Calculate cost breakdown API | P0 | API endpoint to compute and return detailed cost breakdown: rental subtotal, VAT tax (15%), total, and a separate informational deposit line. | 1. Reads base rates and attached tool rates.<br>2. Calculates **Subtotal = rental fee + attachments (deposit EXCLUDED)**.<br>3. Computes 15% VAT on taxable items (unless renter is marked `TaxExempt`).<br>4. Returns Subtotal, VAT, Total, and separate `deposit_informational` line (deposit is informational-only in MVP, excluded from total). | Subtotal = rental fee + attachments only; deposit excluded from Subtotal/VAT/Total/Prepayment and shown as a separate informational line; correct VAT and Total computed in minor units. |
| `BK-04-MOB` | Cost breakdown UI card | P0 | Mobile app renders the cost breakdown summary card (Subtotal, VAT, Total) before booking submission. | 1. Fetches breakdown from API.<br>2. Populates summary card showing Subtotal, 15% VAT line-item, and Total.<br>3. Displays security deposit as separate informational badge: "Security Deposit (Payable on site to owner): [Amount] ETB". | Visual breakdown matching the calculated API response (showing VAT and total) displays correctly; clear deposit disclaimer. |
| `BK-05-BE` | Create booking record endpoint | P0 | Backend creates a Booking record in database with status `PENDING_PAYMENT` using unique idempotency key validation. | 1. Exposes `POST /bookings` endpoint.<br>2. Validates `Idempotency-Key` header against Valkey cache.<br>3. Inserts booking row in `PENDING_PAYMENT` state with initialized `version = 1`.<br>4. Returns booking record, ID, and initial status wrapped in `SuccessEnvelope`. | Booking record created with status `PENDING_PAYMENT` and returned; duplicate requests return original response. |
| `BK-05-MOB` | Submit booking request action | P0 | Mobile app triggers booking creation when renter taps "Request to Book", sending idempotency key. | 1. User reviews summary card.<br>2. Taps "Request to Book".<br>3. Generates UUIDv4 idempotency key and dispatches POST request.<br>4. On HTTP 200/201 success, routes renter to payment instructions and timer screen. | Tapping "Request to Book" successfully triggers API call and redirects; handles network retries safely with idempotency key. |
| `BK-06-BE` | Booking payment details API | P0 | Backend API returns owner's bank account information, required prepayment amount, and expiry timestamp. | 1. Validates that requesting user is the booking renter.<br>2. Fetches fleet owner profile and bank details (bank name, account number, holder name).<br>3. Calculates exact timeout timestamp ($T_{\text{created}} + \text{timeout\_hours}$, default 24h).<br>4. Returns coordinates and prepayment minor-unit amount. | Owner bank details and prepayment amount correctly returned; access blocked for unauthorized callers. |
| `BK-06-MOB` | Payment instructions & timer UI | P0 | Mobile screen displays owner's bank details, prepayment amount, and a circular countdown timer. | 1. Screen loads booking payment details.<br>2. Renders bank coordinates with one-tap clipboard copy buttons.<br>3. Parses expiration timestamp and initializes circular countdown timer.<br>4. Emits warning when timer drops below 1 hour. | Circular timer counts down accurately; correct banking details shown; copy action updates clipboard. |
| `BK-06B-BE` | Pre-booking validation API | P0 | API to validate final booking parameters and retrieve complete confirmation details including cancellation policies. | 1. Exposes `GET /bookings/pre-confirm`.<br>2. Re-validates live asset availability, owner status, and price rates.<br>3. Aggregates complete summary payload: pricing, fleet specs, cancellation policy terms, and operator inclusion details. | Pricing, fleet details, cancellation terms, and operator info returned cleanly; flags any changes since search. |
| `BK-06B-MOB` | Pre-booking summary UI screen | P0 | Mobile app displays non-editable confirmation screen with fleet details, dates, costs, policies, and operator info. | 1. Displays non-editable review screen before final hold commitment.<br>2. Shows fleet photo, dates, itemized costs, cancellation penalty clauses, and operator status.<br>3. "Confirm & Proceed to Payment" initiates booking; "Edit" returns to date/spec selection. | Summary screen shown; explicit confirmation required via button click; navigates accurately. |
| `BK-06C-BE` | Minimum rental duration backend check | P0 | Backend must reject booking requests if duration is less than the fleet's configured minimum rental period. | 1. Compares computed rental duration against fleet's `min_rental_duration` property.<br>2. If duration $<$ minimum, aborts creation with RFC 9457 error code `BOOKING_DURATION_BELOW_MINIMUM`. | Minimum period enforced; database rejects invalid durations with structured 400 error. |
| `BK-06C-MOB` | Minimum rental duration UI validation | P0 | Mobile app validates booking duration locally and disables the booking action if the minimum rental period is not met. | 1. Calculates selected duration in real time.<br>2. Compares with fleet minimum requirement.<br>3. If below minimum, displays warning banner: "Minimum rental period for this machine is X days" and disables "Book Now". | Alert/warning message shown; checkout blocked locally for durations below minimum. |
| `BK-06D-BE` | Backend idempotency check | P0 | Backend must check for existing bookings matching the incoming client-provided UUID idempotency key. | 1. Inspects `Idempotency-Key` HTTP header.<br>2. Queries Valkey/database for existing key record.<br>3. If found, returns HTTP 200 with existing booking details, bypassing duplicate record creation. | Duplicate bookings never created; identical idempotency key returns original booking payload. |
| `BK-06D-MOB` | Idempotency key generation | P0 | Mobile app generates a UUID idempotency key for each unique booking submission attempt. | 1. Generates fresh UUIDv4 when renter opens the pre-booking confirmation screen.<br>2. Attaches token to `Idempotency-Key` header on booking submission.<br>3. Retains identical key across automated network retries. | UUID generated once per booking attempt; sent in request header; preserved during retries. |

---

#### 5.5.2 Payment Evidence Upload & Verification

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `BK-07-BE` | Payment evidence receipt upload API | P0 | Backend processes uploaded receipt image, stores it, and transitions booking status to `PAYMENT_EVIDENCE_UPLOADED`. | 1. Accepts multipart image upload for booking in `PENDING_PAYMENT`.<br>2. Streams file through ClamAV antivirus inspection (`zINSTREAM`).<br>3. Persists file to SeaweedFS S3 storage and records URL in `Payment` table.<br>4. Atomically transitions booking status to `PAYMENT_EVIDENCE_UPLOADED`. | Image stored successfully; antivirus scan passes; status changes to `PAYMENT_EVIDENCE_UPLOADED`; infected files rejected with 422. |
| `BK-07-MOB` | Image picker & crop UI for receipt | P0 | Mobile app enables renter to select or capture a photo of the transfer receipt, crop/preview it, and upload it. | 1. Renter taps "Upload Payment Proof".<br>2. Camera or gallery picker launches.<br>3. Provides square/rectangular crop tool to focus on transaction reference.<br>4. Renders preview confirmation sheet.<br>5. Dispatches multipart upload to backend. | Image picker opens; crop tool displays preview; upload successfully triggers API and transitions view. |
| `BK-07B-BE` | Payment-evidence duplicate & anomaly detection | P1 | Backend detects reused or anomalous payment-evidence images on upload and surfaces flags into owner review and reconciliation. | 1. On receipt upload, computes cryptographic SHA-256 hash and perceptual hash (pHash).<br>2. Scans database for matching hashes across other bookings.<br>3. If match found, sets `evidence_duplicate_flag = TRUE`.<br>4. Optionally runs OCR to extract amount/date and flags divergence via `evidence_amount_mismatch_flag`.<br>5. Attaches flags to evidence for owner review. | Content + perceptual hashes computed and stored; cross-booking reuse sets duplicate flag; flags visible in owner review and admin reconciliation queues. |
| `BK-08-BE` | Trigger payment evidence notifications | P0 | Backend generates push notification and in-app notification for the fleet owner once renter uploads payment evidence. | 1. Triggered by transition to `PAYMENT_EVIDENCE_UPLOADED`.<br>2. Enqueues notification job in event queue.<br>3. Dispatches FCM push notification payload to owner's registered devices.<br>4. Inserts notification row in `Notification` table. | Notifications generated and queued successfully upon receipt upload; push payload formatted with deep link. |
| `BK-08-MOB` | Payment notification receipt UI response | P0 | Mobile app alerts the owner of payment evidence and navigates them to the booking detail review page upon tapping the notification. | 1. Listens for FCM push payload.<br>2. Displays high-priority heads-up notification banner.<br>3. Tapping notification opens app and deep-links directly to booking payment review view. | Notification received on device; tap opens correct booking detail page directly. |
| `BK-09-BE` | Confirm payment API endpoint | P0 | Backend API endpoint for owners to approve payment, setting booking status to `CONFIRMED` and sending notifications. | 1. `PUT /bookings/{id}/confirm-payment` validates caller is owner.<br>2. Validates booking is in `PAYMENT_EVIDENCE_UPLOADED`.<br>3. Starts single `pgx.Tx` transaction.<br>4. Updates status to `CONFIRMED` and increments `version`.<br>5. Debits platform commission from `OwnerWallet` (`COMMISSION_DEBIT`).<br>6. Commits transaction and triggers push alerts. | Status updated to `CONFIRMED`; in-transaction ledger debit executed atomically; both parties notified. |
| `BK-09-MOB` | Owner payment confirmation UI | P0 | Owner mobile screen displays payment slip and has "Confirm Payment" button with a confirmation popup. | 1. Owner opens booking detail screen.<br>2. Renders high-resolution zoomable payment slip image with duplicate warning badges if flagged.<br>3. Taps "Confirm Payment".<br>4. Confirmation dialog prompts: "Confirm receipt of [Amount] ETB?".<br>5. Submits confirmation API call. | Payment slip displays clearly; warning flags visible; confirmation dialog triggers successfully. |
| `BK-10-BE` | Decline payment API endpoint | P0 | API to record decline reason, transition status to `PAYMENT_DECLINED`, and notify the renter. | 1. `PUT /bookings/{id}/decline-payment` receives decline reason code and comment.<br>2. Validates caller is owner.<br>3. Updates booking status to `PAYMENT_DECLINED`.<br>4. Inserts row in `BookingStatusHistory`.<br>5. Dispatches push notification to renter. | Database status updated to `PAYMENT_DECLINED`; decline details recorded; renter notified immediately. |
| `BK-10-MOB` | Owner decline options UI | P0 | Owner mobile screen offers a "Decline" button which prompts a reason selector modal with text comments. | 1. Owner taps "Decline Payment".<br>2. Modal presents structured reason dropdown ("Payment not received", "Incorrect amount", "Illegible receipt", "Other").<br>3. Requires comment text if "Other" is selected.<br>4. Dispatches decline API call. | Reason dropdown and text box function correctly; API payload includes correct inputs. |
| `BK-11-BE` | Handle actions on declined payment APIs | P0 | API logic to support payment retry, free cancellation (setting status to `CANCELLED` without penalties), or opening a dispute. | 1. Endpoints evaluate current `PAYMENT_DECLINED` state.<br>2. Retry upload allows new slip submission (returns to `PAYMENT_EVIDENCE_UPLOADED`).<br>3. Cancel transitions to `CANCELLED` with zero penalty recorded.<br>4. Dispute transitions to `DISPUTE_OPENED`. | Endpoints process retry, cancellation, or dispute transitions correctly; cancellation sets penalty to zero. |
| `BK-11-MOB` | Decline action options UI | P0 | Mobile app displays three actions to the renter on decline: "Retry Upload", "Cancel Booking", or "Dispute Decline". | 1. Renter opens declined booking view.<br>2. UI displays owner decline reason and comment.<br>3. Renders 3 actionable buttons: "Retry Upload", "Cancel Booking (Free)", and "Dispute Decline".<br>4. Each triggers its corresponding navigation/API flow. | All three actions available; cancellation is free and performs immediate checkout cleanup. |
| `BK-12-BE` | Open & resolve dispute APIs | P0 | Backend APIs to support opening a dispute, adding evidence/comments, and allowing admins to resolve the dispute in favor of either party. | 1. `POST /bookings/{id}/dispute` creates `Dispute` record with origin `payment`, setting status to `DISPUTE_OPENED`.<br>2. Multipart upload saves dispute attachments.<br>3. Admin endpoint `POST /admin/disputes/{id}/rule` records binding ruling, transitioning booking to `CONFIRMED` or `CANCELLED`. | Dispute state updates in DB; admin resolution transitions booking state correctly and logs audit record. |
| `BK-12-MOB` | Dispute creation & upload UI | P0 | Mobile app allows renter to file a dispute with description/evidence, and allows owner/renter to view dispute details. | 1. Renter selects "Dispute Decline".<br>2. Form collects dispute description and supporting evidence photos (bank SMS screenshot, statement).<br>3. Submits dispute payload to API.<br>4. Booking banner updates to show "Under Admin Dispute Review". | Dispute form renders; renter can type notes and upload image proof; status updates cleanly. |
| `BK-12-WEB` | Admin dispute management console | P0 | Next.js admin portal to review open disputes, view evidence uploads from both sides, and issue final rulings. | 1. Admin navigates to Operations $\rightarrow$ Disputes.<br>2. Opens pending dispute record.<br>3. Reviews renter slip, owner decline notes, and chat evidence.<br>4. Action buttons: "Rule for Renter (Confirm)" or "Rule for Owner (Cancel)".<br>5. Requires written ruling justification and admin MFA. | Admin console displays all submitted evidence and enables ruling action; updates propagate live. |
| `BK-12B-BE` | Dispute SLA cron & escalation logic | P1 | Backend calculates dispute deadline (3 business days), runs daily cron job to check for overdue disputes, and flags/escalates them. | 1. On dispute creation, computes `dispute_deadline` as 3 business days out (skipping weekends/holidays per `LOC-02-BE`).<br>2. Daily cron job checks for unresolved disputes past deadline.<br>3. Flags dispute as overdue and dispatches high-priority alert to Super Admins.<br>4. Applies across all dispute origins (`payment`, `return_damage`, `hour_meter`). | Deadline calculated correctly; cron job flags overdue entries across all dispute origins. |
| `BK-12B-WEB` | Overdue dispute indicator in Admin Panel | P1 | Next.js admin dashboard highlights overdue disputes with visual alerts and displays escalation warnings to Super Admins. | 1. Admin dispute queue sorts overdue items to top.<br>2. Displays prominent red "SLA Overdue" badge.<br>3. Super Admin header shows counter of overdue disputes. | SLA overdue badge displays correctly; dashboard filtering works; alerts administrators. |
| `BK-12C-BE` | Dispute messaging API | P1 | Backend database schema and API to store dispute messages with target visibility settings. | 1. `POST /disputes/{id}/messages` saves message to `DisputeMessage` table.<br>2. Validates `visibility_scope` enum (`renter-only`, `owner-only`, `both`).<br>3. Admin can send private notes to either party; parties can reply to admin.<br>4. Restricts read endpoints based on caller role. | Messages saved to database; visibility rules validated on load; parties only see allowed messages. |
| `BK-12C-MOB` | Mobile dispute chat screen | P1 | Mobile app displays chat interface within dispute screen for communicating with the admin. | 1. User enters dispute details view.<br>2. Opens embedded "Dispute Chat" tab.<br>3. Renders chat timeline with admin messages.<br>4. Input bar allows sending text replies and attaching images. | Chat interface displays message history; input field sends messages via API in real time. |
| `BK-12C-WEB` | Admin dual-panel chat interface | P1 | Next.js admin console provides side-by-side or separate chat panels to communicate with the renter and owner independently. | 1. Dispute review page renders dual chat panels: "Chat with Renter" and "Chat with Owner".<br>2. Admin types messages in respective panel.<br>3. Backend automatically assigns corresponding `visibility_scope`. | Separate panels show correct message threads; admin replies route only to selected party. |

---

#### 5.5.3 Post-Confirmation & Readiness

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `BK-13-BE` | Update fleet readiness status API | P1 | API to transition booking to `AWAITING_FLEET_READINESS`, save owner's reason, and store ETA date. | 1. `PUT /bookings/{id}/readiness` accepts reason code, notes, and `readiness_eta` timestamp.<br>2. Validates caller is owner and booking is in `CONFIRMED`.<br>3. Transitions booking status to `AWAITING_FLEET_READINESS`.<br>4. Enqueues push notification to renter with ETA details. | Booking status transitions; reason and ETA persisted; renter notified immediately. |
| `BK-13-MOB` | Owner readiness setting & renter ETA display UI | P1 | Mobile app enables owner to set readiness reason/ETA, and displays this ETA warning to the renter. | 1. Owner active booking screen displays "Not Ready Yet" action.<br>2. Modal collects reason and ETA picker.<br>3. Submits readiness payload.<br>4. Renter booking view updates with orange status banner showing owner notes and ETA countdown. | UI supports owner input; renter sees updated status and ETA details clearly. |
| `BK-14-BE` | Additional payment request & resolution APIs | P1 | API to create an additional payment request, update booking status, and process the renter's confirmation or refusal. | 1. `POST /bookings/{id}/additional-payment` creates fee record and sets status to `PENDING_ADDITIONAL_PAYMENT`.<br>2. Renter pay action accepts payment evidence, advancing to `READY_FOR_PICKUP` on confirmation.<br>3. Renter decline triggers cancellation with penalty assessed against owner.<br>4. Idempotent per `BK-17-BE`; evidence deduped per `BK-07B-BE`. | Backend processes additional fee creation, payment updates, and penalty flags on decline; confirmation idempotent. |
| `BK-14-MOB` | Owner payment request & renter response UI | P1 | Mobile screens for owner to request additional funds, and for renter to accept/pay or decline the request. | 1. Owner fills additional fee request form (amount, category, description).<br>2. Renter app displays modal alert with breakdown.<br>3. Renter options: "Upload Payment Proof" or "Decline & Cancel (Owner Penalty)". | UI inputs and options present; transitions and screens function accurately. |
| `BK-14B-BE` | ETA check cron & notification logic | P1 | Cron job checks for AWAITING_FLEET_READINESS bookings with expired ETAs, flags them, and notifies renter of free cancellation. | 1. Hourly cron checks bookings in `AWAITING_FLEET_READINESS` where `readiness_eta < NOW()`.<br>2. Sets `eta_exceeded = TRUE` flag on booking.<br>3. Dispatches high-priority push notification to renter.<br>4. Activates free cancellation rights for renter. | Cron job flags bookings; notifications sent; free cancellation capability activated. |
| `BK-14B-MOB` | Renter free cancellation prompt UI | P1 | Renter's booking detail screen updates to offer a free cancellation action if the fleet ETA has expired. | 1. Detects `eta_exceeded = TRUE` flag.<br>2. Renders prominent banner: "Owner missed the readiness ETA. You can cancel for free."<br>3. Action buttons: "Cancel Booking (Free Refund)" or "Continue Waiting". | Free cancellation button rendered and functional only when backend flag allows. |
| `BK-14B-WEB` | Admin dashboard readiness flags | P1 | Next.js admin dashboard flags bookings that have exceeded readiness ETA for admin review/monitoring. | 1. Queries bookings in operations queue.<br>2. Renders red "Overdue Readiness" badge on affected rows.<br>3. Enables support agent to contact owner. | Overdue listings visible to administrators; filterable in dashboard. |
| `BK-15-BE` | Admin Stuck State Machine Super-Override | P0 | Endpoint allowing Super Admins to force-transition bookings, bypassing normal transition rules. | 1. Receives `POST /admin/bookings/{id}/force-state` with target status and mandatory justification ($\ge 100$ chars).<br>2. Asserts caller has `bookings:override` permission and active MFA session.<br>3. Logs before/after state diff in `AuditLog`.<br>4. Updates booking status, force-clears locks in Valkey, and dispatches push alerts. | Force-transition successfully bypasses state machine logic; logs detailed JSON diff; notifies renter and owner. |
| `BK-15-WEB` | Admin Override Portal UI | P0 | Next.js admin page showing manual override options for active bookings. | 1. Super Admin opens booking details.<br>2. Clicks "Force State Transition".<br>3. Modal requires selecting target status, entering 100-character justification, and verifying TOTP MFA code.<br>4. Dispatches override API request. | Dialog enforces verification inputs; override action executes and displays success alert. |
| `BK-16-BE` | Generate PDF Invoice API | P0 | Endpoint to compile and export formal itemized PDF receipts and invoices for bookings, carrying all legally mandatory fiscal fields. | 1. Exposes `GET /api/v1/bookings/{id}/invoice/pdf`.<br>2. Compiles itemized charges: base rental, attachments, excess hours, damages, 15% VAT, payments received, balance due.<br>3. Inserts legal fiscal fields per `BK-16B-BE`.<br>4. Generates and streams PDF document. | PDF carries invoice number, issue timestamp, issuer name+TIN, buyer name+TIN, itemised VAT breakdown; reprints marked "COPY". |
| `BK-16B-BE` | Legal invoice numbering & fiscal fields | P0 | Backend guarantees every issued invoice carries a gap-free, monotonically increasing per-issuer sequential number and all legally mandatory Ethiopian VAT-invoice fields, retained for 7 years. | 1. Atomic increment on `InvoiceSequence` per legal issuer (owner TIN / platform TIN).<br>2. Persists immutable invoice metadata (sequential number, timestamp, issuer/buyer legal names and TINs, 15% VAT breakdown, immutable UUID).<br>3. First render is original; subsequent renders watermarked "COPY" without consuming numbers.<br>4. Voids/credit notes issued as new sequenced documents referencing original.<br>5. Retained 7 years. | Sequential invoice numbers increment atomically with no gaps/duplicates; reprints marked "COPY"; retained 7 years for tax audit compliance. |
| `BK-16-MOB` | Mobile Download Invoice Button | P0 | Renders invoice download actions on mobile booking completed screen. | 1. Completed booking detail screen displays "Download PDF Invoice" button.<br>2. Downloads PDF to device storage and launches native viewer. | UI displays action button; PDF launches on click cleanly. |
| `BK-16-WEB` | Web Admin & User Invoice Export | P0 | Next.js screens to view and export PDF invoices. | 1. Renters and owners click "Download Invoice" on billing tables.<br>2. Downloads PDF document. | Exports PDF invoice with proper layout and fiscal headers. |
| `BK-17-BE` | Idempotent money-mutating actions | P0 | All financial state transitions require idempotency keys and produce at-most-once ledger effects, applied atomically with the status change. | 1. Every financial action (`confirm-payment`, `decline-payment`, `additional-payment`, `final-payment`, `wallet-topup`) requires `Idempotency-Key` header.<br>2. Enforces unique constraint on `(booking_id, txn_type)` in `financial_ledger_entries`.<br>3. Status change and ledger write execute inside single `pgx.Tx` block guarded by OCC `version`.<br>4. OCC conflict returns HTTP 409 `CONCURRENT_MODIFICATION`. | Replayed financial actions are no-ops returning original response; unique constraint blocks double-posting; status change + ledger write are atomic under OCC guard. |
| `BK-18-BE` | Final/additional payment timeout & escalation | P1 | Configurable timeouts for the two renter-payment wait states, with cron-driven notification, Admin flag, and a configurable terminal action. | 1. Configurable deadlines: `final_payment_deadline` (default 3 business days) for `PENDING_FINAL_PAYMENT` and `additional_payment_deadline` for `PENDING_ADDITIONAL_PAYMENT`.<br>2. Hourly cron identifies overdue bookings, dispatches reminders, and flags admin ops queue.<br>3. Executes configurable terminal action (auto-cancel with penalty or hold for admin adjudication).<br>4. Audit-logs all breaches. | Both states have configurable deadlines; hourly cron notifies + Admin-flags past deadline; configurable terminal action executed; all actions audit-logged. |

---

#### 5.5.4 Cancellation Logic

##### Cancellation Window & Policy Matrix

| Cancellation Window | Who Can Cancel | Penalty | UI Behaviour |
|---|---|---|---|
| **Before CONFIRMED** | Renter or Owner | None | Confirmation popup, no charge. |
| **After CONFIRMED** | Renter | Owner keeps X% of prepayment | Warning: "You will forfeit [exact ETB amount]." |
| **After CONFIRMED** | Owner | Owner pays fixed/% to renter | Warning: "You will be charged [exact ETB amount]." |
| **During PENDING_ADDITIONAL_PAYMENT (renter declines)** | Renter | **Owner penalty** (owner caused the non-readiness) | Warning to renter: "You can cancel for free. The owner will be charged a penalty." |
| **ETA exceeded in AWAITING_FLEET_READINESS** | Renter | None (free cancellation) | Banner: "Owner missed the readiness ETA. You can cancel for free." |
| **Handover timeout in READY_FOR_PICKUP (48h) — neither party acted** | System | None | Auto-cancel 48h after start date; status shown as "Cancelled (Timeout)". |
| **Handover timeout in READY_FOR_PICKUP (48h) — one party completed inspection/attestation (`AR-01`), counterparty did not** | System | No-show penalty on the absent party (renter forfeits prepayment %, owner pays penalty) | Auto-cancel 48h after start date; absent party treated as no-show and charged the corresponding cancellation penalty. |
| **During DISPUTE_OPENED** | Neither | Cannot cancel | Cancel button hidden. |
| **During ACTIVE** | Neither | Early return = normal completion | No cancel option; "End Rental" starts return flow. |

##### Cancellation Penalty Basis (`NE-PAY-05`)
- **Basis:** The renter cancellation penalty is a percentage **of the prepayment** (not of the booking total); the owner cancellation penalty is a **fixed ETB amount OR a percentage of the prepayment**, per the owner's configured penalty rule (`FL-22`/`FL-23`, `PenaltyRuleSet`/`PenaltyRule`). Deposit is excluded from every penalty base (deposit is informational-only in MVP).
- **Configurable per fleet with platform bounds:** Each fleet configures its penalty rate/amount; the platform sets a **default** (e.g. renter forfeits **20%** of prepayment; owner pays a configured fixed/% amount) and enforces **min/max bounds** (e.g. renter forfeit clamped to `[0%, 50%]` of prepayment) so per-fleet values cannot exceed platform limits.
- **Exact ETB derivation:** The "exact ETB amount" shown in warning dialogs (`BK-CX-01`) is derived as `round(penalty_rate × prepayment)` for percentage rules or the configured fixed amount, using `roundMoney()` and integer minor-units per `CW-16-BE`. The forfeit/penalty and the resulting refundable amount feed `BK-CX-04-BE` (refund obligation) and `BK-CX-05-BE` (owner-penalty compensation).

##### Cancellation Requirements Inventory

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `BK-CX-01` | Exact cancellation penalty display | P0 | Mobile and web apps must display the exact calculated minor-unit currency penalty before confirming cancellation. | 1. User taps "Cancel Booking".<br>2. Client calls pre-cancel calculation endpoint.<br>3. Computes exact forfeit amount in ETB/KES minor units based on timeline rules.<br>4. Renders modal: "You will forfeit [Exact Amount] ETB. Proceed?". | Dialog renders exact calculated currency amount rather than generic percentages; blocks cancellation until explicit confirmation. |
| `BK-CX-02` | Immediate availability release on cancel | P0 | Cancelled bookings must release fleet availability immediately so other renters can book. | 1. Booking status transitions to `CANCELLED`.<br>2. Backend releases blocked date ranges in `FleetAvailability`.<br>3. Deletes associated mutual exclusion lock in Valkey.<br>4. Fleet becomes immediately discoverable in search for those dates. | Fleet availability dates freed immediately upon cancellation; other users can reserve the slot without delay. |
| `BK-CX-03` | Cancellation reason capture | P0 | A cancellation reason must be collected from the cancelling party for platform analytics. | 1. User initiates cancellation.<br>2. Form presents structured dropdown of cancellation reasons.<br>3. Allows optional free-text explanation note.<br>4. Persists reason code and text in `BookingStatusHistory`. | Cancellation blocked until reason is selected; reason code and notes persisted in history log. |
| `BK-CX-04-BE` | MVP refund obligation & tracking | P0 | Refunds are an owner obligation; the platform tracks settlement without holding or disbursing renter funds. | 1. On refundable cancellation, computes `refundable_amount = prepayment - forfeit`.<br>2. Creates `RefundObligation` record (`booking_id`, `amount`, `status = PENDING`, `deadline`).<br>3. Owner uploads bank-transfer refund slip (`EVIDENCE_UPLOADED`).<br>4. Renter confirms receipt (`CONFIRMED`).<br>5. Overdue triggers admin escalation and optional wallet debit fallback. | Refund tracked as owner obligation; owner uploads transfer proof; renter confirms receipt; overdue escalates to admin. |
| `BK-CX-05-BE` | Owner-penalty compensation to renter | P0 | When an owner cancellation or readiness failure incurs an owner penalty payable to the renter, create compensation obligation. | 1. Assesses owner penalty on owner cancellation or readiness default.<br>2. Creates `CompensationObligation` record with deadline.<br>3. Records liability on owner ledger.<br>4. Owner uploads transfer evidence to renter; renter confirms.<br>5. Overdue escalates to Admin with optional `OwnerWallet` debit fallback (`PENALTY_DEBIT`). | Compensation obligation created; owner uploads proof; renter confirms; overdue escalates to Admin with wallet debit fallback. |

---

### 5.6 Active Rental Management

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `AR-01-BE` | Pre-rental inspection API | P0 | Backend stores pre-rental inspection photos and hour meter reading from both renter and owner. | 1. `POST /bookings/{id}/inspection/pre` accepts numerical hour meter value, mandatory meter photo URL, and machinery condition photos.<br>2. Verifies file formats and uploads.<br>3. Records inspection in `BookingInspection`. | Database stores readings and photo references; status tracks inspection completion for both parties. |
| `AR-01-MOB` | Pre-rental inspection upload UI | P0 | Mobile app prompts both renter and owner to take/upload machinery photos and input/photograph the hour meter. | 1. App opens pre-rental inspection wizard.<br>2. Enforces camera-only capture for hour meter reading.<br>3. User enters numeric reading matching photo.<br>4. Uploads 4-angle condition photos.<br>5. Submits inspection payload. | Screen blocks proceeding without hour meter photo and numeric input; uploads data correctly to backend. |
| `AR-02-BE` | Confirm handover API | P0 | Backend transitions booking status to `ACTIVE` and records timestamp once both renter and owner confirm handover. | 1. `POST /bookings/{id}/handover` logs confirmation for calling party.<br>2. Asserts pre-rental inspections (`AR-01`) and attestations (`AR-08`) are complete.<br>3. When both parties confirm, updates status to `ACTIVE` and records `started_at` timestamp. | Booking status turns `ACTIVE` only when both users have confirmed handover and completed inspections. |
| `AR-02-MOB` | Confirm handover button UI | P0 | Mobile app displays "Confirm Handover" button to both owner and renter after inspection is complete. | 1. Active booking panel reveals "Confirm Handover" button once inspection is submitted.<br>2. Tapping calls handover API.<br>3. Displays waiting state: "Waiting for counterparty confirmation".<br>4. Transitions view to active operational dashboard upon completion. | Buttons function; UI reflects confirmation status of both users; transitions to active dashboard. |
| `AR-03-BE` | Fetch active rental data API | P0 | Backend API returns active booking details, calculated time remaining, operator details, and allowable actions. | 1. `GET /bookings/{id}/active` retrieves current rental duration, calculated time remaining, operator information, and permission flags.<br>2. Calculates elapsed vs remaining hours.<br>3. Serves assigned operator contact coordinates. | Returns accurate time remaining and operator data payload; restricts data to authorized participants. |
| `AR-03-MOB` | Active rental dashboard screens | P0 | Mobile app shows active booking dashboard with rental countdown, operator info, issue reports, extension options, and operator change controls. | 1. Displays remaining rental time using live circular countdown widget.<br>2. Displays operator avatar, phone dialer, and live GPS map.<br>3. Renders action buttons: "Report Issue", "Request Extension", and "End Rental".<br>4. Owner view displays "Change Operator" action. | Layout differs correctly by role; timers count down accurately; action buttons trigger respective workflows. |
| `AR-04-BE` | Create & process extension request APIs | P1 | Backend endpoints to request a booking extension, calculate cost increments, and handle owner's approval or rejection. | 1. `POST /bookings/{id}/extension` creates extension proposal with new end date.<br>2. Re-runs availability check against existing bookings and blocked dates; conflicting requests rejected.<br>3. Calculates incremental rental fee, VAT, and incremental commission.<br>4. Owner approval updates booking end date and debits incremental commission (`CW-14-BE`). | Extension record created; rates verified; availability re-checked and conflicts rejected; owner approval updates booking timeline. |
| `AR-04-MOB` | Extension request & review UI | P1 | Renter mobile UI to request extensions with new dates, and owner mobile UI to accept/decline extension requests. | 1. Renter opens "Extend Rental" sheet and selects new end date.<br>2. Displays cost breakdown for extension period.<br>3. Submits request.<br>4. Owner receives push alert, reviews proposed dates/fees, and taps "Accept" or "Decline". | Date picker restricts invalid selections; costs update dynamically; buttons trigger appropriate APIs. |
| `AR-05-BE` | Issue reporting API | P1 | API to receive issue report, record in DB, upload images, and notify relevant users. | 1. `POST /bookings/{id}/issues` stores issue record (`severity`, `category`, `description`, photos).<br>2. Saves attachments in SeaweedFS S3.<br>3. Enqueues high-priority alert to owner and support team. | Issue logged in database; notifications sent; visible in administrative monitoring console. |
| `AR-05-MOB` | Report issue form UI | P1 | Mobile app form for reporting active rental issues, including descriptions, severity toggles, and photo attachment pickers. | 1. User taps "Report Issue".<br>2. Selects severity (`Low`, `Medium`, `High`, `Critical / Breakdown`).<br>3. Enters description and attaches photos from camera.<br>4. Submits issue report to backend. | Form fields validate; photos attach; success triggers API and returns to dashboard. |
| `AR-05-WEB` | Admin panel issue monitoring | P1 | Next.js admin dashboard lists reported issues from active bookings to monitor incident resolutions. | 1. Admin navigates to Operations $\rightarrow$ Active Issues.<br>2. Table displays active reports filtered by severity.<br>3. Click opens modal showing photos, machine telematics, and communication thread. | Admin can view reports, filter by severity, and mark issues as resolved. |
| `AR-06-BE` | Scheduled rental reminders cron | P0 | Backend cron job checks for upcoming and ending rentals, sending push reminders 24 hours in advance. | 1. Scheduled cron runs hourly.<br>2. Queries bookings where `start_date` or `end_date` is exactly 24 hours in the future.<br>3. Publishes reminder jobs to notification queue.<br>4. Dispatches push notification and SMS reminders. | Job fires on schedule; reminders successfully queued and dispatched to target users. |
| `AR-06-MOB` | Receive & handle reminder alerts | P0 | Mobile app receives push reminders for upcoming or ending rentals, routing users to details on tap. | 1. Intercepts incoming reminder push payload.<br>2. Renders local notification alert.<br>3. Tapping notification opens app and routes directly to the active booking detail screen. | Reminders display correctly; tapping routes user to correct booking screen. |
| `AR-07-BE` | Handover timeout cron | P1 | Backend worker auto-cancels unconfirmed handovers 48 hours after the start date, applying a no-show penalty when only one party acted. | 1. Cron queries bookings in `READY_FOR_PICKUP` where `start_date + 48h < NOW()`.<br>2. Evaluates inspection records (`AR-01`).<br>3. If neither party inspected, cancels with zero penalty.<br>4. If one party inspected and counterparty did not, assesses no-show penalty against absent party.<br>5. Releases fleet availability and dispatches cancellation alerts. | Auto-cancellation triggers; no penalty when neither acted; no-show penalty applied to absent party when one acted; notifications sent. |
| `AR-07-MOB` | Auto-cancellation UI updates | P1 | Mobile app handles the cancellation status transition gracefully, updating screens to show the timeout status. | 1. Booking view updates status to `CANCELLED` (Timeout).<br>2. Removes handover confirmation buttons.<br>3. Displays explanation banner: "Handover timed out after 48 hours." | Handover controls removed; status reflects timeout cancel clearly. |
| `AR-08-BE` | Handover proof & recipient attestation | P0 | Backend captures verifiable handover proof from the receiving party at both pickup and return, blocking the `ACTIVE` transition until attestation and checklist are complete. | 1. `POST /bookings/{id}/handover/attest` records recipient full name, role, signature image or phone OTP, and completed condition checklist (fuel level, visible defects, attachments present).<br>2. Persists records immutably with cryptographic hash.<br>3. Enforces identical attestation at return via `POST /bookings/{id}/return/attest`.<br>4. Blocks transition to `ACTIVE` until complete. | Handover to `ACTIVE` blocked until both parties sign or OTP-confirm and checklist is complete; records stored immutably. |
| `AR-08-MOB` | Handover attestation & checklist UI | P0 | Mobile app guides the receiving party through capturing name/role, signature or OTP, and the condition checklist before enabling handover. | 1. Screen presents handover checklist (fuel gauge slider 0-100%, tool verification, physical defect toggles).<br>2. Collects recipient name and role.<br>3. Captures signature on digital canvas OR prompts for phone OTP sent to receiving party.<br>4. Enables "Confirm Handover" only upon complete attestation. | Confirm button gated on complete attestation + checklist; signature/OTP captured; uploads to API on confirm. |
| `AR-09-BE` | Incident report & liability capture | P0 | Backend records mid-rental incidents with fault attribution and drives admin escalation and billing adjustments. | 1. `POST /bookings/{id}/incident` accepts type (`accident`, `breakdown`, `theft`, `injury`), description, photos, and `claimed_fault` (`owner`, `renter`, `unknown`).<br>2. High-severity types auto-notify admin with 24h SLA.<br>3. If admin attributes downtime to owner fault, suspends billing clock and excludes hours from excess calculations (`RC-03`). | Incident persisted with type/fault/photos; high-severity notifies admin within 24h SLA; owner-fault downtime suspends billing and excludes excess hours. |
| `AR-09-MOB` | Incident report form UI | P0 | Mobile form for either party to file an incident with type, description, photos, and claimed fault. | 1. User opens "Report Incident".<br>2. Selects incident type.<br>3. Enters description and attaches photos.<br>4. Selects claimed fault.<br>5. Submits incident payload to backend. | Form validates required fields; photos attach; submit calls incident API and shows confirmation. |
| `AR-10-BE` | Operator assignment consent & revocation | P1 | Backend solicits and records operator consent on assignment/change and blocks `ACTIVE` without an accepted required operator. | 1. On operator assignment or reassignment, dispatches notification with Accept/Decline actions.<br>2. Blocks transition to `ACTIVE` while assigned required operator has not accepted.<br>3. Supports operator consent revocation.<br>4. Re-solicits GPS tracking consent upon mid-rental operator changes (`GPS-04`).<br>5. Logs all decisions in `OperatorConsentLog`. | Operator receives Accept/Decline; `ACTIVE` blocked without acceptance; revocation exposed; mid-rental change re-solicits GPS consent; actions logged. |
| `AR-10-MOB` | Operator consent response UI | P1 | Operator mobile UI to accept or decline an assignment and to revoke consent. | 1. Operator receives assignment alert.<br>2. Reviews job site coordinates, machine specs, and dates.<br>3. Taps "Accept Job" or "Decline".<br>4. Settings screen offers "Revoke Job Consent" with confirmation modal. | Accept/Decline calls consent API; revocation available; UI reflects consent status accurately. |

---

### 5.7 Rental Completion, Return & Damage

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `RC-01-BE` | Auto-return & manual return trigger APIs | P0 | Backend API endpoint and cron job to transition active bookings to `RETURN_INITIATED` status. | 1. `POST /bookings/{id}/return/initiate` validates caller is renter or owner.<br>2. Scheduled cron job queries active bookings where `end_date < NOW()`.<br>3. Transitions booking status to `RETURN_INITIATED`.<br>4. Dispatches return inspection reminders to both parties. | Status transitions to `RETURN_INITIATED` via API or cron execution; dispatches inspection notifications. |
| `RC-01-MOB` | Return initiation trigger UI | P0 | Mobile app provides an "Initiate Return" button during active rentals and displays status updates on completion. | 1. Active rental dashboard displays "Initiate Return" button.<br>2. Tapping displays confirmation prompt.<br>3. Calls return trigger API.<br>4. Navigates user directly into return inspection stepper. | Button visible and updates status to `RETURN_INITIATED` on click; routes to inspection stepper. |
| `RC-02-BE` | Store return inspection data API | P0 | Backend API receives final inspection photos and hour meter reading from both owner and renter. | 1. `POST /bookings/{id}/inspection/post` receives return hour meter value, mandatory meter photo, and 4-angle machinery condition photos.<br>2. Validates files and records entries in `BookingInspection`.<br>3. Triggers monotonic meter regression check (`RC-12`). | Images and final reading persisted in DB; blocks progression without valid meter reading and photo. |
| `RC-02-MOB` | Return inspection upload UI | P0 | Mobile app prompts users to upload final condition photos and record final hour meter reading. | 1. Mobile inspection wizard enforces camera-only capture of hour meter.<br>2. Prompts entry of numeric reading.<br>3. Prompts capture of 4-angle condition photos.<br>4. Submits return inspection payload to backend. | Interface enforces inputs; uploads to API upon confirmation; blocks completion without required photos. |
| `RC-03-BE` | Calculate excess hours API logic | P0 | Backend calculates excess hour charges based on standard daily hours and agreed total hours limit. | 1. Computes `agreed_hours = days * standard_daily_hours` (default 8h/day).<br>2. Computes `actual_hours = return_meter - start_meter`.<br>3. If `actual_hours > agreed_hours`, computes `excess_fee = (actual - agreed) * excess_hour_rate`.<br>4. Appends excess fee to booking invoice. | Excess hour fee added correctly based on agreed limits; calculation strictly uses integer minor units. |
| `RC-03-MOB` | Display excess hours breakdown UI | P0 | Mobile app displays excess hour calculation and standard daily hours limits in the billing summary. | 1. Reads standard daily hours limit, agreed total limit, and actual hours recorded.<br>2. Displays itemized line-item: "Excess Operating Hours: X hrs @ [Rate] ETB = [Total] ETB".<br>3. Shows baseline and final hour meter photos side by side. | Itemized excess hours display clearly in UI breakdown along with standard daily limits and photos. |
| `RC-04-BE` | Damage claim submission & response API | P0 | Backend API supports submitting damage claims (amount, photos, text) and recording renter's response (accept or dispute). | 1. `POST /bookings/{id}/damage` accepts owner claim (amount, damage photos, description).<br>2. Saves claim in `DamageClaim` table.<br>3. `PUT /bookings/{id}/damage/respond` records renter response (`ACCEPTED` or `DISPUTED`).<br>4. Disputed claim creates `Dispute` record with origin `return_damage`. | Claims logged in DB; status transitions update accordingly; disputed claims trigger formal dispute workflow. |
| `RC-04-MOB` | Submit & respond to damage claim UI | P0 | Mobile screens for owner to fill out damage claim details, and renter to review, accept, or dispute. | 1. Owner return view displays "Report Damage" action.<br>2. Form captures damage description, estimated repair cost, and photos.<br>3. Renter view presents damage claim alert with photos and cost.<br>4. Renter options: "Accept & Pay" or "Dispute Claim". | Forms and review screens validate inputs; submit actions call API; renders clear photographic evidence. |
| `RC-05-BE` | Final invoice calculation logic | P0 | Backend computes final payable amount using rental base fees, excess hours, accepted damages, prepayments, and 15% VAT on all charges. | 1. Computes total charges: `(rental_fee + excess_hours + accepted_damages)`.<br>2. Calculates 15% VAT on total charges (unless renter is `TaxExempt`).<br>3. Deducts verified prepayment amount.<br>4. Final balance due $= \text{total\_inclusive} - \text{prepayment}$.<br>5. Saves invoice in `BookingInvoice`. | Arithmetic matches calculation rules across positive/negative totals, applying 15% VAT to all post-rental adjustments. |
| `RC-05-MOB` | Display final checkout summary UI | P0 | Mobile app renders the itemized final amount summary showing all adjustments, 15% VAT calculations, and final balance due. | 1. Fetches final invoice payload.<br>2. Itemizes base rental, excess hour charges, accepted damage fees, 15% VAT, and prepayment credit.<br>3. Displays bold final balance due.<br>4. Provides "Pay Final Balance" CTA if balance $> 0$. | Detailed breakdown including VAT is rendered correctly with correct line items; clear final balance. |
| `RC-06-BE` | Final payment upload & confirm API | P0 | Backend APIs for uploading final payment proof and for owners to approve, updating status to `COMPLETED`. | 1. `POST /bookings/{id}/final-payment/proof` uploads bank slip, setting status to `PENDING_FINAL_PAYMENT`.<br>2. Owner verifies slip and calls `PUT /bookings/{id}/final-payment/confirm`.<br>3. Transitions booking status to `COMPLETED`.<br>4. Records transaction in unified financial ledger (`CW-17-BE`).<br>5. Confirm action is idempotent per `BK-17-BE`. | Status updates to `COMPLETED` on owner confirmation; final transaction logged to ledger; confirm idempotent; evidence deduped. |
| `RC-06-MOB` | Renter final payment & owner verification UI | P0 | Mobile screens for renter to upload final payment slip, and owner to verify and confirm receipt. | 1. Renter uploads transfer slip for final balance.<br>2. Owner receives notification, inspects slip image.<br>3. Owner taps "Confirm Final Payment".<br>4. Booking transitions to completed celebration view. | Upload screens function; owner controls trigger confirmation; views update to completed state. |
| `RC-07-BE` | Create reviews & update rating status APIs | P0 | Backend endpoints to submit reviews and stars, updates aggregate scores, and transitions booking status to `RATED` once both rate. | 1. `POST /bookings/{id}/review` stores 1-5 star score and text review.<br>2. Enforces `AB-01-BE` eligibility: booking must be `COMPLETED` or `RATED`.<br>3. Once both parties submit (or 7-day timeout elapses), status transitions to `RATED`.<br>4. Updates rolling average rating on user and fleet profiles. | Review data saved; status turns `RATED` only after both submit; averages update; duplicate/ineligible reviews rejected. |
| `RC-07-MOB` | Rating & review feedback screens | P0 | Mobile rating UI screens with 1-5 star selector widgets and comments field. | 1. Booking completion prompts review modal.<br>2. Star rating widget allows selecting 1 to 5 stars.<br>3. Text area captures written feedback.<br>4. Submits review to backend. | Rating selector works; review submission calls backend successfully; updates local review state. |
| `RC-07B-BE` | Renter-to-operator rating | P1 | Backend allows the renter to rate the assigned operator for a completed booking, distinct from the owner review. | 1. `POST /bookings/{id}/review/operator` accepts 1-5 star score and comment for assigned operator.<br>2. Enforces one review per booking.<br>3. Updates operator aggregate rating in `OperatorProfile`. | Renter can rate operator once per booking; operator average updates; ineligible/duplicate submissions rejected. |
| `RC-08-BE` | Admin partial damage dispute resolution API | P1 | Backend API endpoint allowing admins to set a resolved dispute amount, update invoice, and notify users. | 1. `PUT /disputes/{id}/resolve-damage` accepts `resolved_amount`.<br>2. Transitions booking status to `DAMAGE_PARTIAL_RESOLVED`.<br>3. Appends resolved amount to final invoice charges.<br>4. Auto-advances booking to `PENDING_FINAL_PAYMENT`.<br>5. Dispatches notifications to both parties. | Invoice updated with partial amount; notifications dispatched; booking status set to `DAMAGE_PARTIAL_RESOLVED` and auto-advances to `PENDING_FINAL_PAYMENT`. |
| `RC-08-MOB` | Dispute resolution display on mobile | P1 | Mobile app displays updated resolution amount and status changes on the booking detail and billing screens. | 1. Displays banner: "Damage dispute resolved by Admin: [Amount] ETB approved".<br>2. Final invoice updates to include approved damage line.<br>3. Renter prompted to pay updated final balance. | Status updates to show final amount on mobile client; payment CTA links to final balance settlement. |
| `RC-08-WEB` | Admin damage dispute resolution UI | P1 | Next.js admin dashboard interface displaying damage claim info and a text field to enter a custom resolved amount. | 1. Admin opens return damage dispute.<br>2. Reviews photos, damage estimate, and renter rebuttal.<br>3. Enters binding numeric resolved amount.<br>4. Submits ruling with mandatory justification and admin MFA. | Form validation prevents negative values or amounts exceeding original claim; form triggers API and updates status. |
| `RC-09-BE` | Review moderation API | P1 | API to transition review status to `FLAGGED` and record audit log reason. | 1. `POST /reviews/{id}/flag` sets status to `FLAGGED`.<br>2. Records moderator ID and justification in `AuditLog`.<br>3. Excludes flagged review from public fleet and profile rating queries.<br>4. Recomputes aggregate rating without flagged score. | Status updates; review excluded from public endpoints; audit log created; profile averages recalculate. |
| `RC-09-WEB` | Admin review moderation dashboard | P1 | Next.js admin interface showing list of public reviews with controls to flag or hide items. | 1. Admin navigates to Content $\rightarrow$ Review Moderation.<br>2. Table displays reported or all public reviews.<br>3. Action button: "Flag / Hide Review".<br>4. Prompts justification modal and submits. | Moderation tools trigger API; reviewed items vanish from user-facing screens immediately. |
| `RC-10-BE` | Review response API logic | P1 | API logic to store owner response to a review, validating it is from the fleet owner and does not exceed character limits. | 1. `POST /reviews/{id}/response` checks caller is fleet owner.<br>2. Asserts no prior response exists (blocks duplicates).<br>3. Enforces text length $\le 500$ characters.<br>4. Persists response in `ReviewResponse` table. | Response stored; duplicate attempts rejected; length validated; displays publicly beneath review. |
| `RC-10-MOB` | Owner review reply input UI | P1 | Mobile app displays a "Reply" button under reviews in the owner's fleet detail panel with character-limited text box. | 1. Owner views reviews on fleet listing.<br>2. Taps "Reply" under a renter review.<br>3. Input box displays 500-character counter.<br>4. Submits reply. | Input box validates character limit; reply displays immediately under the parent review. |
| `RC-11-BE` | Rating fallback cron | P1 | Backend cron job auto-advances bookings to `RATED` if ratings are missing 7 days post-completion. | 1. Scheduled cron runs daily.<br>2. Selects bookings in `COMPLETED` older than 7 days lacking reviews.<br>3. Inserts rating record with `NULL` score (preventing score impact).<br>4. Updates booking status to `RATED`. | Status updates to `RATED` automatically; database records NULL values; profile averages are unaffected. |
| `RC-12-BE` | Hour Meter Discrepancy Checker | P1 | Flag bookings for Admin review if renter and owner inspection readings conflict. | 1. Compares renter and owner submitted hour meter readings at handover or return.<br>2. Triggers discrepancy if $|\text{renter\_reading} - \text{owner\_reading}| > \max(1.0\text{h}, 0.05 \times \text{hours})$.<br>3. Sets booking status to `DISPUTED_HOUR_METER` and suspends transitions.<br>4. Generates Admin dispute ticket. Admin ruling is binding. | Conflicting hour readings automatically set to `DISPUTED_HOUR_METER` and locked; only Admin can resolve; resolved reading is final. |
| `RC-12-MOB` | Hour Meter Discrepancy Banner | P1 | Display a warning banner to users when hour meter readings are disputed. | 1. Checks booking status.<br>2. If `DISPUTED_HOUR_METER`, displays prominent warning banner: "Hour meter readings in dispute. Admin is reviewing."<br>3. Disables final payment and completion buttons. | Warning banner renders and locks operations for renters and owners until resolved. |
| `RC-12-WEB` | Admin Hour Meter Dispute UI | P1 | Next.js admin dashboard to resolve hour meter disputes visually. | 1. Admin opens hour meter dispute ticket.<br>2. Displays renter and owner submitted numbers and meter photos side by side.<br>3. Admin inputs binding resolved hour reading.<br>4. Submits ruling, updating final invoice and resuming billing flow. | Form validates entry; updates database with resolved reading; resumes billing flow automatically. |
| `RC-13-BE` | Damage valuation & liability rules | P0 | Backend enforces evidence, valuation basis, and caps on damage claims, and flags renters who fail to pay the final invoice. | 1. For claims $>$ `damage_evidence_threshold`, owner MUST attach $\ge 1$ supporting quote/estimate.<br>2. Admin ruling records `valuation_basis` enum (`repair_quote`, `assessor`, `depreciated_value`, `admin_estimate`).<br>3. Claim capped at asset `replacement_value`; normal wear-and-tear excluded.<br>4. If renter fails to pay final invoice within 3 business days, flags renter `payment_standing = DELINQUENT`, blocks new bookings, and records receivable. | Over-threshold claims rejected without evidence; ruling stores valuation basis; amount $\le$ replacement value; non-payment flags renter DELINQUENT. |
| `RC-14-BE` | Return-dispute lifecycle & SLA | P0 | Backend creates a first-class `Dispute` for return/hour-meter conflicts and enforces the SLA/escalation. | 1. On renter `DISPUTED` damage claim or hour-meter discrepancy, creates `Dispute` row with `origin = return_damage` or `origin = hour_meter`.<br>2. Sets `dispute_deadline = created_at + 3 business days`.<br>3. Booking transitions suspended until resolution.<br>4. Daily cron flags overdue disputes and escalates to Super Admin. | Dispute created with correct origin and deadline; booking suspended until resolved; overdue disputes flagged and escalated by cron. |
| `RC-15-BE` | Early termination & force-majeure | P1 | Backend supports admin-adjudicated early termination with pro-rata recompute and conditional penalty waiver. | 1. Admin adjudicates mutual-agreement or force-majeure termination.<br>2. Recomputes rental fee pro-rata based on actual operational hours up to termination.<br>3. Waives cancellation penalty on verified force-majeure.<br>4. Unused prepayment refunded via standard owner obligation evidence flow (`BK-CX-04-BE`). | Admin can flag mutual/force-majeure termination; fee recomputed to termination point; penalty waived on qualifying force-majeure; refund tracked as owner obligation. |

---

### 5.8 Logistics Transport Integration

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `LG-01-BE` | Logistics vehicle registration API support | P0 | Backend API supports registering vehicles in LOGISTICS category, applying specialized specifications template. | 1. Handles fleet creation POST request.<br>2. Validates category parameters correspond to `LOGISTICS` taxonomy.<br>3. Associates logistics-specific dynamic specs (payload capacity, bed dimensions, ramp type).<br>4. Indexes specs into `FleetSpecIndexed`. | Fleet created with logistics specs; category index updated correctly; payload attributes searchable. |
| `LG-01-MOB` | Mobile logistics registration UI | P0 | Mobile app registration form allows selecting Logistics category and displays spec forms. | 1. Owner selects "Logistics / Transport" in category dropdown.<br>2. Form dynamically updates to display carrier fields (payload capacity, deck length, axle count).<br>3. Enforces required fields.<br>4. Submits carrier registration. | Dynamic form sections display correctly; fields are mandatory; submits carrier listing. |
| `LG-02-BE` | Match logistics API endpoint | P1 | Backend endpoint computes and suggests logistics vehicles matching payload constraints of booked fleet. | 1. `GET /bookings/{id}/logistics-suggestions` extracts operating weight and dimensions of booked machine from `FleetSpecIndexed`.<br>2. Queries active logistics fleets where `payload_capacity >= machine_weight`.<br>3. Filters by proximity to machine depot and job site.<br>4. Returns ranked list of compatible carriers. | API results correctly filter by specs (weight vs payload capacity); returns compatible vehicles only. |
| `LG-02-MOB` | Logistics recommendation list UI | P1 | Mobile app shows a list of suggested logistics fleets to the owner after booking confirmation. | 1. Confirmed booking screen presents "Arrange Transport" card.<br>2. Renders horizontal carousel of recommended carrier fleets with payload match badges.<br>3. Each card displays carrier photo, daily/km rate, and "Book Transport" CTA. | Cards display correct match details and enable booking sub-flow trigger. |
| `LG-03-BE` | Sub-booking creation API | P1 | Backend processes logistics sub-booking, linking parent ID, and enforcing the standard prepayment/payment-evidence and commission rules. | 1. `POST /bookings` accepts `parent_booking_id` reference.<br>2. Machine owner acts as renter; carrier owner acts as fleet owner.<br>3. Links sub-booking to parent.<br>4. Prepayment evidence upload and commission deduction follow standard booking lifecycle. | Sub-booking records parent ID; prepayment-evidence and commission computed according to global platform rates. |
| `LG-03-MOB` | Book logistics sub-flow UI | P1 | Mobile interface lets owner select a carrier, input details, and finalize sub-booking with the standard prepayment-evidence step. | 1. Owner selects recommended carrier.<br>2. Selects mobilization and demobilization transport dates.<br>3. Reviews transport cost breakdown.<br>4. Submits booking request.<br>5. Uploads prepayment transfer evidence to carrier owner. | Sub-flow triggers sub-booking; redirects owner to standard payment upload screen cleanly. |
| `LG-04-BE` | Flag external logistics API | P1 | Endpoint updates booking details to indicate logistics are being managed externally. | 1. `PATCH /bookings/{id}/logistics` accepts `external_logistics = TRUE` and transport notes.<br>2. Updates booking record.<br>3. Removes pending logistics prompts from dashboard. | Boolean flag set; status transitions successfully; prevents redundant logistics alerts. |
| `LG-04-MOB` | External logistics option UI | P1 | Mobile app displays button for owner to bypass logistics selection and record as handled externally. | 1. Logistics prompt displays "I will manage transport externally" button.<br>2. Prompts brief explanation modal (carrier company, driver contact).<br>3. Sends update to API.<br>4. Dashboard updates to show external logistics confirmed. | Option functions; UI updates to show external selection; removes booking prompt. |
| `LG-05-BE` | Sub-booking reverse cascade cancellation logic | P1 | Backend auto-cancels linked logistics bookings on main cancellation, and handles logistics provider cancellation viability locks. | 1. If main booking cancels, automatically cancels linked sub-bookings (`parent_booking_id`).<br>2. If carrier cancels sub-booking, sets `logistics_failed = TRUE` on main booking, notifies renter, and waives machinery cancellation penalty.<br>3. Renter prepayment refunded on machinery auto-cancel.<br>4. `PROCEED_WITHOUT` action clears flag and detaches sub-booking.<br>5. Auto-cancels machinery if unresolved at $\min(\text{failure} + 24\text{h}, \text{rental\_start})$. | Bidirectional cancellation logic works; flags machinery booking on logistics failure; prepayment refunded on machinery auto-cancel; auto-cancels at deadline. |
| `LG-05-MOB` | Reverse cascade cancellation UI alerts | P1 | Mobile UI alerts renter if logistics provider cancels sub-booking, prompting for action. | 1. Displays prominent warning alert on machinery booking screen if `logistics_failed = TRUE`.<br>2. Renders actions: "Re-book Logistics" or "Proceed without Logistics".<br>3. Shows countdown timer to auto-cancellation deadline. | Warning details display correctly; navigation actions allow user to select solutions or cancel. |

---
### 5.9 Commission, Revenue & Owner Payouts

> **MVP Commission Architecture:** The platform operates under a wallet-debit model. The owner pre-funds an `OwnerWallet`. Upon booking confirmation (`CONFIRMED`), platform commission is automatically debited from the owner's available wallet balance (`COMMISSION_DEBIT`). The platform never takes custody of, holds, or disburses renter rental funds under the `MANUAL_TRANSFER` settlement method.

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `CW-00-BE` | Fund-flow model & settlement-method framework | P0 | Canonically define how money moves: a dual settlement-method framework — manual Direct Transfer Model plus pluggable payment gateways (`PAY-GW-01..05`) — governing every financial transaction. | 1. Every payment leg executes via `settlement_method ∈ {MANUAL_TRANSFER, GATEWAY}`.<br>2. `MANUAL_TRANSFER` (Direct Transfer Model): Renter pays owner directly via bank transfer (verified via receipt upload); platform never holds or disburses renter funds; platform debits commission from owner's `OwnerWallet` upon confirmation (`CONFIRMED`); refunds and owner-to-renter penalties are owner obligations tracked by the platform.<br>3. `GATEWAY`: Funds flow per the gateway's declared custody model (`DIRECT_TO_OWNER` or `PLATFORM_HELD`).<br>4. Deposit is informational-only in MVP (excluded from subtotal/VAT/total/prepayment).<br>5. 15% VAT applies to rental fee + attachments only; issued under owner TIN.<br>6. Every movement writes to the unified ledger (`CW-17-BE`) and is idempotent (`BK-17-BE`). | Framework documented as single source of truth; `MANUAL_TRANSFER` semantics active at launch; no platform custody outside `PLATFORM_HELD` gateways; commission path unambiguous; all movements ledgered and idempotent. |
| `PAY-GW-01-BE` | Pluggable Payment Gateway Framework | P1 | A gateway-agnostic adapter interface and registry so payment gateways are plug-and-use via configuration without code changes outside the adapter. | 1. Define `PaymentGatewayAdapter` interface (`initiate`, `verifyCallback`, `queryStatus`, `refund`).<br>2. Gateways stored in `PaymentGateway` registry table (`gateway_key`, `display_name`, `adapter_type`, `config_json`, `is_enabled`, `supported_countries[]`, `supported_currencies[]`, `custody_model ∈ {DIRECT_TO_OWNER, PLATFORM_HELD}`, `capabilities`, `priority`).<br>3. Gateway transaction logs stored in `GatewayTransaction` table.<br>4. Adding gateway = implement adapter + insert registry row via `PAY-GW-05-WEB`. | New gateway usable by configuration and adapter alone; supports multiple concurrent gateways; each attempt logged to `GatewayTransaction`; registry enforces scoping and custody models. |
| `PAY-GW-02-BE` | Checkout settlement-method selection policy | P1 | Backend policy resolving which settlement methods/gateways a renter sees at checkout, admin-configured per country/currency. | 1. `GET /bookings/{id}/payment-options` returns active methods: `MANUAL_TRANSFER` (with owner bank coordinates) plus enabled gateways matching country/currency, ordered by priority.<br>2. Evaluates `payment.manual_transfer_enabled` and gateway scoped rows.<br>3. Rejects configurations that would leave a country with zero settlement methods.<br>4. Persists renter choice on payment leg (`payments.settlement_method`, `payments.gateway_key`). | Options endpoint returns exactly the admin-enabled, country/currency-matched set; renter choice persisted per leg; zero-method configurations blocked at configuration save time. |
| `PAY-GW-02-MOB` | Checkout payment-method selection UI | P1 | Flutter checkout step rendering the resolved settlement options dynamically without hardcoded gateway lists. | 1. Calls `/bookings/{id}/payment-options` and renders returned options (bank transfer with owner coordinates; gateways with logos and display names).<br>2. Selecting `MANUAL_TRANSFER` routes to slip upload flow (`BK-07-MOB`).<br>3. Selecting a gateway launches checkout URL/USSD/deeplink and polls status until callback resolves.<br>4. Pre-selects default option if only one resolves. | Options render dynamically from API; manual path functions as defined; gateway flow launches external app/webview and reflects verified webhook result; failure permits method retry. |
| `PAY-GW-03-BE` | Gateway payment lifecycle & webhook processing | P1 | Initiation, verified callbacks, and state-machine integration for `GATEWAY` payment legs — idempotent and spoof-proof. | 1. `POST /payments/{id}/gateway/initiate` calls adapter `initiate` and returns checkout payload.<br>2. `POST /payments/callback/{gateway_key}` validates cryptographic signature, matches amount and currency, and checks idempotency on `gateway_txn_ref`.<br>3. On verified success: marks payment leg paid; advances booking state automatically (`PENDING_PAYMENT` $\rightarrow$ `CONFIRMED`); writes ledger entry; debits commission per custody model.<br>4. Failed attempts remain `PENDING_PAYMENT` until standard timeout.<br>5. Unverifiable signatures rejected with HTTP 401 and security alerts triggered. | Spoofed/duplicate webhooks rejected via signature and idempotency; verified success auto-advances booking; failures remain retryable; ledger and commission posted correctly. |
| `PAY-GW-04-BE` | PLATFORM_HELD custody accounting & owner payout | P2 | Escrow-style accounting for gateways whose custody model is `PLATFORM_HELD`: platform receives funds, deducts commission at source, and pays out the owner. | 1. On verified settlement: records gross receipt to platform clearing account; deducts commission at source (`TOPUP_PAYMENT_SPLIT`); creates `OwnerPayout` record (`owner_id`, `net_amount`, `status = PENDING`).<br>2. Payouts execute via gateway payout API or manual bank transfer with evidence.<br>3. Reconciles daily against gateway statements.<br>4. Licensing gate: enabling `PLATFORM_HELD` requires recorded regulatory license sign-off in audit log. | Gross/commission/net ledgered separately and reconcile daily; owner payout tracked to execution with evidence; commission deducted at source; blocked without regulatory sign-off. |
| `PAY-GW-05-WEB` | Payment Gateway & Settlement Admin Console | P1 | Single Super-Admin console to manage gateways, credentials, and settlement policy dynamically. | 1. List registered gateways: enable/disable, edit credentials, assign countries/currencies, view custody models.<br>2. Toggle `payment.manual_transfer_enabled` per country.<br>3. Reject invalid configs: zero settlement methods for a country, missing licensing sign-off for `PLATFORM_HELD`.<br>4. Guarded by `config:manage` permission and Super Admin MFA. | Gateway management executed without code deployment; contradictory configurations blocked; all changes audit-logged with immediate propagation. |
| `CW-01-BE` | Transaction processing & commission calculator service | P0 | Backend processes commission calculations and wallet deductions upon booking confirmation events. | 1. Event listener captures `BOOKING_CONFIRMED`.<br>2. Calculates commission over rental base fee (`CW-14-BE`).<br>3. In same database transaction, creates `debit` row in `OwnerWalletTransaction` with type `COMMISSION_DEBIT` and decrements `OwnerWallet.balance`.<br>4. Enforces idempotency via unique constraint on `(booking_id, COMMISSION_DEBIT)`. | Debit posted on confirmation; commission rate matches active tier/rule; exactly one debit per booking enforced; rolls back booking if transaction fails. |
| `CW-01-MOB` | Display wallet debit transactions | P0 | Mobile wallet ledger displays the deducted commission amount with booking details. | 1. Owner opens Wallet tab.<br>2. App fetches paginated transaction history.<br>3. Renders commission debit row showing booking reference, timestamp, minor-unit amount in ETB, and remaining balance. | Displays commission debit type, booking reference, negative amount styling, and post-transaction balance accurately. |
| `CW-02-BE` | Get wallet balance & transaction ledger API | P0 | API to return wallet balance summary (balance, negative limit, status) and paginated list of ledger transactions. | 1. Exposes `GET /wallet`.<br>2. Retrieves owner wallet metrics (`balance_minor`, `negative_limit_minor`, `status`).<br>3. Fetches paginated rows from `OwnerWalletTransaction`.<br>4. Returns payload conforming to `SuccessEnvelope`. | Returns accurate minor-unit balance metrics and paginated transaction history; unauthorized callers rejected. |
| `CW-02-MOB` | Wallet dashboard UI | P0 | Mobile app wallet dashboard screen displaying balance metrics and transaction list. | 1. Owner clicks "Wallet" tab.<br>2. Displays primary balance card (Current Balance, Allowed Credit Limit, Status Badge).<br>3. Renders "Top Up Wallet" CTA.<br>4. Lists transactions with infinite scroll pagination. | Layout renders formatted minor-unit currency; status badges show active/locked state; pagination functions smoothly. |
| `CW-03-BE` | Submit wallet top-up request API | P0 | API endpoint to submit top-up requests (amount, bank slip image) and admin approval endpoint to execute the credit. | 1. `POST /wallet/topup` accepts minor-unit amount and transfer receipt image; stores record in `PENDING` state.<br>2. Admin reviews receipt at `PATCH /admin/wallet/topup/{id}/approve`.<br>3. In single `pgx.Tx`, updates top-up status to `APPROVED`, adds credit row to `OwnerWalletTransaction` (`TOPUP_MANUAL`), and increments wallet balance.<br>4. If wallet was in `LOW_BALANCE_LOCK` and balance rises above limit, restores owner status to `ACTIVE`. | Top-up request persisted; admin approval credits wallet atomically; account lockout cleared if balance restored; duplicate approvals blocked. |
| `CW-03-MOB` | Upload top-up slip UI | P0 | Mobile dashboard for owners to input top-up amounts and upload bank transfer receipts. | 1. Owner taps "Top Up Wallet".<br>2. Displays platform bank coordinates for manual transfer.<br>3. Form captures numerical amount.<br>4. Launches camera/gallery to pick transfer slip.<br>5. Submits top-up request payload to backend. | Form validates numeric bounds; receipt image uploaded securely; shows "Pending Approval" state upon submission. |
| `CW-03-WEB` | Admin wallet top-up verification queue | P0 | Next.js admin page displaying pending wallet top-up requests with verification actions. | 1. Admin navigates to Finance $\rightarrow$ Top-Up Queue.<br>2. Table lists pending requests with date, owner name, claimed amount, and slip thumbnail.<br>3. Clicking slip opens high-resolution preview modal.<br>4. Action buttons: "Approve Top-Up" or "Reject" (with reason). | Receipt preview displays clearly; approval executes backend credit API; rejection captures reason and updates dashboard live. |
| `CW-04-BE` | Admin wallet configuration & audits APIs | P0 | API endpoints to fetch all wallets, update owner negative limits, and log administrative audit changes. | 1. `PUT /admin/owners/{id}/wallet-limit` accepts new `negative_limit_minor`.<br>2. Validates caller has `wallet:manage` permission and active MFA.<br>3. Updates wallet record in DB.<br>4. Logs before/after limit diff in `AuditLog`. | Database limit updates; audit log record created with admin identity and timestamp; unauthorized updates rejected with HTTP 403. |
| `CW-04-WEB` | Admin wallet control panel UI | P0 | Next.js admin dashboard to manage owner wallets, change negative limits, and approve transactions. | 1. Admin opens Finance $\rightarrow$ Owner Wallets.<br>2. Table lists owners, current balances, negative limits, and statuses.<br>3. Clicking "Edit Limit" opens modal to adjust credit threshold.<br>4. Submits update via Server Action. | Table supports search and filtering; modal validates positive integer input; updates reflect immediately in UI. |
| `CW-05-BE` | Tiered commission rate calculator | P1 | Backend commission service applies tiered rates based on booking value ranges stored in commission rules. | 1. Reads active rules from `CommissionRule` table.<br>2. Matches booking rental fee against half-open tier bounds `[min, max)`.<br>3. Computes commission using integer basis points (`commission_bps`).<br>4. Falls back to global base rate if no tier matches. | Tier bounds evaluated correctly; integer basis point math used; calculation matches tier specification exactly. |
| `CW-05-WEB` | Dynamic commission tiers configuration panel | P1 | Next.js admin panel settings to create, edit, or delete tiered commission ranges. | 1. Admin navigates to Settings $\rightarrow$ Commission Tiers.<br>2. Interface displays table of active tiers with value bounds and basis points.<br>3. "Add Tier" modal captures min amount, max amount, and rate (%).<br>4. Validates that adjacent tiers do not overlap or leave gaps. | Form enforces non-overlapping tier bounds; saves to database; changes apply to subsequent bookings without restart. |
| `CW-06-BE` | PDF statement generator service | P1 | Backend worker generates PDF receipt upon ledger transaction and stores file in SeaweedFS S3. | 1. Asynchronous worker triggered upon transaction write.<br>2. Compiles monthly statement or transaction receipt including opening balance, transactions, debits, credits, and closing balance.<br>3. Generates PDF with fiscal metadata.<br>4. Stores file in SeaweedFS and updates `statement_url`. | Statement PDF formatted cleanly; stored in encrypted S3 bucket; download link accessible to authorized owner. |
| `CW-06-MOB` | Owner PDF statement download UI | P1 | Mobile app enables owners to download transaction invoices and monthly statements from the ledger UI. | 1. Owner views wallet dashboard.<br>2. Taps "Export Statement" and selects date range/month.<br>3. Calls statement export API.<br>4. Downloads PDF and launches native device viewer/share sheet. | Statement file downloads cleanly; launches viewer; handles network errors gracefully. |
| `CW-06-WEB` | Admin PDF statement viewer UI | P1 | Next.js admin dashboard supports viewing and downloading PDF receipts from wallet logs. | 1. Admin views owner wallet transaction ledger.<br>2. Each transaction row includes "Download PDF" action icon.<br>3. Clicking streams PDF document from object storage. | Streams statement directly; opens in browser preview tab cleanly. |
| `CW-07-BE` | Restrict commission settings updates API | P0 | Backend rejects any settings modification requests that do not originate from an administrator role. | 1. Middleware inspects caller authentication context.<br>2. Restricts `POST`, `PUT`, `PATCH`, and `DELETE` on `/admin/commission-rules` to users with `super_admin` role and `config:manage` permission.<br>3. Enforces step-up MFA verification within last 15 minutes. | Non-admin calls return HTTP 403 Forbidden; missing MFA returns HTTP 401; authorized calls proceed. |
| `CW-07-MOB` | Hide commission settings configuration UI | P0 | Mobile client ensures owner settings screen does not contain tools to view or edit commission rules. | 1. Mobile layout renderer inspects user permissions.<br>2. Completely omits commission configuration menus and controls from owner application bundles. | Settings menu contains zero commission editing options; API calls to commission management blocked. |
| `CW-07-WEB` | Admin panel commission config page | P0 | Next.js admin dashboard includes management tools for setting up global commission rates. | 1. Super Admin navigates to Settings $\rightarrow$ Platform Commission.<br>2. Renders form displaying default commission basis points.<br>3. Submitting update requires entering current password and TOTP code.<br>4. Persists new rate and logs audit record. | Config page accessible only to Super Admin; requires MFA challenge; updates database global settings. |
| `CW-08-BE` | Account status monitor & listing lockout backend logic | P0 | Backend updates owner status to `LOW_BALANCE_LOCK` and hides listings if wallet balance drops below negative limit. | 1. Following any wallet debit: evaluates `if balance_minor < -negative_limit_minor`.<br>2. If breached, updates owner profile status to `LOW_BALANCE_LOCK`.<br>3. Public search and catalog queries filter out listings belonging to locked owners.<br>4. Existing confirmed bookings continue unaffected. | Listings automatically hidden from search upon credit limit breach; owner status set to `LOW_BALANCE_LOCK`; confirmed bookings protected. |
| `CW-08-MOB` | Low balance lock warning banner UI | P0 | Mobile app alerts owner that listings are deactivated due to low balance, hiding active booking actions. | 1. Owner app detects `LOW_BALANCE_LOCK` status.<br>2. Displays persistent red warning banner across dashboard: "Listings paused due to negative wallet balance. Top up to reactivate."<br>3. Directs owner to top-up screen. | Banner displays prominently; listing creation tools disabled; clear path to top up. |
| `CW-08-WEB` | Admin console account lock indicators | P0 | Next.js admin dashboard identifies locked owners with status flags. | 1. Operations and owner tables render red "Low Balance Lock" badge.<br>2. Filter dropdown allows filtering for all locked accounts. | Status badge clearly visible; filter isolates locked accounts accurately. |
| `CW-09-BE` | Global payment configuration toggles API | P0 | Backend exposes API to update global toggles and enforces payment flow rules based on configuration. | 1. `GET/PATCH /admin/payment-settings` manages flags (`wallet_ledger_enabled`, `direct_transfer_enabled`, `online_gateways_enabled`).<br>2. Checkout middleware evaluates active flags before admitting payment initiation.<br>3. Disabled payment paths return RFC 9457 error `PAYMENT_METHOD_DISABLED`. | Global toggles update in database; checkout middleware enforces active paths; disabled methods rejected. |
| `CW-09-MOB` | Dynamic checkout views by global toggles | P0 | Mobile app fetches configuration flags and shows/hides manual or ledger options dynamically. | 1. Checkout view queries payment options on load.<br>2. Renders only payment pathways enabled by backend toggles.<br>3. If a method is toggled off globally, client omits it without requiring app update. | Checkout dynamically reflects active backend payment pathways; disabled options omitted. |
| `CW-09-WEB` | Admin settings config toggles UI | P0 | Next.js admin dashboard settings page with switches for direct payments and wallet ledger. | 1. Admin navigates to Settings $\rightarrow$ Payment Configuration.<br>2. Renders toggle switches for each payment rail and settlement pathway.<br>3. Toggling switch dispatches Server Action updating `PaymentSetting` table. | Switches reflect active database state; save updates live platform checkout behavior. |
| `CW-12-BE` | Wallet Negative Balance Suspension Cron | P1 | Auto-suspends owner profile and active listings if negative balance is unpaid for $> 90$ days. | 1. Daily cron scans `OwnerWallet` where `balance_minor < 0` and last top-up was $> 90$ days ago.<br>2. Sets owner profile status to `SUSPENDED`.<br>3. Sets all owner fleet listings to `INACTIVE`.<br>4. Logs event in `AuditLog` and dispatches suspension notification. | Delinquent wallets trigger auto-suspension; listings delisted; audit log recorded; notifications dispatched. |
| `CW-12-WEB` | Delinquent Accounts Console View | P1 | Next.js admin table to display delinquent owner accounts. | 1. Finance dashboard displays "Delinquent Accounts ($> 90$ Days)" view.<br>2. Table lists suspended owners, unpaid balance amount, delinquency duration, and contact details.<br>3. Action button: "Contact Owner" or "Send Legal Notice". | Table renders delinquent accounts; accurately sorts by days delinquent; action tools functional. |
| `CW-13-BE` | Daily Accounting Reconciliation Cron | P1 | Cron verifies matching ledgers and database balances against bank statement inputs. | 1. Schedulers run daily at 02:00 UTC.<br>2. Compares bank statement imported credits vs `OwnerWalletTransaction` top-ups vs booking commission debits.<br>3. Mismatches flagged and written to `FinancialMismatch` queue.<br>4. Generates reconciliation summary report. | Discrepancies detected and written to mismatch queue; matching records cleared; summary report logged. |
| `CW-13-WEB` | Financial Mismatch Queue Panel | P1 | Admin Next.js view displaying unmatched billing and payout errors. | 1. Admin opens Finance $\rightarrow$ Reconciliation Mismatches.<br>2. Table displays unmatched deposits, duplicate receipts, and ledger imbalances.<br>3. Action buttons: "Manual Match", "Post Adjustment", or "Flag Fraud". | Panel displays reconciliation discrepancies with diff details; adjustment modal logs to `AuditLog`. |
| `CW-14-BE` | Commission base & post-confirmation recompute | P1 | Define the canonical commission base and how post-confirmation charges trigger incremental commission debits. | 1. Commission base $=$ rental fee only (excludes VAT, excess hours, accepted damages, deposit, and mobilization fees).<br>2. Central `commissionBase(booking)` function used platform-wide.<br>3. Rental extensions (`AR-04`) trigger incremental `COMMISSION_DEBIT` on delta of base fee.<br>4. Excess hours and damages do not attract commission.<br>5. Idempotent per event. | `commissionBase()` returns rental fee only; commission recomputed incrementally on extensions; damage/excess fees excluded; base used consistently. |
| `CW-15-BE` | VAT liability & commission VAT model | P0 | Declare supplier-of-record for VAT and whether commission carries platform VAT, with owner VAT-registration handling. | 1. Owner is supplier-of-record for rental supply: 15% VAT on rental fee + attachments collected and remitted by owner; issued under owner TIN.<br>2. Platform commission is separate supply: commission invoice issued under platform TIN; treated as VAT-inclusive where platform is registered.<br>3. Owner profile stores `vat_registered` (bool) and `vat_registration_no`; if false, VAT lines suppressed on rental invoice.<br>4. `TaxExempt` renter suppresses rental VAT. | Rental VAT attributed to owner TIN; commission VAT attributed to platform TIN; non-registered owners emit no VAT lines; `TaxExempt` renter suppresses VAT. |
| `CW-16-BE` | Currency, precision & rounding standard | P1 | Single canonical money representation, rounding rule, and order of operations for all financial math. | 1. All monetary values stored as 64-bit integer minor units (Santims for ETB, Cents for KES/USD).<br>2. Rounding is round-half-up to 2 decimal places via central `roundMoney()` helper.<br>3. Order of operations: Taxable Base $\rightarrow$ 15% VAT $\rightarrow$ Prepayment $\rightarrow$ Commission.<br>4. Tier bounds in `CommissionRule` are half-open `[min, max)`. | Currency stored as integer minor units; `roundMoney()` is sole rounding path; documented order of operations enforced; tier bounds non-overlapping. |
| `CW-17-BE` | Unified financial ledger & bank-statement import | P1 | Define a single immutable financial ledger and a bank-statement import that reconciliation runs against. | 1. Every financial event writes to `financial_ledger_entries` (`booking_id`, `account_code`, `entry_type`, `amount_minor`, `currency`, `direction`).<br>2. `BankStatementImport` table ingests bank transaction feeds.<br>3. Reconciles bank imports against ledger entries.<br>4. Ledger rows are immutable (append-only) with DB trigger preventing updates/deletes.<br>5. Retained 7 years. | Unified ledger records all money movements; bank statement import supported; rows immutable; retained 7 years for tax compliance. |
| `CW-18-BE` | Booking & top-up amount limits | P2 | Enforce configurable minimum/maximum booking value and top-up amount limits. | 1. Validates booking total against `SystemConfig` bounds (`booking.min_value_minor`, `booking.max_value_minor`).<br>2. Validates wallet top-up against `wallet.min_topup_minor` and `wallet.max_topup_minor`.<br>3. Out-of-bound requests rejected with RFC 9457 validation errors. | Out-of-bound booking values and top-up amounts rejected with clear validation messages; bounds configurable via admin settings. |
| `CW-19-BE` | Owner commission-dispute path | P2 | Allow owners to dispute a specific commission debit for Admin review. | 1. Owner opens dispute against a `COMMISSION_DEBIT` row in wallet ledger.<br>2. Creates administrative review ticket with claimed reason.<br>3. Admin can uphold debit or issue compensating `ADJUSTMENT` (credit) to owner wallet.<br>4. All rulings audit-logged. | Owner can dispute specific commission debit; Admin can issue compensating adjustment credit; all actions logged to audit trail. |
| `CW-20-BE` | Post-confirmation bank-reversal / clawback | P2 | Handle a renter payment that is later reversed by the bank after `CONFIRMED`. | 1. On evidence of post-confirmation bank reversal, Admin opens clawback case.<br>2. Booking flagged as `PAYMENT_REVERSED`.<br>3. Commission debit reversed via `ADJUSTMENT` credit to owner wallet.<br>4. Both parties notified; escalates to legal/dispute queue. | Reversal recorded; commission reversed via adjustment; parties notified; audit trail preserved. |

---

### 5.10 Notifications & In-App Center

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `NF-01-BE` | Notification event publishing & FCM dispatch service | P0 | Backend publishes events to in-process event channel, parses template, and dispatches payload via FCM. | 1. Domain interactors publish events to internal Go channel.<br>2. Worker fetches user FCM device tokens.<br>3. Retrieves localized template from `NotificationTemplate`.<br>4. Dispatches push notification via Firebase Cloud Messaging.<br>5. Inserts notification row in `Notification` table. | FCM push payloads dispatched successfully; in-app notification records stored; handles device token invalidation. |
| `NF-01-MOB` | Mobile FCM integration handler | P0 | Mobile app handles Firebase Cloud Messaging tokens and listens for incoming background/foreground push payloads. | 1. Requests system notification permissions on first launch.<br>2. Registers FCM device token with backend API.<br>3. Listens for foreground, background, and terminated message events.<br>4. Renders native heads-up notification banner. | Native notification alerts appear on device; device token registers with backend; handles token refresh. |
| `NF-02-BE` | Fetch notifications & unread count APIs | P0 | API to return list of notifications and read/unread badge counts. | 1. `GET /notifications` returns paginated notification history.<br>2. `GET /notifications/unread-count` returns total unread count.<br>3. `POST /notifications/read` marks specified or all notifications as read. | API responds with accurate logs; updates unread counts in real time; validates user ownership. |
| `NF-02-MOB` | Mobile notification center UI | P0 | Mobile app displays badged bell icon in home dashboard leading to an inbox list. | 1. App bar renders bell icon with unread count badge.<br>2. Tapping opens Notification Center screen.<br>3. Lists notifications grouped by date.<br>4. Tapping item marks it as read and executes deep link to target screen. | Badge updates dynamically; list renders notification icon, title, timestamp; tap executes deep routing cleanly. |
| `NF-03-BE` | Save & enforce notification preference configurations API | P1 | Backend exposes endpoints to edit preferences and queries preferences table before dispatching messages. | 1. Classifies notifications as `TRANSACTIONAL/SECURITY` (mandatory, cannot be disabled) or `INFORMATIONAL/MARKETING` (user-toggleable).<br>2. `GET/PUT /users/preferences` manages settings.<br>3. `PUT` rejects attempts to disable `TRANSACTIONAL/SECURITY` with validation error.<br>4. Dispatcher checks preferences before sending informational alerts. | Preferences update cleanly; informational notifications skipped if disabled; attempts to disable mandatory alerts rejected. |
| `NF-03-MOB` | Mobile user notification settings UI | P1 | Mobile settings page with toggle buttons for each notification category and delivery channel. | 1. User navigates to Settings $\rightarrow$ Notifications.<br>2. Renders toggle switches for marketing, promos, and system updates.<br>3. Mandatory categories (booking updates, payment alerts, security notices) rendered as locked/disabled toggles with explanation label. | User can toggle non-mandatory categories; mandatory alerts remain locked; saves updates to backend API. |
| `NF-04-BE` | Dynamic Gateway Router Service | P0 | Resolves notification channels, credentials, and routing strategies dynamically from NotificationChannel DB. | 1. Inspects recipient phone number prefix (`+2519` Ethio Telecom, `+2517` Safaricom Ethiopia, `+254` Kenya).<br>2. Resolves active SMS gateway credentials and endpoint from `NotificationChannel`.<br>3. Routes via configured strategy: `PARALLEL` (push + SMS for critical security) or `SEQUENTIAL`. | Messages routed dynamically by prefix; supports Email (SMTP), SMS, and Push; new gateways configured without code changes. |
| `NF-05-BE` | Format deep links in notification payloads | P0 | Backend includes deep link paths in push payloads and DB records. | 1. Injects standardized `deep_link` URI (e.g. `dink://bookings/{id}`, `dink://disputes/{id}`) in JSON push payload.<br>2. Stores URI in `Notification.deep_link` database column. | Payloads contain valid deep link strings; links adhere to unified app routing scheme. |
| `NF-05-MOB` | Handle app deep routing on notification interaction | P0 | Mobile app parses deep links in push payloads on click and routes users to the target screen. | 1. User taps push notification or notification center item.<br>2. App intercepts deep link URI.<br>3. Parses route parameters and navigates directly to target screen (e.g. booking review, payment upload). | Tapping notification routes user directly to target screen without opening root dashboard. |
| `NF-06-BE` | Daily notification archival cron | P1 | Scheduled database job runs daily to migrate notification logs older than 90 days to the archive table. | 1. Scheduled cron runs daily at 03:00 UTC.<br>2. Copies rows from `Notification` to `NotificationArchive` where `created_at < NOW() - 90 days`.<br>3. Deletes archived rows from active table in batches. | Database maintenance job executes cleanly without table locks; archived rows preserved in archive table. |
| `NF-07-BE` | Dispatch webview FCM payload & session validation | P0 | Backend sends FCM push with `action_type: OPEN_WEBVIEW` and target URL, and validates incoming JWT sessions from webview requests. | 1. Publishes push payload with `action_type = OPEN_WEBVIEW` and authenticated URL.<br>2. Next.js middleware extracts and validates bearer token on incoming webview requests. | Push payload formatted correctly; webview authentication middleware admits valid sessions. |
| `NF-07-MOB` | Authenticated in-app webview integration | P0 | Mobile app launches secure webview on tap and injects authorization JWT into the requests. | 1. Detects `OPEN_WEBVIEW` payload action.<br>2. Launches in-app webview modal.<br>3. Injects active JWT bearer token into initial request headers.<br>4. Provides close and refresh controls. | Webview launches seamlessly; session authenticated without re-login prompt; header injection verified. |
| `NF-07-WEB` | Webview session authentication handler | P0 | Next.js landing pages accept and validate authentication tokens from headers/cookies to authenticate users without login prompts. | 1. Webview route intercepts `Authorization` bearer token.<br>2. Validates JWT signature.<br>3. Sets session cookie and renders target view. | Authenticates session automatically; renders responsive mobile-optimized landing page. |
| `NF-08-BE` | Retry Engine & Dead-Letter Queue (DLQ) | P1 | Enforces exponential backoff retries and registers failed dispatches in DLQ. | 1. If push/SMS gateway returns transient error (5xx, timeout), queues retry with exponential backoff (1s, 2s, 4s, 8s).<br>2. If retries exhausted or permanent error occurs, moves message payload to `NotificationDeadLetterQueue`.<br>3. Alerts on-call engineering. | Transient errors retried automatically; permanent failures isolated in DLQ with error context. |
| `NF-09-BE` | Delivery confirmation & critical-event SMS fallback | P0 | For events classified `criticality = CRITICAL`, confirms push acknowledgement and falls back to SMS if unacknowledged within 5 minutes. | 1. Dispatches push notification for critical event (payment confirmed/declined, booking confirmed, dispute opened, security alert).<br>2. Starts 5-minute acknowledgement timer.<br>3. If client sends delivery receipt, cancels timer.<br>4. If unacknowledged after 5 minutes, automatically dispatches SMS fallback via prefix-resolved gateway.<br>5. Logs delivery channel outcome. | Critical alerts reach user on push or SMS within 5 minutes in $\ge 99\%$ of cases; delivery status queryable per notification. |
| `NF-10-BE` | Notification quiet hours | P1 | Users may configure quiet hours (default 22:00–07:00 EAT); non-critical notifications generated during quiet hours are deferred. | 1. `UserNotificationPreference` stores quiet hours window (start/end).<br>2. Non-critical notifications generated inside window are queued in `NotificationDeferred`.<br>3. Deferred notifications dispatch automatically when window ends.<br>4. Critical alerts (`CRITICAL`) bypass quiet hours and send immediately. | Non-critical alerts deferred to end of quiet hours; critical security/booking alerts dispatch immediately. |
| `CW-10-BE` | Policy version tracking & middleware API enforcement | P0 | Backend stores policy documents, user acceptance logs, and rejects critical endpoints with `POLICY_NOT_ACCEPTED` if policies are unapproved. | 1. `PolicyDocument` stores versioned terms and privacy policies.<br>2. `UserPolicyAcceptance` logs user acceptance timestamps.<br>3. Middleware checks caller's accepted version against active mandatory version.<br>4. If unaccepted, blocks mutating operations with HTTP 403 `POLICY_NOT_ACCEPTED`. | Mutating endpoints blocked for accounts with unaccepted mandatory policies; returns structured error. |
| `CW-10-MOB` | Blocking policy acceptance screen UI | P0 | Mobile app checks policy compliance on launch, rendering blocking overlays to force acceptance if required. | 1. Receives `POLICY_NOT_ACCEPTED` or detects unaccepted policy flag on startup.<br>2. Renders full-screen non-dismissible modal displaying updated legal terms in webview.<br>3. "Accept & Continue" submits acceptance API call and unlocks app.<br>4. "Decline" logs user out. | Overlay cannot be dismissed without accepting; acceptance unlocks app immediately; decline terminates session. |
| `CW-10-WEB` | Admin policy management portal | P0 | Next.js admin interface to manage policy texts, set mandatory states, and view user reports. | 1. Admin navigates to Settings $\rightarrow$ Legal Policies.<br>2. Markdown editor allows updating English and Amharic policy texts.<br>3. Toggle "Mandatory Acceptance" creates new active version.<br>4. Publishes policy and tracks user acceptance statistics. | Policy editor saves versioned documents; mandatory toggle activates middleware enforcement; displays acceptance metrics. |
| `CW-11-BE` | Version verification API | P0 | Backend API endpoint checks current client app version and returns minimum and recommended version metadata. | 1. Exposes `GET /app-version/check` accepting `platform` and `client_version`.<br>2. Queries `AppVersion` table.<br>3. Returns payload with `min_version`, `latest_version`, `force_update` (bool), and store download URL. | Accurately returns version bounds; correctly computes force update flag based on minimum version threshold. |
| `CW-11-MOB` | Force/optional update check logic UI | P0 | Mobile app validates client version on startup, displaying mandatory or dismissible update dialogues. | 1. App calls version check API on startup.<br>2. If `client_version < min_version`, renders non-dismissible modal: "Update Required" with direct link to store.<br>3. If `client_version < latest_version`, renders dismissible dialog: "Update Available" with "Later" and "Update Now" options. | Mandatory update blocks application access completely; optional update permits skip; store link launches store app. |
| `CW-11-WEB` | Admin app version settings dashboard | P0 | Next.js admin page tools to change minimum and current version numbers. | 1. Admin opens Settings $\rightarrow$ Mobile App Versions.<br>2. Form allows updating `min_version` and `latest_version` per platform (Android, iOS).<br>3. Saving updates `AppVersion` database table immediately. | Admin updates version bounds; changes reflect in API response instantly without backend redeploy. |

---

### 5.11 Admin Dashboard & Dynamic UI Engine

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `AD-01-BE` | Owner approval queue - Backend | P0 | Endpoints to list pending owner registrations, retrieve uploaded documents, and approve/reject owners with reasons. | 1. `GET /api/v1/admin/owners/pending` returns list of owners in `PENDING_REVIEW`.<br>2. `POST /api/v1/admin/owners/{id}/approve` transitions status to `ACTIVE`.<br>3. `POST /api/v1/admin/owners/{id}/reject` requires text reason; sets status to `REJECTED`.<br>4. Dispatches status update notifications to owner. | Owner status updated; rejection reason logged; audit log created; notifications triggered. |
| `AD-01-WEB` | Owner approval queue - Web | P0 | Next.js admin page displaying a table of pending owner registrations and a detail view with document preview, approval action, and rejection modal. | 1. Admin opens Operations $\rightarrow$ Owner Approvals.<br>2. Table displays pending applications with submission date and TIN.<br>3. Detail view renders trade license PDF and bank coordinates.<br>4. "Approve" button activates owner; "Reject" opens modal requiring explanation. | Documents render inline; action buttons call backend APIs; table refreshes dynamically upon approval/rejection. |
| `AD-02-BE` | Fleet approval queue - Backend | P0 | Provide API endpoints to list pending fleet listings (with spec summaries and photo URLs) and approve or reject fleet submissions. | 1. `GET /api/v1/admin/fleets/pending` returns listings in `PENDING_APPROVAL`.<br>2. `POST /api/v1/admin/fleets/{id}/approve` sets status to `ACTIVE`.<br>3. `POST /api/v1/admin/fleets/{id}/reject` sets status to `REJECTED` with reason.<br>4. Dispatches push alert to owner. | Status transitions accurately; audit log recorded; notifications dispatched to owner. |
| `AD-02-WEB` | Fleet approval queue - Web | P0 | Next.js admin panel view with a list of pending fleets. Detail view showing photo gallery, category, owner, submitted date, and spec summaries. | 1. Admin opens Operations $\rightarrow$ Fleet Approvals.<br>2. Table lists submitted machines.<br>3. Clicking listing opens detail drawer with photo carousel, mechanical specs, and serial numbers.<br>4. Action buttons: "Approve Listing" and "Reject Listing" (with reason). | Images load in gallery; specs display in organized tables; action buttons execute API calls and update queue. |
| `AD-03-BE` | User management - Backend | P0 | Search and filtering APIs for users, detail query API including booking history, and account status toggle endpoints (suspend/activate). | 1. `GET /api/v1/admin/users` supports search by phone, email, name, role, and status.<br>2. `GET /api/v1/admin/users/{id}/details` returns profile, KYC files, and booking history.<br>3. `POST /api/v1/admin/users/{id}/suspend` and `/activate` update user status.<br>4. Suspension revokes active session tokens. | Search query filters accurately; user details return complete history; suspension invalidates sessions immediately. |
| `AD-03-WEB` | User management - Web | P0 | Next.js admin management page for users with search input, user list table, details panel displaying profile/bookings, and suspend/activate toggles. | 1. Admin views Users directory table.<br>2. Search bar filters users in real time.<br>3. Clicking row opens slide-out drawer with user KYC, active bookings, and wallet balance.<br>4. Toggle button allows suspending or reactivating account with confirmation modal. | Search updates dynamically; drawer loads booking history; suspend action calls API and updates status badge. |
| `AD-04-BE` | Dispute resolution center - Backend | P0 | API endpoints to list disputes, fetch dispute evidence/history/messages, and submit a final ruling that updates booking/payment status and logs an audit trail. | 1. `GET /api/v1/admin/disputes` returns active disputes filtered by status/origin.<br>2. `GET /api/v1/admin/disputes/{id}` returns evidence, inspection photos, and tri-party chat history.<br>3. `POST /api/v1/admin/disputes/{id}/rule` records binding ruling, reallocates escrow/deposit in `pgx.Tx`, and logs audit trail. | Dispute ruling resolves booking status in single transaction; balances ledger; writes immutable audit record. |
| `AD-04-WEB` | Dispute resolution center - Web | P0 | Next.js admin view displaying dispute records, interactive chat threads, uploaded files/photos from renter/owner, and resolution action controls. | 1. Admin opens Disputes dashboard.<br>2. Review page displays split view: Renter Claim vs Owner Rebuttal with side-by-side photos.<br>3. Dual chat panel allows communicating with parties.<br>4. Ruling controls: "Rule for Renter", "Rule for Owner", "Split / Custom Amount" with justification input. | Evidence renders side by side; chat threads update; ruling submission executes API and updates status. |
| `AD-05-BE` | Commission tracking - Backend | P0 | Provide aggregate API and list API for commission records categorised by status (earned, pending, paid, overdue). Implement logic to flag overdue payouts. | 1. `GET /api/v1/admin/commissions/summary` aggregates minor-unit totals by status.<br>2. `GET /api/v1/admin/commissions` lists paginated commission transactions with date/owner filters.<br>3. Daily cron checks negative wallet balances $> 30$ days and flags overdue accounts. | Aggregations match database records exactly; overdue flags update automatically based on invoice dates. |
| `AD-05-WEB` | Commission tracking - Web | P0 | Next.js admin page with summary cards displaying totals for each commission status and a table with filtering and pagination capabilities. | 1. Finance dashboard displays summary cards (Total Commission Earned, Unpaid Negative Balances, Overdue Commissions).<br>2. Filterable table lists transactions with owner name, booking reference, and amount.<br>3. Export button exports CSV report. | Summary cards display accurate minor-unit totals; table filters dynamically; export downloads CSV. |
| `AD-06-BE` | Dynamic UI configuration - Backend (`AD-06-BE`) | P0 | Maintain database schema and endpoints to store, update, and serve UI configuration metadata (sidebar, columns, actions, permissions) for admin views. | 1. `AdminPage`, `AdminTableConfig`, `AdminTableColumn`, `AdminAction`, `AdminSidebarItem` tables store UI metadata.<br>2. `GET /api/v1/admin/ui-config` evaluates caller's role permissions.<br>3. Returns customized navigation items, visible table columns, and authorized action triggers.<br>4. `PUT` endpoints allow Super Admins to adjust layout metadata dynamically. | UI metadata served dynamically; response filtered by caller permissions; layout modifications take effect without frontend redeployment. |
| `AD-06-WEB` | Dynamic UI configuration - Web | P0 | Super Admin console to edit sidebar items, configure custom table column lists, and actions via a visual layout builder. Next.js application dynamically renders menus and tables based on backend config. | 1. Main admin shell queries `/api/v1/admin/ui-config` on load.<br>2. Renders sidebar navigation and breadcrumbs dynamically.<br>3. `DynamicTable` component reads column definitions (header, accessor, type, sortable) and compiles table via `@tanstack/react-table`.<br>4. Action buttons render only if permitted in returned config. | Admin panel layout and tables render dynamically from API; Super Admin layout editor persists changes immediately. |
| `AD-07-BE` | Role & permission management - Backend | P0 | Manage Role, Permission, and RolePermission mappings in database. Provide CRUD endpoints. Check and enforce permissions in API gateway/middleware. | 1. `Role`, `Permission`, and `RolePermission` tables maintain RBAC definitions.<br>2. CRUD endpoints: `GET/POST/PUT/DELETE /api/v1/admin/roles` and `/permissions`.<br>3. Chi middleware validates token claims against required permission key (e.g. `fleets:approve`).<br>4. Default-deny: endpoints without explicit matching permission return HTTP 403. | Permissions stored in DB; changes propagate to user session context; unauthorized API requests rejected with 403 Forbidden. |
| `AD-07-WEB` | Role & permission management - Web | P0 | Grid/matrix interface in Next.js displaying roles on one axis and permissions on the other, allowing Super Admin to toggle mapping settings. | 1. Super Admin opens Settings $\rightarrow$ Roles & Permissions.<br>2. Renders interactive grid: Roles (columns) $\times$ Permission Keys (rows).<br>3. Toggling a checkbox issues immediate API update.<br>4. Warns on modifying critical administrative permissions. | Grid displays active mappings; toggling updates backend permissions; changes apply to active sessions. |
| `AD-08-BE` | Category & spec template manager - Backend | P1 | API endpoints for category/subcategory CRUD operations and managing dynamic specification field templates (type, name, options). | 1. `GET/POST/PUT/DELETE /api/v1/categories` manages two-tier hierarchy.<br>2. `/categories/{id}/spec-templates` manages `SpecFieldDefinition` rows (key, type, unit, validation rules, renter-visible flag).<br>3. Adding required spec triggers `status = INCOMPLETE_SPEC` quarantine on existing fleets. | Schema handles nested categories; field validation rules verified during fleet registration; backfill quarantine logic functions. |
| `AD-08-MOB` | Category & spec template manager - Mobile | P1 | Flutter form renderer that queries category specification templates and dynamically builds fields (dropdowns, text inputs, checkboxes) for fleet registration and search filters. | 1. Queries `/categories/{id}/spec-templates` on subcategory selection.<br>2. Form builder dynamically instantiates matching Flutter input widgets based on `field_type` (`STRING`, `INTEGER`, `DECIMAL`, `ENUM`, `BOOL`).<br>3. Applies template validation bounds. | Dynamic inputs render correctly based on data type; form validators applied matching template rules. |
| `AD-08-WEB` | Category & spec template manager - Web | P1 | Next.js admin page featuring a tree view for category management and a structured form builder to define dynamic specs per category. | 1. Admin navigates to Settings $\rightarrow$ Category & Specs.<br>2. Tree view allows adding/editing categories and subcategories.<br>3. Form builder panel allows adding specification fields, configuring data types, units, min/max ranges, and enum option lists. | Saving writes to backend; category changes reflect immediately in mobile and web forms. |
| `AD-09-BE` | Admin audit log - Backend | P0 | Capture and persist all administrative changes in an immutable AuditLog table with actor details, action types, target entities, timestamps, and before/after JSON diffs. | 1. Backend middleware intercepts administrative mutating routes (`/api/v1/admin/*`).<br>2. Captures actor user ID, role, client IP, action type, and target entity.<br>3. Computes JSON diff of modified entity state.<br>4. Inserts row into immutable `AuditLog` table.<br>5. Updates/deletes blocked via database trigger. | Database logs are read-only and immutable; records contain actor, timestamp, and JSON diffs; queryable via admin API. |
| `AD-09-WEB` | Admin audit log - Web | P0 | Next.js log viewer screen with filtering by administrator, date range, action type, and interactive row expansions to view JSON state diffs. | 1. Super Admin opens Settings $\rightarrow$ Audit Logs.<br>2. Searchable table lists all administrative actions.<br>3. Clicking row expands visual diff viewer highlighting added, modified, or deleted JSON attributes in red/green styling. | Table renders audit events; filters by actor and date work cleanly; visual diff provides clear before/after comparison. |
| `AD-10-BE` | Platform analytics dashboard - Backend | P1 | Provide aggregated endpoints or materialized views to calculate and serve platform metrics (revenue, growth, conversion, total bookings) filtered by date ranges. | 1. `GET /api/v1/admin/analytics/summary` queries materialized views with date range filters.<br>2. Aggregates booking volume, GMV, platform commission, active utilization rate, and user growth.<br>3. In-memory caching in Valkey (TTL 10 min) guarantees response time $< 300\text{ ms}$. | Aggregated data matches source ledger records; cached queries return under SLA ceiling. |
| `AD-10-WEB` | Platform analytics dashboard - Web | P1 | Next.js charts and metrics cards integrating with the analytics APIs to show line graphs, bar charts, and key performance indicators with date range controls. | 1. Admin dashboard homepage renders summary KPI cards.<br>2. Interactive charts display booking trends, commission revenue, and fleet category distribution.<br>3. Date range picker allows toggling between 7 days, 30 days, 90 days, and custom ranges. | Charts render smoothly; KPI metrics display formatted minor units; date range switching refreshes data cleanly. |
| `AD-11-BE` | Bulk owner/fleet actions - Backend | P1 | Add bulk endpoints for owner and fleet approval/rejection queues to execute updates and notifications efficiently. | 1. `POST /api/v1/admin/owners/bulk-action` and `/fleets/bulk-action` accept array of IDs and action (`APPROVE`, `REJECT`).<br>2. Executes batch status updates inside single transaction.<br>3. Publishes batch notification events. | Selected records update status atomically; partial failures handled transactional-safely; batch notifications queued. |
| `AD-11-WEB` | Bulk owner/fleet actions - Web | P1 | Next.js admin queues with checkbox selection columns and a sticky bulk action toolbar for approving or rejecting multiple items at once. | 1. Admin selects multiple rows in approval queue tables via checkboxes.<br>2. Floating action toolbar reveals "Approve Selected" and "Reject Selected" buttons with item count badge.<br>3. Confirmation modal executes batch API request. | Checkboxes function correctly; toolbar appears on selection; bulk action completes and refreshes queue. |
| `AD-12-BE` | Custom Advertisement Manager - Backend | P1 | DB schema and endpoints to create/schedule campaigns, upload ad media, serve active ads filtered by user target/region, and track views/clicks. | 1. `ad_campaigns` stores campaigns (`title`, `banner_url`, `target_role`, `target_region`, `start_date`, `end_date`, `cta_url`).<br>2. `GET /api/v1/ads` returns active campaigns matching caller's profile and region.<br>3. `POST /api/v1/ads/{id}/track` increments view/click counters concurrency-safely in Valkey. | Serves active ads matching user criteria; expired ads excluded; tracking increments metrics accurately. |
| `AD-12-MOB` | Custom Advertisement Manager - Mobile | P1 | Flutter widget that fetches targeted active banner ads and displays them, opening the landing URL upon click and logging tracking events. | 1. Homepage queries `/api/v1/ads` on load.<br>2. Renders banner carousel matching active campaigns.<br>3. Impression triggers view tracking call.<br>4. Tapping banner launches landing URL in browser or webview and fires click tracking call. | Banners render cleanly without layout shift; click opens correct URL; tracking calls fire reliably. |
| `AD-12-WEB` | Custom Advertisement Manager - Web | P1 | Next.js "Ad Manager" view containing creation form with banner asset upload, targeting checklists (role, region), scheduling calendar, and performance metrics dashboards. | 1. Admin opens Marketing $\rightarrow$ Campaign Manager.<br>2. Creation form captures campaign title, banner image upload, target personas (`renter`, `owner`), target regions, and flight dates.<br>3. Report dashboard displays total impressions, clicks, and CTR per campaign. | Form validates mandatory fields; image uploads succeed; performance dashboard renders accurate analytics. |
| `AD-13-BE` | Identity Manual-Verification Queue - Backend | P0 | Provider-agnostic admin queue for reviewing and overriding identity verifications. | 1. `GET /api/v1/admin/identity-queue` lists users in `PENDING_REVIEW`, filterable by `provider_key`.<br>2. `POST /api/v1/admin/identity-queue/{userId}/verify` records manual override of `IdentityVerification.status` to `VERIFIED` or `REJECTED`.<br>3. Requires written justification and admin MFA; logs to `AuditLog`. | Queue lists pending verifications across all providers; admin override updates record and logs audit trail. |
| `AD-13-MOB` | Identity Verification Redirect - Mobile | P0 | Flutter client enforcement that redirects users who don't satisfy the identity policy to the provider-driven verification screens (`REG-R07-MOB`). | 1. App evaluates `identity_status` and `missing_providers` post-login.<br>2. If policy unsatisfied and lock active, renders blocking modal listing outstanding providers.<br>3. Shows `PENDING_REVIEW` status badge if under manual review. | Blocking modal displayed when policy unsatisfied; lists missing providers; directs user to complete verification. |
| `AD-13-WEB` | *(superseded)* Identity policy admin | P0 | Superseded by `REG-IDP-03-WEB`. Retained for traceability. | All identity policy and enforcement controls are managed via `REG-IDP-03-WEB`. | Identity policy managed via `REG-IDP-03-WEB`; manual queue served by `AD-13-BE`. |
| `AD-14-BE` | Commission Settings Overrides - Backend | P0 | Database storage for custom commission rules and exception mappings. Payout and commission logic must resolve overrides first. | 1. `CommissionException` stores custom rates mapped to specific `owner_id` or `fleet_id`.<br>2. Commission calculator checks exception table before applying global base rate.<br>3. Modifying exceptions requires Super Admin role and logs to `AuditLog`. | Specific owner/fleet overrides applied accurately; fallback to global rates when no exception exists; audited. |
| `AD-14-WEB` | Commission Settings Overrides - Web | P0 | Next.js UI interface in settings to edit base commission rates and manage an exceptions list mapped to specific owners or fleets. | 1. Super Admin opens Settings $\rightarrow$ Commission Rules.<br>2. Displays global base rate with edit modal.<br>3. "Add Exception" modal allows searching for owner or machine and defining custom commission percentage.<br>4. Saves and applies live. | Form validates percentage bounds (0-100%); saves overrides to database; updates reflect in live commission calculations. |

---

### 5.12 Operator Telematics & GPS Tracking

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `GPS-01-BE` | Location sharing toggle - Backend | P2 | WebSocket server and endpoint to accept location updates from the operator's mobile client during an active rental. Discard pings when deactivated. | 1. WebSocket endpoint `/ws/v1/bookings/{id}/tracking` accepts JSON telemetry packets (`lat`, `lng`, `speed`, `heading`, `timestamp`).<br>2. Validates caller is assigned operator and booking is in `ACTIVE` or `IN_TRANSIT`.<br>3. Verifies operator consent in `OperatorConsentLog`.<br>4. Updates current position in `FleetLocation` and Valkey real-time cache. | WebSocket connection validates session and consent; coordinates update in real time; unauthorized or out-of-rental pings rejected. |
| `GPS-01-MOB` | Location sharing toggle - Mobile | P2 | Flutter switch in the operator interface during active rentals. When switched on, a background service sends location coordinates to backend via WebSocket every 30 seconds. | 1. Active rental dashboard displays "Share Location" switch.<br>2. Toggling on prompts background location permission.<br>3. Initializes background location service.<br>4. Dispatches GPS coordinate pings every 30 seconds via WebSocket.<br>5. Persistent notification indicates active tracking. | Pings dispatch every 30 seconds; persistent notification displayed; toggling off stops background GPS service immediately. |
| `GPS-02-BE` | Real-time location view - Backend | P2 | Provide subscription or query endpoints for active rental location data, allowing real-time tracking feeds. | 1. `GET /api/v1/bookings/{id}/location` returns latest recorded GPS point.<br>2. WebSocket `/ws/v1/bookings/{id}/tracking/watch` broadcasts incoming location pings.<br>3. Restricts access strictly to authorized renter, fleet owner, and admin. | Endpoints restrict location access to authorized booking participants; data returned matches live coordinate state. |
| `GPS-02-MOB` | Real-time location view - Mobile | P2 | Flutter map page showing a marker representing the current location of the fleet, updating dynamically when location data is received. | 1. Renter or owner taps "Track Machine" on active booking.<br>2. Map screen loads, centers on machine coordinates.<br>3. Connects to tracking WebSocket; moves machine marker smoothly as new pings arrive.<br>4. Displays speed and last-updated time. | Map displays marker at correct coordinates; updates render smoothly without reloading map; handles reconnects. |
| `GPS-03-BE` | Privacy & auto-stop - Backend | P2 | Automatically disable location processing for a booking when the booking state changes to COMPLETED. | 1. Booking transition to `COMPLETED` or `CLOSED` triggers tracking termination.<br>2. Closes open WebSocket connections for that booking.<br>3. Discards any subsequent location payloads.<br>4. Schedules raw GPS breadcrumbs for permanent pruning after 30 days (`§9.6`). | WebSocket closed on completion; subsequent coordinates rejected; data pruning scheduled. |
| `GPS-03-MOB` | Privacy & auto-stop - Mobile | P2 | Stop background tracking services immediately when the operator manually disables the toggle or when the local app detects the booking status transitions to COMPLETED. | 1. Listens for booking status updates.<br>2. If status transitions to `RETURN_INITIATED`, `COMPLETED`, or `CLOSED`, terminates background tracking service immediately.<br>3. Removes persistent notification bar.<br>4. Sets UI toggle to off. | Background GPS querying terminates immediately; notification bar dismissed; privacy guarantees verified. |
| `GPS-04-BE` | Location data privacy consent - Backend | P0 | Store operator location tracking consent records in OperatorConsentLog database. Verify consent exists before accepting GPS coordinates. | 1. `POST /api/v1/operator/consent` records operator agreement (`booking_id`, `operator_id`, `version`, `consented_at`).<br>2. WebSocket tracking handler asserts consent record exists.<br>3. Rejects coordinate ingestion if consent is missing. | Consent record persisted immutably; WebSocket rejects tracking without active consent record. |
| `GPS-04-MOB` | Location data privacy consent - Mobile | P0 | Flutter modal shown to the operator when first assigned to an active booking, explaining location tracking data usage and prompting for explicit consent before activating tracking. | 1. Operator opens active job board.<br>2. Modal explains tracking purpose: "Location is shared with renter and owner only during active dispatch and rental hours."<br>3. Action buttons: "Agree & Enable Tracking" or "Decline".<br>4. Submits consent record to backend API. | Modal blocks tracking activation until explicit consent is granted; declining keeps tracking disabled. |

---

### 5.13 Insurance Document Management

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `IN-01-BE` | Insurance upload & data capture - Backend | P0 | Database support for InsuranceDocument. Validation on input fields: provider name, policy number, coverage type, and expiry date. Store file link. | 1. `POST /api/v1/fleets/{id}/insurance` handles multipart certificate upload and metadata.<br>2. Streams file through ClamAV antivirus scan.<br>3. Validates that `expiry_date` is in the future.<br>4. Saves file to SeaweedFS S3 and records metadata in `InsuranceDocument` table. | Document file saved securely; metadata fields populated; past expiry dates or infected files rejected. |
| `IN-01-MOB` | Insurance upload & data capture - Mobile | P0 | Flutter screen for owners to input insurance document details and upload a photo/file of the certificate. | 1. Fleet registration or settings form displays "Insurance Details" section.<br>2. Fields: Insurance Provider, Policy Number, Coverage Type (`Third-Party`, `Comprehensive`), Expiry Date.<br>3. File picker attaches PDF or photo certificate.<br>4. Submits payload to API. | Form validates required fields; file size restricted to 10MB; uploads certificate and displays confirmation. |
| `IN-02-BE` | Expiry alert - Backend | P1 | Background cron job to check insurance document expiry dates daily. Send push notifications at 30-day threshold. At expiry, transition fleet status to INACTIVE and notify owner. | 1. Nightly cron job evaluates all active `InsuranceDocument` records.<br>2. Dispatches warning push notification to owner 30 days prior to policy expiry.<br>3. If policy expires without renewal, sets fleet status to `INACTIVE` (`EXPIRED_INSURANCE`).<br>4. Hides fleet from search. | Warning alerts sent 30 days prior; expired policies set fleet `INACTIVE`; fleet hidden from public search. |
| `IN-02-MOB` | Expiry alert - Mobile | P1 | Display alert banners and receive push notifications on owner app when insurance is nearing expiry or has expired. | 1. Intercepts insurance expiry push notification.<br>2. Renders persistent amber warning banner on owner dashboard: "Insurance for [Fleet] expires in X days."<br>3. Red banner if expired: "Fleet deactivated due to expired insurance." | Warning banners display remaining days accurately; tapping banner directs owner to renewal upload screen. |
| `IN-03-BE` | Insurance visibility to renter - Backend | P0 | Include public insurance info in the fleet detail payload returned to renters. | 1. `GET /api/v1/fleets/{id}` includes `insurance_summary` object.<br>2. Fields: `provider_name`, `coverage_type`, `valid_until`, `is_verified` (bool).<br>3. Omits private certificate document URLs and owner policy costs from renter response. | Renter detail response includes public coverage summary; private documents and costs strictly omitted. |
| `IN-03-MOB` | Insurance visibility to renter - Mobile | P0 | Display insurance details on the fleet details page in the renter app. | 1. Renter opens fleet detail view.<br>2. Overview section renders green verified badge: "Insured: Comprehensive Coverage by [Provider] (Valid through [Date])".<br>3. If unverified, displays warning tag. | Insurance summary rendered legibly; displays coverage type and expiration date clearly to renter. |
| `IN-04-BE` | Multiple insurance documents per fleet - Backend | P1 | Support one-to-many relationship between Fleet and InsuranceDocument. Adjust APIs and cron jobs to support multiple documents. | 1. `InsuranceDocument` maintains foreign key `fleet_id`.<br>2. `GET /api/v1/fleets/{id}/insurance` returns array of all active documents.<br>3. Expiry cron evaluates earliest expiring document to govern fleet active status. | Schema supports multiple policies per asset; cron evaluates earliest expiration; all linked policies queryable. |
| `IN-04-MOB` | Multiple insurance documents per fleet - Mobile | P1 | Update fleet registration and management flows in owner app to allow uploading multiple insurance documents. | 1. Fleet insurance section displays "Add Another Insurance Policy" action.<br>2. Lists currently uploaded certificates with provider and expiry chips.<br>3. Allows deleting or replacing individual documents. | User can add multiple policies; form validates each entry separately; list renders active documents cleanly. |
| `IN-05-BE` | Insurance re-upload and fleet reactivation - Backend | P0 | Backend endpoints to support insurance re-upload on expired/inactive fleets, notify admin, and allow admin approval to restore fleet status to ACTIVE. | 1. `POST /api/v1/fleets/{id}/insurance/reupload` saves renewed certificate; sets review flag.<br>2. Dispatches verification task to admin queue.<br>3. Admin calls `POST /api/v1/admin/fleets/{id}/reactivate`.<br>4. Validates new policy; restores fleet status to `ACTIVE`. | Renewal document persisted; admin approval endpoint restores fleet to `ACTIVE`; owner notified of reactivation. |
| `IN-05-MOB` | Insurance re-upload and fleet reactivation - Mobile | P0 | Owner app shows a prominent reactivation warning card if fleet is inactive due to expired insurance, leading to document upload screen. | 1. Owner dashboard displays red reactivation card for fleets in `EXPIRED_INSURANCE`.<br>2. "Update Insurance Certificate" action opens upload sheet.<br>3. Submitting shows "Pending Admin Review" badge. | Warning card visible only when expired; upload button starts renewal submission; status updates cleanly. |
| `IN-05-WEB` | Insurance re-upload and fleet reactivation - Web | P0 | Next.js admin page displaying a quick-review queue for fleets requiring reactivation due to new insurance uploads. | 1. Admin navigates to Operations $\rightarrow$ Insurance Renewals.<br>2. Table displays fleets awaiting reactivation with old vs new policy side-by-side comparison.<br>3. "Approve Reactivation" button executes backend reactivation API. | Quick-review queue shows relevant entries; side-by-side comparison displays cleanly; approval updates status live. |
| `IN-07-BE` | Insurance validity check & policy binding | P1 | At booking confirmation, assert active insurance coverage and snapshot the covering policy onto the booking. | 1. `confirm-payment` handler checks that fleet has $\ge 1$ non-expired, verified `InsuranceDocument`.<br>2. If no valid policy exists, blocks confirmation with RFC 9457 error `FLEET_INSURANCE_INVALID`.<br>3. Snapshots policy ID, provider, and coverage type onto `Booking` record. | Confirmation blocked if insurance expired; covering policy metadata snapshotted onto booking immutably. |
---
> **Scope note (SC2):** The full insurance **claims** lifecycle below (`IN-06` / insurer portal) is **DEFERRED post-MVP (P2)**. MVP insurance scope is limited to insurance **document management** only (upload, expiry checking, renter visibility — `IN-01`…`IN-05`, `IN-07`). See §13.2 Out of Scope (MVP).

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `IN-06-BE` | Insurance Claim Lifecycle API | P2 (deferred post-MVP) | Backend support and status endpoints for fleet damage claims. | 1. Exposes endpoints to create claims under active booking.<br>2. Transitions through claim statuses: `CLAIM_DRAFT`, `CLAIM_SUBMITTED`, `CLAIM_REVIEW`, `CLAIM_APPROVED`, `CLAIM_REJECTED`, `CLAIM_PAID`.<br>3. Handles insurer portal updates and settlement webhook ingestion. | Database transitions claim status transactional-safely; validates booking details and attached evidence; supports external insurer status sync. |
| `IN-06-MOB` | Mobile Damage Claim Request UI | P2 (deferred post-MVP) | Flutter forms for renters and owners to submit claims. | 1. Completed booking detail screen displays "File Insurance Claim" action.<br>2. Renders document upload form (police report, incident photos, damage estimate).<br>3. Submits claim payload to backend. | Input forms capture reports and file attachments; displays active claim status and progress tracker to user. |
| `IN-06-WEB` | Web Admin & Insurer Claim Manager | P2 (deferred post-MVP) | Next.js interface for support and insurance adjusters. | 1. Portal displays active claims lists with claim ID, booking reference, and claimed amount.<br>2. Insurance adjusters view proof files, input claim resolutions, and edit statuses.<br>3. Logs all decisions to audit trail. | Adjusters can approve, reject, or mark claims as paid; audit log records adjuster identity and resolution justification. |

---

### 5.14 User Onboarding & Help

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `OB-01-BE` | Renter onboarding walkthrough - Backend | P0 | Store the onboarding completion flag on the user profile to prevent showing the onboarding flow repeatedly. | 1. Profile database schema includes `renter_onboarded` boolean.<br>2. `POST /api/v1/user/onboarding/complete` sets value to true.<br>3. `GET /api/v1/user/profile` returns this flag. | Database persists flag; endpoint correctly changes value; subsequent requests return updated flag. |
| `OB-01-MOB` | Renter onboarding walkthrough - Mobile | P0 | Display an overlay onboarding walkthrough with 4 slides (browsing, booking, payment, inspection) upon first login, with a skip button. | 1. On login or app start, read `renter_onboarded` flag.<br>2. If false, show PageView tutorial.<br>3. Skip or complete button calls `/api/v1/user/onboarding/complete` and stores state locally. | Tutorial appears on first launch; skip button functions; tutorial is never shown again after completion. |
| `OB-02-BE` | Owner onboarding walkthrough - Backend | P0 | Store the owner onboarding completion flag on the user profile or owner profile. | 1. Database contains `owner_onboarded` boolean field.<br>2. Expose endpoint `POST /api/v1/owner/onboarding/complete` to set this flag. | Database records state; API responds with updated status. |
| `OB-02-MOB` | Owner onboarding walkthrough - Mobile | P0 | Display a 5-step overlay tutorial explaining owner tools when the owner first accesses the owner dashboard after approval. | 1. App checks `owner_onboarded` flag on owner home navigation.<br>2. If false, presents an overlay walkthrough covering owner screens.<br>3. On skip/finish, call API and set flag. | Walkthrough starts automatically on first visit; handles swipe gestures; saves state on close. |
| `OB-03-BE` | In-app FAQ & Help Center - Backend | P1 | Dynamic database tables for FAQs with language support (Amharic/English) and search APIs. Admin CRUD endpoints. | 1. Create FAQ schema.<br>2. Expose public search API `GET /api/v1/faqs?query=text&lang=am`.<br>3. Expose admin CRUD endpoints `GET/POST/PUT/DELETE /api/v1/admin/faqs`. | Search query handles bilingual inputs and performs basic keyword matching; database stores localized text correctly. |
| `OB-03-MOB` | In-app FAQ & Help Center - Mobile | P1 | Flutter screen showing categorized FAQ lists, search bar, and language support matching the current app locale. | 1. FAQ screen fetches categories from API.<br>2. Local search filters query backend.<br>3. Content displays in active locale. | Search input triggers API fetch; details expand smoothly on click; locale changes translate FAQ items immediately. |
| `OB-03-WEB` | In-app FAQ & Help Center - Web | P1 | Next.js admin page to manage FAQ database with forms for inputting English and Amharic versions of questions and answers. | 1. Admin views table of FAQ items.<br>2. Clicks "Add FAQ", enters questions/answers in both languages, assigns category.<br>3. Save submits to backend. | FAQ list renders in table; forms validate presence of both translations; save pushes changes to production database. |
| `OB-04-MOB` | Contextual help tooltips - Mobile | P1 | Flutter widget that renders a "?" icon next to complex form fields, showing a localized bottom sheet or tooltip on tap. | 1. Bind tooltip icon to localized text keys (e.g. `help_commission_info`).<br>2. Tapping triggers modal bottom sheet or tooltip bubble with explanation. | Tooltips are aligned with fields; text changes according to selected language; tap is responsive. |
| `OB-04-WEB` | Contextual help tooltips - Web | P1 | Next.js component rendering a hoverable/clickable "?" icon next to complex admin config fields, displaying helpful popover explanations in English or Amharic. | 1. Implement popover components using Tailwind or Radix UI next to inputs.<br>2. Pull strings from UI translation catalogs. | Tooltips render on hover or click; support bilingual translation keys; do not obstruct input fields. |
| `OB-05-BE` | In-app support chat - Backend | P2 | Securely provision access tokens or identity hashes for third-party support SDK integrations. | 1. Endpoint `GET /api/v1/support/token` generates JWT or HMAC-SHA256 signature using user profile info to authenticate the chat session. | Tokens generated correctly using API keys; responses are served only to authenticated sessions. |
| `OB-05-MOB` | In-app support chat - Mobile | P2 | Integrate support chat SDK or Webview component in Flutter to allow in-app chatting. | 1. Help screen includes "Chat with Support" button.<br>2. Launches the chat interface using the token received from the backend. | Chat interface loads cleanly; messages send and receive; keyboard is handled correctly on screen. |
| `OB-05-WEB` | In-app support chat - Web | P2 | Next.js integration of support chat widget in the admin portal to allow admin staff to manage support tickets or converse with users. | 1. Load third-party chat dashboard or embed support widget inside the admin workspace layout using API tokens. | Agent UI displays incoming conversations from mobile users; notifications function when new messages arrive. |

---

### 5.15 Search & Discovery Enhancements

*(Note: See functional requirements SR-08 through SR-15 in §5.4 for detailed itemization of Wishlists, Local View Caching, PostGIS Spatial Indexing, Amharic Multilingual Search, and Weighted Organic Relevance Ranking).*

---

### 5.16 Accessibility & Internationalisation

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `I18N-01-BE` | Dynamic Language Registry & OTA Translation Bundles | P0 | Server-driven localization: languages and all UI strings live in the database and are delivered as platform-filtered, versioned bundles — adding a language or fixing a typo reaches clients without app-store release or redeploy. | 1. New `Language` registry table: `lang_code` (BCP 47), `display_name`, `native_name`, `is_enabled`, `is_rtl`, `font_family`, `font_url_mobile`, `font_url_web`, `fallback_lang`, `completeness_pct`, `min_client_version`.<br>2. Every UI string is a dot-namespaced `TranslationKey` with `platform ∈ {mobile, admin, shared}`; values in `Translation` are ICU MessageFormat strings.<br>3. `GET /api/v1/translations/{lang}` with `X-Platform` and `X-App-Version` returns bundle pre-merged recursively over `fallback_lang` chain; supports ETag/If-None-Match (304).<br>4. Lightweight `GET /api/v1/translations/{lang}/version` for 15-min client poll.<br>5. Keys are never deleted (`is_deprecated` flag hides from editor but serves older clients). | New language enabled via admin serves complete fallback-merged, platform-filtered bundle with zero client release; mobile bundles contain no admin-only keys; 304 on unchanged ETag; ICU plural strings round-trip intact. |
| `I18N-02-BE` | Translation Management & Publish Workflow APIs | P0 | Admin CRUD + draft $\rightarrow$ publish workflow for languages/keys/values, with bulk import/export for external translators and safety invariants. | 1. Admin endpoints: language CRUD, key CRUD, per-language value editing.<br>2. Bulk export (`format=csv\|json`) and import with dry-run validation (missing/extra/malformed keys, ICU syntax check).<br>3. Edits accumulate as drafts; explicit Publish atomically bumps `bundle_version`, records change to `AuditLog`, and requires admin MFA + `config:manage`.<br>4. Invariant: the last enabled language cannot be disabled (API rejects).<br>5. Preview mode behind feature flag for test cohorts. | Translator round-trip needs no code; publish is atomic, MFA-gated, audit-logged; disabling the last enabled language is rejected; ICU syntax errors caught at import; preview works on user cohort. |
| `I18N-02-WEB` | Translation & Language Admin Console | P0 | Next.js admin screens for I18N-02-BE, built as dynamic-UI components per ADR-008. | 1. Language list: enable/disable, RTL, fonts, fallback, completeness bar.<br>2. Translation editor: searchable key grid, side-by-side source (en) + target inline edit with ICU placeholder validation.<br>3. Import/export with dry-run results; per-platform Publish button with pending-draft diff count. | Admin can add Oromifa, upload fonts, import translated CSV, review completeness, and publish per platform — all UI-driven; console renders via dynamic-UI engine. |
| `I18N-03-BE` | Smart User Language Fallback | P1 | Per-user language priority list with silent server-side fallback when an admin disables a language mid-flight. | 1. `users.language_priority` (JSON array, e.g. `["am","om","en"]`); `preferred_lang` = head.<br>2. Request middleware: if `preferred_lang` is disabled, resolve first enabled language in priority list, silently update `preferred_lang`, and set `X-Language-Fallback: true` header.<br>3. Client surfaces toast: "Language unavailable. Switched to fallback."<br>4. Notification rendering uses same resolution. | Disabling a language never strands a user: next request silently re-resolves and flags switch; toast shows once; notifications follow resolved language. |
| `A11Y-01-BE` | Language & Locale Preference - Backend | P0 | Store and apply user's preferred language (+ separate formatting locale) for localized notifications, SMS, and emails, validated against Language registry. | 1. `users.preferred_lang` validated against enabled `Language` rows; `users.preferred_locale` separate from language.<br>2. `PUT /api/v1/user/settings` updates both; profile returns them.<br>3. Outbound SMS/push/email templates resolve through TranslationKey catalog + fallback chain. | Both fields persist and round-trip; disabled/unknown codes rejected; notifications render in resolved language; locale defaults correctly when unset. |
| `A11Y-01-MOB` | Dynamic Language Runtime - Mobile | P0 | Flutter i18n runtime on OTA bundles with atomic language activation: language appears in picker only when bundle AND font are fully loaded. | 1. Picker renders from `GET /languages` filtered to ready languages; English pre-bundled as build-time asset.<br>2. Poll version endpoint every 15 min; on change, fetch bundle and font if `font_url_mobile` set.<br>3. Activation gate: bundle parsed OK + font downloaded + `FontLoader` loaded OK before selectable.<br>4. ICU MessageFormat resolved at render time; `is_rtl` wraps app in `Directionality`. | Server-enabled language appears in picker only after bundle+font verified locally; switch is instant, persists, and never restarts app; offline first-run renders English snapshot; RTL mirrors correctly. |
| `A11Y-01-WEB` | Dynamic Language Runtime - Web | P0 | Next.js admin runtime on same OTA bundles (admin-platform-filtered) with atomic activation via `FontFace` / Cache API. | 1. Selector from `GET /languages`; catalogs via bundle endpoint (`X-Platform: admin`) with ETag revalidation.<br>2. Fonts loaded from `font_url_web` (WOFF2) through `FontFace` API and cached.<br>3. `dir="rtl"` on `<html>` from `is_rtl`; styles use CSS logical properties (`ms-*`/`me-*`). | New languages appear without redeploy; admin-only keys never reach mobile; RTL flips layout correctly; font swap-in has no FOUT longer than one frame. |
| `A11Y-02-MOB` | OTA Typography - Mobile | P0 | Hybrid font strategy: Latin (`en`) font pre-bundled; non-Latin scripts (Ethiopic, Arabic) delivered OTA via `font_url_mobile` and hot-loaded. | 1. Ships only Latin font in assets.<br>2. Amharic/Arabic fonts download per `A11Y-01-MOB` activation gate.<br>3. Theme `fontFamily` switches from `Language.font_family` metadata.<br>4. Updating font = new `font_url` + publish (glyph fixes reach users OTA). | Ethiopic/Arabic glyphs render correctly after OTA load; base APK/IPA remains lean; font updates propagate without app store release; line-height/clipping verified per script. |
| `A11Y-02-WEB` | OTA Typography - Web | P0 | Same hybrid strategy for admin panel using WOFF2 + `FontFace`, driven by `font_url_web`. | 1. Latin font self-hosted in app.<br>2. Non-Latin scripts loaded via `FontFace` from `font_url_web` with Cache API persistence and `font-display: swap` semantics. | Ethiopic characters render cleanly; font updates propagate on next revalidation; no layout breakage on script switch. |
| `A11Y-03-MOB` | Minimum tap target size - Mobile | P1 | Ensure all interactive components in Flutter app meet or exceed $48 \times 48\text{ dp}$ tap target sizing. | 1. Set minimum height/width or padding on custom buttons, icon buttons, and list items.<br>2. Run automated widget testing verifying layout boundaries (`androidTapTargetGuideline`). | Mobile widgets have bounds $\ge 48 \times 48\text{ dp}$; interactive icons wrapped in splash/gesture detectors meeting sizing standards. |
| `A11Y-03-WEB` | Minimum tap target size - Web | P1 | Ensure all interactive buttons, links, and forms in Next.js portal have a minimum tap/click area of $44 \times 44\text{ px}$. | 1. Apply Tailwind utility classes to ensure buttons, checkboxes, and navigational links have adequate padding/sizes. | Click targets meet minimum accessibility requirements; elements easily selectable on touchscreen and mouse inputs. |
| `A11Y-04-MOB` | Screen reader compatibility - Mobile | P1 | Implement Flutter Semantics widgets to provide descriptive text for screen readers (TalkBack/VoiceOver). | 1. Wrap decorative images, non-text buttons, and list components with `Semantics` nodes.<br>2. Supply concise description strings. | TalkBack/VoiceOver reads clear text for interactive components; decorative items ignored by reader focus. |
| `A11Y-04-WEB` | Screen reader compatibility - Web | P1 | Next.js code must use semantic HTML tags, ARIA attributes, and image alt text for screen reader compliance (NVDA, VoiceOver). | 1. Use `<button>`, `<nav>`, `<main>` structural tags.<br>2. Add `aria-label` to icon buttons and `alt` tags to images. | Accessibility audit tools (Lighthouse/axe) pass with high score; screen reader reads correct element context. |

---

### 5.17 Offline & Low-Connectivity Support

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `OFF-01-MOB` | Offline detection & banner - Mobile | P0 | Implement network connectivity listener in Flutter using connectivity_plus, displaying a persistent banner when network is lost. | 1. Listen to connection state changes.<br>2. Toggle global state that renders an amber notification banner.<br>3. Dismiss banner and trigger sync queue processing when reconnect is detected. | Banner triggers immediately on network loss; app functions without crashing; banner hides automatically upon reconnection. |
| `OFF-02-MOB` | Cached fleet browsing - Mobile | P1 | Implement local caching for recently viewed fleet details using Isar database, allowing offline retrieval for up to 24 hours. | 1. On loading fleet details, write payload JSON and timestamp to local Isar store.<br>2. If fetching fails due to offline state, read from Isar and display with "Last updated X ago" warning banner.<br>3. Expire cache older than 24 hours. | Cached data displays when offline; banner displays correct last-updated duration; cached data older than 24 hours purged. |
| `OFF-03-MOB` | Photo upload retry - Mobile | P0 | Create a persistent background upload queue in Isar to hold failed image uploads and trigger uploads sequentially on reconnection. | 1. Maintain `OfflineUploadQueue` store containing local file path, target API endpoint, and booking ID.<br>2. Register network state listener to start uploading when online.<br>3. Render persistent task progress bar in app dashboard. | Failed uploads do not lose files; files upload in background automatically once online; upload progress bar updates correctly. |
| `OFF-04-MOB` | Booking action offline guard - Mobile | P0 | Intercept mutating booking actions (request, cancel, pay) when offline and block them, displaying a descriptive modal. | 1. Pre-action helper checks current network status.<br>2. If disconnected, show dialog stating "You need an active internet connection to complete this transaction" instead of dispatching API call. | Mutating actions blocked instantly when offline; user receives modal explanation; API call is not dispatched. |

---

### 5.18 Error Handling & User Feedback

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `ERR-01-BE` | API error taxonomy - Backend | P0 | Build global error handler returning structured JSON error payloads conforming to RFC 9457 Problem Details. Avoid leaking internal traces. | 1. Global exception filter maps runtime and domain errors to predefined error codes.<br>2. Formats response envelope: `{ "success": false, "error": { "type": "...", "title": "...", "status": 400, "detail": "...", "code": "...", "invalidParams": [] }, "meta": { ... } }`. | API returns uniform RFC 9457 JSON formats for 4xx and 5xx statuses; server logs retain stack traces; response payload sanitizes internal details. |
| `ERR-01-MOB` | API error taxonomy - Mobile | P0 | Implement mapping handler in Flutter app to convert backend RFC 9457 error codes into user-friendly localized messages. | 1. Setup error-handling middleware in Dio client.<br>2. Match `error.code` against translation keys.<br>3. Fallback to generic localized message for unhandled errors. | Error dialogs display relevant translated messages; raw code errors logged internally but hidden from user. |
| `ERR-01-WEB` | API error taxonomy - Web | P0 | Implement error response parser in Next.js API client to display user-friendly error banners and inline field validations. | 1. Catch API error responses in Server Actions.<br>2. Map known error codes to localized error messages.<br>3. Bind `invalidParams` array directly into React Hook Form field validation errors. | Admin receives clean error alerts; invalid forms highlight fields from API validation details automatically. |
| `ERR-02-MOB` | Form validation feedback - Mobile | P0 | Flutter form fields validate inputs on focus loss (blur) and on form submission, showing localized red error text. | 1. Enable dynamic form validation keys.<br>2. On field focus loss, run validation function.<br>3. Submit button remains locked or triggers validations. | Incorrect fields display red error labels below input; form submission blocked until validation succeeds; error text respects locale. |
| `ERR-02-WEB` | Form validation feedback - Web | P0 | Next.js React forms validate inputs on blur and submit, rendering localized warning text under fields. | 1. Use React Hook Form + Zod resolvers to validate values.<br>2. Render error messages below input fields in red style. | Form prevents API calls if client-side validation fails; error messages display in active portal language. |
| `ERR-03-MOB` | Loading states - Mobile | P0 | Implement Flutter loading animations and progress indicators for network requests. | 1. Wrap buttons in state checking to replace text with spin indicator on press.<br>2. For requests lasting longer than 3 seconds, display a progress explanation sheet. | Action buttons disabled during network calls; progress bars update for file transfers; full-screen spinner covers page transition. |
| `ERR-03-WEB` | Loading states - Web | P0 | Next.js pages and forms display loading spinners, disable inputs during requests, and show progress bars for file uploads. | 1. Bind form states to loading flags.<br>2. Show skeleton loaders or spinner icons on buttons while async actions resolve. | Double clicking form buttons blocked; skeleton placeholders prevent layout shift during page fetch. |
| `ERR-04-MOB` | Empty states - Mobile | P1 | Display user-friendly placeholder screens when lists (fleet search, booking history) are empty in Flutter. | 1. If list length is 0, render widget with localized empty message, vector graphic illustration, and action button. | Empty screen appears when search matches nothing; action button performs correct recovery step. |
| `ERR-04-WEB` | Empty states - Web | P1 | Render responsive empty states in Next.js dashboards for tables, filters, or inbox panels when no data is returned. | 1. Create custom layout component showing illustration, description, and primary button.<br>2. Render when dataset is empty. | Tables show custom empty state component rather than blank space; clear-all filter button resets view. |
| `ERR-05-MOB` | Crash reporting - Mobile | P0 | Integrate Firebase Crashlytics SDK in Flutter to automatically log and report unhandled client exceptions. | 1. Add Crashlytics dependencies.<br>2. Initialize in main block to catch Flutter framework errors and native crashes.<br>3. Log breadcrumbs and scrub PII before sending. | Uncaught errors successfully transmitted to Crashlytics dashboard; stack trace and device specs included; zero PII leaked. |

---

### 5.19 Account Deletion & Data Subject Rights

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `DEL-01-BE` | Account Anonymisation API | P1 | Backend endpoint to delete/anonymise user accounts while retaining audit records. | 1. Receives confirmed delete request at `POST /user/delete-account`.<br>2. Asserts zero active/incomplete bookings exist.<br>3. Sets account status to `PENDING_DELETION` with a 30-day grace period (`account.deletion_grace_days`).<br>4. If user logs in during grace period and selects "Cancel Deletion", account reverts to active.<br>5. Once grace period elapses, overwrites PII (name, phone, email, TIN, bank details) with irreversible tokens (`deleted-user-{hash}`), nulls free-text fields, and severs person linkability.<br>6. Retains booking and ledger rows in anonymised-linkage form for 7-year audit retention. | PII overwritten/anonymised only after grace window; post-anonymisation records cannot be re-identified; login disabled during `PENDING_DELETION`; account reactivates if cancelled in grace period; financial audit records preserved. |
| `DEL-01-MOB` | Settings Account Deletion UI | P1 | Mobile user settings interface for requesting account deletion. | 1. User opens Settings $\rightarrow$ Account.<br>2. Taps "Delete My Account".<br>3. Displays warnings regarding 30-day grace period and data loss.<br>4. Prompts password and OTP verification.<br>5. Dispatches delete API and logs user out on success. | Enforces confirmation steps; alerts user to grace period; triggers backend deletion call. |
| `DEL-02-BE` | Async Data Portability Export API | P1 | Backend asynchronous job compiling user data into machine-readable format. | 1. Receives request at `POST /user/export-data`.<br>2. Enforces rate limit (max 1 request per 30 days).<br>3. Spawns background worker compiling Profile, Bookings, Wallets, Reviews, and Disputes to JSON/CSV.<br>4. Compresses into encrypted ZIP file.<br>5. Uploads to SeaweedFS S3.<br>6. Sends email with secure 7-day expiring download URL. | Throttling restricts to 1 request per 30 days; data compressed and stored securely; download link sent via email; link expires in 7 days. |
| `DEL-02-MOB` | Request My Data Settings UI | P1 | Mobile settings button and status indicator for data portability exports. | 1. User navigates to Settings $\rightarrow$ Privacy.<br>2. Taps "Export My Data".<br>3. Displays confirmation dialog: "This process compiles all account records. Download link will be emailed within 24 hours."<br>4. Dispatches export request. | Button disabled for 30 days after request; alerts user of email delivery; displays confirmation. |

---

### 5.20 Policy & Terms Management

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `POL-01-BE` | Policy Enforcement Middleware | P2 | API middleware checking if user accepted latest mandatory terms. (Governed canonically by `CW-10-BE`). | 1. Validates incoming requests.<br>2. Checks user's last accepted terms version against active mandatory version in `PolicyDocument`.<br>3. If unaccepted, blocks mutating requests with RFC 9457 `POLICY_NOT_ACCEPTED`. | API blocks mutating calls for users who have not accepted latest mandatory policy. |
| `POL-01-MOB` | Blocking Policy Accept Screen | P2 | Mobile app overlay display forcing user terms acceptance. | 1. App receives `POLICY_NOT_ACCEPTED` from API.<br>2. Renders full-screen sheet displaying policy text in webview.<br>3. User must click "Accept" to dismiss, or "Decline" which logs them out. | Sheet cannot be swiped away or skipped; accepting submits accept API and unlocks app. |
| `POL-01-WEB` | Admin Policy Editor Portal | P2 | Next.js admin dashboard to draft, publish, and enforce policy terms. | 1. Admin navigates to Policy Settings.<br>2. Creates new document with text (English/Amharic) and version number.<br>3. Toggles "Mandatory".<br>4. Publishes, immediately updating active terms. | Admin panel saves and publishes versioned policies; updates database. |

---

### 5.21 Feature Flagging & Remote Configuration

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `FF-01-BE` | Rule Evaluation & Sync API | P0 | Fast, caching API endpoint that dynamically resolves and returns active feature flag values for the user. | 1. Receives request at `GET /api/v1/flags/sync` with headers (`App-Version`, `Platform`, `User-Id`, `Role`).<br>2. Queries active flags from Valkey in-memory cache.<br>3. Evaluates rules sequentially: version range, platform match, role match, deterministic hashing for rollout percentage.<br>4. Returns JSON key-value map of resolved flag values. | Returns resolved flag values; sub-millisecond response time from cache; deterministic percentage rollouts based on hashed user ID. |
| `FF-01-MOB` | Mobile Flag Sync & Offline Cache | P0 | Mobile app synchronization and in-memory caching of feature flags on launch and app resume. | 1. Initiates sync at startup and foreground resume.<br>2. Calls `/api/v1/flags/sync` and stores response in local cache.<br>3. If offline, loads cached values; otherwise falls back to hardcoded defaults.<br>4. Code evaluates flags via `FlagManager.get(key, defaultValue)`. | Syncs flags on launch; falls back to cached values or defaults if offline; does not block UI thread. |
| `FF-01-WEB` | Web Client Flag Sync | P0 | Next.js client synchronisation of feature flags on session load. | 1. Next.js app calls `/api/v1/flags/sync` during hydration or middleware.<br>2. Caches flag configuration in context.<br>3. Custom React hook `useFeatureFlag(key)` renders or hides components dynamically. | Caches flags on page load; supports SSR fallback; hides/shows components dynamically based on resolved flag. |
| `FF-02-BE` | Admin Console Rule Builder & Hygiene Gates | P1 | Backend APIs to manage feature flags, rules, and priority matching, enforcing expiration dates and cleanup tickets. | 1. Admin creates flag, making `expires_at` and `cleanup_ticket` mandatory fields.<br>2. Enforces standard namespace validation (e.g. `corporate.registration.enabled`).<br>3. Daily cron alerts admins of expired flags and lists active tickets to prevent code debt. | Backend validates namespaces, expiration constraints, and cleanup tickets; cron triggers alert logs. |
| `FF-02-WEB` | Admin Feature Flag Management UI | P1 | Next.js dashboard interface to manage flags, configure rule conditions, and track flag debt. | 1. Admin views paginated list of all feature flags with keys, types, expiration dates, and cleanup tickets.<br>2. Highlights expired flags with yellow warning badge.<br>3. Form builder configures rule conditions (operators: `EQUALS`, `IN`, `VERSION_GE`).<br>4. Drag-and-drop rule priority reordering. | Provides intuitive rule editor; reordering triggers immediate API updates; expiration warnings and ticket links clearly visible. |

---

### 5.22 Global Soft-Delete Policy (`SD-01-BE`)

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `SD-01-BE` | Global Soft-Delete Database Scopes | P1 | Middleware and query hooks that append soft-delete filters to all read operations. | 1. All database queries on mutable tables (`fleets`, `users`, `bookings`, `attachments`) automatically append `WHERE deleted_at IS NULL`.<br>2. Mutating delete endpoints execute `UPDATE ... SET deleted_at = NOW()`.<br>3. Hard `DELETE` statements strictly blocked.<br>4. Expose admin endpoints `POST /admin/restore/{entity}/{id}` and `POST /admin/purge/{entity}/{id}`. | Soft-delete prevents physical deletion; default queries exclude soft-deleted rows; admins can restore or purge records. |
| `SD-01-MOB` | Soft-Delete Cascade Logic | P1 | Handles soft-deletion status changes on mobile client. | 1. Client app detects soft-deleted item responses.<br>2. Immediately hides matching UI elements and refetches list. | UI elements disappear smoothly on delete action without requiring full app restart. |
| `SD-01-WEB` | Admin Soft-Deleted Record Management UI | P1 | Next.js admin view to manage deleted records. | 1. Admin dashboard table includes "Show Deleted Records" toggle.<br>2. If active, displays records where `deleted_at IS NOT NULL` (rendered with red strikethrough).<br>3. Action icons: "Restore" (curved arrow) and "Purge" (trash can with confirmation). | Admin can view, restore, or permanently purge soft-deleted records. |

---

### 5.23 Ethiopian Public Holiday Calendar

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `LOC-02-BE` | Holiday & Weekend Excluded Schedulers | P0 | SLA calculations and escalation crons exclude public holidays and weekends. | 1. Schedulers query `public_holidays` table for matching dates.<br>2. Calculate deadlines (e.g. 3-day dispute window, ETA readiness alerts) by incrementing business days only.<br>3. Skips Saturdays, Sundays, and public holiday dates. | SLAs calculated accurately; holidays and weekends correctly bypassed; deadlines extend fairly. |
| `LOC-02-WEB` | Admin Holiday Configuration UI | P1 | Next.js UI to upload CSV and manually manage public holidays. | 1. Admin navigates to Calendar Settings.<br>2. Supports CSV file bulk upload containing holiday dates, English names, and Amharic names.<br>3. Displays editable table to manually add, edit, or delete specific holiday records. | CSV upload populates database; calendar is editable; changes take effect immediately. |

---

### 5.24 Marketplace Integrity & Abuse Controls

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `FR-01-BE` | Fraud Prevention & Scoring Engine | P0 | Backend evaluates profile registration and activity against suspicious duplicate profiles or payment slip OCR hashes. | 1. Checks matching verified identity references (reused `IdentityVerification.provider_ref`, e.g. same Fayda ID on multiple accounts).<br>2. Checks matching phone numbers across owners.<br>3. Compares uploaded payment slip perceptual hashes to catch duplicate receipts.<br>4. Sets status to `UNDER_REVIEW` on positive match. | Detects duplicate identity references and phone profiles immediately; flags matching receipt uploads; blocks action. |
| `FR-02-BE` | Suspicious Activity Engine & Review Lockout | P0 | Middleware blocks account operations if flagged `UNDER_REVIEW` from enumerated suspicious-signal triggers or OTP failures. | 1. Tracks failed OTP counts, payout changes, and 3 concrete suspicious signals: (a) login from new `device_id` combined with new geographic region; (b) $\ge 5$ failed OTP attempts in 1-hour window; (c) impossible travel geovelocity ($> 900\text{ km/h}$).<br>2. Any single signal firing updates user status to `UNDER_REVIEW`.<br>3. Middleware intercepts and rejects mutating API requests. | Account status switches to `UNDER_REVIEW` when any enumerated signal fires; alerts admin queues; rejects mutating calls. |
| `FR-03-BE` | Device Fingerprinting Audit Logging | P0 | Stores user device metadata in database logs for fraud and support investigations. | 1. Client app attaches device identifier, OS type, and version to login headers.<br>2. Saves `device_id`, `device_type`, `first_seen`, and `last_seen` in database logs. | Logs device profiles on auth; records mismatch histories for forensic review. |
| `AB-01-BE` | Review Moderation & Abuse Filters | P1 | Validates that booking transitions match reviews to prevent review bombing or fake ratings. | 1. Validates that reviewer had active, completed booking with target asset/owner.<br>2. Blocks reciprocal reviews within same 7-day window. | Rejects ratings from unrelated profiles; blocks reciprocal collusive ratings. |
| `AB-02-BE` | Messaging Profanity & File Size Filters | P0 | Filters dispute message text content and restricts attachment sizes in websocket sync. | 1. Runs regex filter for blacklisted profanity phrases.<br>2. Checks file sizes on uploads (max 5MB per message attachment). | Profanity filtered; uploads $> 5\text{MB}$ blocked with clear validation error. |
| `AB-02-MOB` | Report User Action | P1 | Mobile chat UI features allowing users to report abuse. | 1. Displays "Report User" action in dispute chat header.<br>2. Opens structured reason form.<br>3. Submits report to admin queues. | User report dispatches successfully; logs incident in administrative queue. |
| `AB-03-BE` | Review abuse reporting & content filter | P1 | Backend runs a content filter on review submission and lets the reviewed party report a review into the moderation queue. | 1. On review submission, runs profanity/PII filter over review text.<br>2. `POST /reviews/{id}/report` lets reviewed party flag a review into moderation queue (`RC-09`).<br>3. Admin can un-flag/restore a review with an audit record. | Profanity/PII filtered on submission; reviewed party can report to moderation queue; admin restore is audited. |

---

## 6. Functional Requirements — Phase 2 (Corporate & Multi-Tenant Accounts)

### 6.1 Corporate Account Management

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `CORP-001-BE` | Organization registration - Backend | P0 (Ph2) | Database schema and endpoints to create organization profile (TIN, license file, primary contact, type) with status PENDING, and handle admin approval. | 1. `POST /api/v1/organizations` accepts organization payload.<br>2. Stores business license and TIN certificate in SeaweedFS.<br>3. Sets status to `PENDING`.<br>4. `POST /api/v1/admin/organizations/{id}/approve` triggers status update and creates initial `org_admin` permissions. | Organization created in `PENDING` state; files stored securely; approval triggers corporate role assignment. |
| `CORP-001-MOB` | Organization registration - Mobile | P0 (Ph2) | Flutter interface for company representatives to input company info, TIN, primary contact, and upload business license files. | 1. Profile screen presents "Register Organization" action.<br>2. Multi-step form collects legal company name, trade registration number, TIN, and contact details.<br>3. Uploads trade license PDF.<br>4. Submits registration to backend. | Form validates required fields; files upload cleanly; screen updates with pending confirmation banner. |
| `CORP-001-WEB` | Organization registration - Web | P0 (Ph2) | Next.js admin page showing a review queue for organization approvals, with document viewer and approve/reject actions. | 1. Admin navigates to Corporate $\rightarrow$ Organization Queue.<br>2. Table displays company name, contact, and submission date.<br>3. Clicking opens license document preview.<br>4. Action buttons: "Approve Organization" or "Reject". | Table renders details accurately; approval updates state in backend and reloads queue. |
| `CORP-002-BE` | Member invitation & roles - Backend | P0 (Ph2) | Manage invitation codes, user-organization links, and role permission matrices in database. Enforce organization boundaries on endpoints. | 1. Endpoint `POST /api/v1/organizations/invitations` generates invitation token for phone number.<br>2. `OrganizationMember` table maps user to organization and role (`org_admin`, `corp_requester`, `corp_approver`, `fleet_manager`).<br>3. Middleware checks active member permissions on all corporate endpoints. | Database tracks invitations; role changes enforce custom permissions; non-members blocked from organization tenant data. |
| `CORP-002-MOB` | Member invitation & roles - Mobile | P0 (Ph2) | Flutter screens for Organization Admin to invite members via phone number and configure roles using a checkbox list of permissions. | 1. "Team Management" screen lists active and pending employees.<br>2. "Invite Member" button captures phone number and role selection.<br>3. Checkbox list allows customizing permissions.<br>4. Dispatches invite. | Validates phone format; list refreshes after sending invitation; permissions checkboxes submit changes successfully. |
| `CORP-003-BE` | Auto-approval limits - Backend | P0 (Ph2) | Implement backend validation rules that evaluate booking requests against member's dynamic auto-approval limits. | 1. On booking request, loads member spending limits from database.<br>2. If booking price, duration, and category are within limits, sets status to `APPROVED` directly.<br>3. Else, transitions status to `CORPORATE_PENDING_APPROVAL`. | Auto-approval triggers when all conditions satisfied; limits evaluated accurately against cost and duration; exceptions enter approval workflow. |
| `CORP-003-MOB` | Auto-approval limits - Mobile | P0 (Ph2) | Flutter form in organization settings where administrators set booking limits (amount, days, category types) for each employee. | 1. Member details screen presents "Configure Approval Limits".<br>2. Inputs for max booking minor-unit amount, max rental days, and category multi-select.<br>3. Submits configuration to backend. | Values persist to backend; bounds checks prevent invalid ranges; screen displays active limits clearly. |
| `CORP-004-BE` | Account switching - Backend | P0 (Ph2) | Support session context switching. Updates authentication tokens to include active organization ID or updates active session context. | 1. `POST /api/v1/auth/switch-context` receives target `organization_id`.<br>2. Asserts caller is active member.<br>3. Re-issues session token scoped with corporate permissions and pricing agreements. | Token contains active corporate scope; subsequent queries filter data by active organization context. |
| `CORP-004-MOB` | Account switching - Mobile | P0 (Ph2) | Flutter dropdown or menu option allowing users to toggle between Personal profile and Corporate profile. | 1. Header widget displays context switcher dropdown.<br>2. Selecting "Corporate" calls switch API, refreshes session token, and updates UI theme to corporate workspace. | Toggling switches session state; API calls fetch context-specific data; layout indicates active context clearly. |

---

### 6.2 Corporate Approval Workflows

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `CORP-005-BE` | Approval group creation - Backend | P0 (Ph2) | Store corporate approval group structures, member association details, and voting rules in the database. | 1. `ApprovalGroup` and `ApprovalGroupMember` tables store group rules.<br>2. `POST /api/v1/organizations/approval-groups` creates group.<br>3. Quorum types: `UNANIMOUS`, `MAJORITY`, `WEIGHTED`, `VETO`. | Group rules persist in DB; JSON structure holds quorum parameters; validates group members belong to organization. |
| `CORP-005-MOB` | Approval group creation - Mobile | P0 (Ph2) | Flutter page to create approval groups, add organization members, and assign quorum types and parameters. | 1. Admin opens Approval Settings $\rightarrow$ New Group.<br>2. Inputs group name, selects quorum type from dropdown, assigns members with weight toggles.<br>3. Submits group configuration. | Form dynamically adjusts inputs based on quorum type; member search filters employees; saves successfully. |
| `CORP-006-BE` | Booking approval execution - Backend | P0 (Ph2) | State machine to block bookings, collect votes from approvers, check quorum calculations, and transition status. | 1. Exceeding limits creates `BookingApprovalVote` rows; sets status to `CORPORATE_PENDING_APPROVAL`.<br>2. `POST /api/v1/bookings/{id}/approve` accepts vote (`APPROVE`, `REJECT`, `ABSTAIN`).<br>3. Quorum met $\rightarrow$ advances to `PENDING_PAYMENT`.<br>4. Veto vote $\rightarrow$ transitions to `CANCELLED`. | Booking blocked until quorum met; veto rejects booking immediately; all votes recorded in audit log. |
| `CORP-006-MOB` | Booking approval execution - Mobile | P0 (Ph2) | Flutter interface for approvers to review pending corporate bookings, inspect details, and cast votes. | 1. Corporate mode renders "Pending Approvals" badge.<br>2. Tapping opens booking details showing cost, project code, and limit excess justification.<br>3. Action buttons: "Approve", "Reject", "Abstain". | Detailed booking info rendered; casting vote calls API and updates list; shows voting tally progress. |
| `CORP-007-BE` | Multi-stage approval - Backend | P1 (Ph2) | Sequence logic to handle multi-stage approvals, checking transitions between stages in the workflow database. | 1. Extends schema with sequential stages (Stage 1: Site PM, Stage 2: CFO).<br>2. Stage advances to next only when current stage quorum is satisfied.<br>3. Push alerts dispatched to next stage approvers.<br>4. Rejection at any stage cancels workflow. | Stage 2 remains locked until Stage 1 approves; rejection halts workflow; transitions execute transactional-safely. |
| `CORP-007-MOB` | Multi-stage approval - Mobile | P1 (Ph2) | Visual representation of approval stages on corporate booking detail screens showing active, completed, and pending stages. | 1. Corporate booking detail screen renders vertical stage progress tracker.<br>2. Displays completed approver names, timestamps, and active stage waiting on action. | Timeline loads stage progress; indicates active stage clearly; updates dynamically as votes are cast. |
| `CORP-008-BE` | Timeout & escalation - Backend | P1 (Ph2) | Cron job to monitor pending corporate approvals, checking for expired stages and auto-escalating to fallback groups. | 1. Hourly cron compares stage creation timestamp against timeout limit.<br>2. If expired, marks stage `EXPIRED`, inserts fallback stage records, and dispatches escalation alerts. | Escalations occur automatically on deadline breach; fallback approvers receive push alerts; audit logged. |
| `CORP-008-MOB` | Timeout & escalation - Mobile | P1 (Ph2) | Display escalation warnings and notify fallback approvers of escalated bookings needing attention. | 1. App receives push alert for escalated task.<br>2. "Pending Approvals" list renders red "Escalated" badge on urgent items. | Badge renders for escalated state; list sorts escalated items to top; tapping navigates to approval screen. |

---

### 6.3 Corporate Pricing & Credit Terms

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `CORP-009-BE` | Negotiated pricing agreements - Backend | P0 (Ph2) | Store discount models (percentage, fixed, tiered) linked to organizations and filter bookings. Apply calculations in invoice generator. | 1. `CorporatePricingAgreement` table stores discount rules mapped to organization and fleet category.<br>2. Booking invoice calculator queries active agreements and recalculates daily rates.<br>3. Prevents final rate from dropping below zero. | Calculations apply discount; bounds checks prevent negative rates; logs agreement reference in invoice record. |
| `CORP-009-WEB` | Negotiated pricing agreements - Web | P0 (Ph2) | Next.js admin screen to create and manage corporate pricing contracts, configuring discount rules, bounds, and date ranges. | 1. Admin navigates to Corporate $\rightarrow$ Pricing Contracts.<br>2. Form captures organization, category, discount percentage, and start/end dates.<br>3. Saves contract to database. | Form validates numeric bounds; saves contract to database; updates active agreements table live. |
| `CORP-010-BE` | Credit limit & payment terms - Backend | P0 (Ph2) | Maintain organization credit accounts, check utilization on checkout, and generate PDF invoices with Net payment terms. | 1. `Organization` schema maintains `credit_limit_minor` and `current_utilization_minor`.<br>2. If `current_utilization + booking_total > credit_limit`, rejects checkout.<br>3. Generates invoice with payment terms (Net 30/Net 60).<br>4. Increments utilization upon confirmation. | Credit check blocks checkout transactional-safely; updates credit usage on confirmation; invoice created with payment terms. |
| `CORP-010-MOB` | Credit limit & payment terms - Mobile | P0 (Ph2) | Enable corporate renters to select "Corporate Credit" as payment method during booking checkout if limit permits. | 1. Checkout view presents "Corporate Credit Account" option.<br>2. Displays available credit balance.<br>3. Disables option with explanation if booking exceeds credit limit. | Credit option disabled if limit exceeded; checkout submits success status; updates remaining credit display. |
| `CORP-010-WEB` | Credit limit & payment terms - Web | P0 (Ph2) | Next.js admin settings to configure organization credit limits and invoicing parameters, and view billing statements. | 1. Admin searches organization in corporate dashboard.<br>2. Adjusts credit limit input and selects payment terms (Net 30, Net 60).<br>3. Displays statement of outstanding credit invoices. | Form validates limit values; changes save to DB; outstanding balances match aggregated booking ledger. |
| `CORP-011-BE` | Purchase Order management - Backend | P1 (Ph2) | Store booking PO references and decrement active PO allocations. | 1. `po_number` and `po_document_url` stored on Booking schema.<br>2. Tracks PO allocation ledger; blocks booking if allocation is depleted. | Database records PO metadata; ledger calculates balances correctly; API rejects requests when allocation exhausted. |
| `CORP-011-MOB` | Purchase Order management - Mobile | P1 (Ph2) | Renter checkout option in Flutter to input PO details and upload PO documents. | 1. Corporate checkout reveals "Purchase Order" fields.<br>2. User enters PO Number and attaches PO authorization PDF.<br>3. Submits booking payload. | Files upload correctly; validation checks input length; fields persist in submitted booking data. |

---

### 6.4 Corporate Fleet Management & 6.5 Reporting

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `CORP-012-BE` | Fleet manager roles & scope - Backend | P0 (Ph2) | Expose endpoints to assign managers to specific fleets. Enforce scope checks on all CRUD operations in the backend. | 1. `OrganizationFleetManager` maps user to array of fleet IDs.<br>2. Middleware verifies caller is assigned to targeted fleet ID before permitting updates.<br>3. Filters fleet lists by manager scope. | Scope checks restrict unauthorized access; queries filter returned lists by manager scope; changes logged in audit history. |
| `CORP-012-MOB` | Fleet manager roles & scope - Mobile | P0 (Ph2) | Flutter screen for Corporate Owner Admin to select fleet managers and assign them to specific fleets; restricted views for managers. | 1. Admin opens Manage Staff $\rightarrow$ Edit Manager.<br>2. Checkbox list allows selecting assigned machines.<br>3. Fleet manager login renders only assigned assets. | Scope edits save to backend; manager dashboard correctly displays subset of listings; non-assigned items inaccessible. |
| `CORP-013-BE` | Spend & utilization reports - Backend | P1 (Ph2) | Create aggregation queries to compile corporate spend and fleet utilization. Generate PDF and Excel file payloads. | 1. `GET /api/v1/reports/export` accepts date range, category, and format (`pdf`, `xlsx`).<br>2. Aggregates booking totals, utilization hours, and tax withholdings.<br>3. Generates and streams document. | Aggregated data matches booking totals; report formats generate without corruption; export secured by corporate permissions. |
| `CORP-013-WEB` | Spend & utilization reports - Web | P1 (Ph2) | Next.js corporate finance panel displaying interactive charts and export triggers for spend and utilization reports. | 1. Corporate portal displays interactive spend charts.<br>2. Filter bar allows selecting cost centers and date ranges.<br>3. "Export PDF" and "Export Excel" buttons download reports. | Filters dynamically update metrics; export calls trigger file downloads; dashboard UI handles high data volume gracefully. |

---

## 7. Functional Requirements — Phase 3 (Advanced Telematics & IoT)

| ID | Requirement | Priority | Description | Detailed Flow | Acceptance Criteria |
|---|---|---|---|---|---|
| `ADV-001-BE` | First payment-gateway adapters - Backend | P0 (Ph3) | Implement concrete `PaymentGatewayAdapter`s (TeleBirr, CBE Birr, MPESA, Stripe) on the `PAY-GW-01-BE` framework. | 1. Implement adapter methods (`initiate`, `verifyCallback`, `queryStatus`, `refund`).<br>2. Insert registry row in `PaymentGateway`.<br>3. Enable via `PAY-GW-05-WEB`. | New gateway goes live as adapter + registry row only; verified callback, idempotency, and auto-advance apply unchanged. |
| `ADV-001-MOB` | Gateway checkout hand-off - Mobile | P0 (Ph3) | Per-gateway mobile hand-off polish on top of `PAY-GW-02-MOB` dynamic method selection. | 1. Verify gateway launch path (deep-link, webview, USSD).<br>2. Return deep link resumes app into polling flow.<br>3. Display gateway-specific logos and localized instructions. | Launching external apps works cleanly; redirects return into polling flow without crashes; zero hardcoded gateway lists. |
| `ADV-002-BE` | Saved payment methods - Backend | P1 (Ph3) | Securely store tokenized payment methods associated with users, without storing raw card data. | 1. `POST /api/v1/payment-methods` persists gateway token reference, masked account, and provider ID.<br>2. `GET /api/v1/payment-methods` returns saved tokens. | No raw card/bank credentials hit server; tokens saved securely; API enables retrieval and deletion of stored payment methods. |
| `ADV-002-MOB` | Saved payment methods - Mobile | P1 (Ph3) | Flutter payment screen options to select a previously saved card/account, or input new details with a "Save payment method" toggle. | 1. Renders saved payment methods list with provider icons.<br>2. "Save payment method" toggle stores token on successful checkout. | Tapping saved method selects it for API call; checkout completes without re-entering account digits; writes token on success. |
| `ADV-003-BE` | Refund processing - Backend | P1 (Ph3) | Implement API to process full or partial refunds through payment gateway APIs. Update financial ledger and booking status. | 1. `POST /api/v1/admin/bookings/{id}/refund` verifies refund eligibility.<br>2. Calls gateway refund API.<br>3. Updates booking status and writes ledger transaction. | Refund amount validated against original transaction; API records refund status; email alerts renter of refund details. |
| `ADV-003-WEB` | Refund processing - Web | P1 (Ph3) | Next.js admin page with booking details showing refund buttons, transaction lists, and status updates. | 1. Admin views transaction details.<br>2. Clicks "Refund Booking".<br>3. Modal captures refund amount and justification.<br>4. Executes refund API call and updates status. | Refund dialog checks bounds; status updates instantly upon return; refund history displays inside transaction overview. |
| `ADV-004-BE` | Commission auto-deduction via PLATFORM_HELD custody | P0 (Ph3) | Automatic commission split at checkout realized by enabling a gateway with `custody_model = PLATFORM_HELD` per `PAY-GW-04-BE`. | 1. Route eligible checkout traffic to `PLATFORM_HELD` gateway.<br>2. Splits ledger as `TOPUP_PAYMENT_SPLIT`.<br>3. Wallet `COMMISSION_DEBIT` skipped for these transactions. | Split executes transactional-safely per `PAY-GW-04`; rounding per `CW-16`; platform and owner reconcile daily; zero double commission. |
| `ADV-005-BE` | Renter-side logistics booking - Backend | P1 (Ph3) | Dynamic checkout search API to match and attach logistics options (transport) to an active fleet booking request. | 1. `GET /api/v1/logistics/options` queries active transport fleets matching date and payload criteria.<br>2. Handles nested logistics booking and updates dependencies. | Search API returns transport fleets matching requirements; system handles nested transactions; cancellation cascades cleanly. |
| `ADV-005-MOB` | Renter-side logistics booking - Mobile | P1 (Ph3) | Add "I need transport" option in Flutter checkout flow, displaying compatible transport fleets and pricing on check. | 1. Checkout view includes logistics toggle.<br>2. Toggling on opens drawer listing available transport carriers.<br>3. Selection aggregates transport fee into checkout total. | Toggle displays options; selection aggregates total checkout price; booking parameters include logistics ID. |
| `ADV-006-BE` | Operator independent booking - Backend | P1 (Ph3) | Database schema and endpoints to manage Operator profiles, dynamic calendars, day rates, and independent booking workflows. | 1. `OperatorProfile` and `OperatorAvailability` models maintain calendars.<br>2. `GET /api/v1/operators` search endpoint.<br>3. `POST /api/v1/bookings/operators` creates booking. | Calendar blocks booked dates; booking cost calculated dynamically; security validation restricts unauthorized calendar modifications. |
| `ADV-006-MOB` | Operator independent booking - Mobile | P1 (Ph3) | Flutter interface for operators to manage availability calendar and day rates, and for renters to search and book operators. | 1. Operator setting: calendar widget to toggle dates, text input for rates.<br>2. Renter search: "Find Operators" tab with filtering, detail screen, and checkout. | Calendar updates sync to backend; renter search displays available operators with profiles; booking submits successfully. |
| `ADV-007-BE` | Dynamic pricing engine - Backend | P2 (Ph3) | AI pricing recommendation service that calculates suggested rates and updates active fleet pricing based on demand triggers. | 1. Recommendation service calculates dynamic rates based on regional demand and supply.<br>2. Respects owner-configured min/max bounds.<br>3. Updates search index prices. | Calculations respect bounds; price changes log in system history; rate updates apply to search queries. |
| `ADV-007-MOB` | Dynamic pricing engine - Mobile | P2 (Ph3) | Owner settings screen in Flutter to toggle AI dynamic pricing opt-in and define rate bounds. | 1. Fleet pricing screen presents "AI Smart Pricing" switch.<br>2. When active, displays min and max daily rate inputs.<br>3. Saves configuration to backend. | Switch toggles visibility of bounds inputs; values submit successfully; validations check that minimum rate is less than maximum. |
| `ADV-008-BE` | Telematics data ingestion - Backend | P2 (Ph3) | Set up MQTT broker to parse incoming IoT telematics data (GPS, hour meter values) and update fleet logs in database. | 1. MQTT broker ingests topic `/devices/{imei}/telematics`.<br>2. Maps IMEI to fleet asset.<br>3. Validates payloads; logs GPS coords and engine hours to database. | Incoming GPS/run-hour telemetry writes to database; values overwrite manual telemetry readings; invalid formats trigger error warnings. |
| `ADV-009-BE` | Public API - Backend | P2 (Ph3) | Developer API endpoints for search and booking. Expose key generator, check API keys in auth middleware, and enforce rate limits. | 1. Public endpoints under `/api/v2/public/*`.<br>2. Middleware reads and validates `X-API-Key` headers.<br>3. Redis rate-limits request frequencies per key. | Middleware blocks invalid keys; rate limiters return 429 Too Many Requests; booking logic matches standard tenant flows. |
| `ADV-009-WEB` | Public API - Web | P2 (Ph3) | Next.js developer page containing API reference documentation, key generation tools, and usage charts. | 1. Developer portal view allows generating API keys.<br>2. Displays active keys and usage charts.<br>3. Searchable OpenAPI reference docs. | Key generation queries API and displays output once; documentation is searchable; analytics charts display key utilization stats. |
| `ADV-009B-BE` | API versioning policy - Backend | P2 (Ph3) | Structure public API routes to use URL version prefixes and include deprecation warnings in headers of deprecated endpoints. | 1. Route mapping: version namespaces `/v1/`, `/v2/`.<br>2. Middleware inserts `Deprecation: date` and `Sunset: date` HTTP response headers for deprecated routes. | URL structure enforces versions; calling deprecated endpoints returns warning headers. |
| `ADV-010-BE` | Regional deployment - Backend | P2 (Ph3) | Deploy multi-region infrastructure templates, configuring databases to support localized data stores and global reads. | 1. Deploys database read replicas across regions.<br>2. Configures routing rules with data residency filters.<br>3. Replicates core catalog globally. | Cross-region replication latency within thresholds; requests map to nearest region; admin logs aggregate globally. |
| `ADV-011-BE` | Multi-currency support - Backend | P2 (Ph3) | Add currency columns to Price fields in the database. Build an exchange rate manager that updates daily and calculates conversions. | 1. Database models maintain `currency_code` (ETB, KES, USD).<br>2. Daily cron updates exchange rates from central bank API.<br>3. Payout and calculation services fetch conversions dynamically. | Exchange rate service logs daily updates; database values store original currency; conversion API returns accurate numbers. |
| `ADV-011-MOB` | Multi-currency support - Mobile | P2 (Ph3) | Display converted pricing values in the selected currency (ETB/USD/KES) throughout renter search and booking screens. | 1. Fetches exchange rates during startup.<br>2. Displays active currency symbols next to pricing amounts.<br>3. Currency selector in settings converts display rates. | Display converts dynamically depending on active locale/currency setting; calculations match exchange service values. |
| `ADV-011-WEB` | Multi-currency support - Web | P2 (Ph3) | Next.js admin page displaying billing, commission tracking, and analytics dashboards with currency switcher controls. | 1. Currency dropdown selector in dashboard header.<br>2. UI converts metric summaries and tables based on selection.<br>3. Preserves underlying minor units. | Changing currency refreshes calculations; formats display correctly with symbols. |
| `ADV-012-BE` | BI dashboards & fraud detection - Backend | P2 (Ph3) | Set up automated data syncs to analytical warehouses and run ML classification models to detect fraudulent booking patterns. | 1. ETL pipelines sync data to analytical store.<br>2. ML service scores bookings for fraud triggers.<br>3. High-risk bookings flagged into admin review queue. | ETL runs automatically on schedule; ML scorer flags abnormal booking behaviors; alerts trigger notifications. |
| `ADV-012-WEB` | BI dashboards & fraud detection - Web | P2 (Ph3) | Next.js admin interface integrating with BI data sources or custom charts showing retention cohorts, churn graphs, and fraud management feeds. | 1. "Fraud Control" dashboard lists flagged bookings with confidence scores.<br>2. Cohorts view displays visual retention and churn charts. | Dashboards render visual models; admin can select and confirm fraud flags, invoking resolution APIs. |

---


## 8. The Canonical 20-State Booking Machine & Concurrency Engine

### 8.1 The 20 Canonical Lifecycle States

| # | Status Key | Formal Name | Operational Definition | Terminal? |
|---|---|---|---|:---:|
| 1 | `DRAFT` | Draft Reservation | Initial booking parameters selected; pre-confirmation review. | No |
| 2 | `HOLD_PENDING` | Payment Hold Active | 60-minute mutual exclusion hold active in Valkey; awaiting payment slip. | No |
| 3 | `HOLD_EXPIRED` | Reservation Hold Expired | 60-minute payment window elapsed without evidence upload; lock freed. | **Yes** |
| 4 | `CANCELLED_BY_RENTER`| Renter Voluntary Abort | Renter cancelled booking prior to payment confirmation. | **Yes** |
| 5 | `PAYMENT_PROCESSING` | Payment Gateway Processing | Pluggable online gateway webhook initiated; awaiting verification. | No |
| 6 | `PAYMENT_FAILED` | Gateway Payment Failed | Online payment rail declined transaction. | No |
| 7 | `PAYMENT_EVIDENCE_UPLOADED` | Slip Uploaded / In Review | Renter uploaded bank transfer slip; awaiting owner verification. | No |
| 8 | `PAYMENT_DECLINED` | Payment Evidence Declined | Owner rejected transfer slip; awaiting retry, dispute, or cancel. | No |
| 9 | `CONFIRMED` | Booking Confirmed | Prepayment verified; cancellation penalty zone active. | No |
| 10 | `AWAITING_FLEET_READINESS`| Fleet Mobilization Delay | Fleet delayed in maintenance/transport; owner provided new ETA. | No |
| 11 | `PENDING_ADDITIONAL_PAYMENT`| Additional Fee Requested | Owner requested extra readiness/mobilization fee; renter deciding. | No |
| 12 | `READY_FOR_PICKUP` | Ready at Site / Handover | Equipment present at site/depot; awaiting joint field inspection. | No |
| 13 | `ACTIVE` | Active Rental Operation | Handover signed; engine meter ticking; equipment working. | No |
| 14 | `EXTENSION_REQUESTED` | Rental Extension Pending | Renter requested additional lease days; awaiting owner response. | No |
| 15 | `RETURN_INITIATED` | Return Handover Initiated | Rental period elapsed; conducting exit field inspection. | No |
| 16 | `DISPUTED` | General Contract Dispute | Dispute opened during payment or active operations; admin reviewing. | No |
| 17 | `DISPUTED_HOUR_METER` | Hour Meter Conflict | Pre/post-inspection engine meters conflict; awaiting admin ruling. | No |
| 18 | `DAMAGE_PARTIAL_RESOLVED` | Partial Damage Ruled | Admin ruled partial return damage liability; added to invoice. | No |
| 19 | `PENDING_FINAL_PAYMENT` | Final Balance Due | Excess hours or damage fees pending final payment slip upload. | No |
| 20 | `CLOSED` | Booking Settled & Closed | All payments verified, reviews submitted/expired, ledger closed. | **Yes** |

---

### 8.2 State Transition Matrix & Rules

| From State | To State | Triggering Actor / Event | Invariant Conditions & Guard Rules |
|---|---|---|---|
| `(Start)` | `DRAFT` | Renter | Fleet available for requested dates. |
| `DRAFT` | `HOLD_PENDING` | Renter | Taps "Request to Book"; acquires Valkey lock `lock:hold:fleet:{id}`; 60-min timer starts. |
| `HOLD_PENDING` | `HOLD_EXPIRED` | System Worker | 60-minute payment timer elapses without receipt upload; lock freed. |
| `HOLD_PENDING` | `CANCELLED_BY_RENTER` | Renter | Renter voluntarily cancels before slip upload; zero penalty. |
| `HOLD_PENDING` | `PAYMENT_EVIDENCE_UPLOADED`| Renter | Bank transfer slip uploaded; passes ClamAV antivirus scan (`MANUAL_TRANSFER`). |
| `HOLD_PENDING` | `PAYMENT_PROCESSING` | System | Online gateway payment checkout initiated (`GATEWAY`). |
| `PAYMENT_PROCESSING` | `CONFIRMED` | System Worker | Gateway webhook verified signature and amount match; auto-advances. |
| `PAYMENT_PROCESSING` | `PAYMENT_FAILED` | System Worker | Gateway payment declined or abandoned. |
| `PAYMENT_FAILED` | `HOLD_PENDING` | Renter | Retries checkout with alternate payment method within original hold window. |
| `PAYMENT_FAILED` | `HOLD_EXPIRED` | System Worker | Hold window elapses following payment failure. |
| `PAYMENT_EVIDENCE_UPLOADED`| `CONFIRMED` | Owner | Owner approves transfer slip; debits `OwnerWallet` commission in same transaction (`BK-17`). |
| `PAYMENT_EVIDENCE_UPLOADED`| `PAYMENT_DECLINED` | Owner | Owner declines receipt; selects reason code and enters notes. |
| `PAYMENT_DECLINED` | `PAYMENT_EVIDENCE_UPLOADED`| Renter | Renter uploads replacement transfer slip (retry). |
| `PAYMENT_DECLINED` | `CANCELLED` | Renter | Renter cancels after decline; zero cancellation penalty assessed. |
| `PAYMENT_DECLINED` | `DISPUTE_OPENED` | Renter | Renter disputes decline; creates `Dispute` record with origin `payment`. |
| `DISPUTE_OPENED` | `CONFIRMED` | Admin | Admin rules receipt is valid; forces confirmation; debits wallet commission. |
| `DISPUTE_OPENED` | `CANCELLED` | Admin | Admin rules payment invalid; cancels booking without penalty. |
| `CONFIRMED` | `READY_FOR_PICKUP` | Owner | Owner marks machine ready at site/depot. |
| `CONFIRMED` | `AWAITING_FLEET_READINESS`| Owner | Owner reports delay; enters reason and `readiness_eta` timestamp. |
| `CONFIRMED` | `CANCELLED` | Renter / Owner | Post-confirmation cancellation; penalty rules apply (`BK-CX-04`, `BK-CX-05`). |
| `AWAITING_FLEET_READINESS`| `READY_FOR_PICKUP` | Owner | Fleet arrives / maintenance complete; owner marks ready. |
| `AWAITING_FLEET_READINESS`| `PENDING_ADDITIONAL_PAYMENT`| Owner | Owner requests additional mobilization fee. |
| `AWAITING_FLEET_READINESS`| `CANCELLED` | Renter | Renter cancels because owner missed readiness ETA; zero penalty. |
| `PENDING_ADDITIONAL_PAYMENT`| `READY_FOR_PICKUP` | Renter + Owner | Renter uploads payment proof; owner confirms receipt. |
| `PENDING_ADDITIONAL_PAYMENT`| `CANCELLED` | Renter | Renter declines additional payment; **owner penalty assessed**. |
| `PENDING_ADDITIONAL_PAYMENT`| `CANCELLED` | System Worker | `additional_payment_deadline` exceeded; auto-cancels with owner penalty. |
| `READY_FOR_PICKUP` | `ACTIVE` | Both Parties | Both complete pre-inspection (`AR-01`), sign attestation (`AR-08`), confirm handover. |
| `READY_FOR_PICKUP` | `CANCELLED` | System Worker | Handover timeout (48h after start date); no-show penalty assessed if one party acted. |
| `READY_FOR_PICKUP` | `DISPUTED_HOUR_METER` | System / Admin | Pre-rental inspection hour meter discrepancy $> \max(1.0\text{h}, 0.05 \times \text{hours})$. |
| `DISPUTED_HOUR_METER` | `READY_FOR_PICKUP` | Admin | Admin resolves pre-rental meter dispute; inputs binding reading. |
| `ACTIVE` | `EXTENSION_REQUESTED` | Renter | Renter requests additional days; re-checks availability schedule. |
| `EXTENSION_REQUESTED` | `ACTIVE` | Owner | Owner accepts extension; updates end date; debits incremental commission. |
| `EXTENSION_REQUESTED` | `RETURN_INITIATED` | Owner / System | Extension rejected or expired; proceeds to return flow. |
| `ACTIVE` | `RETURN_INITIATED` | System / Owner | Scheduled end date reached or manual "Initiate Return" triggered. |
| `RETURN_INITIATED` | `DISPUTED_HOUR_METER` | System / Admin | Return hour meter discrepancy $> \max(1.0\text{h}, 0.05 \times \text{hours})$ or monotonic regression. |
| `RETURN_INITIATED` | `DAMAGE_PARTIAL_RESOLVED`| Admin | Return damage claim disputed; Admin rules partial liability. |
| `RETURN_INITIATED` | `PENDING_FINAL_PAYMENT` | System / Owner | Return inspection complete; excess hours or accepted damage charges $> 0$. |
| `RETURN_INITIATED` | `COMPLETED` | Owner | Return inspection complete; zero additional charges due; deposit release logged. |
| `DISPUTED_HOUR_METER` | `PENDING_FINAL_PAYMENT` | Admin | Admin resolves return meter reading; excess hours charges $> 0$. |
| `DISPUTED_HOUR_METER` | `COMPLETED` | Admin | Admin resolves return meter reading; excess hours charges $= 0$. |
| `DAMAGE_PARTIAL_RESOLVED`| `PENDING_FINAL_PAYMENT` | System | Partial damage amount added to final invoice. |
| `PENDING_FINAL_PAYMENT` | `COMPLETED` | Owner | Renter uploads final balance slip; owner confirms receipt (`RC-06`). |
| `COMPLETED` | `CLOSED` | System Worker | Both parties submit reviews or 7-day rating timeout expires (`RC-11`); ledger closed. |
| `CORPORATE_PENDING_APPROVAL`| `PENDING_PAYMENT` | Corp Approvers | (Phase 2) Internal corporate approval quorum met; advances to payment. |
| `CORPORATE_PENDING_APPROVAL`| `CANCELLED` | Corp Approvers | (Phase 2) Corporate approver rejects booking request. |

---

### 8.3 Cancellation Penalty Summary

- **Before CONFIRMED:** Zero penalty. Renter or owner may cancel freely.
- **After CONFIRMED (Renter Cancel):** Renter forfeits configured percentage of prepayment (default 20%, clamped `[0%, 50%]`). Remainder tracked via `RefundObligation` (`BK-CX-04-BE`).
- **After CONFIRMED (Owner Cancel):** Owner pays fixed penalty or percentage of prepayment to renter. Tracked via `CompensationObligation` (`BK-CX-05-BE`).
- **During DISPUTE or ACTIVE:** Direct voluntary cancellation is prohibited. Early return treated as normal completion.
- **Readiness ETA Exceeded:** Renter may cancel for free. Zero penalty.
- **Renter Declines Additional Payment:** Cancellation with **owner penalty** (owner caused the non-readiness).
- **Handover Timeout (48h):** Auto-cancels 48h after start date. Zero penalty if neither party inspected; no-show penalty assessed against absent party if one party completed inspection.

---

## 9. Non-Functional Requirements

### 9.1 Performance & Resource Budgets
- **Request/Response Compression:** Traefik gateway and Go middleware compress all JSON, CSV, and text payloads exceeding 1KB using Gzip or Brotli based on `Accept-Encoding` headers.
- **Search Latency:** PostGIS radius search response time P95 $\le 500\text{ ms}$ (target $\le 80\text{ ms}$) across 1,000 active listings.
- **Booking Creation:** P95 $\le 2.0\text{ s}$ under standard concurrency.
- **Image Upload:** P95 $\le 8.0\text{ s}$ for 10MB file including synchronous ClamAV stream inspection.
- **Admin Dashboard Load:** P95 $\le 3.0\text{ s}$ for up to 10,000 records.
- **Push Notification Delivery:** P95 $\le 30\text{ seconds}$ from event trigger to device delivery.
- **Load Testing Standards:** Prior to MVP launch, system must pass automated load test simulating 1,000 concurrent users (30% search, 30% view, 20% book, 15% pay, 5% admin) for 1 hour with P95 $\le 1.0\text{ s}$, P99 $\le 2.0\text{ s}$, and error rate $< 1\%$.

### 9.2 Availability & Disaster Recovery
- **Uptime SLA:** $99.5\%$ monthly uptime across core booking and payment APIs.
- **Backup Policy:** Schedulers perform hourly incremental WAL backups, daily database dumps, and weekly full system snapshots retained for 30 days in an isolated secondary region.
- **Recovery Metrics:** RPO $= 1\text{ hour}$ (target $0\text{s}$ via sync WAL streaming); RTO $\le 4\text{ hours}$ (target $< 60\text{s}$). Monthly restoration drills and quarterly disaster recovery simulation exercises.
- **Backup Encryption:** All backups encrypted at rest via AES-256 and in transit.
- **External Dependency Resilience:** Outbound third-party calls (Fayda, SMS, FCM) enforce 5s timeouts and circuit breakers that open after 5 consecutive failures with exponential backoff.

### 9.3 Scalability
- System supports 1,000 concurrent active users and 5,000 listed fleet assets in Phase 1 without architectural mutation (scaling to 10,000 users in Phase 3).
- Search response P95 remains $< 500\text{ ms}$ with 5,000 active fleet records and 1,000,000 stored files.
- Object storage (SeaweedFS S3) scales horizontally to support $> 1,000,000$ files.

### 9.4 Security & ASVS Level 3 Governance
- **Default-Deny Authorization:** Endpoints without explicit matching permissions return HTTP 403. Wildcard grants (`domain:*`) restricted to Super Admin.
- **Transport Hardening:** TLS 1.3 enforced across all ingress points with HSTS preloading (`max-age=63072000`).
- **Secrets Management:** Cryptographic keys and database credentials managed via OpenBao / Vault.
- **Rate Limiting:** Authenticated traffic capped at 100 req/min per user; unauthenticated public traffic capped at 30 req/min per IP; OTP requests capped at 5 requests/hour per IP and 10 requests/hour per phone number.
- **Antivirus Ingress:** Uploaded files scanned for malware via ClamAV daemon (`tcp://clamav:3310`) prior to S3 persistence. Infected files quarantined; upload rejected with HTTP 422.
- **Column Encryption:** Customer PII (phone number, National ID, TIN, bank accounts) encrypted at rest using AES-256-GCM.
- **Mandatory Admin MFA:** All Admin and Super Admin accounts require TOTP MFA. Administrative requests without active MFA session tokens blocked.
- **Breach Notification:** Confirmed personal data breaches reported to national data protection authorities within 72 hours and communicated to affected users without undue delay.
- **Key Rotation:** Column encryption keys and JWT signing keys rotated at least every 12 months.

### 9.5 Maintainability & Indexing Standards
- **Clean Architecture:** Go services strictly isolate domain interactors from transport and database drivers.
- **Database Migrations:** Sequential, deterministic SQL migrations (`.up.sql` and `.down.sql`).
- **Mandatory Indexing Policy:**
  - GIN indexes for JSONB specifications (`FleetSpecIndexed.specs`).
  - GiST PostGIS spatial index for geolocation queries (`fleets.location_point`).
  - B-tree indexes for all foreign keys, status columns, and sorting timestamps.
  - Partial B-tree indexes for `deleted_at WHERE deleted_at IS NULL`.
  - Trigram indexes (`pg_trgm`) for fuzzy text searches (`make`, `model`).
- **Optimistic Concurrency Control:** Stateful tables (`bookings`, `fleets`, `owner_wallets`, `inspections`) enforce `version INT NOT NULL DEFAULT 1`. Updates assert `WHERE id = $1 AND version = $2` and increment `version = version + 1`. Conflicting updates return HTTP 409 `CONCURRENT_MODIFICATION`.
- **Code Coverage Target:** Monotonic coverage ratchet: $\ge 80\%$ line coverage overall; $\ge 95\%$ on financial ledger and authentication modules.

### 9.6 Data Retention, Privacy & Legal Compliance
- **Compliance Baseline:** Adheres strictly to Ethiopian Personal Data Protection Proclamation No. 1321/2024.
- **PII Retention:** User PII retained for maximum 5 years post-account closure, then permanently anonymised.
- **Financial Audit Retention:** All booking records, invoices, and `financial_ledger_entries` retained for exactly 7 years.
- **GPS Telematics Pruning:** Raw high-precision operator GPS coordinates permanently pruned 30 days post-booking closure.
- **Media Storage Lifecycle:** Bank transfer slips: 7 years; inspection photos: 3 years; dispute evidence: 7 years; expired licenses: 5 years.
- **Consent Logging:** Explicit opt-in records logged for GPS tracking (`OperatorConsentLog`), marketing preferences (`UserNotificationPreference`), and identity provider checks (`IdentityVerification`).

### 9.7 API, Logging & QA Development Standards
- **Operational Monitoring & Alerting:** System alerts trigger on key performance metrics: Booking failure rate > 5% in 5 mins, OTP delivery success rate < 90% in 1h, or database storage utilization > 80%. Critical alarms route automatically to support escalation channels (PagerDuty, push, email, SMS) with defined primary, secondary, and executive escalation stages.
- **Observability Baseline:** All services must emit RED metrics (request rate, error %, request duration p50/p95/p99) at a `/metrics` endpoint, propagate W3C `traceparent` trace context end-to-end across service calls, and inject `trace_id`/`correlation_id` into every structured log line (see §9.7.2 Structured Logging Standards). Standard Grafana dashboards for API latency against the §9.1 performance SLAs, database connection-pool health, and core business KPIs must exist before MVP launch.

#### 9.7.1 API Request & Response Standards
All REST APIs must communicate using JSON payloads wrapped in a standard envelope. This ensures client-side parsing and exception handling are consistent across all mobile and web stakeholders.
- **API Versioning Policy:** Endpoints must include version namespaces in the URL route (e.g. `/api/v1/...`). Backward compatibility must be preserved across updates. If breaking changes are required, a new version path (`/api/v2/...`) must be deployed. Deprecated versions must return a header notification (`X-API-Deprecated: True`) and provide a sunset duration of at least 6 months before service termination.

##### A. Standard Success Response Envelope (HTTP 200/201)
```json
{
  "success": true,
  "data": {
    "id": "bk-1029",
    "status": "CONFIRMED",
    "amount": 25000.0,
    "currency": "ETB"
  },
  "meta": {
    "timestamp": "2026-06-21T17:16:20Z",
    "requestId": "req-9876-54321"
  }
}
```

##### B. Standard Error Response Envelope (RFC 7807/9457 problem details style)
```json
{
  "success": false,
  "error": {
    "code": "POLICY_NOT_ACCEPTED",
    "message": "You must accept the latest privacy policy to proceed.",
    "details": [
      {
        "field": "policy_version",
        "issue": "Mandatory version terms-v9 has not been accepted by this user account."
      }
    ]
  },
  "meta": {
    "timestamp": "2026-06-21T17:16:20Z",
    "requestId": "req-9876-54321"
  }
}
```

#### 9.7.2 Structured Logging Standards
All application logs must be printed to standard output in JSON format. This allows our Vector agent to capture and parse logs cleanly before routing them to Grafana Loki.

##### Example Log Structure (INFO/ERROR)
```json
{
  "timestamp": "2026-06-21T17:16:20.123Z",
  "level": "ERROR",
  "service": "booking-service",
  "correlation_id": "req-9876-54321",
  "trace_id": "5b23d9178f8b417c808f2a245647248f",
  "user_id": "usr-9876",
  "message": "Booking checkout failed due to policy non-acceptance",
  "error": {
    "code": "POLICY_NOT_ACCEPTED",
    "message": "User has not accepted mandatory policy terms-v9"
  }
}
```

#### 9.7.3 International QA Standards & Quality Gates
The platform enforces a strict automated quality assurance pipeline in the CI/CD workflow:
- **Test Pyramid Coverage:**
  - **Unit Testing:** Minimum **80% code coverage** required for Go backend services, Flutter widgets, and Next.js utility layers.
  - **Integration Testing:** Automated API contract tests and DB transaction validations.
  - **E2E Testing:** Flutter Integration Tests for core user flows (Login -> Find Fleet -> Book -> Prepayment / Payment-Evidence Upload). ("Escrow" here means the manual prepayment-evidence step — no held-funds account; see glossary.)
- **Static Analysis & Linting:**
  - Backend must pass `golangci-lint` without warnings.
  - Mobile client must pass `fvm flutter analyze --fatal-warnings`.
  - Frontend web panel must pass strict ESLint rules.
- **Security Scans:** Automated dependency vulnerability scans (`govulncheck`, `pnpm audit`) and static security analyses (`gosec`, `SonarQube`) on every PR build. No merge allowed with critical issues.

### 9.8 Locale & Dual Calendar Standards
- **Dual Calendar Integration:** System supports displaying and inputting dates in both Gregorian and Ethiopian (13-month) calendars.
- **Date Storage Standard:** All timestamps stored as standard Gregorian ISO 8601 UTC in database (`TIMESTAMPTZ`).
- **Presentation-Layer Conversion:** Ethiopian calendar conversions resolved entirely at the presentation layers (mobile client and web console) based on user preference.
- **Public Holiday Calendar:** Schedulers and SLAs exclude weekends (Saturday, Sunday) and dates listed in `public_holidays`.
- **System Timezone:** Schedulers resolve timezone dynamically from `SystemConfig` key `system_timezone` (default `Africa/Addis_Ababa`).

### 9.9 Accessibility Standards (WCAG 2.2 AA)
- **Compliance Target:** All mobile client screens and admin web console conform to WCAG 2.2 Level AA.
- **Color Contrast:** Text and interactive elements maintain minimum contrast ratio of $4.5:1$ for normal text and $3:1$ for large text. High-glare outdoor mode maintains $\ge 7.0:1$.
- **Assistive Technologies:** Full support for screen readers (TalkBack on Android, VoiceOver on iOS, NVDA on Windows) with semantic labels and accessibility hints. Visible focus rings on all interactive elements.
- **CI/CD A11y Gates:** Automated `axe-core` Playwright scans on web routes and Flutter semantic tree checks. Any Level AA violation fails CI.

---

## 10. Data Model Overview & Relational Architecture

### 10.1 Entity Summary (~111 tables)

**MVP (Phase 1) — 97 tables**
- `users`, `owner_profiles`, `operator_profiles`, `owner_operator_links`, `fleet_operators`, `fleet_attachments`
- `categories`, `spec_field_groups`, `spec_field_definitions`, `fleets`, `fleet_spec_values`, `fleet_spec_indexed`, `attachments`
- `fleet_availability`, `bookings`, `booking_state_machines`, `booking_state_transitions`, `booking_status_histories`, `payments`
- `inspection_photos`, `damage_claims`, `booking_inspections`, `booking_incidents`, `reviews`
- `commission_rule_sets`, `commission_rules`, `penalty_rule_sets`, `penalty_rules`, `pricing_rule_sets`, `pricing_rules`
- `notification_events`, `notification_templates`, `notification_rules`, `user_notification_preferences`
- `permissions`, `roles`, `role_permissions`, `user_roles`
- `admin_pages`, `admin_page_sections`, `admin_table_configs`, `admin_table_columns`, `admin_actions`, `admin_page_actions`, `admin_sidebar_items`, `admin_dashboard_widgets`
- `document_types`, `required_documents`, `user_documents`, `insurance_documents`, `search_filter_configs`
- `languages`, `translation_keys`, `translations`, `translation_bundle_versions`, `fleet_location_histories`
- `audit_logs`, `security_incidents`, `processing_purposes`, `disputes`, `dispute_messages`, `booking_spec_snapshots`, `saved_fleets`, `operator_consent_logs`, `owner_wallets`, `owner_wallet_transactions`, `refund_obligations`, `compensation_obligations`, `financial_ledger_entries`, `bank_statement_imports`, `invoice_sequences`, `payment_settings`, `app_versions`, `policy_documents`, `user_policy_acceptances`, `feature_flags`, `feature_flag_rules`, `public_holidays`, `system_configs`, `notification_channels`, `user_phone_history`, `user_email_history`, `user_webauthn_credentials`, `user_mfa_backup_codes`, `identity_providers`, `identity_verifications`, `payment_gateways`, `gateway_transactions`, `owner_payouts`, `dispute_reports`, `insurance_claims`
- `notification_archives`, `offline_upload_queues`, `faq_categories`, `faq_items`, `review_responses`, `cancellation_reasons`

**Phase 2 — Additional 14 tables**
- `organizations`, `organization_roles`, `organization_members`, `approval_groups`, `approval_group_members`, `approval_group_scopes`, `booking_approval_workflows`, `booking_approval_votes`, `corporate_pricing_agreements`, `purchase_orders`, `corporate_reports`, `user_active_sessions`, `organization_fleet_managers`, `fleet_manager_scopes`

**Currency Hedge (MVP):** All money-bearing tables include an ISO 4217 `currency` column defaulting to `ETB` (`Payment`, `OwnerWallet`, `OwnerWalletTransaction`, `RefundObligation`, `CompensationObligation`, `FinancialLedgerEntry`, invoices). Values stored as 64-bit integer minor units.

### 10.2 Key Entity Relationships
- `Fleet` $\rightarrow$ `Category` (many-to-one)
- `Fleet` $\rightarrow$ `FleetSpecValue` (one-to-many)
- `Booking` $\rightarrow$ `Fleet` (many-to-one)
- `Booking` $\rightarrow$ `Payment` (one-to-many)
- `Booking` $\rightarrow$ `parent_booking_id` (self-referencing optional FK for logistics sub-bookings)
- `Booking` $\rightarrow$ `BookingSpecSnapshot` (one-to-one; immutable snapshot at confirmation)
- `Review` $\rightarrow$ `ReviewResponse` (one-to-one; owner public response)
- `Booking` $\rightarrow$ `Dispute` (one-to-many; parent of `DisputeMessage`)
- `Booking` $\rightarrow$ `BookingInspection` (one-to-many; pre- and post-rental inspection records)
- `Booking` $\rightarrow$ `RefundObligation` (one-to-one, nullable; tracks owner refund obligations)
- `Booking` $\rightarrow$ `CompensationObligation` (one-to-one, nullable; tracks owner-penalty liabilities)
- `Booking` $\rightarrow$ `FinancialLedgerEntry` (one-to-many; immutable double-entry ledger)
- `FinancialLedgerEntry` $\rightarrow$ `BankStatementImport` (reconciliation match link)
- `User` $\rightarrow$ `UserWebAuthnCredential` (one-to-many; enrolled passkeys)
- `User` $\rightarrow$ `UserEmailHistory` (one-to-many; parallel to `UserPhoneHistory`)

### 10.3 Canonical Status Enumerations

**User Verification Statuses:**
- `verification_status`: `PENDING`, `VERIFIED`, `REJECTED` (general account KYC status).
- `identity_status`: `VERIFIED`, `UNVERIFIED`, `PENDING_REVIEW` (derived by policy engine across identity providers).

**Owner Account Statuses:**
- `PENDING_REVIEW`: Registration submitted, awaiting administrative review.
- `ACTIVE`: Approved and authorized to list machinery.
- `REJECTED`: Application rejected with explanation.
- `DOCUMENT_EXPIRED`: Business license expired; listings paused.
- `LOW_BALANCE_LOCK`: Wallet balance below negative limit; listings hidden.
- `SUSPENDED`: Administratively suspended (unpaid balance $> 90$ days).

**Fleet Statuses:**
- `DRAFT`, `PENDING_APPROVAL`, `ACTIVE`, `REJECTED`, `INCOMPLETE_SPEC`, `PAUSED`, `RENTED`, `INACTIVE`.

**Wallet Transaction Types (`owner_wallet_transactions.type`):**
- `TOPUP_MANUAL`: Owner top-up approved by admin.
- `COMMISSION_DEBIT`: Platform commission debited upon confirmation.
- `PENALTY_DEBIT`: Owner penalty debited for overdue obligations.
- `ADJUSTMENT`: Admin corrective credit or debit.
- `TOPUP_PAYMENT_SPLIT`: Commission deducted at source on `PLATFORM_HELD` gateway transactions.

**Dispute Origins (`Dispute.origin`):**
- `payment`: Payment evidence decline conflict.
- `return_damage`: Return physical condition conflict.
- `hour_meter`: Handover or return hour meter conflict.

**Settlement Enums:**
- `payments.settlement_method`: `MANUAL_TRANSFER`, `GATEWAY`.
- `PaymentGateway.custody_model`: `DIRECT_TO_OWNER`, `PLATFORM_HELD`.
- `GatewayTransaction.status`: `INITIATED`, `SUCCEEDED`, `FAILED`, `REFUNDED`.
- `OwnerPayout.status`: `PENDING`, `EXECUTED`, `FAILED`.

---

## 11. Dynamic Configuration Architecture

| Dynamic Element | Backing Relational Table | Managed By |
|---|---|---|
| Fleet Categories & Taxonomies | `Category` | Admin |
| Spec Fields per Category | `SpecFieldDefinition` | Admin |
| Booking Statuses & Transitions | `BookingStateMachine`, `BookingStateTransition` | Super Admin |
| Cancellation Penalty Rules | `PenaltyRuleSet`, `PenaltyRule` | Fleet Owner |
| Commission Rates & Tiers | `CommissionRuleSet`, `CommissionRule` | Super Admin (MFA Gated) |
| Notification Events & Templates | `NotificationEvent`, `NotificationTemplate` | Super Admin |
| Admin Sidebar & Page Layouts | `AdminSidebarItem`, `AdminPage` | Super Admin |
| Admin Table Columns & Actions | `AdminTableColumn`, `AdminAction`, `AdminPageAction` | Super Admin |
| Roles & Permission Mappings | `Role`, `Permission`, `RolePermission` | Super Admin |
| Corporate Approval Policies | `ApprovalGroup`, `ApprovalGroupMember` | Organization Admin |
| FAQ Content | `FAQCategory`, `FAQItem` | Admin |
| Dynamic UI Translations | `Language`, `TranslationKey`, `Translation` | Admin (OTA Publish) |
| Featured Fleet Highlights | `Fleet.IsFeatured`, `Fleet.FeaturedUntil` | Admin |

### 11.1 Dynamic Specification Schema Versioning & Backfill Rules
- **Migration Trigger:** Adding a mandatory spec field to a category's `SpecFieldDefinition` triggers catalog synchronization.
- **Listing Quarantine:** Existing fleet records in that category are set to `status = INCOMPLETE_SPEC` and hidden from public search indices.
- **User Notification:** Owner receives automated email and push alerts prompting specification parameter updates.
- **Reactivation:** Once owner submits valid values, listing restores to `ACTIVE` (or `PENDING_APPROVAL` if verification configured) and reappears in search.

---

## 12. Technology Stack & Deployment View

```
                              System Deployment Topology
                                           │
  [Mobile App / Web Users]                 ▼
                             ┌───────────────────────────┐
                             │    Traefik 3.x Ingress    │  TLS 1.3, Rate Limiting,
                             │     API Edge Gateway      │  HSTS Preloading
                             └─────────────┬─────────────┘
                                           │
                                           ├──────────────────────────────────┐
                                           ▼                                  ▼
                             ┌───────────────────────────┐      ┌───────────────────────────┐
                             │     Next.js 16 Admin      │      │     Go 1.27.3 Backend     │
                             │     Web Back-Office       │      │     Modular Monolith      │
                             │     (React Server Comp)   │      │     (Chi Router, sqlc)    │
                             └───────────────────────────┘      └─────────────┬─────────────┘
                                                                              │
                                           ┌──────────────────────────────────┼──────────────────────────────────┐
                                           ▼                                  ▼                                  ▼
                             ┌───────────────────────────┐      ┌───────────────────────────┐      ┌───────────────────────────┐
                             │      PostgreSQL 18.6      │      │        Valkey 9.1.3         │      │       SeaweedFS S3        │
                             │     + PostGIS 3.6.4 GiST    │      │     Distributed Cache     │      │     Encrypted Storage     │
                             │     (Relational DB)       │      │     & Mutex Hold Locks    │      │     & ClamAV Daemon       │
                             └───────────────────────────┘      └───────────────────────────┘      └───────────────────────────┘
```

| Subsystem Layer | Technology Pinned | Implementation Rules & Rationale |
|---|---|---|
| **Backend Core** | Go `1.27.3` (`net/http`, `go-chi/v5`) | Modular monolith with Clean Architecture. High concurrency, memory efficiency ($\le 256\text{MB}$ RSS). |
| **API Gateway** | Traefik `3.x` / Envoy Gateway | Lightweight edge routing, TLS 1.3 termination, rate-limiting, and Gzip/Brotli compression. |
| **Mobile Client** | Flutter `3.47.9` (Dart `3.9.x`) | Feature-First Clean Architecture, BLoC state management, Isar Community (`isar_community`) NoSQL offline store. |
| **Admin Web Console** | Next.js `16.x` (React 19 App Router) | React Server Components by default, Tailwind CSS v4, dynamic grids via `@tanstack/react-table`. |
| **Relational Database** | PostgreSQL `18.6` + PostGIS `3.6.4` | Spatial WGS 84 (`SRID 4326`) distance queries via GiST indexing, `TIMESTAMPTZ` UTC storage. |
| **Database Driver** | `jackc/pgx/v5` (`v5.7.1`) | Native binary protocol encoding with connection pooling (`pgxpool`). Heavy ORMs prohibited. |
| **SQL Compiler** | `sqlc` (`v1.27.0`) | Compiles type-safe Go structs and queries directly from plain SQL files. |
| **Cache & Distributed Locks** | Valkey `9.1.3-alpine` | High-throughput open-source Redis replacement for mutex hold locks (`SET NX EX`) and sessions. |
| **Object Persistence** | SeaweedFS S3 | High-performance Go-based object storage for inspection photos, PDF statements, and licenses. |
| **Antivirus Scanning** | ClamAV `clamd` (TCP 3310) | Synchronous stream scanning (`zINSTREAM`) on all uploaded media prior to S3 persistence. |
| **Telemetry & Observability** | OpenTelemetry, Zap, Loki, Vector | Structured JSON logging with Zap PII-masking core, distributed tracing, and RED metrics. |
| **Secrets Management** | OpenBao | Open-source Vault fork managing column-encryption master keys and JWT signing credentials. |

### 12.1 Target Deployment Environments & Promotion Strategy
The system operates across three strictly isolated environments:

| Dimension | Development (`dev`) | Staging (`staging`) | Production (`prod`) |
|---|---|---|---|
| **Backend Monolith** | `http://localhost:8080/api/v1` (Docker Compose, Mock OTP, Debug logs) | `https://staging-api.dinkmachinery.com/api/v1` (Staging host, Sandbox SMS/payments) | `https://api.dinkmachinery.com/api/v1` (Tier III DC, Live telco rails, Zero-PII logs) |
| **Admin Web Console** | `http://localhost:3000` (`.env.development`, React DevTools, Auth mock) | `https://staging-admin.dinkmachinery.com` (`.env.staging`, E2E test harness) | `https://admin.dinkmachinery.com` (Strict CSP, Mandatory MFA, Audit logs) |
| **Mobile Client** | `com.dinkmachinery.dink_mobile.dev` (`lib/main_dev.dart`, `Dink (Dev)`) | `com.dinkmachinery.dink_mobile.staging` (`lib/main_staging.dart`, `Dink (Staging)`) | `com.dinkmachinery.dink_mobile` (`lib/main_prod.dart`, `Dink Machinery`) |
| **Local Storage** | Embedded Isar `isar_dev` | Embedded Isar `isar_staging` | Embedded Isar `isar_prod` (Encrypted) |

### 12.2 Modular OpenAPI Contract Architecture & Per-Endpoint Specifications
To eliminate Git merge collisions and ensure deterministic cross-subsystem type safety:
- **Contract Source Tree (`contracts/`):** Each endpoint is maintained in an isolated file under `contracts/endpoints/<domain>/<action>.yaml` referencing shared domain entities and RFC 9457 Problem Details in `contracts/components/`.
- **Distribution Bundle:** High-speed bundling via `@redocly/cli bundle` generates `contracts/openapi.bundled.yaml` in $< 25\text{ms}$.
- **Downstream Codegen:** `oapi-codegen` compiles Go HTTP transport routers and request/response types; `openapi-typescript` generates Next.js TypeScript definitions.

---

## 13. MVP Scope & Implementation Roadmap

### 13.1 In Scope (MVP Foundation)
- Individual renter, owner, and operator registration with phone OTP and pluggable identity (Fayda reference adapter).
- Dynamic fleet listing with 5–20 multi-angle photos, EAV category specs, and bundled attachments.
- PostGIS spatial radius discovery and keyset cursor-paginated search feeds.
- 20-state rental booking machine with 60-minute payment reservation holds and OCC versioning.
- Direct bank transfer settlement (`MANUAL_TRANSFER`) with receipt verification and owner wallet commission debits.
- Field inspection stepper with camera-enforced hour meter capture, fuel gauge, and ClamAV antivirus inspection.
- Return damage triage, dual-party photographic diffs, and administrative dispute arbitration.
- Server-driven dynamic admin console (`/api/v1/admin/ui-config`) with complete RBAC and immutable audit logging.
- Multilingual dynamic i18n runtime supporting English and Amharic with OTA font delivery (`ADR-013`).
- Offline connectivity guard with local Isar NoSQL draft persistence and background photo upload retry queue.

### 13.2 Out of Scope (Deferred Post-MVP)
- Corporate organizational multi-tenant accounts and multi-stage internal approvals (Phase 2).
- Direct gateway automated commission splits (`PLATFORM_HELD`) (Phase 3).
- Automated CAN-bus hardware IoT telematics ingestion (Phase 3).
- Standalone independent operator hiring marketplace (Phase 3).
- AI dynamic pricing recommendations (Phase 3).
- Formal insurance claims adjudication lifecycle and insurer portal (Phase 3).

### 13.3 Implementation Roadmap (12 Months)

| Sprint Cluster | Focus Area | Deliverables & Milestones | Duration |
|---|---|---|---|
| **Sprints 1–2** | Platform Scaffold & Environments | Base Go modular monolith scaffold, CI/CD pipelines, Docker, Traefik gateway, PostgreSQL + PostGIS setup. | 1.0 Month |
| **Sprints 3–5** | Identity, Auth & KYC Engine | Phone OTP challenges, ES256 JWT lifecycle, Valkey sessions, Fayda identity adapter, owner KYC queues. | 1.5 Months |
| **Sprints 6–9** | Fleet Registry & Dynamic Specs | Category taxonomies, dynamic EAV spec engine, photo upload pipeline with ClamAV, availability calendar. | 2.0 Months |
| **Sprints 10–13**| Discovery & Spatial Search | PostGIS radius search, dynamic category filter facets, sorting queries, compare tables, wishlist. | 2.0 Months |
| **Sprints 14–16**| Booking Engine & Payment Proof | 20-state machine, 60-min hold mutex locks, bank slip upload, owner verification, cancellation logic. | 1.5 Months |
| **Sprints 17–18**| Dispute Center & SLA Escalation | Tri-party dispute chat, photographic damage diff viewer, 3-day SLA cron escalation, binding rulings. | 1.0 Month |
| **Sprints 19–20**| Returns, Damage & Rating | Return inspection wizard, monotonic meter check, excess hours billing, review moderation. | 1.0 Month |
| **Sprints 21–22**| Ledger, Wallets & Notifications | Immutable double-entry ledger, owner wallet top-ups and debits, FCM push alerts, SMS gateway fallback. | 1.0 Month |
| **Sprints 23–24**| Admin Dynamic UI & Audit Logs | Server-driven UI config engine (`/ui-config`), RBAC permission matrix, immutable audit log diff viewer. | 1.0 Month |
| **Sprints 25–26**| Onboarding, Offline & i18n | Multilingual OTA bundle sync, Ethiopic font delivery, Isar offline draft queue, onboarding tutorials. | 1.0 Month |
| **Sprints 27–28**| Security Hardening & ASVS L3 | Zero-PII logging audit, Semgrep SAST, penetration testing drills, secret scanning, Crashlytics. | 1.0 Month |
| **Sprints 29–30**| Verification, Load Testing & Launch | k6 geospatial load testing (P95 $< 80\text{ms}$), E2E user journeys, staging acceptance drills, production launch. | 1.0 Month |

---

## 14. Appendices

### Appendix A: Sample Specification Template — Mini Excavator

| Section | Field Key | Label | Type | Unit | Visible to Renter | Comparable | Required | Min | Max | Enum Values |
|---|---|---|---|---|:---:|:---:|:---:|---|---|---|
| **Identity** | `manufacturer_name` | Manufacturer | `STRING` | — | Yes | Yes | Yes | — | — | — |
| | `model_name` | Model Name | `STRING` | — | Yes | Yes | Yes | — | — | — |
| | `series` | Series | `STRING` | — | Yes | No | No | — | — | — |
| | `introduction_year`| Production Year | `INTEGER` | — | Yes | No | Yes | 1980 | 2026 | — |
| | `weight_class` | Weight Class | `ENUM` | — | Yes | Yes | Yes | — | — | Micro (<2t), Mini (2-5t), Midi (5-10t) |
| **Weights** | `operating_weight_kg`| Operating Weight | `DECIMAL` | kg | Yes | Yes | Yes | 500 | 50000 | — |
| | `shipping_weight_kg` | Shipping Weight | `DECIMAL` | kg | No | No | No | — | — | — |
| | `ground_pressure_kpa`| Ground Pressure | `DECIMAL` | kPa | Yes | No | No | 10 | 200 | — |
| **Engine** | `engine_manufacturer`| Engine Maker | `STRING` | — | Yes | No | Yes | — | — | — |
| | `gross_power_hp` | Gross Power | `DECIMAL` | hp | Yes | Yes | Yes | 5 | 500 | — |
| | `fuel_type` | Fuel Type | `ENUM` | — | Yes | Yes | Yes | — | — | Diesel, Electric, Petrol, Hybrid |
| **Working Ranges**| `max_digging_depth_mm`| Max Digging Depth| `INTEGER` | mm | Yes | Yes | Yes | 1000 | 20000 | — |
| | `max_reach_ground_mm`| Max Ground Reach | `INTEGER` | mm | Yes | Yes | No | 1000 | 30000 | — |
| **Dimensions** | `transport_length_mm`| Transport Length | `INTEGER` | mm | Yes | No | Yes | 2000 | 20000 | — |
| | `transport_width_mm` | Transport Width | `INTEGER` | mm | Yes | No | Yes | 800 | 5000 | — |
| | `transport_height_mm`| Transport Height | `INTEGER` | mm | Yes | No | Yes | 1000 | 5000 | — |
| **Safety & Cab** | `has_rops` | ROPS Certified | `BOOL` | — | Yes | Yes | No | — | — | — |
| | `has_fops` | FOPS Certified | `BOOL` | — | Yes | Yes | No | — | — | — |
| | `has_enclosed_cab` | Enclosed Cabin | `BOOL` | — | Yes | No | No | — | — | — |
| | `has_air_conditioning`| Air Conditioning | `BOOL` | — | Yes | No | No | — | — | — |

---

### Appendix B: Default Roles & Permissions Seed

#### Roles

| Role Key | Role Name | Description |
|---|---|---|
| `super_admin` | Super Administrator | Full system governance; manages roles, permissions, dynamic configs, and emergency state overrides. |
| `admin` | Administrator | Operational management: approvals, disputes, user management, analytics. |
| `support_agent` | Support Agent | Handles disputes and user inquiries only. |
| `finance_admin` | Finance Administrator | Manages commission tracking, owner payouts, wallet credit limits, and financial reports. |

#### Permissions

| Permission Key | Description | Super Admin | Admin | Support Agent | Finance Admin |
|---|---|:---:|:---:|:---:|:---:|
| `dashboard:read` | View admin operational dashboard | Yes | Yes | No | Yes |
| `fleets:list` | View fleet listings catalog | Yes | Yes | Yes | No |
| `fleets:approve` | Approve/reject fleet listings | Yes | Yes | No | No |
| `fleets:update` | Suspend/activate fleet listing | Yes | Yes | No | No |
| `users:read` | View user details and profile records | Yes | Yes | Yes | No |
| `users:update` | Suspend/reactivate user accounts | Yes | Yes | No | No |
| `users:approve` | Verify user trade documents and identity records | Yes | Yes | No | No |
| `disputes:list` | View disputes and support inquiries | Yes | Yes | Yes | No |
| `disputes:resolve` | Issue binding dispute rulings | Yes | Yes | No | No |
| `commissions:read` | View commission ledgers and reports | Yes | Yes | No | Yes |
| `commissions:approve`| Confirm commission payouts and wallet limits | Yes | No | No | Yes |
| `templates:manage` | Manage spec templates and category taxonomies | Yes | Yes | No | No |
| `roles:manage` | Manage roles and permission assignments | Yes | No | No | No |
| `config:manage` | Edit system and jurisdictional configurations | Yes | No | No | No |
| `notifications:manage`| Manage notification templates and gateways | Yes | No | No | No |
| `reports:read` | View advanced analytics and financial exports | Yes | Yes | No | Yes |
| `audit:read` | View immutable audit log records | Yes | Yes | No | No |
| `reviews:manage` | Flag and moderate customer reviews | Yes | Yes | No | No |
| `faq:manage` | Manage multilingual FAQ content | Yes | Yes | No | No |
| `bookings:override` | Force state transition override on bookings | Yes | No | No | No |

---

### Appendix C: Fleet Categories Seed Data

```text
EXCAVATION
  ├── Mini Excavator (1-5 tons)
  ├── Midi Excavator (5-10 tons)
  ├── Large Excavator (10-50 tons)
  └── Long-Reach Excavator
EARTHMOVING
  ├── Bulldozer
  ├── Motor Grader
  ├── Skid Steer Loader
  ├── Wheel Loader
  └── Backhoe Loader
LIFTING
  ├── Mobile Crane
  ├── Tower Crane
  ├── Forklift (Electric)
  ├── Forklift (Diesel)
  ├── Telehandler
  ├── Scissor Lift
  └── Boom Lift
HAULING
  ├── Dump Truck
  ├── Tipper Truck
  ├── Concrete Mixer Truck
  └── Water Tanker
COMPACTION
  ├── Vibratory Roller
  ├── Plate Compactor
  └── Tamping Rammer
CONCRETE & ASPHALT
  ├── Concrete Pump
  ├── Asphalt Paver
  └── Cold Planer
LOGISTICS
  ├── Lowboy Trailer (Light, <10 ton)
  ├── Lowboy Trailer (Medium, 10-30 ton)
  ├── Lowboy Trailer (Heavy, 30+ ton)
  ├── Flatbed Truck
  ├── Crane Truck
  └── Recovery Vehicle
ATTACHMENTS
  ├── Bucket (various sizes)
  ├── Hydraulic Hammer / Breaker
  ├── Earth Auger
  ├── Grapple
  ├── Fork
  ├── Ripper
  └── Thumb
AGRICULTURE
  ├── Tractor
  ├── Combine Harvester
  └── Crop Sprayer
```

---

### Appendix D: Booking State Machine Diagram (Textual)

Canonical set = **20 states** (see §8.1): `PENDING_PAYMENT`, `PAYMENT_EVIDENCE_UPLOADED`, `PAYMENT_DECLINED`, `DISPUTE_OPENED`, `DISPUTE_RESOLVED_RENTER`, `DISPUTE_RESOLVED_OWNER`, `DISPUTE_PARTIAL_RESOLVED`, `CONFIRMED`, `AWAITING_FLEET_READINESS`, `PENDING_ADDITIONAL_PAYMENT`, `READY_FOR_PICKUP`, `ACTIVE`, `RETURN_INITIATED`, `PENDING_FINAL_PAYMENT`, `COMPLETED`, `RATED`, `CANCELLED`, `DISPUTED_HOUR_METER`, `DAMAGE_PARTIAL_RESOLVED`, `CORPORATE_PENDING_APPROVAL`.

```text
[CORPORATE_PENDING_APPROVAL]  (Phase 2 entry)
  |---> PENDING_PAYMENT (approval granted)
  |---> CANCELLED (rejected)

FLEET LISTED (ACTIVE)
  |
  | Renter taps "Book Now" (Valkey lock acquired)
  v
PENDING_PAYMENT (auto-cancel on payment timeout, default 24h)
  |
  | Renter uploads payment slip
  v
PAYMENT_EVIDENCE_UPLOADED
  |
  | Owner reviews slip
  |-----> CONFIRMED (penalty zone starts; commission debited from OwnerWallet)
  |           |
  |           |---> READY_FOR_PICKUP
  |           |           |
  |           |           |---> ACTIVE (both confirm handover; no cancellation)
  |           |           |         |
  |           |           |         | Rental end reached
  |           |           |         v
  |           |           |     RETURN_INITIATED
  |           |           |         |
  |           |           |         |---> DISPUTED_HOUR_METER (discrepancy > max(1.0h, 0.05 × session_hours))
  |           |           |         |          |---> PENDING_FINAL_PAYMENT (charges > 0, admin resolved)
  |           |           |         |          |---> COMPLETED (charges = 0, admin resolved)
  |           |           |         |---> DAMAGE_PARTIAL_RESOLVED (return-time damage, admin partial ruling)
  |           |           |         |          |---> PENDING_FINAL_PAYMENT (auto)
  |           |           |         |---> PENDING_FINAL_PAYMENT (charges > 0)
  |           |           |         |---> COMPLETED (no charges)
  |           |           |                    |
  |           |           |     PENDING_FINAL_PAYMENT ---> COMPLETED (final payment confirmed)
  |           |           |                                     |
  |           |           |                       [both rate or 7-day timeout]
  |           |           |                                     v
  |           |           |                                   RATED
  |           |           |                                     |
  |           |           |                               [ledger settle]
  |           |           |                                     v
  |           |           |                                   CLOSED
  |           |           |
  |           |           |---> CANCELLED (voluntary = penalty; handover timeout 48h = no penalty)
  |           |
  |           |---> AWAITING_FLEET_READINESS
  |                       |
  |                       |---> READY_FOR_PICKUP (owner marks ready)
  |                       |---> PENDING_ADDITIONAL_PAYMENT
  |                       |          |---> READY_FOR_PICKUP (renter pays, owner confirms)
  |                       |          |---> CANCELLED (renter declines; owner penalty)
  |                       |---> CANCELLED (ETA exceeded, free refund)
  |
  |-----> PAYMENT_DECLINED
  |           |---> PAYMENT_EVIDENCE_UPLOADED (retry upload)
  |           |---> CANCELLED (free cancel)
  |           |---> DISPUTE_OPENED
  |                       |---> DISPUTE_RESOLVED_RENTER -> CONFIRMED
  |                       |---> DISPUTE_RESOLVED_OWNER -> CANCELLED
  |                       |---> DISPUTE_PARTIAL_RESOLVED -> PENDING_FINAL_PAYMENT
  |
  |-----> CANCELLED (timeout or manual pre-payment abort)
```

---

### Appendix E: Glossary (Quick Reference)

| Term | Definition |
|---|---|
| **Dink Machinery** | The brand name of the heavy machinery and logistics rental aggregator platform. |
| **Fleet** | A single rentable mechanical asset — heavy machinery or transport carrier. |
| **Attachment** | Add-on equipment rented standalone or bundled with a host machine. |
| **Direct Transfer Model** | MVP fund-flow: renter pays owner directly by bank transfer; platform debits commission from owner's pre-funded `OwnerWallet`. |
| **Prepayment** | The upfront rental portion required to confirm a reservation. |
| **Payment Evidence** | A photo or screenshot of a bank transfer receipt uploaded by the renter. |
| **Booking Timeout** | The window (hours) to upload payment evidence before auto-cancel (default 24h, configurable per fleet). |
| **Cancellation Penalty** | Pre-defined fee charged when cancelling after confirmation. |
| **Commission** | Platform revenue, calculated as a percentage of the rental base fee. |
| **Idempotency Key** | Unique UUIDv4 token sent with mutating API calls preventing duplicate operations. |
| **SLA** | Service Level Agreement — minimum performance or operational commitment. |
| **PostGIS** | Spatial database extension executing spheroidal radius distance queries. |
| **Valkey** | Open-source Redis replacement for mutex hold locks and distributed caching. |

---

### Appendix F: Complete Relational Schema Table Inventory

The complete Dink Machinery database architecture consists of 111 relational tables:

1. `users`
2. `owner_profiles`
3. `operator_profiles`
4. `owner_operator_links`
5. `fleet_operators`
6. `fleet_attachments`
7. `categories`
8. `spec_field_groups`
9. `spec_field_definitions`
10. `fleets`
11. `fleet_spec_values`
12. `fleet_spec_indexed`
13. `attachments`
14. `fleet_availability`
15. `bookings`
16. `booking_state_machines`
17. `booking_state_transitions`
18. `booking_status_histories`
19. `payments`
20. `inspection_photos`
21. `damage_claims`
22. `booking_inspections`
23. `booking_incidents`
24. `reviews`
25. `commission_rule_sets`
26. `commission_rules`
27. `penalty_rule_sets`
28. `penalty_rules`
29. `pricing_rule_sets`
30. `pricing_rules`
31. `notification_events`
32. `notification_templates`
33. `notification_rules`
34. `user_notification_preferences`
35. `permissions`
36. `roles`
37. `role_permissions`
38. `user_roles`
39. `admin_pages`
40. `admin_page_sections`
41. `admin_table_configs`
42. `admin_table_columns`
43. `admin_actions`
44. `admin_page_actions`
45. `admin_sidebar_items`
46. `admin_dashboard_widgets`
47. `document_types`
48. `required_documents`
49. `user_documents`
50. `insurance_documents`
51. `search_filter_configs`
52. `languages`
53. `translation_keys`
54. `translations`
55. `translation_bundle_versions`
56. `fleet_location_histories`
57. `audit_logs`
58. `security_incidents`
59. `processing_purposes`
60. `disputes`
61. `dispute_messages`
62. `booking_spec_snapshots`
63. `saved_fleets`
64. `operator_consent_logs`
65. `owner_wallets`
66. `owner_wallet_transactions`
67. `refund_obligations`
68. `compensation_obligations`
69. `financial_ledger_entries`
70. `bank_statement_imports`
71. `invoice_sequences`
72. `payment_settings`
73. `app_versions`
74. `policy_documents`
75. `user_policy_acceptances`
76. `feature_flags`
77. `feature_flag_rules`
78. `public_holidays`
79. `system_configs`
80. `notification_channels`
81. `user_phone_history`
82. `user_email_history`
83. `user_webauthn_credentials`
84. `user_mfa_backup_codes`
85. `identity_providers`
86. `identity_verifications`
87. `payment_gateways`
88. `gateway_transactions`
89. `owner_payouts`
90. `dispute_reports`
91. `insurance_claims`
92. `notification_archives`
93. `offline_upload_queues`
94. `faq_categories`
95. `faq_items`
96. `review_responses`
97. `cancellation_reasons`
98. `organizations`
99. `organization_roles`
100. `organization_members`
101. `approval_groups`
102. `approval_group_members`
103. `approval_group_scopes`
104. `booking_approval_workflows`
105. `booking_approval_votes`
106. `corporate_pricing_agreements`
107. `purchase_orders`
108. `corporate_reports`
109. `user_active_sessions`
110. `organization_fleet_managers`
111. `fleet_manager_scopes`

---