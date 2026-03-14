package api_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/api"
)

func TestEndpointsHaveBaseURL(t *testing.T) {
	endpoints := []string{
		api.EndpointEmailLogin,
		api.EndpointFetchOps,
		api.EndpointSendMessage,
		api.EndpointOpenChatSearch,
		api.EndpointOpenChatJoin,
	}
	for _, ep := range endpoints {
		assert.True(t, strings.HasPrefix(ep, api.BaseURL),
			"endpoint %q should start with BaseURL", ep)
	}
}
