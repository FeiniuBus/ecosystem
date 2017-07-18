package net

import "net"

// InterfaceIPV4Addrs returns all ip v4 addrs
func InterfaceIPV4Addrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				result = append(result, ip.String())
			}
		}
	}

	return result, nil
}
