package openchat

import (
	"encoding/json"
	"net/http"

	"github.com/vysina/linebotgo/api"
)

// PinMessage pins a message in an OpenChat room.
func (oc *OpenChat) PinMessage(chatID, messageID string) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatPin), map[string]string{
		"chatId":    chatID,
		"messageId": messageID,
	})
}

// SetNotification enables or disables notifications for an OpenChat room.
func (oc *OpenChat) SetNotification(chatID string, enabled bool) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatInfo)+"/notification", map[string]interface{}{
		"chatId":  chatID,
		"enabled": enabled,
	})
}

// GetChatInfo returns metadata about an OpenChat room.
func (oc *OpenChat) GetChatInfo(chatID string) (*api.OpenChatInfo, error) {
	req, err := http.NewRequest("GET", oc.resolveURL(api.EndpointOpenChatInfo)+"?chatId="+chatID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Line-Access", oc.c.AuthToken())
	resp, err := oc.c.HTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info api.OpenChatInfo
	return &info, json.NewDecoder(resp.Body).Decode(&info)
}
