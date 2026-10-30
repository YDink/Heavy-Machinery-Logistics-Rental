# API Contract Testing & Schema Verification Specification — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Specification:** Contract Testing, Wire Conformance & Drift Elimination Protocol  
**Compliance Standard:** ISO/IEC/IEEE 29119 (Software Testing), OpenAPI 3.1.0, RFC 9457 (Problem Details), OWASP ASVS Level 3  
**Target Subsystems:** Go 1.27.3 Backend (Provider), Flutter 3.47.9 Mobile (Consumer), Next.js 16 Admin (Consumer)

---

## 1. Executive Summary & Core Contract Invariants

In **Dink Machinery**, `contracts/openapi.yaml` represents an inviolable legal contract between backend provider implementations and client consumer applications. Under contract-first governance, API divergence, schema guesswork, and manual type mapping are classified as critical system defects.

Contract testing verifies that:
1. **The Provider (Go Backend)** never emits status codes, response headers, or JSON payloads that diverge from `contracts/openapi.yaml`.
2. **The Consumers (Mobile & Admin)** never submit requests that violate boundary schemas and correctly handle all defined 2xx success payloads and RFC 9457 error envelopes.
3. **The Schema Definition** evolves in a backwards-compatible manner, rejecting unapproved breaking changes through automated CI gates.

### Platform Wire Invariants
- **Universal Envelopes Mandatory:** All 2xx endpoints emit `SuccessEnvelope` (`success: true`, `data: {}`, `meta: {}`); all 4xx/5xx endpoints emit `ErrorEnvelope` with RFC 9457 problem fields (`success: false`, `error: {}`, `meta: {}`). Naked objects or flat arrays trigger contract test failures.
- **Integer Minor-Unit Currency:** Monetary properties must resolve to the composite `MoneyAmount` schema (`amountMinor` int64 + `currency` string). Any schema or wire property containing floating-point representations fails contract verification.
- **Idempotency Key Verification:** Mutating endpoints (`POST`, `PUT`, `PATCH`) must validate the presence, format, and replay semantics of the `Idempotency-Key` header.
- **Multi-Country Parameterization (`LOC-01-BE`):** Jurisdictional parameters must accept and parse `X-Country-Code` headers, asserting fallback defaults against `country_configs`.
- **Zero-PII Wire Invariant:** Unmasked phone numbers, National IDs, TINs, and bank accounts must never appear in response envelopes, query strings, or error details.

---

## 2. Multi-Tier Contract Testing Pyramid

```
                                [Level 4: E2E Integration]
                               Full stack live HTTP journeys
                                      (Staging / CI)
                                            │
                             [Level 3: Consumer Wire Tests]
                           Client DTO deserialization & state
                             matching against Stoplight Prism
                                            │
                             [Level 2: Provider Wire Tests]
                           Kin-OpenAPI validation middleware
                           asserting live Go Chi HTTP output
                                            │
                             [Level 1: Static Linter & Diff]
                           Spectral CLI (.spectral.yaml) & oasdiff
                               breaking change verification
```

| Level | Scope | Primary Tooling | Execution Trigger | Gate Criteria |
|---|---|---|---|---|
| **L1: Static Lint & Diff** | Schema structure & compatibility | Spectral CLI, `oasdiff` | Pre-commit (`lefthook`), PR CI | 0 errors, 0 warnings, 0 unapproved breaking changes |
| **L2: Provider Wire** | Go Chi handler output schema validation | `kin-openapi`, `testcontainers-go` | Backend Unit/Integration CI | 100% route coverage matching `openapi.yaml` schemas |
| **L3: Consumer Wire** | Client model deserialization & error handling | Stoplight Prism, Dart test, Vitest | Mobile/Admin test suites | Successful round-trip parse of all schema fixtures |
| **L4: E2E Wire** | Live multi-tenant cross-service flows | Playwright (Admin), Patrol (Mobile) | Staging Deployment Gate | Zero wire schema violations on end-to-end paths |

---

## 3. Level 1: Static Contract Linting & Governance Rules

The canonical contract `contracts/openapi.yaml` is validated via **Spectral CLI** using `.spectral.yaml`. Every schema definition and operation must satisfy custom rules before any code generation is permitted.

### 3.1. Enforced Spectral Rule Set (`.spectral.yaml`)
```yaml
extends: ["spectral:oas", "spectral:asyncapi"]

rules:
  # Structural Invariants
  oas3-schema: true
  operation-operationId: true
  operation-tags: true
  operation-description: true

  # Path Casing
  paths-kebab-case:
    description: "All path segments must use lowercase kebab-case, standard path parameters, and optional RPC action suffixes."
    severity: error
    given: "$.paths[*]~"
    then:
      function: pattern
      functionOptions:
        match: "^(/|[a-zA-Z0-9-~/_{}:]+)*$"

  # Universal Success Envelope Verification
  universal-success-envelope:
    description: "2xx responses must reference an envelope containing success, data, error, and meta properties."
    severity: error
    given: "$.paths.*.*.responses[?(@property >= 200 && @property < 300)].content['application/json'].schema.properties"
    then:
      - field: success
        function: defined
      - field: data
        function: defined
      - field: error
        function: defined
      - field: meta
        function: defined

  # RFC 9457 Problem Details Error Envelope Verification
  rfc9457-error-envelope:
    description: "4xx and 5xx responses must reference an ErrorEnvelope containing success, data, error, and meta properties."
    severity: error
    given: "$.paths.*.*.responses[?(@property >= 400 && @property <= 599)].content['application/json'].schema.properties"
    then:
      - field: success
        function: defined
      - field: data
        function: defined
      - field: error
        function: defined
      - field: meta
        function: defined

  # Integer Minor-Unit Currency Invariant
  no-floating-point-money:
    description: "Financial properties must use 64-bit integer minor units (MoneyAmount). Floats are banned."
    severity: error
    given: "$..properties[?(@property.match(/(amount|price|cost|balance|rate|deposit|fee)/i))]"
    then:
      function: schema
      functionOptions:
        schema:
          not:
            type: number

  # Mandatory Idempotency Header on Mutating Operations
  idempotency-header-required:
    description: "State-mutating operations (POST, PUT, PATCH) must declare the Idempotency-Key header."
    severity: error
    given: "$.paths.*[post,put,patch].parameters[?(@.in == 'header' && @.name == 'Idempotency-Key')]"
    then:
      function: defined
```

### 3.2. Static Execution Drill
```bash
# Execute lint check locally
pnpm dlx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml
```

---

## 4. Breaking Change Detection & Semantic Versioning

Contract modifications must not break deployed mobile applications or administrative dashboards. Every PR touching `contracts/openapi.yaml` executes an automated breaking change analysis against `origin/main` using `oasdiff`.

### 4.1. Breaking Change Rules Matrix

| Operation Category | Permitted in Minor Version? | Rule Criteria |
|---|---|---|
| **Add New Endpoint** | Yes | New path or operation added with valid envelope schemas. |
| **Add Optional Request Property** | Yes | New request body field or query parameter marked `required: false`. |
| **Add Response Property** | Yes | New property added to existing 2xx `data` payload. |
| **Remove Endpoint** | **No (Breaking)** | Deletion of path or HTTP method requires an approved L1 ADR and 90-day sunset. |
| **Rename / Delete Field** | **No (Breaking)** | Renaming properties in request or response envelopes is prohibited. |
| **Add Mandatory Request Field** | **No (Breaking)** | Making an optional field `required: true` breaks active mobile clients. |
| **Change Data Type / Format** | **No (Breaking)** | e.g., Changing `amountMinor` from integer to string or changing UUID format. |
| **Alter Enum Values** | **No (Breaking)** | Removing an enum value or adding a value on request schemas without fallback. |

### 4.2. CI Breaking Change Audit Command
```bash
# Verify backward compatibility against baseline main branch
oasdiff breaking origin/main:contracts/openapi.yaml contracts/openapi.yaml \
  --fail-on-diff \
  --format text
```

---

## 5. Level 2: Provider Wire Verification (Go 1.27.3 Backend)

The Go modular monolith validates actual HTTP responses against `contracts/openapi.yaml` during integration testing. Using `kin-openapi`, handlers are exercised against real PostgreSQL (via `testcontainers-go`), asserting that status codes, response headers, and JSON bodies strictly adhere to the contract.

### 5.1. Automated Kin-OpenAPI Test Helper (`backend/internal/transport/http/contract_test.go`)
```go
package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"[github.com/getkin/kin-openapi/openapi3](https://github.com/getkin/kin-openapi/openapi3)"
	"[github.com/getkin/kin-openapi/openapi3filter](https://github.com/getkin/kin-openapi/openapi3filter)"
	"[github.com/getkin/kin-openapi/routers](https://github.com/getkin/kin-openapi/routers)"
	"[github.com/getkin/kin-openapi/routers/gorillamux](https://github.com/getkin/kin-openapi/routers/gorillamux)"
	"[github.com/stretchr/testify/require](https://github.com/stretchr/testify/require)"
)

// ContractValidator provides reusable wire assertion against openapi.yaml
type ContractValidator struct {
	doc    *openapi3.T
	router routers.Router
}

func NewContractValidator(t *testing.T, specPath string) *ContractValidator {
	t.Helper()
	loader := openapi3filter.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromFile(specPath)
	require.NoError(t, err, "Failed to parse openapi.yaml contract")

	err = doc.Validate(context.Background())
	require.NoError(t, err, "Contract failed self-validation")

	router, err := gorillamux.NewRouter(doc)
	require.NoError(t, err, "Failed to compile GorillaMux router from OpenAPI doc")

	return &ContractValidator{
		doc:    doc,
		router: router,
	}
}

func (cv *ContractValidator) AssertWireConformance(
	t *testing.T,
	req *http.Request,
	status int,
	responseHeaders http.Header,
	responseBody []byte,
) {
	t.Helper()

	// 1. Locate Route in OpenAPI contract
	route, pathParams, err := cv.router.FindRoute(req)
	require.NoError(t, err, "Route %s %s does not exist in openapi.yaml", req.Method, req.URL.Path)

	// 2. Validate Request Conformance
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			IncludeQueryParameters: true,
		},
	}
	err = openapi3filter.ValidateRequest(context.Background(), requestValidationInput)
	require.NoError(t, err, "Outgoing request violates openapi.yaml contract")

	// 3. Validate Response Conformance (Status, Headers, Envelope Schema)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 status,
		Header:                 responseHeaders,
		Body:                   responseBody,
		Options: &openapi3filter.Options{
			IncludeResponseStatus: true,
		},
	}
	err = openapi3filter.ValidateResponse(context.Background(), responseValidationInput)
	require.NoError(t, err, "Backend response violates openapi.yaml schema for %s %s (Status: %d)",
		req.Method, req.URL.Path, status)

	// 4. Assert Mandatory Envelope Invariants
	var envelope struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
		Error   json.RawMessage `json:"error"`
		Meta    struct {
			Version   string `json:"version"`
			RequestID string `json:"requestId"`
			Timestamp string `json:"timestamp"`
		} `json:"meta"`
	}

	err = json.Unmarshal(responseBody, &envelope)
	require.NoError(t, err, "Response body must unmarshal into standard envelope")
	require.NotEmpty(t, envelope.Meta.RequestID, "meta.requestId is mandatory")
	require.NotEmpty(t, envelope.Meta.Timestamp, "meta.timestamp is mandatory")

	if status >= 200 && status < 300 {
		require.True(t, envelope.Success, "2xx response must have success: true")
		require.Equal(t, "null", string(envelope.Error), "2xx response must have error: null")
	} else {
		require.False(t, envelope.Success, "4xx/5xx response must have success: false")
		require.Equal(t, "null", string(envelope.Data), "4xx/5xx response must have data: null")
	}
}
```

### 5.2. Concrete Provider Test Execution Example
```go
func TestBookingHoldCreation_ContractConformance(t *testing.T) {
	validator := NewContractValidator(t, "../../../../contracts/openapi.yaml")

	// Setup Handler & In-Memory or Testcontainer Infrastructure
	handler := setupTestBookingHandler(t)

	// Construct Valid Request Body
	requestBody := []byte(`{
		"fleetId": "7f000001-8b22-1d54-818b-22b0f4910001",
		"startDate": "2026-11-01",
		"endDate": "2026-11-05",
		"siteLocationAddress": "Bole Lemi Industrial Zone Phase 2"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d")
	req.Header.Set("X-Country-Code", "ET")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1Ni...")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Assert strict wire conformance against OpenAPI specification
	validator.AssertWireConformance(t, req, w.Code, w.Header(), w.Body.Bytes())
	require.Equal(t, http.StatusCreated, w.Code)
}
```

---

## 6. Level 3: Consumer Wire Verification (Mobile & Admin)

Consumers verify that their client models, deserialization adapters, and state machines correctly parse canonical payloads and map RFC 9457 errors without crashing.

### 6.1. Next.js 16 Admin Client Contract Test (Vitest + Zod)
The admin web console validates responses received from the mock server against Zod schemas derived from `src/types/api.generated.ts`:

```typescript
import { describe, it, expect } from 'vitest';
import { z } from 'zod';
import { fetchFleetById } from '@/lib/api/fleet-service';

// Runtime Zod Schema Matching contracts/openapi.yaml
const MoneyAmountSchema = z.object({
  amountMinor: z.number().int(),
  currency: z.string().length(3),
});

const FleetAssetSchema = z.object({
  id: z.string().uuid(),
  ownerId: z.string().uuid(),
  categoryId: z.string().uuid(),
  make: z.string().min(1),
  model: z.string().min(1),
  manufacturingYear: z.number().int(),
  registrationNumber: z.string(),
  dailyRate: MoneyAmountSchema,
  status: z.enum(['AVAILABLE', 'RENTED', 'UNDER_MAINTENANCE', 'DECOMMISSIONED']),
  version: z.number().int().positive(),
});

const SuccessResponseSchema = z.object({
  success: z.literal(true),
  data: FleetAssetSchema,
  error: z.null(),
  meta: z.object({
    version: z.string(),
    requestId: z.string(),
    timestamp: z.string().datetime(),
  }),
});

describe('Admin Fleet Client Contract Verification', () => {
  it('correctly deserializes 200 OK fleet payload matching OpenAPI schema', async () => {
    // Queries local Stoplight Prism mock server running on port 4010
    const response = await fetchFleetById('7f000001-8b22-1d54-818b-22b0f4910001', {
      baseUrl: 'http://localhost:4010/api/v1',
    });

    const validationResult = SuccessResponseSchema.safeParse(response);
    expect(validationResult.success).toBe(true);
    if (!validationResult.success) {
      console.error(validationResult.error.format());
    }
  });
});
```

### 6.2. Flutter 3.47.9 Mobile Consumer Test (Dart / Dio)
The mobile app verifies that the Dio HTTP client interceptor properly deserializes an RFC 9457 error into a typed `ProblemDetails` model:

```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:dio/dio.dart';
import 'package:dink_mobile/core/network/problem_details.dart';
import 'package:dink_mobile/core/network/problem_details_interceptor.dart';

void main() {
  group('ProblemDetailsInterceptor Contract Conformance', () {
    late Dio dio;

    setUp(() {
      dio = Dio(BaseOptions(baseUrl: 'http://localhost:4010/api/v1'));
      dio.interceptors.add(ProblemDetailsInterceptor());
    });

    test('deserializes 409 Conflict RFC 9457 Problem Details payload', () async {
      try {
        // Trigger synthetic 409 conflict on Prism mock server
        await dio.post(
          '/bookings/7f000001-8b22-1d54-818b-22b0f4910001:confirm-payment',
          data: {'paymentTransactionId': '7f000001-8b22-1d54-818b-22b0f4910002', 'version': 1},
          options: Options(headers: {
            'Prefer': 'code=409',
            'Idempotency-Key': '9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d',
          }),
        );
        fail('Should have thrown DioException');
      } on DioException catch (e) {
        expect(e.error, isA<ProblemDetails>());
        final problem = e.error as ProblemDetails;

        expect(problem.status, 409);
        expect(problem.code, 'BOOKING_INVALID_STATE_TRANSITION');
        expect(problem.type, contains('errors/booking-invalid-state'));
        expect(problem.detail, isNotEmpty);
      }
    });
  });
}
```

---

## 7. Stoplight Prism Mock Server & Dynamic Simulation

During local development and CI execution, **Stoplight Prism** serves as an active contract validator and mock server, simulating the backend modular monolith.

### 7.1. Starting the Prism Contract Mock Server
```bash
# Start Prism with strict schema enforcement enabled
pnpm dlx @stoplight/prism-cli mock contracts/openapi.yaml \
  --port 4010 \
  --errors \
  --cors
```

### 7.2. Simulation Flags & Dynamic Testing
Prism inspects incoming consumer requests and validates them against `contracts/openapi.yaml`:
- **Strict Validation (`--errors`):** If a client sends a request missing a mandatory property, an invalid UUID, or a floating-point number in money fields, Prism immediately returns `422 Unprocessable Entity` with schema validation error details.
- **Dynamic Status Forcing:** Consumers can test error handling paths by injecting the `Prefer` header:
  - `Prefer: code=400` $\rightarrow$ Simulates Bad Request.
  - `Prefer: code=401` $\rightarrow$ Simulates Expired Token.
  - `Prefer: code=409` $\rightarrow$ Simulates OCC Version Conflict or Hold Expired.
  - `Prefer: code=429` $\rightarrow$ Simulates Rate Limit Exceeded with `Retry-After`.

---

## 8. Continuous Integration Pipeline & Quality Ratchet

Contract validation is enforced in GitHub Actions (`.github/workflows/ci-contract.yml`). Every PR must pass all stages before merge clearance.

```mermaid
flowchart LR
  PR[PR Trigger] --> L1[Spectral Lint]
  L1 --> L2[Oasdiff Breaking Audit]
  L2 --> L3[Regenerate Code]
  L3 --> L4[Assert Zero Git Drift]
  L4 --> L5[Go Provider Tests]
  L5 --> L6[Consumer Prism Tests]
  L6 --> Pass[Gate Clearance]
```

### 8.1. CI Workflow Specification
```yaml
name: Canonical API Contract Verification

on:
  pull_request:
    paths:
      - 'contracts/**'
      - '.spectral.yaml'
      - 'backend/internal/transport/http/**'
      - 'admin/src/lib/api/**'
      - 'mobile/lib/core/network/**'

jobs:
  contract-lint-and-diff:
    name: Lint & Breaking Change Audit
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: 26.13.0

      - name: Install Linting Tooling
        run: pnpm add -g @stoplight/spectral-cli @stoplight/prism-cli

      - name: Execute Spectral Schema Lint
        run: spectral lint contracts/openapi.yaml --ruleset .spectral.yaml

      - name: Install oasdiff
        run: |
          curl -fsSL [https://raw.githubusercontent.com/tufin/oasdiff/main/install.sh](https://raw.githubusercontent.com/tufin/oasdiff/main/install.sh) | sh
          sudo mv ./bin/oasdiff /usr/local/bin/

      - name: Check for Breaking API Changes
        run: |
          oasdiff breaking origin/main:contracts/openapi.yaml contracts/openapi.yaml \
            --fail-on-diff \
            --format text

  drift-verification:
    name: Assert Zero Code-Gen Drift
    needs: contract-lint-and-diff
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Toolchains
        uses: actions/setup-go@v5
        with:
          go-version: '1.27.3'

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: 26.13.0

      - name: Install Code Generators
        run: |
          go install [github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1](https://github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1)
          pnpm add -g openapi-typescript

      - name: Execute Multi-Target Code Generation
        run: bash scripts/generate-contracts.sh

      - name: Assert Zero Drift
        run: |
          if ! git diff --exit-code backend/ internal/ admin/src/types/ mobile/packages/; then
            echo "CRITICAL: Contract drift detected! Committed code does not match contracts/openapi.yaml."
            exit 1
          fi

  provider-wire-test:
    name: Go Provider Wire Conformance
    needs: drift-verification
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.27.3'
      - name: Run Backend Contract Tests
        run: |
          cd backend
          go test -v -race ./internal/transport/http/... -run "ContractConformance"
```

---

## 9. Contract Test Sign-Off Checklist

Before any PR introducing or modifying API operations can be merged into `main`, affirm that all items are verified:

- [ ] **Spectral Lint Clean:** `spectral lint contracts/openapi.yaml` outputs 0 errors, 0 warnings.
- [ ] **No Floating-Point Money:** All monetary attributes verified as integer minor units (`MoneyAmount`).
- [ ] **Envelopes Enforced:** Every 2xx schema wrapped in `SuccessResponse`; every 4xx/5xx wrapped in `ErrorEnvelope`.
- [ ] **Idempotency Annotations:** Mutating operations require `Idempotency-Key` header parameter.
- [ ] **Multi-Country Support:** `X-Country-Code` present on localized endpoints.
- [ ] **Zero Breaking Changes:** `oasdiff breaking origin/main:contracts/openapi.yaml contracts/openapi.yaml` reports 0 breaking modifications.
- [ ] **Zero Code-Gen Drift:** Code generators executed; `git diff --exit-code` confirms generated code matches the contract.
- [ ] **Provider Wire Tests:** Go handler tests validate response envelopes against `contracts/openapi.yaml` via `kin-openapi`.
- [ ] **Consumer Parsing Verified:** Mobile Dio interceptors and Admin Zod schemas deserialize test fixtures without exceptions.
- [ ] **Error Catalog Registry:** Any new domain error codes registered in `ERROR_CATALOG.md`.
- [ ] **Zero PII Exposure (ASVS L3):** Telephone numbers, National IDs, TINs, and credentials omitted from query parameters, URI segments, and error details.