#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# scripts/validate-decision-links.sh
#
# Audits architectural decision records and references across Dink Machinery:
#   1. Resolves all defined L1 ADRs (docs/decisions/ADR-NNN-*.md).
#   2. Resolves all defined L2 Partition Decisions (docs/decisions/partitions/*.md).
#   3. Scans documentation and requirements for referenced decision IDs to guarantee
#      zero dangling references.
#   4. Asserts bidirectional linkage between ADRs and corresponding requirement partitions.
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCS_DIR="${ROOT_DIR}/docs"
DECISIONS_DIR="${DOCS_DIR}/decisions"
PARTITIONS_DIR="${DECISIONS_DIR}/partitions"
REQUIREMENTS_DIR="${DOCS_DIR}/shared/requirements"

ERRORS=0

echo "==> Auditing Architectural Decision Links & References..."

if [ ! -d "${DECISIONS_DIR}" ]; then
  echo "ERROR: Decisions directory does not exist at ${DECISIONS_DIR}"
  exit 1
fi

# 1. Index All Defined ADR IDs
declare -A DEFINED_ADRS
for adr_file in "${DECISIONS_DIR}"/ADR-[0-9][0-9][0-9]-*.md; do
  [ -f "${adr_file}" ] || continue
  filename="$(basename "${adr_file}")"
  adr_id="$(echo "${filename}" | sed -E 's/^(ADR-[0-9]{3})-.*/\1/')"
  DEFINED_ADRS["${adr_id}"]="${adr_file}"
done

echo "Found ${#DEFINED_ADRS[@]} registered L1 ADRs."

# 2. Index All Defined L2 Partition Decision IDs
declare -A DEFINED_L2
if [ -d "${PARTITIONS_DIR}" ]; then
  for part_file in "${PARTITIONS_DIR}"/*.md; do
    [ -f "${part_file}" ] || continue
    while IFS= read -r line; do
      l2_id="$(echo "${line}" | sed -E 's/^##[[:space:]]+(D-[0-9]{2}-[0-9]{2}-[0-9]{3}):?.*/\1/')"
      if [ -n "${l2_id}" ]; then
        DEFINED_L2["${l2_id}"]="${part_file}"
      fi
    done < <(grep -E '^##[[:space:]]+D-[0-9]{2}-[0-9]{2}-[0-9]{3}' "${part_file}" || true)
  done
fi

COUNT_L2=0
for k in "${!DEFINED_L2[@]}"; do
  COUNT_L2=$((COUNT_L2 + 1))
done
echo "Found ${COUNT_L2} registered L2 Partition Decisions."

# 3. Validate Referenced ADRs Across Repository
echo "--> Verifying all ADR references point to existing records..."
while IFS=: read -r file line match; do
  # Extract clean ADR token
  clean_adr="$(echo "${match}" | grep -o -E 'ADR-[0-9]{3}' | head -n 1)"
  if [ -n "${clean_adr}" ]; then
    if [ -z "${DEFINED_ADRS[${clean_adr}]:-}" ]; then
      echo "  [BROKEN LINK] ${file}:${line} references ${clean_adr}, but no corresponding file exists in ${DECISIONS_DIR}/"
      ERRORS=$((ERRORS + 1))
    fi
  fi
done < <(grep -rnE 'ADR-[0-9]{3}' "${DOCS_DIR}" "${ROOT_DIR}/contracts" --exclude="*.swp" 2>/dev/null || true)

# 4. Validate Referenced L2 Decisions Across Repository
echo "--> Verifying all L2 decision references point to valid entries..."
while IFS=: read -r file line match; do
  clean_l2="$(echo "${match}" | grep -o -E 'D-[0-9]{2}-[0-9]{2}-[0-9]{3}' | head -n 1)"
  if [ -n "${clean_l2}" ]; then
    if [ -z "${DEFINED_L2[${clean_l2}]:-}" ]; then
      echo "  [BROKEN LINK] ${file}:${line} references ${clean_l2}, but no entry exists in ${PARTITIONS_DIR}/"
      ERRORS=$((ERRORS + 1))
    fi
  fi
done < <(grep -rnE 'D-[0-9]{2}-[0-9]{2}-[0-9]{3}' "${DOCS_DIR}" "${ROOT_DIR}/contracts" --exclude="*.swp" 2>/dev/null || true)

# 5. Verify Bidirectional Linkage Between ADRs and Affected Partitions
echo "--> Verifying reciprocal links between ADRs and partition specs..."
for adr_id in "${!DEFINED_ADRS[@]}"; do
  adr_path="${DEFINED_ADRS[${adr_id}]}"

  # Extract partition IDs mentioned in 'Affected Partitions' or content (e.g. 05-01, 05-05)
  affected_partitions=$(grep -oE '(partition[s]?[[:space:]]+|[[:space:]]|^)(05-[0-9]{2})' "${adr_path}" | grep -oE '05-[0-9]{2}' | sort -u || true)

  for part_id in ${affected_partitions}; do
    part_req_file="$(find "${REQUIREMENTS_DIR}" -maxdepth 1 -name "${part_id}*.md" | head -n 1 || true)"
    part_dec_file="$(find "${PARTITIONS_DIR}" -maxdepth 1 -name "${part_id}*.md" | head -n 1 || true)"

    has_link=false
    if [ -n "${part_req_file}" ] && [ -f "${part_req_file}" ] && grep -q "${adr_id}" "${part_req_file}"; then
      has_link=true
    fi
    if [ -n "${part_dec_file}" ] && [ -f "${part_dec_file}" ] && grep -q "${adr_id}" "${part_dec_file}"; then
      has_link=true
    fi

    if [ "${has_link}" = false ]; then
      echo "  [UNLINKED PARTITION] ${adr_id} designates partition '${part_id}', but neither ${part_req_file:-${part_id}.md} nor ${part_dec_file:-${part_id}.md} cites ${adr_id}."
      ERRORS=$((ERRORS + 1))
    fi
  done
done

# 6. Final Evaluation
if [ "${ERRORS}" -gt 0 ]; then
  echo "FAILED: Found ${ERRORS} broken or unlinked architectural decision reference(s)."
  exit 1
fi

echo "SUCCESS: All architectural decision references and bidirectional links are valid."
exit 0