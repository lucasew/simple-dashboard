## 2026-01-19 - Fix unused PORT variable and incorrect logging

**Issue:** Unused `PORT` variable in `main.go` and `cmd/simple-dashboardd/main.go`. The latter caused the log message to always report "Listening in port 0".
**Root Cause:** The `PORT` variable was declared but shadowed or ignored in favor of `getlistener.PORT`, leading to dead code and incorrect runtime logging.
**Solution:** Removed the unused `PORT` variables and updated the log message to use the correct `getlistener.PORT` variable.
**Pattern:** Always check if variables declared for flags are actually used or if the flag library binds to a different variable.
