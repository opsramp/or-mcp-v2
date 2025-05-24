#!/usr/bin/env bash

set -euo pipefail

# Secret Detection Scanner
# This script scans for hardcoded credentials, API keys, and sensitive information

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

echo -e "${BLUE}🔐 Secret Detection Scan Starting...${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Output Directory: $OUTPUT_DIR"
echo

cd "$PROJECT_ROOT"

REPORT_FILE="$OUTPUT_DIR/secrets_report_latest"

# Common secret patterns (regex) - using arrays for compatibility
SECRET_TYPES=(
    "API_Keys"
    "AWS_Access_Key"
    "AWS_Secret_Key"
    "JWT_Token"
    "Generic_Secret"
    "Database_URL"
    "Private_Key"
    "SSH_Key"
    "GitHub_Token"
    "Slack_Token"
    "Google_API_Key"
    "OAuth_Token"
)

SECRET_PATTERNS=(
    "(?i)(api[_-]?key|apikey)[[:space:]]*[=:][[:space:]]*['\"]?[a-zA-Z0-9]{16,}['\"]?"
    "AKIA[0-9A-Z]{16}"
    "(?i)(aws[_-]?secret|secret[_-]?key)[[:space:]]*[=:][[:space:]]*['\"]?[a-zA-Z0-9/+=]{40}['\"]?"
    "eyJ[a-zA-Z0-9]{15,}\\.[a-zA-Z0-9]{15,}\\.[a-zA-Z0-9_-]{15,}"
    "(?i)(secret|password|pwd|pass)[[:space:]]*[=:][[:space:]]*['\"]?[a-zA-Z0-9@#$%^&*()_+-=]{8,}['\"]?"
    "(?i)(database[_-]?url|db[_-]?url)[[:space:]]*[=:][[:space:]]*['\"]?[a-zA-Z0-9+://@.-]{20,}['\"]?"
    "-----BEGIN [A-Z ]*PRIVATE KEY-----"
    "ssh-[a-z0-9]+ [A-Za-z0-9+/=]+"
    "ghp_[a-zA-Z0-9]{36}"
    "xox[baprs]-[a-zA-Z0-9-]+"
    "AIza[0-9A-Za-z\\-_]{35}"
    "(?i)(access[_-]?token|oauth[_-]?token)[[:space:]]*[=:][[:space:]]*['\"]?[a-zA-Z0-9_-]{20,}['\"]?"
)

# Files to exclude from scanning
EXCLUDE_PATTERNS=(
    "*.git/*"
    "*.venv/*"
    "*/__pycache__/*"
    "*/node_modules/*"
    "*/vendor/*"
    "*.log"
    "*.pyc"
    "*.so"
    "*.dylib"
    "*.png"
    "*.jpg"
    "*.jpeg"
    "*.gif"
    "*.pdf"
    "*.zip"
    "*.tar.gz"
    "*/tests/security/reports/*"
)

# Build exclude arguments for grep
EXCLUDE_ARGS=""
for pattern in "${EXCLUDE_PATTERNS[@]}"; do
    EXCLUDE_ARGS="$EXCLUDE_ARGS --exclude-dir=${pattern%/*} --exclude=${pattern##*/}"
done

echo -e "${BLUE}🔍 Scanning for secrets and sensitive information...${NC}"

# Initialize report
cat > "$REPORT_FILE.txt" << EOF
SECRET DETECTION REPORT
Generated: $(date)
Project: $PROJECT_ROOT

SCAN RESULTS:
=============

EOF

TOTAL_SECRETS=0
CRITICAL_SECRETS=0

# Scan for each secret pattern
for i in "${!SECRET_TYPES[@]}"; do
    secret_type="${SECRET_TYPES[$i]}"
    pattern="${SECRET_PATTERNS[$i]}"
    
    echo -e "${YELLOW}🔍 Scanning for: $secret_type${NC}"
    
    # Use grep with Perl regex
    if command -v grep &> /dev/null; then
        MATCHES=$(grep -rPn "$pattern" . $EXCLUDE_ARGS 2>/dev/null || true)
    else
        echo -e "${RED}❌ grep not available${NC}"
        continue
    fi
    
    if [[ -n "$MATCHES" ]]; then
        COUNT=$(echo "$MATCHES" | wc -l)
        TOTAL_SECRETS=$((TOTAL_SECRETS + COUNT))
        
        # Mark certain types as critical
        if [[ "$secret_type" =~ (API_Keys|AWS|Private_Key|SSH_Key|GitHub_Token|OAuth_Token) ]]; then
            CRITICAL_SECRETS=$((CRITICAL_SECRETS + COUNT))
            echo -e "${RED}🚨 CRITICAL: Found $COUNT potential $secret_type${NC}"
        else
            echo -e "${YELLOW}⚠️  Found $COUNT potential $secret_type${NC}"
        fi
        
        # Add to report
        cat >> "$REPORT_FILE.txt" << EOF

$secret_type ($COUNT matches):
$(echo "$MATCHES" | head -10)
$(if [[ $(echo "$MATCHES" | wc -l) -gt 10 ]]; then echo "... and $(($(echo "$MATCHES" | wc -l) - 10)) more matches"; fi)

EOF
        
        # Add to JSON report for processing
        if [[ ! -f "$REPORT_FILE.json" ]]; then
            echo '{"secrets": []}' > "$REPORT_FILE.json"
        fi
        
        # Add matches to JSON (simplified)
        python3 << EOF >> /dev/null 2>&1 || true
import json
import re

try:
    with open("$REPORT_FILE.json", "r") as f:
        data = json.load(f)
    
    matches = """$MATCHES""".strip().split('\n') if """$MATCHES""".strip() else []
    
    for match in matches[:10]:  # Limit to first 10 matches
        if ':' in match:
            parts = match.split(':', 2)
            if len(parts) >= 3:
                data["secrets"].append({
                    "type": "$secret_type",
                    "file": parts[0],
                    "line": parts[1],
                    "content": parts[2][:100] + "..." if len(parts[2]) > 100 else parts[2],
                    "severity": "CRITICAL" if "$secret_type" in ["API_Keys", "AWS_Access_Key", "AWS_Secret_Key", "Private_Key", "SSH_Key", "GitHub_Token", "OAuth_Token"] else "MEDIUM"
                })
    
    with open("$REPORT_FILE.json", "w") as f:
        json.dump(data, f, indent=2)
        
except Exception as e:
    pass
EOF
        
    else
        echo -e "${GREEN}✅ No $secret_type found${NC}"
    fi
done

# Additional entropy-based detection for unknown secret patterns
echo -e "${YELLOW}🔍 Running entropy-based detection...${NC}"

# Look for high-entropy strings that might be secrets
ENTROPY_MATCHES=$(find . -type f \( -name "*.py" -o -name "*.go" -o -name "*.js" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.env" -o -name "*.conf" -o -name "*.config" \) \
    -not -path "./.venv/*" -not -path "./.git/*" -not -path "./tests/security/reports/*" \
    -exec grep -l '[a-zA-Z0-9+/=]{32,}' {} \; 2>/dev/null | head -20 || true)

if [[ -n "$ENTROPY_MATCHES" ]]; then
    echo -e "${YELLOW}⚠️  Found files with potential high-entropy strings:${NC}"
    echo "$ENTROPY_MATCHES" | while read -r file; do
        echo "  - $file"
    done
    
    cat >> "$REPORT_FILE.txt" << EOF

High-Entropy Strings:
Files with potential secrets based on entropy analysis:
$ENTROPY_MATCHES

EOF
fi

# Check for common configuration files with secrets
echo -e "${YELLOW}🔍 Checking configuration files...${NC}"

CONFIG_FILES=(
    ".env"
    ".env.local"
    ".env.production"
    "config.yaml"
    "config.yml"
    "secrets.yaml"
    "credentials.json"
    ".aws/credentials"
    ".ssh/id_rsa"
    ".ssh/id_ecdsa"
    ".ssh/id_ed25519"
)

FOUND_CONFIG_FILES=()
for config_file in "${CONFIG_FILES[@]}"; do
    if find . -name "$(basename "$config_file")" -not -path "./.git/*" | grep -q .; then
        FOUND_CONFIG_FILES+=("$config_file")
    fi
done

if [[ ${#FOUND_CONFIG_FILES[@]} -gt 0 ]]; then
    echo -e "${YELLOW}⚠️  Found sensitive configuration files:${NC}"
    printf '%s\n' "${FOUND_CONFIG_FILES[@]}" | while read -r file; do
        echo "  - $file"
    done
    
    cat >> "$REPORT_FILE.txt" << EOF

Sensitive Configuration Files:
${FOUND_CONFIG_FILES[*]}

EOF
fi

# Summary
cat >> "$REPORT_FILE.txt" << EOF

SUMMARY:
========
Total potential secrets found: $TOTAL_SECRETS
Critical secrets found: $CRITICAL_SECRETS
Sensitive config files found: ${#FOUND_CONFIG_FILES[@]}

RECOMMENDATIONS:
================
1. Review all identified potential secrets manually
2. Use environment variables for secrets instead of hardcoding
3. Add sensitive files to .gitignore
4. Consider using secret management tools (HashiCorp Vault, AWS Secrets Manager, etc.)
5. Implement pre-commit hooks to prevent secrets from being committed

EOF

echo
echo -e "${BLUE}📊 Secret Detection Results:${NC}"
echo "Total potential secrets found: $TOTAL_SECRETS"
echo "Critical secrets found: $CRITICAL_SECRETS"
echo "Sensitive config files found: ${#FOUND_CONFIG_FILES[@]}"

echo
echo -e "${BLUE}📁 Reports generated:${NC}"
echo "  - Text report: $REPORT_FILE.txt"
echo "  - JSON report: $REPORT_FILE.json"

# Exit with appropriate code
if [[ $CRITICAL_SECRETS -gt 0 ]]; then
    echo -e "${RED}❌ CRITICAL: Found $CRITICAL_SECRETS critical secrets that need immediate attention!${NC}"
    exit 2
elif [[ $TOTAL_SECRETS -gt 0 ]]; then
    echo -e "${YELLOW}⚠️  Found $TOTAL_SECRETS potential secrets. Review recommended.${NC}"
    exit 1
else
    echo -e "${GREEN}✅ No secrets detected in codebase!${NC}"
    exit 0
fi 