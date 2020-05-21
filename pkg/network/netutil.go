// Package network contains network related utilities
package network

import (
	"net"
	"strings"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

// GetExternalAddress will return the detected external IP address based on the env var, then network interfaces
// (it will look for the first 172.17.0.0/16 address), and finally a fallback on localhost
// TODO: Maybe check for docker subnet first somehow, people can change it from 172.17.0.0/16 to something else after all..
func GetExternalAddress() (string, error) {
	// Check for IP in interface addresses
	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return "", errors.Wrap(err, "failed to get interface addresses")
	}

	// Get first 172.17.0.0/16 address, if any
	for _, address := range addresses {
		if strings.Contains(address.String(), env.DockerAddressPrefix) {
			ip := strings.Split(address.String(), "/")[0]

			return ip, nil
		}
	}

	// Default to localhost
	return "localhost", nil
}
