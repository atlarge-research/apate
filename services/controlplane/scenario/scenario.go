// Package scenario contains some scenario specific utilities
package scenario

import "log"

// This file should be called something else and should be situated somewhere else
// Right now there's no better place though.

// Failed handles the failure of a scenario
func Failed(reason error) {
	// Just print for now
	// TODO figure out what to do exactly
	log.Print(reason)
}
