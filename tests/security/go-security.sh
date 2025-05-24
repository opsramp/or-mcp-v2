#!/bin/bash

set -euo pipefail

# Go Security Scanning with gosec
# This script scans Go code for security vulnerabilities

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/tests/security/reports"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo -e "${BLUE}🔍 Go Security Scan Starting...${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Output Directory: $OUTPUT_DIR"
echo

# Check if gosec is installed
if ! command -v gosec &> /dev/null; then
    echo -e "${RED}❌ gosec not found. Installing...${NC}"
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    if ! command -v gosec &> /dev/null; then
        echo -e "${RED}❌ Failed to install gosec${NC}"
        exit 1
    fi
fi

echo -e "${GREEN}✅ gosec version: $(gosec --version)${NC}"

cd "$PROJECT_ROOT"

# Configuration for gosec
# Exclude G115 (integer overflow) as it has many false positives
# Include critical rules only for high confidence
GOSEC_EXCLUDE="G115,G101" # G101 for hardcoded credentials (we handle separately)
GOSEC_INCLUDE="G102,G103,G104,G106,G107,G108,G109,G110,G111,G112,G114,G201,G202,G203,G204,G301,G302,G303,G304,G305,G306,G401,G402,G403,G404,G405,G501,G502,G503,G504,G505"

echo -e "${YELLOW}📋 Running gosec with configuration:${NC}"
echo "  Excluded rules: $GOSEC_EXCLUDE"
echo "  Included rules: $GOSEC_INCLUDE"
echo

# Run gosec scan
REPORT_FILE="$OUTPUT_DIR/gosec_report_latest"

echo -e "${BLUE}🔍 Scanning Go code for security issues...${NC}"

# Text report for console
if gosec -exclude="$GOSEC_EXCLUDE" -fmt=text -out="$REPORT_FILE.txt" ./...; then
    SCAN_RESULT="SUCCESS"
else
    SCAN_RESULT="ISSUES_FOUND"
fi

# JSON report for processing
gosec -exclude="$GOSEC_EXCLUDE" -fmt=json -out="$REPORT_FILE.json" ./... || true

# SARIF report for CI integration
gosec -exclude="$GOSEC_EXCLUDE" -fmt=sarif -out="$REPORT_FILE.sarif" ./... || true

echo
echo -e "${BLUE}📊 Security Scan Results:${NC}"

# Parse JSON results
if [[ -f "$REPORT_FILE.json" ]]; then
    ISSUES=$(cat "$REPORT_FILE.json" | python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    issues = data.get('Issues', [])
    if issues:
        print(f'Found {len(issues)} security issues:')
        for issue in issues[:10]:  # Show first 10 issues
            print(f'  - {issue.get(\"rule_id\", \"Unknown\")}: {issue.get(\"details\", \"No details\")}')
            print(f'    File: {issue.get(\"file\", \"Unknown\")}:{issue.get(\"line\", \"?\")}')
            print(f'    Severity: {issue.get(\"severity\", \"Unknown\")}, Confidence: {issue.get(\"confidence\", \"Unknown\")}')
            print()
        if len(issues) > 10:
            print(f'  ... and {len(issues) - 10} more issues')
    else:
        print('No security issues found! 🎉')
except Exception as e:
    print(f'Error parsing results: {e}')
" 2>/dev/null) || ISSUES="Error parsing JSON results"

    echo "$ISSUES"
else
    echo -e "${RED}❌ No JSON report generated${NC}"
fi

echo
echo -e "${BLUE}📁 Reports generated:${NC}"
echo "  - Text report: $REPORT_FILE.txt"
echo "  - JSON report: $REPORT_FILE.json"
echo "  - SARIF report: $REPORT_FILE.sarif"

# Display summary
if [[ "$SCAN_RESULT" == "SUCCESS" ]]; then
    echo -e "${GREEN}✅ Go security scan completed successfully!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Go security scan found potential issues. Review reports above.${NC}"
    exit 1
fi 