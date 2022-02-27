package utils

import "net"

//PrivateIP 本地IP4地址
func PrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip := ipNet.IP
				_, cidr24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
				_, cidr20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
				_, cidr16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
				private := cidr24BitBlock.Contains(ip) || cidr20BitBlock.Contains(ip) || cidr16BitBlock.Contains(ip)
				if private {
					return ip.String()
				}
			}
		}
	}
	return ""
}
