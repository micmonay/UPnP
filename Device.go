package UPnP

import (
	"encoding/xml"
)

// Device data
type Device struct {
	upnpRoot *Root
	XMLName  xml.Name `xml:"device"`

	FriendlyName     string     `xml:"friendlyName"`
	Manufacturer     string     `xml:"manufacturer"`
	ManufacturerURL  string     `xml:"manufacturerURL"`
	ModelDescription string     `xml:"modelDescription"`
	ModelName        string     `xml:"modelName"`
	ModelNumber      string     `xml:"modelNumber"`
	ModelURL         string     `xml:"modelURL"`
	SerialNumber     string     `xml:"serialNumber"`
	UDN              string     `xml:"UDN"`
	PresentationURL  string     `xml:"presentationURL"`
	Icons            []*Icon    `xml:"iconList>icon"`
	Devices          []*Device  `xml:"deviceList>device"`
	Services         []*Service `xml:"serviceList>service"`
}

// HasService return true if contain a service
func (d *Device) HasService() bool {
	return len(d.Services) > 0
}

// GetServices return contain service
func (d *Device) GetServices() []*Service {

	for _, service := range d.Services {
		service.unupRoot = d.upnpRoot
	}
	return d.Services
}

// GetServicesByType return service if is a service name
func (d *Device) GetServicesByType(typeName string) []*Service {
	rServices := []*Service{}
	services := d.GetAllService()
	for _, service := range services {
		if service.ServiceType == typeName {
			service.unupRoot = d.upnpRoot
			rServices = append(rServices, service)
		}
	}
	return rServices
}

// GetAllService returnn all services child in child
func (d *Device) GetAllService() []*Service {
	services := d.GetServices()
	if d.HasDevice() {
		for _, device := range d.Devices {
			device.upnpRoot = d.upnpRoot
			childServices := device.GetAllService()
			services = append(services, childServices...)

		}
	}
	return services
}

// HasIcon have icon
func (d *Device) HasIcon() bool {
	return len(d.Icons) > 0
}

// GetIconLink link icon
func (d *Device) GetIconLink() []*Icon {
	return d.Icons
}

// HasDevice have device
func (d *Device) HasDevice() bool {
	return len(d.Devices) > 0
}

// GetDevices return devices child
func (d *Device) GetDevices() []*Device {
	return d.Devices
}
