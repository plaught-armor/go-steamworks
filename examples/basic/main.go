// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The go-steamworks Authors

package main

import (
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

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		steamworks.RunCallbacks()
	}
}
