//go:build !darwin && !linux && !windows && (!freebsd || !amd64)

package process

import "net/netip"

func findProcessName(network string, ip netip.Addr, srcPort int) (string, error) {
	return "", ErrPlatformNotSupport
}

func resolveSocketByNetlink(network string, ip netip.Addr, srcPort int) (int32, int32, error) {
	return 0, 0, ErrPlatformNotSupport
}
