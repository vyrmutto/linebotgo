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

	assert.Equal(t, "jp.naver.line.android", gotHeaders.Get("X-Line-Application"))
	assert.NotEmpty(t, gotHeaders.Get("User-Agent"))
}
