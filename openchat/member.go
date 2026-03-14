package openchat

import (
	"encoding/json"
	"net/http"

	"github.com/vysina/linebotgo/api"
)

// GetMembers returns all members of an OpenChat room.
func (oc *OpenChat) GetMembers(chatID string) ([]api.OpenChatMember, error) {
	req, err := http.NewRequest("GET", oc.resolveURL(api.EndpointOpenChatMembers)+"?chatId="+chatID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Line-Access", oc.c.AuthToken())
	resp, err := oc.c.HTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Members []api.OpenChatMember `json:"members"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Members, nil
}

// Kick removes a member from an OpenChat room.
func (oc *OpenChat) Kick(chatID, memberMID string) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatKick), map[string]string{
		"chatId": chatID,
		"mid":    memberMID,
	})
}

// Ban bans a member from an OpenChat room.
func (oc *OpenChat) Ban(chatID, memberMID string) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatBan), map[string]string{
		"chatId": chatID,
		"mid":    memberMID,
	})
}
