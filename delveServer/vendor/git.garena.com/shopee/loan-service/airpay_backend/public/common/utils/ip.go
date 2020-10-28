package utils

import (
	"fmt"
	"math/big"
	"net"
)

func GetIntranetIp() (string, error) {
	ip := ""
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ip, err
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if IsIntranetIP(ipnet.IP) {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	return ip, nil
}

func IsIntranetIP(IP net.IP) bool {
	if IP.IsLoopback() {
		return false
	}

	intranetCIDRs := []string{
		"192.168.0.0/16",
		"172.16.0.0/12",
		"100.64.0.0/10", // This is preserverd for carrier NAT
		"10.0.0.0/8",
	}

	for _, cidr := range intranetCIDRs {
		_, ipNet, _ := net.ParseCIDR(cidr)
		if ipNet.Contains(IP) {
			return true
		}
	}

	return false
}

//整数转成ip,仅支持 ipv4
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

//ip字符串转整数,仅支持 ipv4
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)

	ipI := net.ParseIP(ip)
	if ipI != nil {
		ret.SetBytes(ipI.To4()) //If ip is not an IPv4 address, To4 returns nil.
	}
	return ret.Int64()
}
