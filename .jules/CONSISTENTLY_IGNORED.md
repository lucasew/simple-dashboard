# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

## IGNORE: Removing unused PORT variable

**- Pattern:** Removing the `PORT` variable in `dashboard.go` (formerly `main.go`) or `cmd/simple-dashboardd/main.go`.
**- Justification:** This change has been consistently rejected (PRs #29, #31, #36, #46). Although unused, the variable is intended to remain in the library file.
**- Files Affected:** `dashboard.go`, `cmd/simple-dashboardd/main.go`, `main.go`

## IGNORE: JSON Marshal for XSS in ServeHTTP

**- Pattern:** Using `json.Marshal` to escape error strings injected into `alert()` within `ServeHTTP`.
**- Justification:** Repeatedly rejected (PRs #30, #32, #46). This specific method of fixing XSS or the reliance on `alert()` is not desired.
**- Files Affected:** `dashboard.go`, `main.go`

## IGNORE: Standardizing CI with Mise

**- Pattern:** Replacing explicit shell steps in `.github/workflows/autorelease.yaml` with `mise run` commands, or migrating to use mise action.
**- Justification:** Large infrastructure changes and tooling standardization have been rejected (PRs #34, #36, #38, #41, #46). The project prefers explicit build steps in CI.
**- Files Affected:** `.github/workflows/autorelease.yaml`, `mise.toml`

## IGNORE: Renaming GitHub Actions Workflow File

**- Pattern:** Renaming `.github/workflows/autorelease.yaml` to `.github/workflows/autorelease.yml`.
**- Justification:** Repeatedly rejected (PRs #34, #36, #37, #41, #42, #46). The project strictly uses the `.yaml` extension for its GitHub Actions workflow file.
**- Files Affected:** `.github/workflows/autorelease.yaml`

## IGNORE: Cosmetic modifications to error messages

**- Pattern:** Fixing typos or altering wording in error messages (e.g., "cant" to "can't", "invalid label template" to "invalid background_color template").
**- Justification:** Rejected (PRs #35, #43). Cosmetic changes to error strings are consistently ignored and should not be proposed.
**- Files Affected:** `block.go`, `dashboard.go`

## IGNORE: Refactoring Section Parsing Logic

**- Pattern:** Extracting integer parsing logic for `size_x` and `size_y` into a helper function (e.g., `getInt`) in `block.go` or making dashboard sections deterministic.
**- Justification:** Refactoring repetitive parsing logic or ordering maps deterministic was rejected (PR #35, #44). The straightforward parsing is preferred over abstraction.
**- Files Affected:** `block.go`, `dashboard.go`, `dashboard_test.go`

## IGNORE: Bundling unrelated CI fixes

**- Pattern:** Including unrelated CI workflow repairs (e.g., changing how `mise` is installed) in Pull Requests meant for documentation, features, or bug fixes.
**- Justification:** Rejected (PRs #37, #42, #44). The project strictly requires separating concerns; unrelated CI fixes must be in their own distinct PRs.
**- Files Affected:** `.github/workflows/autorelease.yaml`

## IGNORE: Adding HSTS Security Header

**- Pattern:** Adding the `Strict-Transport-Security` (HSTS) header.
**- Justification:** This change has been proposed multiple times and rejected. It is considered an operational risk because it can break subdomains that do not support HTTPS.
**- Files Affected:** `src/hooks.server.ts`

## IGNORE: Updating GitHub Action version

**- Pattern:** Upgrading actions checkout versions to v6.
**- Justification:** This change is rejected. Don't downgrade or upgrade github action versions automatically without user instructions.
**- Files Affected:** `.github/workflows/autorelease.yaml`
