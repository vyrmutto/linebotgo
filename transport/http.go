package transport

import (
	"fmt"
	"net/http"
	"time"
)

// LINE Android client identifiers extracted from LINE APK 15.15.1 via jadx decompilation.
//
// X-Line-Application format (from Vt1/d.java, method i()):
//
//	strD + "\t" + appVersion + "\tAndroid OS\t" + Build.VERSION.RELEASE
//
// Where strD comes from C46197f.d() which returns "ANDROID" (primary) or "ANDROIDSECONDARY".
// App version from PackageManager.getPackageInfo().versionName.
//
// IMPORTANT: Build.VERSION.RELEASE is the actual Android OS version of the device at
// runtime — it is NOT hardcoded in the APK. Use NewHTTPClientWithOSVersion to specify
// the exact version string you want to emulate.
//
// To update when LINE bumps their app version:
//   - Check versionName in APK manifest: aapt dump badging line.apk | grep versionName
//   - Update LineAppVersion below; LineUserAgent is derived from it.
const (
	// LineAppVersion is the LINE app version from the APK manifest versionName.
	// Source: aapt dump badging → versionName="15.15.1"
	LineAppVersion = "15.15.1"

	// LineUserAgent is the User-Agent header.
	// Source: Vt1/d.java l(Context) — "Line/" + versionName
	LineUserAgent = "Line/" + LineAppVersion

	// LineSystemName is the OS name sent in X-Line-Application.
	// Source: Vt1/d.java i(Context) — literal string "Android OS"
	LineSystemName = "Android OS"

	// LineAndroidOSVersion is the Android OS version used by default.
	// This is NOT extracted from the APK — it is a runtime value (Build.VERSION.RELEASE).
	// "14" is a reasonable modern default but MUST match the device you are emulating.
	// Use NewHTTPClientWithOSVersion to override.
	LineAndroidOSVersion = "14"

	// LineApplication is the default X-Line-Application header value.
	// Override per-client with NewHTTPClientWithOSVersion if you need a specific OS version.
	LineApplication = "ANDROID\t" + LineAppVersion + "\t" + LineSystemName + "\t" + LineAndroidOSVersion
)

// BuildLineApplication returns the X-Line-Application header value for a given Android OS version.
// Use this when you need to emulate a specific device OS version.
// Source format: Vt1/d.java i(Context)
func BuildLineApplication(androidOSVersion string) string {
	return fmt.Sprintf("ANDROID\t%s\t%s\t%s", LineAppVersion, LineSystemName, androidOSVersion)
}

type HTTPClient struct {
	inner       *http.Client
	application string
}

// NewHTTPClient returns an HTTPClient using the default Android OS version.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		inner:       &http.Client{Timeout: 30 * time.Second},
		application: LineApplication,
	}
}

// NewHTTPClientWithOSVersion returns an HTTPClient that emulates a specific Android OS version.
// androidOSVersion should match Build.VERSION.RELEASE of the device being emulated (e.g. "11", "13", "14").
func NewHTTPClientWithOSVersion(androidOSVersion string) *HTTPClient {
	return &HTTPClient{
		inner:       &http.Client{Timeout: 30 * time.Second},
		application: BuildLineApplication(androidOSVersion),
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Line-Application", c.application)
	req.Header.Set("User-Agent", LineUserAgent)
	req.Header.Set("Content-Type", "application/x-thrift")
	req.Header.Set("Accept", "application/x-thrift")
	return c.inner.Do(req)
}
