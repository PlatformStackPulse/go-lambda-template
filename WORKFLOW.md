# Repository Branch Protection & Workflow Guide

## GitHub Actions Status Checks

Configure the following status checks on your main branch:

### Required Status Checks

1. **CI Pipeline Checks:**
   - `CI Pipeline / Lint & Format Check`
   - `CI Pipeline / Test (1.22)`
   - `CI Pipeline / Test (1.23)`
   - `CI Pipeline / Security Scans`
   - `CI Pipeline / Commit Lint`
   - `CI Pipeline / Build`

2. **Code scanning check:**
   - `CodeQL Analysis / Analyze`

## Branch Protection Rules

### For `main` branch:

```yaml
# Require pull request reviews before merging
Require reviews: 1

# Dismiss stale pull request approvals
Dismiss stale PR approvals: true

# Require status checks to pass before merging
Require status checks:
  - CI Pipeline / Lint & Format Check
  - CI Pipeline / Test (1.22)
  - CI Pipeline / Test (1.23)
  - CI Pipeline / Security Scans
  - CI Pipeline / Commit Lint
  - CI Pipeline / Build
  - CodeQL Analysis / Analyze

# Require branches to be up to date before merging
Require branches up to date: true

# Include administrators
Include administrators: true

# Allow force pushes
Allow force pushes: false

# Allow deletions
Allow deletions: false

# Lockdown
Lockdown: false (or true for restrictive mode)
```

### For other branches:

- Allow direct commits to `develop` for minor updates
- Require PRs for feature branches

## Setup Instructions

1. **Go to Repository Settings** → **Branches**
2. **Click "Add rule"**
3. **Configure for `main` branch:**
   - Apply to administrators: ✅
   - Require pull request reviews: 1 review ✅
   - Dismiss stale reviews: ✅
  - Require status checks: all checks listed above ✅
   - Require branches up to date: ✅

## Quick Apply via API (Script)

You can apply the `main` branch protection policy in one command using:

1. Export token with admin access to this repo:

```bash
export GITHUB_TOKEN=ghp_xxx
```

2. Run script:

```bash
chmod +x scripts/apply-branch-protection.sh
scripts/apply-branch-protection.sh
```

Optional overrides:

```bash
GITHUB_OWNER=PlatformStackPulse GITHUB_REPO=go-template BRANCH=main scripts/apply-branch-protection.sh
```

## Quick Apply Checklist (GitHub UI)

Use this exact list when selecting required status checks for branch protection on `main`:

1. CI Pipeline / Lint & Format Check
2. CI Pipeline / Test (1.22)
3. CI Pipeline / Test (1.23)
4. CI Pipeline / Security Scans
5. CI Pipeline / Commit Lint
6. CI Pipeline / Build
7. CodeQL Analysis / Analyze

Recommended additional protection toggles:

1. Require a pull request before merging
2. Require approvals: 1
3. Dismiss stale pull request approvals when new commits are pushed
4. Require conversation resolution before merging
5. Require branches to be up to date before merging
6. Include administrators
7. Block force pushes
8. Block branch deletion

## Automatic Remediation Workflows

### 1. Auto-Update Dependencies
- **Trigger:** Weekly
- **Action:** Create PR with dependency updates
- **Config:** `.github/workflows/dependencies.yml`

### 2. Auto-Fix Formatting
- **Trigger:** PR submission
- **Action:** Suggest formatting fixes (not auto-commit)
- **Config:** CI Pipeline

### 3. Version Bumping
- **Trigger:** Manual (workflow_dispatch)
- **Action:** Bump version and create release
- **Config:** `.github/workflows/version-bump.yml`

## Recommended Workflow

```
main (protected)
└── develop (semi-protected)
    ├── feature/* (unprotected)
    ├── bugfix/* (unprotected)
    └── hotfix/* (unprotected)

PR Flow:
1. feature/* → PR → develop
2. develop → PR → main (requires approval + checks)
3. hotfix/* → PR → main (direct to main for urgent fixes)
```

## Additional Security Configurations

### Dependabot Settings

Enable in `Settings` → `Code security` → `Dependabot`:

- ✅ Dependabot alerts
- ✅ Dependabot security updates
- ✅ Dependabot version updates

Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
```

### Secret Scanning

Enable in `Settings` → `Security`:

- ✅ Push protection
- ✅ Secret scanning for partner patterns

### Code Scanning

- ✅ CodeQL enabled (see `.github/workflows/codeql.yml`)
- ✅ Alerts reviewed regularly

## Deployment Considerations

### Pre-Deployment Checklist

- [ ] All tests pass
- [ ] Security scans clear
- [ ] Coverage maintained (≥70%)
- [ ] Commits follow Conventional Commits
- [ ] PR has approval
- [ ] Branch is up to date with main

### Release Process

1. Create PR to main
2. Await reviews and checks
3. Merge to main
4. Tag with version (`git tag v1.2.3`)
5. Push tag (`git push origin v1.2.3`)
6. Release workflow auto-triggers
7. Artifacts published to GitHub Releases

For more information, see [CONTRIBUTING.md](CONTRIBUTING.md) and [README.md](README.md).
