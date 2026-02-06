# Janitor Journal

## 2025-05-23 (Arrumador)
- Configured `mise.toml` with standard tasks (`lint`, `fmt`, `test`, `install`, `codegen`, `ci`).
- Added `shellcheck` linter.
- Replaced `.github/workflows/autorelease.yaml` with a standardized workflow using `mise` and `jdx/mise-action`.
- Created `scripts/build_release.sh` to encapsulate build logic.
- Fixed shell script issues in `make_release` and `scripts/build_release.sh`.
