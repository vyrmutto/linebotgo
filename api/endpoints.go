// Package api contains LINE private API endpoint definitions and shared types.
// IMPORTANT: This file is the primary target for the update-line-api Claude skill.
// When LINE changes API endpoints, headers, or parameters — update this file first.
package api

const (
	BaseURL = "https://gw.line.naver.jp"

	// Auth endpoints
	EndpointEmailLogin   = BaseURL + "/api/v4/TalkService/loginWithIdentityCredential"
	EndpointQRCreate     = BaseURL + "/api/v4/TalkService/getAuthQrCode"
	EndpointQRPoll       = BaseURL + "/api/v4/TalkService/checkAuthQrCodeVerified"
	EndpointRefreshToken = BaseURL + "/api/v4/TalkService/refreshToken"

	// Talk / Messaging endpoints
	EndpointFetchOps    = BaseURL + "/SYNC/v3/TalkService/fetchOps"
	EndpointSendMessage = BaseURL + "/S4/TalkService/sendMessage"
	EndpointGetContacts = BaseURL + "/S4/TalkService/getAllContactIds"
	EndpointGetProfile  = BaseURL + "/S4/TalkService/getProfile"

	// OpenChat endpoints
	EndpointOpenChatSearch  = BaseURL + "/api/v1/square/search"
	EndpointOpenChatJoin    = BaseURL + "/api/v1/square/join"
	EndpointOpenChatLeave   = BaseURL + "/api/v1/square/leave"
	EndpointOpenChatSend    = BaseURL + "/api/v1/square/message/send"
	EndpointOpenChatMembers = BaseURL + "/api/v1/square/members"
	EndpointOpenChatKick    = BaseURL + "/api/v1/square/members/kick"
	EndpointOpenChatBan     = BaseURL + "/api/v1/square/members/ban"
	EndpointOpenChatPin     = BaseURL + "/api/v1/square/message/pin"
	EndpointOpenChatInfo    = BaseURL + "/api/v1/square/info"
)
