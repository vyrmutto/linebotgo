package api

// Types extracted from LINE Android APK 15.15.1 via jadx decompilation.
// Source files referenced:
//   TalkService methods:   jp/naver/line/android/thrift/client/impl/LegacyTalkServiceClientImpl.java
//                          jp/naver/line/android/thrift/client/impl/TalkServiceClientImpl.java
//   SquareService methods: com/linecorp/square/protocol/thrift/SquareService.java
//   OpType enum:           lv1/Q6.java
//   ContentType enum:      lv1/H3.java
//   MIDType enum:          lv1/C6.java
//   ContactRelation enum:  lv1/F3.java
//   SyncReason enum:       lv1/G8.java (parameter type for many methods)
//   Profile struct:        lv1/C41072k7.java
//   Contact struct:        lv1/C41249x3.java
//   Message struct:        lv1/E6.java
//   Operation struct:      lv1/R6.java
//   LoginResult struct:    lv1/C41280z6.java
//   Square struct:         com/linecorp/square/protocol/thrift/common/Square.java
//   SquareMember struct:   com/linecorp/square/protocol/thrift/common/SquareMember.java
//   SquareChatMember:      com/linecorp/square/protocol/thrift/common/SquareChatMember.java
//   SquareMemberRole:      com/linecorp/square/protocol/thrift/common/SquareMemberRole.java
//   SquareMembershipState: com/linecorp/square/protocol/thrift/common/SquareMembershipState.java
//   SquareChatType:        com/linecorp/square/protocol/thrift/common/SquareChatType.java
//   SquareChatMembershipState: com/linecorp/square/protocol/thrift/common/SquareChatMembershipState.java
//   SendMessageRequest:    com/linecorp/square/protocol/thrift/SendMessageRequest.java

// ============================================================
// Enums — integer values match Thrift wire encoding
// ============================================================

// OpType is the operation type in the fetchOperations stream.
// Source: lv1/Q6.java (enum Q6)
// Wire type: i32 (TType.I32 = 8), field 3 of Operation struct
type OpType int32

const (
	OpTypeEndOfOperation                    OpType = 0
	OpTypeUpdateProfile                     OpType = 1
	OpTypeNotifiedUpdateProfile             OpType = 2
	OpTypeRegisterUserid                    OpType = 3
	OpTypeAddContact                        OpType = 4
	OpTypeNotifiedAddContact                OpType = 5
	OpTypeBlockContact                      OpType = 6
	OpTypeUnblockContact                    OpType = 7
	OpTypeNotifiedRecommendContact          OpType = 8
	OpTypeCreateGroup                       OpType = 9
	OpTypeUpdateGroup                       OpType = 10
	OpTypeNotifiedUpdateGroup               OpType = 11
	OpTypeInviteIntoGroup                   OpType = 12
	OpTypeNotifiedInviteIntoGroup           OpType = 13
	OpTypeLeaveGroup                        OpType = 14
	OpTypeNotifiedLeaveGroup                OpType = 15
	OpTypeAcceptGroupInvitation             OpType = 16
	OpTypeNotifiedAcceptGroupInvitation     OpType = 17
	OpTypeKickoutFromGroup                  OpType = 18
	OpTypeNotifiedKickoutFromGroup          OpType = 19
	OpTypeCreateRoom                        OpType = 20
	OpTypeInviteIntoRoom                    OpType = 21
	OpTypeNotifiedInviteIntoRoom            OpType = 22
	OpTypeLeaveRoom                         OpType = 23
	OpTypeNotifiedLeaveRoom                 OpType = 24
	OpTypeSendMessage                       OpType = 25
	OpTypeReceiveMessage                    OpType = 26
	OpTypeSendMessageReceipt                OpType = 27
	OpTypeReceiveMessageReceipt             OpType = 28
	OpTypeSendContentReceipt                OpType = 29
	OpTypeReceiveAnnouncement               OpType = 30
	OpTypeCancelInvitationGroup             OpType = 31
	OpTypeNotifiedCancelInvitationGroup     OpType = 32
	OpTypeNotifiedUnregisterUser            OpType = 33
	OpTypeRejectGroupInvitation             OpType = 34
	OpTypeNotifiedRejectGroupInvitation     OpType = 35
	OpTypeUpdateSettings                    OpType = 36
	OpTypeNotifiedRegisterUser              OpType = 37
	OpTypeInviteViaEmail                    OpType = 38
	OpTypeNotifiedRequestRecovery           OpType = 39
	OpTypeSendChatChecked                   OpType = 40
	OpTypeSendChatRemoved                   OpType = 41
	OpTypeNotifiedForceSync                 OpType = 42
	OpTypeSendContent                       OpType = 43
	OpTypeSendMessageMyhome                 OpType = 44
	OpTypeNotifiedUpdateContentPreview      OpType = 45
	OpTypeRemoveAllMessages                 OpType = 46
	OpTypeNotifiedUpdatePurchases           OpType = 47
	OpTypeDummy                             OpType = 48
	OpTypeUpdateContact                     OpType = 49
	OpTypeNotifiedReceivedCall              OpType = 50
	OpTypeCancelCall                        OpType = 51
	OpTypeNotifiedRedirect                  OpType = 52
	OpTypeNotifiedChannelSync               OpType = 53
	OpTypeFailedSendMessage                 OpType = 54
	OpTypeNotifiedReadMessage               OpType = 55
	OpTypeFailedEmailConfirmation           OpType = 56
	OpTypeNotifiedChatContent               OpType = 58
	OpTypeNotifiedPushNoticenterItem        OpType = 59
	OpTypeNotifiedJoinChat                  OpType = 60
	OpTypeNotifiedLeaveChat                 OpType = 61
	OpTypeNotifiedTyping                    OpType = 62
	OpTypeFriendRequestAccepted             OpType = 63
	OpTypeDestroyMessage                    OpType = 64
	OpTypeNotifiedDestroyMessage            OpType = 65
	OpTypeUpdatePublickeychain              OpType = 66
	OpTypeNotifiedUpdatePublickeychain      OpType = 67
	OpTypeNotifiedBlockContact              OpType = 68
	OpTypeNotifiedUnblockContact            OpType = 69
	OpTypeUpdateGrouppreference             OpType = 70
	OpTypeNotifiedPaymentEvent              OpType = 71
	OpTypeRegisterE2eePublickey             OpType = 72
	OpTypeNotifiedE2eeKeyExchangeReq        OpType = 73
	OpTypeNotifiedE2eeKeyExchangeResp       OpType = 74
	OpTypeNotifiedE2eeMessageResendReq      OpType = 75
	OpTypeNotifiedE2eeMessageResendResp     OpType = 76
	OpTypeNotifiedE2eeKeyUpdate             OpType = 77
	OpTypeNotifiedBuddyUpdateProfile        OpType = 78
	OpTypeNotifiedUpdateLineatTabs          OpType = 79
	OpTypeUpdateRoom                        OpType = 80
	OpTypeNotifiedBeaconDetected            OpType = 81
	OpTypeUpdateExtendedProfile             OpType = 82
	OpTypeAddFollow                         OpType = 83
	OpTypeNotifiedAddFollow                 OpType = 84
	OpTypeDeleteFollow                      OpType = 85
	OpTypeNotifiedDeleteFollow              OpType = 86
	OpTypeUpdateTimelineSettings            OpType = 87
	OpTypeNotifiedFriendRequest             OpType = 88
	OpTypeUpdateRingbackTone                OpType = 89
	OpTypeNotifiedPostback                  OpType = 90
	OpTypeReceiveReadWatermark              OpType = 91
	OpTypeNotifiedMessageDelivered          OpType = 92
	OpTypeNotifiedUpdateChatBar             OpType = 93
	OpTypeNotifiedChatappInstalled          OpType = 94
	OpTypeNotifiedChatappUpdated            OpType = 95
	OpTypeNotifiedChatappNewMark            OpType = 96
	OpTypeNotifiedChatappDeleted            OpType = 97
	OpTypeNotifiedChatappSync               OpType = 98
	OpTypeNotifiedUpdateMessage             OpType = 99
	OpTypeUpdateChatroombgm                 OpType = 100
	OpTypeNotifiedUpdateChatroombgm         OpType = 101
	OpTypeUpdateRingtone                    OpType = 102
	OpTypeUpdateUserSettings                OpType = 118
	OpTypeNotifiedUpdateStatusBar           OpType = 119
	OpTypeCreateChat                        OpType = 120
	OpTypeUpdateChat                        OpType = 121
	OpTypeNotifiedUpdateChat                OpType = 122
	OpTypeInviteIntoChat                    OpType = 123
	OpTypeNotifiedInviteIntoChat            OpType = 124
	OpTypeCancelChatInvitation              OpType = 125
	OpTypeNotifiedCancelChatInvitation      OpType = 126
	OpTypeDeleteSelfFromChat                OpType = 127
	OpTypeNotifiedDeleteSelfFromChat        OpType = 128
	OpTypeAcceptChatInvitation              OpType = 129
	OpTypeNotifiedAcceptChatInvitation      OpType = 130
	OpTypeRejectChatInvitation              OpType = 131
	OpTypeDeleteOtherFromChat               OpType = 132
	OpTypeNotifiedDeleteOtherFromChat       OpType = 133
	OpTypeNotifiedContactCalendarEvent      OpType = 134
	OpTypeNotifiedContactCalendarEventAll   OpType = 135
	OpTypeUpdateThingsOperations            OpType = 136
	OpTypeSendChatHidden                    OpType = 137
	OpTypeChatMetaSyncAll                   OpType = 138
	OpTypeSendReaction                      OpType = 139
	OpTypeNotifiedSendReaction              OpType = 140
	OpTypeNotifiedUpdateProfileContent      OpType = 141
	OpTypeFailedDeliveryMessage             OpType = 142
	OpTypeSendEncryptedE2eeKeyRequested     OpType = 143
	OpTypeChannelPaakAuthenticationRequested OpType = 144
	OpTypeUpdatePinState                    OpType = 145
	OpTypeNotifiedPremiumbackupStateChanged OpType = 146
	OpTypeCreateMultiProfile                OpType = 147
	OpTypeMultiProfileStatusChanged         OpType = 148
	OpTypeDeleteMultiProfile                OpType = 149
	OpTypeUpdateProfileMapping              OpType = 150
	OpTypeDeleteProfileMapping              OpType = 151
	OpTypeNotifiedDestroyNoticenterPush     OpType = 152
	OpTypeForceKeyBackupHeaderValidation    OpType = 153
	OpTypeNotifiedGcsReaction               OpType = 154
	OpTypeUpdateMessageRequestBox           OpType = 155
	OpTypeNotifiedUpdateMessageRequestBox   OpType = 156
	OpTypeNotifiedGcsRefreshContent         OpType = 157
)

// ContentType represents the content type of a message.
// Source: lv1/H3.java (enum H3)
// Wire type: i32 (TType.I32 = 8), field 15 of Message struct (E6.java)
type ContentType int32

const (
	ContentTypeNone             ContentType = 0  // plain text
	ContentTypeImage            ContentType = 1
	ContentTypeVideo            ContentType = 2
	ContentTypeAudio            ContentType = 3
	ContentTypeHTML             ContentType = 4
	ContentTypePDF              ContentType = 5
	ContentTypeCall             ContentType = 6
	ContentTypeSticker          ContentType = 7
	ContentTypePresence         ContentType = 8
	ContentTypeGift             ContentType = 9
	ContentTypeGroupboard       ContentType = 10
	ContentTypeApplink          ContentType = 11
	ContentTypeLink             ContentType = 12
	ContentTypeContact          ContentType = 13
	ContentTypeFile             ContentType = 14
	ContentTypeLocation         ContentType = 15
	ContentTypePostnotification ContentType = 16
	ContentTypeRich             ContentType = 17
	ContentTypeChatevent        ContentType = 18
	ContentTypeMusic            ContentType = 19
	ContentTypePayment          ContentType = 20
	ContentTypeExtimage         ContentType = 21
	ContentTypeFlex             ContentType = 22
)

// MIDType represents the type of a LINE MID (message/entity identifier).
// Source: lv1/C6.java (enum C6)
// Wire type: i32 (TType.I32 = 8), field 3 of Message struct (toType, field 3)
type MIDType int32

const (
	MIDTypeUser         MIDType = 0
	MIDTypeRoom         MIDType = 1
	MIDTypeGroup        MIDType = 2
	MIDTypeSquare       MIDType = 3
	MIDTypeSquareChat   MIDType = 4
	MIDTypeSquareMember MIDType = 5
	MIDTypeBot          MIDType = 6
	MIDTypeSquareThread MIDType = 7
)

// ContactRelation represents the relationship with a contact.
// Source: lv1/F3.java (enum F3) — used in Contact struct field 21 "relation"
type ContactRelation int32

const (
	ContactRelationUnspecified      ContactRelation = 0
	ContactRelationFriend           ContactRelation = 1
	ContactRelationFriendBlocked    ContactRelation = 2
	ContactRelationRecommend        ContactRelation = 3
	ContactRelationRecommendBlocked ContactRelation = 4
	ContactRelationDeleted          ContactRelation = 5
	ContactRelationDeletedBlocked   ContactRelation = 6
)

// SyncReason is passed to many TalkService methods to describe why sync is triggered.
// Source: lv1/G8.java (enum G8)
type SyncReason int32

const (
	SyncReasonUnspecified           SyncReason = 0
	SyncReasonUnknown               SyncReason = 1
	SyncReasonInitialization        SyncReason = 2
	SyncReasonOperation             SyncReason = 3
	SyncReasonFullSync              SyncReason = 4
	SyncReasonAutoRepair            SyncReason = 5
	SyncReasonManualRepair          SyncReason = 6
	SyncReasonInternal              SyncReason = 7
	SyncReasonUserInitiated         SyncReason = 8
	SyncReasonPremiumBackupRestore  SyncReason = 9
	SyncReasonPushToLoad            SyncReason = 10
)

// SquareMemberRole is the role of a member in a Square (OpenChat).
// Source: com/linecorp/square/protocol/thrift/common/SquareMemberRole.java
type SquareMemberRole int32

const (
	SquareMemberRoleAdmin   SquareMemberRole = 1
	SquareMemberRoleCoAdmin SquareMemberRole = 2
	SquareMemberRoleMember  SquareMemberRole = 10
)

// SquareMembershipState is the overall membership state in a Square.
// Source: com/linecorp/square/protocol/thrift/common/SquareMembershipState.java
type SquareMembershipState int32

const (
	SquareMembershipStateJoinRequested          SquareMembershipState = 1
	SquareMembershipStateJoined                 SquareMembershipState = 2
	SquareMembershipStateRejected               SquareMembershipState = 3
	SquareMembershipStateLeft                   SquareMembershipState = 4
	SquareMembershipStateKickOut                SquareMembershipState = 5
	SquareMembershipStateBanned                 SquareMembershipState = 6
	SquareMembershipStateDeleted                SquareMembershipState = 7
	SquareMembershipStateJoinRequestWithdrew    SquareMembershipState = 8
	SquareMembershipStateJoinReserved           SquareMembershipState = 9
	SquareMembershipStateJoinReservationExpired SquareMembershipState = 10
)

// SquareChatMembershipState is the chat-level membership state.
// Source: com/linecorp/square/protocol/thrift/common/SquareChatMembershipState.java
type SquareChatMembershipState int32

const (
	SquareChatMembershipStateJoined                 SquareChatMembershipState = 1
	SquareChatMembershipStateLeft                   SquareChatMembershipState = 2
	SquareChatMembershipStateJoinReserved           SquareChatMembershipState = 3
	SquareChatMembershipStateJoinReservationExpired SquareChatMembershipState = 4
)

// SquareChatType represents the type of an OpenChat room.
// Source: com/linecorp/square/protocol/thrift/common/SquareChatType.java
type SquareChatType int32

const (
	SquareChatTypeOpen          SquareChatType = 1
	SquareChatTypeSecret        SquareChatType = 2
	SquareChatTypeOneOnOne      SquareChatType = 3
	SquareChatTypeSquareDefault SquareChatType = 4
)

// SquareType indicates if a Square is open (1) or closed (0).
// Source: com/linecorp/square/protocol/thrift/common/SquareType.java
type SquareType int32

const (
	SquareTypeClosed SquareType = 0
	SquareTypeOpen   SquareType = 1
)

// SquareMemberRelationState represents the relation state (e.g., blocked).
// Source: com/linecorp/square/protocol/thrift/common/SquareMemberRelationState.java
type SquareMemberRelationState int32

const (
	SquareMemberRelationStateNone    SquareMemberRelationState = 1
	SquareMemberRelationStateBlocked SquareMemberRelationState = 2
)

// ============================================================
// Core Thrift struct types
// ============================================================

// Message is a LINE message transferred via TalkService.
// Source: lv1/E6.java — Thrift struct with field IDs
//
// Thrift field map (field_id → name → TType):
//   1  → from          → STRING  (11)
//   2  → to            → STRING  (11)
//   3  → toType        → I32     (8)  — MIDType enum
//   4  → id            → STRING  (11)
//   5  → createdTime   → I64     (10)
//   6  → deliveredTime → I64     (10)
//  10  → text          → STRING  (11)
//  11  → location      → STRUCT  (12)
//  14  → hasContent    → BOOL    (2)
//  15  → contentType   → I32     (8)  — ContentType enum
//  17  → contentPreview → STRING (11)
//  18  → contentMetadata → MAP   (13)
//  19  → sessionId     → BYTE    (3)
//  20  → chunks        → LIST    (15)
//  21  → relatedMessageId → STRING (11)
//  22  → messageRelationType → I32 (8)
//  23  → readCount     → I32     (8)
//  24  → relatedMessageServiceCode → I32 (8)
//  25  → appExtensionType → I32  (8)
//  27  → reactions     → LIST    (15)
type Message struct {
	// Field 1 — sender MID
	From string `json:"from"`
	// Field 2 — recipient MID (user, room, or group MID)
	To string `json:"to"`
	// Field 3 — type of recipient MID; see MIDType constants
	ToType MIDType `json:"toType"`
	// Field 4 — message ID (string representation of i64)
	ID string `json:"id"`
	// Field 5 — milliseconds since epoch
	CreatedTime int64 `json:"createdTime"`
	// Field 6 — delivered time (ms since epoch)
	DeliveredTime int64 `json:"deliveredTime,omitempty"`
	// Field 10 — plain text body (non-empty for NONE content type)
	Text string `json:"text,omitempty"`
	// Field 14 — true when content (image/video/etc.) is attached
	HasContent bool `json:"hasContent,omitempty"`
	// Field 15 — content type; see ContentType constants
	ContentType ContentType `json:"contentType"`
	// Field 17 — preview URL for media content
	ContentPreview string `json:"contentPreview,omitempty"`
	// Field 18 — arbitrary metadata map (key-value string pairs)
	ContentMetadata map[string]string `json:"contentMetadata,omitempty"`
	// Field 21 — for reply messages: the ID of the referenced message
	RelatedMessageID string `json:"relatedMessageId,omitempty"`
}

// Operation is a TalkService operation fetched via fetchOperations/fetchOps.
// Source: lv1/R6.java — Thrift struct
//
// Thrift field map:
//   1  → revision    → I64    (10)
//   2  → createdTime → I64    (10)
//   3  → type        → I32    (8)  — OpType enum
//   4  → reqSeq      → I32    (8)
//   5  → checksum    → STRING (11)
//   7  → status      → I32    (8)
//  10  → param1      → STRING (11)
//  11  → param2      → STRING (11)
//  12  → param3      → STRING (11)
//  20  → message     → STRUCT (12) — Message struct
type Operation struct {
	// Field 1 — monotonically increasing revision; pass as cursor in next poll
	Revision int64 `json:"revision"`
	// Field 2 — milliseconds since epoch
	CreatedTime int64 `json:"createdTime"`
	// Field 3 — operation type; see OpType constants
	Type OpType `json:"type"`
	// Field 4 — request sequence number (for client-initiated ops)
	ReqSeq int32 `json:"reqSeq,omitempty"`
	// Field 5 — checksum (rarely used)
	Checksum string `json:"checksum,omitempty"`
	// Field 7 — operation status
	Status int32 `json:"status,omitempty"`
	// Field 10 — first operation parameter (interpretation depends on OpType)
	Param1 string `json:"param1,omitempty"`
	// Field 11 — second operation parameter
	Param2 string `json:"param2,omitempty"`
	// Field 12 — third operation parameter
	Param3 string `json:"param3,omitempty"`
	// Field 20 — embedded message (for SEND_MESSAGE / RECEIVE_MESSAGE ops)
	Message *Message `json:"message,omitempty"`
}

// Profile is the full user profile returned by TalkService.getProfile.
// Source: lv1/C41072k7.java — Thrift struct
//
// Thrift field map (selected fields):
//   1  → mid           → STRING (11)
//   3  → userid        → STRING (11)
//  10  → phone         → STRING (11)
//  11  → email         → STRING (11)
//  12  → regionCode    → STRING (11)
//  20  → displayName   → STRING (11)
//  21  → phoneticName  → STRING (11)
//  22  → pictureStatus → STRING (11)
//  23  → thumbnailUrl  → STRING (11)
//  24  → statusMessage → STRING (11)
//  31  → allowSearchByUserid → BOOL (2)
//  32  → allowSearchByEmail  → BOOL (2)
//  33  → picturePath   → STRING (11)
//  34  → musicProfile  → STRING (11)
//  35  → videoProfile  → STRING (11)
//  36  → statusMessageContentMetadata → MAP (13)
//  37  → avatarProfile → STRUCT (12)
//  38  → nftProfile    → BOOL (2)
//  39  → pictureSource → I32 (8)
//  40  → profileId     → STRING (11)
//  41  → profileType   → I32 (8)
//  42  → createdTimeMillis → I64 (10)
type Profile struct {
	// Field 1 — user MID (starts with 'u')
	MID string `json:"mid"`
	// Field 3 — user-chosen @username
	UserID string `json:"userid,omitempty"`
	// Field 10 — phone number (E.164 format)
	Phone string `json:"phone,omitempty"`
	// Field 11 — email address
	Email string `json:"email,omitempty"`
	// Field 12 — two-letter region code (e.g., "JP")
	RegionCode string `json:"regionCode,omitempty"`
	// Field 20 — display name shown in chats
	DisplayName string `json:"displayName"`
	// Field 21 — phonetic/reading name
	PhoneticName string `json:"phoneticName,omitempty"`
	// Field 22 — picture status token
	PictureStatus string `json:"pictureStatus,omitempty"`
	// Field 23 — thumbnail image URL
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	// Field 24 — status message text
	StatusMessage string `json:"statusMessage,omitempty"`
	// Field 31 — whether profile is searchable by user ID
	AllowSearchByUserid bool `json:"allowSearchByUserid,omitempty"`
	// Field 32 — whether profile is searchable by email
	AllowSearchByEmail bool `json:"allowSearchByEmail,omitempty"`
	// Field 33 — OBS (object storage) path for profile picture
	PicturePath string `json:"picturePath,omitempty"`
	// Field 34 — music profile data
	MusicProfile string `json:"musicProfile,omitempty"`
	// Field 40 — numeric profile ID string
	ProfileID string `json:"profileId,omitempty"`
	// Field 42 — account creation timestamp (ms since epoch)
	CreatedTimeMillis int64 `json:"createdTimeMillis,omitempty"`
}

// Contact represents a contact in the user's contact list.
// Source: lv1/C41249x3.java — Thrift struct
//
// Thrift field map (selected fields):
//   1  → mid                 → STRING (11)
//   2  → createdTime         → I64    (10)
//  10  → type                → I32    (8)  — MIDType
//  11  → status              → I32    (8)
//  21  → relation            → I32    (8)  — ContactRelation enum
//  22  → displayName         → STRING (11)
//  23  → phoneticName        → STRING (11)
//  24  → pictureStatus       → STRING (11)
//  25  → thumbnailUrl        → STRING (11)
//  26  → statusMessage       → STRING (11)
//  27  → displayNameOverridden → STRING (11)
//  28  → favoriteTime        → I64    (10)
//  31  → capableVoiceCall    → BOOL   (2)
//  32  → capableVideoCall    → BOOL   (2)
//  33  → capableMyhome       → BOOL   (2)
//  34  → capableBuddy        → BOOL   (2)
//  35  → attributes          → I32    (8)
//  36  → settings            → I64    (10)
//  37  → picturePath         → STRING (11)
//  38  → recommendParams     → STRING (11)
//  39  → friendRequestStatus → I32    (8)
//  40  → musicProfile        → STRING (11)
//  42  → videoProfile        → STRING (11)
//  43  → statusMessageContentMetadata → MAP (13)
//  44  → avatarProfile       → STRUCT (12)
//  45  → friendRingtone      → STRING (11)
//  46  → friendRingbackTone  → STRING (11)
//  47  → nftProfile          → BOOL   (2)
//  48  → pictureSource       → I32    (8)
//  49  → profileId           → STRING (11)
type Contact struct {
	// Field 1 — contact MID
	MID string `json:"mid"`
	// Field 2 — time the contact was added (ms since epoch)
	CreatedTime int64 `json:"createdTime,omitempty"`
	// Field 10 — MID type (normally MIDTypeUser = 0)
	Type MIDType `json:"type,omitempty"`
	// Field 11 — contact status integer
	Status int32 `json:"status,omitempty"`
	// Field 21 — relationship with this contact; see ContactRelation constants
	Relation ContactRelation `json:"relation,omitempty"`
	// Field 22 — display name
	DisplayName string `json:"displayName"`
	// Field 23 — phonetic/reading name
	PhoneticName string `json:"phoneticName,omitempty"`
	// Field 24 — picture status token
	PictureStatus string `json:"pictureStatus,omitempty"`
	// Field 25 — thumbnail image URL
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	// Field 26 — status message
	StatusMessage string `json:"statusMessage,omitempty"`
	// Field 27 — display name override (nickname set by this user)
	DisplayNameOverridden string `json:"displayNameOverridden,omitempty"`
	// Field 28 — timestamp when marked as favourite (ms since epoch)
	FavoriteTime int64 `json:"favoriteTime,omitempty"`
	// Field 37 — OBS path for profile picture
	PicturePath string `json:"picturePath,omitempty"`
	// Field 39 — friend request status integer
	FriendRequestStatus int32 `json:"friendRequestStatus,omitempty"`
	// Field 49 — numeric profile ID string
	ProfileID string `json:"profileId,omitempty"`
}

// LoginResult is the result of a TalkService login call.
// Source: lv1/C41280z6.java — Thrift struct
//
// Thrift field map:
//   1  → authToken            → STRING (11)
//   2  → certificate          → STRING (11)
//   3  → verifier             → STRING (11)
//   4  → pinCode              → STRING (11)
//   5  → type                 → I32    (8)  — LoginResultType
//   6  → lastPrimaryBindTime  → I64    (10)
//   7  → displayMessage       → STRING (11)
//   8  → sessionForSMSConfirm → STRUCT (12)
//   9  → tokenV3IssueResult   → STRUCT (12)
//  10  → mid                  → STRING (11)
type LoginResult struct {
	// Field 1 — session auth token for subsequent calls
	AuthToken string `json:"authToken"`
	// Field 2 — device certificate (base64)
	Certificate string `json:"certificate,omitempty"`
	// Field 3 — verifier for two-step auth
	Verifier string `json:"verifier,omitempty"`
	// Field 4 — PIN code for device confirmation
	PinCode string `json:"pinCode,omitempty"`
	// Field 5 — login result type (0=success, 1=verify device, 2=verify SMS, etc.)
	Type int32 `json:"type"`
	// Field 7 — human-readable message (for errors / prompts)
	DisplayMessage string `json:"displayMessage,omitempty"`
	// Field 10 — user MID (available after successful login)
	MID string `json:"mid,omitempty"`
}

// ============================================================
// OpenChat (Square) struct types
// ============================================================

// Square represents a LINE OpenChat community (formerly LINE Square).
// Source: com/linecorp/square/protocol/thrift/common/Square.java
//
// Thrift field map (selected fields):
//   1  → mid                        → STRING (11)
//   2  → name                       → STRING (11)
//   3  → welcomeMessage             → STRING (11)
//   4  → profileImageObsHash        → STRING (11)
//   5  → desc                       → STRING (11)
//   6  → searchable                 → BOOL   (2)
//   7  → type                       → I32    (8)  — SquareType
//   8  → categoryId                 → I32    (8)
//   9  → invitationURL              → STRING (11)
//  10  → revision                   → I64    (10)
//  11  → ableToUseInvitationTicket  → BOOL   (2)
//  12  → state                      → I32    (8)
//  17  → createdAt                  → I64    (10)
//  18  → paidSquare                 → BOOL   (2)
//  21  → expireAt                   → I64    (10)
type Square struct {
	// Field 1 — Square MID (starts with 's')
	MID string `json:"mid"`
	// Field 2 — community name
	Name string `json:"name"`
	// Field 3 — welcome message shown to new members
	WelcomeMessage string `json:"welcomeMessage,omitempty"`
	// Field 4 — OBS hash for the profile image
	ProfileImageObsHash string `json:"profileImageObsHash,omitempty"`
	// Field 5 — community description
	Description string `json:"desc,omitempty"`
	// Field 6 — whether the Square appears in search results
	Searchable bool `json:"searchable,omitempty"`
	// Field 7 — Square type (0=closed, 1=open)
	Type SquareType `json:"type,omitempty"`
	// Field 8 — category ID integer
	CategoryID int32 `json:"categoryId,omitempty"`
	// Field 9 — invitation URL
	InvitationURL string `json:"invitationURL,omitempty"`
	// Field 10 — monotonic revision counter
	Revision int64 `json:"revision,omitempty"`
	// Field 17 — creation timestamp (ms since epoch)
	CreatedAt int64 `json:"createdAt,omitempty"`
}

// SquareMember represents a member of a Square community.
// Source: com/linecorp/square/protocol/thrift/common/SquareMember.java
//
// Thrift field map:
//   1  → squareMemberMid        → STRING (11)
//   2  → squareMid              → STRING (11)
//   3  → displayName            → STRING (11)
//   4  → profileImageObsHash    → STRING (11)
//   5  → ableToReceiveMessage   → BOOL   (2)
//   7  → membershipState        → I32    (8)  — SquareMembershipState
//   8  → role                   → I32    (8)  — SquareMemberRole
//   9  → revision               → I64    (10)
//  10  → preference             → STRUCT (12)
//  11  → joinMessage            → STRING (11)
//  12  → createdAt              → I64    (10)
//  13  → selfIntroduction       → STRING (11)
//  14  → socialMediaAccountUrls → LIST   (15)
type SquareMember struct {
	// Field 1 — member MID within the Square (same as user MID)
	SquareMemberMID string `json:"squareMemberMid"`
	// Field 2 — parent Square MID
	SquareMID string `json:"squareMid"`
	// Field 3 — display name within the Square
	DisplayName string `json:"displayName"`
	// Field 4 — OBS hash for the member profile image
	ProfileImageObsHash string `json:"profileImageObsHash,omitempty"`
	// Field 5 — whether this member can receive messages
	AbleToReceiveMessage bool `json:"ableToReceiveMessage,omitempty"`
	// Field 7 — membership state; see SquareMembershipState constants
	MembershipState SquareMembershipState `json:"membershipState,omitempty"`
	// Field 8 — role; see SquareMemberRole constants
	Role SquareMemberRole `json:"role,omitempty"`
	// Field 9 — monotonic revision counter
	Revision int64 `json:"revision,omitempty"`
	// Field 11 — join message
	JoinMessage string `json:"joinMessage,omitempty"`
	// Field 12 — when the member joined (ms since epoch)
	CreatedAt int64 `json:"createdAt,omitempty"`
	// Field 13 — self-introduction text
	SelfIntroduction string `json:"selfIntroduction,omitempty"`
}

// SquareChatMember represents a member's participation in a specific chat within a Square.
// Source: com/linecorp/square/protocol/thrift/common/SquareChatMember.java
//
// Thrift field map:
//   1  → squareMemberMid          → STRING (11)
//   2  → squareChatMid            → STRING (11)
//   3  → revision                 → I64    (10)
//   4  → membershipState          → I32    (8)  — SquareChatMembershipState
//   5  → notificationForMessage   → BOOL   (2)
//   6  → notificationForNewMember → BOOL   (2)
type SquareChatMember struct {
	// Field 1 — member MID
	SquareMemberMID string `json:"squareMemberMid"`
	// Field 2 — chat MID
	SquareChatMID string `json:"squareChatMid"`
	// Field 3 — monotonic revision
	Revision int64 `json:"revision,omitempty"`
	// Field 4 — chat membership state; see SquareChatMembershipState constants
	MembershipState SquareChatMembershipState `json:"membershipState,omitempty"`
	// Field 5 — whether message notifications are enabled
	NotificationForMessage bool `json:"notificationForMessage,omitempty"`
	// Field 6 — whether new-member notifications are enabled
	NotificationForNewMember bool `json:"notificationForNewMember,omitempty"`
}

// OpenChatInfo represents summary information about a LINE OpenChat room.
// Used as the response type for search and list operations.
type OpenChatInfo struct {
	// Square/chat MID
	ID string `json:"id"`
	// Chat name
	Name string `json:"name"`
	// Chat description
	Description string `json:"description,omitempty"`
	// Approximate member count
	MemberCount int `json:"memberCount,omitempty"`
	// Type of chat; see SquareChatType constants
	ChatType SquareChatType `json:"chatType,omitempty"`
	// Whether the chat is searchable
	Searchable bool `json:"searchable,omitempty"`
}

// OpenChatMember represents a member of an OpenChat room for API responses.
// Simplified view of SquareMember for external usage.
type OpenChatMember struct {
	// User/member MID
	MID string `json:"mid"`
	// Display name within the chat
	DisplayName string `json:"displayName"`
	// Role string derived from SquareMemberRole: "ADMIN", "CO_ADMIN", or "MEMBER"
	Role string `json:"role"`
	// Profile image hash (use with OBS URL to construct image URL)
	ProfileImageObsHash string `json:"profileImageObsHash,omitempty"`
	// Membership state string
	MembershipState SquareMembershipState `json:"membershipState,omitempty"`
}

// ============================================================
// Service request/response types
// ============================================================

// SendMessageRequest is the SquareService.sendMessage request payload.
// Source: com/linecorp/square/protocol/thrift/SendMessageRequest.java
//
// Thrift field map:
//   1  → reqSeq        → I32    (8)
//   2  → squareChatMid → STRING (11)
//   3  → squareMessage → STRUCT (12) — SquareMessage struct
type SendMessageRequest struct {
	// Field 1 — client-side sequence number for deduplication
	ReqSeq int32 `json:"reqSeq"`
	// Field 2 — target Square chat MID
	SquareChatMID string `json:"squareChatMid"`
	// Field 3 — the message to send (embedded Message with text/content)
	Message *Message `json:"squareMessage"`
}

// SyncRequest is the Thrift struct sent as the argument to fetchOperations on /P4.
// Source: lv1/H8.java
//
// Thrift field map (wire types: 10=I64, 8=I32, 13=MAP):
//
//	1  → lastRevision          → I64   (10) — required
//	2  → count                 → I32   (8)  — required, max ops to return
//	3  → lastGlobalRevision    → I64   (10) — optional
//	4  → lastIndividualRevision → I64  (10) — optional
//	5  → fullSyncRequestReason → I32   (8)  — optional (enum SyncReason)
//	6  → lastPartialFullSyncs  → MAP<I32,I64> (13) — optional
type SyncRequest struct {
	// Field 1 — last revision the client has processed
	LastRevision int64 `json:"lastRevision"`
	// Field 2 — max number of operations to return (typically 50)
	Count int32 `json:"count"`
	// Field 3 — optional; last global revision
	LastGlobalRevision int64 `json:"lastGlobalRevision,omitempty"`
	// Field 4 — optional; last individual revision
	LastIndividualRevision int64 `json:"lastIndividualRevision,omitempty"`
	// Field 5 — optional; reason for full sync (enum lv1/U4)
	FullSyncRequestReason int32 `json:"fullSyncRequestReason,omitempty"`
}

// FetchOperationsRequest contains parameters for TalkService.fetchOperations.
// The Thrift method sends these as positional args in the binary protocol:
//   arg 2 (i32): local revision
//   arg 3 (i32): count (max operations to return, typically 50)
type FetchOperationsRequest struct {
	// Last revision seen by the client; operations after this are returned
	Revision int64 `json:"revision"`
	// Maximum number of operations to return (LINE server accepts up to 50)
	Count int32 `json:"count"`
}

// FetchOperationsResponse is the list of operations returned by fetchOperations.
type FetchOperationsResponse struct {
	// Ordered list of operations; update local revision after processing each
	Operations []*Operation `json:"operations"`
}
