#!/bin/bash

set -euo pipefail

# Dependency Vulnerability Scanner
# This script scans Go modules, Python packages, and other dependencies for vulnerabilities

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

echo -e "${BLUE}📦 Dependency Vulnerability Scan Starting...${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Output Directory: $OUTPUT_DIR"
echo

cd "$PROJECT_ROOT"

REPORT_FILE="$OUTPUT_DIR/dependency_report_latest"

# Initialize report
cat > "$REPORT_FILE.txt" << EOF
DEPENDENCY VULNERABILITY REPORT
Generated: $(date)
Project: $PROJECT_ROOT

SCAN RESULTS:
=============

EOF

TOTAL_VULNERABILITIES=0
HIGH_SEVERITY_VULNS=0
SCAN_SUCCESS=true

# Scan Go dependencies
echo -e "${BLUE}🐹 Scanning Go dependencies...${NC}"

if [[ -f "go.mod" ]]; then
    echo -e "${GREEN}📦 Found go.mod file${NC}"
    
    # Check if govulncheck is available
    if ! command -v govulncheck &> /dev/null; then
        echo -e "${YELLOW}📦 Installing govulncheck...${NC}"
        go install golang.org/x/vuln/cmd/govulncheck@latest
    fi
    
    # Run govulncheck if available
    if command -v govulncheck &> /dev/null; then
        echo -e "${BLUE}🔍 Running govulncheck...${NC}"
        
        # Run vulnerability check
        if govulncheck -json ./... > "$REPORT_FILE_govuln.json" 2>/dev/null; then
            echo -e "${GREEN}✅ govulncheck completed${NC}"
            
            # Parse results
            GO_VULNS=$(cat "$REPORT_FILE_govuln.json" | python3 -c "
import json, sys
try:
    vulns = []
    for line in sys.stdin:
        try:
            data = json.loads(line.strip())
            if data.get('message', {}).get('vulnerability'):
                vuln = data['message']['vulnerability']
                vulns.append({
                    'id': vuln.get('id', 'Unknown'),
                    'summary': vuln.get('summary', 'No summary'),
                    'severity': vuln.get('severity', 'Unknown'),
                    'module': data.get('message', {}).get('module', {}).get('path', 'Unknown')
                })
        except:
            continue
    
    print(f'Found {len(vulns)} Go vulnerabilities')
    for vuln in vulns[:10]:
        print(f'  - {vuln[\"id\"]}: {vuln[\"summary\"]}')
        print(f'    Module: {vuln[\"module\"]}, Severity: {vuln[\"severity\"]}')
        print()
    
    if len(vulns) > 10:
        print(f'  ... and {len(vulns) - 10} more vulnerabilities')
        
    # Count high severity
    high_sev = len([v for v in vulns if v.get('severity', '').upper() in ['HIGH', 'CRITICAL']])
    print(f'High/Critical severity: {high_sev}')
    
except Exception as e:
    print(f'Error parsing govulncheck results: {e}')
" 2>/dev/null) || GO_VULNS="Error parsing govulncheck results"

            echo "$GO_VULNS"
            
            # Add to main report
            cat >> "$REPORT_FILE.txt" << EOF

Go Dependencies (govulncheck):
$GO_VULNS

EOF
            
            # Count vulnerabilities
            GO_VULN_COUNT=$(echo "$GO_VULNS" | grep -oP 'Found \K\d+' || echo "0")
            GO_HIGH_COUNT=$(echo "$GO_VULNS" | grep -oP 'High/Critical severity: \K\d+' || echo "0")
            TOTAL_VULNERABILITIES=$((TOTAL_VULNERABILITIES + GO_VULN_COUNT))
            HIGH_SEVERITY_VULNS=$((HIGH_SEVERITY_VULNS + GO_HIGH_COUNT))
            
        else
            echo -e "${RED}❌ govulncheck failed${NC}"
            SCAN_SUCCESS=false
        fi
    else
        echo -e "${YELLOW}⚠️  govulncheck not available. Checking with go list...${NC}"
        
        # Fallback: check go modules for known vulnerable versions
        GO_MODULES=$(go list -m all 2>/dev/null | head -20 || echo "No modules found")
        echo "Go modules in use:"
        echo "$GO_MODULES"
        
        cat >> "$REPORT_FILE.txt" << EOF

Go Dependencies (go list):
$GO_MODULES

EOF
    fi
    
    # Check for outdated Go version
    GO_VERSION=$(go version 2>/dev/null || echo "Go not found")
    echo "Go version: $GO_VERSION"
    
    cat >> "$REPORT_FILE.txt" << EOF

Go Version: $GO_VERSION

EOF
    
else
    echo -e "${YELLOW}⚠️  No go.mod file found. Skipping Go dependency scan.${NC}"
fi

# Scan Python dependencies
echo
echo -e "${BLUE}🐍 Scanning Python dependencies...${NC}"

# Find Python requirements files
PYTHON_REQ_FILES=()
for pattern in "requirements*.txt" "pyproject.toml" "setup.py" "Pipfile"; do
    while IFS= read -r -d '' file; do
        PYTHON_REQ_FILES+=("$file")
    done < <(find . -name "$pattern" -not -path "./.venv/*" -not -path "./.git/*" -print0 2>/dev/null)
done

if [[ ${#PYTHON_REQ_FILES[@]} -gt 0 ]]; then
    echo -e "${GREEN}📦 Found Python requirement files: ${PYTHON_REQ_FILES[*]}${NC}"
    
    # Install pip-audit if not available
    if ! command -v pip-audit &> /dev/null; then
        echo -e "${YELLOW}📦 Installing pip-audit...${NC}"
        if command -v pip3 &> /dev/null; then
            pip3 install pip-audit --quiet --user
        elif [[ -f ".venv/bin/pip" ]]; then
            .venv/bin/pip install pip-audit --quiet
        else
            echo -e "${RED}❌ Cannot install pip-audit${NC}"
        fi
    fi
    
    # Run pip-audit if available
    if command -v pip-audit &> /dev/null; then
        for req_file in "${PYTHON_REQ_FILES[@]}"; do
            if [[ "$req_file" == *requirements*.txt ]]; then
                echo -e "${BLUE}🔍 Scanning $req_file with pip-audit...${NC}"
                
                if pip-audit -r "$req_file" --format json --output "$REPORT_FILE_pipaudit_$(basename "$req_file").json" 2>/dev/null; then
                    # Parse results
                    PYTHON_VULNS=$(cat "$REPORT_FILE_pipaudit_$(basename "$req_file").json" | python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    vulnerabilities = data.get('vulnerabilities', [])
    
    print(f'Found {len(vulnerabilities)} Python vulnerabilities in $(basename "$req_file")')
    for vuln in vulnerabilities[:10]:
        package = vuln.get('package', 'Unknown')
        vuln_id = vuln.get('id', 'Unknown')
        summary = vuln.get('summary', 'No summary')
        severity = vuln.get('severity', 'Unknown')
        
        print(f'  - {vuln_id}: {package}')
        print(f'    {summary}')
        print(f'    Severity: {severity}')
        print()
    
    if len(vulnerabilities) > 10:
        print(f'  ... and {len(vulnerabilities) - 10} more vulnerabilities')
        
    # Count high severity
    high_sev = len([v for v in vulnerabilities if v.get('severity', '').upper() in ['HIGH', 'CRITICAL']])
    print(f'High/Critical severity: {high_sev}')
    
except Exception as e:
    print(f'Error parsing pip-audit results: {e}')
" 2>/dev/null) || PYTHON_VULNS="Error parsing pip-audit results"

                    echo "$PYTHON_VULNS"
                    
                    cat >> "$REPORT_FILE.txt" << EOF

Python Dependencies ($req_file):
$PYTHON_VULNS

EOF
                    
                    # Count vulnerabilities
                    PY_VULN_COUNT=$(echo "$PYTHON_VULNS" | grep -oP 'Found \K\d+' || echo "0")
                    PY_HIGH_COUNT=$(echo "$PYTHON_VULNS" | grep -oP 'High/Critical severity: \K\d+' || echo "0")
                    TOTAL_VULNERABILITIES=$((TOTAL_VULNERABILITIES + PY_VULN_COUNT))
                    HIGH_SEVERITY_VULNS=$((HIGH_SEVERITY_VULNS + PY_HIGH_COUNT))
                    
                else
                    echo -e "${RED}❌ pip-audit failed for $req_file${NC}"
                    SCAN_SUCCESS=false
                fi
            fi
        done
    else
        echo -e "${YELLOW}⚠️  pip-audit not available. Listing Python packages...${NC}"
        
        # Fallback: just list requirements
        for req_file in "${PYTHON_REQ_FILES[@]}"; do
            if [[ -f "$req_file" ]]; then
                echo "Contents of $req_file:"
                head -20 "$req_file"
                echo
            fi
        done
    fi
else
    echo -e "${YELLOW}⚠️  No Python requirement files found. Skipping Python dependency scan.${NC}"
fi

# Scan Node.js dependencies (if present)
echo
echo -e "${BLUE}📦 Scanning Node.js dependencies...${NC}"

if [[ -f "package.json" ]]; then
    echo -e "${GREEN}📦 Found package.json file${NC}"
    
    # Check if npm audit is available
    if command -v npm &> /dev/null; then
        echo -e "${BLUE}🔍 Running npm audit...${NC}"
        
        if npm audit --json > "$REPORT_FILE_npmaudit.json" 2>/dev/null; then
            # Parse npm audit results
            NPM_VULNS=$(cat "$REPORT_FILE_npmaudit.json" | python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    vulnerabilities = data.get('vulnerabilities', {})
    
    vuln_list = []
    for package, vuln_data in vulnerabilities.items():
        severity = vuln_data.get('severity', 'Unknown')
        title = vuln_data.get('title', 'No title')
        vuln_list.append({'package': package, 'severity': severity, 'title': title})
    
    print(f'Found {len(vuln_list)} Node.js vulnerabilities')
    for vuln in vuln_list[:10]:
        print(f'  - {vuln[\"package\"]}: {vuln[\"title\"]}')
        print(f'    Severity: {vuln[\"severity\"]}')
        print()
    
    if len(vuln_list) > 10:
        print(f'  ... and {len(vuln_list) - 10} more vulnerabilities')
        
    # Count high severity
    high_sev = len([v for v in vuln_list if v.get('severity', '').upper() in ['HIGH', 'CRITICAL']])
    print(f'High/Critical severity: {high_sev}')
    
except Exception as e:
    print(f'Error parsing npm audit results: {e}')
" 2>/dev/null) || NPM_VULNS="Error parsing npm audit results"

            echo "$NPM_VULNS"
            
            cat >> "$REPORT_FILE.txt" << EOF

Node.js Dependencies (npm audit):
$NPM_VULNS

EOF
            
            # Count vulnerabilities
            NPM_VULN_COUNT=$(echo "$NPM_VULNS" | grep -oP 'Found \K\d+' || echo "0")
            NPM_HIGH_COUNT=$(echo "$NPM_VULNS" | grep -oP 'High/Critical severity: \K\d+' || echo "0")
            TOTAL_VULNERABILITIES=$((TOTAL_VULNERABILITIES + NPM_VULN_COUNT))
            HIGH_SEVERITY_VULNS=$((HIGH_SEVERITY_VULNS + NPM_HIGH_COUNT))
            
        else
            echo -e "${RED}❌ npm audit failed${NC}"
            SCAN_SUCCESS=false
        fi
    else
        echo -e "${YELLOW}⚠️  npm not available. Skipping Node.js dependency scan.${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  No package.json file found. Skipping Node.js dependency scan.${NC}"
fi

# Summary
cat >> "$REPORT_FILE.txt" << EOF

SUMMARY:
========
Total vulnerabilities found: $TOTAL_VULNERABILITIES
High/Critical severity vulnerabilities: $HIGH_SEVERITY_VULNS
Scan completed successfully: $SCAN_SUCCESS

RECOMMENDATIONS:
================
1. Review and update all vulnerable dependencies to secure versions
2. Consider using dependency pinning for critical applications
3. Set up automated dependency vulnerability monitoring
4. Implement dependency update policies and testing procedures
5. Use tools like Dependabot or Renovate for automated updates

NEXT STEPS:
===========
1. Prioritize fixing HIGH and CRITICAL severity vulnerabilities
2. Test applications after updating dependencies
3. Monitor for new vulnerabilities in your dependency stack
4. Consider using alternative packages if vulnerabilities persist

EOF

echo
echo -e "${BLUE}📊 Dependency Vulnerability Scan Results:${NC}"
echo "Total vulnerabilities found: $TOTAL_VULNERABILITIES"
echo "High/Critical severity vulnerabilities: $HIGH_SEVERITY_VULNS"
echo "Scan completed successfully: $SCAN_SUCCESS"

echo
echo -e "${BLUE}📁 Reports generated:${NC}"
echo "  - Main report: $REPORT_FILE.txt"
if [[ -f "$REPORT_FILE_govuln.json" ]]; then
    echo "  - Go vulnerabilities (JSON): $REPORT_FILE_govuln.json"
fi
for json_file in "$REPORT_FILE"_pipaudit_*.json; do
    if [[ -f "$json_file" ]]; then
        echo "  - Python vulnerabilities (JSON): $json_file"
    fi
done
if [[ -f "$REPORT_FILE_npmaudit.json" ]]; then
    echo "  - Node.js vulnerabilities (JSON): $REPORT_FILE_npmaudit.json"
fi

# Exit with appropriate code
if [[ $HIGH_SEVERITY_VULNS -gt 0 ]]; then
    echo -e "${RED}❌ CRITICAL: Found $HIGH_SEVERITY_VULNS high/critical severity vulnerabilities!${NC}"
    exit 2
elif [[ $TOTAL_VULNERABILITIES -gt 0 ]]; then
    echo -e "${YELLOW}⚠️  Found $TOTAL_VULNERABILITIES vulnerabilities. Review and update recommended.${NC}"
    exit 1
elif [[ "$SCAN_SUCCESS" == "false" ]]; then
    echo -e "${YELLOW}⚠️  Some scans failed. Check the reports for details.${NC}"
    exit 1
else
    echo -e "${GREEN}✅ No vulnerabilities found in dependencies!${NC}"
    exit 0
fi 