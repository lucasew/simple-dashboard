# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

## IGNORE: Removing unused PORT variable

**- Pattern:** Removing the `PORT` variable in `dashboard.go`.
**- Justification:** This change has been consistently rejected (PRs #29, #31, #36). Although unused, the variable is intended to remain in the library file.
**- Files Affected:** `dashboard.go` (formerly `main.go`)

## IGNORE: JSON Marshal for XSS in ServeHTTP

**- Pattern:** Using `json.Marshal` to escape error strings injected into `alert()` within `ServeHTTP`.
**- Justification:** Repeatedly rejected (PRs #30, #32). This specific method of fixing XSS or the reliance on `alert()` is not desired.
**- Files Affected:** `dashboard.go` (formerly `main.go`)

## IGNORE: Standardizing CI with Mise

**- Pattern:** Replacing explicit shell steps in `.github/workflows/autorelease.yaml` with `mise run` commands, or renaming the workflow file to `.yml`.
**- Justification:** Large infrastructure changes and tooling standardization have been rejected (PRs #34, #36, #38). The project prefers explicit build steps in CI.
**- Files Affected:** `.github/workflows/autorelease.yaml`, `mise.toml`
