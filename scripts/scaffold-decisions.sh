#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# scripts/scaffold-decisions.sh
#
# Generates standardized architectural decision records for Dink Machinery:
#   - Tier 1 (L1 - System-wide Architecture Decision Records): docs/decisions/ADR-NNN-<slug>.md
#   - Tier 2 (L2 - Partition Decision Records): docs/decisions/partitions/<partition-id>.md
#
# Enforces the 5-point qualification test and standardized sections per AGENTS.md.
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DECISIONS_DIR="${ROOT_DIR}/docs/decisions"
PARTITIONS_DIR="${DECISIONS_DIR}/partitions"

mkdir -p "${DECISIONS_DIR}"
mkdir -p "${PARTITIONS_DIR}"

usage() {
  cat <<EOF
Usage:
  $(basename "$0") adr <slug> [title]
  $(basename "$0") partition <partition-id> [decision-slug] [title]

Examples:
  $(basename "$0") adr dynamic-font-delivery "Dynamic Ethiopic and Arabic Font Delivery via S3"
  $(basename "$0") partition 05-05 booking-hold-atomicity "In-Transaction Booking Hold and Ledger Atomicity"
EOF
  exit 1
}

if [ $# -lt 2 ]; then
  usage
fi

COMMAND="$1"

case "${COMMAND}" in
  adr|l1)
    SLUG="${2}"
    TITLE="${3:-$(echo "${SLUG}" | tr '-' ' ' | awk '{for(i=1;i<=NF;i++)sub(/./,toupper(substr($i,1,1)),$i)}1')}"

    # Calculate next sequential ADR number
    EXISTING_ADRS=$(find "${DECISIONS_DIR}" -maxdepth 1 -name "ADR-[0-9][0-9][0-9]-*.md" | sort || true)
    if [ -z "${EXISTING_ADRS}" ]; then
      NEXT_NUM=1
    else
      LAST_NUM=$(echo "${EXISTING_ADRS}" | tail -n 1 | sed -E 's/.*ADR-([0-9]+)-.*/\1/' | sed 's/^0*//')
      NEXT_NUM=$((LAST_NUM + 1))
    fi

    ADR_ID=$(printf "ADR-%03d" "${NEXT_NUM}")
    FILE_NAME="${DECISIONS_DIR}/${ADR_ID}-${SLUG}.md"

    if [ -f "${FILE_NAME}" ]; then
      echo "ERROR: Target decision record already exists: ${FILE_NAME}"
      exit 1
    fi

    CURRENT_DATE=$(date -u +"%Y-%m-%d")

    cat <<EOF > "${FILE_NAME}"
# ${ADR_ID}: ${TITLE}

- **Status:** Proposed
- **Date:** ${CURRENT_DATE}
- **Decision Owner:** Dink Machinery Engineering Team
- **Affected Partitions:** Cross-cutting / Platform
- **Related Requirements:** [REQ-XX-YY]

---

## 1. Context & Problem Statement

<!-- Describe the business and technical context, architectural challenge, and why a platform-level decision is required. -->

---

## 2. The 5-Point Qualification Evaluation

- [ ] **Alternative Test:** Multiple viable architectural approaches evaluated.
- [ ] **Surprise Test:** Choice prevents non-obvious confusion for future engineers.
- [ ] **Constraint Test:** Driven by hard constraints (performance, security ASVS L3, compliance).
- [ ] **Irreversibility Test:** High refactor/migration cost (>1 engineer-day or public schema change).
- [ ] **Scope Test:** Impacts multiple requirement partitions or subsystems.

---

## 3. Considered Options

### Option 1: [Option Name]
- **Summary:** <!-- Brief description -->
- **Pros:**
  - 
- **Cons:**
  - 

### Option 2: [Option Name]
- **Summary:** <!-- Brief description -->
- **Pros:**
  - 
- **Cons:**
  - 

---

## 4. Chosen Solution & Architectural Rationale

<!-- State the chosen option and provide the definitive justification. -->

---

## 5. Rejected Alternatives & Technical Justification

- **Option X:** Rejected because...
- **Option Y:** Rejected because...

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - 
- **Negative / Technical Debt:**
  - 
- **Operational / Observability Impacts:**
  - 

---

## 7. Verification Strategy & Invariant Tests

<!-- List automated test suites, benchmarks, or invariant validations enforcing this decision. -->
- **Test Suite:** \`tests/...\`
- **Assertion:** 

---

## 8. Revisit Conditions

<!-- Under what operational metrics, scale, or business pivots should this decision be reconsidered? -->
- 
EOF

    echo "Successfully scaffolded L1 ADR: ${FILE_NAME}"
    ;;

  partition|l2)
    PARTITION_ID="${2}"
    DECISION_SLUG="${3:-decision}"
    TITLE="${4:-$(echo "${DECISION_SLUG}" | tr '-' ' ' | awk '{for(i=1;i<=NF;i++)sub(/./,toupper(substr($i,1,1)),$i)}1')}"

    PARTITION_FILE="${PARTITIONS_DIR}/${PARTITION_ID}.md"
    CURRENT_DATE=$(date -u +"%Y-%m-%d")

    if [ ! -f "${PARTITION_FILE}" ]; then
      cat <<EOF > "${PARTITION_FILE}"
# Partition Decisions: ${PARTITION_ID}

This document contains localized L2 architectural decisions for partition \`${PARTITION_ID}\`.

---
EOF
      NEXT_DECISION_NUM=1
    else
      EXISTING_IDS=$(grep -E "^## D-${PARTITION_ID}-[0-9]+" "${PARTITION_FILE}" || true)
      if [ -z "${EXISTING_IDS}" ]; then
        NEXT_DECISION_NUM=1
      else
        LAST_DECISION_NUM=$(echo "${EXISTING_IDS}" | tail -n 1 | sed -E "s/## D-${PARTITION_ID}-([0-9]+).*/\1/" | sed 's/^0*//')
        NEXT_DECISION_NUM=$((LAST_DECISION_NUM + 1))
      fi
    fi

    DECISION_ID=$(printf "D-%s-%03d" "${PARTITION_ID}" "${NEXT_DECISION_NUM}")

    cat <<EOF >> "${PARTITION_FILE}"

## ${DECISION_ID}: ${TITLE}

- **Status:** Accepted
- **Date:** ${CURRENT_DATE}
- **Requirement Target:** [REQ-${PARTITION_ID}-XX]
- **Subsystem:** [backend | mobile | admin | database]

### 1. Problem Statement & Context
<!-- Describe the partition-specific technical trade-off. -->

### 2. Considered Alternatives
1. **Option A:** 
2. **Option B:** 

### 3. Chosen Solution & Technical Rationale
<!-- Detail why this solution was selected. -->

### 4. Rejected Alternatives
- **Option B:** Rejected because...

### 5. Consequences & Invariants
- **Invariants Maintained:** 
- **Operational Trade-offs:** 

### 6. Verification
- **Test File:** \`..._test.go\`
- **Assertion:** 

### 7. Revisit Trigger
- Revisit if...

---
EOF

    echo "Successfully recorded L2 Decision ${DECISION_ID} in: ${PARTITION_FILE}"
    ;;

  *)
    usage
    ;;
esac