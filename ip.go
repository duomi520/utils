package utils

import "net"

//PrivateIP4 本地IP4地址
func PrivateIP4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				if ipNet.IP.IsPrivate() {
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}
