// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

//go:generate go run gen.go

package steamworks

import (
	"iter"
)

type AppId_t uint32
type CGameID uint64
type CSteamID uint64
type DepotId_t uint32
type InputActionSetHandle_t uint64
type InputAnalogActionHandle_t uint64
type InputDigitalActionHandle_t uint64
type InputHandle_t uint64
type SteamAPICall_t uint64
type HTTPRequestHandle uint32
type HTTPCookieContainerHandle uint32
type PublishedFileId_t uint64
type SteamItemInstanceID_t uint64
type SteamItemDef_t int32
type SteamInventoryResult_t int32
type HSteamUser int32
type HAuthTicket uint32
type HServerListRequest uintptr
type HServerQuery int32

type HSteamNetConnection uint32
type HSteamListenSocket uint32
type HSteamNetPollGroup uint32
type SteamNetworkingMicroseconds int64

type ESteamAPIInitResult int32

type EVoiceResult int32

type EBeginAuthSessionResult int32

type EUserHasLicenseForAppResult int32

type EDurationControlProgress int32

type EDurationControlNotification int32

type EDurationControlOnlineState int32

const (
	ESteamAPIInitResult_OK              ESteamAPIInitResult = 0
	ESteamAPIInitResult_FailedGeneric   ESteamAPIInitResult = 1
	ESteamAPIInitResult_NoSteamClient   ESteamAPIInitResult = 2
	ESteamAPIInitResult_VersionMismatch ESteamAPIInitResult = 3
)

type ESteamInputType int32

const (
	ESteamInputType_Unknown              ESteamInputType = 0
	ESteamInputType_SteamController      ESteamInputType = 1
	ESteamInputType_XBox360Controller    ESteamInputType = 2
	ESteamInputType_XBoxOneController    ESteamInputType = 3
	ESteamInputType_GenericXInput        ESteamInputType = 4
	ESteamInputType_PS4Controller        ESteamInputType = 5
	ESteamInputType_AppleMFiController   ESteamInputType = 6 // Unused
	ESteamInputType_AndroidController    ESteamInputType = 7 // Unused
	ESteamInputType_SwitchJoyConPair     ESteamInputType = 8 // Unused
	ESteamInputType_SwitchJoyConSingle   ESteamInputType = 9 // Unused
	ESteamInputType_SwitchProController  ESteamInputType = 10
	ESteamInputType_MobileTouch          ESteamInputType = 11
	ESteamInputType_PS3Controller        ESteamInputType = 12
	ESteamInputType_PS5Controller        ESteamInputType = 13
	ESteamInputType_SteamDeckController  ESteamInputType = 14
	ESteamInputType_Count                ESteamInputType = 15
	ESteamInputType_MaximumPossibleValue ESteamInputType = 255
)

const (
	_STEAM_INPUT_MAX_COUNT = 16
)

type ESteamInputLEDFlag uint32

const (
	ESteamInputLEDFlag_SetColor           ESteamInputLEDFlag = 0x01
	ESteamInputLEDFlag_RestoreUserDefault ESteamInputLEDFlag = 0x02
)

type ESteamControllerPad int32

const (
	ESteamControllerPad_Left  ESteamControllerPad = 0
	ESteamControllerPad_Right ESteamControllerPad = 1
)

type EInputSourceMode int32

const (
	EInputSourceMode_None           EInputSourceMode = 0
	EInputSourceMode_Dpad           EInputSourceMode = 1
	EInputSourceMode_Buttons        EInputSourceMode = 2
	EInputSourceMode_FourButtons    EInputSourceMode = 3
	EInputSourceMode_AbsoluteMouse  EInputSourceMode = 4
	EInputSourceMode_RelativeMouse  EInputSourceMode = 5
	EInputSourceMode_JoystickMove   EInputSourceMode = 6
	EInputSourceMode_JoystickCamera EInputSourceMode = 7
	EInputSourceMode_ScrollWheel    EInputSourceMode = 8
	EInputSourceMode_Trigger        EInputSourceMode = 9
	EInputSourceMode_TouchMenu      EInputSourceMode = 10
	EInputSourceMode_MouseJoystick  EInputSourceMode = 11
	EInputSourceMode_MouseRegion    EInputSourceMode = 12
	EInputSourceMode_RadialMenu     EInputSourceMode = 13
	EInputSourceMode_SingleButton   EInputSourceMode = 14
	EInputSourceMode_Switches       EInputSourceMode = 15
)

type EInputActionOrigin int32

const (
	EInputActionOrigin_None                                EInputActionOrigin = 0
	EInputActionOrigin_SteamController_A                   EInputActionOrigin = 1
	EInputActionOrigin_SteamController_B                   EInputActionOrigin = 2
	EInputActionOrigin_SteamController_X                   EInputActionOrigin = 3
	EInputActionOrigin_SteamController_Y                   EInputActionOrigin = 4
	EInputActionOrigin_SteamController_LeftBumper          EInputActionOrigin = 5
	EInputActionOrigin_SteamController_RightBumper         EInputActionOrigin = 6
	EInputActionOrigin_SteamController_LeftGrip            EInputActionOrigin = 7
	EInputActionOrigin_SteamController_RightGrip           EInputActionOrigin = 8
	EInputActionOrigin_SteamController_Start               EInputActionOrigin = 9
	EInputActionOrigin_SteamController_Back                EInputActionOrigin = 10
	EInputActionOrigin_SteamController_LeftPad_Touch       EInputActionOrigin = 11
	EInputActionOrigin_SteamController_LeftPad_Swipe       EInputActionOrigin = 12
	EInputActionOrigin_SteamController_LeftPad_Click       EInputActionOrigin = 13
	EInputActionOrigin_SteamController_LeftPad_DPadNorth   EInputActionOrigin = 14
	EInputActionOrigin_SteamController_LeftPad_DPadSouth   EInputActionOrigin = 15
	EInputActionOrigin_SteamController_LeftPad_DPadWest    EInputActionOrigin = 16
	EInputActionOrigin_SteamController_LeftPad_DPadEast    EInputActionOrigin = 17
	EInputActionOrigin_SteamController_RightPad_Touch      EInputActionOrigin = 18
	EInputActionOrigin_SteamController_RightPad_Swipe      EInputActionOrigin = 19
	EInputActionOrigin_SteamController_RightPad_Click      EInputActionOrigin = 20
	EInputActionOrigin_SteamController_RightPad_DPadNorth  EInputActionOrigin = 21
	EInputActionOrigin_SteamController_RightPad_DPadSouth  EInputActionOrigin = 22
	EInputActionOrigin_SteamController_RightPad_DPadWest   EInputActionOrigin = 23
	EInputActionOrigin_SteamController_RightPad_DPadEast   EInputActionOrigin = 24
	EInputActionOrigin_SteamController_LeftTrigger_Pull    EInputActionOrigin = 25
	EInputActionOrigin_SteamController_LeftTrigger_Click   EInputActionOrigin = 26
	EInputActionOrigin_SteamController_RightTrigger_Pull   EInputActionOrigin = 27
	EInputActionOrigin_SteamController_RightTrigger_Click  EInputActionOrigin = 28
	EInputActionOrigin_SteamController_LeftStick_Move      EInputActionOrigin = 29
	EInputActionOrigin_SteamController_LeftStick_Click     EInputActionOrigin = 30
	EInputActionOrigin_SteamController_LeftStick_DPadNorth EInputActionOrigin = 31
	EInputActionOrigin_SteamController_LeftStick_DPadSouth EInputActionOrigin = 32
	EInputActionOrigin_SteamController_LeftStick_DPadWest  EInputActionOrigin = 33
	EInputActionOrigin_SteamController_LeftStick_DPadEast  EInputActionOrigin = 34
	EInputActionOrigin_SteamController_Gyro_Move           EInputActionOrigin = 35
	EInputActionOrigin_SteamController_Gyro_Pitch          EInputActionOrigin = 36
	EInputActionOrigin_SteamController_Gyro_Yaw            EInputActionOrigin = 37
	EInputActionOrigin_SteamController_Gyro_Roll           EInputActionOrigin = 38
)

type InputDigitalActionData struct {
	State  bool
	Active bool
}

type InputAnalogActionData struct {
	Mode   EInputSourceMode
	X      float32
	Y      float32
	Active bool
}

type InputMotionData struct {
	RotQuatX  float32
	RotQuatY  float32
	RotQuatZ  float32
	RotQuatW  float32
	PosAccelX float32
	PosAccelY float32
	PosAccelZ float32
	RotVelX   float32
	RotVelY   float32
	RotVelZ   float32
}

type EFloatingGamepadTextInputMode int32

const (
	EFloatingGamepadTextInputMode_ModeSingleLine    EFloatingGamepadTextInputMode = 0
	EFloatingGamepadTextInputMode_ModeMultipleLines EFloatingGamepadTextInputMode = 1
	EFloatingGamepadTextInputMode_ModeEmail         EFloatingGamepadTextInputMode = 2
	EFloatingGamepadTextInputMode_ModeNumeric       EFloatingGamepadTextInputMode = 3
)

type EOverlayToStoreFlag int32

const (
	EOverlayToStoreFlag_None             EOverlayToStoreFlag = 0
	EOverlayToStoreFlag_AddToCart        EOverlayToStoreFlag = 1
	EOverlayToStoreFlag_AddToCartAndShow EOverlayToStoreFlag = 2
)

type EActivateGameOverlayToWebPageMode int32

const (
	EActivateGameOverlayToWebPageMode_Default EActivateGameOverlayToWebPageMode = 0
	EActivateGameOverlayToWebPageMode_Modal   EActivateGameOverlayToWebPageMode = 1
)

type ELobbyType int32

const (
	ELobbyType_Private     ELobbyType = 0
	ELobbyType_FriendsOnly ELobbyType = 1
	ELobbyType_Public      ELobbyType = 2
	ELobbyType_Invisible   ELobbyType = 3
)

type ELobbyComparison int32

const (
	ELobbyComparisonEqualToOrLessThan    ELobbyComparison = -2
	ELobbyComparisonLessThan             ELobbyComparison = -1
	ELobbyComparisonEqual                ELobbyComparison = 0
	ELobbyComparisonGreaterThan          ELobbyComparison = 1
	ELobbyComparisonEqualToOrGreaterThan ELobbyComparison = 2
	ELobbyComparisonNotEqual             ELobbyComparison = 3
)

type ELobbyDistanceFilter int32

const (
	ELobbyDistanceFilterClose      ELobbyDistanceFilter = 0
	ELobbyDistanceFilterDefault    ELobbyDistanceFilter = 1
	ELobbyDistanceFilterFar        ELobbyDistanceFilter = 2
	ELobbyDistanceFilterWorldwide  ELobbyDistanceFilter = 3
	ELobbyDistanceFilterCompatible ELobbyDistanceFilter = 4
)

type EChatEntryType int32

const (
	EChatEntryTypeInvalid          EChatEntryType = 0
	EChatEntryTypeChatMsg          EChatEntryType = 1
	EChatEntryTypeTyping           EChatEntryType = 2
	EChatEntryTypeInviteGame       EChatEntryType = 3
	EChatEntryTypeEmote            EChatEntryType = 4
	EChatEntryTypeLeftConversation EChatEntryType = 6
	EChatEntryTypeEntered          EChatEntryType = 7
	EChatEntryTypeWasKicked        EChatEntryType = 8
	EChatEntryTypeWasBanned        EChatEntryType = 9
	EChatEntryTypeDisconnected     EChatEntryType = 10
	EChatEntryTypeHistoricalChat   EChatEntryType = 11
	EChatEntryTypeLinkBlocked      EChatEntryType = 14
)

// Steam matchmaking callback IDs for lobby events.
const (
	CallbackIDLobbyDataUpdate CallbackID = 505
	CallbackIDLobbyChatUpdate CallbackID = 506
	CallbackIDLobbyChatMsg    CallbackID = 507

	// CallbackIDSteamRemotePlaySessionAvatarLoaded mirrors SteamRemotePlaySessionAvatarLoaded_t::k_iCallback.
	CallbackIDSteamRemotePlaySessionAvatarLoaded CallbackID = 5704
)

// SteamRemotePlaySessionAvatarLoaded mirrors Steam's SteamRemotePlaySessionAvatarLoaded_t callback payload.
type SteamRemotePlaySessionAvatarLoaded struct {
	SessionID uint32
	Image     int32
	Wide      int32
	Height    int32
}

// LobbyDataUpdate mirrors Steam's LobbyDataUpdate_t callback payload.
type LobbyDataUpdate struct {
	LobbySteamID  CSteamID
	MemberSteamID CSteamID
	Success       uint8
}

// LobbyChatUpdate mirrors Steam's LobbyChatUpdate_t callback payload.
type LobbyChatUpdate struct {
	LobbySteamID          CSteamID
	UserChangedSteamID    CSteamID
	MakingChangeSteamID   CSteamID
	ChatMemberStateChange uint32
}

// LobbyChatMsg mirrors Steam's LobbyChatMsg_t callback payload.
type LobbyChatMsg struct {
	LobbySteamID  CSteamID
	UserSteamID   CSteamID
	ChatEntryType uint8
	_             [3]byte
	ChatID        int32
}

// EntryType returns the callback's chat entry type as EChatEntryType.
func (m LobbyChatMsg) EntryType() EChatEntryType {
	return EChatEntryType(m.ChatEntryType)
}

type EResult int32

const (
	EResultNone         EResult = 0
	EResultOK           EResult = 1
	EResultFail         EResult = 2
	EResultNoConnection EResult = 3
	EResultInvalidParam EResult = 8
	EResultBusy         EResult = 10
	EResultTimeout      EResult = 16
)

type SteamNetworkingSendFlags int32

const (
	SteamNetworkingSend_Unreliable               SteamNetworkingSendFlags = 0
	SteamNetworkingSend_NoNagle                  SteamNetworkingSendFlags = 1
	SteamNetworkingSend_NoDelay                  SteamNetworkingSendFlags = 4
	SteamNetworkingSend_Reliable                 SteamNetworkingSendFlags = 8
	SteamNetworkingSend_UseCurrentThread         SteamNetworkingSendFlags = 16
	SteamNetworkingSend_AutoRestartBrokenSession SteamNetworkingSendFlags = 32
)

type FriendGameInfo struct {
	GameID       CGameID
	GameIP       uint32
	GamePort     uint16
	QueryPort    uint16
	LobbySteamID CSteamID
}

type FavoriteGame struct {
	AppID                  AppId_t
	IP                     uint32
	ConnectionPort         uint16
	QueryPort              uint16
	Flags                  uint32
	LastPlayedOnServerTime uint32
}

type DurationControl struct {
	Progress         EDurationControlProgress
	Notification     EDurationControlNotification
	OnlineState      EDurationControlOnlineState
	SecondsRemaining uint32
	SecondsPlayed    uint32
}

type EHTTPMethod int32

const (
	EHTTPMethodInvalid EHTTPMethod = 0
	EHTTPMethodGET     EHTTPMethod = 1
	EHTTPMethodHEAD    EHTTPMethod = 2
	EHTTPMethodPOST    EHTTPMethod = 3
	EHTTPMethodPUT     EHTTPMethod = 4
	EHTTPMethodDELETE  EHTTPMethod = 5
	EHTTPMethodOPTIONS EHTTPMethod = 6
	EHTTPMethodPATCH   EHTTPMethod = 7
)

type SteamItemDetails struct {
	ItemID     SteamItemInstanceID_t
	Definition SteamItemDef_t
	Quantity   uint16
	Flags      uint16
}

type EFriendRelationship int32

const (
	EFriendRelationshipNone                EFriendRelationship = 0
	EFriendRelationshipBlocked             EFriendRelationship = 1
	EFriendRelationshipRequestRecipient    EFriendRelationship = 2
	EFriendRelationshipFriend              EFriendRelationship = 3
	EFriendRelationshipRequestInitiator    EFriendRelationship = 4
	EFriendRelationshipIgnored             EFriendRelationship = 5
	EFriendRelationshipIgnoredFriend       EFriendRelationship = 6
	EFriendRelationshipSuggestedDeprecated EFriendRelationship = 7
	EFriendRelationshipMax                 EFriendRelationship = 8
)

type EPersonaState int32

const (
	EPersonaStateOffline        EPersonaState = 0
	EPersonaStateOnline         EPersonaState = 1
	EPersonaStateBusy           EPersonaState = 2
	EPersonaStateAway           EPersonaState = 3
	EPersonaStateSnooze         EPersonaState = 4
	EPersonaStateLookingToTrade EPersonaState = 5
	EPersonaStateLookingToPlay  EPersonaState = 6
	EPersonaStateInvisible      EPersonaState = 7
	EPersonaStateMax            EPersonaState = 8
)

type EFriendFlags int32

const (
	EFriendFlagNone                 EFriendFlags = 0x00
	EFriendFlagBlocked              EFriendFlags = 0x01
	EFriendFlagFriendshipRequested  EFriendFlags = 0x02
	EFriendFlagImmediate            EFriendFlags = 0x04
	EFriendFlagClanMember           EFriendFlags = 0x08
	EFriendFlagOnGameServer         EFriendFlags = 0x10
	EFriendFlagRequestingFriendship EFriendFlags = 0x80
	EFriendFlagRequestingInfo       EFriendFlags = 0x100
	EFriendFlagIgnored              EFriendFlags = 0x200
	EFriendFlagIgnoredFriend        EFriendFlags = 0x400
	EFriendFlagChatMember           EFriendFlags = 0x1000
	EFriendFlagAll                  EFriendFlags = 0xFFFF
)

type ESteamAPICallFailure int32

const (
	ESteamAPICallFailureNone               ESteamAPICallFailure = -1
	ESteamAPICallFailureSteamGone          ESteamAPICallFailure = 0
	ESteamAPICallFailureNetworkFailure     ESteamAPICallFailure = 1
	ESteamAPICallFailureInvalidHandle      ESteamAPICallFailure = 2
	ESteamAPICallFailureMismatchedCallback ESteamAPICallFailure = 3
)

type EUniverse int32

const (
	EUniverseInvalid  EUniverse = 0
	EUniversePublic   EUniverse = 1
	EUniverseBeta     EUniverse = 2
	EUniverseInternal EUniverse = 3
	EUniverseDev      EUniverse = 4
	EUniverseMax      EUniverse = 5
)

type ENotificationPosition int32

const (
	ENotificationPositionInvalid     ENotificationPosition = -1
	ENotificationPositionTopLeft     ENotificationPosition = 0
	ENotificationPositionTopRight    ENotificationPosition = 1
	ENotificationPositionBottomLeft  ENotificationPosition = 2
	ENotificationPositionBottomRight ENotificationPosition = 3
)

type SteamNetworkingIPAddr struct {
	data [18]byte
}

func (a *SteamNetworkingIPAddr) SetIPv4(ip uint32, port uint16) {
	a.setIPv4(ip, port)
}

func (a *SteamNetworkingIPAddr) SetIPv6(ip [16]byte, port uint16) {
	copy(a.data[:16], ip[:])
	a.data[16] = byte(port)
	a.data[17] = byte(port >> 8)
}

func (a *SteamNetworkingIPAddr) setIPv4(ip uint32, port uint16) {
	a.data[0] = 0
	a.data[1] = 0
	a.data[2] = 0
	a.data[3] = 0
	a.data[4] = 0
	a.data[5] = 0
	a.data[6] = 0
	a.data[7] = 0
	a.data[8] = 0
	a.data[9] = 0
	a.data[10] = 0xff
	a.data[11] = 0xff
	a.data[12] = byte(ip >> 24)
	a.data[13] = byte(ip >> 16)
	a.data[14] = byte(ip >> 8)
	a.data[15] = byte(ip)
	a.data[16] = byte(port)
	a.data[17] = byte(port >> 8)
}

type ESteamNetworkingIdentityType int32

const (
	SteamNetworkingIdentity_Invalid       ESteamNetworkingIdentityType = 0
	SteamNetworkingIdentity_SteamID       ESteamNetworkingIdentityType = 16
	SteamNetworkingIdentity_IPAddress     ESteamNetworkingIdentityType = 1
	SteamNetworkingIdentity_GenericString ESteamNetworkingIdentityType = 2
	SteamNetworkingIdentity_GenericBytes  ESteamNetworkingIdentityType = 3
)

type SteamNetworkingIdentity struct {
	data [136]byte
}

func (i *SteamNetworkingIdentity) SetSteamID64(steamID uint64) {
	i.setTypeAndSize(SteamNetworkingIdentity_SteamID, 8)
	i.setUint64(8, steamID)
}

func (i *SteamNetworkingIdentity) SetSteamID(steamID CSteamID) {
	i.SetSteamID64(uint64(steamID))
}

func (i *SteamNetworkingIdentity) SetIPv4Addr(ip uint32, port uint16) {
	i.setTypeAndSize(SteamNetworkingIdentity_IPAddress, 18)
	var addr SteamNetworkingIPAddr
	addr.setIPv4(ip, port)
	copy(i.data[8:], addr.data[:])
}

func (i *SteamNetworkingIdentity) setTypeAndSize(identityType ESteamNetworkingIdentityType, size int32) {
	putUint32(i.data[0:4], uint32(identityType))
	putUint32(i.data[4:8], uint32(size))
}

func (i *SteamNetworkingIdentity) setUint64(offset int, value uint64) {
	putUint64(i.data[offset:offset+8], value)
}

type SteamNetworkingMessage struct {
	Data          uintptr
	Size          int32
	Connection    HSteamNetConnection
	IdentityPeer  SteamNetworkingIdentity
	ConnUserData  int64
	TimeReceived  SteamNetworkingMicroseconds
	MessageNumber int64
	FreeDataFunc  uintptr
	ReleaseFunc   uintptr
	Channel       int32
	Flags         int32
	UserData      int64
	Lane          uint16
	_pad1         uint16
}

// ESteamNetworkingConnectionState mirrors the Steamworks SDK enum.
type ESteamNetworkingConnectionState int32

const (
	ESteamNetworkingConnectionState_None                ESteamNetworkingConnectionState = 0
	ESteamNetworkingConnectionState_Connecting          ESteamNetworkingConnectionState = 1
	ESteamNetworkingConnectionState_FindingRoute        ESteamNetworkingConnectionState = 2
	ESteamNetworkingConnectionState_Connected           ESteamNetworkingConnectionState = 3
	ESteamNetworkingConnectionState_ClosedByPeer        ESteamNetworkingConnectionState = 4
	ESteamNetworkingConnectionState_ProblemDetectedLocally ESteamNetworkingConnectionState = 5
)

// SteamNetworkingPOPID is a Point of Presence identifier.
type SteamNetworkingPOPID uint32

// SteamNetConnectionInfo_t mirrors the Steamworks SDK struct (SDK 1.64).
// Total size: ~696 bytes.
type SteamNetConnectionInfo_t struct {
	IdentityRemote SteamNetworkingIdentity         // 136 bytes
	UserData       int64                            // 8 bytes
	ListenSocket   HSteamListenSocket               // 4 bytes
	AddrRemote     SteamNetworkingIPAddr            // 18 bytes
	_pad1          uint16                           // 2 bytes
	IDPOPRemote    SteamNetworkingPOPID             // 4 bytes
	IDPOPRelay     SteamNetworkingPOPID             // 4 bytes
	State          ESteamNetworkingConnectionState  // 4 bytes  (offset 176)
	EndReason      int32                            // 4 bytes
	EndDebug       [128]byte                        // 128 bytes (k_cchSteamNetworkingMaxConnectionCloseReason)
	Description    [128]byte                        // 128 bytes (k_cchSteamNetworkingMaxConnectionDescription)
	Flags          int32                            // 4 bytes
	_reserved      [63]uint32                       // 252 bytes
}

// SteamNetConnectionStatusChangedCallback_t is delivered via callback ID 1221
// when the state of a Steam Networking connection changes.
type SteamNetConnectionStatusChangedCallback_t struct {
	Conn     HSteamNetConnection             // 4 bytes
	_pad     [4]byte                         // 4 bytes alignment
	Info     SteamNetConnectionInfo_t        // ~696 bytes
	OldState ESteamNetworkingConnectionState // 4 bytes
}

// CallbackID for SteamNetConnectionStatusChangedCallback_t.
const CallbackIDSteamNetConnectionStatusChanged CallbackID = 1221

type ISteamApps interface {
	BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool)
	BIsSubscribed() bool
	BIsLowViolence() bool
	BIsCybercafe() bool
	BIsVACBanned() bool
	BIsDlcInstalled(appID AppId_t) bool
	BIsSubscribedApp(appID AppId_t) bool
	BIsSubscribedFromFreeWeekend() bool
	BIsSubscribedFromFamilySharing() bool
	BIsTimedTrial() (allowedSeconds, playedSeconds uint32, ok bool)
	BIsAppInstalled(appID AppId_t) bool
	GetAvailableGameLanguages() string
	GetEarliestPurchaseUnixTime(appID AppId_t) uint32
	GetAppInstallDir(appID AppId_t) string
	GetCurrentGameLanguage() string
	GetDLCCount() int32
	GetCurrentBetaName() (string, bool)
	GetInstalledDepots(appID AppId_t) []DepotId_t
	GetAppOwner() CSteamID
	GetLaunchQueryParam(key string) string
	GetDlcDownloadProgress(appID AppId_t) (downloaded, total uint64, ok bool)
	GetAppBuildId() int32
	GetFileDetails(filename string) SteamAPICall_t
	GetLaunchCommandLine(bufferSize int) string
	GetNumBetas() (total int, available int, private int)
	GetBetaInfo(index int) (flags uint32, buildID uint32, lastUpdated uint32, name string, description string, ok bool)
	InstallDLC(appID AppId_t)
	UninstallDLC(appID AppId_t)
	RequestAppProofOfPurchaseKey(appID AppId_t)
	RequestAllProofOfPurchaseKeys()
	MarkContentCorrupt(missingFilesOnly bool) bool
	SetDlcContext(appID AppId_t) bool
	SetActiveBeta(name string) bool
}

type ISteamHTTP interface {
	CreateHTTPRequest(method EHTTPMethod, absoluteURL string) HTTPRequestHandle
	SetHTTPRequestHeaderValue(request HTTPRequestHandle, headerName, headerValue string) bool
	SendHTTPRequest(request HTTPRequestHandle) (SteamAPICall_t, bool)
	GetHTTPResponseBodySize(request HTTPRequestHandle) (uint32, bool)
	GetHTTPResponseBodyData(request HTTPRequestHandle, buffer []byte) bool
	ReleaseHTTPRequest(request HTTPRequestHandle) bool
}

type ISteamUGC interface {
	GetNumSubscribedItems(includeLocallyDisabled bool) uint32
	GetSubscribedItems(includeLocallyDisabled bool) []PublishedFileId_t
	MarkDownloadedItemAsUnused(publishedFileID PublishedFileId_t) bool
	GetNumDownloadedItems() uint32
	GetDownloadedItems() []PublishedFileId_t
}

type ISteamInventory interface {
	GetResultStatus(result SteamInventoryResult_t) EResult
	GetResultItems(result SteamInventoryResult_t, outItems []SteamItemDetails) (int, bool)
	DestroyResult(result SteamInventoryResult_t)
}

type ISteamNetworkingUtils interface {
	AllocateMessage(size int) *SteamNetworkingMessage
	InitRelayNetworkAccess()
	GetLocalTimestamp() SteamNetworkingMicroseconds
}

type ISteamGameServer interface {
	AssociateWithClan(clanID CSteamID) SteamAPICall_t
	BeginAuthSession(authTicket []byte, steamID CSteamID) EBeginAuthSessionResult
	BLoggedOn() bool
	BSecure() bool
	BUpdateUserData(steamIDUser CSteamID, playerName string, score uint32) bool
	CancelAuthTicket(authTicket HAuthTicket)
	ClearAllKeyValues()
	ComputeNewPlayerCompatibility(steamIDNewPlayer CSteamID, steamIDPlayers []CSteamID, steamIDPlayersInGame []CSteamID, steamIDTeamPlayers []CSteamID) SteamAPICall_t
	CreateUnauthenticatedUserConnection() CSteamID
	EnableHeartbeats(active bool)
	EndAuthSession(steamID CSteamID)
	ForceHeartbeat()
	GetAuthSessionTicket(authTicket []byte) (ticket HAuthTicket, size uint32)
	GetGameplayStats()
	GetNextOutgoingPacket(dest []byte) (size int32, ip uint32, port uint16)
	GetPublicIP() uint32
	GetServerReputation() SteamAPICall_t
	GetSteamID() CSteamID
	HandleIncomingPacket(data []byte, ip uint32, port uint16) bool
	InitGameServer(ip uint32, steamPort uint16, gamePort uint16, queryPort uint16, serverMode uint32, versionString string) bool
	LogOff()
	LogOn(token string)
	LogOnAnonymous()
	RequestUserGroupStatus(steamIDUser CSteamID, steamIDGroup CSteamID) bool
	SendUserConnectAndAuthenticate(ipClient uint32, authBlob []byte) (steamIDUser CSteamID, ok bool)
	SendUserDisconnect(steamIDUser CSteamID)
	SetBotPlayerCount(botPlayers int32)
	SetDedicatedServer(dedicated bool)
	SetGameData(gameData string)
	SetGameDescription(description string)
	SetGameTags(gameTags string)
	SetHeartbeatInterval(interval int)
	SetKeyValue(key string, value string)
	SetMapName(mapName string)
	SetMaxPlayerCount(playersMax int32)
	SetModDir(modDir string)
	SetPasswordProtected(passwordProtected bool)
	SetProduct(product string)
	SetRegion(region string)
	SetServerName(serverName string)
	SetSpectatorPort(spectatorPort uint16)
	SetSpectatorServerName(spectatorServerName string)
	UserHasLicenseForApp(steamID CSteamID, appID AppId_t) EUserHasLicenseForAppResult
	WasRestartRequested() bool
}

type ISteamInput interface {
	GetConnectedControllers() []InputHandle_t
	ConnectedControllers() iter.Seq[InputHandle_t]
	GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType
	Init(bExplicitlyCallRunFrame bool) bool
	Shutdown()
	RunFrame()
	EnableDeviceCallbacks()
	GetActionSetHandle(actionSetName string) InputActionSetHandle_t
	ActivateActionSet(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t)
	GetCurrentActionSet(inputHandle InputHandle_t) InputActionSetHandle_t
	ActivateActionSetLayer(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t)
	DeactivateActionSetLayer(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t)
	DeactivateAllActionSetLayers(inputHandle InputHandle_t)
	GetActiveActionSetLayers(inputHandle InputHandle_t, handles []InputActionSetHandle_t) int
	GetDigitalActionHandle(actionName string) InputDigitalActionHandle_t
	GetDigitalActionData(inputHandle InputHandle_t, actionHandle InputDigitalActionHandle_t) InputDigitalActionData
	GetDigitalActionOrigins(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t, actionHandle InputDigitalActionHandle_t, origins []EInputActionOrigin) int
	GetAnalogActionHandle(actionName string) InputAnalogActionHandle_t
	GetAnalogActionData(inputHandle InputHandle_t, actionHandle InputAnalogActionHandle_t) InputAnalogActionData
	GetAnalogActionOrigins(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t, actionHandle InputAnalogActionHandle_t, origins []EInputActionOrigin) int
	StopAnalogActionMomentum(inputHandle InputHandle_t, actionHandle InputAnalogActionHandle_t)
	GetMotionData(inputHandle InputHandle_t) InputMotionData
	TriggerVibration(inputHandle InputHandle_t, leftSpeed, rightSpeed uint16)
	TriggerVibrationExtended(inputHandle InputHandle_t, leftSpeed, rightSpeed, leftTriggerSpeed, rightTriggerSpeed uint16)
	TriggerSimpleHapticEvent(inputHandle InputHandle_t, pad ESteamControllerPad, durationMicroSec, offMicroSec, repeat uint16)
	SetLEDColor(inputHandle InputHandle_t, red, green, blue uint8, flags ESteamInputLEDFlag)
	ShowBindingPanel(inputHandle InputHandle_t) bool
	GetControllerForGamepadIndex(index int) InputHandle_t
	GetGamepadIndexForController(inputHandle InputHandle_t) int
	GetStringForActionOrigin(origin EInputActionOrigin) string
	GetGlyphForActionOrigin(origin EInputActionOrigin) string
	GetRemotePlaySessionID(inputHandle InputHandle_t) uint32
}

type ISteamRemoteStorage interface {
	FileWrite(file string, data []byte) bool
	FileRead(file string, data []byte) int32
	FileDelete(file string) bool
	GetFileSize(file string) int32
}

type ISteamUser interface {
	AdvertiseGame(gameServerSteamID CSteamID, ip uint32, port uint16)
	BeginAuthSession(authTicket []byte, steamID CSteamID) EBeginAuthSessionResult
	BIsBehindNAT() bool
	BIsPhoneIdentifying() bool
	BIsPhoneRequiringVerification() bool
	BIsPhoneVerified() bool
	BIsTwoFactorEnabled() bool
	BLoggedOn() bool
	BSetDurationControlOnlineState(newState EDurationControlOnlineState) bool
	CancelAuthTicket(authTicket HAuthTicket)
	DecompressVoice(compressedData []byte, destBuffer []byte, desiredSampleRate uint32) (bytesWritten uint32, result EVoiceResult)
	EndAuthSession(steamID CSteamID)
	GetAuthSessionTicket(authTicket []byte, identityRemote *SteamNetworkingIdentity) (ticket HAuthTicket, size uint32)
	GetAuthTicketForWebApi(identity string) HAuthTicket
	GetAvailableVoice() (compressedBytes uint32, uncompressedBytes uint32, result EVoiceResult)
	GetDurationControl() (control DurationControl, ok bool)
	GetEncryptedAppTicket(ticket []byte) (ticketSize uint32, ok bool)
	GetGameBadgeLevel(series int32, foil bool) int32
	GetHSteamUser() HSteamUser
	GetPlayerSteamLevel() int32
	GetSteamID() CSteamID
	GetUserDataFolder() (path string, ok bool)
	GetVoice(wantCompressed bool, compressedData []byte, wantUncompressed bool, uncompressedData []byte, desiredSampleRate uint32) (compressedBytes uint32, uncompressedBytes uint32, result EVoiceResult)
	GetVoiceOptimalSampleRate() uint32
	InitiateGameConnection(authBlob []byte, steamIDGameServer CSteamID, ipServer uint32, portServer uint16, secure bool) int32
	RequestEncryptedAppTicket(dataToInclude []byte) SteamAPICall_t
	RequestStoreAuthURL(redirectURL string) SteamAPICall_t
	StartVoiceRecording()
	StopVoiceRecording()
	TerminateGameConnection(ipServer uint32, portServer uint16)
	TrackAppUsageEvent(gameID CGameID, eventCode int32, extraInfo string)
	UserHasLicenseForApp(steamID CSteamID, appID AppId_t) EUserHasLicenseForAppResult
}

type ISteamUserStats interface {
	GetAchievement(name string) (achieved, success bool)
	SetAchievement(name string) bool
	ClearAchievement(name string) bool
	StoreStats() bool
}

type ISteamUtils interface {
	GetSecondsSinceAppActive() uint32
	GetSecondsSinceComputerActive() uint32
	GetConnectedUniverse() EUniverse
	GetServerRealTime() uint32
	GetIPCountry() string
	GetImageSize(image int) (width, height uint32, ok bool)
	GetImageRGBA(image int, dest []byte) bool
	GetCurrentBatteryPower() uint8
	GetAppID() uint32
	IsOverlayEnabled() bool
	BOverlayNeedsPresent() bool
	IsSteamRunningOnSteamDeck() bool
	SetOverlayNotificationPosition(position ENotificationPosition)
	SetOverlayNotificationInset(horizontal, vertical int32)
	IsAPICallCompleted(call SteamAPICall_t) (failed bool, ok bool)
	GetAPICallFailureReason(call SteamAPICall_t) ESteamAPICallFailure
	GetAPICallResult(call SteamAPICall_t, callback uintptr, callbackSize int32, expectedCallback int32) (failed bool, ok bool)
	GetIPCCallCount() uint32
	ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool
}

type ISteamFriends interface {
	GetPersonaName() string
	GetPersonaState() EPersonaState
	GetFriendCount(flags EFriendFlags) int
	GetFriendByIndex(index int, flags EFriendFlags) CSteamID
	Friends(flags EFriendFlags) iter.Seq[CSteamID]
	GetFriendRelationship(friend CSteamID) EFriendRelationship
	GetFriendPersonaState(friend CSteamID) EPersonaState
	GetFriendPersonaName(friend CSteamID) string
	GetFriendPersonaNameHistory(friend CSteamID, index int) string
	GetFriendSteamLevel(friend CSteamID) int
	GetSmallFriendAvatar(friend CSteamID) int32
	GetMediumFriendAvatar(friend CSteamID) int32
	GetLargeFriendAvatar(friend CSteamID) int32
	SetRichPresence(string, string) bool
	GetFriendGamePlayed(friend CSteamID) (FriendGameInfo, bool)
	InviteUserToGame(friend CSteamID, connectString string) bool
	ActivateGameOverlay(dialog string)
	ActivateGameOverlayToUser(dialog string, steamID CSteamID)
	ActivateGameOverlayToWebPage(url string, mode EActivateGameOverlayToWebPageMode)
	ActivateGameOverlayToStore(appID AppId_t, flag EOverlayToStoreFlag)
	ActivateGameOverlayInviteDialog(lobbyID CSteamID)
	ActivateGameOverlayInviteDialogConnectString(connectString string)
}

type ISteamMatchmaking interface {
	GetFavoriteGameCount() int
	GetFavoriteGame(index int) (FavoriteGame, bool)
	AddFavoriteGame(appID AppId_t, ip uint32, connectionPort, queryPort uint16, flags, lastPlayedOnServerTime uint32) int
	RemoveFavoriteGame(appID AppId_t, ip uint32, connectionPort, queryPort uint16, flags uint32) bool

	RequestLobbyList() SteamAPICall_t
	AddRequestLobbyListStringFilter(key, value string, comparisonType ELobbyComparison)
	AddRequestLobbyListNumericalFilter(key string, value int, comparisonType ELobbyComparison)
	AddRequestLobbyListNearValueFilter(key string, value int)
	AddRequestLobbyListFilterSlotsAvailable(slotsAvailable int)
	AddRequestLobbyListDistanceFilter(distanceFilter ELobbyDistanceFilter)
	AddRequestLobbyListResultCountFilter(maxResults int)
	AddRequestLobbyListCompatibleMembersFilter(lobbyID CSteamID)

	GetLobbyByIndex(index int) CSteamID
	CreateLobby(lobbyType ELobbyType, maxMembers int) SteamAPICall_t
	JoinLobby(lobbyID CSteamID) SteamAPICall_t
	LeaveLobby(lobbyID CSteamID)
	InviteUserToLobby(lobbyID, invitee CSteamID) bool
	SetLobbyMemberLimit(lobbyID CSteamID, maxMembers int) bool
	GetLobbyMemberLimit(lobbyID CSteamID) int
	SetLobbyType(lobbyID CSteamID, lobbyType ELobbyType) bool
	SetLobbyJoinable(lobbyID CSteamID, joinable bool) bool
	GetLobbyOwner(lobbyID CSteamID) CSteamID
	SetLobbyOwner(lobbyID, owner CSteamID) bool
	SetLinkedLobby(lobbyID, lobbyDependent CSteamID) bool

	GetNumLobbyMembers(lobbyID CSteamID) int
	GetLobbyMemberByIndex(lobbyID CSteamID, memberIndex int) CSteamID
	LobbyMembers(lobbyID CSteamID) iter.Seq[CSteamID]
	SetLobbyData(lobbyID CSteamID, key, value string) bool
	GetLobbyData(lobbyID CSteamID, key string) string
	DeleteLobbyData(lobbyID CSteamID, key string) bool
	GetLobbyDataCount(lobbyID CSteamID) int
	GetLobbyDataByIndex(lobbyID CSteamID, lobbyDataIndex int) (key, value string, ok bool)
	SetLobbyMemberData(lobbyID CSteamID, key, value string)
	GetLobbyMemberData(lobbyID, user CSteamID, key string) string
	SendLobbyChatMsg(lobbyID CSteamID, msgBody []byte) bool
	GetLobbyChatEntry(lobbyID CSteamID, chatID int, data []byte) (user CSteamID, entryType EChatEntryType, bytesCopied int)
	RequestLobbyData(lobbyID CSteamID) bool

	SetLobbyGameServer(lobbyID CSteamID, ip uint32, port uint16, server CSteamID)
	GetLobbyGameServer(lobbyID CSteamID) (ip uint32, port uint16, server CSteamID, ok bool)
	CheckForPSNGameBootInvite(lobbyID *CSteamID) bool
}

type ISteamNetworkingMessages interface {
	SendMessageToUser(identity *SteamNetworkingIdentity, data []byte, sendFlags SteamNetworkingSendFlags, remoteChannel int) EResult
	ReceiveMessagesOnChannel(channel int, maxMessages int) []*SteamNetworkingMessage
	AcceptSessionWithUser(identity *SteamNetworkingIdentity) bool
	CloseSessionWithUser(identity *SteamNetworkingIdentity) bool
	CloseChannelWithUser(identity *SteamNetworkingIdentity, channel int) bool
}

type ISteamNetworkingSockets interface {
	CreateListenSocketIP(localAddress *SteamNetworkingIPAddr, options []SteamNetworkingConfigValue) HSteamListenSocket
	CreateListenSocketP2P(localVirtualPort int, options []SteamNetworkingConfigValue) HSteamListenSocket
	ConnectByIPAddress(address *SteamNetworkingIPAddr, options []SteamNetworkingConfigValue) HSteamNetConnection
	ConnectP2P(identity *SteamNetworkingIdentity, remoteVirtualPort int, options []SteamNetworkingConfigValue) HSteamNetConnection
	AcceptConnection(connection HSteamNetConnection) EResult
	CloseConnection(connection HSteamNetConnection, reason int, debug string, enableLinger bool) bool
	CloseListenSocket(socket HSteamListenSocket) bool
	SendMessageToConnection(connection HSteamNetConnection, data []byte, sendFlags SteamNetworkingSendFlags) (EResult, int64)
	ReceiveMessagesOnConnection(connection HSteamNetConnection, maxMessages int) []*SteamNetworkingMessage
	CreatePollGroup() HSteamNetPollGroup
	DestroyPollGroup(group HSteamNetPollGroup) bool
	SetConnectionPollGroup(connection HSteamNetConnection, group HSteamNetPollGroup) bool
	ReceiveMessagesOnPollGroup(group HSteamNetPollGroup, maxMessages int) []*SteamNetworkingMessage
}

type SteamNetworkingConfigValue struct {
	Value    int32
	DataType int32
	Data     uint64
}

type ISteamAppTicket struct{ ptr uintptr }

func (i ISteamAppTicket) Ptr() uintptr { return i.ptr }
func (i ISteamAppTicket) Valid() bool  { return i.ptr != 0 }

type ISteamClient struct{ ptr uintptr }

func (i ISteamClient) Ptr() uintptr { return i.ptr }
func (i ISteamClient) Valid() bool  { return i.ptr != 0 }

type ISteamController struct{ ptr uintptr }

func (i ISteamController) Ptr() uintptr { return i.ptr }
func (i ISteamController) Valid() bool  { return i.ptr != 0 }

type ISteamGameCoordinator struct{ ptr uintptr }

func (i ISteamGameCoordinator) Ptr() uintptr { return i.ptr }
func (i ISteamGameCoordinator) Valid() bool  { return i.ptr != 0 }

type ISteamGameServerStats struct{ ptr uintptr }

func (i ISteamGameServerStats) Ptr() uintptr { return i.ptr }
func (i ISteamGameServerStats) Valid() bool  { return i.ptr != 0 }

type ISteamHTMLSurface struct{ ptr uintptr }

func (i ISteamHTMLSurface) Ptr() uintptr { return i.ptr }
func (i ISteamHTMLSurface) Valid() bool  { return i.ptr != 0 }

type ISteamMatchmakingServers struct{ ptr uintptr }

func (i ISteamMatchmakingServers) Ptr() uintptr { return i.ptr }
func (i ISteamMatchmakingServers) Valid() bool  { return i.ptr != 0 }

type ISteamMusic struct{ ptr uintptr }

func (i ISteamMusic) Ptr() uintptr { return i.ptr }
func (i ISteamMusic) Valid() bool  { return i.ptr != 0 }

type ISteamNetworking struct{ ptr uintptr }

func (i ISteamNetworking) Ptr() uintptr { return i.ptr }
func (i ISteamNetworking) Valid() bool  { return i.ptr != 0 }

type ISteamRemotePlay struct{ ptr uintptr }

func (i ISteamRemotePlay) Ptr() uintptr { return i.ptr }
func (i ISteamRemotePlay) Valid() bool  { return i.ptr != 0 }

type ISteamScreenshots struct{ ptr uintptr }

func (i ISteamScreenshots) Ptr() uintptr { return i.ptr }
func (i ISteamScreenshots) Valid() bool  { return i.ptr != 0 }

type ISteamTimeline struct{ ptr uintptr }

func (i ISteamTimeline) Ptr() uintptr { return i.ptr }
func (i ISteamTimeline) Valid() bool  { return i.ptr != 0 }

type ISteamVideo struct{ ptr uintptr }

func (i ISteamVideo) Ptr() uintptr { return i.ptr }
func (i ISteamVideo) Valid() bool  { return i.ptr != 0 }

type ISteamAPIClient struct{ ptr uintptr }

func (i ISteamAPIClient) Ptr() uintptr { return i.ptr }
func (i ISteamAPIClient) Valid() bool  { return i.ptr != 0 }

type ISteamAPIGameServer struct{ ptr uintptr }

func (i ISteamAPIGameServer) Ptr() uintptr { return i.ptr }
func (i ISteamAPIGameServer) Valid() bool  { return i.ptr != 0 }

const (
	flatAPI_RestartAppIfNecessary      = "SteamAPI_RestartAppIfNecessary"
	flatAPI_InitFlat                   = "SteamAPI_InitFlat"
	flatAPI_RunCallbacks               = "SteamAPI_RunCallbacks"
	flatAPI_Shutdown                   = "SteamAPI_Shutdown"
	flatAPI_IsSteamRunning             = "SteamAPI_IsSteamRunning"
	flatAPI_GetSteamInstallPath        = "SteamAPI_GetSteamInstallPath"
	flatAPI_ReleaseCurrentThreadMemory = "SteamAPI_ReleaseCurrentThreadMemory"

	flatAPI_SteamApps                                 = "SteamAPI_SteamApps_v008"
	flatAPI_SteamAppsV009                             = "SteamAPI_SteamApps_v009"
	flatAPI_SteamAppsUnversioned                      = "SteamAPI_SteamApps"
	flatAPI_ISteamApps_BIsSubscribed                  = "SteamAPI_ISteamApps_BIsSubscribed"
	flatAPI_ISteamApps_BIsLowViolence                 = "SteamAPI_ISteamApps_BIsLowViolence"
	flatAPI_ISteamApps_BIsCybercafe                   = "SteamAPI_ISteamApps_BIsCybercafe"
	flatAPI_ISteamApps_BIsVACBanned                   = "SteamAPI_ISteamApps_BIsVACBanned"
	flatAPI_ISteamApps_BGetDLCDataByIndex             = "SteamAPI_ISteamApps_BGetDLCDataByIndex"
	flatAPI_ISteamApps_BIsDlcInstalled                = "SteamAPI_ISteamApps_BIsDlcInstalled"
	flatAPI_ISteamApps_GetAvailableGameLanguages      = "SteamAPI_ISteamApps_GetAvailableGameLanguages"
	flatAPI_ISteamApps_BIsSubscribedApp               = "SteamAPI_ISteamApps_BIsSubscribedApp"
	flatAPI_ISteamApps_GetEarliestPurchaseUnixTime    = "SteamAPI_ISteamApps_GetEarliestPurchaseUnixTime"
	flatAPI_ISteamApps_BIsSubscribedFromFreeWeekend   = "SteamAPI_ISteamApps_BIsSubscribedFromFreeWeekend"
	flatAPI_ISteamApps_GetAppInstallDir               = "SteamAPI_ISteamApps_GetAppInstallDir"
	flatAPI_ISteamApps_GetCurrentGameLanguage         = "SteamAPI_ISteamApps_GetCurrentGameLanguage"
	flatAPI_ISteamApps_GetDLCCount                    = "SteamAPI_ISteamApps_GetDLCCount"
	flatAPI_ISteamApps_InstallDLC                     = "SteamAPI_ISteamApps_InstallDLC"
	flatAPI_ISteamApps_UninstallDLC                   = "SteamAPI_ISteamApps_UninstallDLC"
	flatAPI_ISteamApps_RequestAppProofOfPurchaseKey   = "SteamAPI_ISteamApps_RequestAppProofOfPurchaseKey"
	flatAPI_ISteamApps_GetCurrentBetaName             = "SteamAPI_ISteamApps_GetCurrentBetaName"
	flatAPI_ISteamApps_MarkContentCorrupt             = "SteamAPI_ISteamApps_MarkContentCorrupt"
	flatAPI_ISteamApps_GetInstalledDepots             = "SteamAPI_ISteamApps_GetInstalledDepots"
	flatAPI_ISteamApps_BIsAppInstalled                = "SteamAPI_ISteamApps_BIsAppInstalled"
	flatAPI_ISteamApps_GetAppOwner                    = "SteamAPI_ISteamApps_GetAppOwner"
	flatAPI_ISteamApps_GetLaunchQueryParam            = "SteamAPI_ISteamApps_GetLaunchQueryParam"
	flatAPI_ISteamApps_GetDlcDownloadProgress         = "SteamAPI_ISteamApps_GetDlcDownloadProgress"
	flatAPI_ISteamApps_GetAppBuildId                  = "SteamAPI_ISteamApps_GetAppBuildId"
	flatAPI_ISteamApps_RequestAllProofOfPurchaseKeys  = "SteamAPI_ISteamApps_RequestAllProofOfPurchaseKeys"
	flatAPI_ISteamApps_GetFileDetails                 = "SteamAPI_ISteamApps_GetFileDetails"
	flatAPI_ISteamApps_GetLaunchCommandLine           = "SteamAPI_ISteamApps_GetLaunchCommandLine"
	flatAPI_ISteamApps_BIsSubscribedFromFamilySharing = "SteamAPI_ISteamApps_BIsSubscribedFromFamilySharing"
	flatAPI_ISteamApps_BIsTimedTrial                  = "SteamAPI_ISteamApps_BIsTimedTrial"
	flatAPI_ISteamApps_SetDlcContext                  = "SteamAPI_ISteamApps_SetDlcContext"
	flatAPI_ISteamApps_GetNumBetas                    = "SteamAPI_ISteamApps_GetNumBetas"
	flatAPI_ISteamApps_GetBetaInfo                    = "SteamAPI_ISteamApps_GetBetaInfo"
	flatAPI_ISteamApps_SetActiveBeta                  = "SteamAPI_ISteamApps_SetActiveBeta"

	flatAPI_SteamFriends                                               = "SteamAPI_SteamFriends_v018"
	flatAPI_ISteamFriends_GetPersonaName                               = "SteamAPI_ISteamFriends_GetPersonaName"
	flatAPI_ISteamFriends_GetPersonaState                              = "SteamAPI_ISteamFriends_GetPersonaState"
	flatAPI_ISteamFriends_GetFriendCount                               = "SteamAPI_ISteamFriends_GetFriendCount"
	flatAPI_ISteamFriends_GetFriendByIndex                             = "SteamAPI_ISteamFriends_GetFriendByIndex"
	flatAPI_ISteamFriends_GetFriendRelationship                        = "SteamAPI_ISteamFriends_GetFriendRelationship"
	flatAPI_ISteamFriends_GetFriendPersonaState                        = "SteamAPI_ISteamFriends_GetFriendPersonaState"
	flatAPI_ISteamFriends_GetFriendPersonaName                         = "SteamAPI_ISteamFriends_GetFriendPersonaName"
	flatAPI_ISteamFriends_GetFriendPersonaNameHistory                  = "SteamAPI_ISteamFriends_GetFriendPersonaNameHistory"
	flatAPI_ISteamFriends_GetFriendSteamLevel                          = "SteamAPI_ISteamFriends_GetFriendSteamLevel"
	flatAPI_ISteamFriends_GetSmallFriendAvatar                         = "SteamAPI_ISteamFriends_GetSmallFriendAvatar"
	flatAPI_ISteamFriends_GetMediumFriendAvatar                        = "SteamAPI_ISteamFriends_GetMediumFriendAvatar"
	flatAPI_ISteamFriends_GetLargeFriendAvatar                         = "SteamAPI_ISteamFriends_GetLargeFriendAvatar"
	flatAPI_ISteamFriends_SetRichPresence                              = "SteamAPI_ISteamFriends_SetRichPresence"
	flatAPI_ISteamFriends_GetFriendGamePlayed                          = "SteamAPI_ISteamFriends_GetFriendGamePlayed"
	flatAPI_ISteamFriends_InviteUserToGame                             = "SteamAPI_ISteamFriends_InviteUserToGame"
	flatAPI_ISteamFriends_ActivateGameOverlay                          = "SteamAPI_ISteamFriends_ActivateGameOverlay"
	flatAPI_ISteamFriends_ActivateGameOverlayToUser                    = "SteamAPI_ISteamFriends_ActivateGameOverlayToUser"
	flatAPI_ISteamFriends_ActivateGameOverlayToWebPage                 = "SteamAPI_ISteamFriends_ActivateGameOverlayToWebPage"
	flatAPI_ISteamFriends_ActivateGameOverlayToStore                   = "SteamAPI_ISteamFriends_ActivateGameOverlayToStore"
	flatAPI_ISteamFriends_ActivateGameOverlayInviteDialog              = "SteamAPI_ISteamFriends_ActivateGameOverlayInviteDialog"
	flatAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString = "SteamAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString"

	flatAPI_SteamMatchmaking                                             = "SteamAPI_SteamMatchmaking_v009"
	flatAPI_ISteamMatchmaking_GetFavoriteGameCount                       = "SteamAPI_ISteamMatchmaking_GetFavoriteGameCount"
	flatAPI_ISteamMatchmaking_GetFavoriteGame                            = "SteamAPI_ISteamMatchmaking_GetFavoriteGame"
	flatAPI_ISteamMatchmaking_AddFavoriteGame                            = "SteamAPI_ISteamMatchmaking_AddFavoriteGame"
	flatAPI_ISteamMatchmaking_RemoveFavoriteGame                         = "SteamAPI_ISteamMatchmaking_RemoveFavoriteGame"
	flatAPI_ISteamMatchmaking_RequestLobbyList                           = "SteamAPI_ISteamMatchmaking_RequestLobbyList"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter            = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter         = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter         = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable    = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter          = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter       = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter"
	flatAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter = "SteamAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter"
	flatAPI_ISteamMatchmaking_GetLobbyByIndex                            = "SteamAPI_ISteamMatchmaking_GetLobbyByIndex"
	flatAPI_ISteamMatchmaking_CreateLobby                                = "SteamAPI_ISteamMatchmaking_CreateLobby"
	flatAPI_ISteamMatchmaking_JoinLobby                                  = "SteamAPI_ISteamMatchmaking_JoinLobby"
	flatAPI_ISteamMatchmaking_LeaveLobby                                 = "SteamAPI_ISteamMatchmaking_LeaveLobby"
	flatAPI_ISteamMatchmaking_InviteUserToLobby                          = "SteamAPI_ISteamMatchmaking_InviteUserToLobby"
	flatAPI_ISteamMatchmaking_GetLobbyMemberLimit                        = "SteamAPI_ISteamMatchmaking_GetLobbyMemberLimit"
	flatAPI_ISteamMatchmaking_SetLobbyMemberLimit                        = "SteamAPI_ISteamMatchmaking_SetLobbyMemberLimit"
	flatAPI_ISteamMatchmaking_SetLobbyType                               = "SteamAPI_ISteamMatchmaking_SetLobbyType"
	flatAPI_ISteamMatchmaking_SetLobbyJoinable                           = "SteamAPI_ISteamMatchmaking_SetLobbyJoinable"
	flatAPI_ISteamMatchmaking_GetLobbyOwner                              = "SteamAPI_ISteamMatchmaking_GetLobbyOwner"
	flatAPI_ISteamMatchmaking_SetLobbyOwner                              = "SteamAPI_ISteamMatchmaking_SetLobbyOwner"
	flatAPI_ISteamMatchmaking_SetLinkedLobby                             = "SteamAPI_ISteamMatchmaking_SetLinkedLobby"
	flatAPI_ISteamMatchmaking_GetNumLobbyMembers                         = "SteamAPI_ISteamMatchmaking_GetNumLobbyMembers"
	flatAPI_ISteamMatchmaking_GetLobbyMemberByIndex                      = "SteamAPI_ISteamMatchmaking_GetLobbyMemberByIndex"
	flatAPI_ISteamMatchmaking_SetLobbyData                               = "SteamAPI_ISteamMatchmaking_SetLobbyData"
	flatAPI_ISteamMatchmaking_GetLobbyData                               = "SteamAPI_ISteamMatchmaking_GetLobbyData"
	flatAPI_ISteamMatchmaking_DeleteLobbyData                            = "SteamAPI_ISteamMatchmaking_DeleteLobbyData"
	flatAPI_ISteamMatchmaking_GetLobbyDataCount                          = "SteamAPI_ISteamMatchmaking_GetLobbyDataCount"
	flatAPI_ISteamMatchmaking_GetLobbyDataByIndex                        = "SteamAPI_ISteamMatchmaking_GetLobbyDataByIndex"
	flatAPI_ISteamMatchmaking_SetLobbyMemberData                         = "SteamAPI_ISteamMatchmaking_SetLobbyMemberData"
	flatAPI_ISteamMatchmaking_GetLobbyMemberData                         = "SteamAPI_ISteamMatchmaking_GetLobbyMemberData"
	flatAPI_ISteamMatchmaking_SendLobbyChatMsg                           = "SteamAPI_ISteamMatchmaking_SendLobbyChatMsg"
	flatAPI_ISteamMatchmaking_GetLobbyChatEntry                          = "SteamAPI_ISteamMatchmaking_GetLobbyChatEntry"
	flatAPI_ISteamMatchmaking_RequestLobbyData                           = "SteamAPI_ISteamMatchmaking_RequestLobbyData"
	flatAPI_ISteamMatchmaking_SetLobbyGameServer                         = "SteamAPI_ISteamMatchmaking_SetLobbyGameServer"
	flatAPI_ISteamMatchmaking_GetLobbyGameServer                         = "SteamAPI_ISteamMatchmaking_GetLobbyGameServer"
	flatAPI_ISteamMatchmaking_CheckForPSNGameBootInvite                  = "SteamAPI_ISteamMatchmaking_CheckForPSNGameBootInvite"

	flatAPI_SteamMatchmakingServers_RequestInternetServerList  = "SteamAPI_ISteamMatchmakingServers_RequestInternetServerList"
	flatAPI_SteamMatchmakingServers_RequestLANServerList       = "SteamAPI_ISteamMatchmakingServers_RequestLANServerList"
	flatAPI_SteamMatchmakingServers_RequestFriendsServerList   = "SteamAPI_ISteamMatchmakingServers_RequestFriendsServerList"
	flatAPI_SteamMatchmakingServers_RequestFavoritesServerList = "SteamAPI_ISteamMatchmakingServers_RequestFavoritesServerList"
	flatAPI_SteamMatchmakingServers_RequestHistoryServerList   = "SteamAPI_ISteamMatchmakingServers_RequestHistoryServerList"
	flatAPI_SteamMatchmakingServers_RequestSpectatorServerList = "SteamAPI_ISteamMatchmakingServers_RequestSpectatorServerList"
	flatAPI_SteamMatchmakingServers_ReleaseRequest             = "SteamAPI_ISteamMatchmakingServers_ReleaseRequest"
	flatAPI_SteamMatchmakingServers_GetServerDetails           = "SteamAPI_ISteamMatchmakingServers_GetServerDetails"
	flatAPI_SteamMatchmakingServers_CancelQuery                = "SteamAPI_ISteamMatchmakingServers_CancelQuery"
	flatAPI_SteamMatchmakingServers_RefreshQuery               = "SteamAPI_ISteamMatchmakingServers_RefreshQuery"
	flatAPI_SteamMatchmakingServers_IsRefreshing               = "SteamAPI_ISteamMatchmakingServers_IsRefreshing"
	flatAPI_SteamMatchmakingServers_GetServerCount             = "SteamAPI_ISteamMatchmakingServers_GetServerCount"
	flatAPI_SteamMatchmakingServers_RefreshServer              = "SteamAPI_ISteamMatchmakingServers_RefreshServer"
	flatAPI_SteamMatchmakingServers_PingServer                 = "SteamAPI_ISteamMatchmakingServers_PingServer"
	flatAPI_SteamMatchmakingServers_PlayerDetails              = "SteamAPI_ISteamMatchmakingServers_PlayerDetails"
	flatAPI_SteamMatchmakingServers_ServerRules                = "SteamAPI_ISteamMatchmakingServers_ServerRules"
	flatAPI_SteamMatchmakingServers_CancelServerQuery          = "SteamAPI_ISteamMatchmakingServers_CancelServerQuery"

	flatAPI_SteamHTTP                            = "SteamAPI_SteamHTTP_v003"
	flatAPI_ISteamHTTP_CreateHTTPRequest         = "SteamAPI_ISteamHTTP_CreateHTTPRequest"
	flatAPI_ISteamHTTP_SetHTTPRequestHeaderValue = "SteamAPI_ISteamHTTP_SetHTTPRequestHeaderValue"
	flatAPI_ISteamHTTP_SendHTTPRequest           = "SteamAPI_ISteamHTTP_SendHTTPRequest"
	flatAPI_ISteamHTTP_GetHTTPResponseBodySize   = "SteamAPI_ISteamHTTP_GetHTTPResponseBodySize"
	flatAPI_ISteamHTTP_GetHTTPResponseBodyData   = "SteamAPI_ISteamHTTP_GetHTTPResponseBodyData"
	flatAPI_ISteamHTTP_ReleaseHTTPRequest        = "SteamAPI_ISteamHTTP_ReleaseHTTPRequest"

	flatAPI_SteamUGC                             = "SteamAPI_SteamUGC_v021"
	flatAPI_ISteamUGC_GetNumSubscribedItems      = "SteamAPI_ISteamUGC_GetNumSubscribedItems"
	flatAPI_ISteamUGC_GetSubscribedItems         = "SteamAPI_ISteamUGC_GetSubscribedItems"
	flatAPI_ISteamUGC_MarkDownloadedItemAsUnused = "SteamAPI_ISteamUGC_MarkDownloadedItemAsUnused"
	flatAPI_ISteamUGC_GetNumDownloadedItems      = "SteamAPI_ISteamUGC_GetNumDownloadedItems"
	flatAPI_ISteamUGC_GetDownloadedItems         = "SteamAPI_ISteamUGC_GetDownloadedItems"

	flatAPI_SteamInventory                  = "SteamAPI_SteamInventory_v003"
	flatAPI_ISteamInventory_GetResultStatus = "SteamAPI_ISteamInventory_GetResultStatus"
	flatAPI_ISteamInventory_GetResultItems  = "SteamAPI_ISteamInventory_GetResultItems"
	flatAPI_ISteamInventory_DestroyResult   = "SteamAPI_ISteamInventory_DestroyResult"

	flatAPI_SteamInput                               = "SteamAPI_SteamInput_v006"
	flatAPI_ISteamInput_GetConnectedControllers      = "SteamAPI_ISteamInput_GetConnectedControllers"
	flatAPI_ISteamInput_GetInputTypeForHandle        = "SteamAPI_ISteamInput_GetInputTypeForHandle"
	flatAPI_ISteamInput_Init                         = "SteamAPI_ISteamInput_Init"
	flatAPI_ISteamInput_Shutdown                     = "SteamAPI_ISteamInput_Shutdown"
	flatAPI_ISteamInput_RunFrame                     = "SteamAPI_ISteamInput_RunFrame"
	flatAPI_ISteamInput_EnableDeviceCallbacks        = "SteamAPI_ISteamInput_EnableDeviceCallbacks"
	flatAPI_ISteamInput_GetActionSetHandle           = "SteamAPI_ISteamInput_GetActionSetHandle"
	flatAPI_ISteamInput_ActivateActionSet            = "SteamAPI_ISteamInput_ActivateActionSet"
	flatAPI_ISteamInput_GetCurrentActionSet          = "SteamAPI_ISteamInput_GetCurrentActionSet"
	flatAPI_ISteamInput_ActivateActionSetLayer       = "SteamAPI_ISteamInput_ActivateActionSetLayer"
	flatAPI_ISteamInput_DeactivateActionSetLayer     = "SteamAPI_ISteamInput_DeactivateActionSetLayer"
	flatAPI_ISteamInput_DeactivateAllActionSetLayers = "SteamAPI_ISteamInput_DeactivateAllActionSetLayers"
	flatAPI_ISteamInput_GetActiveActionSetLayers     = "SteamAPI_ISteamInput_GetActiveActionSetLayers"
	flatAPI_ISteamInput_GetDigitalActionHandle       = "SteamAPI_ISteamInput_GetDigitalActionHandle"
	flatAPI_ISteamInput_GetDigitalActionData         = "SteamAPI_ISteamInput_GetDigitalActionData"
	flatAPI_ISteamInput_GetDigitalActionOrigins      = "SteamAPI_ISteamInput_GetDigitalActionOrigins"
	flatAPI_ISteamInput_GetAnalogActionHandle        = "SteamAPI_ISteamInput_GetAnalogActionHandle"
	flatAPI_ISteamInput_GetAnalogActionData          = "SteamAPI_ISteamInput_GetAnalogActionData"
	flatAPI_ISteamInput_GetAnalogActionOrigins       = "SteamAPI_ISteamInput_GetAnalogActionOrigins"
	flatAPI_ISteamInput_StopAnalogActionMomentum     = "SteamAPI_ISteamInput_StopAnalogActionMomentum"
	flatAPI_ISteamInput_GetMotionData                = "SteamAPI_ISteamInput_GetMotionData"
	flatAPI_ISteamInput_TriggerVibration             = "SteamAPI_ISteamInput_TriggerVibration"
	flatAPI_ISteamInput_TriggerVibrationExtended     = "SteamAPI_ISteamInput_TriggerVibrationExtended"
	flatAPI_ISteamInput_TriggerSimpleHapticEvent     = "SteamAPI_ISteamInput_TriggerSimpleHapticEvent"
	flatAPI_ISteamInput_SetLEDColor                  = "SteamAPI_ISteamInput_SetLEDColor"
	flatAPI_ISteamInput_ShowBindingPanel             = "SteamAPI_ISteamInput_ShowBindingPanel"
	flatAPI_ISteamInput_GetControllerForGamepadIndex = "SteamAPI_ISteamInput_GetControllerForGamepadIndex"
	flatAPI_ISteamInput_GetGamepadIndexForController = "SteamAPI_ISteamInput_GetGamepadIndexForController"
	flatAPI_ISteamInput_GetStringForActionOrigin     = "SteamAPI_ISteamInput_GetStringForActionOrigin"
	flatAPI_ISteamInput_GetGlyphForActionOrigin      = "SteamAPI_ISteamInput_GetGlyphForActionOrigin"
	flatAPI_ISteamInput_GetRemotePlaySessionID       = "SteamAPI_ISteamInput_GetRemotePlaySessionID"

	flatAPI_SteamRemotePlay                             = "SteamAPI_SteamRemotePlay_v001"
	flatAPI_ISteamRemotePlay_BSessionRemotePlayTogether = "SteamAPI_ISteamRemotePlay_BSessionRemotePlayTogether"
	flatAPI_ISteamRemotePlay_GetSessionGuestID          = "SteamAPI_ISteamRemotePlay_GetSessionGuestID"
	flatAPI_ISteamRemotePlay_GetSmallSessionAvatar      = "SteamAPI_ISteamRemotePlay_GetSmallSessionAvatar"
	flatAPI_ISteamRemotePlay_GetMediumSessionAvatar     = "SteamAPI_ISteamRemotePlay_GetMediumSessionAvatar"
	flatAPI_ISteamRemotePlay_GetLargeSessionAvatar      = "SteamAPI_ISteamRemotePlay_GetLargeSessionAvatar"

	flatAPI_SteamRemoteStorage              = "SteamAPI_SteamRemoteStorage_v016"
	flatAPI_ISteamRemoteStorage_FileWrite   = "SteamAPI_ISteamRemoteStorage_FileWrite"
	flatAPI_ISteamRemoteStorage_FileRead    = "SteamAPI_ISteamRemoteStorage_FileRead"
	flatAPI_ISteamRemoteStorage_FileDelete  = "SteamAPI_ISteamRemoteStorage_FileDelete"
	flatAPI_ISteamRemoteStorage_GetFileSize = "SteamAPI_ISteamRemoteStorage_GetFileSize"

	flatAPI_SteamUser                                 = "SteamAPI_SteamUser_v023"
	flatAPI_ISteamUser_AdvertiseGame                  = "SteamAPI_ISteamUser_AdvertiseGame"
	flatAPI_ISteamUser_BeginAuthSession               = "SteamAPI_ISteamUser_BeginAuthSession"
	flatAPI_ISteamUser_BIsBehindNAT                   = "SteamAPI_ISteamUser_BIsBehindNAT"
	flatAPI_ISteamUser_BIsPhoneIdentifying            = "SteamAPI_ISteamUser_BIsPhoneIdentifying"
	flatAPI_ISteamUser_BIsPhoneRequiringVerification  = "SteamAPI_ISteamUser_BIsPhoneRequiringVerification"
	flatAPI_ISteamUser_BIsPhoneVerified               = "SteamAPI_ISteamUser_BIsPhoneVerified"
	flatAPI_ISteamUser_BIsTwoFactorEnabled            = "SteamAPI_ISteamUser_BIsTwoFactorEnabled"
	flatAPI_ISteamUser_BLoggedOn                      = "SteamAPI_ISteamUser_BLoggedOn"
	flatAPI_ISteamUser_BSetDurationControlOnlineState = "SteamAPI_ISteamUser_BSetDurationControlOnlineState"
	flatAPI_ISteamUser_CancelAuthTicket               = "SteamAPI_ISteamUser_CancelAuthTicket"
	flatAPI_ISteamUser_DecompressVoice                = "SteamAPI_ISteamUser_DecompressVoice"
	flatAPI_ISteamUser_EndAuthSession                 = "SteamAPI_ISteamUser_EndAuthSession"
	flatAPI_ISteamUser_GetAuthSessionTicket           = "SteamAPI_ISteamUser_GetAuthSessionTicket"
	flatAPI_ISteamUser_GetAuthTicketForWebApi         = "SteamAPI_ISteamUser_GetAuthTicketForWebApi"
	flatAPI_ISteamUser_GetAvailableVoice              = "SteamAPI_ISteamUser_GetAvailableVoice"
	flatAPI_ISteamUser_GetDurationControl             = "SteamAPI_ISteamUser_GetDurationControl"
	flatAPI_ISteamUser_GetEncryptedAppTicket          = "SteamAPI_ISteamUser_GetEncryptedAppTicket"
	flatAPI_ISteamUser_GetGameBadgeLevel              = "SteamAPI_ISteamUser_GetGameBadgeLevel"
	flatAPI_ISteamUser_GetHSteamUser                  = "SteamAPI_ISteamUser_GetHSteamUser"
	flatAPI_ISteamUser_GetPlayerSteamLevel            = "SteamAPI_ISteamUser_GetPlayerSteamLevel"
	flatAPI_ISteamUser_GetSteamID                     = "SteamAPI_ISteamUser_GetSteamID"
	flatAPI_ISteamUser_GetUserDataFolder              = "SteamAPI_ISteamUser_GetUserDataFolder"
	flatAPI_ISteamUser_GetVoice                       = "SteamAPI_ISteamUser_GetVoice"
	flatAPI_ISteamUser_GetVoiceOptimalSampleRate      = "SteamAPI_ISteamUser_GetVoiceOptimalSampleRate"
	flatAPI_ISteamUser_InitiateGameConnection         = "SteamAPI_ISteamUser_InitiateGameConnection_DEPRECATED"
	flatAPI_ISteamUser_RequestEncryptedAppTicket      = "SteamAPI_ISteamUser_RequestEncryptedAppTicket"
	flatAPI_ISteamUser_RequestStoreAuthURL            = "SteamAPI_ISteamUser_RequestStoreAuthURL"
	flatAPI_ISteamUser_StartVoiceRecording            = "SteamAPI_ISteamUser_StartVoiceRecording"
	flatAPI_ISteamUser_StopVoiceRecording             = "SteamAPI_ISteamUser_StopVoiceRecording"
	flatAPI_ISteamUser_TerminateGameConnection        = "SteamAPI_ISteamUser_TerminateGameConnection_DEPRECATED"
	flatAPI_ISteamUser_TrackAppUsageEvent             = "SteamAPI_ISteamUser_TrackAppUsageEvent"
	flatAPI_ISteamUser_UserHasLicenseForApp           = "SteamAPI_ISteamUser_UserHasLicenseForApp"

	flatAPI_SteamUserStats                   = "SteamAPI_SteamUserStats_v013"
	flatAPI_ISteamUserStats_GetAchievement   = "SteamAPI_ISteamUserStats_GetAchievement"
	flatAPI_ISteamUserStats_SetAchievement   = "SteamAPI_ISteamUserStats_SetAchievement"
	flatAPI_ISteamUserStats_ClearAchievement = "SteamAPI_ISteamUserStats_ClearAchievement"
	flatAPI_ISteamUserStats_StoreStats       = "SteamAPI_ISteamUserStats_StoreStats"

	flatAPI_SteamUtils                                 = "SteamAPI_SteamUtils_v010"
	flatAPI_ISteamUtils_GetSecondsSinceAppActive       = "SteamAPI_ISteamUtils_GetSecondsSinceAppActive"
	flatAPI_ISteamUtils_GetSecondsSinceComputerActive  = "SteamAPI_ISteamUtils_GetSecondsSinceComputerActive"
	flatAPI_ISteamUtils_GetConnectedUniverse           = "SteamAPI_ISteamUtils_GetConnectedUniverse"
	flatAPI_ISteamUtils_GetServerRealTime              = "SteamAPI_ISteamUtils_GetServerRealTime"
	flatAPI_ISteamUtils_GetIPCountry                   = "SteamAPI_ISteamUtils_GetIPCountry"
	flatAPI_ISteamUtils_GetImageSize                   = "SteamAPI_ISteamUtils_GetImageSize"
	flatAPI_ISteamUtils_GetImageRGBA                   = "SteamAPI_ISteamUtils_GetImageRGBA"
	flatAPI_ISteamUtils_GetCurrentBatteryPower         = "SteamAPI_ISteamUtils_GetCurrentBatteryPower"
	flatAPI_ISteamUtils_GetAppID                       = "SteamAPI_ISteamUtils_GetAppID"
	flatAPI_ISteamUtils_SetOverlayNotificationPosition = "SteamAPI_ISteamUtils_SetOverlayNotificationPosition"
	flatAPI_ISteamUtils_IsAPICallCompleted             = "SteamAPI_ISteamUtils_IsAPICallCompleted"
	flatAPI_ISteamUtils_GetAPICallFailureReason        = "SteamAPI_ISteamUtils_GetAPICallFailureReason"
	flatAPI_ISteamUtils_GetAPICallResult               = "SteamAPI_ISteamUtils_GetAPICallResult"
	flatAPI_ISteamUtils_GetIPCCallCount                = "SteamAPI_ISteamUtils_GetIPCCallCount"
	flatAPI_ISteamUtils_IsOverlayEnabled               = "SteamAPI_ISteamUtils_IsOverlayEnabled"
	flatAPI_ISteamUtils_BOverlayNeedsPresent           = "SteamAPI_ISteamUtils_BOverlayNeedsPresent"
	flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck      = "SteamAPI_ISteamUtils_IsSteamRunningOnSteamDeck"
	flatAPI_ISteamUtils_ShowFloatingGamepadTextInput   = "SteamAPI_ISteamUtils_ShowFloatingGamepadTextInput"
	flatAPI_ISteamUtils_SetOverlayNotificationInset    = "SteamAPI_ISteamUtils_SetOverlayNotificationInset"

	flatAPI_SteamNetworkingUtils                         = "SteamAPI_SteamNetworkingUtils_SteamAPI_v004"
	flatAPI_ISteamNetworkingUtils_AllocateMessage        = "SteamAPI_ISteamNetworkingUtils_AllocateMessage"
	flatAPI_ISteamNetworkingUtils_InitRelayNetworkAccess = "SteamAPI_ISteamNetworkingUtils_InitRelayNetworkAccess"
	flatAPI_ISteamNetworkingUtils_GetLocalTimestamp      = "SteamAPI_ISteamNetworkingUtils_GetLocalTimestamp"

	flatAPI_SteamNetworkingMessages                           = "SteamAPI_SteamNetworkingMessages_SteamAPI_v002"
	flatAPI_ISteamNetworkingMessages_SendMessageToUser        = "SteamAPI_ISteamNetworkingMessages_SendMessageToUser"
	flatAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel = "SteamAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel"
	flatAPI_ISteamNetworkingMessages_AcceptSessionWithUser    = "SteamAPI_ISteamNetworkingMessages_AcceptSessionWithUser"
	flatAPI_ISteamNetworkingMessages_CloseSessionWithUser     = "SteamAPI_ISteamNetworkingMessages_CloseSessionWithUser"
	flatAPI_ISteamNetworkingMessages_CloseChannelWithUser     = "SteamAPI_ISteamNetworkingMessages_CloseChannelWithUser"

	flatAPI_SteamNetworkingSockets                              = "SteamAPI_SteamNetworkingSockets_SteamAPI_v012"
	flatAPI_ISteamNetworkingSockets_CreateListenSocketIP        = "SteamAPI_ISteamNetworkingSockets_CreateListenSocketIP"
	flatAPI_ISteamNetworkingSockets_CreateListenSocketP2P       = "SteamAPI_ISteamNetworkingSockets_CreateListenSocketP2P"
	flatAPI_ISteamNetworkingSockets_ConnectByIPAddress          = "SteamAPI_ISteamNetworkingSockets_ConnectByIPAddress"
	flatAPI_ISteamNetworkingSockets_ConnectP2P                  = "SteamAPI_ISteamNetworkingSockets_ConnectP2P"
	flatAPI_ISteamNetworkingSockets_AcceptConnection            = "SteamAPI_ISteamNetworkingSockets_AcceptConnection"
	flatAPI_ISteamNetworkingSockets_CloseConnection             = "SteamAPI_ISteamNetworkingSockets_CloseConnection"
	flatAPI_ISteamNetworkingSockets_CloseListenSocket           = "SteamAPI_ISteamNetworkingSockets_CloseListenSocket"
	flatAPI_ISteamNetworkingSockets_SendMessageToConnection     = "SteamAPI_ISteamNetworkingSockets_SendMessageToConnection"
	flatAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection = "SteamAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection"
	flatAPI_ISteamNetworkingSockets_CreatePollGroup             = "SteamAPI_ISteamNetworkingSockets_CreatePollGroup"
	flatAPI_ISteamNetworkingSockets_DestroyPollGroup            = "SteamAPI_ISteamNetworkingSockets_DestroyPollGroup"
	flatAPI_ISteamNetworkingSockets_SetConnectionPollGroup      = "SteamAPI_ISteamNetworkingSockets_SetConnectionPollGroup"
	flatAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup  = "SteamAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup"

	flatAPI_SteamGameServer                                      = "SteamAPI_SteamGameServer_v015"
	flatAPI_ISteamGameServer_AssociateWithClan                   = "SteamAPI_ISteamGameServer_AssociateWithClan"
	flatAPI_ISteamGameServer_BeginAuthSession                    = "SteamAPI_ISteamGameServer_BeginAuthSession"
	flatAPI_ISteamGameServer_BLoggedOn                           = "SteamAPI_ISteamGameServer_BLoggedOn"
	flatAPI_ISteamGameServer_BSecure                             = "SteamAPI_ISteamGameServer_BSecure"
	flatAPI_ISteamGameServer_BUpdateUserData                     = "SteamAPI_ISteamGameServer_BUpdateUserData"
	flatAPI_ISteamGameServer_CancelAuthTicket                    = "SteamAPI_ISteamGameServer_CancelAuthTicket"
	flatAPI_ISteamGameServer_ClearAllKeyValues                   = "SteamAPI_ISteamGameServer_ClearAllKeyValues"
	flatAPI_ISteamGameServer_ComputeNewPlayerCompatibility       = "SteamAPI_ISteamGameServer_ComputeNewPlayerCompatibility"
	flatAPI_ISteamGameServer_CreateUnauthenticatedUserConnection = "SteamAPI_ISteamGameServer_CreateUnauthenticatedUserConnection"
	flatAPI_ISteamGameServer_EnableHeartbeats                    = "SteamAPI_ISteamGameServer_EnableHeartbeats"
	flatAPI_ISteamGameServer_EndAuthSession                      = "SteamAPI_ISteamGameServer_EndAuthSession"
	flatAPI_ISteamGameServer_ForceHeartbeat                      = "SteamAPI_ISteamGameServer_ForceHeartbeat"
	flatAPI_ISteamGameServer_GetAuthSessionTicket                = "SteamAPI_ISteamGameServer_GetAuthSessionTicket"
	flatAPI_ISteamGameServer_GetGameplayStats                    = "SteamAPI_ISteamGameServer_GetGameplayStats"
	flatAPI_ISteamGameServer_GetNextOutgoingPacket               = "SteamAPI_ISteamGameServer_GetNextOutgoingPacket"
	flatAPI_ISteamGameServer_GetPublicIP                         = "SteamAPI_ISteamGameServer_GetPublicIP"
	flatAPI_ISteamGameServer_GetServerReputation                 = "SteamAPI_ISteamGameServer_GetServerReputation"
	flatAPI_ISteamGameServer_GetSteamID                          = "SteamAPI_ISteamGameServer_GetSteamID"
	flatAPI_ISteamGameServer_HandleIncomingPacket                = "SteamAPI_ISteamGameServer_HandleIncomingPacket"
	flatAPI_ISteamGameServer_InitGameServer                      = "SteamAPI_ISteamGameServer_InitGameServer"
	flatAPI_ISteamGameServer_LogOff                              = "SteamAPI_ISteamGameServer_LogOff"
	flatAPI_ISteamGameServer_LogOn                               = "SteamAPI_ISteamGameServer_LogOn"
	flatAPI_ISteamGameServer_LogOnAnonymous                      = "SteamAPI_ISteamGameServer_LogOnAnonymous"
	flatAPI_ISteamGameServer_RequestUserGroupStatus              = "SteamAPI_ISteamGameServer_RequestUserGroupStatus"
	flatAPI_ISteamGameServer_SendUserConnectAndAuthenticate      = "SteamAPI_ISteamGameServer_SendUserConnectAndAuthenticate_DEPRECATED"
	flatAPI_ISteamGameServer_SendUserDisconnect                  = "SteamAPI_ISteamGameServer_SendUserDisconnect_DEPRECATED"
	flatAPI_ISteamGameServer_SetBotPlayerCount                   = "SteamAPI_ISteamGameServer_SetBotPlayerCount"
	flatAPI_ISteamGameServer_SetDedicatedServer                  = "SteamAPI_ISteamGameServer_SetDedicatedServer"
	flatAPI_ISteamGameServer_SetGameData                         = "SteamAPI_ISteamGameServer_SetGameData"
	flatAPI_ISteamGameServer_SetGameDescription                  = "SteamAPI_ISteamGameServer_SetGameDescription"
	flatAPI_ISteamGameServer_SetGameTags                         = "SteamAPI_ISteamGameServer_SetGameTags"
	flatAPI_ISteamGameServer_SetHeartbeatInterval                = "SteamAPI_ISteamGameServer_SetHeartbeatInterval"
	flatAPI_ISteamGameServer_SetKeyValue                         = "SteamAPI_ISteamGameServer_SetKeyValue"
	flatAPI_ISteamGameServer_SetMapName                          = "SteamAPI_ISteamGameServer_SetMapName"
	flatAPI_ISteamGameServer_SetMaxPlayerCount                   = "SteamAPI_ISteamGameServer_SetMaxPlayerCount"
	flatAPI_ISteamGameServer_SetModDir                           = "SteamAPI_ISteamGameServer_SetModDir"
	flatAPI_ISteamGameServer_SetPasswordProtected                = "SteamAPI_ISteamGameServer_SetPasswordProtected"
	flatAPI_ISteamGameServer_SetProduct                          = "SteamAPI_ISteamGameServer_SetProduct"
	flatAPI_ISteamGameServer_SetRegion                           = "SteamAPI_ISteamGameServer_SetRegion"
	flatAPI_ISteamGameServer_SetServerName                       = "SteamAPI_ISteamGameServer_SetServerName"
	flatAPI_ISteamGameServer_SetSpectatorPort                    = "SteamAPI_ISteamGameServer_SetSpectatorPort"
	flatAPI_ISteamGameServer_SetSpectatorServerName              = "SteamAPI_ISteamGameServer_SetSpectatorServerName"
	flatAPI_ISteamGameServer_UserHasLicenseForApp                = "SteamAPI_ISteamGameServer_UserHasLicenseForApp"
	flatAPI_ISteamGameServer_WasRestartRequested                 = "SteamAPI_ISteamGameServer_WasRestartRequested"
)

type steamErrMsg [1024]byte

func (s *steamErrMsg) String() string {
	for i, b := range s {
		if b == 0 {
			return string(s[:i])
		}
	}
	return ""
}

func putUint32(dst []byte, v uint32) {
	dst[0] = byte(v)
	dst[1] = byte(v >> 8)
	dst[2] = byte(v >> 16)
	dst[3] = byte(v >> 24)
}

func putUint64(dst []byte, v uint64) {
	dst[0] = byte(v)
	dst[1] = byte(v >> 8)
	dst[2] = byte(v >> 16)
	dst[3] = byte(v >> 24)
	dst[4] = byte(v >> 32)
	dst[5] = byte(v >> 40)
	dst[6] = byte(v >> 48)
	dst[7] = byte(v >> 56)
}
