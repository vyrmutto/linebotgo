package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vysina/linebotgo/api"
)

// Listen starts long polling and dispatches incoming messages to OnMessage handlers.
// Blocks until ctx is cancelled. Retries on transient errors.
func (c *Client) Listen(ctx context.Context) error {
	revision := int64(0)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ops, err := c.fetchOps(ctx, revision)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
				continue
			}
		}

		for _, op := range ops {
			if op.Message != nil {
				c.dispatchMessage(op.Message)
			}
			if op.Revision > revision {
				revision = op.Revision
			}
		}
	}
}

type operation struct {
	Revision int64
	Message  *Message
}

func (c *Client) fetchOps(ctx context.Context, revision int64) ([]operation, error) {
	url := api.EndpointFetchOps
	if c.baseURL != "" {
		url = c.baseURL + "/fetchOps"
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Line-Access", c.AuthToken())
	q := req.URL.Query()
	q.Set("revision", fmt.Sprintf("%d", revision))
	q.Set("count", "50")
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Operations []struct {
			Revision int64    `json:"revision"`
			Message  *Message `json:"message,omitempty"`
		} `json:"operations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ops := make([]operation, len(result.Operations))
	for i, o := range result.Operations {
		ops[i] = operation{Revision: o.Revision, Message: o.Message}
	}
	return ops, nil
}
