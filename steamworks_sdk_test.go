// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The go-steamworks Authors

package steamworks

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"
)

var (
	initOnce   sync.Once
	initResult initStatus
	libHandle  uintptr
)

func setupSteamAPI(t *testing.T) initStatus {
	t.Helper()
	initOnce.Do(func() {
		libHandle = loadSDKLibrary(t)
		registerFunctions(libHandle)
		registerInputStructReturns(libHandle)
		initResult = initSteamAPI(t)
	})
	if !initResult.ok {
		t.Fatalf("Steam API failed to initialize: %s", initResult.message)
	}
	return initResult
}

func TestMain(m *testing.M) {
	code := m.Run()

	// Shutdown if we initialized
	if initResult.ok && ptrAPI_Shutdown != nil {
		fmt.Println("Shutting down Steam API...")
		ptrAPI_Shutdown()
	}

	os.Exit(code)
}

func sdkLibraryPath() (string, error) {
	envPath := os.Getenv(steamworksLibEnv)
	if envPath == "" {
		return "", fmt.Errorf("%s must be set to the Steamworks SDK zip URL or library path", steamworksLibEnv)
	}
	if isRemoteLocation(envPath) {
		return downloadLibrary(envPath)
	}
	if _, err := os.Stat(envPath); err != nil {
		return "", fmt.Errorf("%s points to missing file: %w", steamworksLibEnv, err)
	}
	isZip, err := isZipArchive(envPath)
	if err != nil {
		return "", err
	}
	if isZip {
		tmpDir, err := os.MkdirTemp("", "steamworks-sdk-*")
		if err != nil {
			return "", err
		}
		return extractSDKZip(envPath, tmpDir)
	}
	return envPath, nil
}

func isRemoteLocation(path string) bool {
	parsed, err := url.Parse(path)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func downloadLibrary(location string) (string, error) {
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Get(location)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	tmpDir, err := os.MkdirTemp("", "steamworks-sdk-*")
	if err != nil {
		return "", err
	}
	zipPath := filepath.Join(tmpDir, "steamworks_sdk.zip")
	zipFile, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(zipFile, resp.Body); err != nil {
		zipFile.Close()
		return "", err
	}
	if err := zipFile.Close(); err != nil {
		return "", err
	}

	return extractSDKZip(zipPath, tmpDir)
}

func isZipArchive(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	header := make([]byte, 4)
	if _, err := io.ReadFull(file, header); err != nil {
		return false, err
	}

	return bytes.HasPrefix(header, []byte("PK")), nil
}

func extractSDKZip(zipPath, destDir string) (string, error) {
	entryName, err := sdkLibraryEntry()
	if err != nil {
		return "", err
	}
	return extractZipFile(zipPath, entryName, destDir)
}

func sdkLibraryEntry() (string, error) {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "sdk/redistributable_bin/linux64/libsteam_api.so", nil
		case "386":
			return "sdk/redistributable_bin/linux32/libsteam_api.so", nil
		}
	case "darwin":
		return "sdk/redistributable_bin/osx/libsteam_api.dylib", nil
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return "sdk/redistributable_bin/win64/steam_api64.dll", nil
		case "386":
			return "sdk/redistributable_bin/steam_api.dll", nil
		}
	}
	return "", fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}

func extractZipFile(zipPath, entryName, destDir string) (string, error) {
	zipFile, err := os.Open(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	info, err := zipFile.Stat()
	if err != nil {
		return "", err
	}

	reader, err := zip.NewReader(zipFile, info.Size())
	if err != nil {
		return "", err
	}

	for _, file := range reader.File {
		if file.Name != entryName {
			continue
		}

		src, err := file.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()

		filename := filepath.Base(entryName)
		if filename == "." || filename == "/" || filename == "" {
			filename = "libsteam_api"
		}
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		outPath := filepath.Join(destDir, filename)
		outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(outFile, src); err != nil {
			outFile.Close()
			return "", err
		}

		if err := outFile.Close(); err != nil {
			return "", err
		}

		return outPath, nil
	}

	return "", fmt.Errorf("sdk library %s not found in archive", entryName)
}

func TestSDKSymbolResolution(t *testing.T) {
	initOnce.Do(func() {
		libHandle = loadSDKLibrary(t)
		registerFunctions(libHandle)
		registerInputStructReturns(libHandle)
		initResult = initSteamAPI(t)
	})

	expectedMissing := map[string]struct{}{
		flatAPI_ISteamInput_GetGlyphForActionOrigin: {},
		flatAPI_SteamApps: {},
	}
	optionalSymbols := map[string]struct{}{
		flatAPI_ISteamMatchmaking_CheckForPSNGameBootInvite: {},
		flatAPI_ISteamGameServer_EnableHeartbeats:           {},
		flatAPI_ISteamGameServer_ForceHeartbeat:             {},
		flatAPI_ISteamGameServer_InitGameServer:             {},
		flatAPI_ISteamGameServer_SetHeartbeatInterval:       {},
		flatAPI_ISteamUGC_MarkDownloadedItemAsUnused:        {},
		flatAPI_ISteamUGC_GetNumDownloadedItems:             {},
		flatAPI_ISteamUGC_GetDownloadedItems:                {},
		flatAPI_SteamRemotePlay:                             {},
		flatAPI_ISteamRemotePlay_BSessionRemotePlayTogether: {},
		flatAPI_ISteamRemotePlay_GetSessionGuestID:          {},
		flatAPI_ISteamRemotePlay_GetSmallSessionAvatar:      {},
		flatAPI_ISteamRemotePlay_GetMediumSessionAvatar:     {},
		flatAPI_ISteamRemotePlay_GetLargeSessionAvatar:      {},
	}
	for _, symbol := range allFlatAPISymbols() {
		t.Run(symbol, func(t *testing.T) {
			ptr, err := lookupSymbolAddr(libHandle, symbol)
			if err != nil {
				if _, ok := expectedMissing[symbol]; ok {
					t.Logf("expected missing symbol: %s", symbol)
					return
				}
				if _, ok := optionalSymbols[symbol]; ok {
					t.Logf("optional symbol not exported by this SDK: %s", symbol)
					return
				}
				t.Fatalf("lookupSymbolAddr(%s): %v", symbol, err)
			}
			if _, ok := expectedMissing[symbol]; ok {
				t.Fatalf("expected missing symbol %s, but it was present", symbol)
			}
			if ptr == 0 {
				t.Fatalf("lookupSymbolAddr(%s) returned 0", symbol)
			}
			t.Logf("resolved symbol %s: 0x%x", symbol, ptr)
		})
	}
}

func TestSDKCallSymbol(t *testing.T) {
	initOnce.Do(func() {
		libHandle = loadSDKLibrary(t)
		registerFunctions(libHandle)
		registerInputStructReturns(libHandle)
		initResult = initSteamAPI(t)
	})
	ptr, err := lookupSymbolAddr(libHandle, flatAPI_IsSteamRunning)
	if err != nil {
		t.Fatalf("lookupSymbolAddr(%s): %v", flatAPI_IsSteamRunning, err)
	}
	result := CallSymbolPtr(ptr)
	if result == 0 {
		t.Log("Steam is not running")
	} else {
		t.Logf("Steam is running (result: %d)", result)
	}
}

func TestSDKFunctionSignatures(t *testing.T) {
	initOnce.Do(func() {
		libHandle = loadSDKLibrary(t)
		registerFunctions(libHandle)
		registerInputStructReturns(libHandle)
		initResult = initSteamAPI(t)
	})
	actuals := make(map[string]interface{})
	for _, item := range allRegisteredFunctions() {
		actuals[item.name] = item.value
	}

	for _, expectation := range signatureExpectations() {
		t.Run(expectation.name, func(t *testing.T) {
			actual, ok := actuals[expectation.name]
			if !ok {
				t.Fatalf("missing registered function %s", expectation.name)
			}
			assertSignature(t, expectation.name, actual, expectation.expected)
			t.Logf("signature ok: %s is %T", expectation.name, actual)
		})
	}
}

func TestSDKFunctionExecution(t *testing.T) {
	initState := setupSteamAPI(t)
	interfacePtrs := interfacePointers()

	expectedMissing := map[string]struct{}{
		"ptrAPI_ISteamInput_GetGlyphForActionOrigin":         {},
		"ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite": {},
		"ptrAPI_ISteamGameServer_EnableHeartbeats":           {},
		"ptrAPI_ISteamGameServer_ForceHeartbeat":             {},
		"ptrAPI_ISteamGameServer_InitGameServer":             {},
		"ptrAPI_ISteamGameServer_SetHeartbeatInterval":       {},
		"ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused":        {},
		"ptrAPI_ISteamUGC_GetNumDownloadedItems":             {},
		"ptrAPI_ISteamUGC_GetDownloadedItems":                {},
		"ptrAPI_SteamRemotePlay":                             {},
		"ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether": {},
		"ptrAPI_ISteamRemotePlay_GetSessionGuestID":          {},
		"ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar":      {},
		"ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar":     {},
		"ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar":      {},
	}

	actuals := make(map[string]interface{})
	for _, item := range allRegisteredFunctions() {
		actuals[item.name] = item.value
	}

	for _, expectation := range signatureExpectations() {
		t.Run(expectation.name, func(t *testing.T) {
			actual, ok := actuals[expectation.name]
			if !ok {
				t.Fatalf("missing registered function %s", expectation.name)
			}

			if _, ok := expectation.expected.(uintptr); ok {
				ptr, ok := actual.(uintptr)
				if !ok {
					t.Fatalf("%s has type %T, want uintptr", expectation.name, actual)
				}
				if ptr == 0 {
					t.Fatalf("%s ffi pointer is 0", expectation.name)
				}
				runFFIInputCall(t, expectation.name, ptr, interfacePtrs["ISteamInput"])
				t.Logf("ffi call ok for %s: 0x%x", expectation.name, ptr)
				return
			}

			value := reflect.ValueOf(actual)
			if value.Kind() != reflect.Func {
				t.Fatalf("%s has type %T, want func", expectation.name, actual)
			}
			if value.IsNil() {
				if _, ok := expectedMissing[expectation.name]; ok {
					t.Logf("expected missing function: %s", expectation.name)
					return
				}
				t.Fatalf("%s is nil after registration", expectation.name)
			}

			if expectation.name == "ptrAPI_Shutdown" {
				t.Logf("skipping %s during main test execution", expectation.name)
				return
			}

			validateSignatureTypes(t, expectation.name, value.Type())
			callRegisteredFunction(t, expectation.name, value, interfacePtrs, initState)
		})
	}
}

func validateSignatureTypes(t *testing.T, name string, fnType reflect.Type) {
	t.Helper()

	for i := 0; i < fnType.NumIn(); i++ {
		arg := fnType.In(i)
		if !isSupportedType(arg) {
			t.Fatalf("%s arg[%d] unsupported type: %s", name, i, arg)
		}
	}
	for i := 0; i < fnType.NumOut(); i++ {
		out := fnType.Out(i)
		if !isSupportedType(out) {
			t.Fatalf("%s return[%d] unsupported type: %s", name, i, out)
		}
	}
}

type initStatus struct {
	ok      bool
	message string
}

func initSteamAPI(t *testing.T) initStatus {
	t.Helper()

	if appID := os.Getenv("STEAM_APPID"); appID != "" {
		if err := os.WriteFile("steam_appid.txt", []byte(appID), 0644); err != nil {
			t.Logf("failed to write steam_appid.txt: %v", err)
		} else {
			t.Cleanup(func() {
				_ = os.Remove("steam_appid.txt")
			})
		}
	}

	var msg steamErrMsg
	result := ptrAPI_InitFlat(uintptr(unsafe.Pointer(&msg)))
	if result != ESteamAPIInitResult_OK {
		message := fmt.Sprintf("InitFlat failed (%d): %s", result, msg.String())
		t.Log(message)
		return initStatus{ok: false, message: message}
	}
	t.Log("InitFlat succeeded")
	return initStatus{ok: true}
}

func steamAppsInterfacePtr() uintptr {
	if ptrAPI_SteamApps != nil {
		if ptr := ptrAPI_SteamApps(); ptr != 0 {
			return ptr
		}
	}
	return resolveInterfaceFactory("SteamAPI_SteamApps_v009", "SteamAPI_SteamApps_v008", "SteamAPI_SteamApps")
}

func interfacePointers() map[string]uintptr {
	return map[string]uintptr{
		"ISteamApps":               steamAppsInterfacePtr(),
		"ISteamFriends":            ptrAPI_SteamFriends(),
		"ISteamMatchmaking":        ptrAPI_SteamMatchmaking(),
		"ISteamMatchmakingServers": SteamMatchmakingServersRaw().Ptr(),
		"ISteamHTTP":               ptrAPI_SteamHTTP(),
		"ISteamUGC":                ptrAPI_SteamUGC(),
		"ISteamInventory":          ptrAPI_SteamInventory(),
		"ISteamInput":              ptrAPI_SteamInput(),
		"ISteamRemotePlay":         SteamRemotePlay().Ptr(),
		"ISteamRemoteStorage":      ptrAPI_SteamRemoteStorage(),
		"ISteamUser":               ptrAPI_SteamUser(),
		"ISteamUserStats":          ptrAPI_SteamUserStats(),
		"ISteamUtils":              ptrAPI_SteamUtils(),
		"ISteamNetworkingUtils":    ptrAPI_SteamNetworkingUtils(),
		"ISteamNetworkingMessages": ptrAPI_SteamNetworkingMessages(),
		"ISteamNetworkingSockets":  ptrAPI_SteamNetworkingSockets(),
		"ISteamGameServer":         ptrAPI_SteamGameServer(),
	}
}

func runFFIInputCall(t *testing.T, name string, ptr uintptr, steamInput uintptr) {
	t.Helper()

	switch name {
	case "ptrAPI_ISteamInput_GetDigitalActionData":
		result := callInputDigitalActionData(ptr, steamInput, 0, 0)
		validateFFIResult(t, name, result)
	case "ptrAPI_ISteamInput_GetAnalogActionData":
		result := callInputAnalogActionData(ptr, steamInput, 0, 0)
		validateFFIResult(t, name, result)
	case "ptrAPI_ISteamInput_GetMotionData":
		result := callInputMotionData(ptr, steamInput, 0)
		validateFFIResult(t, name, result)
	default:
		t.Fatalf("unknown ffi pointer %s", name)
	}
}

func validateFFIResult(t *testing.T, name string, result interface{}) {
	t.Helper()

	value := reflect.ValueOf(result)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.Float32, reflect.Float64:
			if math.IsNaN(field.Convert(reflect.TypeOf(float64(0))).Float()) {
				t.Fatalf("%s returned NaN for field %d", name, i)
			}
		}
	}
}

func callRegisteredFunction(t *testing.T, name string, fn reflect.Value, interfacePtrs map[string]uintptr, initState initStatus) {
	t.Helper()

	fnType := fn.Type()
	args, keepAlive, ok := buildArgs(t, name, fnType, interfacePtrs, initState)
	if !ok {
		return
	}
	for _, value := range keepAlive {
		_ = value
	}

	results := fn.Call(args)
	validateResults(t, name, results, initState)
	t.Logf("call ok for %s", name)
}

func buildArgs(t *testing.T, name string, fnType reflect.Type, interfacePtrs map[string]uintptr, initState initStatus) ([]reflect.Value, []interface{}, bool) {
	t.Helper()

	args := make([]reflect.Value, 0, fnType.NumIn())
	keepAlive := make([]interface{}, 0, fnType.NumIn())
	var lastBufferSize int
	lastWasBuffer := false

	if name == "ptrAPI_InitFlat" && fnType.NumIn() == 1 && fnType.In(0).Kind() == reflect.Uintptr {
		var msg steamErrMsg
		args = append(args, reflect.ValueOf(uintptr(unsafe.Pointer(&msg))))
		keepAlive = append(keepAlive, &msg)
		return args, keepAlive, true
	}

	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)
		if i == 0 {
			if ifaceName := interfaceNameFor(name); ifaceName != "" && argType.Kind() == reflect.Uintptr {
				iface := interfacePtrs[ifaceName]
				if iface == 0 {
					t.Logf("skipping %s because interface pointer %s is 0", name, ifaceName)
					return nil, nil, false
				}
				args = append(args, reflect.ValueOf(iface))
				lastWasBuffer = false
				continue
			}
		}

		switch argType.Kind() {
		case reflect.Bool:
			args = append(args, reflect.ValueOf(true))
			lastWasBuffer = false
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := int64(1)
			if lastWasBuffer && (argType.Kind() == reflect.Int32 || argType.Kind() == reflect.Int64) {
				value = int64(lastBufferSize)
			}
			args = append(args, reflect.ValueOf(value).Convert(argType))
			lastWasBuffer = false
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if argType.Kind() == reflect.Uintptr {
				buf := make([]byte, 256)
				ptr := uintptr(unsafe.Pointer(&buf[0]))
				args = append(args, reflect.ValueOf(ptr))
				keepAlive = append(keepAlive, buf)
				lastBufferSize = len(buf)
				lastWasBuffer = true
				continue
			}
			value := uint64(1)
			if lastWasBuffer && (argType.Kind() == reflect.Uint32 || argType.Kind() == reflect.Uint64) {
				value = uint64(lastBufferSize)
			}
			args = append(args, reflect.ValueOf(value).Convert(argType))
			lastWasBuffer = false
		case reflect.Float32, reflect.Float64:
			args = append(args, reflect.Zero(argType))
			lastWasBuffer = false
		case reflect.String:
			args = append(args, reflect.ValueOf("test"))
			lastWasBuffer = false
		default:
			t.Fatalf("%s has unsupported arg type %s", name, argType)
		}
	}

	return args, keepAlive, true
}

func interfaceNameFor(name string) string {
	switch {
	case strings.HasPrefix(name, "ptrAPI_ISteamApps_"):
		return "ISteamApps"
	case strings.HasPrefix(name, "ptrAPI_ISteamFriends_"):
		return "ISteamFriends"
	case strings.HasPrefix(name, "ptrAPI_ISteamMatchmaking_"):
		return "ISteamMatchmaking"
	case strings.HasPrefix(name, "ptrAPI_ISteamMatchmakingServers_"):
		return "ISteamMatchmakingServers"
	case strings.HasPrefix(name, "ptrAPI_ISteamHTTP_"):
		return "ISteamHTTP"
	case strings.HasPrefix(name, "ptrAPI_ISteamUGC_"):
		return "ISteamUGC"
	case strings.HasPrefix(name, "ptrAPI_ISteamInventory_"):
		return "ISteamInventory"
	case strings.HasPrefix(name, "ptrAPI_ISteamInput_"):
		return "ISteamInput"
	case strings.HasPrefix(name, "ptrAPI_ISteamRemotePlay_"):
		return "ISteamRemotePlay"
	case strings.HasPrefix(name, "ptrAPI_ISteamRemoteStorage_"):
		return "ISteamRemoteStorage"
	case strings.HasPrefix(name, "ptrAPI_ISteamUser_"):
		return "ISteamUser"
	case strings.HasPrefix(name, "ptrAPI_ISteamUserStats_"):
		return "ISteamUserStats"
	case strings.HasPrefix(name, "ptrAPI_ISteamUtils_"):
		return "ISteamUtils"
	case strings.HasPrefix(name, "ptrAPI_ISteamNetworkingUtils_"):
		return "ISteamNetworkingUtils"
	case strings.HasPrefix(name, "ptrAPI_ISteamNetworkingMessages_"):
		return "ISteamNetworkingMessages"
	case strings.HasPrefix(name, "ptrAPI_ISteamNetworkingSockets_"):
		return "ISteamNetworkingSockets"
	case strings.HasPrefix(name, "ptrAPI_ISteamGameServer_"):
		return "ISteamGameServer"
	default:
		return ""
	}
}

func validateResults(t *testing.T, name string, results []reflect.Value, initState initStatus) {
	t.Helper()

	for idx, result := range results {
		switch result.Kind() {
		case reflect.Bool:
			t.Logf("%s return[%d]=%v", name, idx, result.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			t.Logf("%s return[%d]=%d", name, idx, result.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			t.Logf("%s return[%d]=%d", name, idx, result.Uint())
		case reflect.Uintptr:
			if initState.ok && result.Uint() == 0 && strings.HasPrefix(name, "ptrAPI_Steam") {
				// Some interfaces might return 0 depending on the SDK version, configuration, or environment
				// (e.g., SteamGameServer on a client, or SteamNetworkingSockets on older SDKs).
				allowedZero := map[string]bool{
					"ptrAPI_SteamGameServer":         true,
					"ptrAPI_SteamNetworkingSockets":  true,
					"ptrAPI_SteamNetworkingMessages": true,
					"ptrAPI_SteamNetworkingUtils":    true,
				}
				if !allowedZero[name] {
					t.Errorf("%s return[%d]=0 with initialized API", name, idx)
				}
			}
			t.Logf("%s return[%d]=0x%x", name, idx, result.Uint())
		case reflect.Float32, reflect.Float64:
			value := result.Convert(reflect.TypeOf(float64(0))).Float()
			if math.IsNaN(value) || math.IsInf(value, 0) {
				t.Errorf("%s return[%d] invalid float %f", name, idx, value)
			}
			t.Logf("%s return[%d]=%f", name, idx, value)
		case reflect.String:
			value := result.String()
			if !utf8.ValidString(value) {
				t.Errorf("%s return[%d] invalid utf8", name, idx)
			}
			t.Logf("%s return[%d]=%q", name, idx, value)
		default:
			t.Errorf("%s return[%d] unsupported kind %s", name, idx, result.Kind())
		}
	}
}

func isSupportedType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func TestSDKInputStructReturns(t *testing.T) {
	initOnce.Do(func() {
		libHandle = loadSDKLibrary(t)
		registerFunctions(libHandle)
		registerInputStructReturns(libHandle)
		initResult = initSteamAPI(t)
	})

	ptrs := []struct {
		name string
		ptr  uintptr
	}{
		{name: flatAPI_ISteamInput_GetDigitalActionData, ptr: ptrAPI_ISteamInput_GetDigitalActionData},
		{name: flatAPI_ISteamInput_GetAnalogActionData, ptr: ptrAPI_ISteamInput_GetAnalogActionData},
		{name: flatAPI_ISteamInput_GetMotionData, ptr: ptrAPI_ISteamInput_GetMotionData},
	}

	for _, item := range ptrs {
		if item.ptr == 0 {
			t.Fatalf("%s not registered for struct returns", item.name)
		}
	}

	_ = callInputDigitalActionData(ptrAPI_ISteamInput_GetDigitalActionData, 0, 0, 0)
	_ = callInputAnalogActionData(ptrAPI_ISteamInput_GetAnalogActionData, 0, 0, 0)
	_ = callInputMotionData(ptrAPI_ISteamInput_GetMotionData, 0, 0)
}

type signatureExpectation struct {
	name     string
	expected interface{}
}

type registeredFunction struct {
	name  string
	value interface{}
}

func allRegisteredFunctions() []registeredFunction {
	return []registeredFunction{
		{name: "ptrAPI_RestartAppIfNecessary", value: ptrAPI_RestartAppIfNecessary},
		{name: "ptrAPI_InitFlat", value: ptrAPI_InitFlat},
		{name: "ptrAPI_RunCallbacks", value: ptrAPI_RunCallbacks},
		{name: "ptrAPI_Shutdown", value: ptrAPI_Shutdown},
		{name: "ptrAPI_IsSteamRunning", value: ptrAPI_IsSteamRunning},
		{name: "ptrAPI_GetSteamInstallPath", value: ptrAPI_GetSteamInstallPath},
		{name: "ptrAPI_ReleaseCurrentThreadMemory", value: ptrAPI_ReleaseCurrentThreadMemory},

		{name: "ptrAPI_SteamApps", value: ptrAPI_SteamApps},
		{name: "ptrAPI_ISteamApps_BIsSubscribed", value: ptrAPI_ISteamApps_BIsSubscribed},
		{name: "ptrAPI_ISteamApps_BIsLowViolence", value: ptrAPI_ISteamApps_BIsLowViolence},
		{name: "ptrAPI_ISteamApps_BIsCybercafe", value: ptrAPI_ISteamApps_BIsCybercafe},
		{name: "ptrAPI_ISteamApps_BIsVACBanned", value: ptrAPI_ISteamApps_BIsVACBanned},
		{name: "ptrAPI_ISteamApps_BGetDLCDataByIndex", value: ptrAPI_ISteamApps_BGetDLCDataByIndex},
		{name: "ptrAPI_ISteamApps_BIsDlcInstalled", value: ptrAPI_ISteamApps_BIsDlcInstalled},
		{name: "ptrAPI_ISteamApps_GetAvailableGameLanguages", value: ptrAPI_ISteamApps_GetAvailableGameLanguages},
		{name: "ptrAPI_ISteamApps_BIsSubscribedApp", value: ptrAPI_ISteamApps_BIsSubscribedApp},
		{name: "ptrAPI_ISteamApps_GetEarliestPurchaseUnixTime", value: ptrAPI_ISteamApps_GetEarliestPurchaseUnixTime},
		{name: "ptrAPI_ISteamApps_BIsSubscribedFromFreeWeekend", value: ptrAPI_ISteamApps_BIsSubscribedFromFreeWeekend},
		{name: "ptrAPI_ISteamApps_GetAppInstallDir", value: ptrAPI_ISteamApps_GetAppInstallDir},
		{name: "ptrAPI_ISteamApps_GetCurrentGameLanguage", value: ptrAPI_ISteamApps_GetCurrentGameLanguage},
		{name: "ptrAPI_ISteamApps_GetDLCCount", value: ptrAPI_ISteamApps_GetDLCCount},
		{name: "ptrAPI_ISteamApps_InstallDLC", value: ptrAPI_ISteamApps_InstallDLC},
		{name: "ptrAPI_ISteamApps_UninstallDLC", value: ptrAPI_ISteamApps_UninstallDLC},
		{name: "ptrAPI_ISteamApps_RequestAppProofOfPurchaseKey", value: ptrAPI_ISteamApps_RequestAppProofOfPurchaseKey},
		{name: "ptrAPI_ISteamApps_GetCurrentBetaName", value: ptrAPI_ISteamApps_GetCurrentBetaName},
		{name: "ptrAPI_ISteamApps_MarkContentCorrupt", value: ptrAPI_ISteamApps_MarkContentCorrupt},
		{name: "ptrAPI_ISteamApps_GetInstalledDepots", value: ptrAPI_ISteamApps_GetInstalledDepots},
		{name: "ptrAPI_ISteamApps_BIsAppInstalled", value: ptrAPI_ISteamApps_BIsAppInstalled},
		{name: "ptrAPI_ISteamApps_GetAppOwner", value: ptrAPI_ISteamApps_GetAppOwner},
		{name: "ptrAPI_ISteamApps_GetLaunchQueryParam", value: ptrAPI_ISteamApps_GetLaunchQueryParam},
		{name: "ptrAPI_ISteamApps_GetDlcDownloadProgress", value: ptrAPI_ISteamApps_GetDlcDownloadProgress},
		{name: "ptrAPI_ISteamApps_GetAppBuildId", value: ptrAPI_ISteamApps_GetAppBuildId},
		{name: "ptrAPI_ISteamApps_RequestAllProofOfPurchaseKeys", value: ptrAPI_ISteamApps_RequestAllProofOfPurchaseKeys},
		{name: "ptrAPI_ISteamApps_GetFileDetails", value: ptrAPI_ISteamApps_GetFileDetails},
		{name: "ptrAPI_ISteamApps_GetLaunchCommandLine", value: ptrAPI_ISteamApps_GetLaunchCommandLine},
		{name: "ptrAPI_ISteamApps_BIsSubscribedFromFamilySharing", value: ptrAPI_ISteamApps_BIsSubscribedFromFamilySharing},
		{name: "ptrAPI_ISteamApps_BIsTimedTrial", value: ptrAPI_ISteamApps_BIsTimedTrial},
		{name: "ptrAPI_ISteamApps_SetDlcContext", value: ptrAPI_ISteamApps_SetDlcContext},
		{name: "ptrAPI_ISteamApps_GetNumBetas", value: ptrAPI_ISteamApps_GetNumBetas},
		{name: "ptrAPI_ISteamApps_GetBetaInfo", value: ptrAPI_ISteamApps_GetBetaInfo},
		{name: "ptrAPI_ISteamApps_SetActiveBeta", value: ptrAPI_ISteamApps_SetActiveBeta},

		{name: "ptrAPI_SteamFriends", value: ptrAPI_SteamFriends},
		{name: "ptrAPI_ISteamFriends_GetPersonaName", value: ptrAPI_ISteamFriends_GetPersonaName},
		{name: "ptrAPI_ISteamFriends_GetPersonaState", value: ptrAPI_ISteamFriends_GetPersonaState},
		{name: "ptrAPI_ISteamFriends_GetFriendCount", value: ptrAPI_ISteamFriends_GetFriendCount},
		{name: "ptrAPI_ISteamFriends_GetFriendByIndex", value: ptrAPI_ISteamFriends_GetFriendByIndex},
		{name: "ptrAPI_ISteamFriends_GetFriendRelationship", value: ptrAPI_ISteamFriends_GetFriendRelationship},
		{name: "ptrAPI_ISteamFriends_GetFriendPersonaState", value: ptrAPI_ISteamFriends_GetFriendPersonaState},
		{name: "ptrAPI_ISteamFriends_GetFriendPersonaName", value: ptrAPI_ISteamFriends_GetFriendPersonaName},
		{name: "ptrAPI_ISteamFriends_GetFriendPersonaNameHistory", value: ptrAPI_ISteamFriends_GetFriendPersonaNameHistory},
		{name: "ptrAPI_ISteamFriends_GetFriendSteamLevel", value: ptrAPI_ISteamFriends_GetFriendSteamLevel},
		{name: "ptrAPI_ISteamFriends_GetSmallFriendAvatar", value: ptrAPI_ISteamFriends_GetSmallFriendAvatar},
		{name: "ptrAPI_ISteamFriends_GetMediumFriendAvatar", value: ptrAPI_ISteamFriends_GetMediumFriendAvatar},
		{name: "ptrAPI_ISteamFriends_GetLargeFriendAvatar", value: ptrAPI_ISteamFriends_GetLargeFriendAvatar},
		{name: "ptrAPI_ISteamFriends_SetRichPresence", value: ptrAPI_ISteamFriends_SetRichPresence},
		{name: "ptrAPI_ISteamFriends_GetFriendGamePlayed", value: ptrAPI_ISteamFriends_GetFriendGamePlayed},
		{name: "ptrAPI_ISteamFriends_InviteUserToGame", value: ptrAPI_ISteamFriends_InviteUserToGame},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlay", value: ptrAPI_ISteamFriends_ActivateGameOverlay},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayToUser", value: ptrAPI_ISteamFriends_ActivateGameOverlayToUser},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayToWebPage", value: ptrAPI_ISteamFriends_ActivateGameOverlayToWebPage},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayToStore", value: ptrAPI_ISteamFriends_ActivateGameOverlayToStore},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialog", value: ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialog},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString", value: ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString},

		{name: "ptrAPI_SteamMatchmaking", value: ptrAPI_SteamMatchmaking},
		{name: "ptrAPI_ISteamMatchmaking_GetFavoriteGameCount", value: ptrAPI_ISteamMatchmaking_GetFavoriteGameCount},
		{name: "ptrAPI_ISteamMatchmaking_GetFavoriteGame", value: ptrAPI_ISteamMatchmaking_GetFavoriteGame},
		{name: "ptrAPI_ISteamMatchmaking_AddFavoriteGame", value: ptrAPI_ISteamMatchmaking_AddFavoriteGame},
		{name: "ptrAPI_ISteamMatchmaking_RemoveFavoriteGame", value: ptrAPI_ISteamMatchmaking_RemoveFavoriteGame},
		{name: "ptrAPI_ISteamMatchmaking_RequestLobbyList", value: ptrAPI_ISteamMatchmaking_RequestLobbyList},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter", value: ptrAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyByIndex", value: ptrAPI_ISteamMatchmaking_GetLobbyByIndex},
		{name: "ptrAPI_ISteamMatchmaking_CreateLobby", value: ptrAPI_ISteamMatchmaking_CreateLobby},
		{name: "ptrAPI_ISteamMatchmaking_JoinLobby", value: ptrAPI_ISteamMatchmaking_JoinLobby},
		{name: "ptrAPI_ISteamMatchmaking_LeaveLobby", value: ptrAPI_ISteamMatchmaking_LeaveLobby},
		{name: "ptrAPI_ISteamMatchmaking_InviteUserToLobby", value: ptrAPI_ISteamMatchmaking_InviteUserToLobby},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyMemberLimit", value: ptrAPI_ISteamMatchmaking_GetLobbyMemberLimit},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyMemberLimit", value: ptrAPI_ISteamMatchmaking_SetLobbyMemberLimit},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyType", value: ptrAPI_ISteamMatchmaking_SetLobbyType},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyJoinable", value: ptrAPI_ISteamMatchmaking_SetLobbyJoinable},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyOwner", value: ptrAPI_ISteamMatchmaking_GetLobbyOwner},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyOwner", value: ptrAPI_ISteamMatchmaking_SetLobbyOwner},
		{name: "ptrAPI_ISteamMatchmaking_SetLinkedLobby", value: ptrAPI_ISteamMatchmaking_SetLinkedLobby},
		{name: "ptrAPI_ISteamMatchmaking_GetNumLobbyMembers", value: ptrAPI_ISteamMatchmaking_GetNumLobbyMembers},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyMemberByIndex", value: ptrAPI_ISteamMatchmaking_GetLobbyMemberByIndex},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyData", value: ptrAPI_ISteamMatchmaking_SetLobbyData},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyData", value: ptrAPI_ISteamMatchmaking_GetLobbyData},
		{name: "ptrAPI_ISteamMatchmaking_DeleteLobbyData", value: ptrAPI_ISteamMatchmaking_DeleteLobbyData},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyDataCount", value: ptrAPI_ISteamMatchmaking_GetLobbyDataCount},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyDataByIndex", value: ptrAPI_ISteamMatchmaking_GetLobbyDataByIndex},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyMemberData", value: ptrAPI_ISteamMatchmaking_SetLobbyMemberData},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyMemberData", value: ptrAPI_ISteamMatchmaking_GetLobbyMemberData},
		{name: "ptrAPI_ISteamMatchmaking_SendLobbyChatMsg", value: ptrAPI_ISteamMatchmaking_SendLobbyChatMsg},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyChatEntry", value: ptrAPI_ISteamMatchmaking_GetLobbyChatEntry},
		{name: "ptrAPI_ISteamMatchmaking_RequestLobbyData", value: ptrAPI_ISteamMatchmaking_RequestLobbyData},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyGameServer", value: ptrAPI_ISteamMatchmaking_SetLobbyGameServer},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyGameServer", value: ptrAPI_ISteamMatchmaking_GetLobbyGameServer},
		{name: "ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite", value: ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite},

		{name: "ptrAPI_ISteamMatchmakingServers_RequestInternetServerList", value: ptrAPI_ISteamMatchmakingServers_RequestInternetServerList},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestLANServerList", value: ptrAPI_ISteamMatchmakingServers_RequestLANServerList},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestFriendsServerList", value: ptrAPI_ISteamMatchmakingServers_RequestFriendsServerList},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestFavoritesServerList", value: ptrAPI_ISteamMatchmakingServers_RequestFavoritesServerList},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestHistoryServerList", value: ptrAPI_ISteamMatchmakingServers_RequestHistoryServerList},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestSpectatorServerList", value: ptrAPI_ISteamMatchmakingServers_RequestSpectatorServerList},
		{name: "ptrAPI_ISteamMatchmakingServers_ReleaseRequest", value: ptrAPI_ISteamMatchmakingServers_ReleaseRequest},
		{name: "ptrAPI_ISteamMatchmakingServers_GetServerDetails", value: ptrAPI_ISteamMatchmakingServers_GetServerDetails},
		{name: "ptrAPI_ISteamMatchmakingServers_CancelQuery", value: ptrAPI_ISteamMatchmakingServers_CancelQuery},
		{name: "ptrAPI_ISteamMatchmakingServers_RefreshQuery", value: ptrAPI_ISteamMatchmakingServers_RefreshQuery},
		{name: "ptrAPI_ISteamMatchmakingServers_IsRefreshing", value: ptrAPI_ISteamMatchmakingServers_IsRefreshing},
		{name: "ptrAPI_ISteamMatchmakingServers_GetServerCount", value: ptrAPI_ISteamMatchmakingServers_GetServerCount},
		{name: "ptrAPI_ISteamMatchmakingServers_RefreshServer", value: ptrAPI_ISteamMatchmakingServers_RefreshServer},
		{name: "ptrAPI_ISteamMatchmakingServers_PingServer", value: ptrAPI_ISteamMatchmakingServers_PingServer},
		{name: "ptrAPI_ISteamMatchmakingServers_PlayerDetails", value: ptrAPI_ISteamMatchmakingServers_PlayerDetails},
		{name: "ptrAPI_ISteamMatchmakingServers_ServerRules", value: ptrAPI_ISteamMatchmakingServers_ServerRules},
		{name: "ptrAPI_ISteamMatchmakingServers_CancelServerQuery", value: ptrAPI_ISteamMatchmakingServers_CancelServerQuery},

		{name: "ptrAPI_SteamHTTP", value: ptrAPI_SteamHTTP},
		{name: "ptrAPI_ISteamHTTP_CreateHTTPRequest", value: ptrAPI_ISteamHTTP_CreateHTTPRequest},
		{name: "ptrAPI_ISteamHTTP_SetHTTPRequestHeaderValue", value: ptrAPI_ISteamHTTP_SetHTTPRequestHeaderValue},
		{name: "ptrAPI_ISteamHTTP_SendHTTPRequest", value: ptrAPI_ISteamHTTP_SendHTTPRequest},
		{name: "ptrAPI_ISteamHTTP_GetHTTPResponseBodySize", value: ptrAPI_ISteamHTTP_GetHTTPResponseBodySize},
		{name: "ptrAPI_ISteamHTTP_GetHTTPResponseBodyData", value: ptrAPI_ISteamHTTP_GetHTTPResponseBodyData},
		{name: "ptrAPI_ISteamHTTP_ReleaseHTTPRequest", value: ptrAPI_ISteamHTTP_ReleaseHTTPRequest},

		{name: "ptrAPI_SteamUGC", value: ptrAPI_SteamUGC},
		{name: "ptrAPI_ISteamUGC_GetNumSubscribedItems", value: ptrAPI_ISteamUGC_GetNumSubscribedItems},
		{name: "ptrAPI_ISteamUGC_GetSubscribedItems", value: ptrAPI_ISteamUGC_GetSubscribedItems},
		{name: "ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused", value: ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused},
		{name: "ptrAPI_ISteamUGC_GetNumDownloadedItems", value: ptrAPI_ISteamUGC_GetNumDownloadedItems},
		{name: "ptrAPI_ISteamUGC_GetDownloadedItems", value: ptrAPI_ISteamUGC_GetDownloadedItems},

		{name: "ptrAPI_SteamInventory", value: ptrAPI_SteamInventory},
		{name: "ptrAPI_ISteamInventory_GetResultStatus", value: ptrAPI_ISteamInventory_GetResultStatus},
		{name: "ptrAPI_ISteamInventory_GetResultItems", value: ptrAPI_ISteamInventory_GetResultItems},
		{name: "ptrAPI_ISteamInventory_DestroyResult", value: ptrAPI_ISteamInventory_DestroyResult},

		{name: "ptrAPI_SteamInput", value: ptrAPI_SteamInput},
		{name: "ptrAPI_ISteamInput_GetConnectedControllers", value: ptrAPI_ISteamInput_GetConnectedControllers},
		{name: "ptrAPI_ISteamInput_GetInputTypeForHandle", value: ptrAPI_ISteamInput_GetInputTypeForHandle},
		{name: "ptrAPI_ISteamInput_Init", value: ptrAPI_ISteamInput_Init},
		{name: "ptrAPI_ISteamInput_Shutdown", value: ptrAPI_ISteamInput_Shutdown},
		{name: "ptrAPI_ISteamInput_RunFrame", value: ptrAPI_ISteamInput_RunFrame},
		{name: "ptrAPI_ISteamInput_EnableDeviceCallbacks", value: ptrAPI_ISteamInput_EnableDeviceCallbacks},
		{name: "ptrAPI_ISteamInput_GetActionSetHandle", value: ptrAPI_ISteamInput_GetActionSetHandle},
		{name: "ptrAPI_ISteamInput_ActivateActionSet", value: ptrAPI_ISteamInput_ActivateActionSet},
		{name: "ptrAPI_ISteamInput_GetCurrentActionSet", value: ptrAPI_ISteamInput_GetCurrentActionSet},
		{name: "ptrAPI_ISteamInput_ActivateActionSetLayer", value: ptrAPI_ISteamInput_ActivateActionSetLayer},
		{name: "ptrAPI_ISteamInput_DeactivateActionSetLayer", value: ptrAPI_ISteamInput_DeactivateActionSetLayer},
		{name: "ptrAPI_ISteamInput_DeactivateAllActionSetLayers", value: ptrAPI_ISteamInput_DeactivateAllActionSetLayers},
		{name: "ptrAPI_ISteamInput_GetActiveActionSetLayers", value: ptrAPI_ISteamInput_GetActiveActionSetLayers},
		{name: "ptrAPI_ISteamInput_GetDigitalActionHandle", value: ptrAPI_ISteamInput_GetDigitalActionHandle},
		{name: "ptrAPI_ISteamInput_GetDigitalActionData", value: ptrAPI_ISteamInput_GetDigitalActionData},
		{name: "ptrAPI_ISteamInput_GetDigitalActionOrigins", value: ptrAPI_ISteamInput_GetDigitalActionOrigins},
		{name: "ptrAPI_ISteamInput_GetAnalogActionHandle", value: ptrAPI_ISteamInput_GetAnalogActionHandle},
		{name: "ptrAPI_ISteamInput_GetAnalogActionData", value: ptrAPI_ISteamInput_GetAnalogActionData},
		{name: "ptrAPI_ISteamInput_GetAnalogActionOrigins", value: ptrAPI_ISteamInput_GetAnalogActionOrigins},
		{name: "ptrAPI_ISteamInput_StopAnalogActionMomentum", value: ptrAPI_ISteamInput_StopAnalogActionMomentum},
		{name: "ptrAPI_ISteamInput_GetMotionData", value: ptrAPI_ISteamInput_GetMotionData},
		{name: "ptrAPI_ISteamInput_TriggerVibration", value: ptrAPI_ISteamInput_TriggerVibration},
		{name: "ptrAPI_ISteamInput_TriggerVibrationExtended", value: ptrAPI_ISteamInput_TriggerVibrationExtended},
		{name: "ptrAPI_ISteamInput_TriggerSimpleHapticEvent", value: ptrAPI_ISteamInput_TriggerSimpleHapticEvent},
		{name: "ptrAPI_ISteamInput_SetLEDColor", value: ptrAPI_ISteamInput_SetLEDColor},
		{name: "ptrAPI_ISteamInput_ShowBindingPanel", value: ptrAPI_ISteamInput_ShowBindingPanel},
		{name: "ptrAPI_ISteamInput_GetControllerForGamepadIndex", value: ptrAPI_ISteamInput_GetControllerForGamepadIndex},
		{name: "ptrAPI_ISteamInput_GetGamepadIndexForController", value: ptrAPI_ISteamInput_GetGamepadIndexForController},
		{name: "ptrAPI_ISteamInput_GetStringForActionOrigin", value: ptrAPI_ISteamInput_GetStringForActionOrigin},
		{name: "ptrAPI_ISteamInput_GetGlyphForActionOrigin", value: ptrAPI_ISteamInput_GetGlyphForActionOrigin},
		{name: "ptrAPI_ISteamInput_GetRemotePlaySessionID", value: ptrAPI_ISteamInput_GetRemotePlaySessionID},

		{name: "ptrAPI_SteamRemotePlay", value: ptrAPI_SteamRemotePlay},
		{name: "ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether", value: ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether},
		{name: "ptrAPI_ISteamRemotePlay_GetSessionGuestID", value: ptrAPI_ISteamRemotePlay_GetSessionGuestID},
		{name: "ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar", value: ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar},
		{name: "ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar", value: ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar},
		{name: "ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar", value: ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar},

		{name: "ptrAPI_SteamRemoteStorage", value: ptrAPI_SteamRemoteStorage},
		{name: "ptrAPI_ISteamRemoteStorage_FileWrite", value: ptrAPI_ISteamRemoteStorage_FileWrite},
		{name: "ptrAPI_ISteamRemoteStorage_FileRead", value: ptrAPI_ISteamRemoteStorage_FileRead},
		{name: "ptrAPI_ISteamRemoteStorage_FileDelete", value: ptrAPI_ISteamRemoteStorage_FileDelete},
		{name: "ptrAPI_ISteamRemoteStorage_GetFileSize", value: ptrAPI_ISteamRemoteStorage_GetFileSize},

		{name: "ptrAPI_SteamUser", value: ptrAPI_SteamUser},
		{name: "ptrAPI_ISteamUser_AdvertiseGame", value: ptrAPI_ISteamUser_AdvertiseGame},
		{name: "ptrAPI_ISteamUser_BeginAuthSession", value: ptrAPI_ISteamUser_BeginAuthSession},
		{name: "ptrAPI_ISteamUser_BIsBehindNAT", value: ptrAPI_ISteamUser_BIsBehindNAT},
		{name: "ptrAPI_ISteamUser_BIsPhoneIdentifying", value: ptrAPI_ISteamUser_BIsPhoneIdentifying},
		{name: "ptrAPI_ISteamUser_BIsPhoneRequiringVerification", value: ptrAPI_ISteamUser_BIsPhoneRequiringVerification},
		{name: "ptrAPI_ISteamUser_BIsPhoneVerified", value: ptrAPI_ISteamUser_BIsPhoneVerified},
		{name: "ptrAPI_ISteamUser_BIsTwoFactorEnabled", value: ptrAPI_ISteamUser_BIsTwoFactorEnabled},
		{name: "ptrAPI_ISteamUser_BLoggedOn", value: ptrAPI_ISteamUser_BLoggedOn},
		{name: "ptrAPI_ISteamUser_BSetDurationControlOnlineState", value: ptrAPI_ISteamUser_BSetDurationControlOnlineState},
		{name: "ptrAPI_ISteamUser_CancelAuthTicket", value: ptrAPI_ISteamUser_CancelAuthTicket},
		{name: "ptrAPI_ISteamUser_DecompressVoice", value: ptrAPI_ISteamUser_DecompressVoice},
		{name: "ptrAPI_ISteamUser_EndAuthSession", value: ptrAPI_ISteamUser_EndAuthSession},
		{name: "ptrAPI_ISteamUser_GetAuthSessionTicket", value: ptrAPI_ISteamUser_GetAuthSessionTicket},
		{name: "ptrAPI_ISteamUser_GetAuthTicketForWebApi", value: ptrAPI_ISteamUser_GetAuthTicketForWebApi},
		{name: "ptrAPI_ISteamUser_GetAvailableVoice", value: ptrAPI_ISteamUser_GetAvailableVoice},
		{name: "ptrAPI_ISteamUser_GetDurationControl", value: ptrAPI_ISteamUser_GetDurationControl},
		{name: "ptrAPI_ISteamUser_GetEncryptedAppTicket", value: ptrAPI_ISteamUser_GetEncryptedAppTicket},
		{name: "ptrAPI_ISteamUser_GetGameBadgeLevel", value: ptrAPI_ISteamUser_GetGameBadgeLevel},
		{name: "ptrAPI_ISteamUser_GetHSteamUser", value: ptrAPI_ISteamUser_GetHSteamUser},
		{name: "ptrAPI_ISteamUser_GetPlayerSteamLevel", value: ptrAPI_ISteamUser_GetPlayerSteamLevel},
		{name: "ptrAPI_ISteamUser_GetSteamID", value: ptrAPI_ISteamUser_GetSteamID},
		{name: "ptrAPI_ISteamUser_GetUserDataFolder", value: ptrAPI_ISteamUser_GetUserDataFolder},
		{name: "ptrAPI_ISteamUser_GetVoice", value: ptrAPI_ISteamUser_GetVoice},
		{name: "ptrAPI_ISteamUser_GetVoiceOptimalSampleRate", value: ptrAPI_ISteamUser_GetVoiceOptimalSampleRate},
		{name: "ptrAPI_ISteamUser_InitiateGameConnection", value: ptrAPI_ISteamUser_InitiateGameConnection},
		{name: "ptrAPI_ISteamUser_RequestEncryptedAppTicket", value: ptrAPI_ISteamUser_RequestEncryptedAppTicket},
		{name: "ptrAPI_ISteamUser_RequestStoreAuthURL", value: ptrAPI_ISteamUser_RequestStoreAuthURL},
		{name: "ptrAPI_ISteamUser_StartVoiceRecording", value: ptrAPI_ISteamUser_StartVoiceRecording},
		{name: "ptrAPI_ISteamUser_StopVoiceRecording", value: ptrAPI_ISteamUser_StopVoiceRecording},
		{name: "ptrAPI_ISteamUser_TerminateGameConnection", value: ptrAPI_ISteamUser_TerminateGameConnection},
		{name: "ptrAPI_ISteamUser_TrackAppUsageEvent", value: ptrAPI_ISteamUser_TrackAppUsageEvent},
		{name: "ptrAPI_ISteamUser_UserHasLicenseForApp", value: ptrAPI_ISteamUser_UserHasLicenseForApp},

		{name: "ptrAPI_SteamUserStats", value: ptrAPI_SteamUserStats},
		{name: "ptrAPI_ISteamUserStats_GetAchievement", value: ptrAPI_ISteamUserStats_GetAchievement},
		{name: "ptrAPI_ISteamUserStats_SetAchievement", value: ptrAPI_ISteamUserStats_SetAchievement},
		{name: "ptrAPI_ISteamUserStats_ClearAchievement", value: ptrAPI_ISteamUserStats_ClearAchievement},
		{name: "ptrAPI_ISteamUserStats_StoreStats", value: ptrAPI_ISteamUserStats_StoreStats},

		{name: "ptrAPI_SteamUtils", value: ptrAPI_SteamUtils},
		{name: "ptrAPI_ISteamUtils_GetSecondsSinceAppActive", value: ptrAPI_ISteamUtils_GetSecondsSinceAppActive},
		{name: "ptrAPI_ISteamUtils_GetSecondsSinceComputerActive", value: ptrAPI_ISteamUtils_GetSecondsSinceComputerActive},
		{name: "ptrAPI_ISteamUtils_GetConnectedUniverse", value: ptrAPI_ISteamUtils_GetConnectedUniverse},
		{name: "ptrAPI_ISteamUtils_GetServerRealTime", value: ptrAPI_ISteamUtils_GetServerRealTime},
		{name: "ptrAPI_ISteamUtils_GetIPCountry", value: ptrAPI_ISteamUtils_GetIPCountry},
		{name: "ptrAPI_ISteamUtils_GetImageSize", value: ptrAPI_ISteamUtils_GetImageSize},
		{name: "ptrAPI_ISteamUtils_GetImageRGBA", value: ptrAPI_ISteamUtils_GetImageRGBA},
		{name: "ptrAPI_ISteamUtils_GetCurrentBatteryPower", value: ptrAPI_ISteamUtils_GetCurrentBatteryPower},
		{name: "ptrAPI_ISteamUtils_GetAppID", value: ptrAPI_ISteamUtils_GetAppID},
		{name: "ptrAPI_ISteamUtils_SetOverlayNotificationPosition", value: ptrAPI_ISteamUtils_SetOverlayNotificationPosition},
		{name: "ptrAPI_ISteamUtils_IsAPICallCompleted", value: ptrAPI_ISteamUtils_IsAPICallCompleted},
		{name: "ptrAPI_ISteamUtils_GetAPICallFailureReason", value: ptrAPI_ISteamUtils_GetAPICallFailureReason},
		{name: "ptrAPI_ISteamUtils_GetAPICallResult", value: ptrAPI_ISteamUtils_GetAPICallResult},
		{name: "ptrAPI_ISteamUtils_GetIPCCallCount", value: ptrAPI_ISteamUtils_GetIPCCallCount},
		{name: "ptrAPI_ISteamUtils_IsOverlayEnabled", value: ptrAPI_ISteamUtils_IsOverlayEnabled},
		{name: "ptrAPI_ISteamUtils_BOverlayNeedsPresent", value: ptrAPI_ISteamUtils_BOverlayNeedsPresent},
		{name: "ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck", value: ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck},
		{name: "ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput", value: ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput},
		{name: "ptrAPI_ISteamUtils_SetOverlayNotificationInset", value: ptrAPI_ISteamUtils_SetOverlayNotificationInset},

		{name: "ptrAPI_SteamNetworkingUtils", value: ptrAPI_SteamNetworkingUtils},
		{name: "ptrAPI_ISteamNetworkingUtils_AllocateMessage", value: ptrAPI_ISteamNetworkingUtils_AllocateMessage},
		{name: "ptrAPI_ISteamNetworkingUtils_InitRelayNetworkAccess", value: ptrAPI_ISteamNetworkingUtils_InitRelayNetworkAccess},
		{name: "ptrAPI_ISteamNetworkingUtils_GetLocalTimestamp", value: ptrAPI_ISteamNetworkingUtils_GetLocalTimestamp},

		{name: "ptrAPI_SteamGameServer", value: ptrAPI_SteamGameServer},
		{name: "ptrAPI_ISteamGameServer_AssociateWithClan", value: ptrAPI_ISteamGameServer_AssociateWithClan},
		{name: "ptrAPI_ISteamGameServer_BeginAuthSession", value: ptrAPI_ISteamGameServer_BeginAuthSession},
		{name: "ptrAPI_ISteamGameServer_BLoggedOn", value: ptrAPI_ISteamGameServer_BLoggedOn},
		{name: "ptrAPI_ISteamGameServer_BSecure", value: ptrAPI_ISteamGameServer_BSecure},
		{name: "ptrAPI_ISteamGameServer_BUpdateUserData", value: ptrAPI_ISteamGameServer_BUpdateUserData},
		{name: "ptrAPI_ISteamGameServer_CancelAuthTicket", value: ptrAPI_ISteamGameServer_CancelAuthTicket},
		{name: "ptrAPI_ISteamGameServer_ClearAllKeyValues", value: ptrAPI_ISteamGameServer_ClearAllKeyValues},
		{name: "ptrAPI_ISteamGameServer_ComputeNewPlayerCompatibility", value: ptrAPI_ISteamGameServer_ComputeNewPlayerCompatibility},
		{name: "ptrAPI_ISteamGameServer_CreateUnauthenticatedUserConnection", value: ptrAPI_ISteamGameServer_CreateUnauthenticatedUserConnection},
		{name: "ptrAPI_ISteamGameServer_EnableHeartbeats", value: ptrAPI_ISteamGameServer_EnableHeartbeats},
		{name: "ptrAPI_ISteamGameServer_EndAuthSession", value: ptrAPI_ISteamGameServer_EndAuthSession},
		{name: "ptrAPI_ISteamGameServer_ForceHeartbeat", value: ptrAPI_ISteamGameServer_ForceHeartbeat},
		{name: "ptrAPI_ISteamGameServer_GetAuthSessionTicket", value: ptrAPI_ISteamGameServer_GetAuthSessionTicket},
		{name: "ptrAPI_ISteamGameServer_GetGameplayStats", value: ptrAPI_ISteamGameServer_GetGameplayStats},
		{name: "ptrAPI_ISteamGameServer_GetNextOutgoingPacket", value: ptrAPI_ISteamGameServer_GetNextOutgoingPacket},
		{name: "ptrAPI_ISteamGameServer_GetPublicIP", value: ptrAPI_ISteamGameServer_GetPublicIP},
		{name: "ptrAPI_ISteamGameServer_GetServerReputation", value: ptrAPI_ISteamGameServer_GetServerReputation},
		{name: "ptrAPI_ISteamGameServer_GetSteamID", value: ptrAPI_ISteamGameServer_GetSteamID},
		{name: "ptrAPI_ISteamGameServer_HandleIncomingPacket", value: ptrAPI_ISteamGameServer_HandleIncomingPacket},
		{name: "ptrAPI_ISteamGameServer_InitGameServer", value: ptrAPI_ISteamGameServer_InitGameServer},
		{name: "ptrAPI_ISteamGameServer_LogOff", value: ptrAPI_ISteamGameServer_LogOff},
		{name: "ptrAPI_ISteamGameServer_LogOn", value: ptrAPI_ISteamGameServer_LogOn},
		{name: "ptrAPI_ISteamGameServer_LogOnAnonymous", value: ptrAPI_ISteamGameServer_LogOnAnonymous},
		{name: "ptrAPI_ISteamGameServer_RequestUserGroupStatus", value: ptrAPI_ISteamGameServer_RequestUserGroupStatus},
		{name: "ptrAPI_ISteamGameServer_SendUserConnectAndAuthenticate", value: ptrAPI_ISteamGameServer_SendUserConnectAndAuthenticate},
		{name: "ptrAPI_ISteamGameServer_SendUserDisconnect", value: ptrAPI_ISteamGameServer_SendUserDisconnect},
		{name: "ptrAPI_ISteamGameServer_SetBotPlayerCount", value: ptrAPI_ISteamGameServer_SetBotPlayerCount},
		{name: "ptrAPI_ISteamGameServer_SetDedicatedServer", value: ptrAPI_ISteamGameServer_SetDedicatedServer},
		{name: "ptrAPI_ISteamGameServer_SetGameData", value: ptrAPI_ISteamGameServer_SetGameData},
		{name: "ptrAPI_ISteamGameServer_SetGameDescription", value: ptrAPI_ISteamGameServer_SetGameDescription},
		{name: "ptrAPI_ISteamGameServer_SetGameTags", value: ptrAPI_ISteamGameServer_SetGameTags},
		{name: "ptrAPI_ISteamGameServer_SetHeartbeatInterval", value: ptrAPI_ISteamGameServer_SetHeartbeatInterval},
		{name: "ptrAPI_ISteamGameServer_SetKeyValue", value: ptrAPI_ISteamGameServer_SetKeyValue},
		{name: "ptrAPI_ISteamGameServer_SetMapName", value: ptrAPI_ISteamGameServer_SetMapName},
		{name: "ptrAPI_ISteamGameServer_SetMaxPlayerCount", value: ptrAPI_ISteamGameServer_SetMaxPlayerCount},
		{name: "ptrAPI_ISteamGameServer_SetModDir", value: ptrAPI_ISteamGameServer_SetModDir},
		{name: "ptrAPI_ISteamGameServer_SetPasswordProtected", value: ptrAPI_ISteamGameServer_SetPasswordProtected},
		{name: "ptrAPI_ISteamGameServer_SetProduct", value: ptrAPI_ISteamGameServer_SetProduct},
		{name: "ptrAPI_ISteamGameServer_SetRegion", value: ptrAPI_ISteamGameServer_SetRegion},
		{name: "ptrAPI_ISteamGameServer_SetServerName", value: ptrAPI_ISteamGameServer_SetServerName},
		{name: "ptrAPI_ISteamGameServer_SetSpectatorPort", value: ptrAPI_ISteamGameServer_SetSpectatorPort},
		{name: "ptrAPI_ISteamGameServer_SetSpectatorServerName", value: ptrAPI_ISteamGameServer_SetSpectatorServerName},
		{name: "ptrAPI_ISteamGameServer_UserHasLicenseForApp", value: ptrAPI_ISteamGameServer_UserHasLicenseForApp},
		{name: "ptrAPI_ISteamGameServer_WasRestartRequested", value: ptrAPI_ISteamGameServer_WasRestartRequested},

		{name: "ptrAPI_SteamNetworkingMessages", value: ptrAPI_SteamNetworkingMessages},
		{name: "ptrAPI_ISteamNetworkingMessages_SendMessageToUser", value: ptrAPI_ISteamNetworkingMessages_SendMessageToUser},
		{name: "ptrAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel", value: ptrAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel},
		{name: "ptrAPI_ISteamNetworkingMessages_AcceptSessionWithUser", value: ptrAPI_ISteamNetworkingMessages_AcceptSessionWithUser},
		{name: "ptrAPI_ISteamNetworkingMessages_CloseSessionWithUser", value: ptrAPI_ISteamNetworkingMessages_CloseSessionWithUser},
		{name: "ptrAPI_ISteamNetworkingMessages_CloseChannelWithUser", value: ptrAPI_ISteamNetworkingMessages_CloseChannelWithUser},

		{name: "ptrAPI_SteamNetworkingSockets", value: ptrAPI_SteamNetworkingSockets},
		{name: "ptrAPI_ISteamNetworkingSockets_CreateListenSocketIP", value: ptrAPI_ISteamNetworkingSockets_CreateListenSocketIP},
		{name: "ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P", value: ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P},
		{name: "ptrAPI_ISteamNetworkingSockets_ConnectByIPAddress", value: ptrAPI_ISteamNetworkingSockets_ConnectByIPAddress},
		{name: "ptrAPI_ISteamNetworkingSockets_ConnectP2P", value: ptrAPI_ISteamNetworkingSockets_ConnectP2P},
		{name: "ptrAPI_ISteamNetworkingSockets_AcceptConnection", value: ptrAPI_ISteamNetworkingSockets_AcceptConnection},
		{name: "ptrAPI_ISteamNetworkingSockets_CloseConnection", value: ptrAPI_ISteamNetworkingSockets_CloseConnection},
		{name: "ptrAPI_ISteamNetworkingSockets_CloseListenSocket", value: ptrAPI_ISteamNetworkingSockets_CloseListenSocket},
		{name: "ptrAPI_ISteamNetworkingSockets_SendMessageToConnection", value: ptrAPI_ISteamNetworkingSockets_SendMessageToConnection},
		{name: "ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection", value: ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection},
		{name: "ptrAPI_ISteamNetworkingSockets_CreatePollGroup", value: ptrAPI_ISteamNetworkingSockets_CreatePollGroup},
		{name: "ptrAPI_ISteamNetworkingSockets_DestroyPollGroup", value: ptrAPI_ISteamNetworkingSockets_DestroyPollGroup},
		{name: "ptrAPI_ISteamNetworkingSockets_SetConnectionPollGroup", value: ptrAPI_ISteamNetworkingSockets_SetConnectionPollGroup},
		{name: "ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup", value: ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup},
	}
}

func loadSDKLibrary(t *testing.T) uintptr {
	t.Helper()

	path, err := sdkLibraryPath()
	if err != nil {
		t.Fatalf("sdkLibraryPath: %v", err)
	}
	_ = os.Setenv(steamworksLibEnv, path)

	lib, err := loadLib()
	if err != nil {
		t.Fatalf("loadLib: %v", err)
	}

	return lib
}

func assertSignature(t *testing.T, name string, actual interface{}, expected interface{}) {
	t.Helper()

	if _, ok := expected.(uintptr); ok {
		if _, ok := actual.(uintptr); !ok {
			t.Fatalf("%s has type %T, want uintptr", name, actual)
		}
		return
	}

	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(expected)
	if actualType != expectedType {
		t.Fatalf("%s has type %v, want %v", name, actualType, expectedType)
	}
}

func signatureExpectations() []signatureExpectation {
	return []signatureExpectation{
		{name: "ptrAPI_RestartAppIfNecessary", expected: (func(uint32) bool)(nil)},
		{name: "ptrAPI_InitFlat", expected: (func(uintptr) ESteamAPIInitResult)(nil)},
		{name: "ptrAPI_RunCallbacks", expected: (func())(nil)},
		{name: "ptrAPI_Shutdown", expected: (func())(nil)},
		{name: "ptrAPI_IsSteamRunning", expected: (func() bool)(nil)},
		{name: "ptrAPI_GetSteamInstallPath", expected: (func() string)(nil)},
		{name: "ptrAPI_ReleaseCurrentThreadMemory", expected: (func())(nil)},

		{name: "ptrAPI_SteamApps", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamApps_BIsSubscribed", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_BIsLowViolence", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_BIsCybercafe", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_BIsVACBanned", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_BGetDLCDataByIndex", expected: (func(uintptr, int32, uintptr, uintptr, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamApps_BIsDlcInstalled", expected: (func(uintptr, AppId_t) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetAvailableGameLanguages", expected: (func(uintptr) string)(nil)},
		{name: "ptrAPI_ISteamApps_BIsSubscribedApp", expected: (func(uintptr, AppId_t) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetEarliestPurchaseUnixTime", expected: (func(uintptr, AppId_t) uint32)(nil)},
		{name: "ptrAPI_ISteamApps_BIsSubscribedFromFreeWeekend", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetAppInstallDir", expected: (func(uintptr, AppId_t, uintptr, int32) int32)(nil)},
		{name: "ptrAPI_ISteamApps_GetCurrentGameLanguage", expected: (func(uintptr) string)(nil)},
		{name: "ptrAPI_ISteamApps_GetDLCCount", expected: (func(uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamApps_InstallDLC", expected: (func(uintptr, AppId_t))(nil)},
		{name: "ptrAPI_ISteamApps_UninstallDLC", expected: (func(uintptr, AppId_t))(nil)},
		{name: "ptrAPI_ISteamApps_RequestAppProofOfPurchaseKey", expected: (func(uintptr, AppId_t))(nil)},
		{name: "ptrAPI_ISteamApps_GetCurrentBetaName", expected: (func(uintptr, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamApps_MarkContentCorrupt", expected: (func(uintptr, bool) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetInstalledDepots", expected: (func(uintptr, AppId_t, uintptr, uint32) uint32)(nil)},
		{name: "ptrAPI_ISteamApps_BIsAppInstalled", expected: (func(uintptr, AppId_t) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetAppOwner", expected: (func(uintptr) CSteamID)(nil)},
		{name: "ptrAPI_ISteamApps_GetLaunchQueryParam", expected: (func(uintptr, string) string)(nil)},
		{name: "ptrAPI_ISteamApps_GetDlcDownloadProgress", expected: (func(uintptr, AppId_t, uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetAppBuildId", expected: (func(uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamApps_RequestAllProofOfPurchaseKeys", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamApps_GetFileDetails", expected: (func(uintptr, string) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamApps_GetLaunchCommandLine", expected: (func(uintptr, uintptr, int32) int32)(nil)},
		{name: "ptrAPI_ISteamApps_BIsSubscribedFromFamilySharing", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_BIsTimedTrial", expected: (func(uintptr, uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamApps_SetDlcContext", expected: (func(uintptr, AppId_t) bool)(nil)},
		{name: "ptrAPI_ISteamApps_GetNumBetas", expected: (func(uintptr, uintptr, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamApps_GetBetaInfo", expected: (func(uintptr, int32, uintptr, uintptr, uintptr, uintptr, int32, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamApps_SetActiveBeta", expected: (func(uintptr, string) bool)(nil)},

		{name: "ptrAPI_SteamFriends", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamFriends_GetPersonaName", expected: (func(uintptr) string)(nil)},
		{name: "ptrAPI_ISteamFriends_GetPersonaState", expected: (func(uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendCount", expected: (func(uintptr, int32) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendByIndex", expected: (func(uintptr, int32, int32) CSteamID)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendRelationship", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendPersonaState", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendPersonaName", expected: (func(uintptr, CSteamID) string)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendPersonaNameHistory", expected: (func(uintptr, CSteamID, int32) string)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendSteamLevel", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetSmallFriendAvatar", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetMediumFriendAvatar", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_GetLargeFriendAvatar", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamFriends_SetRichPresence", expected: (func(uintptr, string, string) bool)(nil)},
		{name: "ptrAPI_ISteamFriends_GetFriendGamePlayed", expected: (func(uintptr, CSteamID, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamFriends_InviteUserToGame", expected: (func(uintptr, CSteamID, string) bool)(nil)},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlay", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayToUser", expected: (func(uintptr, string, CSteamID))(nil)},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayToWebPage", expected: (func(uintptr, string, EActivateGameOverlayToWebPageMode))(nil)},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayToStore", expected: (func(uintptr, AppId_t, EOverlayToStoreFlag))(nil)},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialog", expected: (func(uintptr, CSteamID))(nil)},
		{name: "ptrAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString", expected: (func(uintptr, string))(nil)},

		{name: "ptrAPI_SteamMatchmaking", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetFavoriteGameCount", expected: (func(uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetFavoriteGame", expected: (func(uintptr, int32, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddFavoriteGame", expected: (func(uintptr, AppId_t, uint32, uint16, uint16, uint32, uint32) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_RemoveFavoriteGame", expected: (func(uintptr, AppId_t, uint32, uint16, uint16, uint32) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_RequestLobbyList", expected: (func(uintptr) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter", expected: (func(uintptr, string, string, ELobbyComparison))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter", expected: (func(uintptr, string, int32, ELobbyComparison))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter", expected: (func(uintptr, string, int32))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable", expected: (func(uintptr, int32))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter", expected: (func(uintptr, ELobbyDistanceFilter))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter", expected: (func(uintptr, int32))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter", expected: (func(uintptr, CSteamID))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyByIndex", expected: (func(uintptr, int32) CSteamID)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_CreateLobby", expected: (func(uintptr, ELobbyType, int32) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_JoinLobby", expected: (func(uintptr, CSteamID) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_LeaveLobby", expected: (func(uintptr, CSteamID))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_InviteUserToLobby", expected: (func(uintptr, CSteamID, CSteamID) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyMemberLimit", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyMemberLimit", expected: (func(uintptr, CSteamID, int32) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyType", expected: (func(uintptr, CSteamID, ELobbyType) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyJoinable", expected: (func(uintptr, CSteamID, bool) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyOwner", expected: (func(uintptr, CSteamID) CSteamID)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyOwner", expected: (func(uintptr, CSteamID, CSteamID) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLinkedLobby", expected: (func(uintptr, CSteamID, CSteamID) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetNumLobbyMembers", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyMemberByIndex", expected: (func(uintptr, CSteamID, int32) CSteamID)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyData", expected: (func(uintptr, CSteamID, string, string) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyData", expected: (func(uintptr, CSteamID, string) string)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_DeleteLobbyData", expected: (func(uintptr, CSteamID, string) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyDataCount", expected: (func(uintptr, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyDataByIndex", expected: (func(uintptr, CSteamID, int32, uintptr, int32, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyMemberData", expected: (func(uintptr, CSteamID, string, string))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyMemberData", expected: (func(uintptr, CSteamID, CSteamID, string) string)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SendLobbyChatMsg", expected: (func(uintptr, CSteamID, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyChatEntry", expected: (func(uintptr, CSteamID, int32, uintptr, uintptr, int32, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_RequestLobbyData", expected: (func(uintptr, CSteamID) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_SetLobbyGameServer", expected: (func(uintptr, CSteamID, uint32, uint16, CSteamID))(nil)},
		{name: "ptrAPI_ISteamMatchmaking_GetLobbyGameServer", expected: (func(uintptr, CSteamID, uintptr, uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmaking_CheckForPSNGameBootInvite", expected: (func(uintptr, uintptr) bool)(nil)},

		{name: "ptrAPI_ISteamMatchmakingServers_RequestInternetServerList", expected: (func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestLANServerList", expected: (func(uintptr, AppId_t, uintptr) HServerListRequest)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestFriendsServerList", expected: (func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestFavoritesServerList", expected: (func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestHistoryServerList", expected: (func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RequestSpectatorServerList", expected: (func(uintptr, AppId_t, uintptr, uint32, uintptr) HServerListRequest)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_ReleaseRequest", expected: (func(uintptr, HServerListRequest))(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_GetServerDetails", expected: (func(uintptr, HServerListRequest, int32) uintptr)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_CancelQuery", expected: (func(uintptr, HServerListRequest))(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RefreshQuery", expected: (func(uintptr, HServerListRequest))(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_IsRefreshing", expected: (func(uintptr, HServerListRequest) bool)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_GetServerCount", expected: (func(uintptr, HServerListRequest) int32)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_RefreshServer", expected: (func(uintptr, HServerListRequest, int32))(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_PingServer", expected: (func(uintptr, uint32, uint16, uintptr) HServerQuery)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_PlayerDetails", expected: (func(uintptr, uint32, uint16, uintptr) HServerQuery)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_ServerRules", expected: (func(uintptr, uint32, uint16, uintptr) HServerQuery)(nil)},
		{name: "ptrAPI_ISteamMatchmakingServers_CancelServerQuery", expected: (func(uintptr, HServerQuery))(nil)},

		{name: "ptrAPI_SteamHTTP", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamHTTP_CreateHTTPRequest", expected: (func(uintptr, int32, string) HTTPRequestHandle)(nil)},
		{name: "ptrAPI_ISteamHTTP_SetHTTPRequestHeaderValue", expected: (func(uintptr, HTTPRequestHandle, string, string) bool)(nil)},
		{name: "ptrAPI_ISteamHTTP_SendHTTPRequest", expected: (func(uintptr, HTTPRequestHandle, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamHTTP_GetHTTPResponseBodySize", expected: (func(uintptr, HTTPRequestHandle, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamHTTP_GetHTTPResponseBodyData", expected: (func(uintptr, HTTPRequestHandle, uintptr, uint32) bool)(nil)},
		{name: "ptrAPI_ISteamHTTP_ReleaseHTTPRequest", expected: (func(uintptr, HTTPRequestHandle) bool)(nil)},

		{name: "ptrAPI_SteamUGC", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamUGC_GetNumSubscribedItems", expected: (func(uintptr, bool) uint32)(nil)},
		{name: "ptrAPI_ISteamUGC_GetSubscribedItems", expected: (func(uintptr, uintptr, uint32, bool) uint32)(nil)},
		{name: "ptrAPI_ISteamUGC_MarkDownloadedItemAsUnused", expected: (func(uintptr, PublishedFileId_t) bool)(nil)},
		{name: "ptrAPI_ISteamUGC_GetNumDownloadedItems", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUGC_GetDownloadedItems", expected: (func(uintptr, uintptr, uint32) uint32)(nil)},

		{name: "ptrAPI_SteamInventory", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamInventory_GetResultStatus", expected: (func(uintptr, SteamInventoryResult_t) int32)(nil)},
		{name: "ptrAPI_ISteamInventory_GetResultItems", expected: (func(uintptr, SteamInventoryResult_t, uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamInventory_DestroyResult", expected: (func(uintptr, SteamInventoryResult_t))(nil)},

		{name: "ptrAPI_SteamInput", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamInput_GetConnectedControllers", expected: (func(uintptr, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamInput_GetInputTypeForHandle", expected: (func(uintptr, InputHandle_t) int32)(nil)},
		{name: "ptrAPI_ISteamInput_Init", expected: (func(uintptr, bool) bool)(nil)},
		{name: "ptrAPI_ISteamInput_Shutdown", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamInput_RunFrame", expected: (func(uintptr, bool))(nil)},
		{name: "ptrAPI_ISteamInput_EnableDeviceCallbacks", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamInput_GetActionSetHandle", expected: (func(uintptr, string) InputActionSetHandle_t)(nil)},
		{name: "ptrAPI_ISteamInput_ActivateActionSet", expected: (func(uintptr, InputHandle_t, InputActionSetHandle_t))(nil)},
		{name: "ptrAPI_ISteamInput_GetCurrentActionSet", expected: (func(uintptr, InputHandle_t) InputActionSetHandle_t)(nil)},
		{name: "ptrAPI_ISteamInput_ActivateActionSetLayer", expected: (func(uintptr, InputHandle_t, InputActionSetHandle_t))(nil)},
		{name: "ptrAPI_ISteamInput_DeactivateActionSetLayer", expected: (func(uintptr, InputHandle_t, InputActionSetHandle_t))(nil)},
		{name: "ptrAPI_ISteamInput_DeactivateAllActionSetLayers", expected: (func(uintptr, InputHandle_t))(nil)},
		{name: "ptrAPI_ISteamInput_GetActiveActionSetLayers", expected: (func(uintptr, InputHandle_t, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamInput_GetDigitalActionHandle", expected: (func(uintptr, string) InputDigitalActionHandle_t)(nil)},
		{name: "ptrAPI_ISteamInput_GetDigitalActionData", expected: uintptr(0)},
		{name: "ptrAPI_ISteamInput_GetDigitalActionOrigins", expected: (func(uintptr, InputHandle_t, InputActionSetHandle_t, InputDigitalActionHandle_t, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamInput_GetAnalogActionHandle", expected: (func(uintptr, string) InputAnalogActionHandle_t)(nil)},
		{name: "ptrAPI_ISteamInput_GetAnalogActionData", expected: uintptr(0)},
		{name: "ptrAPI_ISteamInput_GetAnalogActionOrigins", expected: (func(uintptr, InputHandle_t, InputActionSetHandle_t, InputAnalogActionHandle_t, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamInput_StopAnalogActionMomentum", expected: (func(uintptr, InputHandle_t, InputAnalogActionHandle_t))(nil)},
		{name: "ptrAPI_ISteamInput_GetMotionData", expected: uintptr(0)},
		{name: "ptrAPI_ISteamInput_TriggerVibration", expected: (func(uintptr, InputHandle_t, uint16, uint16))(nil)},
		{name: "ptrAPI_ISteamInput_TriggerVibrationExtended", expected: (func(uintptr, InputHandle_t, uint16, uint16, uint16, uint16))(nil)},
		{name: "ptrAPI_ISteamInput_TriggerSimpleHapticEvent", expected: (func(uintptr, InputHandle_t, ESteamControllerPad, uint16, uint16, uint16))(nil)},
		{name: "ptrAPI_ISteamInput_SetLEDColor", expected: (func(uintptr, InputHandle_t, uint8, uint8, uint8, ESteamInputLEDFlag))(nil)},
		{name: "ptrAPI_ISteamInput_ShowBindingPanel", expected: (func(uintptr, InputHandle_t) bool)(nil)},
		{name: "ptrAPI_ISteamInput_GetControllerForGamepadIndex", expected: (func(uintptr, int32) InputHandle_t)(nil)},
		{name: "ptrAPI_ISteamInput_GetGamepadIndexForController", expected: (func(uintptr, InputHandle_t) int32)(nil)},
		{name: "ptrAPI_ISteamInput_GetStringForActionOrigin", expected: (func(uintptr, EInputActionOrigin) string)(nil)},
		{name: "ptrAPI_ISteamInput_GetGlyphForActionOrigin", expected: (func(uintptr, EInputActionOrigin) string)(nil)},
		{name: "ptrAPI_ISteamInput_GetRemotePlaySessionID", expected: (func(uintptr, InputHandle_t) uint32)(nil)},

		{name: "ptrAPI_SteamRemotePlay", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamRemotePlay_BSessionRemotePlayTogether", expected: (func(uintptr, uint32) bool)(nil)},
		{name: "ptrAPI_ISteamRemotePlay_GetSessionGuestID", expected: (func(uintptr, uint32) uint32)(nil)},
		{name: "ptrAPI_ISteamRemotePlay_GetSmallSessionAvatar", expected: (func(uintptr, uint32) int32)(nil)},
		{name: "ptrAPI_ISteamRemotePlay_GetMediumSessionAvatar", expected: (func(uintptr, uint32) int32)(nil)},
		{name: "ptrAPI_ISteamRemotePlay_GetLargeSessionAvatar", expected: (func(uintptr, uint32) int32)(nil)},

		{name: "ptrAPI_SteamRemoteStorage", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamRemoteStorage_FileWrite", expected: (func(uintptr, string, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamRemoteStorage_FileRead", expected: (func(uintptr, string, uintptr, int32) int32)(nil)},
		{name: "ptrAPI_ISteamRemoteStorage_FileDelete", expected: (func(uintptr, string) bool)(nil)},
		{name: "ptrAPI_ISteamRemoteStorage_GetFileSize", expected: (func(uintptr, string) int32)(nil)},

		{name: "ptrAPI_SteamUser", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamUser_AdvertiseGame", expected: (func(uintptr, CSteamID, uint32, uint16))(nil)},
		{name: "ptrAPI_ISteamUser_BeginAuthSession", expected: (func(uintptr, uintptr, int32, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamUser_BIsBehindNAT", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_BIsPhoneIdentifying", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_BIsPhoneRequiringVerification", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_BIsPhoneVerified", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_BIsTwoFactorEnabled", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_BLoggedOn", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_BSetDurationControlOnlineState", expected: (func(uintptr, EDurationControlOnlineState) bool)(nil)},
		{name: "ptrAPI_ISteamUser_CancelAuthTicket", expected: (func(uintptr, HAuthTicket))(nil)},
		{name: "ptrAPI_ISteamUser_DecompressVoice", expected: (func(uintptr, uintptr, uint32, uintptr, uint32, uintptr, uint32) int32)(nil)},
		{name: "ptrAPI_ISteamUser_EndAuthSession", expected: (func(uintptr, CSteamID))(nil)},
		{name: "ptrAPI_ISteamUser_GetAuthSessionTicket", expected: (func(uintptr, uintptr, int32, uintptr, uintptr) HAuthTicket)(nil)},
		{name: "ptrAPI_ISteamUser_GetAuthTicketForWebApi", expected: (func(uintptr, string, uintptr) HAuthTicket)(nil)},
		{name: "ptrAPI_ISteamUser_GetAvailableVoice", expected: (func(uintptr, uintptr, uintptr, uint32) int32)(nil)},
		{name: "ptrAPI_ISteamUser_GetDurationControl", expected: (func(uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_GetEncryptedAppTicket", expected: (func(uintptr, uintptr, int32, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUser_GetGameBadgeLevel", expected: (func(uintptr, int32, bool) int32)(nil)},
		{name: "ptrAPI_ISteamUser_GetHSteamUser", expected: (func(uintptr) HSteamUser)(nil)},
		{name: "ptrAPI_ISteamUser_GetPlayerSteamLevel", expected: (func(uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamUser_GetSteamID", expected: (func(uintptr) CSteamID)(nil)},
		{name: "ptrAPI_ISteamUser_GetUserDataFolder", expected: (func(uintptr, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamUser_GetVoice", expected: (func(uintptr, bool, uintptr, uint32, uintptr, bool, uintptr, uint32, uintptr, uint32) int32)(nil)},
		{name: "ptrAPI_ISteamUser_GetVoiceOptimalSampleRate", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUser_InitiateGameConnection", expected: (func(uintptr, uintptr, int32, CSteamID, uint32, uint16, bool) int32)(nil)},
		{name: "ptrAPI_ISteamUser_RequestEncryptedAppTicket", expected: (func(uintptr, uintptr, int32) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamUser_RequestStoreAuthURL", expected: (func(uintptr, string) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamUser_StartVoiceRecording", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamUser_StopVoiceRecording", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamUser_TerminateGameConnection", expected: (func(uintptr, uint32, uint16))(nil)},
		{name: "ptrAPI_ISteamUser_TrackAppUsageEvent", expected: (func(uintptr, CGameID, int32, string))(nil)},
		{name: "ptrAPI_ISteamUser_UserHasLicenseForApp", expected: (func(uintptr, CSteamID, AppId_t) int32)(nil)},

		{name: "ptrAPI_SteamUserStats", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamUserStats_GetAchievement", expected: (func(uintptr, string, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUserStats_SetAchievement", expected: (func(uintptr, string) bool)(nil)},
		{name: "ptrAPI_ISteamUserStats_ClearAchievement", expected: (func(uintptr, string) bool)(nil)},
		{name: "ptrAPI_ISteamUserStats_StoreStats", expected: (func(uintptr) bool)(nil)},

		{name: "ptrAPI_SteamUtils", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamUtils_GetSecondsSinceAppActive", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUtils_GetSecondsSinceComputerActive", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUtils_GetConnectedUniverse", expected: (func(uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamUtils_GetServerRealTime", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUtils_GetIPCountry", expected: (func(uintptr) string)(nil)},
		{name: "ptrAPI_ISteamUtils_GetImageSize", expected: (func(uintptr, int32, uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_GetImageRGBA", expected: (func(uintptr, int32, uintptr, int32) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_GetCurrentBatteryPower", expected: (func(uintptr) uint8)(nil)},
		{name: "ptrAPI_ISteamUtils_GetAppID", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUtils_SetOverlayNotificationPosition", expected: (func(uintptr, ENotificationPosition))(nil)},
		{name: "ptrAPI_ISteamUtils_IsAPICallCompleted", expected: (func(uintptr, SteamAPICall_t, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_GetAPICallFailureReason", expected: (func(uintptr, SteamAPICall_t) int32)(nil)},
		{name: "ptrAPI_ISteamUtils_GetAPICallResult", expected: (func(uintptr, SteamAPICall_t, uintptr, int32, int32, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_GetIPCCallCount", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamUtils_IsOverlayEnabled", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_BOverlayNeedsPresent", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput", expected: (func(uintptr, EFloatingGamepadTextInputMode, int32, int32, int32, int32) bool)(nil)},
		{name: "ptrAPI_ISteamUtils_SetOverlayNotificationInset", expected: (func(uintptr, int32, int32))(nil)},

		{name: "ptrAPI_SteamNetworkingUtils", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamNetworkingUtils_AllocateMessage", expected: (func(uintptr, int32) uintptr)(nil)},
		{name: "ptrAPI_ISteamNetworkingUtils_InitRelayNetworkAccess", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamNetworkingUtils_GetLocalTimestamp", expected: (func(uintptr) SteamNetworkingMicroseconds)(nil)},

		{name: "ptrAPI_SteamGameServer", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamGameServer_AssociateWithClan", expected: (func(uintptr, CSteamID) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamGameServer_BeginAuthSession", expected: (func(uintptr, uintptr, int32, CSteamID) int32)(nil)},
		{name: "ptrAPI_ISteamGameServer_BLoggedOn", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_BSecure", expected: (func(uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_BUpdateUserData", expected: (func(uintptr, CSteamID, string, uint32) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_CancelAuthTicket", expected: (func(uintptr, HAuthTicket))(nil)},
		{name: "ptrAPI_ISteamGameServer_ClearAllKeyValues", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamGameServer_ComputeNewPlayerCompatibility", expected: (func(uintptr, CSteamID, uintptr, uint32, uintptr, uint32, uintptr, uint32) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamGameServer_CreateUnauthenticatedUserConnection", expected: (func(uintptr) CSteamID)(nil)},
		{name: "ptrAPI_ISteamGameServer_EnableHeartbeats", expected: (func(uintptr, bool))(nil)},
		{name: "ptrAPI_ISteamGameServer_EndAuthSession", expected: (func(uintptr, CSteamID))(nil)},
		{name: "ptrAPI_ISteamGameServer_ForceHeartbeat", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamGameServer_GetAuthSessionTicket", expected: (func(uintptr, uintptr, int32, uintptr) HAuthTicket)(nil)},
		{name: "ptrAPI_ISteamGameServer_GetGameplayStats", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamGameServer_GetNextOutgoingPacket", expected: (func(uintptr, uintptr, int32, uintptr, uintptr) int32)(nil)},
		{name: "ptrAPI_ISteamGameServer_GetPublicIP", expected: (func(uintptr) uint32)(nil)},
		{name: "ptrAPI_ISteamGameServer_GetServerReputation", expected: (func(uintptr) SteamAPICall_t)(nil)},
		{name: "ptrAPI_ISteamGameServer_GetSteamID", expected: (func(uintptr) CSteamID)(nil)},
		{name: "ptrAPI_ISteamGameServer_HandleIncomingPacket", expected: (func(uintptr, uintptr, int32, uint32, uint16) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_InitGameServer", expected: (func(uintptr, uint32, uint16, uint16, uint16, uint32, string) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_LogOff", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamGameServer_LogOn", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_LogOnAnonymous", expected: (func(uintptr))(nil)},
		{name: "ptrAPI_ISteamGameServer_RequestUserGroupStatus", expected: (func(uintptr, CSteamID, CSteamID) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_SendUserConnectAndAuthenticate", expected: (func(uintptr, uint32, uintptr, uint32, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamGameServer_SendUserDisconnect", expected: (func(uintptr, CSteamID))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetBotPlayerCount", expected: (func(uintptr, int32))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetDedicatedServer", expected: (func(uintptr, bool))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetGameData", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetGameDescription", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetGameTags", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetHeartbeatInterval", expected: (func(uintptr, int32))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetKeyValue", expected: (func(uintptr, string, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetMapName", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetMaxPlayerCount", expected: (func(uintptr, int32))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetModDir", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetPasswordProtected", expected: (func(uintptr, bool))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetProduct", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetRegion", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetServerName", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetSpectatorPort", expected: (func(uintptr, uint16))(nil)},
		{name: "ptrAPI_ISteamGameServer_SetSpectatorServerName", expected: (func(uintptr, string))(nil)},
		{name: "ptrAPI_ISteamGameServer_UserHasLicenseForApp", expected: (func(uintptr, CSteamID, AppId_t) int32)(nil)},
		{name: "ptrAPI_ISteamGameServer_WasRestartRequested", expected: (func(uintptr) bool)(nil)},

		{name: "ptrAPI_SteamNetworkingMessages", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamNetworkingMessages_SendMessageToUser", expected: (func(uintptr, uintptr, uintptr, uint32, int32, int32) EResult)(nil)},
		{name: "ptrAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel", expected: (func(uintptr, int32, uintptr, int32) int32)(nil)},
		{name: "ptrAPI_ISteamNetworkingMessages_AcceptSessionWithUser", expected: (func(uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamNetworkingMessages_CloseSessionWithUser", expected: (func(uintptr, uintptr) bool)(nil)},
		{name: "ptrAPI_ISteamNetworkingMessages_CloseChannelWithUser", expected: (func(uintptr, uintptr, int32) bool)(nil)},

		{name: "ptrAPI_SteamNetworkingSockets", expected: (func() uintptr)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_CreateListenSocketIP", expected: (func(uintptr, uintptr, int32, uintptr) HSteamListenSocket)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P", expected: (func(uintptr, int32, int32, uintptr) HSteamListenSocket)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_ConnectByIPAddress", expected: (func(uintptr, uintptr, int32, uintptr) HSteamNetConnection)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_ConnectP2P", expected: (func(uintptr, uintptr, int32, int32, uintptr) HSteamNetConnection)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_AcceptConnection", expected: (func(uintptr, HSteamNetConnection) EResult)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_CloseConnection", expected: (func(uintptr, HSteamNetConnection, int32, string, bool) bool)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_CloseListenSocket", expected: (func(uintptr, HSteamListenSocket) bool)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_SendMessageToConnection", expected: (func(uintptr, HSteamNetConnection, uintptr, uint32, int32, uintptr) EResult)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection", expected: (func(uintptr, HSteamNetConnection, uintptr, int32) int32)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_CreatePollGroup", expected: (func(uintptr) HSteamNetPollGroup)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_DestroyPollGroup", expected: (func(uintptr, HSteamNetPollGroup) bool)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_SetConnectionPollGroup", expected: (func(uintptr, HSteamNetConnection, HSteamNetPollGroup) bool)(nil)},
		{name: "ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup", expected: (func(uintptr, HSteamNetPollGroup, uintptr, int32) int32)(nil)},
	}
}

func allFlatAPISymbols() []string {
	return []string{
		flatAPI_RestartAppIfNecessary,
		flatAPI_InitFlat,
		flatAPI_RunCallbacks,
		flatAPI_Shutdown,
		flatAPI_IsSteamRunning,
		flatAPI_GetSteamInstallPath,
		flatAPI_ReleaseCurrentThreadMemory,

		flatAPI_SteamApps,
		flatAPI_ISteamApps_BIsSubscribed,
		flatAPI_ISteamApps_BIsLowViolence,
		flatAPI_ISteamApps_BIsCybercafe,
		flatAPI_ISteamApps_BIsVACBanned,
		flatAPI_ISteamApps_BGetDLCDataByIndex,
		flatAPI_ISteamApps_BIsDlcInstalled,
		flatAPI_ISteamApps_GetAvailableGameLanguages,
		flatAPI_ISteamApps_BIsSubscribedApp,
		flatAPI_ISteamApps_GetEarliestPurchaseUnixTime,
		flatAPI_ISteamApps_BIsSubscribedFromFreeWeekend,
		flatAPI_ISteamApps_GetAppInstallDir,
		flatAPI_ISteamApps_GetCurrentGameLanguage,
		flatAPI_ISteamApps_GetDLCCount,
		flatAPI_ISteamApps_InstallDLC,
		flatAPI_ISteamApps_UninstallDLC,
		flatAPI_ISteamApps_RequestAppProofOfPurchaseKey,
		flatAPI_ISteamApps_GetCurrentBetaName,
		flatAPI_ISteamApps_MarkContentCorrupt,
		flatAPI_ISteamApps_GetInstalledDepots,
		flatAPI_ISteamApps_BIsAppInstalled,
		flatAPI_ISteamApps_GetAppOwner,
		flatAPI_ISteamApps_GetLaunchQueryParam,
		flatAPI_ISteamApps_GetDlcDownloadProgress,
		flatAPI_ISteamApps_GetAppBuildId,
		flatAPI_ISteamApps_RequestAllProofOfPurchaseKeys,
		flatAPI_ISteamApps_GetFileDetails,
		flatAPI_ISteamApps_GetLaunchCommandLine,
		flatAPI_ISteamApps_BIsSubscribedFromFamilySharing,
		flatAPI_ISteamApps_BIsTimedTrial,
		flatAPI_ISteamApps_SetDlcContext,
		flatAPI_ISteamApps_GetNumBetas,
		flatAPI_ISteamApps_GetBetaInfo,
		flatAPI_ISteamApps_SetActiveBeta,

		flatAPI_SteamFriends,
		flatAPI_ISteamFriends_GetPersonaName,
		flatAPI_ISteamFriends_GetPersonaState,
		flatAPI_ISteamFriends_GetFriendCount,
		flatAPI_ISteamFriends_GetFriendByIndex,
		flatAPI_ISteamFriends_GetFriendRelationship,
		flatAPI_ISteamFriends_GetFriendPersonaState,
		flatAPI_ISteamFriends_GetFriendPersonaName,
		flatAPI_ISteamFriends_GetFriendPersonaNameHistory,
		flatAPI_ISteamFriends_GetFriendSteamLevel,
		flatAPI_ISteamFriends_GetSmallFriendAvatar,
		flatAPI_ISteamFriends_GetMediumFriendAvatar,
		flatAPI_ISteamFriends_GetLargeFriendAvatar,
		flatAPI_ISteamFriends_SetRichPresence,
		flatAPI_ISteamFriends_GetFriendGamePlayed,
		flatAPI_ISteamFriends_InviteUserToGame,
		flatAPI_ISteamFriends_ActivateGameOverlay,
		flatAPI_ISteamFriends_ActivateGameOverlayToUser,
		flatAPI_ISteamFriends_ActivateGameOverlayToWebPage,
		flatAPI_ISteamFriends_ActivateGameOverlayToStore,
		flatAPI_ISteamFriends_ActivateGameOverlayInviteDialog,
		flatAPI_ISteamFriends_ActivateGameOverlayInviteDialogConnectString,

		flatAPI_SteamMatchmaking,
		flatAPI_ISteamMatchmaking_GetFavoriteGameCount,
		flatAPI_ISteamMatchmaking_GetFavoriteGame,
		flatAPI_ISteamMatchmaking_AddFavoriteGame,
		flatAPI_ISteamMatchmaking_RemoveFavoriteGame,
		flatAPI_ISteamMatchmaking_RequestLobbyList,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListStringFilter,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListNumericalFilter,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListNearValueFilter,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListFilterSlotsAvailable,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListDistanceFilter,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListResultCountFilter,
		flatAPI_ISteamMatchmaking_AddRequestLobbyListCompatibleMembersFilter,
		flatAPI_ISteamMatchmaking_GetLobbyByIndex,
		flatAPI_ISteamMatchmaking_CreateLobby,
		flatAPI_ISteamMatchmaking_JoinLobby,
		flatAPI_ISteamMatchmaking_LeaveLobby,
		flatAPI_ISteamMatchmaking_InviteUserToLobby,
		flatAPI_ISteamMatchmaking_GetLobbyMemberLimit,
		flatAPI_ISteamMatchmaking_SetLobbyMemberLimit,
		flatAPI_ISteamMatchmaking_SetLobbyType,
		flatAPI_ISteamMatchmaking_SetLobbyJoinable,
		flatAPI_ISteamMatchmaking_GetLobbyOwner,
		flatAPI_ISteamMatchmaking_SetLobbyOwner,
		flatAPI_ISteamMatchmaking_SetLinkedLobby,
		flatAPI_ISteamMatchmaking_GetNumLobbyMembers,
		flatAPI_ISteamMatchmaking_GetLobbyMemberByIndex,
		flatAPI_ISteamMatchmaking_SetLobbyData,
		flatAPI_ISteamMatchmaking_GetLobbyData,
		flatAPI_ISteamMatchmaking_DeleteLobbyData,
		flatAPI_ISteamMatchmaking_GetLobbyDataCount,
		flatAPI_ISteamMatchmaking_GetLobbyDataByIndex,
		flatAPI_ISteamMatchmaking_SetLobbyMemberData,
		flatAPI_ISteamMatchmaking_GetLobbyMemberData,
		flatAPI_ISteamMatchmaking_SendLobbyChatMsg,
		flatAPI_ISteamMatchmaking_GetLobbyChatEntry,
		flatAPI_ISteamMatchmaking_RequestLobbyData,
		flatAPI_ISteamMatchmaking_SetLobbyGameServer,
		flatAPI_ISteamMatchmaking_GetLobbyGameServer,
		flatAPI_ISteamMatchmaking_CheckForPSNGameBootInvite,

		flatAPI_SteamMatchmakingServers_RequestInternetServerList,
		flatAPI_SteamMatchmakingServers_RequestLANServerList,
		flatAPI_SteamMatchmakingServers_RequestFriendsServerList,
		flatAPI_SteamMatchmakingServers_RequestFavoritesServerList,
		flatAPI_SteamMatchmakingServers_RequestHistoryServerList,
		flatAPI_SteamMatchmakingServers_RequestSpectatorServerList,
		flatAPI_SteamMatchmakingServers_ReleaseRequest,
		flatAPI_SteamMatchmakingServers_GetServerDetails,
		flatAPI_SteamMatchmakingServers_CancelQuery,
		flatAPI_SteamMatchmakingServers_RefreshQuery,
		flatAPI_SteamMatchmakingServers_IsRefreshing,
		flatAPI_SteamMatchmakingServers_GetServerCount,
		flatAPI_SteamMatchmakingServers_RefreshServer,
		flatAPI_SteamMatchmakingServers_PingServer,
		flatAPI_SteamMatchmakingServers_PlayerDetails,
		flatAPI_SteamMatchmakingServers_ServerRules,
		flatAPI_SteamMatchmakingServers_CancelServerQuery,

		flatAPI_SteamHTTP,
		flatAPI_ISteamHTTP_CreateHTTPRequest,
		flatAPI_ISteamHTTP_SetHTTPRequestHeaderValue,
		flatAPI_ISteamHTTP_SendHTTPRequest,
		flatAPI_ISteamHTTP_GetHTTPResponseBodySize,
		flatAPI_ISteamHTTP_GetHTTPResponseBodyData,
		flatAPI_ISteamHTTP_ReleaseHTTPRequest,

		flatAPI_SteamUGC,
		flatAPI_ISteamUGC_GetNumSubscribedItems,
		flatAPI_ISteamUGC_GetSubscribedItems,
		flatAPI_ISteamUGC_MarkDownloadedItemAsUnused,
		flatAPI_ISteamUGC_GetNumDownloadedItems,
		flatAPI_ISteamUGC_GetDownloadedItems,

		flatAPI_SteamInventory,
		flatAPI_ISteamInventory_GetResultStatus,
		flatAPI_ISteamInventory_GetResultItems,
		flatAPI_ISteamInventory_DestroyResult,

		flatAPI_SteamInput,
		flatAPI_ISteamInput_GetConnectedControllers,
		flatAPI_ISteamInput_GetInputTypeForHandle,
		flatAPI_ISteamInput_Init,
		flatAPI_ISteamInput_Shutdown,
		flatAPI_ISteamInput_RunFrame,
		flatAPI_ISteamInput_EnableDeviceCallbacks,
		flatAPI_ISteamInput_GetActionSetHandle,
		flatAPI_ISteamInput_ActivateActionSet,
		flatAPI_ISteamInput_GetCurrentActionSet,
		flatAPI_ISteamInput_ActivateActionSetLayer,
		flatAPI_ISteamInput_DeactivateActionSetLayer,
		flatAPI_ISteamInput_DeactivateAllActionSetLayers,
		flatAPI_ISteamInput_GetActiveActionSetLayers,
		flatAPI_ISteamInput_GetDigitalActionHandle,
		flatAPI_ISteamInput_GetDigitalActionData,
		flatAPI_ISteamInput_GetDigitalActionOrigins,
		flatAPI_ISteamInput_GetAnalogActionHandle,
		flatAPI_ISteamInput_GetAnalogActionData,
		flatAPI_ISteamInput_GetAnalogActionOrigins,
		flatAPI_ISteamInput_StopAnalogActionMomentum,
		flatAPI_ISteamInput_GetMotionData,
		flatAPI_ISteamInput_TriggerVibration,
		flatAPI_ISteamInput_TriggerVibrationExtended,
		flatAPI_ISteamInput_TriggerSimpleHapticEvent,
		flatAPI_ISteamInput_SetLEDColor,
		flatAPI_ISteamInput_ShowBindingPanel,
		flatAPI_ISteamInput_GetControllerForGamepadIndex,
		flatAPI_ISteamInput_GetGamepadIndexForController,
		flatAPI_ISteamInput_GetStringForActionOrigin,
		flatAPI_ISteamInput_GetGlyphForActionOrigin,
		flatAPI_ISteamInput_GetRemotePlaySessionID,

		flatAPI_SteamRemotePlay,
		flatAPI_ISteamRemotePlay_BSessionRemotePlayTogether,
		flatAPI_ISteamRemotePlay_GetSessionGuestID,
		flatAPI_ISteamRemotePlay_GetSmallSessionAvatar,
		flatAPI_ISteamRemotePlay_GetMediumSessionAvatar,
		flatAPI_ISteamRemotePlay_GetLargeSessionAvatar,

		flatAPI_SteamRemoteStorage,
		flatAPI_ISteamRemoteStorage_FileWrite,
		flatAPI_ISteamRemoteStorage_FileRead,
		flatAPI_ISteamRemoteStorage_FileDelete,
		flatAPI_ISteamRemoteStorage_GetFileSize,

		flatAPI_SteamUser,
		flatAPI_ISteamUser_GetSteamID,

		flatAPI_SteamUserStats,
		flatAPI_ISteamUserStats_GetAchievement,
		flatAPI_ISteamUserStats_SetAchievement,
		flatAPI_ISteamUserStats_ClearAchievement,
		flatAPI_ISteamUserStats_StoreStats,

		flatAPI_SteamUtils,
		flatAPI_ISteamUtils_GetSecondsSinceAppActive,
		flatAPI_ISteamUtils_GetSecondsSinceComputerActive,
		flatAPI_ISteamUtils_GetConnectedUniverse,
		flatAPI_ISteamUtils_GetServerRealTime,
		flatAPI_ISteamUtils_GetIPCountry,
		flatAPI_ISteamUtils_GetImageSize,
		flatAPI_ISteamUtils_GetImageRGBA,
		flatAPI_ISteamUtils_GetCurrentBatteryPower,
		flatAPI_ISteamUtils_GetAppID,
		flatAPI_ISteamUtils_SetOverlayNotificationPosition,
		flatAPI_ISteamUtils_IsAPICallCompleted,
		flatAPI_ISteamUtils_GetAPICallFailureReason,
		flatAPI_ISteamUtils_GetAPICallResult,
		flatAPI_ISteamUtils_GetIPCCallCount,
		flatAPI_ISteamUtils_IsOverlayEnabled,
		flatAPI_ISteamUtils_BOverlayNeedsPresent,
		flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck,
		flatAPI_ISteamUtils_ShowFloatingGamepadTextInput,
		flatAPI_ISteamUtils_SetOverlayNotificationInset,

		flatAPI_SteamNetworkingUtils,
		flatAPI_ISteamNetworkingUtils_AllocateMessage,
		flatAPI_ISteamNetworkingUtils_InitRelayNetworkAccess,
		flatAPI_ISteamNetworkingUtils_GetLocalTimestamp,

		flatAPI_SteamNetworkingMessages,
		flatAPI_ISteamNetworkingMessages_SendMessageToUser,
		flatAPI_ISteamNetworkingMessages_ReceiveMessagesOnChannel,
		flatAPI_ISteamNetworkingMessages_AcceptSessionWithUser,
		flatAPI_ISteamNetworkingMessages_CloseSessionWithUser,
		flatAPI_ISteamNetworkingMessages_CloseChannelWithUser,

		flatAPI_SteamNetworkingSockets,
		flatAPI_ISteamNetworkingSockets_CreateListenSocketIP,
		flatAPI_ISteamNetworkingSockets_CreateListenSocketP2P,
		flatAPI_ISteamNetworkingSockets_ConnectByIPAddress,
		flatAPI_ISteamNetworkingSockets_ConnectP2P,
		flatAPI_ISteamNetworkingSockets_AcceptConnection,
		flatAPI_ISteamNetworkingSockets_CloseConnection,
		flatAPI_ISteamNetworkingSockets_CloseListenSocket,
		flatAPI_ISteamNetworkingSockets_SendMessageToConnection,
		flatAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection,
		flatAPI_ISteamNetworkingSockets_CreatePollGroup,
		flatAPI_ISteamNetworkingSockets_DestroyPollGroup,
		flatAPI_ISteamNetworkingSockets_SetConnectionPollGroup,
		flatAPI_ISteamNetworkingSockets_ReceiveMessagesOnPollGroup,

		flatAPI_SteamGameServer,
		flatAPI_ISteamGameServer_AssociateWithClan,
		flatAPI_ISteamGameServer_BeginAuthSession,
		flatAPI_ISteamGameServer_BLoggedOn,
		flatAPI_ISteamGameServer_BSecure,
		flatAPI_ISteamGameServer_BUpdateUserData,
		flatAPI_ISteamGameServer_CancelAuthTicket,
		flatAPI_ISteamGameServer_ClearAllKeyValues,
		flatAPI_ISteamGameServer_ComputeNewPlayerCompatibility,
		flatAPI_ISteamGameServer_CreateUnauthenticatedUserConnection,
		flatAPI_ISteamGameServer_EnableHeartbeats,
		flatAPI_ISteamGameServer_EndAuthSession,
		flatAPI_ISteamGameServer_ForceHeartbeat,
		flatAPI_ISteamGameServer_GetAuthSessionTicket,
		flatAPI_ISteamGameServer_GetGameplayStats,
		flatAPI_ISteamGameServer_GetNextOutgoingPacket,
		flatAPI_ISteamGameServer_GetPublicIP,
		flatAPI_ISteamGameServer_GetServerReputation,
		flatAPI_ISteamGameServer_GetSteamID,
		flatAPI_ISteamGameServer_HandleIncomingPacket,
		flatAPI_ISteamGameServer_InitGameServer,
		flatAPI_ISteamGameServer_LogOff,
		flatAPI_ISteamGameServer_LogOn,
		flatAPI_ISteamGameServer_LogOnAnonymous,
		flatAPI_ISteamGameServer_RequestUserGroupStatus,
		flatAPI_ISteamGameServer_SendUserConnectAndAuthenticate,
		flatAPI_ISteamGameServer_SendUserDisconnect,
		flatAPI_ISteamGameServer_SetBotPlayerCount,
		flatAPI_ISteamGameServer_SetDedicatedServer,
		flatAPI_ISteamGameServer_SetGameData,
		flatAPI_ISteamGameServer_SetGameDescription,
		flatAPI_ISteamGameServer_SetGameTags,
		flatAPI_ISteamGameServer_SetHeartbeatInterval,
		flatAPI_ISteamGameServer_SetKeyValue,
		flatAPI_ISteamGameServer_SetMapName,
		flatAPI_ISteamGameServer_SetMaxPlayerCount,
		flatAPI_ISteamGameServer_SetModDir,
		flatAPI_ISteamGameServer_SetPasswordProtected,
		flatAPI_ISteamGameServer_SetProduct,
		flatAPI_ISteamGameServer_SetRegion,
		flatAPI_ISteamGameServer_SetServerName,
		flatAPI_ISteamGameServer_SetSpectatorPort,
		flatAPI_ISteamGameServer_SetSpectatorServerName,
		flatAPI_ISteamGameServer_UserHasLicenseForApp,
		flatAPI_ISteamGameServer_WasRestartRequested,
	}
}
