package UPnP

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Icon icon info
type Icon struct {
	XMLName  xml.Name `xml:"icon"`
	Mimetype string   `xml:"mimetype"`
	Width    int      `xml:"width"`
	Height   int      `xml:"height"`
	Depth    int      `xml:"depth"`
	URL      string   `xml:"url"`
}

// Root root description file
type Root struct {
	XMLName  xml.Name `xml:"root"`
	Location *url.URL
	Device   *Device `xml:"device"`
}

// NewXMLUPNPFile get description file
func NewXMLUPNPFile(_url string) (*Root, error) {
	rep, err := http.Get(_url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(rep.Body)
	r := Root{}
	err = xml.Unmarshal(bytes, &r)
	if err != nil {
		return nil, err
	}
	r.Location, _ = url.Parse(_url)
	return &r, nil
}

// GetDevice return child device
func (r *Root) GetDevice() *Device {
	device := r.Device
	device.upnpRoot = r
	return device
}

// GetLocation return location url
func (r *Root) GetLocation() *url.URL {
	return r.Location
}
