package http_utils

import (
	"net"
	"strings"
)

// ip_util.go provide ip relevant utilities

// get all ip binding at ethernet interface and filter by prefix
func GetIpByPrefix(prefix string) ([]string, error) {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, addr := range addrs {
		ip := addr.String()
		if strings.HasPrefix(ip, prefix) {
			ips = append(ips, strings.Split(ip, "/")[0])
		}
	}
	return ips, nil
}

// get all ip by interface name
func GetIpByInterfaceName(netInterfaceName string) ([]string, error) {
	ips := []string{}

	ief, err := net.InterfaceByName(netInterfaceName)
	if err != nil {
		return ips, err
	}

	addrs, err := ief.Addrs()
	if err != nil {
		return ips, err
	}

	for _, addr := range addrs {
		ip := addr.(*net.IPNet).IP.String()
		ips = append(ips, ip)
	}

	return ips, nil
}

// config pattern is ["eth0:10.111", "en5:192.168"], invoke GetLocalIp
func GetIpWithIpConfig(ipConfigs []string) []string {
	ips := []string{}

	for _, ipConfig := range ipConfigs {
		netwInfo := strings.Split(ipConfig, ":")
		ipOnOneInterface, _ := GetLocalIp(netwInfo[0], netwInfo[1])
		ips = append(ips, ipOnOneInterface...)
	}

	return ips
}

// legacy method
func GetLocalIp(netInterface, prefix string) ([]string, error) {
	ips := []string{}

	allIps, err := GetIpByInterfaceName(netInterface)
	if err != nil {
		return ips, err
	}

	for _, ip := range allIps {
		if strings.HasPrefix(ip, prefix) {
			ips = append(ips, ip)
		}
	}

	return ips, nil
}
