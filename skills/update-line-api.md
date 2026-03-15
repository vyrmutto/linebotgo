---
name: update-line-api
description: Update linebotgo private LINE API definitions when LINE releases a new APK version. Extraction method is APK decompilation via jadx ONLY — never guess or infer values.
type: project
---

# update-line-api

Use this skill when:
- LINE releases a new APK and API constants may have changed
- LINE returns unexpected 4xx/5xx errors suggesting endpoint or header changes
- User runs `/update-line-api`

**METHOD: APK decompilation only.** Never guess, infer, or copy from external sources.
Wrong constants produce silent integration failures (valid binary, wrong semantics).

---

## Step 0 — Download and Decompile the New APK

```bash
# Check current version in transport/http.go
grep LineAppVersion transport/http.go

# Decompile the new APK (replace path with actual APK location)
jadx -d /tmp/line_decompile /path/to/line.apk

# Verify the app version from APK manifest
aapt dump badging /path/to/line.apk | grep versionName
```

The decompiled sources will be at `/tmp/line_decompile/sources/`.

---

## Step 1 — Extract Server Hostnames

**Source file:** `assets/default-connection-info/default.json`

```bash
cat "/tmp/line_decompile/sources/assets/default-connection-info/default.json"
```

Look for:
- `servers.legy.mobile[0].host` → `BaseURL` in `api/endpoints.go`
- `servers.obs.mobile[0].host` → `ObjectStorageURL` in `api/endpoints.go`

---

## Step 2 — Extract All Service Endpoint Paths

**Source file:** `zi/EnumC55943a.java`

```bash
grep 'new EnumC55943a' /tmp/line_decompile/sources/zi/EnumC55943a.java
```

Each entry has the form:
```java
new EnumC55943a(INDEX, "SERVICE_NAME", "/path", timeoutMs)
```

Map every path to the appropriate constant in `api/endpoints.go`.
Also check `zi/V.java` for secondary device login paths (`/Q`, `/LF1`):

```bash
grep 'prefixUrl\|SECONDARY_DEVICE' /tmp/line_decompile/sources/zi/V.java
```

---

## Step 3 — Extract HTTP Header Constants

**Source file:** `Vt1/d.java`

```bash
grep -A5 'strD\|versionName\|Line/\|Android OS' /tmp/line_decompile/sources/Vt1/d.java | head -40
```

- Method `i(Context)` → X-Line-Application format: `strD + "\t" + appVersion + "\tAndroid OS\t" + Build.VERSION.RELEASE`
- Method `l(Context)` → User-Agent: `"Line/" + versionName`

**Source file:** `qZ/C46197f.java` (device type string)

```bash
grep -A10 'public final String d()' /tmp/line_decompile/sources/qZ/C46197f.java
```

Returns `"ANDROID"` for primary device, `"ANDROIDSECONDARY"` for secondary.

Update in `transport/http.go`:
- `LineAppVersion` — from `aapt dump badging` versionName
- `LineUserAgent` — derived from LineAppVersion (`"Line/" + version`)
- `LineApplication` — derived from LineAppVersion + LineSystemName + LineAndroidOSVersion

**Note:** `LineAndroidOSVersion` is `Build.VERSION.RELEASE` at runtime — it is NOT in the APK.
Keep the default as a common recent Android version. Never hardcode a specific device's version.

---

## Step 4 — Extract Thrift Struct Field IDs

Domain model structs are in the `lv1/` package (obfuscated class names).
Find them by searching for known field names:

```bash
# Find Message struct
grep -rn '"id"\|"text"\|"contentType"\|"toType"' /tmp/line_decompile/sources/lv1/ | grep 'jy1.c f' | head -20

# Find Operation struct
grep -rn '"revision"\|"type"\|"param1"' /tmp/line_decompile/sources/lv1/ | grep 'jy1.c f' | head -20

# Find SyncRequest struct
grep -rn '"lastRevision"\|"count"\|"lastGlobalRevision"' /tmp/line_decompile/sources/lv1/ | head -10
```

Each field is declared as:
```java
public static final jy1.c fXXX = new jy1.c("fieldName", (byte) WIRE_TYPE, FIELD_ID);
```

Thrift wire types: 2=BOOL, 3=I8, 6=I16, 8=I32, 10=I64, 11=STRING, 12=STRUCT, 13=MAP, 14=SET, 15=LIST

Update structs in `api/types.go` with new field IDs.

---

## Step 5 — Check for New Enum Values

```bash
# OpType enum (operation types in fetchOperations stream)
wc -l /tmp/line_decompile/sources/lv1/Q6.java
grep -c 'OpType\|= OpType' /Users/vysina/my_workspace/linebotgo/api/types.go

# ContentType enum
grep -c 'new.*EnumC\|= new' /tmp/line_decompile/sources/lv1/H3.java
```

Compare counts. If the APK has more entries, add the new ones to `api/types.go`.
Wire integer values MUST be read directly from the Java source — never assumed.

---

## Step 6 — Make Changes

Priority order:

1. **`transport/http.go`** — `LineAppVersion` (if LINE app version changed)
2. **`api/endpoints.go`** — paths from Step 1 & 2
3. **`api/types.go`** — new/changed struct fields, new enum values

---

## Step 7 — Run Tests

```bash
go test ./...
```

Fix all failures before committing.

---

## Step 8 — Commit

```bash
git add api/endpoints.go api/types.go transport/http.go transport/http_test.go
git commit -m "feat: update LINE API constants for APK vX.XX.X"
```

---

## Source File Map (quick reference)

| What | APK source file |
|---|---|
| Server hostnames | `assets/default-connection-info/default.json` |
| All Thrift service paths | `zi/EnumC55943a.java` |
| Secondary device login paths | `zi/V.java` |
| X-Line-Application format | `Vt1/d.java` method `i()` |
| User-Agent format | `Vt1/d.java` method `l()` |
| Device type string ("ANDROID") | `qZ/C46197f.java` method `d()` |
| OpType enum | `lv1/Q6.java` |
| ContentType enum | `lv1/H3.java` |
| MIDType enum | `lv1/C6.java` |
| Message struct field IDs | `lv1/E6.java` |
| Operation struct field IDs | `lv1/R6.java` |
| SyncRequest struct field IDs | `lv1/H8.java` |
| Profile struct field IDs | `lv1/C41072k7.java` |
| Contact struct field IDs | `lv1/C41249x3.java` |
| LoginResult struct field IDs | `lv1/C41280z6.java` |
| SquareService structs | `com/linecorp/square/protocol/thrift/common/*.java` |
| TalkService client impl | `jp/naver/line/android/thrift/client/impl/LegacyTalkServiceClientImpl.java` |
