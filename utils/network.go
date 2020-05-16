package utils

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
)

var (
	errNoAvailableInterface = errors.New("no available interface")
	errNoAvailableAddress   = errors.New("no available address")
)
// RoutedInterface returns a network interface that can route IP
// traffic and satisfies flags.
//
// The provided network must be "ip", "ip4" or "ip6".
func RoutedInterface(network string, flags net.Flags) (*net.Interface, error) {
	switch network {
	case "ip", "ip4", "ip6":
	default:
		return nil, errNoAvailableInterface
	}
	ift, err := net.Interfaces()
	if err != nil {
		return nil, errNoAvailableInterface
	}
	for _, ifi := range ift {
		if ifi.Flags&flags != flags {
			continue
		}
		if _, ok := hasRoutableIP(network, &ifi); !ok {
			continue
		}
		return &ifi, nil
	}
	return nil, errNoAvailableInterface
}


func hasRoutableIP(network string, ifi *net.Interface) (net.IP, bool) {
	ifat, err := ifi.Addrs()
	if err != nil {
		return nil, false
	}
	for _, ifa := range ifat {
		log.Debug("Intf name ",ifa.String())
		switch ifa := ifa.(type) {
		case *net.IPAddr:
			if ip, ok := routableIP(network, ifa.IP); ok {
				return ip, true
			}
		case *net.IPNet:
			if ip, ok := routableIP(network, ifa.IP); ok {
				return ip, true
			}
		}
	}
	return nil, false
}

func routableIP(network string, ip net.IP) (net.IP, bool) {
	if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() && !ip.IsGlobalUnicast() {
		return nil, false
	}
	log.Debug("IP is routable ",ip.IsLoopback())
	switch network {
	case "ip4":
		if ip := ip.To4(); ip != nil {
			return ip, true
		}
	case "ip6":
		if ip.IsLoopback() { // addressing scope of the loopback address depends on each implementation
			return nil, false
		}
		if ip := ip.To16(); ip != nil && ip.To4() == nil {
			return ip, true
		}
	default:
		if ip.IsLoopback() { // addressing scope of the loopback address depends on each implementation
			return nil, false
		}
		if ip := ip.To4(); ip != nil {
			return ip, true
		}
		if ip := ip.To16(); ip != nil {
			return ip, true
		}
	}
	return nil, false
}
