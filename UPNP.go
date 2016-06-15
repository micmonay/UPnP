/*
	Easy for dialog with a UPnP device in golang
*/
package UPnP

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	ssdpDiscover   = `"ssdp:discover"`
	ssdpUDP4Addr   = "239.255.255.250:1900"
	ssdpSearchPort = 1900
	methodSearch   = "M-SEARCH"
	methodNotify   = "NOTIFY"
)

//Exemple of service research
const (
	ALL_DEVICE              = "upnp:rootdevice"
	SERVICE_GATEWAY_STATE   = "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1"
	SERVICE_GATEWAY_IPV4_V1 = "urn:schemas-upnp-org:service:WANIPConnection:1"
	SERVICE_GATEWAY_IPV4_V2 = "urn:schemas-upnp-org:service:WANIPConnection:2"
	SERVICE_GATEWAY_IPV6    = "urn:schemas-upnp-org:service:WANIPv6FirewallControl:1"
)

// UPNP is a struct for dialog on UPNP
type UPNP struct {
	_serivceType string
}
type dialResponse struct {
	IP   string
	UUID string
}
type dialog struct {
	stop   bool
	socket *net.UDPConn
}

// NewUPNP for create new UPNP request for type device
func NewUPNP(serviceType string) *UPNP {
	return &UPNP{_serivceType: serviceType}
}

// NewUPNPAllService for create new UPNP request for all type device
func NewUPNPAllService() *UPNP {
	return &UPNP{_serivceType: ALL_DEVICE}
}
func (d *dialog) reader(cresponse chan *http.Response, crequest chan *http.Request) {

	message := make([]byte, 4096)
	for !d.stop {
		n, _, err := d.socket.ReadFromUDP(message)
		if err != nil {
			break
		}
		response, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(message[:n])), nil)
		if err != nil {
			request, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(message[:n])))
			if err != nil {
				log.Printf("ReadResponse: ERROR reading - %s", err)
			} else {
				crequest <- request
			}
		} else {
			cresponse <- response
		}
	}
	d.socket.Close()
}
func (d *dialog) write(request *http.Request, addr *net.UDPAddr) {
	var buf bytes.Buffer
	if _, err := fmt.Fprintf(&buf, "%s %s HTTP/1.1\r\n", request.Method, request.URL.RequestURI()); err != nil {
		log.Println("Socket write error ! ", err)
	}
	request.Header.Write(&buf)
	if _, err := buf.Write([]byte{'\r', '\n'}); err != nil {
		log.Println("Socket write error ! ", err)
	}
	n, err := d.socket.WriteToUDP(buf.Bytes(), addr)
	if err != nil {
		log.Println("Socket write error ! ", err)
	}
	if n <= 0 {
		log.Println("Socket write error ! no sending data")
	}
}

// GetAllCompatibleDevice return All devices
func (u *UPNP) GetAllCompatibleDevice(_interface *net.Interface, timeout int) []*Device {
	devices := []*Device{}
	listsLink := u.GetLinkDeviceOnly(_interface, timeout)
	for _, link := range listsLink {
		upnpRoot, err := NewXMLUPNPFile(link)
		if err != nil {
			log.Println("Error load : '", link, "' error :", err)
			continue
		}
		devices = append(devices, upnpRoot.GetDevice())
	}
	return devices
}

// GetLinkDeviceOnly For get all link description file
func (u *UPNP) GetLinkDeviceOnly(_interface *net.Interface, timeout int) []string {
	d := dialog{stop: false}
	ip := GetIPAdress(_interface)
	udpAddr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		log.Fatalln("1 : ResolveUDPAddr failed: ", err)
	}
	udpAddr2, err := net.ResolveUDPAddr("udp4", ip+":0")
	if err != nil {
		log.Fatalln("2 : ResolveUDPAddr failed: ", ip, " -- ", err)
	}
	d.socket, err = net.ListenUDP("udp4", udpAddr2)
	if err != nil {
		log.Fatalln(err)
	}
	if timeout <= 0 {
		timeout = 1
	}
	var cresponse = make(chan *http.Response, 10)
	var crequest = make(chan *http.Request, 10)
	go d.reader(cresponse, crequest)
	req := http.Request{
		Method: methodSearch,
		Host:   ssdpUDP4Addr,
		URL:    &url.URL{Opaque: "*"},
		Header: http.Header{
			// Putting headers in here avoids them being title-cased.
			// (The UPnP discovery protocol uses case-sensitive headers)
			"HOST": []string{ssdpUDP4Addr},
			"MX":   []string{strconv.FormatInt(int64(timeout), 10)},
			"MAN":  []string{ssdpDiscover},
			"ST":   []string{u._serivceType},
		},
	}
	d.write(&req, udpAddr)
	ctimeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		ctimeout <- true
	}()
	var UpnpLink []string
	for !d.stop {
		select {
		case response := <-cresponse:
			UpnpLink = append(UpnpLink, response.Header.Get("Location"))
		case <-ctimeout:
			d.socket.Close()
			d.stop = true
		}
	}
	return UpnpLink
}
