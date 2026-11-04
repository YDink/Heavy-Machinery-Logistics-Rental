# Dink Machinery Database Test Fixtures (Multi-Machine Parity)

## Purpose & Mandate
Developers and agents work across multiple workstations (e.g., Windows at the office, macOS at home). To eliminate environment divergence, broken manual tests, and missing database states when switching machines, **all test seed data, mock entities, and manual verification records must be committed as idempotent SQL scripts in this directory**.

Ad-hoc, uncommitted database rows that exist only on one local machine are strictly prohibited per `AGENTS.md § 1` and `docs/quality/MANUAL_TESTING_PROTOCOL.md § 2.5`.

---

## Directory Organization

```
scripts/fixtures/
├── README.md                           # Architecture and usage guidelines
├── apply-fixtures.sh                   # macOS / Linux / Git Bash runner
├── apply-fixtures.ps1                  # Windows PowerShell runner
├── 00_base_system_fixtures.sql         # Base roles, permissions, system admins, test users
├── 05_01_auth_fixtures.sql             # Test contractor and owner accounts, session tokens
├── 05_02_fleet_fixtures.sql            # Test machinery, categories, attachments, photos
├── 05_05_booking_fixtures.sql          # Test confirmed bookings, escrow holds, dates
├── 05_06_active_rental_fixtures.sql    # Active telemetry, check-in inspections, PIN challenges
└── 05_07_completion_fixtures.sql       # Return inspections, damage reports, settlements
```

---

## Core Rules for Fixture Scripts

1. **Strict Idempotency:**
   - Every statement MUST use `INSERT INTO ... ON CONFLICT (id) DO NOTHING` or `DO UPDATE SET ...`.
   - Re-running the scripts $N$ times must never fail with unique constraint violations or duplicate key errors.

2. **Deterministic Static UUIDs:**
   - Never use `gen_random_uuid()` in test fixtures.
   - Use fixed, deterministic UUIDs so foreign keys, HTTP tests, and API probes reference the exact same IDs on both Windows and macOS:
     - System/Admin: `00000000-0000-4000-8000-000000000001`
     - Test Contractor: `11111111-1111-4111-8111-111111111111`
     - Test Equipment Owner: `22222222-2222-4222-8222-222222222222`
     - Test Machine (CAT 320D Excavator): `33333333-3333-4333-8333-333333333333`
     - Test Booking (Confirmed): `55555555-5555-4555-8555-555555555555`

3. **Session Start Synchronization:**
   - When starting a new session or switching machines (Office Windows ↔ Home Mac), execute:
     - **macOS / Linux:** `bash scripts/fixtures/apply-fixtures.sh`
     - **Windows PowerShell:** `.\scripts\fixtures\apply-fixtures.ps1`
   - This ensures all test records are present in the local PostgreSQL instance without altering schema migrations.
