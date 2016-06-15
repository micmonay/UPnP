package UPnP

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	return strings.TrimSpace(str)
}
func selectInterface() (*net.Interface, error) {
	addrs, _ := net.Interfaces()
	for n, addr := range addrs {
		ip, err := addr.Addrs()
		if err != nil {
			continue
		}
		fmt.Println(n, " : ", addr.Name, " {", addr.HardwareAddr, "} ", ip)
	}
	fmt.Print("Please select interface number : ")
	interfaceNum, err := strconv.ParseUint(getInput(), 10, 64)
	if err != nil {
		return nil, err
	}
	return &addrs[interfaceNum], nil
}

//ExampleNewUPNP example for get external ipv4 address from gateway
func ExampleNewUPNP() {
	up := NewUPNP(SERVICE_GATEWAY_IPV4_V2)
	_interface, err := selectInterface()
	if err != nil {
		panic(err)
	}
	devices := up.GetAllCompatibleDevice(_interface, 1)
	if len(devices) == 0 {
		return
	}
	services := devices[0].GetServicesByType(SERVICE_GATEWAY_IPV4_V2)
	if len(services) == 0 {
		return
	}
	service := services[0]
	response, err := service.GetAction("GetExternalIPAddress").Send()
	if err != nil {
		panic(err)
	}
	fmt.Println(response.ToString())
	fmt.Print("Press enter")
	getInput()
	// Output: external ip
}
