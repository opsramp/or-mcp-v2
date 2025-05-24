#!/bin/bash

set -euo pipefail

# Python Security Scanning with bandit and safety
# This script scans Python code for security vulnerabilities

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

echo -e "${BLUE}🐍 Python Security Scan Starting...${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Output Directory: $OUTPUT_DIR"
echo

cd "$PROJECT_ROOT"

# Find Python directories
PYTHON_DIRS=$(find . -name "*.py" -not -path "./.venv/*" -not -path "./.git/*" -not -path "./build/*" -not -path "./__pycache__/*" | head -1 | xargs dirname 2>/dev/null || echo "")

if [[ -z "$PYTHON_DIRS" ]]; then
    echo -e "${YELLOW}⚠️  No Python files found. Skipping Python security scan.${NC}"
    exit 0
fi

echo -e "${GREEN}📁 Found Python code in project${NC}"

# Install bandit if not available
if ! command -v bandit &> /dev/null; then
    echo -e "${YELLOW}📦 Installing bandit...${NC}"
    if command -v pip3 &> /dev/null; then
        pip3 install bandit[toml] --quiet
    elif [[ -f "$PROJECT_ROOT/.venv/bin/pip" ]]; then
        "$PROJECT_ROOT/.venv/bin/pip" install bandit[toml] --quiet
    else
        echo -e "${RED}❌ Cannot install bandit. Please install manually.${NC}"
        exit 1
    fi
fi

# Install safety if not available
if ! command -v safety &> /dev/null; then
    echo -e "${YELLOW}📦 Installing safety...${NC}"
    if command -v pip3 &> /dev/null; then
        pip3 install safety --quiet
    elif [[ -f "$PROJECT_ROOT/.venv/bin/pip" ]]; then
        "$PROJECT_ROOT/.venv/bin/pip" install safety --quiet
    else
        echo -e "${YELLOW}⚠️  Cannot install safety. Skipping dependency vulnerability scan.${NC}"
    fi
fi

# Run bandit security scan
echo -e "${BLUE}🔍 Running bandit security scan...${NC}"

BANDIT_REPORT="$OUTPUT_DIR/bandit_report_latest"
BANDIT_CONFIG="$PROJECT_ROOT/.bandit"

# Create bandit configuration to reduce false positives
cat > "$BANDIT_CONFIG" << EOF
[bandit]
exclude_dirs = [".venv", ".git", "build", "__pycache__", "tests"]
skips = ["B101", "B601"]  # Skip assert_used_check and shell_check for testing

# Confidence levels: HIGH, MEDIUM, LOW
# Severity levels: HIGH, MEDIUM, LOW
confidence_level = MEDIUM
severity_level = MEDIUM
EOF

# Run bandit with different output formats
if bandit -r . -f txt -o "$BANDIT_REPORT.txt" --config "$BANDIT_CONFIG" 2>/dev/null; then
    BANDIT_RESULT="SUCCESS"
else
    BANDIT_RESULT="ISSUES_FOUND"
fi

# JSON format for processing
bandit -r . -f json -o "$BANDIT_REPORT.json" --config "$BANDIT_CONFIG" 2>/dev/null || true

# CSV format for spreadsheet analysis
bandit -r . -f csv -o "$BANDIT_REPORT.csv" --config "$BANDIT_CONFIG" 2>/dev/null || true

echo -e "${BLUE}📊 Bandit Results:${NC}"

# Parse bandit JSON results
if [[ -f "$BANDIT_REPORT.json" ]]; then
    BANDIT_SUMMARY=$(cat "$BANDIT_REPORT.json" | python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    results = data.get('results', [])
    metrics = data.get('metrics', {})
    
    print(f'Files scanned: {metrics.get(\"_totals\", {}).get(\"loc\", \"Unknown\")} lines of code')
    print(f'Security issues found: {len(results)}')
    
    if results:
        high_severity = [r for r in results if r.get('issue_severity') == 'HIGH']
        medium_severity = [r for r in results if r.get('issue_severity') == 'MEDIUM'] 
        low_severity = [r for r in results if r.get('issue_severity') == 'LOW']
        
        print(f'  - HIGH severity: {len(high_severity)}')
        print(f'  - MEDIUM severity: {len(medium_severity)}')
        print(f'  - LOW severity: {len(low_severity)}')
        
        print('')
        print('Top 5 issues:')
        for issue in results[:5]:
            print(f'  - {issue.get(\"test_id\", \"Unknown\")}: {issue.get(\"issue_text\", \"No description\")}')
            print(f'    File: {issue.get(\"filename\", \"Unknown\")}:{issue.get(\"line_number\", \"?\")}')
            print(f'    Severity: {issue.get(\"issue_severity\", \"Unknown\")}, Confidence: {issue.get(\"issue_confidence\", \"Unknown\")}')
            print()
    else:
        print('No security issues found! 🎉')
        
except Exception as e:
    print(f'Error parsing bandit results: {e}')
" 2>/dev/null) || BANDIT_SUMMARY="Error parsing bandit results"

    echo "$BANDIT_SUMMARY"
else
    echo -e "${RED}❌ No bandit JSON report generated${NC}"
fi

# Run safety dependency vulnerability scan
echo
echo -e "${BLUE}🔍 Running safety dependency vulnerability scan...${NC}"

SAFETY_RESULT="SUCCESS"
SAFETY_REPORT="$OUTPUT_DIR/safety_report_latest"

# Check for Python requirements files
REQ_FILES=()
for file in requirements.txt requirements-dev.txt pyproject.toml setup.py; do
    if [[ -f "$PROJECT_ROOT/$file" ]]; then
        REQ_FILES+=("$file")
    fi
done

# Also check in subdirectories
for file in $(find . -name "requirements*.txt" -not -path "./.venv/*" -not -path "./.git/*" 2>/dev/null); do
    REQ_FILES+=("$file")
done

if [[ ${#REQ_FILES[@]} -eq 0 ]]; then
    echo -e "${YELLOW}⚠️  No Python requirements files found. Skipping dependency scan.${NC}"
else
    echo -e "${GREEN}📦 Found requirements files: ${REQ_FILES[*]}${NC}"
    
    if command -v safety &> /dev/null; then
        for req_file in "${REQ_FILES[@]}"; do
            echo "Scanning $req_file..."
            if ! safety check -r "$req_file" --json --output "$SAFETY_REPORT_${req_file//\//_}.json" 2>/dev/null; then
                SAFETY_RESULT="ISSUES_FOUND"
            fi
            safety check -r "$req_file" --output "$SAFETY_REPORT_${req_file//\//_}.txt" 2>/dev/null || true
        done
    else
        echo -e "${YELLOW}⚠️  Safety not available. Skipping dependency vulnerability scan.${NC}"
    fi
fi

# Clean up bandit config
rm -f "$BANDIT_CONFIG"

echo
echo -e "${BLUE}📁 Reports generated:${NC}"
echo "  - Bandit text report: $BANDIT_REPORT.txt"
echo "  - Bandit JSON report: $BANDIT_REPORT.json"
echo "  - Bandit CSV report: $BANDIT_REPORT.csv"
if [[ ${#REQ_FILES[@]} -gt 0 ]] && command -v safety &> /dev/null; then
    echo "  - Safety reports: $SAFETY_REPORT*.txt and *.json"
fi

# Overall result
if [[ "$BANDIT_RESULT" == "SUCCESS" && "$SAFETY_RESULT" == "SUCCESS" ]]; then
    echo -e "${GREEN}✅ Python security scan completed successfully!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Python security scan found potential issues. Review reports above.${NC}"
    exit 1
fi 