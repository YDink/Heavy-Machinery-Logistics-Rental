# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification — 07.01. In-App Payment Processing, Multi-Gateway Integration & Non-Custodial Direct Settlement Architecture (`docs/shared/requirements/07-01-in-app-payment-processing.md`)

**Document Identifier:** `SRS-DINK-CORE-V7.0-SEC-07-01`  
**Classification:** Canonical System Requirements Specification & Payment Integration Baseline (Single Source of Truth)  
**Governing Standard:** IEEE Std 830-1998 / IEEE/ISO/IEC 29148:2018 (Systems and Software Engineering — Requirements Engineering)  
**Associated Artifacts:** `SRS-DINK-CORE-V7.0` (§5.2, §5.5, §5.7, §5.9, §6.3, §7.0, §7.1, §8, §10, §12, §13.2), `FLOW-BOOKING-STATE-MACHINE-DINK-CORE-V7.0`, `FLOW-COMMISSION-DEBT-DINK-CORE-V7.0`, `FLOW-WIRE-DISPUTE-DINK-CORE-V7.0`, `ARCH-BACKEND-GO-DINK-CORE-V7.0`, `ARCH-FLUTTER-DINK-CORE-V7.0`, `ARCH-ADMIN-NEXT-DINK-CORE-V7.0`, `ARCH-SEC-DINK-CORE-V7.0`, `DS-DINK-CORE-V7.0`, `DS-PERSONAS-DINK-CORE-V7.0`, `contracts/openapi.yaml`, `ADR-001`, `ADR-002`, `ADR-004`, `ADR-005`, `ADR-008`, `ADR-010`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3, ClamAV Daemon (:3310), SeaweedFS S3, Flutter 3.47.8 / Dart 3.9 (Cross-Platform Mobile), Next.js 16.x (latest CVE security patch) / React 19 (Web Admin & Corporate Portal), Isar NoSQL 3.1+  
**Target Personas:** Civil Engineering Contractor (*Dawit Haile-Mariam* / *Patrick Kariuki*), Equipment Fleet Owner (*W/ro Almaz Tesfaye* / *Kamau Ndegwa*), Corporate Finance Controller (*Solomon Kassa*), Compliance Officer & Legal Counsel (*Yohannes Bekele*), Super Administrator (*Rahel Tadesse*)  
**Target Operating Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3 (Cryptographic Attestation, Payment Gateway Idempotency & Webhook Signature Verification), RFC 9457 (Problem Details), ISO 9241-110, W3C WCAG 2.2 Level AAA (Field Mobile) & Level AA (Admin Web), National Bank of Ethiopia (NBE) Payment Instrument Issuer Directives (ONPS/01/2020 & ONPS/09/2023), Central Bank of Kenya (CBK) National Payment System Act (Cap 491B), Safaricom Daraja 2.0 API Security Standards, Ethio Telecom TeleBirr SuperApp Partner Security Standard, Commercial Bank of Ethiopia CBEBirr Plus Integration Standard v3  
**Status:** Approved Engineering Baseline

---

## Table of Contents

- 1. [Document Scope & Architectural Mandate](#1-document-scope--architectural-mandate)
  - 1.1 [Purpose & Functional Scope](#11-purpose--functional-scope)
  - 1.2 [Core Financial & Payment Invariants](#12-core-financial--payment-invariants)
- 2. [Multi-Gateway Integration Architecture (`PAY-01`)](#2-multi-gateway-integration-architecture-pay-01)
  - 2.1 [Sovereign Multi-Country Payment Rail Matrix](#21-sovereign-multi-country-payment-rail-matrix)
  - 2.2 [Gateway Abstraction Layer & Payment Driver Topology](#22-gateway-abstraction-layer--payment-driver-topology)
  - 2.3 [High-Availability Failover, Latency Routing & Circuit Breakers](#23-high-availability-failover-latency-routing--circuit-breakers)
- 3. [Ethiopian Sovereign Payment Rails (`PAY-ET`)](#3-ethiopian-sovereign-payment-rails-pay-et)
  - 3.1 [TeleBirr (Ethio Telecom) Mobile Money & H5 SuperApp Deep-Linking](#31-telebirr-ethio-telecom-mobile-money--h5-superapp-deep-linking)
  - 3.2 [CBEBirr Plus & Commercial Bank of Ethiopia Direct Payment Gateway](#32-cbebirr-plus--commercial-bank-of-ethiopia-direct-payment-gateway)
  - 3.3 [Chapa Payment Orchestration (Local Cards, TeleBirr, CBE, Awash)](#33-chapa-payment-orchestration-local-cards-telebirr-cbe-awash)
  - 3.4 [Direct P2P Bank Wire Ingestion (CBE, Awash, Dashen, BoA)](#34-direct-p2p-bank-wire-ingestion-cbe-awash-dashen-boa)
- 4. [Kenyan Sovereign Payment Rails (`PAY-KE`)](#4-kenyan-sovereign-payment-rails-pay-ke)
  - 4.1 [Safaricom M-PESA Daraja 2.0 (STK Push / Lipa Na M-PESA Online)](#41-safaricom-m-pesa-daraja-20-stk-push--lipa-na-m-pesa-online)
  - 4.2 [M-PESA C2B & B2C Automated Owner Payout Rail](#42-m-pesa-c2b--b2c-automated-owner-payout-rail)
  - 4.3 [Card Acquiring Rails (Pesapal, I&M Bank, Standard Chartered KES/USD)](#43-card-acquiring-rails-pesapal-im-bank-standard-chartered-kesusd)
  - 4.4 [Kenyatta National Electronic Settlement (KETS / RTGS / Pesalink)](#44-kenyatta-national-electronic-settlement-kets--rtgs--pesalink)
- 5. [Payment Execution State Machine & Lifecycle Topology](#5-payment-execution-state-machine--lifecycle-topology)
  - 5.1 [End-to-End Payment Processing Flowchart](#51-end-to-end-payment-processing-flowchart)
  - 5.2 [Payment Transaction State Taxonomy (`PaymentState`)](#52-payment-transaction-state-taxonomy-paymentstate)
  - 5.3 [Permitted Payment State Transitions](#53-permitted-payment-state-transitions)
- 6. [Webhook Processing, Reconciliation & Cryptographic Proof (`PAY-02`)](#6-webhook-processing-reconciliation--cryptographic-proof-pay-02)
  - 6.1 [Universal Webhook Ingress & Dynamic Signature Verification](#61-universal-webhook-ingress--dynamic-signature-verification)
  - 6.2 [Strict Idempotency Protection (`Idempotency-Key` / Valkey Locks)](#62-strict-idempotency-protection-idempotency-key--valkey-locks)
  - 6.3 [Automated Webhook Retry Schedules & Dead-Letter Queue (DLQ)](#63-automated-webhook-retry-schedules--dead-letter-queue-dlq)
  - 6.4 [Double-Entry Payment Ledger & Commission Debt Accrual](#64-double-entry-payment-ledger--commission-debt-accrual)
- 7. [Relational Database Schema & Data Modeling (`PAY-DB`)](#7-relational-database-schema--data-modeling-pay-db)
  - 7.1 [Entity Relationship Diagram](#71-entity-relationship-diagram)
  - 7.2 [Data Definition Language (DDL) Specifications](#72-data-definition-language-ddl-specifications)
- 8. [Wire Protocol Contracts & RFC 9457 Error Schemas](#8-wire-protocol-contracts--rfc-9457-error-schemas)
  - 8.1 [`POST /api/v1/payments/intent`](#81-post-apiv1paymentsintent)
  - 8.2 [`POST /api/v1/payments/confirm-wire`](#82-post-apiv1paymentsconfirm-wire)
  - 8.3 [`POST /api/v1/payments/webhooks/{gateway}`](#83-post-apiv1paymentswebhooksgateway)
  - 8.4 [`GET /api/v1/payments/{paymentId}/status`](#84-get-apiv1paymentspaymentidstatus)
  - 8.5 [RFC 9457 Universal Problem Details Payloads](#85-rfc-9457-universal-problem-details-payloads)
- 9. [Automated Testing Suite & Verification Gates](#9-automated-testing-suite--verification-gates)

---

## 1. Document Scope & Architectural Mandate

### 1.1 Purpose & Functional Scope

This specification establishes **Section 07.01 (In-App Payment Processing, Multi-Gateway Integration & Non-Custodial Direct Settlement Architecture)** of the canonical Software Requirements Specification (SRS) for **Dink Machinery Version 7.0**.

In the East African commercial equipment hire and heavy haulage sectors across Ethiopia and Kenya, financial transactions present complex regulatory and operational constraints:
* **The Non-Custodial Regulatory Mandate:** Under National Bank of Ethiopia (NBE) payment direct licensing directives (ONPS/01/2020) and Central Bank of Kenya (CBK) regulations, software aggregators operating without a full trust-bank banking license are legally prohibited from pooling third-party gross rental funds into custodial escrow balances. Dink Machinery cannot act as a custodial clearing house for gross rental capital ($100,000\text{ to }300,000\text{ USD}$ machinery values; $50,000\text{ to }5,000,000\text{ ETB}$ contract values).
* **Direct Peer-to-Peer Settlement Model:** Gross rental payments (50% mobilization advance and final reconciliation balances) flow directly from the contractor's digital wallet or commercial bank account into the verified corporate account of the equipment owner.
* **Platform Monetization via Commission Debt (`FLOW-COMMISSION-DEBT`):** The platform generates revenue by assessing a standard platform service fee (e.g., $10.0\%$ machinery base hire, $2.5\%$ freight haulage). When an owner verifies receipt of funds from a contractor, the platform records a matching commission debt against the owner's virtual wallet (`owner_wallets`). Fleet owners settle their accumulated platform commission debt using in-app digital payment rails.
* **Extreme Payment Fragmentation:** Transactions span diverse sovereign rails:
  - **Ethiopia (`ET`):** TeleBirr (Ethio Telecom), CBEBirr Plus (Commercial Bank of Ethiopia), Chapa payment aggregation, and direct core banking RTGS/ACH transfers.
  - **Kenya (`KE`):** Safaricom M-PESA Daraja 2.0 (STK Push, C2B, B2C automated payouts), Pesapal acquiring, and Pesalink/RTGS banking networks.

This document formally specifies:
1. The multi-gateway payment orchestration engine abstracting regional payment rails into a unified driver interface.
2. The sovereign payment adapters for Ethiopia (TeleBirr, CBEBirr Plus, Chapa, P2P Wire) and Kenya (M-PESA Daraja 2.0, Pesapal, Pesalink).
3. The cryptographic webhook ingestion, HMAC signature verification, and idempotency protection engine.
4. The double-entry payment ledger and automated commission debt accumulation engine.

### 1.2 Core Financial & Payment Invariants

1. **Strict Non-Custodial Invariant (NBE / CBK Regulatory Floor):** Gross equipment rental capital must never touch Dink Machinery corporate bank accounts. Renter payments execute directly into the designated fleet owner's verified bank account or mobile money till. Dink Machinery only processes in-app transactions for direct platform commission debt settlement, booking administrative fees, or credit buffer top-ups.
2. **Double-Entry Financial Immutability Floor:** All payment transactions, ledger entries, and wallet balance adjustments must balance to zero minor units across debit and credit legs:
   $$\sum \text{Debits} \equiv \sum \text{Credits} \quad (\Delta = 0\text{ Minor Units})$$
   Tables `payment_transactions`, `owner_wallet_transactions`, and `commission_debt_ledger` are strictly append-only. Database triggers enforce `ON UPDATE DO INSTEAD NOTHING` and `ON DELETE DO INSTEAD NOTHING`.
3. **Strict Integer Minor Unit Representation:** All monetary amounts, transaction fees, withholding taxes, and commission debts are modeled as 64-bit signed integers representing currency minor units ($100\text{ Cents} = 1\text{ ETB}$ or $1\text{ KES}$). Floating-point arithmetic is prohibited across all payment runtimes.
4. **Idempotent Payment Execution SLO ($\le 5.0\text{ ms}$ Lock Resolution):** Every payment initiation and webhook callback must supply or extract an explicit `Idempotency-Key`. Re-submitting the same key resolves to the previous transaction state without triggering duplicate charges or duplicate wallet credits.
5. **Universal Webhook Cryptographic Verification:** Webhook callbacks from payment gateways must be authenticated via HMAC-SHA256 signatures or asymmetric public key validation before payload deserialization. Unsigned or mismatched webhooks are rejected with **RFC 9457 HTTP 401 Unauthorized (`ERR_WEBHOOK_SIGNATURE_INVALID`)**.

---

## 2. Multi-Gateway Integration Architecture (`PAY-01`)

### 2.1 Sovereign Multi-Country Payment Rail Matrix

```text
+----------------------------------------------------------------------------------------------------+
| DINK MACHINERY CANONICAL PAYMENT RAIL MATRIX (PAY-01)                                              |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| PAYMENT RAIL KEY    | COUNTRY     | CURRENCY          | TRANSACTION SCOPE  | SETTLEMENT MECHANISM  |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| TELEBIRR_DIRECT     | ET          | ETB (Minor Units) | Commission Settle, | Direct API / H5 Deep- |
|                     |             |                   | Advance Hold       | Link (Ethio Telecom)  |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| CBEBIRR_PLUS        | ET          | ETB (Minor Units) | Commission Settle, | CBEBirr MiniApp / H5  |
|                     |             |                   | Direct Rent Advance| Webview Channel :8888 |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| CHAPA_AGGREGATOR    | ET          | ETB (Minor Units) | Visa/Mastercard ET,| Hosted Checkout Modal |
|                     |             |                   | Awash/Dashen/BoA   | Webhook Callback      |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| P2P_BANK_WIRE       | ET          | ETB (Minor Units) | Gross Machinery    | Direct Peer-to-Peer   |
|                     |             |                   | Rental Settlements | Bank Slip OCR Audit   |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| MPESA_DARAJA_STK    | KE          | KES (Minor Units) | Commission Settle, | Safaricom Daraja 2.0  |
|                     |             |                   | Instant Advance    | STK Push (USSD Prompt)|
+---------------------+-------------+-------------------+--------------------+-----------------------+
| MPESA_B2C_PAYOUT    | KE          | KES (Minor Units) | Automated Carrier  | Safaricom Daraja B2C  |
|                     |             |                   | & Owner Payouts    | Payout Engine         |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| PESAPAL_V3          | KE          | KES / USD         | Visa, Mastercard,  | Pesapal 3.0 API /     |
|                     |             | (Minor Units)     | Mobile Money KE    | IPN Notification      |
+---------------------+-------------+-------------------+--------------------+-----------------------+
| PESALINK_RTGS       | KE          | KES (Minor Units) | Corporate Gross    | Direct Corporate Bank |
|                     |             |                   | Wire Settlements   | Wire Reconciliation   |
+----------------------------------------------------------------------------------------------------+
```

---

### 2.2 Gateway Abstraction Layer & Payment Driver Topology

The Go backend monolith isolates external payment APIs behind a unified internal interface:

```go
// PaymentDriver defines the canonical contract for all payment gateway integrations
type PaymentDriver interface {
    GetGatewayKey() GatewayKey
    InitiatePayment(ctx context.Context, intent *PaymentIntent) (*InitiateResult, error)
    VerifyWebhook(r *http.Request) (*WebhookPayload, error)
    QueryStatus(ctx context.Context, transactionRef string) (*TransactionStatus, error)
    Refund(ctx context.Context, refundReq *RefundRequest) (*RefundResult, error)
}
```

* **Pluggable Factory:** The payment engine instantiates the appropriate driver based on the user's country code (`ET` vs. `KE`), requested rail key, and operational policy rules.

---

### 2.3 High-Availability Failover, Latency Routing & Circuit Breakers

To guarantee uninterrupted checkout operations in regions experiencing intermittent telecom instability:
* **Circuit Breaker Pattern:** Each payment driver is wrapped in a Valkey-backed circuit breaker (`go-resiliency/breaker`). If a gateway experiences $> 30\%$ error rates or p95 latency $> 10.0\text{ seconds}$ over a 3-minute sliding window, the circuit trips to `OPEN`.
* **Dynamic Failover Cascade:**
  - In Ethiopia: If `TELEBIRR_DIRECT` trips, checkout automatically routes to `CHAPA_AGGREGATOR` (TeleBirr channel) or offers `CBEBIRR_PLUS`.
  - In Kenya: If `MPESA_DARAJA_STK` experiences gateway timeouts, the UI falls back to `PESAPAL_V3` (M-PESA channel) or manual Paybill instructions.

---

## 3. Ethiopian Sovereign Payment Rails (`PAY-ET`)

### 3.1 TeleBirr (Ethio Telecom) Mobile Money & H5 SuperApp Deep-Linking

TeleBirr integration operates across two distinct modalities:
1. **In-App Direct API (USSD Push):** Invokes Ethio Telecom's B2B payment API, triggering an immediate USSD payment prompt on the contractor's phone.
2. **TeleBirr SuperApp H5 MiniApp Deep-Link:** When executed inside the TeleBirr mobile container, the checkout invokes the native JavaScript communication bridge:
   ```javascript
   window.myJsChannel.postMessage(JSON.stringify({
     command: "pay",
     orderId: "ORD-2026-99412",
     amount: "106750.00",
     shortCode: "88120"
   }));
   ```
* **RSA-2048 Encryption:** All request payloads dispatched to Ethio Telecom are signed and encrypted using Ethio Telecom's public key certificate and Dink Machinery's private signing key.

---

### 3.2 CBEBirr Plus & Commercial Bank of Ethiopia Direct Payment Gateway

CBEBirr Plus integration adheres strictly to the official Commercial Bank of Ethiopia (CBE) MiniApp Integration Standard:

```mermaid
sequenceDiagram
    autonumber
    actor User as Contractor (Dawit)
    participant App as Dink Flutter App / Webview
    participant Core as Dink Go Monolith Core
    participant CBE as CBEBirr Gateway (:8888)

    App->>Core: Request Payment Checkout (Token Ingestion)
    Core->>CBE: GET /auth/user (Validate Bearer Token)<br/>Header: Authorization: Bearer {token}
    CBE-->>Core: HTTP 200 OK {"phone": "251911200491"}
    
    Core->>Core: Sort Payload Alphabetically & Compute Signature<br/>Key=amount&callBackURL=...&companyName=...&key={hashingKey}
    Core->>CBE: POST /auth/pay (Submit Payment Request)<br/>Headers: Authorization: Bearer {token}
    CBE-->>Core: HTTP 200 OK {"token": "ProcessedPaymentToken"}
    
    Core-->>App: Emit Payment Token to Webview Client
    App->>App: window.myJsChannel.postMessage(data)<br/>(CBEBirr Mobile App takes over payment confirmation)
    
    CBE->>Core: POST {callBackURL} (Notify Transaction Status)<br/>Payload: paidAmount, paidByNumber, txnRef, signature
    Core->>Core: Verify CBE Signature & Validate Token
    Core-->>CBE: HTTP 200 OK {"status": "SUCCESS"}
    Core->>Core: Transition PaymentState to SETTLED & Update Booking
```

* **Deterministic Sorting Invariant:** Payload keys must be sorted alphabetically before computing the SHA-256 HMAC signature using the private merchant hashing key.

---

### 3.3 Chapa Payment Orchestration (Local Cards, TeleBirr, CBE, Awash)

For contractors paying with domestic bank-issued debit cards or multi-channel wallets:
* Dink Machinery integrates the **Chapa API**:
  - Initializes payment session via `POST https://api.chapa.co/v1/transaction/initialize`.
  - Injects custom metadata: `booking_id`, `organization_id`, `payment_intent_id`.
  - Redirects mobile client to Chapa hosted payment sheet or renders native card elements.
* **Webhook Authentication:** Verifies the `Chapa-Signature` header using HMAC-SHA256 computed against the raw JSON webhook request body using the merchant secret key.

---

### 3.4 Direct P2P Bank Wire Ingestion (CBE, Awash, Dashen, BoA)

For high-value corporate rental advances that exceed mobile wallet transaction limits ($> 100,000.00\text{ ETB}$):
* Contractor executes a direct core banking RTGS/ACH transfer to the fleet owner's verified commercial bank account.
* Contractor uploads the stamped bank deposit voucher or electronic debit advice via `POST /api/v1/payments/confirm-wire` (`FLOW-07`).
* Ingestion triggers the automated perceptual hashing (`dHash <= 4`) and OCR text extraction engine (`SRS 05.24 §4`), verifying transaction references, amounts, and timestamps before notifying the fleet owner for manual statement confirmation.

---

## 4. Kenyan Sovereign Payment Rails (`PAY-KE`)

### 4.1 Safaricom M-PESA Daraja 2.0 (STK Push / Lipa Na M-PESA Online)

For instant mobile money collections in Kenya:
* **STK Push Invocation:** The Go monolith calls the Safaricom Daraja 2.0 Express endpoint (`POST /mpesa/stkpush/v1/processrequest`).
* **Password Generation:**
  $$\text{Password} = \text{Base64Encode}(\text{BusinessShortCode} \parallel \text{Passkey} \parallel \text{Timestamp})$$
* **Contractor Experience:** Contractor *Patrick Kariuki* receives an immediate SIM toolkit prompt on his Safaricom handset requesting his M-PESA PIN to authorize the payment.

---

### 4.2 M-PESA C2B & B2C Automated Owner Payout Rail

* **Customer to Business (C2B):** For contractors paying via Paybill or Buy Goods till numbers directly from their M-PESA menu, Dink Machinery exposes validated `ValidationURL` and `ConfirmationURL` endpoints.
* **Business to Customer (B2C Payouts):** When reimbursing safety deposits or disbursing logistics fees to lowboy drivers, the platform executes automated B2C payout calls (`POST /mpesa/b2c/v1/paymentrequest`), disbursing funds directly to the recipient's mobile number within $< 3.0\text{ seconds}$.

---

### 4.3 Card Acquiring Rails (Pesapal, I&M Bank, Standard Chartered KES/USD)

For cross-border contractors and international joint ventures operating in Kenya:
* Integrates **Pesapal v3 API** supporting 3D-Secure Visa and Mastercard debit/credit cards.
* Enforces IPN (Instant Payment Notification) listener endpoints authenticating notifications using Pesapal consumer keys and secrets.

---

### 4.4 Kenyatta National Electronic Settlement (KETS / RTGS / Pesalink)

High-value transactions in Kenya ($> 150,000\text{ KES}$) route via interbank clearing:
* Integrates bank instruction parsing for **Pesalink** (instant account-to-account transfer up to $999,999\text{ KES}$) and **KETS / RTGS** for multi-million shilling corporate construction contracts.

---

## 5. Payment Execution State Machine & Lifecycle Topology

### 5.1 End-to-End Payment Processing Flowchart

```mermaid
flowchart TD
    Checkout(["Contractor Initiates Payment Checkout<br/>POST /api/v1/payments/intent<br/>(Context: BookingID, Amount, RailKey, IdempotencyKey)"]) --> IdempotencyCheck{"Valkey Idempotency Lock<br/>SET payment:idempotency:{key} NX"}
    
    IdempotencyCheck -->|"Key Exists (In-Flight)"| ReturnCached["Return Cached Response / In-Flight Status<br/>HTTP 409: ERR_DUPLICATE_PAYMENT_REQUEST"]
    
    IdempotencyCheck -->|"Lock Acquired"| ResolveGateway["Resolve Target Gateway Driver<br/>(TeleBirr / CBEBirr / Chapa / MPESA / Wire)"]
    
    ResolveGateway --> CircuitBreakerCheck{"Is Gateway Circuit Breaker Closed?"}
    
    CircuitBreakerCheck -->|"Open (Outage Detected)"| FailoverRoute["Trigger Automated Failover Route<br/>(e.g., TeleBirr -> Chapa Fallback)"]
    CircuitBreakerCheck -->|"Closed (Healthy)"| CreateRecord["INSERT INTO payment_transactions<br/>Status: INITIATED"]
    
    FailoverRoute --> CreateRecord
    
    CreateRecord --> DispatchGateway["Dispatch API Call to Payment Rail<br/>(STK Push / MiniApp Token / Hosted URL)"]
    
    DispatchGateway --> CheckImmediate{"Immediate Sync Result?"}
    
    CheckImmediate -->|"Redirect / STK Prompted"| PendingClient["Status -> PENDING_USER_ACTION<br/>Return Client Directives (H5 / USSD / Form)"]
    CheckImmediate -->|"Synchronous Failure"| FailPayment["Status -> FAILED<br/>Release Valkey Idempotency Lock"]
    
    PendingClient --> AwaitWebhook["Step 2: Await Gateway Webhook Callback<br/>(POST /api/v1/payments/webhooks/{gateway})"]
    
    AwaitWebhook --> VerifySig{"Verify Cryptographic Signature<br/>(HMAC-SHA256 / RSA Public Key)"}
    
    VerifySig -->|"Invalid Signature"| RejectWebhook["HTTP 401: ERR_WEBHOOK_SIGNATURE_INVALID<br/>Log Security Incident & Alert SecOps"]
    
    VerifySig -->|"Signature Verified"| CheckStatusPayload{"Transaction Status in Payload?"}
    
    CheckStatusPayload -->|"Success / Settled"| CommitSettle["Step 3: Atomic Settlement Transaction<br/>1. UPDATE payment_transactions Status -> SETTLED<br/>2. Append to Double-Entry Ledger (PAY-DB)<br/>3. Assess Commission Debt on Owner Wallet<br/>4. Advance Booking Machine to CONFIRMED"]
    
    CheckStatusPayload -->|"Failed / Cancelled"| MarkFailed["Status -> FAILED<br/>Notify Contractor via Push & SMS"]
    
    CommitSettle --> EmitSSE["Broadcast Real-Time SSE Invalidation Frame<br/>event: payment_settled<br/>Mobile HUD Updates in < 100ms"]
```

---

### 5.2 Payment Transaction State Taxonomy (`PaymentState`)

Every payment transaction record in PostgreSQL 18.6 is governed by an explicit lifecycle enumeration:

| Status Key | Is Funds Transferred? | System & Booking Behavior |
|---|:---:|---|
| **`INITIATED`** | ❌ No | Intent generated; payment driver initializing gateway handshake. |
| **`PENDING_USER_ACTION`** | ❌ No | USSD prompt dispatched / H5 webview open; awaiting user PIN entry. |
| **`PAYMENT_VERIFYING`** | ⚠️ In Flight | Bank wire voucher uploaded; awaiting perceptual hash & statement audit. |
| **`SETTLED`** | ✅ Yes | Funds confirmed received; double-entry ledger updated; booking unlocked. |
| **`FAILED`** | ❌ No | Gateway declined, user cancelled prompt, or bank transfer bounced. |
| **`EXPIRED_TIMEOUT`** | ❌ No | Hold lease expired ($3,600\text{ seconds}$) without payment evidence. |
| **`REFUNDED`** | 🔄 Reversed | Transaction refunded to contractor pursuant to formal arbitration ruling. |

---

### 5.3 Permitted Payment State Transitions

```text
+----------------------------------------------------------------------------------------------------+
| PAYMENT TRANSACTION STATE TRANSITION MATRIX                                                        |
+------------------------------------+---------------------------------------------------------------+
| SOURCE STATE                       | PERMITTED TARGET STATES                                       |
+------------------------------------+---------------------------------------------------------------+
| INITIATED                          | PENDING_USER_ACTION, PAYMENT_VERIFYING, FAILED                |
| PENDING_USER_ACTION                | SETTLED, FAILED, EXPIRED_TIMEOUT                              |
| PAYMENT_VERIFYING                  | SETTLED, FAILED, EXPIRED_TIMEOUT                              |
| SETTLED                            | REFUNDED                                                      |
| FAILED                             | [Terminal Failure State — Allows Re-Attempt with New Intent]  |
| EXPIRED_TIMEOUT                    | [Terminal State — Hold Lease Released]                        |
| REFUNDED                           | [Terminal Accounting State — Ledgers Balanced]                |
+------------------------------------+---------------------------------------------------------------+
```

---

## 6. Webhook Processing, Reconciliation & Cryptographic Proof (`PAY-02`)

### 6.1 Universal Webhook Ingress & Dynamic Signature Verification

Payment gateways notify Dink Machinery of transaction completions via asynchronous HTTP webhooks:
* **The Verification Standard:**
  ```go
  // Universal Webhook Signature Verification Pipeline
  func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
      bodyBytes, _ := io.ReadAll(r.Body)
      r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore for downstream consumption

      gatewayKey := chi.URLParam(r, "gateway")
      driver := h.registry.Get(gatewayKey)

      if err := driver.VerifySignature(r, bodyBytes); err != nil {
          rfc9457.WriteProblem(w, http.StatusUnauthorized, "ERR_WEBHOOK_SIGNATURE_INVALID", err.Error())
          return
      }

      // Process Payload Idempotently
      h.processor.ProcessAsync(gatewayKey, bodyBytes)
      w.WriteHeader(http.StatusOK)
  }
  ```

---

### 6.2 Strict Idempotency Protection (`Idempotency-Key` / Valkey Locks)

Gateways frequently retry webhook deliveries. To prevent duplicate double-entry allocations:
1. The webhook handler extracts the gateway's unique transaction identifier (*e.g., Safaricom `MpesaReceiptNumber`, Chapa `tx_ref`, CBEBirr `txnRef`*).
2. Executes an atomic check-and-set in Valkey 9.1.3:
   ```text
   SET lock:webhook:{txnRef} "PROCESSING" EX 60 NX
   ```
3. If the key exists, the request returns `HTTP 200 OK` immediately, terminating duplicate processing.
4. If acquired, the processor verifies that `payment_transactions` has not already transitioned to `SETTLED`.

---

### 6.3 Automated Webhook Retry Schedules & Dead-Letter Queue (DLQ)

If an internal database outage or network glitch prevents webhook processing:
* The transaction payload is committed to an encrypted **Dead-Letter Queue (DLQ)** in Valkey: `queue:payments:dlq`.
* An asynchronous exponential backoff recovery worker retries failed events at intervals:
  $$T_{\text{retry}} = 2^n \times 15\text{ seconds} \quad (n = 1, 2, 3, 4, 5)$$
* If 5 retries fail, an urgent P0 alert routes to the on-call Site Reliability Engineer (*Yared Birhanu*) and Compliance Auditor (*Solomon Kassa*).

---

### 6.4 Double-Entry Payment Ledger & Commission Debt Accrual

When a payment settles (`SETTLED`):
1. **Direct Gross Rental Settlement (Contractor $\longrightarrow$ Fleet Owner):**
   - The platform verifies that the gross capital reached the owner's bank account.
   - The booking transitions to `CONFIRMED`.
2. **Platform Commission Debt Assessment (`FLOW-COMMISSION-DEBT`):**
   - The Go monolith executes an atomic transaction calculating the platform service fee:
     $$\mathcal{C}_{\text{debt\_minor}} = \left\lfloor \mathcal{V}_{\text{gross\_base\_hire}} \times \left( \frac{\text{BaselineBps}}{10,000} \right) \right\rfloor$$
   - Inserts balancing debit and credit entries into `commission_debt_ledger`.
   - Decrements `owner_wallets.current_balance_minor`, reflecting commission owed by the owner.

---

## 7. Relational Database Schema & Data Modeling (`PAY-DB`)

### 7.1 Entity Relationship Diagram

```mermaid
erDiagram
    Booking ||--o{ PaymentTransaction : "settled by"
    Organization ||--o{ PaymentTransaction : "originates"
    PaymentTransaction ||--o{ PaymentLedgerEntry : "generates"
    PaymentTransaction ||--o{ PaymentWebhookAudit : "receives"
    OwnerWallet ||--o{ CommissionDebtEntry : "tracks"
    PaymentTransaction ||--o{ CommissionDebtEntry : "triggers"

    PaymentTransaction {
        uuid id PK
        uuid booking_id FK
        uuid payer_user_id FK
        uuid payer_organization_id FK NULL
        uuid recipient_organization_id FK
        string payment_rail -- 'TELEBIRR_DIRECT', 'CBEBIRR_PLUS', 'CHAPA', 'MPESA_DARAJA', 'P2P_WIRE'
        string gateway_transaction_ref UK
        bigint amount_minor
        string currency -- 'ETB', 'KES'
        string status -- 'INITIATED', 'PENDING_USER_ACTION', 'PAYMENT_VERIFYING', 'SETTLED', 'FAILED', 'REFUNDED'
        string idempotency_key UK
        jsonb gateway_payload_json
        timestamptz settled_at
        timestamptz created_at
        timestamptz updated_at
    }

    PaymentLedgerEntry {
        uuid id PK
        uuid transaction_id FK
        string account_code -- 'RENTER_RECEIVABLE', 'OWNER_CLEARING', 'COMMISSION_REVENUE', 'TAX_PAYABLE'
        string entry_type -- 'DEBIT', 'CREDIT'
        bigint amount_minor
        string currency
        timestamptz created_at
    }

    PaymentWebhookAudit {
        uuid id PK
        string gateway_key
        string external_event_id
        text raw_payload_json
        string signature_header
        boolean signature_verified
        int response_status_code
        timestamptz processed_at
    }

    OwnerWallet {
        uuid id PK
        uuid owner_organization_id FK UK
        bigint current_balance_minor
        bigint credit_buffer_limit_minor
        string currency
        string status -- 'ACTIVE', 'DELINQUENT_WARNING', 'LOCKED_SUSPENDED'
        timestamptz updated_at
    }

    CommissionDebtEntry {
        uuid id PK
        uuid wallet_id FK
        uuid payment_transaction_id FK
        bigint debt_amount_minor
        bigint resulting_balance_minor
        string status -- 'ACCRUED', 'SETTLED', 'WRITTEN_OFF'
        timestamptz created_at
    }
```

---

### 7.2 Data Definition Language (DDL) Specifications

```sql
-- 1. Master Payment Transactions Table
CREATE TABLE payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    payer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    payer_organization_id UUID NULL REFERENCES organizations(id) ON DELETE RESTRICT,
    recipient_organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT,
    payment_rail VARCHAR(32) NOT NULL,
    gateway_transaction_ref VARCHAR(128) NULL,
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('ETB', 'KES')),
    status VARCHAR(32) NOT NULL DEFAULT 'INITIATED',
    idempotency_key VARCHAR(128) NOT NULL UNIQUE,
    gateway_payload_json JSONB NOT NULL DEFAULT '{}'::JSONB,
    settled_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_payment_status CHECK (status IN ('INITIATED', 'PENDING_USER_ACTION', 'PAYMENT_VERIFYING', 'SETTLED', 'FAILED', 'EXPIRED_TIMEOUT', 'REFUNDED')),
    CONSTRAINT chk_payment_rail CHECK (payment_rail IN ('TELEBIRR_DIRECT', 'CBEBIRR_PLUS', 'CHAPA', 'MPESA_DARAJA', 'MPESA_B2C', 'PESAPAL', 'P2P_WIRE', 'CORPORATE_CREDIT'))
);

CREATE INDEX idx_payment_transactions_booking ON payment_transactions (booking_id);
CREATE INDEX idx_payment_transactions_status ON payment_transactions (status);
CREATE UNIQUE INDEX uq_payment_gateway_ref ON payment_transactions (payment_rail, gateway_transaction_ref) WHERE gateway_transaction_ref IS NOT NULL;

-- 2. Double-Entry Payment Ledger Table (Append-Only Journal)
CREATE TABLE payment_ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES payment_transactions(id) ON DELETE RESTRICT,
    account_code VARCHAR(64) NOT NULL,
    entry_type VARCHAR(16) NOT NULL CHECK (entry_type IN ('DEBIT', 'CREDIT')),
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('ETB', 'KES')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_ledger_tx ON payment_ledger_entries (transaction_id);

CREATE RULE prevent_payment_ledger_update AS ON UPDATE TO payment_ledger_entries DO INSTEAD NOTHING;
CREATE RULE prevent_payment_ledger_delete AS ON DELETE TO payment_ledger_entries DO INSTEAD NOTHING;

-- 3. Webhook Audit Log Table (Cryptographic Verification Records)
CREATE TABLE payment_webhook_audits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gateway_key VARCHAR(32) NOT NULL,
    external_event_id VARCHAR(128) NULL,
    raw_payload_json TEXT NOT NULL,
    signature_header VARCHAR(512) NOT NULL,
    signature_verified BOOLEAN NOT NULL,
    response_status_code INT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_audits_gateway ON payment_webhook_audits (gateway_key, processed_at DESC);

-- 4. Owner Virtual Wallets Table (Tracks Platform Commission Debt)
CREATE TABLE owner_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_organization_id UUID NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    current_balance_minor BIGINT NOT NULL DEFAULT 0, -- Negative balance indicates commission debt owed
    credit_buffer_limit_minor BIGINT NOT NULL DEFAULT 5000000, -- 50,000.00 ETB default buffer
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('ETB', 'KES')),
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_wallet_status CHECK (status IN ('ACTIVE', 'DELINQUENT_WARNING', 'LOCKED_SUSPENDED'))
);

CREATE INDEX idx_owner_wallets_status ON owner_wallets (status);

-- 5. Commission Debt Journal Table (Append-Only)
CREATE TABLE commission_debt_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES owner_wallets(id) ON DELETE RESTRICT,
    payment_transaction_id UUID NOT NULL REFERENCES payment_transactions(id) ON DELETE RESTRICT,
    debt_amount_minor BIGINT NOT NULL, -- Positive represents new debt; Negative represents commission settlement
    resulting_balance_minor BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'ACCRUED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_debt_status CHECK (status IN ('ACCRUED', 'SETTLED', 'WRITTEN_OFF'))
);

CREATE RULE prevent_commission_debt_update AS ON UPDATE TO commission_debt_ledger DO INSTEAD NOTHING;
CREATE RULE prevent_commission_debt_delete AS ON DELETE TO commission_debt_ledger DO INSTEAD NOTHING;
```

---

## 8. Wire Protocol Contracts & RFC 9457 Error Schemas

### 8.1 `POST /api/v1/payments/intent`

Initializes a payment intent, resolving sovereign gateway adapters and acquiring locks.

#### Request Headers
```http
POST /api/v1/payments/intent HTTP/1.1
Host: api.dinkmachinery.com
Content-Type: application/json
Authorization: Bearer <userAccessToken>
Idempotency-Key: 99182a01-4b12-4f9a-8812-7b9a4c10bc21
```

#### Request Payload
```json
{
  "bookingId": "bk_019a8f60-9912-7b9a",
  "amountMinor": 10675000,
  "currency": "ETB",
  "paymentRail": "CBEBIRR_PLUS",
  "callBackUrl": "[https://api.dinkmachinery.com/api/v1/payments/webhooks/cbebirr](https://api.dinkmachinery.com/api/v1/payments/webhooks/cbebirr)",
  "payerPhoneNumber": "+251911200491"
}
```

#### Response: `HTTP 201 Created`
```json
{
  "success": true,
  "data": {
    "paymentId": "pay_019a8fe0-8812-7b9a-bc21-4491204812aa",
    "bookingId": "bk_019a8f60-9912-7b9a",
    "amountMinor": 10675000,
    "currency": "ETB",
    "paymentRail": "CBEBIRR_PLUS",
    "status": "PENDING_USER_ACTION",
    "gatewayDirectives": {
      "channel": "CBEBIRR_MINIAPP",
      "processedPaymentToken": "cbe_token_99142a8fe08812bc21...",
      "jsPostMessagePayload": {
        "command": "pay",
        "token": "cbe_token_99142a8fe08812bc21..."
      }
    },
    "expiresAt": "2026-10-29T14:40:00Z",
    "createdAt": "2026-10-29T13:40:00Z"
  },
  "timestamp": "2026-10-29T13:40:01Z"
}
```

---

### 8.2 `POST /api/v1/payments/confirm-wire`

Submits bank wire transfer evidence for non-custodial direct peer-to-peer rental settlement.

#### Request Headers
```http
POST /api/v1/payments/confirm-wire HTTP/1.1
Host: api.dinkmachinery.com
Content-Type: multipart/form-data
Authorization: Bearer <contractorAccessToken>
Idempotency-Key: 88182a01-5c12-4f9a-8812-8c0b5d21cd32
```

#### Multipart Form Body
* `bookingId`: String (`"bk_019a8f60-9912-7b9a"`)
* `bankCarrier`: String (`"COMMERCIAL_BANK_OF_ETHIOPIA"`)
* `transactionReference`: String (`"FT26301ABC88"`)
* `depositedAmountMinor`: BigInt (`10675000`)
* `wireSlipDocument`: Binary Stream (WebP/PDF, scanned in-flight via ClamAV socket `:3310`)

#### Response: `HTTP 200 OK`
```json
{
  "success": true,
  "data": {
    "paymentId": "pay_019a8fe5-1123-7b9a",
    "bookingId": "bk_019a8f60-9912-7b9a",
    "status": "PAYMENT_VERIFYING",
    "dHashFingerprint": "dhash_e8f2c019b78e3312",
    "ocrVerified": true,
    "ownerVerificationSlaHours": 24,
    "submittedAt": "2026-10-29T13:45:00Z"
  },
  "timestamp": "2026-10-29T13:45:01Z"
}
```

---

### 8.3 `POST /api/v1/payments/webhooks/{gateway}`

Universal webhook ingestion endpoint receiving gateway settlement notifications.

#### Request Headers
```http
POST /api/v1/payments/webhooks/cbebirr HTTP/1.1
Host: api.dinkmachinery.com
Content-Type: application/json
Authorization: Bearer cbe_bearer_token_here...
X-Signature: 8f4c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7
```

#### Request Payload
```json
{
  "paidAmount": "106750.00",
  "paidByNumber": "251911200491",
  "txnRef": "FT26301ABC88",
  "transactionId": "TXN-CBE-2026-99412",
  "transactionTime": "2026-10-29 13:42:15",
  "tillCode": "88120",
  "token": "cbe_token_99142a8fe08812bc21...",
  "signature": "8f4c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7"
}
```

#### Response: `HTTP 200 OK`
```json
{
  "status": "SUCCESS",
  "message": "Transaction verified and settled into double-entry ledger."
}
```

---

### 8.4 `GET /api/v1/payments/{paymentId}/status`

Polls real-time transaction status across gateway drivers.

#### Request Headers
```http
GET /api/v1/payments/pay_019a8fe0-8812-7b9a-bc21-4491204812aa/status HTTP/1.1
Host: api.dinkmachinery.com
Authorization: Bearer <userAccessToken>
```

#### Response: `HTTP 200 OK`
```json
{
  "success": true,
  "data": {
    "paymentId": "pay_019a8fe0-8812-7b9a-bc21-4491204812aa",
    "bookingId": "bk_019a8f60-9912-7b9a",
    "status": "SETTLED",
    "amountMinor": 10675000,
    "currency": "ETB",
    "gatewayTransactionRef": "FT26301ABC88",
    "settledAt": "2026-10-29T13:42:20Z"
  },
  "timestamp": "2026-10-29T13:43:00Z"
}
```

---

### 8.5 RFC 9457 Universal Problem Details Payloads

#### Example: Duplicate In-Flight Payment (`HTTP 409 Conflict`)
```json
{
  "type": "[https://api.dinkmachinery.com/errors/duplicate-payment-request](https://api.dinkmachinery.com/errors/duplicate-payment-request)",
  "title": "Duplicate Payment Request In Flight",
  "status": 409,
  "detail": "A payment intent with Idempotency-Key '99182a01-4b12-4f9a-8812-7b9a4c10bc21' is currently being processed by the CBEBirr gateway. Duplicate charge prevented.",
  "instance": "/api/v1/payments/intent",
  "code": "ERR_DUPLICATE_PAYMENT_REQUEST",
  "correlationId": "cor_019a8fe0-8812-7b9a-bc21-4491204812aa",
  "timestamp": "2026-10-29T13:40:05Z"
}
```

#### Example: Invalid Webhook Signature (`HTTP 401 Unauthorized`)
```json
{
  "type": "[https://api.dinkmachinery.com/errors/webhook-signature-invalid](https://api.dinkmachinery.com/errors/webhook-signature-invalid)",
  "title": "Gateway Webhook Cryptographic Verification Failed",
  "status": 401,
  "detail": "HMAC signature computed from payload does not match X-Signature header. Potential tampering or misconfigured merchant key.",
  "instance": "/api/v1/payments/webhooks/cbebirr",
  "code": "ERR_WEBHOOK_SIGNATURE_INVALID",
  "correlationId": "cor_019a8fe0-9914-4c12-8812-7b9a4c10bc21",
  "timestamp": "2026-10-29T13:42:18Z"
}
```

---

## 9. Automated Testing Suite & Verification Gates

```mermaid
flowchart TD
    PR["Pull Request Submitted"] --> Gate1{"Gate 1: OpenAPI Schema Parity<br/>kin-openapi validates payment wire schemas"}
    Gate1 -->|Pass| Gate2{"Gate 2: Double-Entry Ledger Invariant Gate<br/>Asserts sum of debits == credits on all rails"}
    Gate2 -->|Pass| Gate3{"Gate 3: Webhook Idempotency Fuzzing Gate<br/>Dispatches 100 duplicate webhooks simultaneously"}
    Gate3 -->|Pass| Gate4{"Gate 4: Gateway HMAC Signature Gate<br/>Asserts HTTP 401 on mutated payload bytes"}
    Gate4 -->|Pass| Gate5{"Gate 5: Non-Custodial Commission Accrual Gate<br/>Verifies wallet debt matches 10% base hire exactly"}
    Gate5 -->|Pass| Approved["Payment Processing Module Passed: Deploy Permitted"]

    Gate1 -->|Fail| Abort["Reject PR: Wire Schema Contract Mismatch"]
    Gate2 -->|Fail| Abort["Reject PR: Financial Imbalance Flaw"]
    Gate3 -->|Fail| Abort["Reject PR: Duplicate Settlement Vulnerability"]
    Gate4 -->|Fail| Abort["Reject PR: Cryptographic Security Failure"]
    Gate5 -->|Fail| Abort["Reject PR: Commission Accounting Deviation"]
```

1. **OpenAPI Schema Contract Gate (`kin-openapi`):** Validates all payment intent creation, wire confirmation, webhook callback, and status polling endpoints against `contracts/openapi.yaml`.
2. **Double-Entry Ledger Invariant Test:** Executes 1,000 synthetic payment cycles across all payment rails. Asserts that the double-entry accounting ledger maintains absolute balance parity down to zero minor units:
   $$\sum \text{Debits} - \sum \text{Credits} = 0\text{ Minor Units}$$
3. **Webhook Idempotency Stress Benchmark:** Dispatches 100 parallel, concurrent webhook callbacks sharing an identical `txnRef` across multiple goroutines. Asserts that exactly one transaction executes while 99 return `HTTP 200 OK` with zero duplicate ledger allocations.
4. **Cryptographic HMAC Signature Fuzzing Gate:** Mutates a single byte of incoming webhook payloads across TeleBirr, CBEBirr, Chapa, and M-PESA simulators. Asserts that $100\%$ of tampered payloads are intercepted with HTTP 401 `ERR_WEBHOOK_SIGNATURE_INVALID`.
5. **Non-Custodial Commission Accrual Assertion:** Simulates completed peer-to-peer wire confirmations, asserting that the Go backend accurately credits gross capital to the booking and debits the owner's virtual wallet by exactly $10.0\%$ base hire commission without touching gross funds.

---

*This document establishes the canonical functional requirements, sovereign multi-gateway integrations, non-custodial fund flows, cryptographic webhook verifications, double-entry ledgers, and wire protocol contracts for In-App Payment Processing in Dink Machinery Version 7.0. All downstream Flutter mobile apps, Next.js corporate consoles, Go backend payment microservices, and PostgreSQL financial databases must strictly conform to this specification.*