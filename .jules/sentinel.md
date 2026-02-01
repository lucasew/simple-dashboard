## 2026-02-01 - Fix XSS in block rendering

**Vulnerability:** The application used `text/template` to render dashboard blocks. This library does not escape HTML output, allowing Stored XSS if configuration contains malicious templates or Reflected/Stored XSS if templates interpolate user-controlled data (e.g., usernames).

**Learning:** `text/template` should never be used to generate HTML.

**Prevention:** Switch to `html/template`, which automatically escapes variables based on context (HTML body, attributes, etc.).
