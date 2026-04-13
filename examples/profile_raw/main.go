// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The go-steamworks Authors

package main

import (
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/plaught-armor/go-steamworks"
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

	userPtr := mustSymbol("SteamAPI_SteamUser_v023")
	friendsPtr := mustSymbol("SteamAPI_SteamFriends_v018")
	utilsPtr := mustSymbol("SteamAPI_SteamUtils_v010")

	steamID := steamworks.CSteamID(callSymbol("SteamAPI_ISteamUser_GetSteamID", userPtr))
	personaName := cStringToGo(callSymbol("SteamAPI_ISteamFriends_GetPersonaName", friendsPtr))
	personaState := steamworks.EPersonaState(callSymbol("SteamAPI_ISteamFriends_GetPersonaState", friendsPtr))
	steamLevel := int(callSymbol("SteamAPI_ISteamFriends_GetFriendSteamLevel", friendsPtr, uintptr(steamID)))

	log.Printf("Signed in as %s (%d)", personaName, steamID)
	log.Printf("Persona state: %v", personaState)
	log.Printf("Steam level: %d", steamLevel)

	avatarHandle, err := fetchLargeAvatarHandle(friendsPtr, steamID, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to get avatar handle: %v", err)
	}
	if avatarHandle == 0 {
		log.Fatal("no avatar available for this user")
	}

	width, height, ok := getImageSize(utilsPtr, avatarHandle)
	if !ok {
		log.Fatal("failed to get avatar image size")
	}

	pixels := make([]byte, int(width*height*4))
	if ok := getImageRGBA(utilsPtr, avatarHandle, pixels); !ok {
		log.Fatal("failed to read avatar image data")
	}

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	copy(img.Pix, pixels)

	file, err := os.Create("avatar.png")
	if err != nil {
		log.Fatalf("failed to create avatar.png: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Fatalf("failed to write avatar.png: %v", err)
	}

	log.Printf("Saved avatar for %s to avatar.png", personaName)
}

func mustSymbol(name string) uintptr {
	ptr, err := steamworks.CallSymbol(name)
	if err != nil {
		log.Fatalf("failed to resolve %s: %v", name, err)
	}
	return ptr
}

func callSymbol(name string, args ...uintptr) uintptr {
	result, err := steamworks.CallSymbol(name, args...)
	if err != nil {
		log.Fatalf("failed to call %s: %v", name, err)
	}
	return result
}

func cStringToGo(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	buffer := make([]byte, 0, 64)
	for i := uintptr(0); ; i++ {
		b := *(*byte)(unsafe.Pointer(ptr + i))
		if b == 0 {
			break
		}
		buffer = append(buffer, b)
	}
	return string(buffer)
}

func fetchLargeAvatarHandle(friendsPtr uintptr, steamID steamworks.CSteamID, timeout time.Duration) (int32, error) {
	deadline := time.Now().Add(timeout)
	for {
		avatarHandle := int32(callSymbol(
			"SteamAPI_ISteamFriends_GetLargeFriendAvatar",
			friendsPtr,
			uintptr(steamID),
		))

		if avatarHandle != -1 {
			return avatarHandle, nil
		}

		if time.Now().After(deadline) {
			return 0, errors.New("timed out waiting for avatar")
		}

		steamworks.RunCallbacks()
		time.Sleep(50 * time.Millisecond)
	}
}

func getImageSize(utilsPtr uintptr, imageHandle int32) (uint32, uint32, bool) {
	var width uint32
	var height uint32
	ok := callSymbol(
		"SteamAPI_ISteamUtils_GetImageSize",
		utilsPtr,
		uintptr(imageHandle),
		uintptr(unsafe.Pointer(&width)),
		uintptr(unsafe.Pointer(&height)),
	)
	return width, height, ok != 0
}

func getImageRGBA(utilsPtr uintptr, imageHandle int32, dest []byte) bool {
	if len(dest) == 0 {
		return false
	}
	ok := callSymbol(
		"SteamAPI_ISteamUtils_GetImageRGBA",
		utilsPtr,
		uintptr(imageHandle),
		uintptr(unsafe.Pointer(&dest[0])),
		uintptr(len(dest)),
	)
	return ok != 0
}
