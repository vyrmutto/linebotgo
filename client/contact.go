package client

import (
	"encoding/json"
	"net/http"

	"github.com/vysina/linebotgo/api"
)

// GetProfile returns profile information for a given LINE MID.
func (c *Client) GetProfile(mid string) (*api.Contact, error) {
	url := api.EndpointGetProfile + "?mid=" + mid
	if c.baseURL != "" {
		url = c.baseURL + "/profile?mid=" + mid
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Line-Access", c.AuthToken())
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var contact api.Contact
	return &contact, json.NewDecoder(resp.Body).Decode(&contact)
}
