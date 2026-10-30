# Security & Threat Architecture — Dink Machinery

**Compliance Baseline:** OWASP ASVS Level 3, OWASP Mobile Top 10, RFC 9457[cite: 24]  
**Architecture Source:** `docs/shared/architecture/09-security-architecture.md`[cite: 32]

---

## 1. Authentication & Session Lifecycle

```mermaid
sequenceDiagram
  autonumber
  actor User as Client (Flutter / Next.js)
  participant Gateway as Traefik Ingress
  participant Auth as Auth Module
  participant Valkey as Valkey 9.1.3 Cache
  participant DB as PostgreSQL 18.6

  User->>Gateway: POST /auth/login (phone/email + password)
  Gateway->>Auth: Validate Credentials
  Auth->>DB: Query user by phone/email index (deleted_at IS NULL)
  alt Credentials Invalid
    Auth->>Valkey: Increment Failed Attempts Counter (30m TTL)
    Auth-->>User: 401 Unauthorized (attempts remaining)
  end
  alt Failed Attempts >= 10
    Auth-->>User: 429 Too Many Requests (Account Locked 30m)
  end
  alt MFA Enabled
    Auth-->>User: 200 OK (pre_auth_token + factor challenge)
    User->>Auth: POST /auth/mfa (pre_auth_token + TOTP/Passkey)
  end
  Auth->>DB: Record UserActiveSession (enforce max concurrent sessions)
  Auth-->>User: 200 OK (15m Access JWT + 30d Opaque Refresh Token)
```[cite: 32, 34]

---

## 2. Role-Based Access Control (RBAC) Matrix

Permissions follow the colon-separated format `resource:action` (Appendix B)[cite: 29]. Wildcard assignments are restricted to `super_admin`[cite: 24].

| Permission Key[cite: 29] | Description[cite: 29] | Super Admin[cite: 29] | Admin[cite: 29] | Support Agent[cite: 29] | Finance Admin[cite: 29] |
|---|---|:---:|:---:|:---:|:---:|
| `dashboard:read`[cite: 29] | View metrics & dashboards[cite: 29] | Yes[cite: 29] | Yes[cite: 29] | No[cite: 29] | Yes[cite: 29] |
| `fleets:approve`[cite: 29] | Approve/reject fleet listings[cite: 29] | Yes[cite: 29] | Yes[cite: 29] | No[cite: 29] | No[cite: 29] |
| `users:approve`[cite: 29] | Verify KYC documents & identity[cite: 29] | Yes[cite: 29] | Yes[cite: 29] | No[cite: 29] | No[cite: 29] |
| `disputes:resolve`[cite: 29] | Adjudicate dispute rulings[cite: 29] | Yes[cite: 29] | Yes[cite: 29] | No[cite: 29] | No[cite: 29] |
| `commissions:approve`[cite: 29] | Approve manual wallet top-up slips[cite: 29] | Yes[cite: 29] | No[cite: 29] | No[cite: 29] | Yes[cite: 29] |
| `bookings:override`[cite: 29] | Super Admin emergency state override[cite: 29] | Yes[cite: 29] | No[cite: 29] | No[cite: 29] | No[cite: 29] |
| `config:manage`[cite: 29] | Edit SystemConfig & country_configs[cite: 29] | Yes[cite: 29] | No[cite: 29] | No[cite: 29] | No[cite: 29] |

---

## 3. PII Inventory & Data Protection Controls

| Field | Entity | Storage Format | Access Scope |
|---|---|---|---|
| **National ID Number** | `User` / `IdentityVerification` | AES-256-GCM Encrypted[cite: 24] | Admin with `users:approve` only[cite: 29] |
| **Bank Account Number** | `OwnerProfile` | AES-256-GCM Encrypted[cite: 24] | Renter with active `PENDING_PAYMENT` booking[cite: 35] |
| **Phone Number** | `User` | AES-256-GCM Encrypted[cite: 24] | Authenticated counterparty post-handover[cite: 35] |
| **TIN Number** | `OwnerProfile` | AES-256-GCM Encrypted[cite: 32] | Financial invoice generation only[cite: 35] |

---

## 4. Rate Limiting Tiers

- **Unauthenticated Public Endpoints (Search/Catalog):** Max 30 requests/minute per source IP[cite: 24].
- **OTP Request Endpoints:** Max 5 requests/hour per source IP; max 10 requests/hour per phone number[cite: 24].
- **Authenticated User Endpoints:** Max 100 requests/minute per authenticated user session[cite: 24].
- **Violation Response:** HTTP 429 with standard `Retry-After: <seconds>` header[cite: 24].