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

**IMPORTANT:** Obfuscated class names (`lv1/Q6.java`, `Vt1/d.java`, `qZ/C46197f.java`, etc.)
change with every APK build. Always find files by grep-ing for stable string content,
not by filename. Only asset paths and full-package class names are stable across builds.

---

## Step 0 — Download and Decompile the New APK

```bash
# Check current version in transport/http.go
grep LineAppVersion transport/http.go

# Decompile the new APK
jadx -d /tmp/line_decompile /path/to/line.apk

# Verify the app version from the APK manifest
aapt dump badging /path/to/line.apk | grep versionName
```

Decompiled sources will be at `/tmp/line_decompile/sources/`.

---

## Step 1 — Extract Server Hostnames

**Stable asset path — filename does not change between builds.**

```bash
cat "/tmp/line_decompile/sources/assets/default-connection-info/default.json"
```

Look for:
- `servers.legy.mobile[0].host` → `BaseURL` in `api/endpoints.go`
- `servers.obs.mobile[0].host` → `ObjectStorageURL` in `api/endpoints.go`

---

## Step 2 — Extract All Service Endpoint Paths

**Find the service path enum by searching for known stable path strings:**

```bash
# The enum contains all Thrift service paths as string literals
grep -rn '"/S4"\|"/P4"\|"/SQ1"\|"LONG_POLLING"\|"NORMAL"' \
  /tmp/line_decompile/sources/ | grep -v '.class:' | head -5
```

This will reveal the file (e.g. `zi/EnumC55943a.java` or whatever it's named this build).
Then dump the full file to extract all entries:

```bash
# Replace FILE with the path found above
grep 'new EnumC' FILE
```

Each entry has the form:
```java
new EnumCXXXXX(INDEX, "SERVICE_NAME", "/path", timeoutMs)
```

Map every path to the appropriate constant in `api/endpoints.go`.

**Also find the secondary device login paths:**

```bash
# Search for the secondary device login enum by its stable string content
grep -rn '"verify_pin"\|SECONDARY_DEVICE_LOGIN' \
  /tmp/line_decompile/sources/ | grep -v '.class:' | head -5
```

Look for `prefixUrl` values `/Q` and `/LF1`.

---

## Step 3 — Extract HTTP Header Constants

**Find the X-Line-Application construction logic:**

```bash
# The header is built by concatenating device type + version + "Android OS" + Build.VERSION.RELEASE
grep -rn 'Android OS\|Build.VERSION.RELEASE\|X-Line-Application' \
  /tmp/line_decompile/sources/ | grep -v '.class:' | grep -v 'test\|Test' | head -10
```

Reveals the file containing method `i()`. Open it and read:
- The format string: `strD + "\t" + appVersion + "\tAndroid OS\t" + Build.VERSION.RELEASE`

**Find the device type string:**

```bash
# Search for the literal "ANDROID" device type constant
grep -rn '"ANDROIDSECONDARY"\|"ANDROID"' \
  /tmp/line_decompile/sources/ | grep -v '.class:' | grep -v 'import\|//\|test' | head -10
```

Should return `"ANDROID"` for primary device and `"ANDROIDSECONDARY"` for secondary.

Update in `transport/http.go`:
- `LineAppVersion` — from `aapt dump badging` versionName (Step 0)
- `LineUserAgent` — derived: `"Line/" + LineAppVersion`
- `LineApplication` — derived from LineAppVersion + "Android OS" + LineAndroidOSVersion

**Note:** `LineAndroidOSVersion` is `Build.VERSION.RELEASE` at runtime — it is NOT in the APK.
Keep the default as a common recent Android version. Never hardcode a specific device's version.

---

## Step 4 — Extract Thrift Struct Field IDs

Structs are in an obfuscated package (currently `lv1/`). Find them by known field names:

```bash
# Find Message struct — search for its stable Thrift field names
grep -rn '"text"\|"contentType"\|"toType"\|"createdTime"' \
  /tmp/line_decompile/sources/ | grep 'jy1.c f\|new jy1.c' | head -10

# Find Operation struct
grep -rn '"revision"\|"param1"\|"param2"' \
  /tmp/line_decompile/sources/ | grep 'new jy1.c' | head -10

# Find SyncRequest struct
grep -rn '"lastRevision"\|"lastGlobalRevision"\|"fullSyncRequestReason"' \
  /tmp/line_decompile/sources/ | grep 'new jy1.c' | head -10

# Find Profile struct
grep -rn '"displayName"\|"pictureStatus"\|"statusMessage"' \
  /tmp/line_decompile/sources/ | grep 'new jy1.c' | head -10
```

Each Thrift field is declared as:
```java
public static final jy1.c fXXX = new jy1.c("fieldName", (byte) WIRE_TYPE, FIELD_ID);
```

Thrift wire types: 2=BOOL, 3=I8, 6=I16, **8=I32**, **10=I64**, **11=STRING**, **12=STRUCT**, 13=MAP, 14=SET, 15=LIST

Once you find the file, open it to read all field IDs. Update structs in `api/types.go`.

---

## Step 5 — Check for New Enum Values

```bash
# Find OpType enum — search for stable value names
grep -rn '"END_OF_OPERATION"\|"RECEIVE_MESSAGE"\|"SEND_MESSAGE"' \
  /tmp/line_decompile/sources/ | grep -v '.class:' | head -5
# Then open that file and count entries vs current api/types.go

# Find ContentType enum
grep -rn '"NONE"\|"IMAGE"\|"VIDEO"\|"AUDIO"' \
  /tmp/line_decompile/sources/ | grep 'new jy1.c\|ContentType\|getValue' | head -10
```

If the APK has more enum entries, add them to `api/types.go`.
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
git add api/endpoints.go api/types.go transport/http.go
git commit -m "feat: update LINE API constants for APK vX.XX.X"
```

---

## Stable vs Obfuscated Paths

| Type | Example | Stable? |
|---|---|---|
| Asset files | `assets/default-connection-info/default.json` | ✅ stable |
| Full package classes | `com/linecorp/square/protocol/thrift/common/Square.java` | ✅ stable |
| Full package classes | `jp/naver/line/android/thrift/client/impl/LegacyTalkServiceClientImpl.java` | ✅ stable |
| Obfuscated short classes | `lv1/Q6.java`, `Vt1/d.java`, `qZ/C46197f.java`, `zi/EnumC55943a.java` | ❌ changes every build |

Always locate obfuscated files by grep-ing for stable string content inside them,
not by filename. The source file map below shows which files were used **last build**
as a starting point, but verify each one before trusting it.

## Source File Map (from APK 15.15.1 — verify on each new build)

| What | Last known file |
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
