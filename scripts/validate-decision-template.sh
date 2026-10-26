#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# scripts/validate-decision-template.sh
#
# Audits all architectural decision records (L1 ADRs and L2 Partition Decisions)
# against the mandatory structural specification defined in AGENTS.md:
#   1. Every L1 ADR must contain all 8 standardized sections, metadata headers,
#      completed 5-point evaluation, and minimum 2 considered alternatives.
#   2. Every L2 Partition decision must contain all 7 standardized sub-sections,
#      metadata attributes, and minimum 2 considered alternatives.
#   3. Flags unpopulated template placeholders or HTML comments.
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DECISIONS_DIR="${ROOT_DIR}/docs/decisions"
PARTITIONS_DIR="${DECISIONS_DIR}/partitions"

ERRORS=0

echo "==> Auditing Architectural Decision Template Conformance..."

if [ ! -d "${DECISIONS_DIR}" ]; then
  echo "ERROR: Decisions directory does not exist at ${DECISIONS_DIR}"
  exit 1
fi

# ------------------------------------------------------------------------------
# 1. Audit L1 ADR Records (docs/decisions/ADR-NNN-*.md)
# ------------------------------------------------------------------------------
ADR_FILES=$(find "${DECISIONS_DIR}" -maxdepth 1 -name "ADR-[0-9][0-9][0-9]-*.md" | sort || true)

MANDATORY_ADR_SECTIONS=(
  "## 1. Context & Problem Statement"
  "## 2. The 5-Point Qualification Evaluation"
  "## 3. Considered Options"
  "## 4. Chosen Solution & Architectural Rationale"
  "## 5. Rejected Alternatives & Technical Justification"
  "## 6. Consequences & Operational Trade-offs"
  "## 7. Verification Strategy & Invariant Tests"
  "## 8. Revisit Conditions"
)

MANDATORY_ADR_METADATA=(
  "- **Status:**"
  "- **Date:**"
  "- **Decision Owner:**"
  "- **Affected Partitions:**"
  "- **Related Requirements:**"
)

for file in ${ADR_FILES}; do
  [ -f "${file}" ] || continue
  filename="$(basename "${file}")"
  echo "--> Validating L1 ADR: ${filename}"

  # Check metadata fields
  for meta in "${MANDATORY_ADR_METADATA[@]}"; do
    if ! grep -qF "${meta}" "${file}"; then
      echo "  [MISSING METADATA] ${filename} is missing metadata: '${meta}'"
      ERRORS=$((ERRORS + 1))
    fi
  done

  # Check top-level mandatory sections
  for section in "${MANDATORY_ADR_SECTIONS[@]}"; do
    if ! grep -qF "${section}" "${file}"; then
      echo "  [MISSING SECTION] ${filename} is missing section: '${section}'"
      ERRORS=$((ERRORS + 1))
    fi
  done

  # Validate considered options (minimum 2 options required)
  OPTION_COUNT=$(grep -cE '### Option [0-9]+:' "${file}" || true)
  if [ "${OPTION_COUNT}" -lt 2 ]; then
    echo "  [INSUFFICIENT OPTIONS] ${filename} defines ${OPTION_COUNT} options. Minimum 2 required under 'Considered Options'."
    ERRORS=$((ERRORS + 1))
  fi

  # Check for raw placeholder HTML comments
  if grep -qE '<!--.*-->' "${file}"; then
    echo "  [UNRESOLVED PLACEHOLDER] ${filename} contains raw HTML template comment placeholders."
    ERRORS=$((ERRORS + 1))
  fi
done

# ------------------------------------------------------------------------------
# 2. Audit L2 Partition Decisions (docs/decisions/partitions/*.md)
# ------------------------------------------------------------------------------
MANDATORY_L2_SECTIONS=(
  "### 1. Problem Statement & Context"
  "### 2. Considered Alternatives"
  "### 3. Chosen Solution & Technical Rationale"
  "### 4. Rejected Alternatives"
  "### 5. Consequences & Invariants"
  "### 6. Verification"
  "### 7. Revisit Trigger"
)

MANDATORY_L2_METADATA=(
  "- **Status:**"
  "- **Date:**"
  "- **Requirement Target:**"
  "- **Subsystem:**"
)

if [ -d "${PARTITIONS_DIR}" ]; then
  PARTITION_FILES=$(find "${PARTITIONS_DIR}" -maxdepth 1 -name "*.md" | sort || true)

  for file in ${PARTITION_FILES}; do
    [ -f "${file}" ] || continue
    filename="$(basename "${file}")"

    # Find line numbers of all decision headers (## D-XX-YY-ZZZ)
    DECISION_LINES=$(grep -nE '^##[[:space:]]+D-[0-9]{2}-[0-9]{2}-[0-9]{3}' "${file}" | cut -d: -f1 || true)

    if [ -z "${DECISION_LINES}" ]; then
      continue
    fi

    TOTAL_LINES=$(wc -l < "${file}")
    DEC_ARRAY=(${DECISION_LINES})

    for i in "${!DEC_ARRAY[@]}"; do
      START_LINE="${DEC_ARRAY[$i]}"
      if [ $((i + 1)) -lt "${#DEC_ARRAY[@]}" ]; then
        END_LINE=$((DEC_ARRAY[$((i + 1))] - 1))
      else
        END_LINE="${TOTAL_LINES}"
      fi

      DEC_HEADER=$(sed -n "${START_LINE}p" "${file}")
      DEC_ID=$(echo "${DEC_HEADER}" | grep -oE 'D-[0-9]{2}-[0-9]{2}-[0-9]{3}')
      echo "--> Validating L2 Decision: ${DEC_ID} in ${filename}"

      DEC_CONTENT=$(sed -n "${START_LINE},${END_LINE}p" "${file}")

      # Verify metadata
      for meta in "${MANDATORY_L2_METADATA[@]}"; do
        if ! echo "${DEC_CONTENT}" | grep -qF "${meta}"; then
          echo "  [MISSING METADATA] ${filename} (${DEC_ID}) is missing metadata: '${meta}'"
          ERRORS=$((ERRORS + 1))
        fi
      done

      # Verify sections
      for section in "${MANDATORY_L2_SECTIONS[@]}"; do
        if ! echo "${DEC_CONTENT}" | grep -qF "${section}"; then
          echo "  [MISSING SECTION] ${filename} (${DEC_ID}) is missing section: '${section}'"
          ERRORS=$((ERRORS + 1))
        fi
      done

      # Verify at least 2 considered alternatives
      ALTERNATIVES_COUNT=$(echo "${DEC_CONTENT}" | grep -cE '^[0-9]+\.[[:space:]]+\*\*Option[[:space:]]+[A-Z0-9]+:\*\*' || true)
      if [ "${ALTERNATIVES_COUNT}" -lt 2 ]; then
        echo "  [INSUFFICIENT ALTERNATIVES] ${filename} (${DEC_ID}) specifies ${ALTERNATIVES_COUNT} alternatives. Minimum 2 required."
        ERRORS=$((ERRORS + 1))
      fi

      # Check for raw placeholder HTML comments
      if echo "${DEC_CONTENT}" | grep -qE '<!--.*-->'; then
        echo "  [UNRESOLVED PLACEHOLDER] ${filename} (${DEC_ID}) contains raw HTML template comment placeholders."
        ERRORS=$((ERRORS + 1))
      fi
    done
  done
fi

# ------------------------------------------------------------------------------
# 3. Final Status Report
# ------------------------------------------------------------------------------
if [ "${ERRORS}" -gt 0 ]; then
  echo "FAILED: Found ${ERRORS} decision template conformance error(s)."
  exit 1
fi

echo "SUCCESS: All architectural decision templates strictly adhere to governance specifications."
exit 0