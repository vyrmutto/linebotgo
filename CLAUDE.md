# linebotgo — Claude Instructions

## Project Context

Go SDK wrapping LINE's private (non-official) Thrift RPC API. All API constants are reverse-engineered from the LINE Android APK via jadx decompilation.

Module: `github.com/vysina/linebotgo`
LINE APK version: 15.15.1 (armeabi-v7a)
Decompiled source: `/tmp/line_decompile/sources/`

---

## CRITICAL RULE: Never Guess API Values

**Do not guess, estimate, or infer any of the following from general knowledge:**

- API endpoint paths (e.g. `/S4`, `/SQ1`, `/acct/pais/v1`)
- HTTP header values (`X-Line-Application`, `User-Agent`)
- Thrift struct field IDs (field numbers used in binary serialization)
- Enum wire values (integers sent over the wire in Thrift binary)
- Server hostnames
- Timeout values

**Why this matters:** Integration tests that receive a wrong interface fail silently or produce corrupt binary payloads. A wrong Thrift field ID means a structurally valid but semantically incorrect message — no compile error, no obvious runtime error, but the server rejects it or returns garbage.

If the real value is unknown, say so explicitly. Do not substitute a plausible-looking value.

---

## How to Find Real Values

All values must come from one of these two sources:

### 1. APK Decompilation (primary)

Decompiled sources are at `/tmp/line_decompile/sources/`. Key files:

| What you need | Where to look |
|---|---|
| Server hostnames | `assets/default-connection-info/default.json` |
| Thrift service paths | `zi/EnumC55943a.java` |
| X-Line-Application format | `Vt1/d.java` method `i()` |
| Device type string ("ANDROID") | `qZ/C46197f.java` method `d()` |
| TalkService method names + field IDs | `jp/naver/line/android/thrift/client/impl/LegacyTalkServiceClientImpl.java` |
| SquareService structs + field IDs | `com/linecorp/square/protocol/thrift/SquareService.java` |
| Domain model structs (Message, Contact, etc.) | `lv1/*.java` (obfuscated package) |

To re-decompile: `jadx -d /tmp/line_decompile /path/to/line.apk`

### 2. Real traffic capture (secondary)

Use mitmproxy with SSL unpinning (Frida) on an Android device with the LINE app. Captured traffic reveals the exact binary payloads and any values not visible from static analysis.

---

## Known Values That Are Still Uncertain

| Value | Location | Issue |
|---|---|---|
| Android OS version in `LineApplication` | `transport/http.go` | `LineAndroidOSVersion = "14"` is a default, not extracted from APK. It is `Build.VERSION.RELEASE` at runtime. Use `NewHTTPClientWithOSVersion("X")` to set the correct value for your target device. |
| Auth request/response Thrift format | `auth/` package | Login flow structs are heavily obfuscated — exact field IDs not confirmed from APK source |

These must be resolved before writing integration tests that depend on them.

---

## Confirmed Sources for Current Constants

| Constant | Source file | Extraction method |
|---|---|---|
| `BaseURL = "https://legy.line-apps.com"` | `assets/default-connection-info/default.json` | Direct read |
| `ObjectStorageURL = "https://obs.line-apps.com"` | same | Direct read |
| All 46 service paths in `api/endpoints.go` | `zi/EnumC55943a.java` (complete enum, all entries) | Direct read |
| `EndpointCertificate = "/Q"` | `zi/V.java` SECONDARY_DEVICE_LOGIN_VERIFY_PIN | Direct read |
| `EndpointCertificateE2E = "/LF1"` | `zi/V.java` SECONDARY_DEVICE_LOGIN_VERIFY_PIN_WITH_E2EE | Direct read |
| `EndpointOpenChat = "/SQ1"` | `zi/EnumC55943a.java` SQUARE entry | Direct read (was `/SQS1` in older versions) |
| `EndpointOpenChatBot = "/BP1"` | `zi/EnumC55943a.java` SQUARE_BOT entry (index 34) | Direct read |
| `EndpointUnifiedShop = "/TSHOP4"` | `zi/EnumC55943a.java` UNIFIED_SHOP entry (index 13) | Direct read |
| `LineApplication` format | `Vt1/d.java` method `i()` | Direct read |
| `LineUserAgent = "Line/15.15.1"` | `Vt1/d.java` method `l()` | Direct read |
| `LineAppVersion = "15.15.1"` | APK manifest `versionName` via `aapt dump badging` | Direct read |
| `SyncRequest` field IDs in `api/types.go` | `lv1/H8.java` | Direct read |
| Thrift struct field IDs (`Message`, `Operation`, etc.) | `lv1/*.java` decompiled model classes | Direct read |
