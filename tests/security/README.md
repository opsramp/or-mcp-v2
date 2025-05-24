# 🛡️ Security Testing Framework

Comprehensive security vulnerability scanning for the HPE OpsRamp MCP integration platform.

## 📋 Overview

This security testing framework provides **zero-tolerance security scanning** with professional-grade tools and comprehensive reporting. All scans are orchestrated via Makefile targets for consistency and automation.

## 🚀 Quick Start

```bash
# Run comprehensive security scan (recommended)
make security-full

# Run quick security scan (Go code + secrets)
make security-scan

# Show detailed security help
make security-help

# Clean previous reports (optional - reports auto-overwrite)
make security-clean
```

## 🔧 Security Scan Types

### 1. **Go Code Security** (`security-go`)
- **Tool**: `gosec` - Industry-standard Go security analyzer
- **Scans**: 60+ security rules (G101-G602)
- **Detects**: SQL injection, command injection, hardcoded credentials, weak cryptography, file permissions, integer overflows
- **Configuration**: Excludes G115 (integer overflow false positives), G101 (handled separately)

### 2. **Python Code Security** (`security-python`) 
- **Tool**: `bandit` - Python security linter
- **Scans**: Python AST for security vulnerabilities
- **Detects**: SQL injection, shell injection, weak crypto, hardcoded passwords, unsafe YAML/pickle
- **Dependencies**: `safety` for Python package vulnerability scanning

### 3. **Secret Detection** (`security-secrets`)
- **Tool**: Custom regex-based scanner
- **Detects**: API keys, AWS credentials, JWT tokens, private keys, SSH keys, OAuth tokens, database URLs
- **Features**: High-entropy string detection, configuration file scanning
- **Patterns**: 12+ secret types with severity classification

### 4. **Dependency Vulnerabilities** (`security-deps`)
- **Tools**: `govulncheck`, `pip-audit`, `npm audit`
- **Scans**: Go modules, Python packages, Node.js dependencies
- **Detects**: Known CVEs in dependencies
- **Reports**: JSON/SARIF formats for CI integration

## 📊 Security Scoring System

| Score | Level | Status | Action Required |
|-------|-------|--------|-----------------|
| 90-100 | 🟢 EXCELLENT | All scans passed | Maintain current security posture |
| 70-89 | 🟡 GOOD | Minor issues | Review and address warnings |
| 50-69 | 🟠 MODERATE | Several issues | Plan remediation activities |
| 0-49 | 🔴 POOR | Critical issues | **IMMEDIATE ACTION REQUIRED** |

## 📁 Output Structure

**Clean Report Management**: Reports use fixed filenames and auto-overwrite on each run to prevent clutter.

```
tests/security/reports/
├── security_master_report_latest.txt       # Comprehensive report
├── security_master_report_latest.json      # JSON summary for automation
├── gosec_report_latest.{txt,json,sarif}    # Go security scan results
├── bandit_report_latest.{txt,json,csv}     # Python security scan results
├── secrets_report_latest.{txt,json}        # Secret detection results
└── dependency_report_latest.txt            # Dependency vulnerability results
```

**Note**: Each scan overwrites previous reports with the same name, maintaining only the latest results. No timestamped files accumulate.

## 🚨 Exit Codes & CI Integration

```bash
# Exit codes for automation
0  = ✅ PASS - No security issues found
1  = ⚠️  WARNINGS - Issues found, review recommended  
2  = 🚨 CRITICAL - Critical issues, immediate action required
>2 = ❌ FAILED - Scan execution failed
```

### CI/CD Integration Example

```yaml
# .github/workflows/security.yml
- name: Security Scan
  run: make security-full
  continue-on-error: false  # Fail build on security issues
  
- name: Upload Security Reports
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: security-reports
    path: tests/security/reports/
```

## 🔧 Advanced Usage

### Individual Scans

```bash
# Run specific scans
make security-go       # Go code only
make security-python   # Python code only  
make security-secrets  # Secret detection only
make security-deps     # Dependencies only
```

### Cleanup

```bash
# Clean all security reports
make security-clean
```

### Custom Configuration

Security scanners can be configured by modifying the scripts in `tests/security/`:

- `go-security.sh` - Modify `GOSEC_EXCLUDE` and `GOSEC_INCLUDE` rules
- `python-security.sh` - Update bandit configuration in generated `.bandit` file
- `secret-scan.sh` - Add custom secret patterns to `SECRET_PATTERNS` array
- `dependency-scan.sh` - Configure vulnerability thresholds and scanning scope

## 🛠 Tool Installation

All security tools are **automatically installed** when first run:

- `gosec`: `go install github.com/securego/gosec/v2/cmd/gosec@latest`
- `govulncheck`: `go install golang.org/x/vuln/cmd/govulncheck@latest`
- `bandit`: `pip install bandit[toml]`
- `pip-audit`: `pip install pip-audit`
- `safety`: `pip install safety`

## 📋 Security Best Practices

### 1. **Regular Scanning**
- Run `make security-full` before every commit
- Integrate into CI/CD pipeline
- Schedule weekly comprehensive scans

### 2. **Issue Remediation**
- **CRITICAL** issues: Fix immediately before deployment
- **WARNING** issues: Address within sprint cycle
- **INFO** issues: Consider during refactoring

### 3. **Dependency Management**
- Pin dependency versions in production
- Regular dependency updates with testing
- Monitor security advisories for dependencies

### 4. **Secret Management**
- Never commit secrets to version control
- Use environment variables or secret management tools
- Implement pre-commit hooks for secret detection

### 5. **Code Review**
- Include security considerations in code reviews
- Review security scan reports with team
- Document security decisions and exceptions

## 🚫 Known Limitations

### G115 Integer Overflow Rule
- **Issue**: Go 1.22+ causes widespread false positives
- **Solution**: Rule G115 is excluded by default
- **Monitoring**: Community working on improved implementation

### Python Virtual Environments
- **Requirement**: Some tools require proper Python environment setup
- **Solution**: Scripts auto-detect `.venv` and system Python
- **Fallback**: Manual tool installation instructions provided

## 📞 Support & Troubleshooting

### Common Issues

1. **"Tool not found" errors**
   - Solution: Scripts auto-install tools; ensure Go/Python are available

2. **Permission denied**
   - Solution: `chmod +x tests/security/*.sh`

3. **High false positive rate**
   - Solution: Review and customize scanner configurations

4. **Large codebases timeout**
   - Solution: Run individual scans (`security-go`, `security-python`, etc.)

### Debugging

```bash
# Enable verbose mode for debugging
DEBUG=1 make security-full

# Run individual scripts for detailed output
./tests/security/go-security.sh
./tests/security/secret-scan.sh
```

## 📊 Security Metrics & Reporting

The security framework generates metrics suitable for:

- **Security dashboards** (JSON output)
- **Compliance reporting** (detailed text reports)
- **Trend analysis** (timestamped reports)
- **CI/CD integration** (exit codes and SARIF format)

## 🎯 Zero-Tolerance Security Philosophy

This framework implements a **zero-tolerance approach** to security:

- ✅ **Proactive scanning** before vulnerabilities reach production
- 🚨 **Immediate alerts** for critical security issues
- 📊 **Comprehensive reporting** for audit and compliance
- 🔄 **Continuous monitoring** integrated into development workflow
- 🛡️ **Defense in depth** with multiple scanning tools and techniques

---

**Remember**: Security is not a one-time check but a continuous process. Regular scanning, prompt remediation, and security-conscious development practices are essential for maintaining a secure codebase. 