// Package api contains LINE private API endpoint definitions and shared types.
// IMPORTANT: This file is the primary target for the update-line-api Claude skill.
//
// Values extracted from LINE Android APK 15.15.1 via jadx decompilation:
//   - Server host: assets/default-connection-info/default.json → legy.line-apps.com
//   - Service paths: zi/EnumC55943a.java (complete service path enum, 46 entries)
//   - Secondary device login: zi/V.java
//   - OpenChat path: "/SQ1" (changed from legacy "/SQS1")
//
// To update when LINE changes their API:
//  1. Run the update-line-api skill with mitmproxy traffic, OR
//  2. Decompile the new APK with jadx and re-read this file's sources
package api

const (
	// BaseURL is the primary LINE server (LEGY — LINE Edge Gateway, HTTP/2).
	// Source: assets/default-connection-info/default.json → servers.legy.mobile[0].host
	BaseURL = "https://legy.line-apps.com"

	// ObjectStorageURL is used for sending/receiving media.
	// Source: assets/default-connection-info/default.json → servers.obs.mobile[0].host
	ObjectStorageURL = "https://obs.line-apps.com"

	// ── Auth / Registration ────────────────────────────────────────────────
	// Source: zi/EnumC55943a.java

	// EndpointLogin is the legacy REGISTRATION path used in the primary login flow.
	EndpointLogin = BaseURL + "/api/v4/TalkService.do" // enum index 6, REGISTRATION

	// EndpointNewRegistration is the new account registration/login endpoint (acct).
	EndpointNewRegistration = BaseURL + "/acct/pais/v1" // enum index 41, NEW_REGISTRATION

	// EndpointSecondaryQRLogin is used for QR-code login on secondary devices.
	EndpointSecondaryQRLogin = BaseURL + "/ACCT/lgn/sq/v1" // enum index 42, SECONDARY_QR_LOGIN

	// EndpointCertificate is the secondary device PIN verification endpoint.
	// Source: zi/V.java — SECONDARY_DEVICE_LOGIN_VERIFY_PIN (prefixUrl "/Q")
	// NOTE: zi/V.java also has SECONDARY_DEVICE_LOGIN_VERIFY_PIN_WITH_E2EE at "/LF1"
	EndpointCertificate    = BaseURL + "/Q"   // secondary device PIN verify
	EndpointCertificateE2E = BaseURL + "/LF1" // secondary device PIN verify with E2EE

	// ── Talk Service (Thrift RPC) ──────────────────────────────────────────
	// All TalkService methods POST to a single path; the method name is
	// encoded in the Thrift binary payload, not in the URL.
	// Source: zi/EnumC55943a.java

	EndpointTalk          = BaseURL + "/S4"   // enum index 2,  NORMAL
	EndpointPoll          = BaseURL + "/P4"   // enum index 0,  LONG_POLLING
	EndpointNormalPolling = BaseURL + "/NP4"  // enum index 1,  NORMAL_POLLING
	EndpointCall          = BaseURL + "/V4"   // enum index 20, CALL
	EndpointChannel       = BaseURL + "/CH4"  // enum index 15, CHANNEL
	EndpointCancelPoll    = BaseURL + "/CP4"  // enum index 16, CANCEL_LONGPOLLING

	// Compact message paths — used for lighter message payloads.
	EndpointCompactMessage      = BaseURL + "/C5"    // enum index 3,  COMPACT_MESSAGE
	EndpointCompactPlainMessage = BaseURL + "/CA5"   // enum index 4,  COMPACT_PLAIN_MESSAGE
	EndpointCompactE2EEMessage  = BaseURL + "/ECA5"  // enum index 5,  COMPACT_E2EE_MESSAGE

	// ── OpenChat (formerly LINE Square) ───────────────────────────────────
	// Source: zi/EnumC55943a.java — SQUARE is "/SQ1" (NOT the legacy "/SQS1")
	EndpointOpenChat    = BaseURL + "/SQ1"  // enum index 33, SQUARE
	EndpointOpenChatBot = BaseURL + "/BP1"  // enum index 34, SQUARE_BOT

	// ── Shop / Commerce ───────────────────────────────────────────────────
	EndpointShop             = BaseURL + "/SHOP4"                      // enum index 9,  SHOP
	EndpointShopAuth         = BaseURL + "/SHOPA"                      // enum index 10, SHOP_AUTH
	EndpointUnifiedShop      = BaseURL + "/TSHOP4"                     // enum index 13, UNIFIED_SHOP
	EndpointSticon           = BaseURL + "/SC4"                        // enum index 14, STICON
	EndpointShopRecommend    = BaseURL + "/EXT/sapi/sapir/v1p/strs"    // enum index 11, SHOP_RECOMMENDATION
	EndpointShopLFLPremium   = BaseURL + "/EXT/sapi/sapil/v1p/lps"    // enum index 12, SHOP_LFL_PREMIUM
	EndpointPoint            = BaseURL + "/POINT4"                     // enum index 35, POINT
	EndpointCoin             = BaseURL + "/COIN4"                      // enum index 36, COIN
	EndpointPay              = BaseURL + "/PY4"                        // enum index 25, PAY
	EndpointWallet           = BaseURL + "/WALLET4"                    // enum index 26, WALLET

	// ── Social / Content ──────────────────────────────────────────────────
	EndpointBuddy           = BaseURL + "/BUDDY4"                     // enum index 8,  BUDDY
	EndpointSpot            = BaseURL + "/SP4"                        // enum index 19, SPOT
	EndpointBeacon          = BaseURL + "/BEACON4"                    // enum index 31, BEACON
	EndpointPersona         = BaseURL + "/PS4"                        // enum index 32, PERSONA
	EndpointLiff            = BaseURL + "/LIFF1"                      // enum index 37, LIFF
	EndpointChatApp         = BaseURL + "/CAPP1"                      // enum index 38, CHAT_APP
	EndpointIOT             = BaseURL + "/IOT1"                       // enum index 39, IOT
	EndpointUserProvidedData = BaseURL + "/UPD4"                      // enum index 40, USER_PROVIDED_DATA
	EndpointBotExternal     = BaseURL + "/BOTE"                       // enum index 44, BOT_EXTERNAL

	// ── Search ────────────────────────────────────────────────────────────
	EndpointSearchV2             = BaseURL + "/search/v2"         // enum index 29, SEARCH_V2
	EndpointSearchV3             = BaseURL + "/search/v3"         // enum index 30, SEARCH_V3
	EndpointSearchCollectionV1   = BaseURL + "/collection/v1"     // enum index 27, SEARCH_COLLECTION_MENU_V1
	EndpointSearchCollectionV2   = BaseURL + "/collection/v2"     // enum index 28, SEARCH_COLLECTION_MENU_V2

	// ── External / Misc ───────────────────────────────────────────────────
	EndpointNotifyBackground = BaseURL + "/B"                         // enum index 7,  NOTIFY_BACKGROUND
	EndpointConnInfo         = BaseURL + "/R2"                        // enum index 22, CONN_INFO
	EndpointTyping           = BaseURL + "/TS"                        // enum index 21, TYPING (timeout 0)
	EndpointUserBehaviorLog  = BaseURL + "/L1"                        // enum index 18, USER_BEHAVIOR_LOG
	EndpointLineSpot         = BaseURL + "/ex/spot"                   // enum index 43, LINE_SPOT
	EndpointOAMembership     = BaseURL + "/EXT/oafan/api"             // enum index 45, OA_MEMBERSHIP
)
