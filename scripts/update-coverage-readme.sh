#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
README_FILE="$REPO_ROOT/README.md"
COVERAGE_FILE="$REPO_ROOT/coverage.out"
TEST_OUTPUT_FILE="$REPO_ROOT/.coverage-test-output.tmp"
COVERAGE_FUNC_FILE="$REPO_ROOT/.coverage-func-output.tmp"
COVERAGE_SUMMARY_FILE="$REPO_ROOT/.coverage-summary.tmp"

cd "$REPO_ROOT"

go test ./... -coverprofile="$COVERAGE_FILE" | tee "$TEST_OUTPUT_FILE"
go tool cover -func="$COVERAGE_FILE" > "$COVERAGE_FUNC_FILE"

awk '
  /^ok[[:space:]]+github.com\/Junx27\/junxzy/ && /coverage:/ {
    pkg = $2
    cov = ""
    for (i = 1; i <= NF; i++) {
      if ($i == "coverage:") {
        cov = $(i+1)
        break
      }
    }
    if (cov != "") {
      printf "%s %s\n", pkg, cov
    }
  }
' "$TEST_OUTPUT_FILE" > "$COVERAGE_SUMMARY_FILE"

total_cov="$(awk '/^total:/ {print $3; exit}' "$COVERAGE_FUNC_FILE")"
printf "TOTAL %s\n" "$total_cov" >> "$COVERAGE_SUMMARY_FILE"

TODAY="$(date +%F)"

awk -v coverage_file="$COVERAGE_SUMMARY_FILE" -v today="$TODAY" '
BEGIN {
  while ((getline line < coverage_file) > 0) {
    coverage = coverage line "\n"
  }
}
{
  if ($0 ~ /^Last checked:/) {
    print "Last checked: " today
    next
  }

  if ($0 == "### Coverage summary") {
    print
    in_summary = 1
    next
  }

  if (in_summary && $0 == "```text") {
    print
    printf "%s", coverage
    replacing_block = 1
    in_summary = 0
    next
  }

  if (replacing_block) {
    if ($0 == "```") {
      print
      replacing_block = 0
    }
    next
  }

  print
}
' "$README_FILE" > "$README_FILE.tmp"

mv "$README_FILE.tmp" "$README_FILE"

rm -f "$TEST_OUTPUT_FILE" "$COVERAGE_FUNC_FILE" "$COVERAGE_SUMMARY_FILE"

echo "README coverage summary updated from go test output."
