# Testing Strategy & Verification Standards — Dink Machinery

**Standard Reference:** ISO/IEC/IEEE 29119 (Software Testing), ISO/IEC 25010 (Software Quality), OWASP ASVS L3[cite: 24].  
**Coverage Invariant:** Non-decreasing coverage ratchet. Minimum 80% line coverage on changed code; 95% branch coverage required on ledger, payments, auth, and state-machine engines[cite: 24].

---

## 1. Test Pyramid & Execution Tiers

```
         / \
        / E2E \         Synthetic User Journeys & Multi-Role Workflows (< 5%)[cite: 24]
       /-------\
      / Contract\       OpenAPI Schema & Consumer-Driven Verification (15%)[cite: 24]
     /-----------\
    / Integration \     testcontainers-go + PostgreSQL 18 + Valkey 8 (30%)[cite: 24, 32]
   /---------------\
  /   Unit Tests    \   Pure Domain, BLoC States, Invariants, & Math (50%)[cite: 24]
 /-------------------\
```

| Level | Scope | Execution Target | Tooling | Isolation Mechanism |
|---|---|---|---|---|
| **Unit** | Domain entities, aggregates, state transitions, BLoC events, pure calculations[cite: 24, 32]. | < 50ms per test | Go `testing`, `testify`, Dart `test`, `bloc_test`, Jest[cite: 25, 32]. | Pure memory, mocked interfaces, zero external I/O[cite: 24]. |
| **Integration** | Repositories, database queries, cache semantics, outbox relays, middleware pipelines[cite: 15, 32]. | < 2s per suite | `testcontainers-go`, `pgx/v5`, `isar_community`[cite: 25, 32]. | Ephemeral Docker containers, transactional test rollbacks[cite: 25, 32]. |
| **Contract** | Request/response envelopes, RFC 9457 Problem Details, client DTO serialization[cite: 24, 37]. | < 500ms per endpoint | Spectral, Stoplight, custom contract runners[cite: 37]. | `contracts/openapi.yaml` schema assertions[cite: 37]. |
| **E2E / System** | Complete reservation lifecycles (Booking → Attestation → Return → Ledger)[cite: 35]. | < 30s per scenario | Flutter integration test, Playwright[cite: 24, 32]. | Dedicated isolated staging / local Docker network[cite: 24, 37]. |

---

## 2. Core Testing Principles & Non-Negotiables

### 2.1 Determinism & Zero Flakiness
- **Flaky Test Quarantine:** Any test exhibiting non-deterministic behavior must be quarantined within 24 hours of first detection[cite: 24]. Quarantined tests do not block CI but create an immediate P1 remediation ticket[cite: 24].
- **No System Wall-Clock Dependencies:** Domain code must never call `time.Now()` directly. All time-sensitive operations (hold timeouts, dispute SLAs, cancellation grace periods) must accept an injected `Clock` interface[cite: 24]:
  ```go
  type Clock interface {
      Now() time.Time
  }
  ```
- **Hermetic State:** Tests must not share mutable state. Database integration tests execute inside an isolated transaction rolled back on cleanup or run against disposable test databases provisioned per test process[cite: 24, 32].

### 2.2 Financial & Balance Invariants (`BK-17-BE`, `CW-17-BE`)
Every financial test modifying balances or transactions must assert the double-entry accounting identity:
```text
ASSERT: SUM(financial_ledger_entries.amount WHERE direction = 'debit') == 
        SUM(financial_ledger_entries.amount WHERE direction = 'credit')
```
- **Integer Minor-Unit Integrity:** Tests must assert that no monetary calculation produces fractional minor units[cite: 35]. All monetary parameters must evaluate as signed `int64` minor units (`santim`/`cents`)[cite: 35].
- **Idempotency Replay Invariance:** Re-executing any financial command with an identical `Idempotency-Key` must return an identical response payload without creating duplicate ledger or wallet rows[cite: 35].

### 2.3 Zero-PII Leakage Assertion
Automated regression tests must inspect log output buffers, vector sinks, and RFC 9457 error response payloads to assert the total absence of plain-text PII (names, phone numbers, National ID numbers, bank accounts)[cite: 24].

---

## 3. The 7-Step Contract Verification Protocol

Executed whenever an engineering increment affects endpoints shared between `backend/`, `mobile/`, and `admin/`[cite: 37]:

1. **Local Infrastructure Verification:** Confirm local Docker infrastructure is healthy[cite: 37]:
   ```bash
   docker compose ps
   ```
2. **Backend Execution:** Launch backend API service using `.env.dev`[cite: 37]:
   ```bash
   cd backend && go run cmd/api/main.go
   ```
3. **Spec Alignment Check:** Verify Swagger UI (`/swagger/index.html`) correctly serves the updated `contracts/openapi.yaml`[cite: 37].
4. **Live Integration Execution:** Run client integration tests against the live local backend (mocked servers are prohibited for contract sign-off)[cite: 37].
5. **Envelope Validation:** Validate that both success payloads and error payloads match standard envelopes (`meta.requestId`, `meta.version`, RFC 9457 problem fields)[cite: 24, 37].
6. **Drift Resolution:**
   - If client fails: Update client code to match `contracts/openapi.yaml`[cite: 37].
   - If contract requires modification: Author an ADR or partition decision doc before altering `contracts/openapi.yaml`[cite: 37].
7. **Sign-off:** Update requirement status to `done` in `ROADMAP.md` only after Step 5 passes[cite: 37].

---

## 4. Backend Testing Standards (Go 1.27.3)

### 4.1 Unit Testing Strategy
- Unit tests reside in the same package directory as the code under test (`*_test.go`)[cite: 15].
- Pure domain entities and aggregates are tested without external mocks:
  ```go
  func TestBooking_CalculateExcessHours(t *testing.T) {
      t.Parallel()
      booking := domain.NewBooking(domain.BookingParams{
          DailyRate:          500000, // 5,000.00 ETB
          AgreedDays:         2,
          StandardDailyHours: 8,
          ExcessHourRate:     75000,  // 750.00 ETB
      })

      // Actual rental hours: 20 hours (4 excess hours over 16 agreed hours)
      excessFee := booking.CalculateExcessCharge(20)
      require.Equal(t, int64(300000), excessFee) // 3,000.00 ETB
  }
  ```

### 4.2 Integration Testing with `testcontainers-go`
Integration tests assert real PostgreSQL 18.4 + PostGIS and Valkey 8.0 interactions[cite: 25, 27, 32]:
```go
func TestPaymentRepository_ConfirmPayment_AtomicTx(t *testing.T) {
    ctx := context.Background()
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgis/postgis:18-3.4-alpine"),
        postgres.WithDatabase("dink_test"),
        testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections")),
    )
    require.NoError(t, err)
    defer func() { require.NoError(t, pgContainer.Terminate(ctx)) }()

    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    pool, err := pgxpool.New(ctx, connStr)
    require.NoError(t, err)
    defer pool.Close()

    // 1. Assert optimistic concurrency lock (version check)
    // 2. Assert COMMISSION_DEBIT row creation in owner_wallet_transactions
    // 3. Assert status transition to CONFIRMED
    // 4. Assert domain_events_outbox insertion within identical pgx.Tx
}
```

### 4.3 Soft-Delete & Scoping Verification (`SD-01-BE`)
Every database query compiled via `sqlc` must be validated to confirm soft-deleted records are filtered out[cite: 15, 35]:
```go
func TestFleetRepository_ExcludesSoftDeletedRecords(t *testing.T) {
    ctx := context.Background()
    repo := persistence.NewFleetRepository(testPool)

    fleetID := createTestFleet(t, ctx, testPool)
    err := repo.SoftDeleteFleet(ctx, fleetID)
    require.NoError(t, err)

    // Assert query using standard scope returns ErrNotFound
    _, err = repo.GetFleetByID(ctx, fleetID)
    require.ErrorIs(t, err, domain.ErrFleetNotFound)
}
```

---

## 5. Mobile Testing Standards (Flutter 3.44)

### 5.1 BLoC Unit Testing
All UI state mutations are tested via `bloc_test` asserting deterministic event-to-state streams[cite: 32]:
```dart
void main() {
  group('AuthBloc - OTP Verification', () {
    late MockAuthRepository mockAuthRepository;

    setUp(() {
      mockAuthRepository = MockAuthRepository();
    });

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthSuccess] when OTP verification succeeds',
      build: () {
        when(() => mockAuthRepository.verifyOTP(any(), any()))
            .thenAnswer((_) async => Right(AuthSession.dummy()));
        return AuthBloc(authRepository: mockAuthRepository);
      },
      act: (bloc) => bloc.add(const SubmitOTPEvent(otp: '123456')),
      expect: () => [
        const AuthLoadingState(),
        isA<AuthSuccessState>(),
      ],
      verify: (_) {
        verify(() => mockAuthRepository.verifyOTP(any(), '123456')).called(1);
      },
    );
  });
}
```

### 5.2 Widget & Design Token Drift Testing
- Widgets must be tested against the compiled tokens generated from `design-tokens.json`[cite: 33, 37, 38].
- Automated assertions confirm interactive tap targets meet the 44x44 dp minimum bounding size[cite: 33].
- Golden tests must run on both light and dark surface configurations asserting high contrast text token usage[cite: 33].

### 5.3 Offline Storage & Upload Queue Testing
- Local caching tests run against in-memory `Isar` instances[cite: 27, 32].
- Assert that network disconnect events trigger the offline banner, preserve local queue state, and process sequential retry pipelines upon reconnect[cite: 35].

---

## 6. Admin Web Testing Standards (Next.js 16)

### 6.1 Server Action & Envelope Parsing Testing
Server Actions must be tested using mocked backend adapters asserting correct RFC 9457 error parsing and token forwarding[cite: 24, 32]:
```typescript
describe('approveOwnerAction', () => {
  it('correctly maps 403 Forbidden to permission problem detail', async () => {
    const mockApiClient = jest.spyOn(apiClient, 'apiFetch').mockRejectedValue(
      new ProblemDetailsError({
        type: '[https://api.dinkmachinery.com/errors/forbidden](https://api.dinkmachinery.com/errors/forbidden)',
        status: 403,
        code: 'INSUFFICIENT_PERMISSIONS',
        detail: 'User lacks owners:approve scope.'
      })
    );

    const result = await approveOwnerAction('usr_test123', 'Valid business license');
    expect(result.success).toBe(false);
    expect(result.error?.code).toBe('INSUFFICIENT_PERMISSIONS');
  });
});
```

### 6.2 Accessibility Automated Scans (`axe-core`)
- Every page rendered via React Server Components must pass automated `axe-core` audits during PR builds[cite: 24].
- Any **Critical** or **Serious** violation immediately fails the build[cite: 24].

---

## 7. Multi-Country Parameterization Test Suite (`ADR-015`)

Tests must assert dynamic behavior against the `country_configs` table without hardcoded country assumptions[cite: 35]:

```go
func TestCountryConfig_DynamicValidation(t *testing.T) {
    cases := []struct {
        countryCode string
        phoneInput  string
        expectedOK  bool
    }{
        {"ET", "+251911234567", true},
        {"ET", "+251711234567", true},
        {"ET", "+254711234567", false}, // Invalid prefix for Ethiopia
        {"KE", "+254712345678", true},
    }

    for _, tc := range cases {
        cfg := getCountryConfig(t, tc.countryCode)
        ok := regexp.MustCompile(cfg.PhoneValidationRegex).MatchString(tc.phoneInput)
        assert.Equal(t, tc.expectedOK, ok, "Validation failed for country %s", tc.countryCode)
    }
}
```

- **SLA Calculation Test:** Assert that SLA calculators exclude both weekends and country-specific rows in `public_holidays` for the active `country_code`[cite: 35].
- **Tax Model Test:** Assert that VAT/Tax is suppressed (0%) when `owner.vat_registered == false`, and that accepted damages are excluded from the taxable base when `tax_applies_to_damages == false`[cite: 35].

---

## 8. Quality Gate Matrix & Execution Commands

### 8.1 Gate Pipeline Criteria

| Gate | Stage | Verification Commands | Blocking Rules |
|---|---|---|---|
| **G0** | Local Pre-Commit | `lefthook run pre-commit` | Any linter warning, unstaged `CONTEXT.md`, missing test file[cite: 24]. |
| **G1** | Pull Request | `gh workflow run ci.yml` | Coverage drop, contract schema mismatch, static scan finding[cite: 24]. |
| **G2** | Pre-Merge | `bash scripts/validate-decision-links.sh` | Missing ADR/Partition link, non-linear commit history[cite: 24]. |
| **G3** | Release / UAT | `k6 run tests/load/booking_soak.js` | Latency p95 > 1.0s, p99 > 2.0s, error rate >= 1%[cite: 24]. |

### 8.2 Standard Test Commands

```bash
# Backend Test Execution (Unit & Integration)
cd backend
go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
golangci-lint run ./...
govulncheck ./...

# Mobile Test Execution (Flutter)
cd mobile
dart run tool/generate_tokens.dart
flutter test --coverage
flutter analyze --fatal-warnings

# Admin Web Execution (Next.js)
cd admin
npm run lint
npm run test:ci
npx axe-core-cli http://localhost:3000/admin/dashboard

# Contract Specification Linter
npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml
```