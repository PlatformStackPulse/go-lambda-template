# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability, please use GitHub Security Advisories for this repository instead of opening a public issue.

Please include:

- Description of the vulnerability
- Steps to reproduce (if applicable)
- Potential impact
- Suggested fix (if any)

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | ✅ Yes    |
| N-1     | ✅ Yes    |
| Older   | ❌ No     |

## Security Scanning

This project uses the following security scanning tools:

- **gosec** — Go security checker (detects unsafe code patterns)
- **govulncheck** — Vulnerability checker for Go dependencies
- **CodeQL** — SIEM code analysis and vulnerability detection
- **Dependabot** — Automated dependency updates and vulnerability alerts

### Running Locally

```bash
# Security scan (gosec)
make security

# Vulnerability check (govulncheck)
make sec-update

# All security and dev setup
make dev-setup      # Installs all tools
make security       # Run security scan
make sec-update     # Check for vulnerabilities
```

### CI/CD Integration

Security scans run automatically on:
- Every pull request
- Before merge to main
- Weekly scheduled scan
- On push to main

## Disclosure Timeline

- Notify maintainers
- Wait for acknowledgment (within 48 hours)
- Provide reasonable time for fix (typically 30-90 days)
- Coordinated disclosure

## Compliance

- Follows responsible disclosure practices
- Reports processed with urgency
- Security patches released as soon as possible
