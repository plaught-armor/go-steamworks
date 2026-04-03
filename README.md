# go-steamworks

Go bindings covering all Steamworks API interface families, with typed wrappers for implemented methods and purego/ffi handles for the remaining surfaces.

> [!WARNING]
> 32-bit OSes are not supported.

## Steamworks SDK version

164

> [!NOTE]
> If newer Steamworks SDK releases add or update symbols that are not yet in these bindings, use the [raw symbol access](#raw-symbol-access) method to call them directly.

## Getting started

### Requirements

Before using this library, make sure Steam's redistributable binaries are
available on your runtime machine. This repository no longer ships the
precompiled Steamworks shared libraries; provide them alongside your
application at runtime (for example, next to your executable).

Common locations and filenames:

* Linux (64-bit): `libsteam_api.so`
* macOS: `libsteam_api.dylib`
* Windows (64-bit): `steam_api64.dll`

On Windows, copy the DLL into the working directory:

* `steam_api64.dll` (copy from `redistribution_bin\\win64\\steam_api64.dll` in the SDK)

For local development, ensure `steam_appid.txt` is available next to the
executable (or run Steam with your app ID configured).

### Initialization

The Steamworks client must be running and the API must be initialized before
calling most interfaces. `Load` is optional, but allows you to surface missing
redistributables early.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/badhex/go-steamworks"
)

const appID = 480 // Replace with your own App ID.

func main() {
	if steamworks.RestartAppIfNecessary(appID) {
		os.Exit(0)
	}
	if err := steamworks.Load(); err != nil {
		log.Fatalf("failed to load steamworks: %v", err)
	}
	if err := steamworks.Init(); err != nil {
		log.Fatalf("steamworks.Init failed: %v", err)
	}

	fmt.Printf("SteamID: %v\n", steamworks.SteamUser().GetSteamID())
}
```

### Callback pump

Steamworks expects you to poll callbacks regularly on your main thread.

```go
for running {
	steamworks.RunCallbacks()
	// ...your game loop...
}
```

### Example: language selection

```go
package steamapi

import (
	"github.com/badhex/go-steamworks"
	"golang.org/x/text/language"
)

func SystemLang() language.Tag {
	switch steamworks.SteamApps().GetCurrentGameLanguage() {
	case "english":
		return language.English
	case "japanese":
		return language.Japanese
	}
	return language.Und
}
```

### Example: achievements

```go
if achieved, ok := steamworks.SteamUserStats().GetAchievement("FIRST_WIN"); ok && !achieved {
	steamworks.SteamUserStats().SetAchievement("FIRST_WIN")
	steamworks.SteamUserStats().StoreStats()
}
```

### Example: async call results

```go
call := steamworks.SteamHTTP().CreateHTTPRequest(steamworks.EHTTPMethodGET, "https://example.com")
callHandle, ok := steamworks.SteamHTTP().SendHTTPRequest(call)
if !ok {
	// handle request creation failure
}

type HTTPRequestCompleted struct {
	Request steamworks.HTTPRequestHandle
	Context uint64
	Status  int32
}

// Define the struct to mirror the Steamworks callback payload you expect.
// Use the SDK's callback ID for the expected payload.
result := steamworks.NewCallResult[HTTPRequestCompleted](callHandle, 2101)

if _, failed, err := result.Wait(context.Background(), 0); err == nil && !failed {
	// process response
}
```

## SDK-aligned helpers

This repository ships typed helpers for async call results and manual callback
dispatch, plus additional interface accessors to align with common Steamworks
flows.

* Use `NewCallResult` to await async call results with typed payloads.
* Use `NewCallbackDispatcher` + `RegisterCallback` for manual callback registration and dispatch.
* Use versioned accessors such as `SteamAppsV008()` when you need explicit
  interface versions.

## Build tags and runtime loading

By default, the package expects Steam redistributables to be available on the
runtime library path. You can also opt into embedding redistributables with a
build tag:

* Runtime loading (default): rely on `libsteam_api.so` / `libsteam_api.dylib`
  being in the dynamic linker path or alongside your executable.
* Embedded loading: build with `-tags steamworks_embedded` to embed the SDK
  redistributables and load them from a temporary file at runtime.

Use `STEAMWORKS_LIB_PATH` to point at a custom shared library location when
runtime loading.

## Repository layout

* `gen.go` — code generator for parsing the SDK and building bindings.
* `examples/` — runnable samples for common startup flows.

### Steamworks API coverage and methods

The package now exposes all Steamworks API interface families through friendly
Go accessors. Interfaces are either:

* fully/partially typed wrappers (method-by-method Go bindings), or
* handle-backed wrappers exposing native Go structs with `Ptr() uintptr` and `Valid() bool`.


**General**

* `RestartAppIfNecessary(appID uint32) bool`
* `Init() error`
* `RunCallbacks()`
* `Shutdown()`
* `IsSteamRunning() bool`
* `GetSteamInstallPath() string`
* `ReleaseCurrentThreadMemory()`

**ISteamApps** (`SteamApps() ISteamApps`) — typed wrappers

* `BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, name string, success bool)`
* `BIsSubscribed() bool`
* `BIsLowViolence() bool`
* `BIsCybercafe() bool`
* `BIsVACBanned() bool`
* `BIsDlcInstalled(appID AppId_t) bool`
* `BIsSubscribedApp(appID AppId_t) bool`
* `BIsSubscribedFromFreeWeekend() bool`
* `BIsSubscribedFromFamilySharing() bool`
* `BIsTimedTrial() (allowedSeconds, playedSeconds uint32, ok bool)`
* `BIsAppInstalled(appID AppId_t) bool`
* `GetAvailableGameLanguages() string`
* `GetEarliestPurchaseUnixTime(appID AppId_t) uint32`
* `GetAppInstallDir(appID AppId_t) string`
* `GetCurrentGameLanguage() string`
* `GetDLCCount() int32`
* `GetCurrentBetaName() (string, bool)`
* `GetInstalledDepots(appID AppId_t) []DepotId_t`
* `GetAppOwner() CSteamID`
* `GetLaunchQueryParam(key string) string`
* `GetDlcDownloadProgress(appID AppId_t) (downloaded, total uint64, ok bool)`
* `GetAppBuildId() int32`
* `GetFileDetails(filename string) SteamAPICall_t`
* `GetLaunchCommandLine(bufferSize int) string`
* `GetNumBetas() (total int, available int, private int)`
* `GetBetaInfo(index int) (flags uint32, buildID uint32, lastUpdated uint32, name string, description string, ok bool)`
* `InstallDLC(appID AppId_t)`
* `UninstallDLC(appID AppId_t)`
* `RequestAppProofOfPurchaseKey(appID AppId_t)`
* `RequestAllProofOfPurchaseKeys()`
* `MarkContentCorrupt(missingFilesOnly bool) bool`
* `SetDlcContext(appID AppId_t) bool`
* `SetActiveBeta(name string) bool`

**ISteamAppTicket** (`SteamAppTicket() ISteamAppTicket`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.
* `BSessionRemotePlayTogether(sessionID uint32) bool`
* `GetSessionGuestID(sessionID uint32) CSteamID`
* `GetSmallSessionAvatar(sessionID uint32) int32`
* `GetMediumSessionAvatar(sessionID uint32) int32`
* `GetLargeSessionAvatar(sessionID uint32) int32`

**ISteamClient** (`SteamClient() ISteamClient`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamController** (`SteamController() ISteamController`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamFriends** (`SteamFriends() ISteamFriends`) — typed wrappers

* `GetPersonaName() string`
* `GetPersonaState() EPersonaState`
* `GetFriendCount(flags EFriendFlags) int`
* `GetFriendByIndex(index int, flags EFriendFlags) CSteamID`
* `GetFriendRelationship(friend CSteamID) EFriendRelationship`
* `GetFriendPersonaState(friend CSteamID) EPersonaState`
* `GetFriendPersonaName(friend CSteamID) string`
* `GetFriendPersonaNameHistory(friend CSteamID, index int) string`
* `GetFriendSteamLevel(friend CSteamID) int`
* `GetSmallFriendAvatar(friend CSteamID) int32`
* `GetMediumFriendAvatar(friend CSteamID) int32`
* `GetLargeFriendAvatar(friend CSteamID) int32`
* `SetRichPresence(key, value string) bool`
* `GetFriendGamePlayed(friend CSteamID) (FriendGameInfo, bool)`
  * Returns `FriendGameInfo` mapped from SDK `FriendGameInfo_t` (see field breakdown below).
* `InviteUserToGame(friend CSteamID, connectString string) bool`
* `ActivateGameOverlay(dialog string)`
* `ActivateGameOverlayToUser(dialog string, steamID CSteamID)`
* `ActivateGameOverlayToWebPage(url string, mode EActivateGameOverlayToWebPageMode)`
* `ActivateGameOverlayToStore(appID AppId_t, flag EOverlayToStoreFlag)`
* `ActivateGameOverlayInviteDialog(lobbyID CSteamID)`
* `ActivateGameOverlayInviteDialogConnectString(connectString string)`

Returned structure details:

* `FriendGameInfo` (SDK `FriendGameInfo_t`) fields:
  * `GameID CGameID`
  * `GameIP uint32`
  * `GamePort uint16`
  * `QueryPort uint16`
  * `LobbySteamID CSteamID`

**ISteamGameCoordinator** (`SteamGameCoordinator() ISteamGameCoordinator`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamGameServer** (`SteamGameServer() ISteamGameServer`) — typed wrappers

* `AssociateWithClan(clanID CSteamID) SteamAPICall_t`
* `BeginAuthSession(authTicket []byte, steamID CSteamID) EBeginAuthSessionResult`
* `BLoggedOn() bool`
* `BSecure() bool`
* `BUpdateUserData(steamIDUser CSteamID, playerName string, score uint32) bool`
* `CancelAuthTicket(authTicket HAuthTicket)`
* `ClearAllKeyValues()`
* `ComputeNewPlayerCompatibility(steamIDNewPlayer CSteamID, steamIDPlayers []CSteamID, steamIDPlayersInGame []CSteamID, steamIDTeamPlayers []CSteamID) SteamAPICall_t`
* `CreateUnauthenticatedUserConnection() CSteamID`
* `EnableHeartbeats(active bool)`
* `EndAuthSession(steamID CSteamID)`
* `ForceHeartbeat()`
* `GetAuthSessionTicket(authTicket []byte) (ticket HAuthTicket, size uint32)`
* `GetGameplayStats()`
* `GetNextOutgoingPacket(dest []byte) (size int32, ip uint32, port uint16)`
* `GetPublicIP() uint32`
* `GetServerReputation() SteamAPICall_t`
* `GetSteamID() CSteamID`
* `HandleIncomingPacket(data []byte, ip uint32, port uint16) bool`
* `InitGameServer(ip uint32, steamPort uint16, gamePort uint16, queryPort uint16, serverMode uint32, versionString string) bool`
* `LogOff()`
* `LogOn(token string)`
* `LogOnAnonymous()`
* `RequestUserGroupStatus(steamIDUser CSteamID, steamIDGroup CSteamID) bool`
* `SendUserConnectAndAuthenticate(ipClient uint32, authBlob []byte) (steamIDUser CSteamID, ok bool)`
* `SendUserDisconnect(steamIDUser CSteamID)`
* `SetBotPlayerCount(botPlayers int32)`
* `SetDedicatedServer(dedicated bool)`
* `SetGameData(gameData string)`
* `SetGameDescription(description string)`
* `SetGameTags(gameTags string)`
* `SetHeartbeatInterval(interval int)`
* `SetKeyValue(key string, value string)`
* `SetMapName(mapName string)`
* `SetMaxPlayerCount(playersMax int32)`
* `SetModDir(modDir string)`
* `SetPasswordProtected(passwordProtected bool)`
* `SetProduct(product string)`
* `SetRegion(region string)`
* `SetServerName(serverName string)`
* `SetSpectatorPort(spectatorPort uint16)`
* `SetSpectatorServerName(spectatorServerName string)`
* `UserHasLicenseForApp(steamID CSteamID, appID AppId_t) EUserHasLicenseForAppResult`
* `WasRestartRequested() bool`

**ISteamGameServerStats** (`SteamGameServerStats() ISteamGameServerStats`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamHTMLSurface** (`SteamHTMLSurface() ISteamHTMLSurface`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamHTTP** (`SteamHTTP() ISteamHTTP`) — typed wrappers

* `CreateHTTPRequest(method EHTTPMethod, absoluteURL string) HTTPRequestHandle`
* `SetHTTPRequestHeaderValue(request HTTPRequestHandle, headerName, headerValue string) bool`
* `SendHTTPRequest(request HTTPRequestHandle) (SteamAPICall_t, bool)`
* `GetHTTPResponseBodySize(request HTTPRequestHandle) (uint32, bool)`
* `GetHTTPResponseBodyData(request HTTPRequestHandle, buffer []byte) bool`
* `ReleaseHTTPRequest(request HTTPRequestHandle) bool`

**ISteamInput** (`SteamInput() ISteamInput`) — typed wrappers

* `GetConnectedControllers() []InputHandle_t`
* `GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType`
* `Init(bExplicitlyCallRunFrame bool) bool`
* `Shutdown()`
* `RunFrame()`
* `EnableDeviceCallbacks()`
* `GetActionSetHandle(actionSetName string) InputActionSetHandle_t`
* `ActivateActionSet(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t)`
* `GetCurrentActionSet(inputHandle InputHandle_t) InputActionSetHandle_t`
* `ActivateActionSetLayer(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t)`
* `DeactivateActionSetLayer(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t)`
* `DeactivateAllActionSetLayers(inputHandle InputHandle_t)`
* `GetActiveActionSetLayers(inputHandle InputHandle_t, handles []InputActionSetHandle_t) int`
* `GetDigitalActionHandle(actionName string) InputDigitalActionHandle_t`
* `GetDigitalActionData(inputHandle InputHandle_t, actionHandle InputDigitalActionHandle_t) InputDigitalActionData`
  * Returns `InputDigitalActionData` mapped from SDK `InputDigitalActionData_t`.
* `GetDigitalActionOrigins(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t, actionHandle InputDigitalActionHandle_t, origins []EInputActionOrigin) int`
* `GetAnalogActionHandle(actionName string) InputAnalogActionHandle_t`
* `GetAnalogActionData(inputHandle InputHandle_t, actionHandle InputAnalogActionHandle_t) InputAnalogActionData`
  * Returns `InputAnalogActionData` mapped from SDK `InputAnalogActionData_t`.
* `GetAnalogActionOrigins(inputHandle InputHandle_t, actionSetHandle InputActionSetHandle_t, actionHandle InputAnalogActionHandle_t, origins []EInputActionOrigin) int`
* `StopAnalogActionMomentum(inputHandle InputHandle_t, actionHandle InputAnalogActionHandle_t)`
* `GetMotionData(inputHandle InputHandle_t) InputMotionData`
  * Returns `InputMotionData` mapped from SDK `InputMotionData_t`.
* `TriggerVibration(inputHandle InputHandle_t, leftSpeed, rightSpeed uint16)`
* `TriggerVibrationExtended(inputHandle InputHandle_t, leftSpeed, rightSpeed, leftTriggerSpeed, rightTriggerSpeed uint16)`
* `TriggerSimpleHapticEvent(inputHandle InputHandle_t, pad ESteamControllerPad, durationMicroSec, offMicroSec, repeat uint16)`
* `SetLEDColor(inputHandle InputHandle_t, red, green, blue uint8, flags ESteamInputLEDFlag)`
* `ShowBindingPanel(inputHandle InputHandle_t) bool`
* `GetControllerForGamepadIndex(index int) InputHandle_t`
* `GetGamepadIndexForController(inputHandle InputHandle_t) int`
* `GetStringForActionOrigin(origin EInputActionOrigin) string`
* `GetGlyphForActionOrigin(origin EInputActionOrigin) string`
* `GetRemotePlaySessionID(inputHandle InputHandle_t) uint32`

Returned structure details:

* `InputDigitalActionData` (SDK `InputDigitalActionData_t`) fields:
  * `State bool`
  * `Active bool`
* `InputAnalogActionData` (SDK `InputAnalogActionData_t`) fields:
  * `Mode EInputSourceMode`
  * `X float32`
  * `Y float32`
  * `Active bool`
* `InputMotionData` (SDK `InputMotionData_t`) fields:
  * `RotQuatX float32`, `RotQuatY float32`, `RotQuatZ float32`, `RotQuatW float32`
  * `PosAccelX float32`, `PosAccelY float32`, `PosAccelZ float32`
  * `RotVelX float32`, `RotVelY float32`, `RotVelZ float32`

**ISteamInventory** (`SteamInventory() ISteamInventory`) — typed wrappers

* `GetResultStatus(result SteamInventoryResult_t) EResult`
* `GetResultItems(result SteamInventoryResult_t, outItems []SteamItemDetails) (int, bool)`
  * Populates `outItems` with `SteamItemDetails` entries mapped from SDK `SteamItemDetails_t`.
* `DestroyResult(result SteamInventoryResult_t)`

Returned structure details:

* `SteamItemDetails` (SDK `SteamItemDetails_t`) fields:
  * `ItemID SteamItemInstanceID_t`
  * `Definition SteamItemDef_t`
  * `Quantity uint16`
  * `Flags uint16`

**ISteamMatchmaking** (`SteamMatchmaking() ISteamMatchmaking`) — typed wrappers

* `GetFavoriteGameCount() int`
* `GetFavoriteGame(index int) (FavoriteGame, bool)`
* `AddFavoriteGame(appID AppId_t, ip uint32, connectionPort, queryPort uint16, flags, lastPlayedOnServerTime uint32) int`
* `RemoveFavoriteGame(appID AppId_t, ip uint32, connectionPort, queryPort uint16, flags uint32) bool`
* `RequestLobbyList() SteamAPICall_t`
* `AddRequestLobbyListStringFilter(key, value string, comparisonType ELobbyComparison)`
* `AddRequestLobbyListNumericalFilter(key string, value int, comparisonType ELobbyComparison)`
* `AddRequestLobbyListNearValueFilter(key string, value int)`
* `AddRequestLobbyListFilterSlotsAvailable(slotsAvailable int)`
* `AddRequestLobbyListDistanceFilter(distanceFilter ELobbyDistanceFilter)`
* `AddRequestLobbyListResultCountFilter(maxResults int)`
* `AddRequestLobbyListCompatibleMembersFilter(lobbyID CSteamID)`
* `GetLobbyByIndex(index int) CSteamID`
* `CreateLobby(lobbyType ELobbyType, maxMembers int) SteamAPICall_t`
* `JoinLobby(lobbyID CSteamID) SteamAPICall_t`
* `LeaveLobby(lobbyID CSteamID)`
* `InviteUserToLobby(lobbyID, invitee CSteamID) bool`
* `SetLobbyMemberLimit(lobbyID CSteamID, maxMembers int) bool`
* `GetLobbyMemberLimit(lobbyID CSteamID) int`
* `SetLobbyType(lobbyID CSteamID, lobbyType ELobbyType) bool`
* `SetLobbyJoinable(lobbyID CSteamID, joinable bool) bool`
* `GetLobbyOwner(lobbyID CSteamID) CSteamID`
* `SetLobbyOwner(lobbyID, owner CSteamID) bool`
* `SetLinkedLobby(lobbyID, lobbyDependent CSteamID) bool`
* `GetNumLobbyMembers(lobbyID CSteamID) int`
* `GetLobbyMemberByIndex(lobbyID CSteamID, memberIndex int) CSteamID`
* `SetLobbyData(lobbyID CSteamID, key, value string) bool`
* `GetLobbyData(lobbyID CSteamID, key string) string`
* `DeleteLobbyData(lobbyID CSteamID, key string) bool`
* `GetLobbyDataCount(lobbyID CSteamID) int`
* `GetLobbyDataByIndex(lobbyID CSteamID, lobbyDataIndex int) (key, value string, ok bool)`
* `SetLobbyMemberData(lobbyID CSteamID, key, value string)`
* `GetLobbyMemberData(lobbyID, user CSteamID, key string) string`
* `SendLobbyChatMsg(lobbyID CSteamID, msgBody []byte) bool`
* `GetLobbyChatEntry(lobbyID CSteamID, chatID int, data []byte) (user CSteamID, entryType EChatEntryType, bytesCopied int)`
* `RequestLobbyData(lobbyID CSteamID) bool`
* `SetLobbyGameServer(lobbyID CSteamID, ip uint32, port uint16, server CSteamID)`
* `GetLobbyGameServer(lobbyID CSteamID) (ip uint32, port uint16, server CSteamID, ok bool)`
* `CheckForPSNGameBootInvite(lobbyID *CSteamID) bool`

**ISteamMatchmakingServers** (`SteamMatchmakingServers() ISteamMatchmakingServers`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.
* `RequestFavoritesServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest`
* `RequestFriendsServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest`
* `RequestHistoryServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest`
* `RequestInternetServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest`
* `RequestLANServerList(appID AppId_t, response uintptr) HServerListRequest`
* `RequestSpectatorServerList(appID AppId_t, filters []uintptr, response uintptr) HServerListRequest`
* `ReleaseRequest(request HServerListRequest)`
* `GetServerDetails(request HServerListRequest, server int) uintptr`
* `CancelQuery(request HServerListRequest)`
* `RefreshQuery(request HServerListRequest)`
* `IsRefreshing(request HServerListRequest) bool`
* `GetServerCount(request HServerListRequest) int`
* `RefreshServer(request HServerListRequest, server int)`
* `PingServer(ip uint32, port uint16, response uintptr) HServerQuery`
* `PlayerDetails(ip uint32, port uint16, response uintptr) HServerQuery`
* `ServerRules(ip uint32, port uint16, response uintptr) HServerQuery`
* `CancelServerQuery(query HServerQuery)`

**ISteamMusic** (`SteamMusic() ISteamMusic`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamNetworking** (`SteamNetworking() ISteamNetworking`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamNetworkingMessages** (`SteamNetworkingMessages() ISteamNetworkingMessages`) — typed wrappers

* `SendMessageToUser(identity *SteamNetworkingIdentity, data []byte, sendFlags SteamNetworkingSendFlags, remoteChannel int) EResult`
* `ReceiveMessagesOnChannel(channel int, maxMessages int) []*SteamNetworkingMessage`
  * Returns a slice of `*SteamNetworkingMessage` wrappers over SDK `SteamNetworkingMessage_t`.
* `AcceptSessionWithUser(identity *SteamNetworkingIdentity) bool`
* `CloseSessionWithUser(identity *SteamNetworkingIdentity) bool`
* `CloseChannelWithUser(identity *SteamNetworkingIdentity, channel int) bool`

Returned structure details:

* `SteamNetworkingIdentity` fields:
  * `IdentityType int32`
  * `Reserved [3]int32`
  * `Data [128]byte`
* `SteamNetworkingMessage` (SDK `SteamNetworkingMessage_t`) pointer wrapper:
  * `Data uintptr`
  * `Size int32`
  * `Conn HSteamNetConnection`
  * `IdentityPeer SteamNetworkingIdentity`
  * `ConnUserData int64`
  * `TimeReceived int64`
  * `MessageNumber int64`
  * `ReleaseFunc uintptr` (invoked by `Release()`)

**ISteamNetworkingSockets** (`SteamNetworkingSockets() ISteamNetworkingSockets`) — typed wrappers

* `CreateListenSocketIP(localAddress *SteamNetworkingIPAddr, options []SteamNetworkingConfigValue) HSteamListenSocket`
* `CreateListenSocketP2P(localVirtualPort int, options []SteamNetworkingConfigValue) HSteamListenSocket`
* `ConnectByIPAddress(address *SteamNetworkingIPAddr, options []SteamNetworkingConfigValue) HSteamNetConnection`
* `ConnectP2P(identity *SteamNetworkingIdentity, remoteVirtualPort int, options []SteamNetworkingConfigValue) HSteamNetConnection`
* `AcceptConnection(connection HSteamNetConnection) EResult`
* `CloseConnection(connection HSteamNetConnection, reason int, debug string, enableLinger bool) bool`
* `CloseListenSocket(socket HSteamListenSocket) bool`
* `SendMessageToConnection(connection HSteamNetConnection, data []byte, sendFlags SteamNetworkingSendFlags) (EResult, int64)`
* `ReceiveMessagesOnConnection(connection HSteamNetConnection, maxMessages int) []*SteamNetworkingMessage`
  * Returns a slice of `*SteamNetworkingMessage` wrappers over SDK `SteamNetworkingMessage_t`.
* `CreatePollGroup() HSteamNetPollGroup`
* `DestroyPollGroup(group HSteamNetPollGroup) bool`
* `SetConnectionPollGroup(connection HSteamNetConnection, group HSteamNetPollGroup) bool`
* `ReceiveMessagesOnPollGroup(group HSteamNetPollGroup, maxMessages int) []*SteamNetworkingMessage`
  * Returns a slice of `*SteamNetworkingMessage` wrappers over SDK `SteamNetworkingMessage_t`.

Returned structure details:

* `SteamNetworkingIPAddr` fields:
  * `IP [16]byte`
  * `Port uint16`
* `SteamNetworkingIdentity` fields:
  * `IdentityType int32`
  * `Reserved [3]int32`
  * `Data [128]byte`
* `SteamNetworkingMessage` (SDK `SteamNetworkingMessage_t`) pointer wrapper:
  * `Data uintptr`
  * `Size int32`
  * `Conn HSteamNetConnection`
  * `IdentityPeer SteamNetworkingIdentity`
  * `ConnUserData int64`
  * `TimeReceived int64`
  * `MessageNumber int64`
  * `ReleaseFunc uintptr` (invoked by `Release()`)

**ISteamNetworkingUtils** (`SteamNetworkingUtils() ISteamNetworkingUtils`) — typed wrappers

* `AllocateMessage(size int) *SteamNetworkingMessage`
  * Returns a `*SteamNetworkingMessage` wrapper over SDK `SteamNetworkingMessage_t`.
* `InitRelayNetworkAccess()`
* `GetLocalTimestamp() SteamNetworkingMicroseconds`

**ISteamRemotePlay** (`SteamRemotePlay() ISteamRemotePlay`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamRemoteStorage** (`SteamRemoteStorage() ISteamRemoteStorage`) — typed wrappers

* `FileWrite(file string, data []byte) bool`
* `FileRead(file string, data []byte) int32`
* `FileDelete(file string) bool`
* `GetFileSize(file string) int32`

**ISteamScreenshots** (`SteamScreenshots() ISteamScreenshots`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamTimeline** (`SteamTimeline() ISteamTimeline`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**ISteamUGC** (`SteamUGC() ISteamUGC`) — typed wrappers

* `GetNumSubscribedItems(includeLocallyDisabled bool) uint32`
* `GetSubscribedItems(includeLocallyDisabled bool) []PublishedFileId_t`
* `MarkDownloadedItemAsUnused(publishedFileID PublishedFileId_t) bool`
* `GetNumDownloadedItems() uint32`
* `GetDownloadedItems() []PublishedFileId_t`

**ISteamUser** (`SteamUser() ISteamUser`) — typed wrappers

* `AdvertiseGame(gameServerSteamID CSteamID, ip uint32, port uint16)`
* `BeginAuthSession(authTicket []byte, steamID CSteamID) EBeginAuthSessionResult`
* `BIsBehindNAT() bool`
* `BIsPhoneIdentifying() bool`
* `BIsPhoneRequiringVerification() bool`
* `BIsPhoneVerified() bool`
* `BIsTwoFactorEnabled() bool`
* `BLoggedOn() bool`
* `BSetDurationControlOnlineState(newState EDurationControlOnlineState) bool`
* `CancelAuthTicket(authTicket HAuthTicket)`
* `DecompressVoice(compressedData []byte, destBuffer []byte, desiredSampleRate uint32) (bytesWritten uint32, result EVoiceResult)`
* `EndAuthSession(steamID CSteamID)`
* `GetAuthSessionTicket(authTicket []byte, identityRemote *SteamNetworkingIdentity) (ticket HAuthTicket, size uint32)`
* `GetAuthTicketForWebApi(identity string) HAuthTicket`
* `GetAvailableVoice() (compressedBytes uint32, uncompressedBytes uint32, result EVoiceResult)`
* `GetDurationControl() (control DurationControl, ok bool)`
* `GetEncryptedAppTicket(ticket []byte) (ticketSize uint32, ok bool)`
* `GetGameBadgeLevel(series int32, foil bool) int32`
* `GetHSteamUser() HSteamUser`
* `GetPlayerSteamLevel() int32`
* `GetSteamID() CSteamID`
* `GetUserDataFolder() (path string, ok bool)`
* `GetVoice(wantCompressed bool, compressedData []byte, wantUncompressed bool, uncompressedData []byte, desiredSampleRate uint32) (compressedBytes uint32, uncompressedBytes uint32, result EVoiceResult)`
* `GetVoiceOptimalSampleRate() uint32`
* `InitiateGameConnection(authBlob []byte, steamIDGameServer CSteamID, ipServer uint32, portServer uint16, secure bool) int32`
* `RequestEncryptedAppTicket(dataToInclude []byte) SteamAPICall_t`
* `RequestStoreAuthURL(redirectURL string) SteamAPICall_t`
* `StartVoiceRecording()`
* `StopVoiceRecording()`
* `TerminateGameConnection(ipServer uint32, portServer uint16)`
* `TrackAppUsageEvent(gameID CGameID, eventCode int32, extraInfo string)`
* `UserHasLicenseForApp(steamID CSteamID, appID AppId_t) EUserHasLicenseForAppResult`

**ISteamUserStats** (`SteamUserStats() ISteamUserStats`) — typed wrappers

* `GetAchievement(name string) (achieved, success bool)`
* `SetAchievement(name string) bool`
* `ClearAchievement(name string) bool`
* `StoreStats() bool`

**ISteamUtils** (`SteamUtils() ISteamUtils`) — typed wrappers

* `GetSecondsSinceAppActive() uint32`
* `GetSecondsSinceComputerActive() uint32`
* `GetConnectedUniverse() EUniverse`
* `GetServerRealTime() uint32`
* `GetIPCountry() string`
* `GetImageSize(image int) (width, height uint32, ok bool)`
* `GetImageRGBA(image int, dest []byte) bool`
* `GetCurrentBatteryPower() uint8`
* `GetAppID() uint32`
* `IsOverlayEnabled() bool`
* `BOverlayNeedsPresent() bool`
* `IsSteamRunningOnSteamDeck() bool`
* `SetOverlayNotificationPosition(position ENotificationPosition)`
* `SetOverlayNotificationInset(horizontal, vertical int32)`
* `IsAPICallCompleted(call SteamAPICall_t) (failed bool, ok bool)`
* `GetAPICallFailureReason(call SteamAPICall_t) ESteamAPICallFailure`
* `GetAPICallResult(call SteamAPICall_t, callback uintptr, callbackSize int32, expectedCallback int32) (failed bool, ok bool)`
* `GetIPCCallCount() uint32`
* `ShowFloatingGamepadTextInput(...) bool`

**ISteamVideo** (`SteamVideo() ISteamVideo`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**SteamEncryptedAppTicket** — typed utility wrappers

* `SteamEncryptedAppTicketBDecryptTicket(ticket, decrypted, key []byte) (decryptedSize uint32, ok bool)`
* `SteamEncryptedAppTicketBIsTicketForApp(decryptedTicket []byte, appID AppId_t) bool`
* `SteamEncryptedAppTicketGetTicketIssueTime(decryptedTicket []byte) uint32`
* `SteamEncryptedAppTicketGetTicketSteamID(decryptedTicket []byte) (CSteamID, bool)`

**steam_api** (`SteamAPIClient() ISteamAPIClient`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

**steam_gameserver** (`SteamAPIGameServer() ISteamAPIGameServer`) — handle-backed

* Returned wrapper struct shape: `{ ptr uintptr }` with methods `Ptr() uintptr` and `Valid() bool`.

### Additional raw symbol helpers (purego/ffi)

Raw helpers return concrete Go wrapper structs per interface family; use `Ptr()` for FFI call entry and `Valid()` before invocation.

Interface raw helpers resolve the `SteamAPI_*` exported factory symbol and invoke it (zero-arg) to obtain the actual ISteam* instance pointer. If `purego.Dlsym` cannot resolve the factory symbol, the package attempts an ffi-based fallback lookup via `github.com/jupiterrider/ffi`.

* `SteamAppTicketRaw() ISteamAppTicket` (ISteamAppTicket)
* `SteamClientRaw() ISteamClient` (ISteamClient)
* `SteamControllerRaw() ISteamController` (ISteamController)
* `SteamGameCoordinatorRaw() ISteamGameCoordinator` (ISteamGameCoordinator)
* `SteamGameServerStatsRaw() ISteamGameServerStats` (ISteamGameServerStats)
* `SteamHTMLSurfaceRaw() ISteamHTMLSurface` (ISteamHTMLSurface)
* `SteamMatchmakingServersRaw() ISteamMatchmakingServers` (ISteamMatchmakingServers)
* `SteamMusicRaw() ISteamMusic` (ISteamMusic)
* `SteamNetworkingRaw() ISteamNetworking` (ISteamNetworking)
* `SteamRemotePlayRaw() ISteamRemotePlay` (ISteamRemotePlay)
* `SteamScreenshotsRaw() ISteamScreenshots` (ISteamScreenshots)
* `SteamTimelineRaw() ISteamTimeline` (ISteamTimeline)
* `SteamVideoRaw() ISteamVideo` (ISteamVideo)
* `SteamAPIClientRaw() ISteamAPIClient` (`steam_api` foundation; wrapped by `SteamAPIClient()`)
* `SteamAPIGameServerRaw() ISteamAPIGameServer` (`steam_gameserver` foundation; wrapped by `SteamAPIGameServer()`)

Encrypted ticket utilities are also exposed with direct symbol wrappers:

* `SteamEncryptedAppTicketBDecryptTicket(ticket, decrypted, key []byte) (decryptedSize uint32, ok bool)`
* `SteamEncryptedAppTicketBIsTicketForApp(decryptedTicket []byte, appID AppId_t) bool`
* `SteamEncryptedAppTicketGetTicketIssueTime(decryptedTicket []byte) uint32`
* `SteamEncryptedAppTicketGetTicketSteamID(decryptedTicket []byte) (CSteamID, bool)`

### Raw symbol access

To access newer or unsupported Steamworks SDK methods, you can call raw symbols
directly:

```go
// Look up a symbol and call it directly (advanced usage).
ptr, err := steamworks.LookupSymbol("SteamAPI_ISteamFriends_GetPersonaName")
if err != nil {
	panic(err)
}
result := steamworks.CallSymbolPtr(ptr)
_ = result
```

Or use `CallSymbol` to combine lookup + call:

```go
result, err := steamworks.CallSymbol("SteamAPI_SteamApps_v008")
if err != nil {
	panic(err)
}
_ = result
```

## License

All the source code files are licensed under Apache License 2.0.

## Resources

 * [Steamworks SDK](https://partner.steamgames.com/doc/sdk)
