package transport

import (
	"net/http"
	"time"
)

const (
	LineApplication = "jp.naver.line.android"
	LineUserAgent   = "Line/13.18.1"
	LineSystemName  = "Android OS"
	LineAppVersion  = "13.18.1"
)

type HTTPClient struct {
	inner *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		inner: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Line-Application", LineApplication)
	req.Header.Set("User-Agent", LineUserAgent)
	req.Header.Set("X-Line-System-Name", LineSystemName)
	req.Header.Set("Content-Type", "application/x-thrift")
	req.Header.Set("Accept", "application/x-thrift")
	return c.inner.Do(req)
}
