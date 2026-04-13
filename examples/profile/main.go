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

	steamID := steamworks.SteamUser().GetSteamID()
	personaName := steamworks.SteamFriends().GetPersonaName()
	log.Printf("Signed in as %s (%d)", personaName, steamID)

	avatarHandle, err := fetchLargeAvatarHandle(steamID, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to get avatar handle: %v", err)
	}
	if avatarHandle == 0 {
		log.Fatal("no avatar available for this user")
	}

	width, height, ok := steamworks.SteamUtils().GetImageSize(int(avatarHandle))
	if !ok {
		log.Fatal("failed to get avatar image size")
	}
	pixels := make([]byte, int(width*height*4))
	if ok := steamworks.SteamUtils().GetImageRGBA(int(avatarHandle), pixels); !ok {
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

func fetchLargeAvatarHandle(steamID steamworks.CSteamID, timeout time.Duration) (int32, error) {
	deadline := time.Now().Add(timeout)
	for {
		avatar := steamworks.SteamFriends().GetLargeFriendAvatar(steamID)
		if avatar != -1 {
			return avatar, nil
		}

		if time.Now().After(deadline) {
			return 0, errors.New("timed out waiting for avatar")
		}

		steamworks.RunCallbacks()
		time.Sleep(50 * time.Millisecond)
	}
}
