# Partition Decisions: 07-08 — Phase 3 Multi-Region Multi-Currency & Cross-Border Logistics Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `07-08`.

- **Governing System ADRs:** `ADR-012`, `ADR-015`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** REG-01-BE to REG-12-BE

---

## D-07-08-001: Real-Time Interbank FX Rate Ingestion Engine with Currency Conversion Freezing at Booking

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [REG-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Expanding across East Africa (Ethiopia ETB, Kenya KES, Uganda UGX, Tanzania TZS, Rwanda RWF) introduces multi-currency transactions. Foreign exchange volatility between booking initiation and payment confirmation creates severe financial reconciliation discrepancies if conversion rates fluctuate in-flight.

### 2. Considered Alternatives
1. **Option A:** Daily official central bank FX rate ingestion with rate locking at booking confirmation: exchange_rate_snapshot frozen on booking row and applied across entire lease lifecycle
   - *Pros:* Guaranteed price stability for contractor; zero currency drift dispute between owner and renter; regulatory compliance.
   - *Cons:* Platform assumes minor FX fluctuation risk between quote and confirmation.
2. **Option B:** Dynamic floating exchange rates re-calculated at final invoice settlement
   - *Pros:* Always uses current spot market rate.
   - *Cons:* Contractors furious when currency devaluations increase rental invoice amounts unexpectedly post-job.
3. **Option C:** Forcing all international transactions into US Dollars (USD)
   - *Pros:* Single currency accounting.
   - *Cons:* Explicitly illegal under domestic legal tenders in Ethiopia and Kenya; severe foreign currency shortages.

### 3. Chosen Solution & Technical Rationale
Option A: Central bank FX ingestion with exchange rate freezing at booking creation. The exchange_rate is snapshot into the booking record. All subsequent overtime calculations, damage disputes, and ledger entries reference the frozen snapshot rate.

### 4. Rejected Alternatives
- Option B: Rejected because floating post-rental currency adjustments cause bitter legal contract disputes.
- Option C: Rejected because transacting in foreign currency domestically violates central bank currency sovereignty laws.

### 5. Consequences & Invariants
- **Invariants Maintained:** Once a booking is created, its exchange_rate_snapshot is immutable and governs all child financial legs.
- **Operational Trade-offs:** FX snapshot records must be maintained in currency_exchange_rates table with source central bank attribution.

### 6. Verification
- **Test File:** `backend/internal/modules/currency/application/fx_snapshot_test.go`
- **Assertion:** TestFXSnapshot_FreezesRateAtBookingCreation asserts subsequent rate updates do not alter existing booking currency amounts.

### 7. Revisit Trigger
- Revisit if regional central banks launch unified East African monetary union currency.

---

## D-07-08-002: Cross-Border Customs Escrow and Tariff Calculation Architecture for Trans-National Rentals

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [REG-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Moving heavy machinery across national borders (e.g. Kenya to Ethiopia for corridor infrastructure construction) requires temporary importation carnets (COMESA Carnet de Passages), cross-border road permits, and customs bond escrow to guarantee re-exportation.

### 2. Considered Alternatives
1. **Option A:** Cross-border customs workflow module tracking temporary import permits, COMESA carnet document verification, and customs bond escrow ledger allocations
   - *Pros:* Enables high-margin international cross-border rentals; full compliance with EAC and COMESA trade treaties; automated permit tracking.
   - *Cons:* Requires specialized compliance verification workflows in admin console.
2. **Option B:** Prohibiting all cross-border rentals completely
   - *Pros:* Zero international trade complexity.
   - *Cons:* Leaves high-value regional border projects (e.g. LAPSSET corridor, Moyale border post) unserved.
3. **Option C:** Leaving all customs compliance entirely to the contractor offline
   - *Pros:* No software code needed.
   - *Cons:* Machinery impounded at border crossings, creating stranded asset disputes and insurance nullification.

### 3. Chosen Solution & Technical Rationale
Option A: Cross-border customs module. The booking flow detects trans-national rentals (owner country != site country), activates mandatory COMESA customs carnet document requirements, and allocates a customs bond deposit in the double-entry ledger.

### 4. Rejected Alternatives
- Option B: Rejected because regional cross-border infrastructure corridors represent immense rental revenue potential.
- Option C: Rejected because unguided customs crossings lead to impounded machinery and total operational paralysis.

### 5. Consequences & Invariants
- **Invariants Maintained:** A cross-border rental cannot transition to READY_FOR_PICKUP without verified customs carnet and transit permit documentation.
- **Operational Trade-offs:** Cross-border booking confirmation SLA extends to 5 business days for permit clearance.

### 6. Verification
- **Test File:** `backend/internal/modules/logistics/application/cross_border_test.go`
- **Assertion:** TestCrossBorder_RequiresCustomsCarnet asserts international booking is blocked without approved carnet document.

### 7. Revisit Trigger
- Revisit if African Continental Free Trade Area (AfCFTA) digital transit carnets automate border clearance.

---

