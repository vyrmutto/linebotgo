---
name: update-line-api
description: Update linebotgo private LINE API definitions when LINE changes their API. Use when encountering 4xx/5xx errors, or to analyze mitmproxy/Charles HTTP traffic dumps. Targets api/endpoints.go, api/types.go, and transport/http.go.
type: project
---

# update-line-api

Use this skill when:
- LINE returns unexpected 4xx/5xx errors from linebotgo
- User pastes an error log or stack trace
- User pastes HTTP traffic (mitmproxy/Charles/Proxyman export)
- User pastes a Thrift/protobuf schema diff or APK-decompiled IDL
- User runs `/update-line-api`

## Workflow

### Step 1 — Read Current State
```bash
# Read the primary target files first
cat api/endpoints.go
cat api/types.go
cat transport/http.go
```

### Step 2 — Analyze Input

**Error log input:**
- Find which endpoint returned the error (URL in log)
- Check if it's a: URL path change | header change | payload schema change | auth change
- Compare with constants in `api/endpoints.go`

**HTTP traffic dump (mitmproxy/Charles):**
- Extract: method, URL path, request headers, request body structure, response body
- Compare URL paths to `api/endpoints.go` constants
- Compare header names/values to `transport/http.go` constants
- Note new/changed request/response JSON fields → update `api/types.go`

**Thrift/protobuf schema diff:**
- Map changed field IDs/types to structs in `api/types.go`
- Add/rename/remove fields accordingly

### Step 3 — Make Changes (in priority order)

1. **`api/endpoints.go`** — URL or path changes
2. **`api/types.go`** — request/response struct field changes
3. **`transport/http.go`** — `LineApplication`, `LineUserAgent`, `LineAppVersion` constants (when LINE app version bumps)
4. **`auth/email.go`** or **`auth/qrcode.go`** — auth flow changes (endpoint or payload)

### Step 4 — Run Tests
```bash
go test ./... -v
```
Fix any failures before proceeding to Step 5.

### Step 5 — Update Docs
- Add entry under `## Unreleased` in `CHANGELOG.md`:
  ```markdown
  ## Unreleased
  ### Fixed
  - Update LINE API endpoint: `EndpointXxx` → new path
  ```
- If public API changed (new endpoints or types), update README
- Suggest version bump:
  - Patch: endpoint URL fix, header version bump
  - Minor: new endpoints or new message types added

## Reference
- Original Python implementation: https://github.com/fadhiilrachman/line-py
- LINE Android APK Thrift IDL can be extracted via jadx decompiler
- Headers in `transport/http.go` mirror the LINE Android app version string
