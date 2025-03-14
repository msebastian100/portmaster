package network

import (
	"github.com/safing/portbase/log"
	"github.com/safing/portbase/rng"
)

// GetUnusedLocalPort returns a local port of the specified protocol that is
// currently unused and is unlikely to be used within the next seconds.
func GetUnusedLocalPort(protocol uint8) (port uint16, ok bool) {
	allConns := conns.clone()
	tries := 1000

	// Try up to 1000 times to find an unused port.
nextPort:
	for i := 0; i < tries; i++ {
		// Generate random port between 10000 and 65535
		rN, err := rng.Number(55535)
		if err != nil {
			log.Warningf("network: failed to generate random port: %s", err)
			return 0, false
		}
		port := uint16(rN + 10000)

		// Shrink range when we chew through the tries.
		portRangeStart := port - 10

		// Check if the generated port is unused.
	nextConnection:
		for _, conn := range allConns {
			// Skip connection if the protocol does not match the protocol of interest.
			if conn.Entity.Protocol != protocol {
				continue nextConnection
			}
			// Skip port if the local port is in dangerous proximity.
			// Consecutive port numbers are very common.
			if conn.LocalPort <= port && conn.LocalPort >= portRangeStart {
				continue nextPort
			}
		}

		// Log if it took more than 10 attempts.
		if i >= 10 {
			log.Warningf("network: took %d attempts to find a suitable unused port for pre-auth", i+1)
		}

		// The checks have passed. We have found a good unused port.
		return port, true
	}

	return 0, false
}
