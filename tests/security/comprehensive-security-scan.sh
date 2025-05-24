#!/bin/bash

set -euo pipefail

# Comprehensive Security Scanner
# This script orchestrates all security scans and produces a consolidated report

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/tests/security/reports"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo -e "${PURPLE}🛡️  COMPREHENSIVE SECURITY SCAN${NC}"
echo -e "${PURPLE}=================================${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Output Directory: $OUTPUT_DIR"
echo "Scan Started: $(date)"
echo

MASTER_REPORT="$OUTPUT_DIR/security_master_report_latest"

# Initialize master report
cat > "$MASTER_REPORT.txt" << EOF
COMPREHENSIVE SECURITY SCAN REPORT
===================================
Generated: $(date)
Project: $PROJECT_ROOT
Scan ID: latest

EXECUTIVE SUMMARY:
==================

EOF

# Scan results tracking
SCAN_RESULTS=()
TOTAL_CRITICAL_ISSUES=0
TOTAL_HIGH_ISSUES=0
TOTAL_MEDIUM_ISSUES=0
TOTAL_LOW_ISSUES=0
FAILED_SCANS=()

# Function to run individual scan and track results
run_security_scan() {
    local scan_name="$1"
    local scan_script="$2"
    local scan_description="$3"
    
    echo -e "${CYAN}🔍 Running $scan_name...${NC}"
    echo -e "${BLUE}   $scan_description${NC}"
    
    local start_time=$(date +%s)
    local result_code=0
    local scan_output=""
    
    # Make script executable
    chmod +x "$SCRIPT_DIR/$scan_script"
    
    # Run the scan and capture output
    if scan_output=$("$SCRIPT_DIR/$scan_script" 2>&1); then
        result_code=$?
    else
        result_code=$?
    fi
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    # Determine result status
    local status
    case $result_code in
        0) status="✅ PASS" ;;
        1) status="⚠️  ISSUES FOUND" ;;
        2) status="🚨 CRITICAL ISSUES" ;;
        *) status="❌ FAILED" ;;
    esac
    
    echo -e "${BLUE}   Duration: ${duration}s | Status: $status${NC}"
    
    # Track results
    SCAN_RESULTS+=("$scan_name|$status|$duration|$result_code")
    
    if [[ $result_code -gt 2 ]]; then
        FAILED_SCANS+=("$scan_name")
    fi
    
    # Add to master report
    cat >> "$MASTER_REPORT.txt" << EOF

$scan_name SCAN:
$(printf '=%.0s' {1..50})
Status: $status
Duration: ${duration} seconds
Description: $scan_description

Output:
$scan_output

EOF
    
    echo
}

# Execute all security scans
echo -e "${PURPLE}🚀 Starting Security Scan Suite...${NC}"
echo

# 1. Go Security Scan
run_security_scan \
    "Go Code Security" \
    "go-security.sh" \
    "Static analysis of Go code for security vulnerabilities using gosec"

# 2. Python Security Scan  
run_security_scan \
    "Python Code Security" \
    "python-security.sh" \
    "Static analysis of Python code for security vulnerabilities using bandit"

# 3. Secret Detection Scan
run_security_scan \
    "Secret Detection" \
    "secret-scan.sh" \
    "Detection of hardcoded credentials, API keys, and sensitive information"

# 4. Dependency Vulnerability Scan
run_security_scan \
    "Dependency Vulnerabilities" \
    "dependency-scan.sh" \
    "Vulnerability scanning of Go modules, Python packages, and Node.js dependencies"

# Generate comprehensive analysis
echo -e "${PURPLE}📊 Generating Comprehensive Analysis...${NC}"

# Parse individual scan reports for metrics
TOTAL_SCANS=${#SCAN_RESULTS[@]}
PASSED_SCANS=0
FAILED_SCANS_COUNT=0
CRITICAL_SCANS=0
WARNING_SCANS=0

for result in "${SCAN_RESULTS[@]}"; do
    IFS='|' read -ra PARTS <<< "$result"
    status="${PARTS[1]}"
    
    case "$status" in
        "✅ PASS") ((PASSED_SCANS++)) ;;
        "⚠️  ISSUES FOUND") ((WARNING_SCANS++)) ;;
        "🚨 CRITICAL ISSUES") ((CRITICAL_SCANS++)) ;;
        "❌ FAILED") ((FAILED_SCANS_COUNT++)) ;;
    esac
done

# Calculate security score
SECURITY_SCORE=100
if [[ $CRITICAL_SCANS -gt 0 ]]; then
    SECURITY_SCORE=$((SECURITY_SCORE - (CRITICAL_SCANS * 40)))
fi
if [[ $WARNING_SCANS -gt 0 ]]; then
    SECURITY_SCORE=$((SECURITY_SCORE - (WARNING_SCANS * 20)))
fi
if [[ $FAILED_SCANS_COUNT -gt 0 ]]; then
    SECURITY_SCORE=$((SECURITY_SCORE - (FAILED_SCANS_COUNT * 10)))
fi
SECURITY_SCORE=$((SECURITY_SCORE > 0 ? SECURITY_SCORE : 0))

# Determine overall security level
if [[ $SECURITY_SCORE -ge 90 ]]; then
    SECURITY_LEVEL="🟢 EXCELLENT"
elif [[ $SECURITY_SCORE -ge 70 ]]; then
    SECURITY_LEVEL="🟡 GOOD"
elif [[ $SECURITY_SCORE -ge 50 ]]; then
    SECURITY_LEVEL="🟠 MODERATE"
else
    SECURITY_LEVEL="🔴 POOR"
fi

# Update executive summary
sed -i.bak '/EXECUTIVE SUMMARY:/,/^$/c\
EXECUTIVE SUMMARY:\
==================\
\
Security Score: '"$SECURITY_SCORE"'/100 ('"$SECURITY_LEVEL"')\
Total Scans Performed: '"$TOTAL_SCANS"'\
Scans Passed: '"$PASSED_SCANS"'\
Scans with Warnings: '"$WARNING_SCANS"'\
Scans with Critical Issues: '"$CRITICAL_SCANS"'\
Failed Scans: '"$FAILED_SCANS_COUNT"'\
\
SCAN RESULTS SUMMARY:\
=====================' "$MASTER_REPORT.txt"

# Add scan results table
for result in "${SCAN_RESULTS[@]}"; do
    IFS='|' read -ra PARTS <<< "$result"
    name="${PARTS[0]}"
    status="${PARTS[1]}"
    duration="${PARTS[2]}"
    
    echo "$name: $status (${duration}s)" >> "$MASTER_REPORT.txt"
done

# Add recommendations
cat >> "$MASTER_REPORT.txt" << EOF

SECURITY RECOMMENDATIONS:
=========================

IMMEDIATE ACTION REQUIRED:
EOF

if [[ $CRITICAL_SCANS -gt 0 ]]; then
    cat >> "$MASTER_REPORT.txt" << EOF
🚨 CRITICAL: Address all critical security issues immediately
   - Review all scans marked as "CRITICAL ISSUES"
   - Implement fixes before deployment
   - Consider security incident response procedures
EOF
fi

if [[ $WARNING_SCANS -gt 0 ]]; then
    cat >> "$MASTER_REPORT.txt" << EOF
⚠️  WARNINGS: Review and address warning-level issues
   - Plan remediation for identified vulnerabilities
   - Update dependencies with known vulnerabilities
   - Remove or secure hardcoded secrets
EOF
fi

if [[ $FAILED_SCANS_COUNT -gt 0 ]]; then
    cat >> "$MASTER_REPORT.txt" << EOF
❌ FAILURES: Investigate and resolve scan failures
   - Check tool installations and configurations
   - Review scan logs for specific error messages
   - Ensure all required dependencies are available
EOF
fi

cat >> "$MASTER_REPORT.txt" << EOF

ONGOING SECURITY MEASURES:
1. Implement automated security scanning in CI/CD pipeline
2. Regular dependency updates and vulnerability monitoring
3. Code review processes with security focus
4. Security training for development team
5. Incident response plan and procedures

NEXT STEPS:
1. Address all CRITICAL issues immediately
2. Create tickets for all WARNING-level issues
3. Implement automated scanning in development workflow
4. Schedule regular security reviews and updates
5. Consider penetration testing for production systems

COMPLIANCE CONSIDERATIONS:
- Document all security findings and remediation efforts
- Maintain security scan history for audit purposes
- Ensure adherence to relevant security standards (OWASP, NIST, etc.)
- Regular security assessments and reviews

EOF

# Generate JSON summary for automation
cat > "$MASTER_REPORT.json" << EOF
{
  "timestamp": "$(date -Iseconds)",
  "project_root": "$PROJECT_ROOT",
  "security_score": $SECURITY_SCORE,
  "security_level": "$SECURITY_LEVEL",
  "total_scans": $TOTAL_SCANS,
  "passed_scans": $PASSED_SCANS,
  "warning_scans": $WARNING_SCANS,
  "critical_scans": $CRITICAL_SCANS,
  "failed_scans": $FAILED_SCANS_COUNT,
  "scan_results": [
EOF

first=true
for result in "${SCAN_RESULTS[@]}"; do
    IFS='|' read -ra PARTS <<< "$result"
    name="${PARTS[0]}"
    status="${PARTS[1]}"
    duration="${PARTS[2]}"
    code="${PARTS[3]}"
    
    if [[ "$first" == "true" ]]; then
        first=false
    else
        echo "," >> "$MASTER_REPORT.json"
    fi
    
    cat >> "$MASTER_REPORT.json" << EOF
    {
      "name": "$name",
      "status": "$status",
      "duration": $duration,
      "exit_code": $code
    }
EOF
done

cat >> "$MASTER_REPORT.json" << EOF

  ],
  "recommendations": {
    "immediate_action_required": $(if [[ $CRITICAL_SCANS -gt 0 ]]; then echo "true"; else echo "false"; fi),
    "warnings_to_address": $(if [[ $WARNING_SCANS -gt 0 ]]; then echo "true"; else echo "false"; fi),
    "failed_scans_to_investigate": $(if [[ $FAILED_SCANS_COUNT -gt 0 ]]; then echo "true"; else echo "false"; fi)
  }
}
EOF

# Final summary
echo -e "${PURPLE}📋 COMPREHENSIVE SECURITY SCAN COMPLETE${NC}"
echo -e "${PURPLE}=========================================${NC}"
echo
echo -e "${BLUE}🎯 Security Score: $SECURITY_SCORE/100 ($SECURITY_LEVEL)${NC}"
echo -e "${BLUE}📊 Scan Summary:${NC}"
echo "   Total Scans: $TOTAL_SCANS"
echo "   ✅ Passed: $PASSED_SCANS"
echo "   ⚠️  Warnings: $WARNING_SCANS"
echo "   🚨 Critical: $CRITICAL_SCANS"
echo "   ❌ Failed: $FAILED_SCANS_COUNT"
echo
echo -e "${BLUE}📁 Reports Generated:${NC}"
echo "   - Master Report: $MASTER_REPORT.txt"
echo "   - JSON Summary: $MASTER_REPORT.json"
echo "   - Individual Reports: $OUTPUT_DIR/"
echo
echo -e "${BLUE}🕒 Scan Completed: $(date)${NC}"

# Determine exit code based on results
if [[ $CRITICAL_SCANS -gt 0 ]]; then
    echo -e "${RED}🚨 CRITICAL SECURITY ISSUES FOUND - IMMEDIATE ACTION REQUIRED!${NC}"
    exit 2
elif [[ $WARNING_SCANS -gt 0 ]]; then
    echo -e "${YELLOW}⚠️  Security issues found. Review and address recommended.${NC}"
    exit 1
elif [[ $FAILED_SCANS_COUNT -gt 0 ]]; then
    echo -e "${YELLOW}⚠️  Some scans failed. Investigation required.${NC}"
    exit 1
else
    echo -e "${GREEN}✅ All security scans passed successfully!${NC}"
    exit 0
fi 
