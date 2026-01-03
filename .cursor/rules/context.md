# AI Coding Agent Operating Contract (Edit-First)

## Purpose
This document defines the mandatory operating rules for any AI coding agent
(e.g., Cursor, Claude, Codex) working in this repository.

All agents MUST follow the Edit-First workflow and repository conventions
defined below.

---

## Primary Policy — Edit-First
- ALWAYS prefer `Edit()` over `Write()`.
- `Write()` is permitted ONLY when no existing file logically fits the change.
- Incremental changes are preferred over large rewrites.

---

## Execution Loop (Mandatory)
Before ANY change:
1. `Read()` relevant files
2. `Search()` for related logic
3. `Plan()` the minimal correct change
4. `Edit()` the existing file

After EACH change:
1. `Run()` tests
2. `Lint`
3. Summarize the change
4. Log progress in `progress.md`

---

## Repository Conventions
- Production code MUST live in one of:
  - `src/`
  - `code/`
  - `internal/`
  - parent directory
- Tests MUST live in:
  - `tests/`
  - or files matching `*_test.*`
- Every change SHOULD include a related test when feasible.

---

## Engineering Principles
The following principles ALWAYS apply:
- DRY (Don't Repeat Yourself)
- KISS (Keep It Simple)
- YAGNI (You Aren't Gonna Need It)
- SoC (Separation of Concerns)
- SRP (Single Responsibility Principle)

Additional rules:
- No dead code
- No unused or unreferenced imports
- Imports MUST NOT appear inside functions or methods
- Security ALWAYS takes priority over convenience

---

## Allowed Automation Tools
Only the following automation tools may be used:
- `files.read`
- `files.search`
- `git.diff`
- `git.apply`
- `process.run` (tests, linters, validation scripts)

---

## Logging & Traceability
- All meaningful progress MUST be documented in `progress.md`
- Planning decisions SHOULD be reflected in `PLAN.md`

