# Janitor's Journal

## 2026-01-21 - Fix Unused PORT Variables

**Issue:** Two `PORT` variables were defined but unused, leading to misleading logs ("Listening in port 0").
**Root Cause:** The flag parsing logic was updating `getlistener.PORT`, but the log message was using a local, zero-valued `PORT` variable. There was also an unused global `PORT` in the library package.
**Solution:** Removed the unused variables and updated the log message to reference the correct source of truth (`getlistener.PORT`).
**Pattern:** Always verify that variables used in log messages are the ones actually holding the configuration values, especially when using flag libraries.
