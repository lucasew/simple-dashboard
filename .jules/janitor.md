## 2026-01-27 - Refactor section parsing and fix typos in block.go

### Issue
`SectionAsRenderBlock` has repetitive integer parsing logic and typo in error messages.

### Root Cause
Copy-paste coding for `size_x` and `size_y` parsing, and lack of spell check.

### Solution
Introduced `getInt` helper function and corrected "cant" to "can't".

### Pattern
Refactoring repetitive parsing logic.
