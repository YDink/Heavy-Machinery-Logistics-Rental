# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification — 14. Appendices (`docs/shared/requirements/14-appendices.md`)

**Document Identifier:** `SRS-DINK-CORE-V7.0-SEC-14`  
**Classification:** Canonical System Requirements Specification & Reference Appendices (Single Source of Truth)  
**Governing Standard:** IEEE Std 830-1998 / IEEE/ISO/IEC 29148:2018 (Systems and Software Engineering — Requirements Engineering)  
**Associated Artifacts:** `SRS-DINK-CORE-V7.0` (§14), `FLOW-BOOKING-STATE-MACHINE-DINK-CORE-V7.0`, `ARCH-DB-DINK-CORE-V7.0`, `contracts/openapi.yaml`, `ADR-001`, `ADR-002`, `ADR-003`, `ADR-004`, `ADR-008`, `ADR-010`, `ADR-011`, `ADR-015`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, Flutter 3.47.9 Mobile Client, Next.js 16.x Admin Console (latest CVE security patch), PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3, SeaweedFS S3  
**Target Operating Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), ISO 4217, ISO 3166-1 alpha-2  
**Status:** Approved Engineering Baseline

---

## Appendix A: Sample Specification Template — Mini Excavator

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

## Appendix B: Default Roles & Permissions Seed

### Roles

| Role Key | Role Name | Description |
|---|---|---|
| `super_admin` | Super Administrator | Full system governance; manages roles, permissions, dynamic configs, and emergency state overrides. |
| `admin` | Administrator | Operational management: approvals, disputes, user management, analytics. |
| `support_agent` | Support Agent | Handles disputes and user inquiries only. |
| `finance_admin` | Finance Administrator | Manages commission tracking, owner payouts, wallet credit limits, and financial reports. |

### Permissions Matrix

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

## Appendix C: Fleet Categories Seed Data

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

## Appendix D: Booking State Machine Diagram (Textual)

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

## Appendix E: Universal Glossary (Quick Reference)

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

## Appendix F: Complete Relational Schema Table Inventory (111 Tables)

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
