# Partition Decisions: 06-02 — Corporate Approval Workflows & Multi-Tier Authorization Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `06-02`.

- **Governing System ADRs:** `ADR-004`, `ADR-007`
- **Subsystem Scope:** backend
- **Requirement Targets:** APPR-01-BE to APPR-08-BE

---

## D-06-02-001: Multi-Tiered Dual-Authorization Workflow Engine for High-Value Equipment Reservations

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [APPR-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Enterprise machinery rentals involve contracts valued at millions of Birr. Corporate procurement governance mandates multi-tiered approvals (e.g. Project Engineer requisitions -> Project Manager approves < 200,000 ETB -> Procurement Director approves >= 200,000 ETB) before a reservation can proceed to payment.

### 2. Considered Alternatives
1. **Option A:** Flexible multi-tier workflow state engine (booking_approval_workflows) evaluating spend thresholds and routing to approval groups
   - *Pros:* Enforces corporate governance rules; prevents unauthorized rogue equipment rentals; full audit record of approval signatures.
   - *Cons:* Adds intermediate CORPORATE_PENDING_APPROVAL state before booking reservation hold.
2. **Option B:** Single approval step by organization owner for all transactions
   - *Pros:* Simple workflow.
   - *Cons:* Overwhelms corporate executives with approval requests for minor utility rentals.
3. **Option C:** Offline paper approval process where user uploads a signed requisition form
   - *Pros:* Zero workflow engine code.
   - *Cons:* Slows down procurement by days; unverified document fraud risks.

### 3. Chosen Solution & Technical Rationale
Option A: Multi-tiered workflow engine. When a corporate booking is submitted, the engine evaluates organization spend thresholds and routes the request to approval groups. The booking remains in CORPORATE_PENDING_APPROVAL until all required approvers cast positive votes.

### 4. Rejected Alternatives
- Option B: Rejected because enterprise construction firms require delegation of authority based on financial thresholds.
- Option C: Rejected due to human delay and lack of tamper-evident digital authorization.

### 5. Consequences & Invariants
- **Invariants Maintained:** A corporate booking cannot transition to PENDING_PAYMENT or CONFIRMED without meeting all tier approval conditions.
- **Operational Trade-offs:** Booking holds are only acquired in Valkey after internal corporate approval is finalized.

### 6. Verification
- **Test File:** `backend/internal/modules/corporate/application/approval_workflow_test.go`
- **Assertion:** TestApprovalWorkflow_RequiresDirectorForHighValue asserts a 500,000 ETB rental requires Director-level vote.

### 7. Revisit Trigger
- Revisit if automated ERP integration (SAP, Oracle) handles internal approvals via webhook.

---

## D-06-02-002: Parallel vs Sequential Approval Voting with Automated Procurement Escalation Timeouts

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [APPR-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
When multiple managers must approve a requisition, sequential approvals create cascading delays if one manager is traveling on site. If approvals stall, project start dates are missed and equipment reservations are lost.

### 2. Considered Alternatives
1. **Option A:** Configurable parallel voting groups with automated 24-hour escalation reminders and delegate proxy approvals
   - *Pros:* Fastest approval turnaround; transparent voting progress; automated escalation to backup approvers.
   - *Cons:* Requires cron scheduler to monitor pending approval timeouts.
2. **Option B:** Strict sequential approval chain without delegation or timeouts
   - *Pros:* Rigid hierarchy enforcement.
   - *Cons:* A single absent manager halts the entire company equipment requisition pipeline.
3. **Option C:** Majority vote rule across all corporate members regardless of department
   - *Pros:* Fast resolution.
   - *Cons:* Violates financial segregation of duties.

### 3. Chosen Solution & Technical Rationale
Option A: Parallel voting groups with automated timeout escalation. Approvers in the same tier can vote in parallel. If an approver fails to act within 24 hours, the system notifies their designated delegate or escalates to the corporate admin.

### 4. Rejected Alternatives
- Option B: Rejected because sequential bottlenecks paralyze emergency field equipment mobilizations.
- Option C: Rejected because departmental budget authority cannot be diluted by general staff voting.

### 5. Consequences & Invariants
- **Invariants Maintained:** Approval votes capture user UUID, timestamp, IP, and vote disposition (APPROVED, REJECTED).
- **Operational Trade-offs:** Requires maintaining active corporate delegate relationships.

### 6. Verification
- **Test File:** `backend/internal/modules/corporate/application/approval_timeout_test.go`
- **Assertion:** TestApproval_EscalatesAfterTimeout asserts a pending vote escalates to the designated proxy after 24 hours.

### 7. Revisit Trigger
- Revisit if push-notification interactive actionable buttons are approved for executive mobile devices.

---

