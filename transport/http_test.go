package transport_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/transport"
)

func TestHTTPClientSetsLINEHeaders(t *testing.T) {
	var gotHeaders http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := transport.NewHTTPClient()
	req, _ := http.NewRequest("GET", srv.URL, nil)
	c.Do(req)

	assert.Equal(t, transport.LineApplication, gotHeaders.Get("X-Line-Application"))
	assert.Equal(t, transport.LineUserAgent, gotHeaders.Get("User-Agent"))
	// Android LINE app does not send X-Line-Carrier; removed in v15
	assert.Empty(t, gotHeaders.Get("X-Line-Carrier"))
}

func TestHTTPClientWithOSVersionSetsCorrectHeader(t *testing.T) {
	var gotHeaders http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := transport.NewHTTPClientWithOSVersion("13")
	req, _ := http.NewRequest("GET", srv.URL, nil)
	c.Do(req)

	want := transport.BuildLineApplication("13")
	assert.Equal(t, want, gotHeaders.Get("X-Line-Application"))
	assert.Contains(t, gotHeaders.Get("X-Line-Application"), "\t13")
}

func TestBuildLineApplication(t *testing.T) {
	got := transport.BuildLineApplication("11")
	assert.Equal(t, "ANDROID\t15.15.1\tAndroid OS\t11", got)
}
