// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The go-steamworks Authors

package steamworks

import (
	"bytes"
	"fmt"
	"iter"
	"os"
	"runtime"
	"sync"
	"unique"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

type lib struct {
	lib uintptr
}

var (
	// General
	ptrAPI_RestartAppIfNecessary      func(uint32) bool
	ptrAPI_InitFlat                   func(uintptr) ESteamAPIInitResult
	ptrAPI_RunCallbacks               func()
	ptrAPI_Shutdown                   func()
	ptrAPI_IsSteamRunning             func() bool
	ptrAPI_GetSteamInstallPath        func() string
	ptrAPI_ReleaseCurrentThreadMemory func()

	// ISteamApps
	ptrAPI_SteamApps                                 func() uintptr
	ptrAPI_ISteamApps_BIsSubscribed                  func(uintptr) bool
	ptrAPI_ISteamApps_BIsLowViolence                 func(uintptr) bool
	ptrAPI_ISteamApps_BIsCybercafe                   func(uintptr) bool
	ptrAPI_ISteamApps_BIsVACBanned                   func(uintptr) bool
	ptrAPI_ISteamApps_BGetDLCDataByIndex             func(uintptr, int32, uintptr, uintptr, uintptr, int32) bool
	ptrAPI_ISteamApps_BIsDlcInstalled                func(uintptr, AppId_t) bool
	ptrAPI_ISteamApps_GetAvailableGameLanguages      func(uintptr) string
	ptrAPI_ISteamApps_BIsSubscribedApp               func(uintptr, AppId_t) bool
	ptrAPI_ISteamApps_GetEarliestPurchaseUnixTime    func(uintptr, AppId_t) uint32
	ptrAPI_ISteamApps_BIsSubscribedFromFreeWeekend   func(uintptr) bool
	ptrAPI_ISteamApps_GetAppInstallDir               func(uintptr, AppId_t, uintptr, int32) int32
	ptrAPI_ISteamApps_GetCurrentGameLanguage         func(uintptr) string
	ptrAPI_ISteamApps_GetDLCCount                    func(uintptr) int32
	ptrAPI_ISteamApps_InstallDLC                     func(uintptr, AppId_t)
	ptrAPI_ISteamApps_UninstallDLC                   func(uintptr, AppId_t)
	ptrAPI_ISteamApps_RequestAppProofOfPurchaseKey   func(uintptr, AppId_t)
	ptrAPI_ISteamApps_GetCurrentBetaName             func(uintptr, uintptr, int32) bool
	ptrAPI_ISteamApps_MarkContentCorrupt             func(uintptr, bool) bool
	ptrAPI_ISteamApps_GetInstalledDepots             func(uintptr, AppId_t, uintptr, uint32) uint32
	ptrAPI_ISteamApps_BIsAppInstalled                func(uintptr, AppId_t) bool
	ptrAPI_ISteamApps_GetAppOwner                    func(uintptr) CSteamID
	ptrAPI_ISteamApps_GetLaunchQueryParam            func(uintptr, string) string
	ptrAPI_ISteamApps_GetDlcDownloadProgress         func(uintptr, AppId_t, uintptr, uintptr) bool
	ptrAPI_ISteamApps_GetAppBuildId                  func(uintptr) int32
	ptrAPI_ISteamApps_RequestAllProofOfPurchaseKeys  func(uintptr)
	ptrAPI_ISteamApps_GetFileDetails                 func(uintptr, string) SteamAPICall_t
	ptrAPI_ISteamApps_GetLaunchCommandLine           func(uintptr, uintptr, int32) int32
	ptrAPI_ISteamApps_BIsSubscribedFromFamilySharing func(uintptr) bool
	ptrAPI_ISteamApps_BIsTimedTrial                  func(uintptr, uintptr, uintptr) bool
	ptrAPI_ISteamApps_SetDlcContext                  func(uintptr, AppId_t) bool
	ptrAPI_ISteamApps_GetNumBetas                    func(uintptr, uintptr, uintptr) int32
	ptrAPI_ISteamApps_GetBetaInfo                    func(uintptr, int32, uintptr, uintptr, uintptr, uintptr, int32, uintptr, int32) bool
	ptrAPI_ISteamApps_SetActiveBeta                  func(uintptr, string) bool

	// ISteamFriends
	ptrAPI_SteamFriends                                               func() uintptr
	ptrAPI_ISteamFriends_GetPersonaName                               func(uintptr) string
	ptrAPI_ISteamFriends_GetPersonaState                              func(uintptr) int32
	ptrAPI_ISteamFriends_GetFriendCount                               func(uintptr, int32) int32
	ptrAPI_ISteamFriends_GetFriendByIndex                             func(uintptr, int32, int32) CSteamID
	ptrAPI_ISteamFriends_GetFriendRelationship                        func(uintptr, CSteamID) int32
	ptrAPI_ISteamFriends_GetFriendPersonaState                        func(uintptr, CSteamID) int32
	ptrAPI_ISteamFriends_GetFriendPersonaName                         func(uintptr, CSteamID) string
	ptrAPI_ISteamFriends_GetFriendPersonaNameHistory                  func(uintptr, CSteamID, int32) string
	ptrAPI_ISteamFriends_GetFriendSteamLevel                          func(uintptr, CSteamID) int32
	ptrAPI_ISteamFriends_GetSmallFriendAvatar                         func(uintptr, CSteamID) int32
	ptrAPI_ISteamFriends_GetMediumFriendAvatar                        func(uintptr, CSteamID) int32
	ptrAPI_ISteamFriends_GetLargeFriendAvatar                         func(uintptr, CSteamID) int32
	ptrAPI_ISteamFriends_SetRichPresence                              func(uintptr, string, string) bool
	ptrAPI_ISteamFriends_GetFriendGamePlayed                          func(uintptr, CSteamID, uintptr) bool
	ptrAPI_ISteamFriends_InviteUserToGame                             func(uintptr, CSteamID, string) bool
	ptrAPI_ISteamFriends_ActivateGameOverlay                          func(uintptr, string)
	ptrAPI_ISteamFriends_ActivateGameOverlayToUser                    func(uintptr, string, CSteamID)
	ptrAPI_ISteamFriends_ActivateGameOverlayToWebPage                 func(uintptr, string, EActivateGameOverlayToWebPageMode)
	ptrAPI_ISteamFriends_ActivateGameOverlayToStore                   func(uintptr, AppId_t, EOverlayToStoreFlag)
	ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialog              func(uintptr, CSteamID)
	ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString func(uintptr, string)

	// ISteamMatchmaking
	ptrAPI_SteamMatchmaking                                             func() uintptr
	ptrAPI_ISteamMatchmaking_GetFavoriteGameCount                       func(uintptr) int32
	ptrAPI_ISteamMatchmaking_GetFavoriteGame                            func(uintptr, int32, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) bool
	ptrAPI_ISteamMatchmaking_AddFavoriteGame                            func(uintptr, AppId_t, uint32, uint16, uint16, uint32, uint32) int32
	ptrAPI_ISteamMatchmaking_RemoveFavoriteGame                         func(uintptr, AppId_t, uint32, uint16, uint16, uint32) bool
	ptrAPI_ISteamMatchmaking_RequestLobbyList                           func(uintptr) SteamAPICall_t
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter            func(uintptr, string, string, ELobbyComparison)
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter         func(uintptr, string, int32, ELobbyComparison)
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter         func(uintptr, string, int32)
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable    func(uintptr, int32)
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter          func(uintptr, ELobbyDistanceFilter)
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter       func(uintptr, int32)
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter func(uintptr, CSteamID)
	ptrAPI_ISteamMatchmaking_GetLobbyByIndex                            func(uintptr, int32) CSteamID
	ptrAPI_ISteamMatchmaking_CreateLobby                                func(uintptr, ELobbyType, int32) SteamAPICall_t
	ptrAPI_ISteamMatchmaking_JoinLobby                                  func(uintptr, CSteamID) SteamAPICall_t
	ptrAPI_ISteamMatchmaking_LeaveLobby                                 func(uintptr, CSteamID)
	ptrAPI_ISteamMatchmaking_InviteUserToLobby                          func(uintptr, CSteamID, CSteamID) bool
	ptrAPI_ISteamMatchmaking_GetLobbyMemberLimit                        func(uintptr, CSteamID) int32
	ptrAPI_ISteamMatchmaking_SetLobbyMemberLimit                        func(uintptr, CSteamID, int32) bool
	ptrAPI_ISteamMatchmaking_SetLobbyType                               func(uintptr, CSteamID, ELobbyType) bool
	ptrAPI_ISteamMatchmaking_SetLobbyJoinable                           func(uintptr, CSteamID, bool) bool
	ptrAPI_ISteamMatchmaking_GetLobbyOwner                              func(uintptr, CSteamID) CSteamID
	ptrAPI_ISteamMatchmaking_SetLobbyOwner                              func(uintptr, CSteamID, CSteamID) bool
	ptrAPI_ISteamMatchmaking_SetLinkedLobby                             func(uintptr, CSteamID, CSteamID) bool
	ptrAPI_ISteamMatchmaking_GetNumLobbyMembers                         func(uintptr, CSteamID) int32
	ptrAPI_ISteamMatchmaking_GetLobbyMemberByIndex                      func(uintptr, CSteamID, int32) CSteamID
	ptrAPI_ISteamMatchmaking_SetLobbyData                               func(uintptr, CSteamID, string, string) bool
	ptrAPI_ISteamMatchmaking_GetLobbyData                               func(uintptr, CSteamID, string) string
	ptrAPI_ISteamMatchmaking_DeleteLobbyData                            func(uintptr, CSteamID, string) bool
	ptrAPI_ISteamMatchmaking_GetLobbyDataCount                          func(uintptr, CSteamID) int32
	ptrAPI_ISteamMatchmaking_GetLobbyDataByIndex                        func(uintptr, CSteamID, int32, uintptr, int32, uintptr, int32) bool
	ptrAPI_ISteamMatchmaking_SetLobbyMemberData                         func(uintptr, CSteamID, string, string)
	ptrAPI_ISteamMatchmaking_GetLobbyMemberData                         func(uintptr, CSteamID, CSteamID, string) string
	ptrAPI_ISteamMatchmaking_SendLobbyChatMsg                           func(uintptr, CSteamID, uintptr, int32) bool
	ptrAPI_ISteamMatchmaking_GetLobbyChatEntry                          func(uintptr, CSteamID, int32, uintptr, uintptr, int32, uintptr) int32
	ptrAPI_ISteamMatchmaking_RequestLobbyData                           func(uintptr, CSteamID) bool
	ptrAPI_ISteamMatchmaking_SetLobbyGameServer                         func(uintptr, CSteamID, uint32, uint16, CSteamID)
	ptrAPI_ISteamMatchmaking_GetLobbyGameServer                         func(uintptr, CSteamID, uintptr, uintptr, uintptr) bool
	ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite                  func(uintptr, uintptr) bool

	// ISteamMatchmakingServers
	ptrAPI_ISteamMatchmakingServers_RequestInternetServerList  func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest
	ptrAPI_ISteamMatchmakingServers_RequestLANServerList       func(uintptr, AppId_t, uintptr) HServerListRequest
	ptrAPI_ISteamMatchmakingServers_RequestFriendsServerList   func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest
	ptrAPI_ISteamMatchmakingServers_RequestFavoritesServerList func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest
	ptrAPI_ISteamMatchmakingServers_RequestHistoryServerList   func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest
	ptrAPI_ISteamMatchmakingServers_RequestSpectatorServerList func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest
	ptrAPI_ISteamMatchmakingServers_ReleaseRequest             func(uintptr, HServerListRequest)
	ptrAPI_ISteamMatchmakingServers_GetServerDetails           func(uintptr, HServerListRequest, int32) uintptr
	ptrAPI_ISteamMatchmakingServers_CancelQuery                func(uintptr, HServerListRequest)
	ptrAPI_ISteamMatchmakingServers_RefreshQuery               func(uintptr, HServerListRequest)
	ptrAPI_ISteamMatchmakingServers_IsRefreshing               func(uintptr, HServerListRequest) bool
	ptrAPI_ISteamMatchmakingServers_GetServerCount             func(uintptr, HServerListRequest) int32
	ptrAPI_ISteamMatchmakingServers_RefreshServer              func(uintptr, HServerListRequest, int32)
	ptrAPI_ISteamMatchmakingServers_PingServer                 func(uintptr, uint32, uint16, uintptr) HServerQuery
	ptrAPI_ISteamMatchmakingServers_PlayerDetails              func(uintptr, uint32, uint16, uintptr) HServerQuery
	ptrAPI_ISteamMatchmakingServers_ServerRules                func(uintptr, uint32, uint16, uintptr) HServerQuery
	ptrAPI_ISteamMatchmakingServers_CancelServerQuery          func(uintptr, HServerQuery)

	// ISteamHTTP
	ptrAPI_SteamHTTP                            func() uintptr
	ptrAPI_ISteamHTTP_CreateHTTPRequest         func(uintptr, int32, string) HTTPRequestHandle
	ptrAPI_ISteamHTTP_SetHTTPRequestHeaderValue func(uintptr, HTTPRequestHandle, string, string) bool
	ptrAPI_ISteamHTTP_SendHTTPRequest           func(uintptr, HTTPRequestHandle, uintptr) bool
	ptrAPI_ISteamHTTP_GetHTTPResponseBodySize   func(uintptr, HTTPRequestHandle, uintptr) bool
	ptrAPI_ISteamHTTP_GetHTTPResponseBodyData   func(uintptr, HTTPRequestHandle, uintptr, uint32) bool
	ptrAPI_ISteamHTTP_ReleaseHTTPRequest        func(uintptr, HTTPRequestHandle) bool

	// ISteamUGC
	ptrAPI_SteamUGC                             func() uintptr
	ptrAPI_ISteamUGC_GetNumSubscribedItems      func(uintptr, bool) uint32
	ptrAPI_ISteamUGC_GetSubscribedItems         func(uintptr, uintptr, uint32, bool) uint32
	ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused func(uintptr, PublishedFileId_t) bool
	ptrAPI_ISteamUGC_GetNumDownloadedItems      func(uintptr) uint32
	ptrAPI_ISteamUGC_GetDownloadedItems         func(uintptr, uintptr, uint32) uint32

	// ISteamInventory
	ptrAPI_SteamInventory                  func() uintptr
	ptrAPI_ISteamInventory_GetResultStatus func(uintptr, SteamInventoryResult_t) int32
	ptrAPI_ISteamInventory_GetResultItems  func(uintptr, SteamInventoryResult_t, uintptr, uintptr) bool
	ptrAPI_ISteamInventory_DestroyResult   func(uintptr, SteamInventoryResult_t)

	// ISteamInput
	ptrAPI_SteamInput                               func() uintptr
	ptrAPI_ISteamInput_GetConnectedControllers      func(uintptr, uintptr) int32
	ptrAPI_ISteamInput_GetInputTypeForHandle        func(uintptr, InputHandle_t) int32
	ptrAPI_ISteamInput_Init                         func(uintptr, bool) bool
	ptrAPI_ISteamInput_Shutdown                     func(uintptr)
	ptrAPI_ISteamInput_RunFrame                     func(uintptr, bool)
	ptrAPI_ISteamInput_EnableDeviceCallbacks        func(uintptr)
	ptrAPI_ISteamInput_GetActionSetHandle           func(uintptr, string) InputActionSetHandle_t
	ptrAPI_ISteamInput_ActivateActionSet            func(uintptr, InputHandle_t, InputActionSetHandle_t)
	ptrAPI_ISteamInput_GetCurrentActionSet          func(uintptr, InputHandle_t) InputActionSetHandle_t
	ptrAPI_ISteamInput_ActivateActionSetLayer       func(uintptr, InputHandle_t, InputActionSetHandle_t)
	ptrAPI_ISteamInput_DeactivateActionSetLayer     func(uintptr, InputHandle_t, InputActionSetHandle_t)
	ptrAPI_ISteamInput_DeactivateAllActionSetLayers func(uintptr, InputHandle_t)
	ptrAPI_ISteamInput_GetActiveActionSetLayers     func(uintptr, InputHandle_t, uintptr) int32
	ptrAPI_ISteamInput_GetDigitalActionHandle       func(uintptr, string) InputDigitalActionHandle_t
	ptrAPI_ISteamInput_GetDigitalActionData         uintptr
	ptrAPI_ISteamInput_GetDigitalActionOrigins      func(uintptr, InputHandle_t, InputActionSetHandle_t, InputDigitalActionHandle_t, uintptr) int32
	ptrAPI_ISteamInput_GetAnalogActionHandle        func(uintptr, string) InputAnalogActionHandle_t
	ptrAPI_ISteamInput_GetAnalogActionData          uintptr
	ptrAPI_ISteamInput_GetAnalogActionOrigins       func(uintptr, InputHandle_t, InputActionSetHandle_t, InputAnalogActionHandle_t, uintptr) int32
	ptrAPI_ISteamInput_StopAnalogActionMomentum     func(uintptr, InputHandle_t, InputAnalogActionHandle_t)
	ptrAPI_ISteamInput_GetMotionData                uintptr
	ptrAPI_ISteamInput_TriggerVibration             func(uintptr, InputHandle_t, uint16, uint16)
	ptrAPI_ISteamInput_TriggerVibrationExtended     func(uintptr, InputHandle_t, uint16, uint16, uint16, uint16)
	ptrAPI_ISteamInput_TriggerSimpleHapticEvent     func(uintptr, InputHandle_t, ESteamControllerPad, uint16, uint16, uint16)
	ptrAPI_ISteamInput_SetLEDColor                  func(uintptr, InputHandle_t, uint8, uint8, uint8, ESteamInputLEDFlag)
	ptrAPI_ISteamInput_ShowBindingPanel             func(uintptr, InputHandle_t) bool
	ptrAPI_ISteamInput_GetControllerForGamepadIndex func(uintptr, int32) InputHandle_t
	ptrAPI_ISteamInput_GetGamepadIndexForController func(uintptr, InputHandle_t) int32
	ptrAPI_ISteamInput_GetStringForActionOrigin     func(uintptr, EInputActionOrigin) string
	ptrAPI_ISteamInput_GetGlyphForActionOrigin      func(uintptr, EInputActionOrigin) string
	ptrAPI_ISteamInput_GetRemotePlaySessionID       func(uintptr, InputHandle_t) uint32

	// ISteamRemotePlay
	ptrAPI_SteamRemotePlay                             func() uintptr
	ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether func(uintptr, uint32) bool
	ptrAPI_ISteamRemotePlay_GetSessionGuestID          func(uintptr, uint32) CSteamID
	ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar      func(uintptr, uint32) int32
	ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar     func(uintptr, uint32) int32
	ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar      func(uintptr, uint32) int32

	// ISteamRemoteStorage
	ptrAPI_SteamRemoteStorage              func() uintptr
	ptrAPI_ISteamRemoteStorage_FileWrite   func(uintptr, string, uintptr, int32) bool
	ptrAPI_ISteamRemoteStorage_FileRead    func(uintptr, string, uintptr, int32) int32
	ptrAPI_ISteamRemoteStorage_FileDelete  func(uintptr, string) bool
	ptrAPI_ISteamRemoteStorage_GetFileSize func(uintptr, string) int32

	// ISteamUser
	ptrAPI_SteamUser                                 func() uintptr
	ptrAPI_ISteamUser_AdvertiseGame                  func(uintptr, CSteamID, uint32, uint16)
	ptrAPI_ISteamUser_BeginAuthSession               func(uintptr, uintptr, int32, CSteamID) int32
	ptrAPI_ISteamUser_BIsBehindNAT                   func(uintptr) bool
	ptrAPI_ISteamUser_BIsPhoneIdentifying            func(uintptr) bool
	ptrAPI_ISteamUser_BIsPhoneRequiringVerification  func(uintptr) bool
	ptrAPI_ISteamUser_BIsPhoneVerified               func(uintptr) bool
	ptrAPI_ISteamUser_BIsTwoFactorEnabled            func(uintptr) bool
	ptrAPI_ISteamUser_BLoggedOn                      func(uintptr) bool
	ptrAPI_ISteamUser_BSetDurationControlOnlineState func(uintptr, EDurationControlOnlineState) bool
	ptrAPI_ISteamUser_CancelAuthTicket               func(uintptr, HAuthTicket)
	ptrAPI_ISteamUser_DecompressVoice                func(uintptr, uintptr, uint32, uintptr, uint32, uintptr, uint32) int32
	ptrAPI_ISteamUser_EndAuthSession                 func(uintptr, CSteamID)
	ptrAPI_ISteamUser_GetAuthSessionTicket           func(uintptr, uintptr, int32, uintptr, uintptr) HAuthTicket
	ptrAPI_ISteamUser_GetAuthTicketForWebApi         func(uintptr, string, uintptr) HAuthTicket
	ptrAPI_ISteamUser_GetAvailableVoice              func(uintptr, uintptr, uintptr, uint32) int32
	ptrAPI_ISteamUser_GetDurationControl             func(uintptr, uintptr) bool
	ptrAPI_ISteamUser_GetEncryptedAppTicket          func(uintptr, uintptr, int32, uintptr) bool
	ptrAPI_ISteamUser_GetGameBadgeLevel              func(uintptr, int32, bool) int32
	ptrAPI_ISteamUser_GetHSteamUser                  func(uintptr) HSteamUser
	ptrAPI_ISteamUser_GetPlayerSteamLevel            func(uintptr) int32
	ptrAPI_ISteamUser_GetSteamID                     func(uintptr) CSteamID
	ptrAPI_ISteamUser_GetUserDataFolder              func(uintptr, uintptr, int32) bool
	ptrAPI_ISteamUser_GetVoice                       func(uintptr, bool, uintptr, uint32, uintptr, bool, uintptr, uint32, uintptr, uint32) int32
	ptrAPI_ISteamUser_GetVoiceOptimalSampleRate      func(uintptr) uint32
	ptrAPI_ISteamUser_InitiateGameConnection         func(uintptr, uintptr, int32, CSteamID, uint32, uint16, bool) int32
	ptrAPI_ISteamUser_RequestEncryptedAppTicket      func(uintptr, uintptr, int32) SteamAPICall_t
	ptrAPI_ISteamUser_RequestStoreAuthURL            func(uintptr, string) SteamAPICall_t
	ptrAPI_ISteamUser_StartVoiceRecording            func(uintptr)
	ptrAPI_ISteamUser_StopVoiceRecording             func(uintptr)
	ptrAPI_ISteamUser_TerminateGameConnection        func(uintptr, uint32, uint16)
	ptrAPI_ISteamUser_TrackAppUsageEvent             func(uintptr, CGameID, int32, string)
	ptrAPI_ISteamUser_UserHasLicenseForApp           func(uintptr, CSteamID, AppId_t) int32

	// ISteamUserStats
	ptrAPI_SteamUserStats                   func() uintptr
	ptrAPI_ISteamUserStats_GetAchievement   func(uintptr, string, uintptr) bool
	ptrAPI_ISteamUserStats_SetAchievement   func(uintptr, string) bool
	ptrAPI_ISteamUserStats_ClearAchievement func(uintptr, string) bool
	ptrAPI_ISteamUserStats_StoreStats       func(uintptr) bool

	// ISteamUtils
	ptrAPI_SteamUtils                                 func() uintptr
	ptrAPI_ISteamUtils_GetSecondsSinceAppActive       func(uintptr) uint32
	ptrAPI_ISteamUtils_GetSecondsSinceComputerActive  func(uintptr) uint32
	ptrAPI_ISteamUtils_GetConnectedUniverse           func(uintptr) int32
	ptrAPI_ISteamUtils_GetServerRealTime              func(uintptr) uint32
	ptrAPI_ISteamUtils_GetIPCountry                   func(uintptr) string
	ptrAPI_ISteamUtils_GetImageSize                   func(uintptr, int32, uintptr, uintptr) bool
	ptrAPI_ISteamUtils_GetImageRGBA                   func(uintptr, int32, uintptr, int32) bool
	ptrAPI_ISteamUtils_GetCurrentBatteryPower         func(uintptr) uint8
	ptrAPI_ISteamUtils_GetAppID                       func(uintptr) uint32
	ptrAPI_ISteamUtils_SetOverlayNotificationPosition func(uintptr, ENotificationPosition)
	ptrAPI_ISteamUtils_IsAPICallCompleted             func(uintptr, SteamAPICall_t, uintptr) bool
	ptrAPI_ISteamUtils_GetAPICallFailureReason        func(uintptr, SteamAPICall_t) int32
	ptrAPI_ISteamUtils_GetAPICallResult               func(uintptr, SteamAPICall_t, uintptr, int32, int32, uintptr) bool
	ptrAPI_ISteamUtils_GetIPCCallCount                func(uintptr) uint32
	ptrAPI_ISteamUtils_IsOverlayEnabled               func(uintptr) bool
	ptrAPI_ISteamUtils_BOverlayNeedsPresent           func(uintptr) bool
	ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck      func(uintptr) bool
	ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput   func(uintptr, EFloatingGamepadTextInputMode, int32, int32, int32, int32) bool
	ptrAPI_ISteamUtils_SetOverlayNotificationInset    func(uintptr, int32, int32)

	// ISteamNetworkingUtils
	ptrAPI_SteamNetworkingUtils                         func() uintptr
	ptrAPI_ISteamNetworkingUtils_AllocateMessage        func(uintptr, int32) uintptr
	ptrAPI_ISteamNetworkingUtils_InitRelayNetworkAccess func(uintptr)
	ptrAPI_ISteamNetworkingUtils_GetLocalTimestamp      func(uintptr) SteamNetworkingMicroseconds

	// ISteamGameServer
	ptrAPI_SteamGameServer                                      func() uintptr
	ptrAPI_ISteamGameServer_AssociateWithClan                   func(uintptr, CSteamID) SteamAPICall_t
	ptrAPI_ISteamGameServer_BeginAuthSession                    func(uintptr, uintptr, int32, CSteamID) int32
	ptrAPI_ISteamGameServer_BLoggedOn                           func(uintptr) bool
	ptrAPI_ISteamGameServer_BSecure                             func(uintptr) bool
	ptrAPI_ISteamGameServer_BUpdateUserData                     func(uintptr, CSteamID, string, uint32) bool
	ptrAPI_ISteamGameServer_CancelAuthTicket                    func(uintptr, HAuthTicket)
	ptrAPI_ISteamGameServer_ClearAllKeyValues                   func(uintptr)
	ptrAPI_ISteamGameServer_ComputeNewPlayerCompatibility       func(uintptr, CSteamID, uintptr, uint32, uintptr, uint32, uintptr, uint32) SteamAPICall_t
	ptrAPI_ISteamGameServer_CreateUnauthenticatedUserConnection func(uintptr) CSteamID
	ptrAPI_ISteamGameServer_EnableHeartbeats                    func(uintptr, bool)
	ptrAPI_ISteamGameServer_EndAuthSession                      func(uintptr, CSteamID)
	ptrAPI_ISteamGameServer_ForceHeartbeat                      func(uintptr)
	ptrAPI_ISteamGameServer_GetAuthSessionTicket                func(uintptr, uintptr, int32, uintptr) HAuthTicket
	ptrAPI_ISteamGameServer_GetGameplayStats                    func(uintptr)
	ptrAPI_ISteamGameServer_GetNextOutgoingPacket               func(uintptr, uintptr, int32, uintptr, uintptr) int32
	ptrAPI_ISteamGameServer_GetPublicIP                         func(uintptr) uint32
	ptrAPI_ISteamGameServer_GetServerReputation                 func(uintptr) SteamAPICall_t
	ptrAPI_ISteamGameServer_GetSteamID                          func(uintptr) CSteamID
	ptrAPI_ISteamGameServer_HandleIncomingPacket                func(uintptr, uintptr, int32, uint32, uint16) bool
	ptrAPI_ISteamGameServer_InitGameServer                      func(uintptr, uint32, uint16, uint16, uint16, uint32, string) bool
	ptrAPI_ISteamGameServer_LogOff                              func(uintptr)
	ptrAPI_ISteamGameServer_LogOn                               func(uintptr, string)
	ptrAPI_ISteamGameServer_LogOnAnonymous                      func(uintptr)
	ptrAPI_ISteamGameServer_RequestUserGroupStatus              func(uintptr, CSteamID, CSteamID) bool
	ptrAPI_ISteamGameServer_SendUserConnectAndAuthenticate      func(uintptr, uint32, uintptr, uint32, uintptr) bool
	ptrAPI_ISteamGameServer_SendUserDisconnect                  func(uintptr, CSteamID)
	ptrAPI_ISteamGameServer_SetBotPlayerCount                   func(uintptr, int32)
	ptrAPI_ISteamGameServer_SetDedicatedServer                  func(uintptr, bool)
	ptrAPI_ISteamGameServer_SetGameData                         func(uintptr, string)
	ptrAPI_ISteamGameServer_SetGameDescription                  func(uintptr, string)
	ptrAPI_ISteamGameServer_SetGameTags                         func(uintptr, string)
	ptrAPI_ISteamGameServer_SetHeartbeatInterval                func(uintptr, int32)
	ptrAPI_ISteamGameServer_SetKeyValue                         func(uintptr, string, string)
	ptrAPI_ISteamGameServer_SetMapName                          func(uintptr, string)
	ptrAPI_ISteamGameServer_SetMaxPlayerCount                   func(uintptr, int32)
	ptrAPI_ISteamGameServer_SetModDir                           func(uintptr, string)
	ptrAPI_ISteamGameServer_SetPasswordProtected                func(uintptr, bool)
	ptrAPI_ISteamGameServer_SetProduct                          func(uintptr, string)
	ptrAPI_ISteamGameServer_SetRegion                           func(uintptr, string)
	ptrAPI_ISteamGameServer_SetServerName                       func(uintptr, string)
	ptrAPI_ISteamGameServer_SetSpectatorPort                    func(uintptr, uint16)
	ptrAPI_ISteamGameServer_SetSpectatorServerName              func(uintptr, string)
	ptrAPI_ISteamGameServer_UserHasLicenseForApp                func(uintptr, CSteamID, AppId_t) int32
	ptrAPI_ISteamGameServer_WasRestartRequested                 func(uintptr) bool

	// ISteamNetworkingMessages
	ptrAPI_SteamNetworkingMessages                           func() uintptr
	ptrAPI_ISteamNetworkingMessages_SendMessageToUser        func(uintptr, uintptr, uintptr, uint32, int32, int32) EResult
	ptrAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel func(uintptr, int32, uintptr, int32) int32
	ptrAPI_ISteamNetworkingMessages_AcceptSessionWithUser    func(uintptr, uintptr) bool
	ptrAPI_ISteamNetworkingMessages_CloseSessionWithUser     func(uintptr, uintptr) bool
	ptrAPI_ISteamNetworkingMessages_CloseChannelWithUser     func(uintptr, uintptr, int32) bool

	// ISteamNetworkingSockets
	ptrAPI_SteamNetworkingSockets                              func() uintptr
	ptrAPI_ISteamNetworkingSockets_CreateListenSocketIP        func(uintptr, uintptr, int32, uintptr) HSteamListenSocket
	ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P       func(uintptr, int32, int32, uintptr) HSteamListenSocket
	ptrAPI_ISteamNetworkingSockets_ConnectByIPAddress          func(uintptr, uintptr, int32, uintptr) HSteamNetConnection
	ptrAPI_ISteamNetworkingSockets_ConnectP2P                  func(uintptr, uintptr, int32, int32, uintptr) HSteamNetConnection
	ptrAPI_ISteamNetworkingSockets_AcceptConnection            func(uintptr, HSteamNetConnection) EResult
	ptrAPI_ISteamNetworkingSockets_CloseConnection             func(uintptr, HSteamNetConnection, int32, string, bool) bool
	ptrAPI_ISteamNetworkingSockets_CloseListenSocket           func(uintptr, HSteamListenSocket) bool
	ptrAPI_ISteamNetworkingSockets_SendMessageToConnection     func(uintptr, HSteamNetConnection, uintptr, uint32, int32, uintptr) EResult
	ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection func(uintptr, HSteamNetConnection, uintptr, int32) int32
	ptrAPI_ISteamNetworkingSockets_CreatePollGroup             func(uintptr) HSteamNetPollGroup
	ptrAPI_ISteamNetworkingSockets_DestroyPollGroup            func(uintptr, HSteamNetPollGroup) bool
	ptrAPI_ISteamNetworkingSockets_SetConnectionPollGroup      func(uintptr, HSteamNetConnection, HSteamNetPollGroup) bool
	ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup  func(uintptr, HSteamNetPollGroup, uintptr, int32) int32
)

var ffiLibOnce = sync.OnceValues(func() (ffi.Lib, error) {
	candidates := []string{}
	if customPath := os.Getenv("STEAMWORKS_LIB_PATH"); customPath != "" {
		candidates = append(candidates, customPath)
	}
	switch runtime.GOOS {
	case "windows":
		candidates = append(candidates, "steam_api64.dll")
	case "darwin":
		candidates = append(candidates, "libsteam_api.dylib")
	default:
		candidates = append(candidates, "libsteam_api.so")
	}

	for _, candidate := range candidates {
		lib, err := ffi.Load(candidate)
		if err == nil {
			return lib, nil
		}
	}
	return ffi.Lib{}, fmt.Errorf("steamworks: ffi loader could not open steam_api library")
})

func fallbackResolveInterfaceSymbol(symbol string) uintptr {
	lib, err := ffiLibOnce()
	if err != nil {
		return 0
	}
	ptr, err := lib.Get(symbol)
	if err != nil {
		return 0
	}
	return ptr
}

func mustLookupSymbol(lib uintptr, name string) uintptr {
	ptr, err := lookupSymbolAddr(lib, name)
	if err != nil {
		panic(fmt.Errorf("steamworks: symbol lookup failed for %s: %w", name, err))
	}
	return ptr
}

func registerInputStructReturns(lib uintptr) {
	ptrAPI_ISteamInput_GetDigitalActionData = mustLookupSymbol(lib, flatAPI_ISteamInput_GetDigitalActionData)
	ptrAPI_ISteamInput_GetAnalogActionData = mustLookupSymbol(lib, flatAPI_ISteamInput_GetAnalogActionData)
	ptrAPI_ISteamInput_GetMotionData = mustLookupSymbol(lib, flatAPI_ISteamInput_GetMotionData)
}

func registerOptionalFunc(fptr any, lib uintptr, name string) {
	ptr, err := lookupSymbolAddr(lib, name)
	if err != nil {
		return
	}
	purego.RegisterFunc(fptr, ptr)
}

func registerFunctions(lib uintptr) {
	// General
	purego.RegisterLibFunc(&ptrAPI_RestartAppIfNecessary, lib, flatAPI_RestartAppIfNecessary)
	purego.RegisterLibFunc(&ptrAPI_InitFlat, lib, flatAPI_InitFlat)
	purego.RegisterLibFunc(&ptrAPI_RunCallbacks, lib, flatAPI_RunCallbacks)
	purego.RegisterLibFunc(&ptrAPI_Shutdown, lib, flatAPI_Shutdown)
	purego.RegisterLibFunc(&ptrAPI_IsSteamRunning, lib, flatAPI_IsSteamRunning)
	purego.RegisterLibFunc(&ptrAPI_GetSteamInstallPath, lib, flatAPI_GetSteamInstallPath)
	purego.RegisterLibFunc(&ptrAPI_ReleaseCurrentThreadMemory, lib, flatAPI_ReleaseCurrentThreadMemory)

	// ISteamApps
	purego.RegisterLibFunc(&ptrAPI_SteamApps, lib, flatAPI_SteamApps)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsSubscribed, lib, flatAPI_ISteamApps_BIsSubscribed)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsLowViolence, lib, flatAPI_ISteamApps_BIsLowViolence)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsCybercafe, lib, flatAPI_ISteamApps_BIsCybercafe)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsVACBanned, lib, flatAPI_ISteamApps_BIsVACBanned)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BGetDLCDataByIndex, lib, flatAPI_ISteamApps_BGetDLCDataByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsDlcInstalled, lib, flatAPI_ISteamApps_BIsDlcInstalled)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetAvailableGameLanguages, lib, flatAPI_ISteamApps_GetAvailableGameLanguages)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsSubscribedApp, lib, flatAPI_ISteamApps_BIsSubscribedApp)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetEarliestPurchaseUnixTime, lib, flatAPI_ISteamApps_GetEarliestPurchaseUnixTime)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsSubscribedFromFreeWeekend, lib, flatAPI_ISteamApps_BIsSubscribedFromFreeWeekend)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetAppInstallDir, lib, flatAPI_ISteamApps_GetAppInstallDir)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetCurrentGameLanguage, lib, flatAPI_ISteamApps_GetCurrentGameLanguage)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetDLCCount, lib, flatAPI_ISteamApps_GetDLCCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_InstallDLC, lib, flatAPI_ISteamApps_InstallDLC)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_UninstallDLC, lib, flatAPI_ISteamApps_UninstallDLC)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_RequestAppProofOfPurchaseKey, lib, flatAPI_ISteamApps_RequestAppProofOfPurchaseKey)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetCurrentBetaName, lib, flatAPI_ISteamApps_GetCurrentBetaName)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_MarkContentCorrupt, lib, flatAPI_ISteamApps_MarkContentCorrupt)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetInstalledDepots, lib, flatAPI_ISteamApps_GetInstalledDepots)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsAppInstalled, lib, flatAPI_ISteamApps_BIsAppInstalled)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetAppOwner, lib, flatAPI_ISteamApps_GetAppOwner)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetLaunchQueryParam, lib, flatAPI_ISteamApps_GetLaunchQueryParam)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetDlcDownloadProgress, lib, flatAPI_ISteamApps_GetDlcDownloadProgress)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetAppBuildId, lib, flatAPI_ISteamApps_GetAppBuildId)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_RequestAllProofOfPurchaseKeys, lib, flatAPI_ISteamApps_RequestAllProofOfPurchaseKeys)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetFileDetails, lib, flatAPI_ISteamApps_GetFileDetails)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetLaunchCommandLine, lib, flatAPI_ISteamApps_GetLaunchCommandLine)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsSubscribedFromFamilySharing, lib, flatAPI_ISteamApps_BIsSubscribedFromFamilySharing)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsTimedTrial, lib, flatAPI_ISteamApps_BIsTimedTrial)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_SetDlcContext, lib, flatAPI_ISteamApps_SetDlcContext)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetNumBetas, lib, flatAPI_ISteamApps_GetNumBetas)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetBetaInfo, lib, flatAPI_ISteamApps_GetBetaInfo)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_SetActiveBeta, lib, flatAPI_ISteamApps_SetActiveBeta)

	// ISteamFriends
	purego.RegisterLibFunc(&ptrAPI_SteamFriends, lib, flatAPI_SteamFriends)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetPersonaName, lib, flatAPI_ISteamFriends_GetPersonaName)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetPersonaState, lib, flatAPI_ISteamFriends_GetPersonaState)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendCount, lib, flatAPI_ISteamFriends_GetFriendCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendByIndex, lib, flatAPI_ISteamFriends_GetFriendByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendRelationship, lib, flatAPI_ISteamFriends_GetFriendRelationship)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendPersonaState, lib, flatAPI_ISteamFriends_GetFriendPersonaState)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendPersonaName, lib, flatAPI_ISteamFriends_GetFriendPersonaName)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendPersonaNameHistory, lib, flatAPI_ISteamFriends_GetFriendPersonaNameHistory)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendSteamLevel, lib, flatAPI_ISteamFriends_GetFriendSteamLevel)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetSmallFriendAvatar, lib, flatAPI_ISteamFriends_GetSmallFriendAvatar)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetMediumFriendAvatar, lib, flatAPI_ISteamFriends_GetMediumFriendAvatar)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetLargeFriendAvatar, lib, flatAPI_ISteamFriends_GetLargeFriendAvatar)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_SetRichPresence, lib, flatAPI_ISteamFriends_SetRichPresence)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetFriendGamePlayed, lib, flatAPI_ISteamFriends_GetFriendGamePlayed)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_InviteUserToGame, lib, flatAPI_ISteamFriends_InviteUserToGame)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_ActivateGameOverlay, lib, flatAPI_ISteamFriends_ActivateGameOverlay)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_ActivateGameOverlayToUser, lib, flatAPI_ISteamFriends_ActivateGameOverlayToUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_ActivateGameOverlayToWebPage, lib, flatAPI_ISteamFriends_ActivateGameOverlayToWebPage)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_ActivateGameOverlayToStore, lib, flatAPI_ISteamFriends_ActivateGameOverlayToStore)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialog, lib, flatAPI_ISteamFriends_ActivateGameOverlayInviteDialog)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString, lib, flatAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString)

	// ISteamMatchmaking
	purego.RegisterLibFunc(&ptrAPI_SteamMatchmaking, lib, flatAPI_SteamMatchmaking)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetFavoriteGameCount, lib, flatAPI_ISteamMatchmaking_GetFavoriteGameCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetFavoriteGame, lib, flatAPI_ISteamMatchmaking_GetFavoriteGame)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddFavoriteGame, lib, flatAPI_ISteamMatchmaking_AddFavoriteGame)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_RemoveFavoriteGame, lib, flatAPI_ISteamMatchmaking_RemoveFavoriteGame)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_RequestLobbyList, lib, flatAPI_ISteamMatchmaking_RequestLobbyList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter, lib, flatAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyByIndex, lib, flatAPI_ISteamMatchmaking_GetLobbyByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_CreateLobby, lib, flatAPI_ISteamMatchmaking_CreateLobby)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_JoinLobby, lib, flatAPI_ISteamMatchmaking_JoinLobby)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_LeaveLobby, lib, flatAPI_ISteamMatchmaking_LeaveLobby)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_InviteUserToLobby, lib, flatAPI_ISteamMatchmaking_InviteUserToLobby)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyMemberLimit, lib, flatAPI_ISteamMatchmaking_GetLobbyMemberLimit)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyMemberLimit, lib, flatAPI_ISteamMatchmaking_SetLobbyMemberLimit)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyType, lib, flatAPI_ISteamMatchmaking_SetLobbyType)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyJoinable, lib, flatAPI_ISteamMatchmaking_SetLobbyJoinable)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyOwner, lib, flatAPI_ISteamMatchmaking_GetLobbyOwner)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyOwner, lib, flatAPI_ISteamMatchmaking_SetLobbyOwner)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLinkedLobby, lib, flatAPI_ISteamMatchmaking_SetLinkedLobby)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetNumLobbyMembers, lib, flatAPI_ISteamMatchmaking_GetNumLobbyMembers)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyMemberByIndex, lib, flatAPI_ISteamMatchmaking_GetLobbyMemberByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyData, lib, flatAPI_ISteamMatchmaking_SetLobbyData)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyData, lib, flatAPI_ISteamMatchmaking_GetLobbyData)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_DeleteLobbyData, lib, flatAPI_ISteamMatchmaking_DeleteLobbyData)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyDataCount, lib, flatAPI_ISteamMatchmaking_GetLobbyDataCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyDataByIndex, lib, flatAPI_ISteamMatchmaking_GetLobbyDataByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyMemberData, lib, flatAPI_ISteamMatchmaking_SetLobbyMemberData)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyMemberData, lib, flatAPI_ISteamMatchmaking_GetLobbyMemberData)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SendLobbyChatMsg, lib, flatAPI_ISteamMatchmaking_SendLobbyChatMsg)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyChatEntry, lib, flatAPI_ISteamMatchmaking_GetLobbyChatEntry)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_RequestLobbyData, lib, flatAPI_ISteamMatchmaking_RequestLobbyData)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_SetLobbyGameServer, lib, flatAPI_ISteamMatchmaking_SetLobbyGameServer)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmaking_GetLobbyGameServer, lib, flatAPI_ISteamMatchmaking_GetLobbyGameServer)
	registerOptionalFunc(&ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite, lib, flatAPI_ISteamMatchmaking_CheckForPSNGameBootInvite)

	// ISteamMatchmakingServers
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RequestInternetServerList, lib, flatAPI_SteamMatchmakingServers_RequestInternetServerList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RequestLANServerList, lib, flatAPI_SteamMatchmakingServers_RequestLANServerList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RequestFriendsServerList, lib, flatAPI_SteamMatchmakingServers_RequestFriendsServerList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RequestFavoritesServerList, lib, flatAPI_SteamMatchmakingServers_RequestFavoritesServerList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RequestHistoryServerList, lib, flatAPI_SteamMatchmakingServers_RequestHistoryServerList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RequestSpectatorServerList, lib, flatAPI_SteamMatchmakingServers_RequestSpectatorServerList)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_ReleaseRequest, lib, flatAPI_SteamMatchmakingServers_ReleaseRequest)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_GetServerDetails, lib, flatAPI_SteamMatchmakingServers_GetServerDetails)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_CancelQuery, lib, flatAPI_SteamMatchmakingServers_CancelQuery)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RefreshQuery, lib, flatAPI_SteamMatchmakingServers_RefreshQuery)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_IsRefreshing, lib, flatAPI_SteamMatchmakingServers_IsRefreshing)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_GetServerCount, lib, flatAPI_SteamMatchmakingServers_GetServerCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_RefreshServer, lib, flatAPI_SteamMatchmakingServers_RefreshServer)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_PingServer, lib, flatAPI_SteamMatchmakingServers_PingServer)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_PlayerDetails, lib, flatAPI_SteamMatchmakingServers_PlayerDetails)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_ServerRules, lib, flatAPI_SteamMatchmakingServers_ServerRules)
	purego.RegisterLibFunc(&ptrAPI_ISteamMatchmakingServers_CancelServerQuery, lib, flatAPI_SteamMatchmakingServers_CancelServerQuery)

	// ISteamHTTP
	purego.RegisterLibFunc(&ptrAPI_SteamHTTP, lib, flatAPI_SteamHTTP)
	purego.RegisterLibFunc(&ptrAPI_ISteamHTTP_CreateHTTPRequest, lib, flatAPI_ISteamHTTP_CreateHTTPRequest)
	purego.RegisterLibFunc(&ptrAPI_ISteamHTTP_SetHTTPRequestHeaderValue, lib, flatAPI_ISteamHTTP_SetHTTPRequestHeaderValue)
	purego.RegisterLibFunc(&ptrAPI_ISteamHTTP_SendHTTPRequest, lib, flatAPI_ISteamHTTP_SendHTTPRequest)
	purego.RegisterLibFunc(&ptrAPI_ISteamHTTP_GetHTTPResponseBodySize, lib, flatAPI_ISteamHTTP_GetHTTPResponseBodySize)
	purego.RegisterLibFunc(&ptrAPI_ISteamHTTP_GetHTTPResponseBodyData, lib, flatAPI_ISteamHTTP_GetHTTPResponseBodyData)
	purego.RegisterLibFunc(&ptrAPI_ISteamHTTP_ReleaseHTTPRequest, lib, flatAPI_ISteamHTTP_ReleaseHTTPRequest)

	// ISteamUGC
	purego.RegisterLibFunc(&ptrAPI_SteamUGC, lib, flatAPI_SteamUGC)
	purego.RegisterLibFunc(&ptrAPI_ISteamUGC_GetNumSubscribedItems, lib, flatAPI_ISteamUGC_GetNumSubscribedItems)
	purego.RegisterLibFunc(&ptrAPI_ISteamUGC_GetSubscribedItems, lib, flatAPI_ISteamUGC_GetSubscribedItems)
	registerOptionalFunc(&ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused, lib, flatAPI_ISteamUGC_MarkDownloadedItemAsUnused)
	registerOptionalFunc(&ptrAPI_ISteamUGC_GetNumDownloadedItems, lib, flatAPI_ISteamUGC_GetNumDownloadedItems)
	registerOptionalFunc(&ptrAPI_ISteamUGC_GetDownloadedItems, lib, flatAPI_ISteamUGC_GetDownloadedItems)

	// ISteamInventory
	purego.RegisterLibFunc(&ptrAPI_SteamInventory, lib, flatAPI_SteamInventory)
	purego.RegisterLibFunc(&ptrAPI_ISteamInventory_GetResultStatus, lib, flatAPI_ISteamInventory_GetResultStatus)
	purego.RegisterLibFunc(&ptrAPI_ISteamInventory_GetResultItems, lib, flatAPI_ISteamInventory_GetResultItems)
	purego.RegisterLibFunc(&ptrAPI_ISteamInventory_DestroyResult, lib, flatAPI_ISteamInventory_DestroyResult)

	// ISteamInput
	purego.RegisterLibFunc(&ptrAPI_SteamInput, lib, flatAPI_SteamInput)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetConnectedControllers, lib, flatAPI_ISteamInput_GetConnectedControllers)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetInputTypeForHandle, lib, flatAPI_ISteamInput_GetInputTypeForHandle)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_Init, lib, flatAPI_ISteamInput_Init)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_Shutdown, lib, flatAPI_ISteamInput_Shutdown)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_RunFrame, lib, flatAPI_ISteamInput_RunFrame)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_EnableDeviceCallbacks, lib, flatAPI_ISteamInput_EnableDeviceCallbacks)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetActionSetHandle, lib, flatAPI_ISteamInput_GetActionSetHandle)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_ActivateActionSet, lib, flatAPI_ISteamInput_ActivateActionSet)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetCurrentActionSet, lib, flatAPI_ISteamInput_GetCurrentActionSet)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_ActivateActionSetLayer, lib, flatAPI_ISteamInput_ActivateActionSetLayer)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_DeactivateActionSetLayer, lib, flatAPI_ISteamInput_DeactivateActionSetLayer)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_DeactivateAllActionSetLayers, lib, flatAPI_ISteamInput_DeactivateAllActionSetLayers)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetActiveActionSetLayers, lib, flatAPI_ISteamInput_GetActiveActionSetLayers)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetDigitalActionHandle, lib, flatAPI_ISteamInput_GetDigitalActionHandle)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetDigitalActionOrigins, lib, flatAPI_ISteamInput_GetDigitalActionOrigins)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetAnalogActionHandle, lib, flatAPI_ISteamInput_GetAnalogActionHandle)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetAnalogActionOrigins, lib, flatAPI_ISteamInput_GetAnalogActionOrigins)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_StopAnalogActionMomentum, lib, flatAPI_ISteamInput_StopAnalogActionMomentum)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_TriggerVibration, lib, flatAPI_ISteamInput_TriggerVibration)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_TriggerVibrationExtended, lib, flatAPI_ISteamInput_TriggerVibrationExtended)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_TriggerSimpleHapticEvent, lib, flatAPI_ISteamInput_TriggerSimpleHapticEvent)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_SetLEDColor, lib, flatAPI_ISteamInput_SetLEDColor)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_ShowBindingPanel, lib, flatAPI_ISteamInput_ShowBindingPanel)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetControllerForGamepadIndex, lib, flatAPI_ISteamInput_GetControllerForGamepadIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetGamepadIndexForController, lib, flatAPI_ISteamInput_GetGamepadIndexForController)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetStringForActionOrigin, lib, flatAPI_ISteamInput_GetStringForActionOrigin)
	registerOptionalFunc(&ptrAPI_ISteamInput_GetGlyphForActionOrigin, lib, flatAPI_ISteamInput_GetGlyphForActionOrigin)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetRemotePlaySessionID, lib, flatAPI_ISteamInput_GetRemotePlaySessionID)

	// ISteamRemotePlay
	registerOptionalFunc(&ptrAPI_SteamRemotePlay, lib, flatAPI_SteamRemotePlay)
	registerOptionalFunc(&ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether, lib, flatAPI_ISteamRemotePlay_BSessionRemotePlayTogether)
	registerOptionalFunc(&ptrAPI_ISteamRemotePlay_GetSessionGuestID, lib, flatAPI_ISteamRemotePlay_GetSessionGuestID)
	registerOptionalFunc(&ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar, lib, flatAPI_ISteamRemotePlay_GetSmallSessionAvatar)
	registerOptionalFunc(&ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar, lib, flatAPI_ISteamRemotePlay_GetMediumSessionAvatar)
	registerOptionalFunc(&ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar, lib, flatAPI_ISteamRemotePlay_GetLargeSessionAvatar)

	// ISteamRemoteStorage
	purego.RegisterLibFunc(&ptrAPI_SteamRemoteStorage, lib, flatAPI_SteamRemoteStorage)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileWrite, lib, flatAPI_ISteamRemoteStorage_FileWrite)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileRead, lib, flatAPI_ISteamRemoteStorage_FileRead)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileDelete, lib, flatAPI_ISteamRemoteStorage_FileDelete)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_GetFileSize, lib, flatAPI_ISteamRemoteStorage_GetFileSize)

	// ISteamUser
	purego.RegisterLibFunc(&ptrAPI_SteamUser, lib, flatAPI_SteamUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_AdvertiseGame, lib, flatAPI_ISteamUser_AdvertiseGame)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BeginAuthSession, lib, flatAPI_ISteamUser_BeginAuthSession)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BIsBehindNAT, lib, flatAPI_ISteamUser_BIsBehindNAT)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BIsPhoneIdentifying, lib, flatAPI_ISteamUser_BIsPhoneIdentifying)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BIsPhoneRequiringVerification, lib, flatAPI_ISteamUser_BIsPhoneRequiringVerification)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BIsPhoneVerified, lib, flatAPI_ISteamUser_BIsPhoneVerified)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BIsTwoFactorEnabled, lib, flatAPI_ISteamUser_BIsTwoFactorEnabled)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BLoggedOn, lib, flatAPI_ISteamUser_BLoggedOn)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_BSetDurationControlOnlineState, lib, flatAPI_ISteamUser_BSetDurationControlOnlineState)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_CancelAuthTicket, lib, flatAPI_ISteamUser_CancelAuthTicket)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_DecompressVoice, lib, flatAPI_ISteamUser_DecompressVoice)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_EndAuthSession, lib, flatAPI_ISteamUser_EndAuthSession)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetAuthSessionTicket, lib, flatAPI_ISteamUser_GetAuthSessionTicket)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetAuthTicketForWebApi, lib, flatAPI_ISteamUser_GetAuthTicketForWebApi)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetAvailableVoice, lib, flatAPI_ISteamUser_GetAvailableVoice)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetDurationControl, lib, flatAPI_ISteamUser_GetDurationControl)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetEncryptedAppTicket, lib, flatAPI_ISteamUser_GetEncryptedAppTicket)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetGameBadgeLevel, lib, flatAPI_ISteamUser_GetGameBadgeLevel)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetHSteamUser, lib, flatAPI_ISteamUser_GetHSteamUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetPlayerSteamLevel, lib, flatAPI_ISteamUser_GetPlayerSteamLevel)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetSteamID, lib, flatAPI_ISteamUser_GetSteamID)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetUserDataFolder, lib, flatAPI_ISteamUser_GetUserDataFolder)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetVoice, lib, flatAPI_ISteamUser_GetVoice)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetVoiceOptimalSampleRate, lib, flatAPI_ISteamUser_GetVoiceOptimalSampleRate)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_InitiateGameConnection, lib, flatAPI_ISteamUser_InitiateGameConnection)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_RequestEncryptedAppTicket, lib, flatAPI_ISteamUser_RequestEncryptedAppTicket)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_RequestStoreAuthURL, lib, flatAPI_ISteamUser_RequestStoreAuthURL)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_StartVoiceRecording, lib, flatAPI_ISteamUser_StartVoiceRecording)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_StopVoiceRecording, lib, flatAPI_ISteamUser_StopVoiceRecording)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_TerminateGameConnection, lib, flatAPI_ISteamUser_TerminateGameConnection)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_TrackAppUsageEvent, lib, flatAPI_ISteamUser_TrackAppUsageEvent)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_UserHasLicenseForApp, lib, flatAPI_ISteamUser_UserHasLicenseForApp)

	// ISteamUserStats
	purego.RegisterLibFunc(&ptrAPI_SteamUserStats, lib, flatAPI_SteamUserStats)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_GetAchievement, lib, flatAPI_ISteamUserStats_GetAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_SetAchievement, lib, flatAPI_ISteamUserStats_SetAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_ClearAchievement, lib, flatAPI_ISteamUserStats_ClearAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_StoreStats, lib, flatAPI_ISteamUserStats_StoreStats)

	// ISteamUtils
	purego.RegisterLibFunc(&ptrAPI_SteamUtils, lib, flatAPI_SteamUtils)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetSecondsSinceAppActive, lib, flatAPI_ISteamUtils_GetSecondsSinceAppActive)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetSecondsSinceComputerActive, lib, flatAPI_ISteamUtils_GetSecondsSinceComputerActive)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetConnectedUniverse, lib, flatAPI_ISteamUtils_GetConnectedUniverse)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetServerRealTime, lib, flatAPI_ISteamUtils_GetServerRealTime)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetIPCountry, lib, flatAPI_ISteamUtils_GetIPCountry)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetImageSize, lib, flatAPI_ISteamUtils_GetImageSize)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetImageRGBA, lib, flatAPI_ISteamUtils_GetImageRGBA)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetCurrentBatteryPower, lib, flatAPI_ISteamUtils_GetCurrentBatteryPower)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetAppID, lib, flatAPI_ISteamUtils_GetAppID)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_SetOverlayNotificationPosition, lib, flatAPI_ISteamUtils_SetOverlayNotificationPosition)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsAPICallCompleted, lib, flatAPI_ISteamUtils_IsAPICallCompleted)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetAPICallFailureReason, lib, flatAPI_ISteamUtils_GetAPICallFailureReason)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetAPICallResult, lib, flatAPI_ISteamUtils_GetAPICallResult)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_GetIPCCallCount, lib, flatAPI_ISteamUtils_GetIPCCallCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsOverlayEnabled, lib, flatAPI_ISteamUtils_IsOverlayEnabled)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_BOverlayNeedsPresent, lib, flatAPI_ISteamUtils_BOverlayNeedsPresent)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck, lib, flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput, lib, flatAPI_ISteamUtils_ShowFloatingGamepadTextInput)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_SetOverlayNotificationInset, lib, flatAPI_ISteamUtils_SetOverlayNotificationInset)

	// ISteamNetworkingUtils
	purego.RegisterLibFunc(&ptrAPI_SteamNetworkingUtils, lib, flatAPI_SteamNetworkingUtils)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingUtils_AllocateMessage, lib, flatAPI_ISteamNetworkingUtils_AllocateMessage)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingUtils_InitRelayNetworkAccess, lib, flatAPI_ISteamNetworkingUtils_InitRelayNetworkAccess)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingUtils_GetLocalTimestamp, lib, flatAPI_ISteamNetworkingUtils_GetLocalTimestamp)

	// ISteamGameServer
	purego.RegisterLibFunc(&ptrAPI_SteamGameServer, lib, flatAPI_SteamGameServer)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_AssociateWithClan, lib, flatAPI_ISteamGameServer_AssociateWithClan)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_BeginAuthSession, lib, flatAPI_ISteamGameServer_BeginAuthSession)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_BLoggedOn, lib, flatAPI_ISteamGameServer_BLoggedOn)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_BSecure, lib, flatAPI_ISteamGameServer_BSecure)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_BUpdateUserData, lib, flatAPI_ISteamGameServer_BUpdateUserData)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_CancelAuthTicket, lib, flatAPI_ISteamGameServer_CancelAuthTicket)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_ClearAllKeyValues, lib, flatAPI_ISteamGameServer_ClearAllKeyValues)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_ComputeNewPlayerCompatibility, lib, flatAPI_ISteamGameServer_ComputeNewPlayerCompatibility)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_CreateUnauthenticatedUserConnection, lib, flatAPI_ISteamGameServer_CreateUnauthenticatedUserConnection)
	registerOptionalFunc(&ptrAPI_ISteamGameServer_EnableHeartbeats, lib, flatAPI_ISteamGameServer_EnableHeartbeats)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_EndAuthSession, lib, flatAPI_ISteamGameServer_EndAuthSession)
	registerOptionalFunc(&ptrAPI_ISteamGameServer_ForceHeartbeat, lib, flatAPI_ISteamGameServer_ForceHeartbeat)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_GetAuthSessionTicket, lib, flatAPI_ISteamGameServer_GetAuthSessionTicket)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_GetGameplayStats, lib, flatAPI_ISteamGameServer_GetGameplayStats)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_GetNextOutgoingPacket, lib, flatAPI_ISteamGameServer_GetNextOutgoingPacket)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_GetPublicIP, lib, flatAPI_ISteamGameServer_GetPublicIP)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_GetServerReputation, lib, flatAPI_ISteamGameServer_GetServerReputation)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_GetSteamID, lib, flatAPI_ISteamGameServer_GetSteamID)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_HandleIncomingPacket, lib, flatAPI_ISteamGameServer_HandleIncomingPacket)
	registerOptionalFunc(&ptrAPI_ISteamGameServer_InitGameServer, lib, flatAPI_ISteamGameServer_InitGameServer)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_LogOff, lib, flatAPI_ISteamGameServer_LogOff)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_LogOn, lib, flatAPI_ISteamGameServer_LogOn)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_LogOnAnonymous, lib, flatAPI_ISteamGameServer_LogOnAnonymous)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_RequestUserGroupStatus, lib, flatAPI_ISteamGameServer_RequestUserGroupStatus)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SendUserConnectAndAuthenticate, lib, flatAPI_ISteamGameServer_SendUserConnectAndAuthenticate)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SendUserDisconnect, lib, flatAPI_ISteamGameServer_SendUserDisconnect)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetBotPlayerCount, lib, flatAPI_ISteamGameServer_SetBotPlayerCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetDedicatedServer, lib, flatAPI_ISteamGameServer_SetDedicatedServer)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetGameData, lib, flatAPI_ISteamGameServer_SetGameData)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetGameDescription, lib, flatAPI_ISteamGameServer_SetGameDescription)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetGameTags, lib, flatAPI_ISteamGameServer_SetGameTags)
	registerOptionalFunc(&ptrAPI_ISteamGameServer_SetHeartbeatInterval, lib, flatAPI_ISteamGameServer_SetHeartbeatInterval)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetKeyValue, lib, flatAPI_ISteamGameServer_SetKeyValue)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetMapName, lib, flatAPI_ISteamGameServer_SetMapName)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetMaxPlayerCount, lib, flatAPI_ISteamGameServer_SetMaxPlayerCount)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetModDir, lib, flatAPI_ISteamGameServer_SetModDir)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetPasswordProtected, lib, flatAPI_ISteamGameServer_SetPasswordProtected)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetProduct, lib, flatAPI_ISteamGameServer_SetProduct)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetRegion, lib, flatAPI_ISteamGameServer_SetRegion)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetServerName, lib, flatAPI_ISteamGameServer_SetServerName)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetSpectatorPort, lib, flatAPI_ISteamGameServer_SetSpectatorPort)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_SetSpectatorServerName, lib, flatAPI_ISteamGameServer_SetSpectatorServerName)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_UserHasLicenseForApp, lib, flatAPI_ISteamGameServer_UserHasLicenseForApp)
	purego.RegisterLibFunc(&ptrAPI_ISteamGameServer_WasRestartRequested, lib, flatAPI_ISteamGameServer_WasRestartRequested)

	// ISteamNetworkingMessages
	purego.RegisterLibFunc(&ptrAPI_SteamNetworkingMessages, lib, flatAPI_SteamNetworkingMessages)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingMessages_SendMessageToUser, lib, flatAPI_ISteamNetworkingMessages_SendMessageToUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel, lib, flatAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingMessages_AcceptSessionWithUser, lib, flatAPI_ISteamNetworkingMessages_AcceptSessionWithUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingMessages_CloseSessionWithUser, lib, flatAPI_ISteamNetworkingMessages_CloseSessionWithUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingMessages_CloseChannelWithUser, lib, flatAPI_ISteamNetworkingMessages_CloseChannelWithUser)

	// ISteamNetworkingSockets
	purego.RegisterLibFunc(&ptrAPI_SteamNetworkingSockets, lib, flatAPI_SteamNetworkingSockets)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CreateListenSocketIP, lib, flatAPI_ISteamNetworkingSockets_CreateListenSocketIP)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P, lib, flatAPI_ISteamNetworkingSockets_CreateListenSocketP2P)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ConnectByIPAddress, lib, flatAPI_ISteamNetworkingSockets_ConnectByIPAddress)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ConnectP2P, lib, flatAPI_ISteamNetworkingSockets_ConnectP2P)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_AcceptConnection, lib, flatAPI_ISteamNetworkingSockets_AcceptConnection)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CloseConnection, lib, flatAPI_ISteamNetworkingSockets_CloseConnection)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CloseListenSocket, lib, flatAPI_ISteamNetworkingSockets_CloseListenSocket)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_SendMessageToConnection, lib, flatAPI_ISteamNetworkingSockets_SendMessageToConnection)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection, lib, flatAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CreatePollGroup, lib, flatAPI_ISteamNetworkingSockets_CreatePollGroup)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_DestroyPollGroup, lib, flatAPI_ISteamNetworkingSockets_DestroyPollGroup)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_SetConnectionPollGroup, lib, flatAPI_ISteamNetworkingSockets_SetConnectionPollGroup)
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup, lib, flatAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup)

	registerInputStructReturns(lib)
}

func RestartAppIfNecessary(appID uint32) bool {
	mustLoad()
	return ptrAPI_RestartAppIfNecessary(appID)
}

func Init() error {
	l, err := ensureLoaded()
	if err != nil {
		return err
	}
	theLib = l

	if appID := os.Getenv("STEAM_APPID"); appID != "" {
		if err := os.WriteFile("steam_appid.txt", []byte(appID), 0644); err != nil {
			return fmt.Errorf("steamworks: failed to write steam_appid.txt: %w", err)
		}
	}

	var msg steamErrMsg
	if ptrAPI_InitFlat(uintptr(unsafe.Pointer(&msg))) != ESteamAPIInitResult_OK {
		return fmt.Errorf("steamworks: InitFlat failed: %s", msg.String())
	}
	return nil
}

func RunCallbacks() {
	mustLoad()
	ptrAPI_RunCallbacks()
}

// Shutdown shuts down the Steamworks API.
func Shutdown() {
	mustLoad()
	ptrAPI_Shutdown()
}

// IsSteamRunning reports whether the Steam client is currently running.
func IsSteamRunning() bool {
	mustLoad()
	return ptrAPI_IsSteamRunning()
}

// GetSteamInstallPath returns the Steam installation directory, if available.
func GetSteamInstallPath() string {
	mustLoad()
	return ptrAPI_GetSteamInstallPath()
}

// ReleaseCurrentThreadMemory releases per-thread memory used by the Steamworks API.
func ReleaseCurrentThreadMemory() {
	mustLoad()
	ptrAPI_ReleaseCurrentThreadMemory()
}

func SteamAppTicket() ISteamAppTicket {
	return SteamAppTicketRaw()
}

func SteamClient() ISteamClient {
	return SteamClientRaw()
}

func SteamController() ISteamController {
	return SteamControllerRaw()
}

func SteamGameCoordinator() ISteamGameCoordinator {
	return SteamGameCoordinatorRaw()
}

func SteamGameServerStats() ISteamGameServerStats {
	return SteamGameServerStatsRaw()
}

func SteamHTMLSurface() ISteamHTMLSurface {
	return SteamHTMLSurfaceRaw()
}

func SteamMatchmakingServers() ISteamMatchmakingServers {
	return SteamMatchmakingServersRaw()
}

func SteamMusic() ISteamMusic {
	return SteamMusicRaw()
}

func SteamNetworking() ISteamNetworking {
	return SteamNetworkingRaw()
}

func SteamRemotePlay() ISteamRemotePlay {
	mustLoad()
	if ptrAPI_SteamRemotePlay != nil {
		if ptr := ptrAPI_SteamRemotePlay(); ptr != 0 {
			return ISteamRemotePlay{ptr: ptr}
		}
	}
	return SteamRemotePlayRaw()
}

func SteamScreenshots() ISteamScreenshots {
	return SteamScreenshotsRaw()
}

func SteamTimeline() ISteamTimeline {
	return SteamTimelineRaw()
}

func SteamVideo() ISteamVideo {
	return SteamVideoRaw()
}

func SteamAPIClient() ISteamAPIClient {
	return SteamAPIClientRaw()
}

func SteamAPIGameServer() ISteamAPIGameServer {
	return SteamAPIGameServerRaw()
}

func resolveInterfaceFactory(symbols ...string) uintptr {
	mustLoad()
	for _, sym := range symbols {
		fn, err := lookupSymbolAddr(theLib.lib, sym)
		if err == nil && fn != 0 {
			ptr, _, _ := purego.SyscallN(fn)
			if ptr != 0 {
				return ptr
			}
		}
		if fn := fallbackResolveInterfaceSymbol(sym); fn != 0 {
			ptr, _, _ := purego.SyscallN(fn)
			if ptr != 0 {
				return ptr
			}
		}
	}
	return 0
}

// SteamAppTicketRaw returns the ISteamAppTicket interface pointer for purego/ffi calls.
func SteamAppTicketRaw() ISteamAppTicket {
	return ISteamAppTicket{ptr: resolveInterfaceFactory("SteamAPI_SteamAppTicket_v001")}
}

// SteamClientRaw returns the ISteamClient interface pointer for purego/ffi calls.
func SteamClientRaw() ISteamClient {
	return ISteamClient{ptr: resolveInterfaceFactory("SteamAPI_SteamClient_v022", "SteamAPI_SteamClient_v021", "SteamAPI_SteamClient_v020")}
}

// SteamControllerRaw returns the ISteamController interface pointer for purego/ffi calls.
func SteamControllerRaw() ISteamController {
	return ISteamController{ptr: resolveInterfaceFactory("SteamAPI_SteamController_v008", "SteamAPI_SteamController_v007")}
}

// SteamGameCoordinatorRaw returns the ISteamGameCoordinator interface pointer for purego/ffi calls.
func SteamGameCoordinatorRaw() ISteamGameCoordinator {
	return ISteamGameCoordinator{ptr: resolveInterfaceFactory("SteamAPI_SteamGameCoordinator_v001")}
}

// SteamGameServerStatsRaw returns the ISteamGameServerStats interface pointer for purego/ffi calls.
func SteamGameServerStatsRaw() ISteamGameServerStats {
	return ISteamGameServerStats{ptr: resolveInterfaceFactory("SteamAPI_SteamGameServerStats_v001")}
}

// SteamHTMLSurfaceRaw returns the ISteamHTMLSurface interface pointer for purego/ffi calls.
func SteamHTMLSurfaceRaw() ISteamHTMLSurface {
	return ISteamHTMLSurface{ptr: resolveInterfaceFactory("SteamAPI_SteamHTMLSurface_v005")}
}

// SteamMatchmakingServersRaw returns the ISteamMatchmakingServers interface pointer for purego/ffi calls.
func SteamMatchmakingServersRaw() ISteamMatchmakingServers {
	return ISteamMatchmakingServers{ptr: resolveInterfaceFactory("SteamAPI_SteamMatchmakingServers_v002")}
}

func ptrSlice(items []uintptr) uintptr {
	if len(items) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&items[0]))
}

func steamIDSlicePtr(items []CSteamID) uintptr {
	if len(items) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&items[0]))
}

func (s ISteamMatchmakingServers) RequestFavoritesServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest {
	return ptrAPI_ISteamMatchmakingServers_RequestFavoritesServerList(s.ptr, appID, ptrSlice(filters), uint32(len(filters)), response)
}

func (s ISteamMatchmakingServers) RequestFriendsServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest {
	return ptrAPI_ISteamMatchmakingServers_RequestFriendsServerList(s.ptr, appID, ptrSlice(filters), uint32(len(filters)), response)
}

func (s ISteamMatchmakingServers) RequestHistoryServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest {
	return ptrAPI_ISteamMatchmakingServers_RequestHistoryServerList(s.ptr, appID, ptrSlice(filters), uint32(len(filters)), response)
}

func (s ISteamMatchmakingServers) RequestInternetServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest {
	return ptrAPI_ISteamMatchmakingServers_RequestInternetServerList(s.ptr, appID, ptrSlice(filters), uint32(len(filters)), response)
}

func (s ISteamMatchmakingServers) RequestLANServerList(appID AppId_t, response uintptr) HServerListRequest {
	return ptrAPI_ISteamMatchmakingServers_RequestLANServerList(s.ptr, appID, response)
}

func (s ISteamMatchmakingServers) RequestSpectatorServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest {
	return ptrAPI_ISteamMatchmakingServers_RequestSpectatorServerList(s.ptr, appID, ptrSlice(filters), uint32(len(filters)), response)
}

func (s ISteamMatchmakingServers) ReleaseRequest(request HServerListRequest) {
	ptrAPI_ISteamMatchmakingServers_ReleaseRequest(s.ptr, request)
}

func (s ISteamMatchmakingServers) GetServerDetails(request HServerListRequest, server int) uintptr {
	return ptrAPI_ISteamMatchmakingServers_GetServerDetails(s.ptr, request, int32(server))
}

func (s ISteamMatchmakingServers) CancelQuery(request HServerListRequest) {
	ptrAPI_ISteamMatchmakingServers_CancelQuery(s.ptr, request)
}

func (s ISteamMatchmakingServers) RefreshQuery(request HServerListRequest) {
	ptrAPI_ISteamMatchmakingServers_RefreshQuery(s.ptr, request)
}

func (s ISteamMatchmakingServers) IsRefreshing(request HServerListRequest) bool {
	return ptrAPI_ISteamMatchmakingServers_IsRefreshing(s.ptr, request)
}

func (s ISteamMatchmakingServers) GetServerCount(request HServerListRequest) int {
	return int(ptrAPI_ISteamMatchmakingServers_GetServerCount(s.ptr, request))
}

func (s ISteamMatchmakingServers) RefreshServer(request HServerListRequest, server int) {
	ptrAPI_ISteamMatchmakingServers_RefreshServer(s.ptr, request, int32(server))
}

func (s ISteamMatchmakingServers) PingServer(ip uint32, port uint16, response uintptr) HServerQuery {
	return ptrAPI_ISteamMatchmakingServers_PingServer(s.ptr, ip, port, response)
}

func (s ISteamMatchmakingServers) PlayerDetails(ip uint32, port uint16, response uintptr) HServerQuery {
	return ptrAPI_ISteamMatchmakingServers_PlayerDetails(s.ptr, ip, port, response)
}

func (s ISteamMatchmakingServers) ServerRules(ip uint32, port uint16, response uintptr) HServerQuery {
	return ptrAPI_ISteamMatchmakingServers_ServerRules(s.ptr, ip, port, response)
}

func (s ISteamMatchmakingServers) CancelServerQuery(query HServerQuery) {
	ptrAPI_ISteamMatchmakingServers_CancelServerQuery(s.ptr, query)
}

// SteamMusicRaw returns the ISteamMusic interface pointer for purego/ffi calls.
func SteamMusicRaw() ISteamMusic {
	return ISteamMusic{ptr: resolveInterfaceFactory("SteamAPI_SteamMusic_v001")}
}

// SteamNetworkingRaw returns the legacy ISteamNetworking interface pointer for purego/ffi calls.
func SteamNetworkingRaw() ISteamNetworking {
	return ISteamNetworking{ptr: resolveInterfaceFactory("SteamAPI_SteamNetworking_v006", "SteamAPI_SteamNetworking_v005")}
}

// SteamRemotePlayRaw returns the ISteamRemotePlay interface pointer for purego/ffi calls.
func SteamRemotePlayRaw() ISteamRemotePlay {
	return ISteamRemotePlay{ptr: resolveInterfaceFactory("SteamAPI_SteamRemotePlay_v001")}
}

func (s ISteamRemotePlay) BSessionRemotePlayTogether(sessionID uint32) bool {
	if ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether == nil {
		return false
	}
	return ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether(s.ptr, sessionID)
}

func (s ISteamRemotePlay) GetSessionGuestID(sessionID uint32) CSteamID {
	if ptrAPI_ISteamRemotePlay_GetSessionGuestID == nil {
		return 0
	}
	return ptrAPI_ISteamRemotePlay_GetSessionGuestID(s.ptr, sessionID)
}

func (s ISteamRemotePlay) GetSmallSessionAvatar(sessionID uint32) int32 {
	if ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar == nil {
		return -1
	}
	return ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar(s.ptr, sessionID)
}

func (s ISteamRemotePlay) GetMediumSessionAvatar(sessionID uint32) int32 {
	if ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar == nil {
		return -1
	}
	return ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar(s.ptr, sessionID)
}

func (s ISteamRemotePlay) GetLargeSessionAvatar(sessionID uint32) int32 {
	if ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar == nil {
		return -1
	}
	return ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar(s.ptr, sessionID)
}

// SteamScreenshotsRaw returns the ISteamScreenshots interface pointer for purego/ffi calls.
func SteamScreenshotsRaw() ISteamScreenshots {
	return ISteamScreenshots{ptr: resolveInterfaceFactory("SteamAPI_SteamScreenshots_v003")}
}

// SteamTimelineRaw returns the ISteamTimeline interface pointer for purego/ffi calls.
func SteamTimelineRaw() ISteamTimeline {
	return ISteamTimeline{ptr: resolveInterfaceFactory("SteamAPI_SteamTimeline_v001")}
}

// SteamVideoRaw returns the ISteamVideo interface pointer for purego/ffi calls.
func SteamVideoRaw() ISteamVideo {
	return ISteamVideo{ptr: resolveInterfaceFactory("SteamAPI_SteamVideo_v002")}
}

// SteamAPIClientRaw returns the steam_api client foundation handle for purego/ffi calls.
func SteamAPIClientRaw() ISteamAPIClient {
	return ISteamAPIClient{ptr: resolveInterfaceFactory("SteamAPI_SteamClient_v022", "SteamAPI_SteamClient_v021", "SteamAPI_SteamClient_v020")}
}

// SteamAPIGameServerRaw returns the steam_gameserver foundation handle for purego/ffi calls.
func SteamAPIGameServerRaw() ISteamAPIGameServer {
	return ISteamAPIGameServer{ptr: resolveInterfaceFactory("SteamAPI_SteamGameServer_v015")}
}

func SteamApps() ISteamApps {
	mustLoad()
	return steamApps(ptrAPI_SteamApps())
}

// SteamAppsV008 returns the v008 apps interface.
func SteamAppsV008() ISteamApps {
	return SteamApps()
}

type steamApps uintptr

func (s steamApps) BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool) {
	var name [4096]byte
	v := ptrAPI_ISteamApps_BGetDLCDataByIndex(uintptr(s), int32(iDLC), uintptr(unsafe.Pointer(&appID)), uintptr(unsafe.Pointer(&available)), uintptr(unsafe.Pointer(&name[0])), int32(len(name)))
	return appID, available, cStringToGo(name[:]), v
}

func (s steamApps) BIsSubscribed() bool {
	return ptrAPI_ISteamApps_BIsSubscribed(uintptr(s))
}

func (s steamApps) BIsLowViolence() bool {
	return ptrAPI_ISteamApps_BIsLowViolence(uintptr(s))
}

func (s steamApps) BIsCybercafe() bool {
	return ptrAPI_ISteamApps_BIsCybercafe(uintptr(s))
}

func (s steamApps) BIsVACBanned() bool {
	return ptrAPI_ISteamApps_BIsVACBanned(uintptr(s))
}

func (s steamApps) BIsDlcInstalled(appID AppId_t) bool {
	return ptrAPI_ISteamApps_BIsDlcInstalled(uintptr(s), appID)
}

func (s steamApps) BIsSubscribedApp(appID AppId_t) bool {
	return ptrAPI_ISteamApps_BIsSubscribedApp(uintptr(s), appID)
}

func (s steamApps) BIsSubscribedFromFreeWeekend() bool {
	return ptrAPI_ISteamApps_BIsSubscribedFromFreeWeekend(uintptr(s))
}

func (s steamApps) BIsSubscribedFromFamilySharing() bool {
	return ptrAPI_ISteamApps_BIsSubscribedFromFamilySharing(uintptr(s))
}

func (s steamApps) BIsTimedTrial() (allowedSeconds, playedSeconds uint32, ok bool) {
	ok = ptrAPI_ISteamApps_BIsTimedTrial(uintptr(s), uintptr(unsafe.Pointer(&allowedSeconds)), uintptr(unsafe.Pointer(&playedSeconds)))
	return
}

func (s steamApps) BIsAppInstalled(appID AppId_t) bool {
	return ptrAPI_ISteamApps_BIsAppInstalled(uintptr(s), appID)
}

func (s steamApps) GetAvailableGameLanguages() string {
	return unique.Make(ptrAPI_ISteamApps_GetAvailableGameLanguages(uintptr(s))).Value()
}

func (s steamApps) GetEarliestPurchaseUnixTime(appID AppId_t) uint32 {
	return ptrAPI_ISteamApps_GetEarliestPurchaseUnixTime(uintptr(s), appID)
}

func (s steamApps) GetAppInstallDir(appID AppId_t) string {
	var path [4096]byte
	v := ptrAPI_ISteamApps_GetAppInstallDir(uintptr(s), appID, uintptr(unsafe.Pointer(&path[0])), int32(len(path)))
	if v == 0 {
		return ""
	}
	return string(path[:v-1])
}

func (s steamApps) GetCurrentGameLanguage() string {
	return unique.Make(ptrAPI_ISteamApps_GetCurrentGameLanguage(uintptr(s))).Value()
}

func (s steamApps) GetDLCCount() int32 {
	return ptrAPI_ISteamApps_GetDLCCount(uintptr(s))
}

func (s steamApps) GetCurrentBetaName() (string, bool) {
	var name [4096]byte
	ok := ptrAPI_ISteamApps_GetCurrentBetaName(uintptr(s), uintptr(unsafe.Pointer(&name[0])), int32(len(name)))
	if !ok {
		return "", false
	}
	return cStringToGo(name[:]), true
}

func (s steamApps) GetInstalledDepots(appID AppId_t) []DepotId_t {
	depots := make([]DepotId_t, 32)
	for {
		count := ptrAPI_ISteamApps_GetInstalledDepots(uintptr(s), appID, uintptr(unsafe.Pointer(&depots[0])), uint32(len(depots)))
		if int(count) <= len(depots) {
			return depots[:count]
		}
		depots = make([]DepotId_t, count)
	}
}

func (s steamApps) GetAppOwner() CSteamID {
	return ptrAPI_ISteamApps_GetAppOwner(uintptr(s))
}

func (s steamApps) GetLaunchQueryParam(key string) string {
	return ptrAPI_ISteamApps_GetLaunchQueryParam(uintptr(s), key)
}

func (s steamApps) GetDlcDownloadProgress(appID AppId_t) (downloaded, total uint64, ok bool) {
	ok = ptrAPI_ISteamApps_GetDlcDownloadProgress(uintptr(s), appID, uintptr(unsafe.Pointer(&downloaded)), uintptr(unsafe.Pointer(&total)))
	return
}

func (s steamApps) GetAppBuildId() int32 {
	return ptrAPI_ISteamApps_GetAppBuildId(uintptr(s))
}

func (s steamApps) GetFileDetails(filename string) SteamAPICall_t {
	return ptrAPI_ISteamApps_GetFileDetails(uintptr(s), filename)
}

func (s steamApps) GetLaunchCommandLine(bufferSize int) string {
	if bufferSize <= 0 {
		bufferSize = 4096
	}
	buf := make([]byte, bufferSize)
	v := ptrAPI_ISteamApps_GetLaunchCommandLine(uintptr(s), uintptr(unsafe.Pointer(&buf[0])), int32(len(buf)))
	if v == 0 {
		return ""
	}
	return string(buf[:v])
}

func (s steamApps) GetNumBetas() (total int, available int, private int) {
	total = int(ptrAPI_ISteamApps_GetNumBetas(uintptr(s), uintptr(unsafe.Pointer(&available)), uintptr(unsafe.Pointer(&private))))
	return
}

func (s steamApps) GetBetaInfo(index int) (flags uint32, buildID uint32, lastUpdated uint32, name string, description string, ok bool) {
	var nameBuf [4096]byte
	var descBuf [4096]byte
	ok = ptrAPI_ISteamApps_GetBetaInfo(uintptr(s), int32(index), uintptr(unsafe.Pointer(&flags)), uintptr(unsafe.Pointer(&buildID)), uintptr(unsafe.Pointer(&lastUpdated)), uintptr(unsafe.Pointer(&nameBuf[0])), int32(len(nameBuf)), uintptr(unsafe.Pointer(&descBuf[0])), int32(len(descBuf)))
	if !ok {
		return 0, 0, 0, "", "", false
	}
	return flags, buildID, lastUpdated, cStringToGo(nameBuf[:]), cStringToGo(descBuf[:]), true
}

func (s steamApps) InstallDLC(appID AppId_t) {
	ptrAPI_ISteamApps_InstallDLC(uintptr(s), appID)
}

func (s steamApps) UninstallDLC(appID AppId_t) {
	ptrAPI_ISteamApps_UninstallDLC(uintptr(s), appID)
}

func (s steamApps) RequestAppProofOfPurchaseKey(appID AppId_t) {
	ptrAPI_ISteamApps_RequestAppProofOfPurchaseKey(uintptr(s), appID)
}

func (s steamApps) RequestAllProofOfPurchaseKeys() {
	ptrAPI_ISteamApps_RequestAllProofOfPurchaseKeys(uintptr(s))
}

func (s steamApps) MarkContentCorrupt(missingFilesOnly bool) bool {
	return ptrAPI_ISteamApps_MarkContentCorrupt(uintptr(s), missingFilesOnly)
}

func (s steamApps) SetDlcContext(appID AppId_t) bool {
	return ptrAPI_ISteamApps_SetDlcContext(uintptr(s), appID)
}

func (s steamApps) SetActiveBeta(name string) bool {
	return ptrAPI_ISteamApps_SetActiveBeta(uintptr(s), name)
}

func SteamFriends() ISteamFriends {
	mustLoad()
	return steamFriends(ptrAPI_SteamFriends())
}

// SteamFriendsV018 returns the v018 friends interface.
func SteamFriendsV018() ISteamFriends {
	return SteamFriends()
}

type steamFriends uintptr

func (s steamFriends) GetPersonaName() string {
	return unique.Make(ptrAPI_ISteamFriends_GetPersonaName(uintptr(s))).Value()
}

func (s steamFriends) GetPersonaState() EPersonaState {
	return EPersonaState(ptrAPI_ISteamFriends_GetPersonaState(uintptr(s)))
}

func (s steamFriends) GetFriendCount(flags EFriendFlags) int {
	return int(ptrAPI_ISteamFriends_GetFriendCount(uintptr(s), int32(flags)))
}

func (s steamFriends) GetFriendByIndex(index int, flags EFriendFlags) CSteamID {
	return ptrAPI_ISteamFriends_GetFriendByIndex(uintptr(s), int32(index), int32(flags))
}

func (s steamFriends) Friends(flags EFriendFlags) iter.Seq[CSteamID] {
	return func(yield func(CSteamID) bool) {
		count := s.GetFriendCount(flags)
		for i := 0; i < count; i++ {
			if !yield(s.GetFriendByIndex(i, flags)) {
				return
			}
		}
	}
}

func (s steamFriends) GetFriendRelationship(friend CSteamID) EFriendRelationship {
	return EFriendRelationship(ptrAPI_ISteamFriends_GetFriendRelationship(uintptr(s), friend))
}

func (s steamFriends) GetFriendPersonaState(friend CSteamID) EPersonaState {
	return EPersonaState(ptrAPI_ISteamFriends_GetFriendPersonaState(uintptr(s), friend))
}

func (s steamFriends) GetFriendPersonaName(friend CSteamID) string {
	return unique.Make(ptrAPI_ISteamFriends_GetFriendPersonaName(uintptr(s), friend)).Value()
}

func (s steamFriends) GetFriendPersonaNameHistory(friend CSteamID, index int) string {
	return unique.Make(ptrAPI_ISteamFriends_GetFriendPersonaNameHistory(uintptr(s), friend, int32(index))).Value()
}

func (s steamFriends) GetFriendSteamLevel(friend CSteamID) int {
	return int(ptrAPI_ISteamFriends_GetFriendSteamLevel(uintptr(s), friend))
}

func (s steamFriends) GetSmallFriendAvatar(friend CSteamID) int32 {
	return ptrAPI_ISteamFriends_GetSmallFriendAvatar(uintptr(s), friend)
}

func (s steamFriends) GetMediumFriendAvatar(friend CSteamID) int32 {
	return ptrAPI_ISteamFriends_GetMediumFriendAvatar(uintptr(s), friend)
}

func (s steamFriends) GetLargeFriendAvatar(friend CSteamID) int32 {
	return ptrAPI_ISteamFriends_GetLargeFriendAvatar(uintptr(s), friend)
}

func (s steamFriends) SetRichPresence(key, value string) bool {
	return ptrAPI_ISteamFriends_SetRichPresence(uintptr(s), key, value)
}

func (s steamFriends) GetFriendGamePlayed(friend CSteamID) (FriendGameInfo, bool) {
	var info FriendGameInfo
	ok := ptrAPI_ISteamFriends_GetFriendGamePlayed(uintptr(s), friend, uintptr(unsafe.Pointer(&info)))
	return info, ok
}

func (s steamFriends) InviteUserToGame(friend CSteamID, connectString string) bool {
	return ptrAPI_ISteamFriends_InviteUserToGame(uintptr(s), friend, connectString)
}

func (s steamFriends) ActivateGameOverlay(dialog string) {
	ptrAPI_ISteamFriends_ActivateGameOverlay(uintptr(s), dialog)
}

func (s steamFriends) ActivateGameOverlayToUser(dialog string, steamID CSteamID) {
	ptrAPI_ISteamFriends_ActivateGameOverlayToUser(uintptr(s), dialog, steamID)
}

func (s steamFriends) ActivateGameOverlayToWebPage(url string, mode EActivateGameOverlayToWebPageMode) {
	ptrAPI_ISteamFriends_ActivateGameOverlayToWebPage(uintptr(s), url, mode)
}

func (s steamFriends) ActivateGameOverlayToStore(appID AppId_t, flag EOverlayToStoreFlag) {
	ptrAPI_ISteamFriends_ActivateGameOverlayToStore(uintptr(s), appID, flag)
}

func (s steamFriends) ActivateGameOverlayInviteDialog(lobbyID CSteamID) {
	ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialog(uintptr(s), lobbyID)
}

func (s steamFriends) ActivateGameOverlayInviteDialogConnectString(connectString string) {
	ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString(uintptr(s), connectString)
}

func SteamMatchmaking() ISteamMatchmaking {
	mustLoad()
	return steamMatchmaking(ptrAPI_SteamMatchmaking())
}

// SteamMatchmakingV009 returns the v009 matchmaking interface.
func SteamMatchmakingV009() ISteamMatchmaking {
	return SteamMatchmaking()
}

type steamMatchmaking uintptr

func (s steamMatchmaking) GetFavoriteGameCount() int {
	return int(ptrAPI_ISteamMatchmaking_GetFavoriteGameCount(uintptr(s)))
}

func (s steamMatchmaking) GetFavoriteGame(index int) (FavoriteGame, bool) {
	var favorite FavoriteGame
	ok := ptrAPI_ISteamMatchmaking_GetFavoriteGame(
		uintptr(s),
		int32(index),
		uintptr(unsafe.Pointer(&favorite.AppID)),
		uintptr(unsafe.Pointer(&favorite.IP)),
		uintptr(unsafe.Pointer(&favorite.ConnectionPort)),
		uintptr(unsafe.Pointer(&favorite.QueryPort)),
		uintptr(unsafe.Pointer(&favorite.Flags)),
		uintptr(unsafe.Pointer(&favorite.LastPlayedOnServerTime)),
	)
	return favorite, ok
}

func (s steamMatchmaking) AddFavoriteGame(appID AppId_t, ip uint32, connectionPort, queryPort uint16, flags, lastPlayedOnServerTime uint32) int {
	return int(ptrAPI_ISteamMatchmaking_AddFavoriteGame(uintptr(s), appID, ip, connectionPort, queryPort, flags, lastPlayedOnServerTime))
}

func (s steamMatchmaking) RemoveFavoriteGame(appID AppId_t, ip uint32, connectionPort, queryPort uint16, flags uint32) bool {
	return ptrAPI_ISteamMatchmaking_RemoveFavoriteGame(uintptr(s), appID, ip, connectionPort, queryPort, flags)
}

func (s steamMatchmaking) RequestLobbyList() SteamAPICall_t {
	return ptrAPI_ISteamMatchmaking_RequestLobbyList(uintptr(s))
}

func (s steamMatchmaking) AddRequestLobbyListStringFilter(key, value string, comparisonType ELobbyComparison) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter(uintptr(s), key, value, comparisonType)
}

func (s steamMatchmaking) AddRequestLobbyListNumericalFilter(key string, value int, comparisonType ELobbyComparison) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter(uintptr(s), key, int32(value), comparisonType)
}

func (s steamMatchmaking) AddRequestLobbyListNearValueFilter(key string, value int) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter(uintptr(s), key, int32(value))
}

func (s steamMatchmaking) AddRequestLobbyListFilterSlotsAvailable(slotsAvailable int) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable(uintptr(s), int32(slotsAvailable))
}

func (s steamMatchmaking) AddRequestLobbyListDistanceFilter(distanceFilter ELobbyDistanceFilter) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter(uintptr(s), distanceFilter)
}

func (s steamMatchmaking) AddRequestLobbyListResultCountFilter(maxResults int) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter(uintptr(s), int32(maxResults))
}

func (s steamMatchmaking) AddRequestLobbyListCompatibleMembersFilter(lobbyID CSteamID) {
	ptrAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter(uintptr(s), lobbyID)
}

func (s steamMatchmaking) GetLobbyByIndex(index int) CSteamID {
	return ptrAPI_ISteamMatchmaking_GetLobbyByIndex(uintptr(s), int32(index))
}

func (s steamMatchmaking) CreateLobby(lobbyType ELobbyType, maxMembers int) SteamAPICall_t {
	return ptrAPI_ISteamMatchmaking_CreateLobby(uintptr(s), lobbyType, int32(maxMembers))
}

func (s steamMatchmaking) JoinLobby(lobbyID CSteamID) SteamAPICall_t {
	return ptrAPI_ISteamMatchmaking_JoinLobby(uintptr(s), lobbyID)
}

func (s steamMatchmaking) LeaveLobby(lobbyID CSteamID) {
	ptrAPI_ISteamMatchmaking_LeaveLobby(uintptr(s), lobbyID)
}

func (s steamMatchmaking) InviteUserToLobby(lobbyID, invitee CSteamID) bool {
	return ptrAPI_ISteamMatchmaking_InviteUserToLobby(uintptr(s), lobbyID, invitee)
}

func (s steamMatchmaking) SetLobbyMemberLimit(lobbyID CSteamID, maxMembers int) bool {
	return ptrAPI_ISteamMatchmaking_SetLobbyMemberLimit(uintptr(s), lobbyID, int32(maxMembers))
}

func (s steamMatchmaking) GetLobbyMemberLimit(lobbyID CSteamID) int {
	return int(ptrAPI_ISteamMatchmaking_GetLobbyMemberLimit(uintptr(s), lobbyID))
}

func (s steamMatchmaking) SetLobbyType(lobbyID CSteamID, lobbyType ELobbyType) bool {
	return ptrAPI_ISteamMatchmaking_SetLobbyType(uintptr(s), lobbyID, lobbyType)
}

func (s steamMatchmaking) SetLobbyJoinable(lobbyID CSteamID, joinable bool) bool {
	return ptrAPI_ISteamMatchmaking_SetLobbyJoinable(uintptr(s), lobbyID, joinable)
}

func (s steamMatchmaking) GetLobbyOwner(lobbyID CSteamID) CSteamID {
	return ptrAPI_ISteamMatchmaking_GetLobbyOwner(uintptr(s), lobbyID)
}

func (s steamMatchmaking) SetLobbyOwner(lobbyID, owner CSteamID) bool {
	return ptrAPI_ISteamMatchmaking_SetLobbyOwner(uintptr(s), lobbyID, owner)
}

func (s steamMatchmaking) SetLinkedLobby(lobbyID, lobbyDependent CSteamID) bool {
	return ptrAPI_ISteamMatchmaking_SetLinkedLobby(uintptr(s), lobbyID, lobbyDependent)
}

func (s steamMatchmaking) GetNumLobbyMembers(lobbyID CSteamID) int {
	return int(ptrAPI_ISteamMatchmaking_GetNumLobbyMembers(uintptr(s), lobbyID))
}

func (s steamMatchmaking) GetLobbyMemberByIndex(lobbyID CSteamID, memberIndex int) CSteamID {
	return ptrAPI_ISteamMatchmaking_GetLobbyMemberByIndex(uintptr(s), lobbyID, int32(memberIndex))
}

func (s steamMatchmaking) LobbyMembers(lobbyID CSteamID) iter.Seq[CSteamID] {
	return func(yield func(CSteamID) bool) {
		count := s.GetNumLobbyMembers(lobbyID)
		for i := 0; i < count; i++ {
			if !yield(s.GetLobbyMemberByIndex(lobbyID, i)) {
				return
			}
		}
	}
}

func (s steamMatchmaking) SetLobbyData(lobbyID CSteamID, key, value string) bool {
	return ptrAPI_ISteamMatchmaking_SetLobbyData(uintptr(s), lobbyID, key, value)
}

func (s steamMatchmaking) GetLobbyData(lobbyID CSteamID, key string) string {
	return ptrAPI_ISteamMatchmaking_GetLobbyData(uintptr(s), lobbyID, key)
}

func (s steamMatchmaking) DeleteLobbyData(lobbyID CSteamID, key string) bool {
	return ptrAPI_ISteamMatchmaking_DeleteLobbyData(uintptr(s), lobbyID, key)
}

func (s steamMatchmaking) GetLobbyDataCount(lobbyID CSteamID) int {
	return int(ptrAPI_ISteamMatchmaking_GetLobbyDataCount(uintptr(s), lobbyID))
}

func (s steamMatchmaking) GetLobbyDataByIndex(lobbyID CSteamID, lobbyDataIndex int) (key, value string, ok bool) {
	var keyBuf [256]byte
	var valueBuf [4096]byte
	ok = ptrAPI_ISteamMatchmaking_GetLobbyDataByIndex(
		uintptr(s),
		lobbyID,
		int32(lobbyDataIndex),
		uintptr(unsafe.Pointer(&keyBuf[0])),
		int32(len(keyBuf)),
		uintptr(unsafe.Pointer(&valueBuf[0])),
		int32(len(valueBuf)),
	)
	if !ok {
		return "", "", false
	}
	key = cStringToGo(keyBuf[:])
	value = cStringToGo(valueBuf[:])
	return key, value, true
}

func (s steamMatchmaking) SetLobbyMemberData(lobbyID CSteamID, key, value string) {
	ptrAPI_ISteamMatchmaking_SetLobbyMemberData(uintptr(s), lobbyID, key, value)
}

func (s steamMatchmaking) GetLobbyMemberData(lobbyID, user CSteamID, key string) string {
	return ptrAPI_ISteamMatchmaking_GetLobbyMemberData(uintptr(s), lobbyID, user, key)
}

func (s steamMatchmaking) SendLobbyChatMsg(lobbyID CSteamID, msgBody []byte) bool {
	var ptr uintptr
	if len(msgBody) != 0 {
		ptr = uintptr(unsafe.Pointer(&msgBody[0]))
	}
	return ptrAPI_ISteamMatchmaking_SendLobbyChatMsg(uintptr(s), lobbyID, ptr, int32(len(msgBody)))
}

func (s steamMatchmaking) GetLobbyChatEntry(lobbyID CSteamID, chatID int, data []byte) (user CSteamID, entryType EChatEntryType, bytesCopied int) {
	var ptr uintptr
	if len(data) != 0 {
		ptr = uintptr(unsafe.Pointer(&data[0]))
	}
	var rawEntryType int32
	bytesCopied = int(ptrAPI_ISteamMatchmaking_GetLobbyChatEntry(
		uintptr(s),
		lobbyID,
		int32(chatID),
		uintptr(unsafe.Pointer(&user)),
		ptr,
		int32(len(data)),
		uintptr(unsafe.Pointer(&rawEntryType)),
	))
	entryType = EChatEntryType(rawEntryType)
	return
}

func (s steamMatchmaking) RequestLobbyData(lobbyID CSteamID) bool {
	return ptrAPI_ISteamMatchmaking_RequestLobbyData(uintptr(s), lobbyID)
}

func (s steamMatchmaking) SetLobbyGameServer(lobbyID CSteamID, ip uint32, port uint16, server CSteamID) {
	ptrAPI_ISteamMatchmaking_SetLobbyGameServer(uintptr(s), lobbyID, ip, port, server)
}

func (s steamMatchmaking) GetLobbyGameServer(lobbyID CSteamID) (ip uint32, port uint16, server CSteamID, ok bool) {
	ok = ptrAPI_ISteamMatchmaking_GetLobbyGameServer(uintptr(s), lobbyID, uintptr(unsafe.Pointer(&ip)), uintptr(unsafe.Pointer(&port)), uintptr(unsafe.Pointer(&server)))
	return
}

func (s steamMatchmaking) CheckForPSNGameBootInvite(lobbyID *CSteamID) bool {
	if ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite == nil {
		return false
	}
	var ptr uintptr
	if lobbyID != nil {
		ptr = uintptr(unsafe.Pointer(lobbyID))
	}
	return ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite(uintptr(s), ptr)
}

func SteamHTTP() ISteamHTTP {
	mustLoad()
	return steamHTTP(ptrAPI_SteamHTTP())
}

// SteamHTTPV003 returns the v003 HTTP interface.
func SteamHTTPV003() ISteamHTTP {
	return SteamHTTP()
}

type steamHTTP uintptr

func (s steamHTTP) CreateHTTPRequest(method EHTTPMethod, absoluteURL string) HTTPRequestHandle {
	return ptrAPI_ISteamHTTP_CreateHTTPRequest(uintptr(s), int32(method), absoluteURL)
}

func (s steamHTTP) SetHTTPRequestHeaderValue(request HTTPRequestHandle, headerName, headerValue string) bool {
	return ptrAPI_ISteamHTTP_SetHTTPRequestHeaderValue(uintptr(s), request, headerName, headerValue)
}

func (s steamHTTP) SendHTTPRequest(request HTTPRequestHandle) (SteamAPICall_t, bool) {
	var call SteamAPICall_t
	ok := ptrAPI_ISteamHTTP_SendHTTPRequest(uintptr(s), request, uintptr(unsafe.Pointer(&call)))
	return call, ok
}

func (s steamHTTP) GetHTTPResponseBodySize(request HTTPRequestHandle) (uint32, bool) {
	var size uint32
	ok := ptrAPI_ISteamHTTP_GetHTTPResponseBodySize(uintptr(s), request, uintptr(unsafe.Pointer(&size)))
	return size, ok
}

func (s steamHTTP) GetHTTPResponseBodyData(request HTTPRequestHandle, buffer []byte) bool {
	if len(buffer) == 0 {
		return false
	}
	return ptrAPI_ISteamHTTP_GetHTTPResponseBodyData(uintptr(s), request, uintptr(unsafe.Pointer(&buffer[0])), uint32(len(buffer)))
}

func (s steamHTTP) ReleaseHTTPRequest(request HTTPRequestHandle) bool {
	return ptrAPI_ISteamHTTP_ReleaseHTTPRequest(uintptr(s), request)
}

func SteamUGC() ISteamUGC {
	mustLoad()
	return steamUGC(ptrAPI_SteamUGC())
}

// SteamUGCV021 returns the v021 UGC interface.
func SteamUGCV021() ISteamUGC {
	return SteamUGC()
}

type steamUGC uintptr

func (s steamUGC) GetNumSubscribedItems(includeLocallyDisabled bool) uint32 {
	return ptrAPI_ISteamUGC_GetNumSubscribedItems(uintptr(s), includeLocallyDisabled)
}

func (s steamUGC) GetSubscribedItems(includeLocallyDisabled bool) []PublishedFileId_t {
	count := ptrAPI_ISteamUGC_GetNumSubscribedItems(uintptr(s), includeLocallyDisabled)
	if count == 0 {
		return nil
	}
	items := make([]PublishedFileId_t, count)
	written := ptrAPI_ISteamUGC_GetSubscribedItems(uintptr(s), uintptr(unsafe.Pointer(&items[0])), count, includeLocallyDisabled)
	return items[:written]
}

func (s steamUGC) MarkDownloadedItemAsUnused(publishedFileID PublishedFileId_t) bool {
	if ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused == nil {
		return false
	}
	return ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused(uintptr(s), publishedFileID)
}

func (s steamUGC) GetNumDownloadedItems() uint32 {
	if ptrAPI_ISteamUGC_GetNumDownloadedItems == nil {
		return 0
	}
	return ptrAPI_ISteamUGC_GetNumDownloadedItems(uintptr(s))
}

func (s steamUGC) GetDownloadedItems() []PublishedFileId_t {
	if ptrAPI_ISteamUGC_GetDownloadedItems == nil {
		return nil
	}
	count := s.GetNumDownloadedItems()
	if count == 0 {
		return nil
	}
	items := make([]PublishedFileId_t, count)
	written := ptrAPI_ISteamUGC_GetDownloadedItems(uintptr(s), uintptr(unsafe.Pointer(&items[0])), count)
	return items[:written]
}

func SteamInventory() ISteamInventory {
	mustLoad()
	return steamInventory(ptrAPI_SteamInventory())
}

// SteamInventoryV003 returns the v003 inventory interface.
func SteamInventoryV003() ISteamInventory {
	return SteamInventory()
}

type steamInventory uintptr

func (s steamInventory) GetResultStatus(result SteamInventoryResult_t) EResult {
	return EResult(ptrAPI_ISteamInventory_GetResultStatus(uintptr(s), result))
}

func (s steamInventory) GetResultItems(result SteamInventoryResult_t, outItems []SteamItemDetails) (int, bool) {
	if len(outItems) == 0 {
		return 0, false
	}
	outSize := uint32(len(outItems))
	ok := ptrAPI_ISteamInventory_GetResultItems(uintptr(s), result, uintptr(unsafe.Pointer(&outItems[0])), uintptr(unsafe.Pointer(&outSize)))
	return int(outSize), ok
}

func (s steamInventory) DestroyResult(result SteamInventoryResult_t) {
	ptrAPI_ISteamInventory_DestroyResult(uintptr(s), result)
}

func SteamInput() ISteamInput {
	mustLoad()
	return steamInput(ptrAPI_SteamInput())
}

// SteamInputV006 returns the v006 input interface.
func SteamInputV006() ISteamInput {
	return SteamInput()
}

type steamInput uintptr

func (s steamInput) GetConnectedControllers() []InputHandle_t {
	var handles [_STEAM_INPUT_MAX_COUNT]InputHandle_t
	v := ptrAPI_ISteamInput_GetConnectedControllers(uintptr(s), uintptr(unsafe.Pointer(&handles[0])))
	return handles[:int(v)]
}

func (s steamInput) ConnectedControllers() iter.Seq[InputHandle_t] {
	return func(yield func(InputHandle_t) bool) {
		for _, h := range s.GetConnectedControllers() {
			if !yield(h) {
				return
			}
		}
	}
}

func (s steamInput) GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType {
	v := ptrAPI_ISteamInput_GetInputTypeForHandle(uintptr(s), inputHandle)
	return ESteamInputType(v)
}

func (s steamInput) Init(bExplicitlyCallRunFrame bool) bool {
	return ptrAPI_ISteamInput_Init(uintptr(s), bExplicitlyCallRunFrame)
}

func (s steamInput) Shutdown() {
	ptrAPI_ISteamInput_Shutdown(uintptr(s))
}

func (s steamInput) RunFrame() {
	ptrAPI_ISteamInput_RunFrame(uintptr(s), false)
}

func (s steamInput) EnableDeviceCallbacks() {
	ptrAPI_ISteamInput_EnableDeviceCallbacks(uintptr(s))
}

func (s steamInput) GetActionSetHandle(actionSetName string) InputActionSetHandle_t {
	return ptrAPI_ISteamInput_GetActionSetHandle(uintptr(s), actionSetName)
}

func (s steamInput) ActivateActionSet(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t) {
	ptrAPI_ISteamInput_ActivateActionSet(uintptr(s), inputHandle, actionSetHandle)
}

func (s steamInput) GetCurrentActionSet(inputHandle InputHandle_t) InputActionSetHandle_t {
	return ptrAPI_ISteamInput_GetCurrentActionSet(uintptr(s), inputHandle)
}

func (s steamInput) ActivateActionSetLayer(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t) {
	ptrAPI_ISteamInput_ActivateActionSetLayer(uintptr(s), inputHandle, actionSetHandle)
}

func (s steamInput) DeactivateActionSetLayer(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t) {
	ptrAPI_ISteamInput_DeactivateActionSetLayer(uintptr(s), inputHandle, actionSetHandle)
}

func (s steamInput) DeactivateAllActionSetLayers(inputHandle InputHandle_t) {
	ptrAPI_ISteamInput_DeactivateAllActionSetLayers(uintptr(s), inputHandle)
}

func (s steamInput) GetActiveActionSetLayers(inputHandle InputHandle_t, handles []InputActionSetHandle_t) int {
	if len(handles) == 0 {
		return 0
	}
	return int(ptrAPI_ISteamInput_GetActiveActionSetLayers(uintptr(s), inputHandle, uintptr(unsafe.Pointer(&handles[0]))))
}

func (s steamInput) GetDigitalActionHandle(actionName string) InputDigitalActionHandle_t {
	return ptrAPI_ISteamInput_GetDigitalActionHandle(uintptr(s), actionName)
}

func (s steamInput) GetDigitalActionData(inputHandle InputHandle_t, actionHandle InputDigitalActionHandle_t) InputDigitalActionData {
	data := callInputDigitalActionData(
		ptrAPI_ISteamInput_GetDigitalActionData,
		uintptr(s),
		uint64(inputHandle),
		uint64(actionHandle),
	)
	return InputDigitalActionData{
		State:  data.State,
		Active: data.Active,
	}
}

func (s steamInput) GetDigitalActionOrigins(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t, actionHandle InputDigitalActionHandle_t, origins []EInputActionOrigin) int {
	if len(origins) == 0 {
		return 0
	}
	return int(ptrAPI_ISteamInput_GetDigitalActionOrigins(uintptr(s), inputHandle, actionSetHandle, actionHandle, uintptr(unsafe.Pointer(&origins[0]))))
}

func (s steamInput) GetAnalogActionHandle(actionName string) InputAnalogActionHandle_t {
	return ptrAPI_ISteamInput_GetAnalogActionHandle(uintptr(s), actionName)
}

func (s steamInput) GetAnalogActionData(inputHandle InputHandle_t, actionHandle InputAnalogActionHandle_t) InputAnalogActionData {
	data := callInputAnalogActionData(
		ptrAPI_ISteamInput_GetAnalogActionData,
		uintptr(s),
		uint64(inputHandle),
		uint64(actionHandle),
	)
	return InputAnalogActionData{
		Mode:   EInputSourceMode(data.Mode),
		X:      data.X,
		Y:      data.Y,
		Active: data.Active,
	}
}

func (s steamInput) GetAnalogActionOrigins(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t, actionHandle InputAnalogActionHandle_t, origins []EInputActionOrigin) int {
	if len(origins) == 0 {
		return 0
	}
	return int(ptrAPI_ISteamInput_GetAnalogActionOrigins(uintptr(s), inputHandle, actionSetHandle, actionHandle, uintptr(unsafe.Pointer(&origins[0]))))
}

func (s steamInput) StopAnalogActionMomentum(inputHandle InputHandle_t, actionHandle InputAnalogActionHandle_t) {
	ptrAPI_ISteamInput_StopAnalogActionMomentum(uintptr(s), inputHandle, actionHandle)
}

func (s steamInput) GetMotionData(inputHandle InputHandle_t) InputMotionData {
	data := callInputMotionData(
		ptrAPI_ISteamInput_GetMotionData,
		uintptr(s),
		uint64(inputHandle),
	)
	return InputMotionData{
		RotQuatX:  data.RotQuatX,
		RotQuatY:  data.RotQuatY,
		RotQuatZ:  data.RotQuatZ,
		RotQuatW:  data.RotQuatW,
		PosAccelX: data.PosAccelX,
		PosAccelY: data.PosAccelY,
		PosAccelZ: data.PosAccelZ,
		RotVelX:   data.RotVelX,
		RotVelY:   data.RotVelY,
		RotVelZ:   data.RotVelZ,
	}
}

func (s steamInput) TriggerVibration(inputHandle InputHandle_t, leftSpeed, rightSpeed uint16) {
	ptrAPI_ISteamInput_TriggerVibration(uintptr(s), inputHandle, leftSpeed, rightSpeed)
}

func (s steamInput) TriggerVibrationExtended(inputHandle InputHandle_t, leftSpeed, rightSpeed, leftTriggerSpeed, rightTriggerSpeed uint16) {
	ptrAPI_ISteamInput_TriggerVibrationExtended(uintptr(s), inputHandle, leftSpeed, rightSpeed, leftTriggerSpeed, rightTriggerSpeed)
}

func (s steamInput) TriggerSimpleHapticEvent(inputHandle InputHandle_t, pad ESteamControllerPad, durationMicroSec, offMicroSec, repeat uint16) {
	ptrAPI_ISteamInput_TriggerSimpleHapticEvent(uintptr(s), inputHandle, pad, durationMicroSec, offMicroSec, repeat)
}

func (s steamInput) SetLEDColor(inputHandle InputHandle_t, red, green, blue uint8, flags ESteamInputLEDFlag) {
	ptrAPI_ISteamInput_SetLEDColor(uintptr(s), inputHandle, red, green, blue, flags)
}

func (s steamInput) ShowBindingPanel(inputHandle InputHandle_t) bool {
	return ptrAPI_ISteamInput_ShowBindingPanel(uintptr(s), inputHandle)
}

func (s steamInput) GetControllerForGamepadIndex(index int) InputHandle_t {
	return ptrAPI_ISteamInput_GetControllerForGamepadIndex(uintptr(s), int32(index))
}

func (s steamInput) GetGamepadIndexForController(inputHandle InputHandle_t) int {
	return int(ptrAPI_ISteamInput_GetGamepadIndexForController(uintptr(s), inputHandle))
}

func (s steamInput) GetStringForActionOrigin(origin EInputActionOrigin) string {
	return ptrAPI_ISteamInput_GetStringForActionOrigin(uintptr(s), origin)
}

func (s steamInput) GetGlyphForActionOrigin(origin EInputActionOrigin) string {
	return ptrAPI_ISteamInput_GetGlyphForActionOrigin(uintptr(s), origin)
}

func (s steamInput) GetRemotePlaySessionID(inputHandle InputHandle_t) uint32 {
	return ptrAPI_ISteamInput_GetRemotePlaySessionID(uintptr(s), inputHandle)
}

func SteamRemoteStorage() ISteamRemoteStorage {
	mustLoad()
	return steamRemoteStorage(ptrAPI_SteamRemoteStorage())
}

// SteamRemoteStorageV016 returns the v016 remote storage interface.
func SteamRemoteStorageV016() ISteamRemoteStorage {
	return SteamRemoteStorage()
}

type steamRemoteStorage uintptr

func (s steamRemoteStorage) FileWrite(file string, data []byte) bool {
	return ptrAPI_ISteamRemoteStorage_FileWrite(uintptr(s), file, uintptr(unsafe.Pointer(&data[0])), int32(len(data)))
}

func (s steamRemoteStorage) FileRead(file string, data []byte) int32 {
	return ptrAPI_ISteamRemoteStorage_FileRead(uintptr(s), file, uintptr(unsafe.Pointer(&data[0])), int32(len(data)))
}

func (s steamRemoteStorage) FileDelete(file string) bool {
	return ptrAPI_ISteamRemoteStorage_FileDelete(uintptr(s), file)
}

func (s steamRemoteStorage) GetFileSize(file string) int32 {
	return ptrAPI_ISteamRemoteStorage_GetFileSize(uintptr(s), file)
}

func SteamUser() ISteamUser {
	mustLoad()
	return steamUser(ptrAPI_SteamUser())
}

// SteamUserV023 returns the v023 user interface.
func SteamUserV023() ISteamUser {
	return SteamUser()
}

type steamUser uintptr

func (s steamUser) AdvertiseGame(gameServerSteamID CSteamID, ip uint32, port uint16) {
	ptrAPI_ISteamUser_AdvertiseGame(uintptr(s), gameServerSteamID, ip, port)
}

func (s steamUser) BeginAuthSession(authTicket []byte, steamID CSteamID) EBeginAuthSessionResult {
	return EBeginAuthSessionResult(ptrAPI_ISteamUser_BeginAuthSession(uintptr(s), uintptr(unsafe.Pointer(&authTicket[0])), int32(len(authTicket)), steamID))
}

func (s steamUser) BIsBehindNAT() bool {
	return ptrAPI_ISteamUser_BIsBehindNAT(uintptr(s))
}

func (s steamUser) BIsPhoneIdentifying() bool {
	return ptrAPI_ISteamUser_BIsPhoneIdentifying(uintptr(s))
}

func (s steamUser) BIsPhoneRequiringVerification() bool {
	return ptrAPI_ISteamUser_BIsPhoneRequiringVerification(uintptr(s))
}

func (s steamUser) BIsPhoneVerified() bool {
	return ptrAPI_ISteamUser_BIsPhoneVerified(uintptr(s))
}

func (s steamUser) BIsTwoFactorEnabled() bool {
	return ptrAPI_ISteamUser_BIsTwoFactorEnabled(uintptr(s))
}

func (s steamUser) BLoggedOn() bool {
	return ptrAPI_ISteamUser_BLoggedOn(uintptr(s))
}

func (s steamUser) BSetDurationControlOnlineState(newState EDurationControlOnlineState) bool {
	return ptrAPI_ISteamUser_BSetDurationControlOnlineState(uintptr(s), newState)
}

func (s steamUser) CancelAuthTicket(authTicket HAuthTicket) {
	ptrAPI_ISteamUser_CancelAuthTicket(uintptr(s), authTicket)
}

func (s steamUser) DecompressVoice(compressedData []byte, destBuffer []byte, desiredSampleRate uint32) (bytesWritten uint32, result EVoiceResult) {
	result = EVoiceResult(ptrAPI_ISteamUser_DecompressVoice(uintptr(s), uintptr(unsafe.Pointer(&compressedData[0])), uint32(len(compressedData)), uintptr(unsafe.Pointer(&destBuffer[0])), uint32(len(destBuffer)), uintptr(unsafe.Pointer(&bytesWritten)), desiredSampleRate))
	return
}

func (s steamUser) EndAuthSession(steamID CSteamID) {
	ptrAPI_ISteamUser_EndAuthSession(uintptr(s), steamID)
}

func (s steamUser) GetAuthSessionTicket(authTicket []byte, identityRemote *SteamNetworkingIdentity) (ticket HAuthTicket, size uint32) {
	var remotePtr uintptr
	if identityRemote != nil {
		remotePtr = uintptr(unsafe.Pointer(identityRemote))
	}
	ticket = ptrAPI_ISteamUser_GetAuthSessionTicket(uintptr(s), uintptr(unsafe.Pointer(&authTicket[0])), int32(len(authTicket)), uintptr(unsafe.Pointer(&size)), remotePtr)
	return
}

func (s steamUser) GetAuthTicketForWebApi(identity string) HAuthTicket {
	var handle HAuthTicket
	return ptrAPI_ISteamUser_GetAuthTicketForWebApi(uintptr(s), identity, uintptr(unsafe.Pointer(&handle)))
}

func (s steamUser) GetAvailableVoice() (compressedBytes uint32, uncompressedBytes uint32, result EVoiceResult) {
	result = EVoiceResult(ptrAPI_ISteamUser_GetAvailableVoice(uintptr(s), uintptr(unsafe.Pointer(&compressedBytes)), uintptr(unsafe.Pointer(&uncompressedBytes)), 0))
	return
}

func (s steamUser) GetDurationControl() (control DurationControl, ok bool) {
	ok = ptrAPI_ISteamUser_GetDurationControl(uintptr(s), uintptr(unsafe.Pointer(&control)))
	return
}

func (s steamUser) GetEncryptedAppTicket(ticket []byte) (ticketSize uint32, ok bool) {
	ok = ptrAPI_ISteamUser_GetEncryptedAppTicket(uintptr(s), uintptr(unsafe.Pointer(&ticket[0])), int32(len(ticket)), uintptr(unsafe.Pointer(&ticketSize)))
	return
}

func (s steamUser) GetGameBadgeLevel(series int32, foil bool) int32 {
	return ptrAPI_ISteamUser_GetGameBadgeLevel(uintptr(s), series, foil)
}

func (s steamUser) GetHSteamUser() HSteamUser {
	return ptrAPI_ISteamUser_GetHSteamUser(uintptr(s))
}

func (s steamUser) GetPlayerSteamLevel() int32 {
	return ptrAPI_ISteamUser_GetPlayerSteamLevel(uintptr(s))
}

func (s steamUser) GetSteamID() CSteamID {
	return CSteamID(ptrAPI_ISteamUser_GetSteamID(uintptr(s)))
}

func (s steamUser) GetUserDataFolder() (path string, ok bool) {
	buf := make([]byte, 4096)
	ok = ptrAPI_ISteamUser_GetUserDataFolder(uintptr(s), uintptr(unsafe.Pointer(&buf[0])), int32(len(buf)))
	if i := bytes.IndexByte(buf, 0); i >= 0 {
		path = string(buf[:i])
	} else {
		path = string(buf)
	}
	return
}

func (s steamUser) GetVoice(wantCompressed bool, compressedData []byte, wantUncompressed bool, uncompressedData []byte, desiredSampleRate uint32) (compressedBytes uint32, uncompressedBytes uint32, result EVoiceResult) {
	var compressedPtr uintptr
	if len(compressedData) > 0 {
		compressedPtr = uintptr(unsafe.Pointer(&compressedData[0]))
	}
	var uncompressedPtr uintptr
	if len(uncompressedData) > 0 {
		uncompressedPtr = uintptr(unsafe.Pointer(&uncompressedData[0]))
	}
	result = EVoiceResult(ptrAPI_ISteamUser_GetVoice(uintptr(s), wantCompressed, compressedPtr, uint32(len(compressedData)), uintptr(unsafe.Pointer(&compressedBytes)), wantUncompressed, uncompressedPtr, uint32(len(uncompressedData)), uintptr(unsafe.Pointer(&uncompressedBytes)), desiredSampleRate))
	return
}

func (s steamUser) GetVoiceOptimalSampleRate() uint32 {
	return ptrAPI_ISteamUser_GetVoiceOptimalSampleRate(uintptr(s))
}

func (s steamUser) InitiateGameConnection(authBlob []byte, steamIDGameServer CSteamID, ipServer uint32, portServer uint16, secure bool) int32 {
	return ptrAPI_ISteamUser_InitiateGameConnection(uintptr(s), uintptr(unsafe.Pointer(&authBlob[0])), int32(len(authBlob)), steamIDGameServer, ipServer, portServer, secure)
}

func (s steamUser) RequestEncryptedAppTicket(dataToInclude []byte) SteamAPICall_t {
	return ptrAPI_ISteamUser_RequestEncryptedAppTicket(uintptr(s), uintptr(unsafe.Pointer(&dataToInclude[0])), int32(len(dataToInclude)))
}

func (s steamUser) RequestStoreAuthURL(redirectURL string) SteamAPICall_t {
	return ptrAPI_ISteamUser_RequestStoreAuthURL(uintptr(s), redirectURL)
}

func (s steamUser) StartVoiceRecording() {
	ptrAPI_ISteamUser_StartVoiceRecording(uintptr(s))
}

func (s steamUser) StopVoiceRecording() {
	ptrAPI_ISteamUser_StopVoiceRecording(uintptr(s))
}

func (s steamUser) TerminateGameConnection(ipServer uint32, portServer uint16) {
	ptrAPI_ISteamUser_TerminateGameConnection(uintptr(s), ipServer, portServer)
}

func (s steamUser) TrackAppUsageEvent(gameID CGameID, eventCode int32, extraInfo string) {
	ptrAPI_ISteamUser_TrackAppUsageEvent(uintptr(s), gameID, eventCode, extraInfo)
}

func (s steamUser) UserHasLicenseForApp(steamID CSteamID, appID AppId_t) EUserHasLicenseForAppResult {
	return EUserHasLicenseForAppResult(ptrAPI_ISteamUser_UserHasLicenseForApp(uintptr(s), steamID, appID))
}

func SteamUserStats() ISteamUserStats {
	mustLoad()
	return steamUserStats(ptrAPI_SteamUserStats())
}

// SteamUserStatsV013 returns the v013 user stats interface.
func SteamUserStatsV013() ISteamUserStats {
	return SteamUserStats()
}

type steamUserStats uintptr

func (s steamUserStats) GetAchievement(name string) (achieved, success bool) {
	success = ptrAPI_ISteamUserStats_GetAchievement(uintptr(s), name, uintptr(unsafe.Pointer(&achieved)))
	return
}

func (s steamUserStats) SetAchievement(name string) bool {
	return ptrAPI_ISteamUserStats_SetAchievement(uintptr(s), name)
}

func (s steamUserStats) ClearAchievement(name string) bool {
	return ptrAPI_ISteamUserStats_ClearAchievement(uintptr(s), name)
}

func (s steamUserStats) StoreStats() bool {
	return ptrAPI_ISteamUserStats_StoreStats(uintptr(s))
}

func SteamUtils() ISteamUtils {
	mustLoad()
	return steamUtils(ptrAPI_SteamUtils())
}

// SteamUtilsV010 returns the v010 utils interface.
func SteamUtilsV010() ISteamUtils {
	return SteamUtils()
}

type steamUtils uintptr

func (s steamUtils) GetSecondsSinceAppActive() uint32 {
	return ptrAPI_ISteamUtils_GetSecondsSinceAppActive(uintptr(s))
}

func (s steamUtils) GetSecondsSinceComputerActive() uint32 {
	return ptrAPI_ISteamUtils_GetSecondsSinceComputerActive(uintptr(s))
}

func (s steamUtils) GetConnectedUniverse() EUniverse {
	return EUniverse(ptrAPI_ISteamUtils_GetConnectedUniverse(uintptr(s)))
}

func (s steamUtils) GetServerRealTime() uint32 {
	return ptrAPI_ISteamUtils_GetServerRealTime(uintptr(s))
}

func (s steamUtils) GetIPCountry() string {
	return unique.Make(ptrAPI_ISteamUtils_GetIPCountry(uintptr(s))).Value()
}

func (s steamUtils) GetImageSize(image int) (width, height uint32, ok bool) {
	ok = ptrAPI_ISteamUtils_GetImageSize(uintptr(s), int32(image), uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return
}

func (s steamUtils) GetImageRGBA(image int, dest []byte) bool {
	if len(dest) == 0 {
		return false
	}
	return ptrAPI_ISteamUtils_GetImageRGBA(uintptr(s), int32(image), uintptr(unsafe.Pointer(&dest[0])), int32(len(dest)))
}

func (s steamUtils) GetCurrentBatteryPower() uint8 {
	return ptrAPI_ISteamUtils_GetCurrentBatteryPower(uintptr(s))
}

func (s steamUtils) GetAppID() uint32 {
	return ptrAPI_ISteamUtils_GetAppID(uintptr(s))
}

func (s steamUtils) SetOverlayNotificationPosition(position ENotificationPosition) {
	ptrAPI_ISteamUtils_SetOverlayNotificationPosition(uintptr(s), position)
}

func (s steamUtils) SetOverlayNotificationInset(horizontal, vertical int32) {
	ptrAPI_ISteamUtils_SetOverlayNotificationInset(uintptr(s), horizontal, vertical)
}

func (s steamUtils) IsAPICallCompleted(call SteamAPICall_t) (failed bool, ok bool) {
	ok = ptrAPI_ISteamUtils_IsAPICallCompleted(uintptr(s), call, uintptr(unsafe.Pointer(&failed)))
	return
}

func (s steamUtils) GetAPICallFailureReason(call SteamAPICall_t) ESteamAPICallFailure {
	return ESteamAPICallFailure(ptrAPI_ISteamUtils_GetAPICallFailureReason(uintptr(s), call))
}

func (s steamUtils) GetAPICallResult(call SteamAPICall_t, callback uintptr, callbackSize int32, expectedCallback int32) (failed bool, ok bool) {
	ok = ptrAPI_ISteamUtils_GetAPICallResult(uintptr(s), call, callback, callbackSize, expectedCallback, uintptr(unsafe.Pointer(&failed)))
	return
}

func (s steamUtils) GetIPCCallCount() uint32 {
	return ptrAPI_ISteamUtils_GetIPCCallCount(uintptr(s))
}

func (s steamUtils) IsOverlayEnabled() bool {
	return ptrAPI_ISteamUtils_IsOverlayEnabled(uintptr(s))
}

func (s steamUtils) BOverlayNeedsPresent() bool {
	return ptrAPI_ISteamUtils_BOverlayNeedsPresent(uintptr(s))
}

func (s steamUtils) IsSteamRunningOnSteamDeck() bool {
	return ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck(uintptr(s))
}

func (s steamUtils) ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool {
	return ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput(uintptr(s), keyboardMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight)
}

func SteamNetworkingUtils() ISteamNetworkingUtils {
	mustLoad()
	return steamNetworkingUtils(ptrAPI_SteamNetworkingUtils())
}

// SteamNetworkingUtilsV004 returns the v004 networking utils interface.
func SteamNetworkingUtilsV004() ISteamNetworkingUtils {
	return SteamNetworkingUtils()
}

type steamNetworkingUtils uintptr

func (s steamNetworkingUtils) AllocateMessage(size int) *SteamNetworkingMessage {
	if size <= 0 {
		return nil
	}
	ptr := ptrAPI_ISteamNetworkingUtils_AllocateMessage(uintptr(s), int32(size))
	if ptr == 0 {
		return nil
	}
	return (*SteamNetworkingMessage)(unsafe.Pointer(ptr))
}

func (s steamNetworkingUtils) InitRelayNetworkAccess() {
	ptrAPI_ISteamNetworkingUtils_InitRelayNetworkAccess(uintptr(s))
}

func (s steamNetworkingUtils) GetLocalTimestamp() SteamNetworkingMicroseconds {
	return ptrAPI_ISteamNetworkingUtils_GetLocalTimestamp(uintptr(s))
}

func SteamGameServer() ISteamGameServer {
	mustLoad()
	return steamGameServer(ptrAPI_SteamGameServer())
}

// SteamGameServerV015 returns the v015 game server interface.
func SteamGameServerV015() ISteamGameServer {
	return SteamGameServer()
}

type steamGameServer uintptr

func (s steamGameServer) AssociateWithClan(clanID CSteamID) SteamAPICall_t {
	return ptrAPI_ISteamGameServer_AssociateWithClan(uintptr(s), clanID)
}

func (s steamGameServer) BeginAuthSession(authTicket []byte, steamID CSteamID) EBeginAuthSessionResult {
	return EBeginAuthSessionResult(ptrAPI_ISteamGameServer_BeginAuthSession(uintptr(s), uintptr(unsafe.Pointer(&authTicket[0])), int32(len(authTicket)), steamID))
}

func (s steamGameServer) BLoggedOn() bool {
	return ptrAPI_ISteamGameServer_BLoggedOn(uintptr(s))
}

func (s steamGameServer) BSecure() bool {
	return ptrAPI_ISteamGameServer_BSecure(uintptr(s))
}

func (s steamGameServer) BUpdateUserData(steamIDUser CSteamID, playerName string, score uint32) bool {
	return ptrAPI_ISteamGameServer_BUpdateUserData(uintptr(s), steamIDUser, playerName, score)
}

func (s steamGameServer) CancelAuthTicket(authTicket HAuthTicket) {
	ptrAPI_ISteamGameServer_CancelAuthTicket(uintptr(s), authTicket)
}

func (s steamGameServer) ClearAllKeyValues() {
	ptrAPI_ISteamGameServer_ClearAllKeyValues(uintptr(s))
}

func (s steamGameServer) ComputeNewPlayerCompatibility(steamIDNewPlayer CSteamID, steamIDPlayers []CSteamID, steamIDPlayersInGame []CSteamID, steamIDTeamPlayers []CSteamID) SteamAPICall_t {
	return ptrAPI_ISteamGameServer_ComputeNewPlayerCompatibility(
		uintptr(s),
		steamIDNewPlayer,
		steamIDSlicePtr(steamIDPlayers),
		uint32(len(steamIDPlayers)),
		steamIDSlicePtr(steamIDPlayersInGame),
		uint32(len(steamIDPlayersInGame)),
		steamIDSlicePtr(steamIDTeamPlayers),
		uint32(len(steamIDTeamPlayers)),
	)
}

func (s steamGameServer) CreateUnauthenticatedUserConnection() CSteamID {
	return ptrAPI_ISteamGameServer_CreateUnauthenticatedUserConnection(uintptr(s))
}

func (s steamGameServer) EnableHeartbeats(active bool) {
	if ptrAPI_ISteamGameServer_EnableHeartbeats != nil {
		ptrAPI_ISteamGameServer_EnableHeartbeats(uintptr(s), active)
	}
}

func (s steamGameServer) EndAuthSession(steamID CSteamID) {
	ptrAPI_ISteamGameServer_EndAuthSession(uintptr(s), steamID)
}

func (s steamGameServer) ForceHeartbeat() {
	if ptrAPI_ISteamGameServer_ForceHeartbeat != nil {
		ptrAPI_ISteamGameServer_ForceHeartbeat(uintptr(s))
	}
}

func (s steamGameServer) GetAuthSessionTicket(authTicket []byte) (ticket HAuthTicket, size uint32) {
	ticket = ptrAPI_ISteamGameServer_GetAuthSessionTicket(uintptr(s), uintptr(unsafe.Pointer(&authTicket[0])), int32(len(authTicket)), uintptr(unsafe.Pointer(&size)))
	return
}

func (s steamGameServer) GetGameplayStats() {
	ptrAPI_ISteamGameServer_GetGameplayStats(uintptr(s))
}

func (s steamGameServer) GetNextOutgoingPacket(dest []byte) (size int32, ip uint32, port uint16) {
	size = ptrAPI_ISteamGameServer_GetNextOutgoingPacket(uintptr(s), uintptr(unsafe.Pointer(&dest[0])), int32(len(dest)), uintptr(unsafe.Pointer(&ip)), uintptr(unsafe.Pointer(&port)))
	return
}

func (s steamGameServer) GetPublicIP() uint32 {
	return ptrAPI_ISteamGameServer_GetPublicIP(uintptr(s))
}

func (s steamGameServer) GetServerReputation() SteamAPICall_t {
	return ptrAPI_ISteamGameServer_GetServerReputation(uintptr(s))
}

func (s steamGameServer) GetSteamID() CSteamID {
	return ptrAPI_ISteamGameServer_GetSteamID(uintptr(s))
}

func (s steamGameServer) HandleIncomingPacket(data []byte, ip uint32, port uint16) bool {
	return ptrAPI_ISteamGameServer_HandleIncomingPacket(uintptr(s), uintptr(unsafe.Pointer(&data[0])), int32(len(data)), ip, port)
}

func (s steamGameServer) InitGameServer(ip uint32, steamPort uint16, gamePort uint16, queryPort uint16, serverMode uint32, versionString string) bool {
	if ptrAPI_ISteamGameServer_InitGameServer == nil {
		return false
	}
	return ptrAPI_ISteamGameServer_InitGameServer(uintptr(s), ip, steamPort, gamePort, queryPort, serverMode, versionString)
}

func (s steamGameServer) LogOff() {
	ptrAPI_ISteamGameServer_LogOff(uintptr(s))
}

func (s steamGameServer) LogOn(token string) {
	ptrAPI_ISteamGameServer_LogOn(uintptr(s), token)
}

func (s steamGameServer) LogOnAnonymous() {
	ptrAPI_ISteamGameServer_LogOnAnonymous(uintptr(s))
}

func (s steamGameServer) RequestUserGroupStatus(steamIDUser CSteamID, steamIDGroup CSteamID) bool {
	return ptrAPI_ISteamGameServer_RequestUserGroupStatus(uintptr(s), steamIDUser, steamIDGroup)
}

func (s steamGameServer) SendUserConnectAndAuthenticate(ipClient uint32, authBlob []byte) (steamIDUser CSteamID, ok bool) {
	ok = ptrAPI_ISteamGameServer_SendUserConnectAndAuthenticate(uintptr(s), ipClient, uintptr(unsafe.Pointer(&authBlob[0])), uint32(len(authBlob)), uintptr(unsafe.Pointer(&steamIDUser)))
	return
}

func (s steamGameServer) SendUserDisconnect(steamIDUser CSteamID) {
	ptrAPI_ISteamGameServer_SendUserDisconnect(uintptr(s), steamIDUser)
}

func (s steamGameServer) SetBotPlayerCount(botPlayers int32) {
	ptrAPI_ISteamGameServer_SetBotPlayerCount(uintptr(s), botPlayers)
}

func (s steamGameServer) SetDedicatedServer(dedicated bool) {
	ptrAPI_ISteamGameServer_SetDedicatedServer(uintptr(s), dedicated)
}

func (s steamGameServer) SetGameData(gameData string) {
	ptrAPI_ISteamGameServer_SetGameData(uintptr(s), gameData)
}

func (s steamGameServer) SetGameDescription(description string) {
	ptrAPI_ISteamGameServer_SetGameDescription(uintptr(s), description)
}

func (s steamGameServer) SetGameTags(gameTags string) {
	ptrAPI_ISteamGameServer_SetGameTags(uintptr(s), gameTags)
}

func (s steamGameServer) SetHeartbeatInterval(interval int) {
	if ptrAPI_ISteamGameServer_SetHeartbeatInterval != nil {
		ptrAPI_ISteamGameServer_SetHeartbeatInterval(uintptr(s), int32(interval))
	}
}

func (s steamGameServer) SetKeyValue(key string, value string) {
	ptrAPI_ISteamGameServer_SetKeyValue(uintptr(s), key, value)
}

func (s steamGameServer) SetMapName(mapName string) {
	ptrAPI_ISteamGameServer_SetMapName(uintptr(s), mapName)
}

func (s steamGameServer) SetMaxPlayerCount(playersMax int32) {
	ptrAPI_ISteamGameServer_SetMaxPlayerCount(uintptr(s), playersMax)
}

func (s steamGameServer) SetModDir(modDir string) {
	ptrAPI_ISteamGameServer_SetModDir(uintptr(s), modDir)
}

func (s steamGameServer) SetPasswordProtected(passwordProtected bool) {
	ptrAPI_ISteamGameServer_SetPasswordProtected(uintptr(s), passwordProtected)
}

func (s steamGameServer) SetProduct(product string) {
	ptrAPI_ISteamGameServer_SetProduct(uintptr(s), product)
}

func (s steamGameServer) SetRegion(region string) {
	ptrAPI_ISteamGameServer_SetRegion(uintptr(s), region)
}

func (s steamGameServer) SetServerName(serverName string) {
	ptrAPI_ISteamGameServer_SetServerName(uintptr(s), serverName)
}

func (s steamGameServer) SetSpectatorPort(spectatorPort uint16) {
	ptrAPI_ISteamGameServer_SetSpectatorPort(uintptr(s), spectatorPort)
}

func (s steamGameServer) SetSpectatorServerName(spectatorServerName string) {
	ptrAPI_ISteamGameServer_SetSpectatorServerName(uintptr(s), spectatorServerName)
}

func (s steamGameServer) UserHasLicenseForApp(steamID CSteamID, appID AppId_t) EUserHasLicenseForAppResult {
	return EUserHasLicenseForAppResult(ptrAPI_ISteamGameServer_UserHasLicenseForApp(uintptr(s), steamID, appID))
}

func (s steamGameServer) WasRestartRequested() bool {
	return ptrAPI_ISteamGameServer_WasRestartRequested(uintptr(s))
}

func SteamNetworkingMessages() ISteamNetworkingMessages {
	mustLoad()
	return steamNetworkingMessages(ptrAPI_SteamNetworkingMessages())
}

// SteamNetworkingMessagesV002 returns the v002 networking messages interface.
func SteamNetworkingMessagesV002() ISteamNetworkingMessages {
	return SteamNetworkingMessages()
}

type steamNetworkingMessages uintptr

func (s steamNetworkingMessages) SendMessageToUser(identity *SteamNetworkingIdentity, data []byte, sendFlags SteamNetworkingSendFlags, remoteChannel int) EResult {
	var dataPtr uintptr
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}
	return ptrAPI_ISteamNetworkingMessages_SendMessageToUser(uintptr(s), uintptr(unsafe.Pointer(identity)), dataPtr, uint32(len(data)), int32(sendFlags), int32(remoteChannel))
}

func (s steamNetworkingMessages) ReceiveMessagesOnChannel(channel int, maxMessages int) []*SteamNetworkingMessage {
	if maxMessages <= 0 {
		return nil
	}
	messages := make([]*SteamNetworkingMessage, maxMessages)
	count := ptrAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel(uintptr(s), int32(channel), uintptr(unsafe.Pointer(&messages[0])), int32(maxMessages))
	if count <= 0 {
		return nil
	}
	return messages[:count]
}

func (s steamNetworkingMessages) AcceptSessionWithUser(identity *SteamNetworkingIdentity) bool {
	return ptrAPI_ISteamNetworkingMessages_AcceptSessionWithUser(uintptr(s), uintptr(unsafe.Pointer(identity)))
}

func (s steamNetworkingMessages) CloseSessionWithUser(identity *SteamNetworkingIdentity) bool {
	return ptrAPI_ISteamNetworkingMessages_CloseSessionWithUser(uintptr(s), uintptr(unsafe.Pointer(identity)))
}

func (s steamNetworkingMessages) CloseChannelWithUser(identity *SteamNetworkingIdentity, channel int) bool {
	return ptrAPI_ISteamNetworkingMessages_CloseChannelWithUser(uintptr(s), uintptr(unsafe.Pointer(identity)), int32(channel))
}

func SteamNetworkingSockets() ISteamNetworkingSockets {
	mustLoad()
	return steamNetworkingSockets(ptrAPI_SteamNetworkingSockets())
}

// SteamNetworkingSocketsV012 returns the v012 networking sockets interface.
func SteamNetworkingSocketsV012() ISteamNetworkingSockets {
	return SteamNetworkingSockets()
}

type steamNetworkingSockets uintptr

func (s steamNetworkingSockets) CreateListenSocketIP(localAddress *SteamNetworkingIPAddr, options []SteamNetworkingConfigValue) HSteamListenSocket {
	return ptrAPI_ISteamNetworkingSockets_CreateListenSocketIP(uintptr(s), uintptr(unsafe.Pointer(localAddress)), int32(len(options)), optionsPtr(options))
}

func (s steamNetworkingSockets) CreateListenSocketP2P(localVirtualPort int, options []SteamNetworkingConfigValue) HSteamListenSocket {
	return ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P(uintptr(s), int32(localVirtualPort), int32(len(options)), optionsPtr(options))
}

func (s steamNetworkingSockets) ConnectByIPAddress(address *SteamNetworkingIPAddr, options []SteamNetworkingConfigValue) HSteamNetConnection {
	return ptrAPI_ISteamNetworkingSockets_ConnectByIPAddress(uintptr(s), uintptr(unsafe.Pointer(address)), int32(len(options)), optionsPtr(options))
}

func (s steamNetworkingSockets) ConnectP2P(identity *SteamNetworkingIdentity, remoteVirtualPort int, options []SteamNetworkingConfigValue) HSteamNetConnection {
	return ptrAPI_ISteamNetworkingSockets_ConnectP2P(uintptr(s), uintptr(unsafe.Pointer(identity)), int32(remoteVirtualPort), int32(len(options)), optionsPtr(options))
}

func (s steamNetworkingSockets) AcceptConnection(connection HSteamNetConnection) EResult {
	return ptrAPI_ISteamNetworkingSockets_AcceptConnection(uintptr(s), connection)
}

func (s steamNetworkingSockets) CloseConnection(connection HSteamNetConnection, reason int, debug string, enableLinger bool) bool {
	return ptrAPI_ISteamNetworkingSockets_CloseConnection(uintptr(s), connection, int32(reason), debug, enableLinger)
}

func (s steamNetworkingSockets) CloseListenSocket(socket HSteamListenSocket) bool {
	return ptrAPI_ISteamNetworkingSockets_CloseListenSocket(uintptr(s), socket)
}

func (s steamNetworkingSockets) SendMessageToConnection(connection HSteamNetConnection, data []byte, sendFlags SteamNetworkingSendFlags) (EResult, int64) {
	var dataPtr uintptr
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}
	var messageNumber int64
	result := ptrAPI_ISteamNetworkingSockets_SendMessageToConnection(uintptr(s), connection, dataPtr, uint32(len(data)), int32(sendFlags), uintptr(unsafe.Pointer(&messageNumber)))
	return result, messageNumber
}

func (s steamNetworkingSockets) ReceiveMessagesOnConnection(connection HSteamNetConnection, maxMessages int) []*SteamNetworkingMessage {
	if maxMessages <= 0 {
		return nil
	}
	messages := make([]*SteamNetworkingMessage, maxMessages)
	count := ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection(uintptr(s), connection, uintptr(unsafe.Pointer(&messages[0])), int32(maxMessages))
	if count <= 0 {
		return nil
	}
	return messages[:count]
}

func (s steamNetworkingSockets) CreatePollGroup() HSteamNetPollGroup {
	return ptrAPI_ISteamNetworkingSockets_CreatePollGroup(uintptr(s))
}

func (s steamNetworkingSockets) DestroyPollGroup(group HSteamNetPollGroup) bool {
	return ptrAPI_ISteamNetworkingSockets_DestroyPollGroup(uintptr(s), group)
}

func (s steamNetworkingSockets) SetConnectionPollGroup(connection HSteamNetConnection, group HSteamNetPollGroup) bool {
	return ptrAPI_ISteamNetworkingSockets_SetConnectionPollGroup(uintptr(s), connection, group)
}

func (s steamNetworkingSockets) ReceiveMessagesOnPollGroup(group HSteamNetPollGroup, maxMessages int) []*SteamNetworkingMessage {
	if maxMessages <= 0 {
		return nil
	}
	messages := make([]*SteamNetworkingMessage, maxMessages)
	count := ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup(uintptr(s), group, uintptr(unsafe.Pointer(&messages[0])), int32(maxMessages))
	if count <= 0 {
		return nil
	}
	return messages[:count]
}

func (m *SteamNetworkingMessage) Release() {
	if m == nil || m.ReleaseFunc == 0 {
		return
	}
	purego.SyscallN(m.ReleaseFunc, uintptr(unsafe.Pointer(m)))
}

func optionsPtr(options []SteamNetworkingConfigValue) uintptr {
	if len(options) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&options[0]))
}

func cStringToGo(name []byte) string {
	idx := bytes.IndexByte(name, 0)
	if idx < 0 {
		return string(name)
	}
	return string(name[:idx])
}
