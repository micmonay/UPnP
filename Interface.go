package UPnP

import (
	"errors"
	"net"
	"strings"
)

// GetInterfaceByName get interface by name
func GetInterfaceByName(name string) (*net.Interface, error) {
	Interfaces, _ := net.Interfaces()
	for _, Interface := range Interfaces {
		if Interface.Name == name {
			return &Interface, nil
		}
	}

	return nil, errors.New("Not found interface : " + name)
}

// GetInterfaceByIP get interface by ip 192.168.1.1 (search by contain not exactly)
func GetInterfaceByIP(ip string) (*net.Interface, error) {
	Interfaces, _ := net.Interfaces()
	for _, Interface := range Interfaces {
		addrs, err := Interface.Addrs()
		if err != nil {
			continue
		}
		for _, add := range addrs {
			if strings.Contains(add.String(), ip) {
				return &Interface, nil
			}
		}
	}
	return nil, errors.New("Not found interface with ip : " + ip)
}

// GetIPAdress for get ip adress
func GetIPAdress(_interface *net.Interface) string {
	var ip net.IP
	ipv4 := ""
	addrs, _ := _interface.Addrs()
	for _, interfaceIP := range addrs {
		switch v := interfaceIP.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip.To4() != nil {
			ipv4 = ip.To4().String()
		}
	}
	return ipv4
}
