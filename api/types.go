package api

import "time"

// MessageType represents the LINE message content type.
type MessageType string

const (
	MessageTypeText     MessageType = "NONE" // LINE uses NONE for plain text
	MessageTypeSticker  MessageType = "STICKER"
	MessageTypeImage    MessageType = "IMAGE"
	MessageTypeVideo    MessageType = "VIDEO"
	MessageTypeAudio    MessageType = "AUDIO"
	MessageTypeLocation MessageType = "LOCATION"
	MessageTypeContact  MessageType = "CONTACT"
)

// Message represents a LINE message (talk or OpenChat).
type Message struct {
	ID         string      `json:"id"`
	Type       MessageType `json:"type"`
	From       string      `json:"from"`
	To         string      `json:"to"`
	Text       string      `json:"text,omitempty"`
	ContentURL string      `json:"contentUrl,omitempty"`
	CreatedAt  time.Time   `json:"createdAt"`
}

// Contact represents a LINE contact/friend.
type Contact struct {
	MID           string `json:"mid"`
	DisplayName   string `json:"displayName"`
	StatusMessage string `json:"statusMessage"`
	PictureURL    string `json:"pictureUrl"`
}

// OpenChatInfo represents an OpenChat room.
type OpenChatInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MemberCount int    `json:"memberCount"`
}

// OpenChatMember represents a member of an OpenChat room.
type OpenChatMember struct {
	MID         string `json:"mid"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"` // "ADMIN" or "MEMBER"
}
