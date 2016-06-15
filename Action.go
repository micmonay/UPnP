package UPnP

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	scpd "github.com/micmonay/UPnP/SCPD"
)

// Action usued for generate and read action
type Action struct {
	service        *Service
	stateVariables []*scpd.StateVariable
	action         *scpd.Action
	variables      map[string]string
	lastResponse   string
	lastRequest    string
}

// GetName for get action name
func (a *Action) GetName() string {
	return a.action.Name
}

// GetInArguments for get arguments
func (a *Action) GetInArguments() []*Argument {
	arguments := []*Argument{}
	argsSCPD := a.action.GetInArguments()
	for _, argSCPD := range argsSCPD {
		arg := Argument{}
		arg.argSCPD = argSCPD
		arg.scpd = a.service.scpd
		arguments = append(arguments, &arg)
	}
	return arguments
}

// AddVariable add data for argument
func (a *Action) AddVariable(argument string, value string) {
	if a.variables == nil {
		a.variables = make(map[string]string)
	}
	a.variables[argument] = value
}

// Send data and return response or error
func (a *Action) Send() (*Response, error) {
	xmlActionHead := XMLActionHead{}
	xmlActionHead.Xmlns = "http://schemas.xmlsoap.org/soap/envelope/"
	xmlActionHead.EncodingStyle = "http://schemas.xmlsoap.org/soap/encoding/"
	var xmlAction XMLAction
	xmlAction.Xmlns = a.service.ServiceType
	xmlAction.XMLName.Local = "u:" + a.action.Name
	xmlActionHead.Body.Action = xmlAction
	//add argument
	for key, value := range a.variables {
		xmlActionHead.Body.Action.Variables = append(xmlActionHead.Body.Action.Variables, XMLVariable{XMLName: xml.Name{Local: key}, Value: value})
	}
	out, err := xml.MarshalIndent(&xmlActionHead, " ", " ")
	a.lastRequest = string(out)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), out...))
	host := a.service.unupRoot.Location.Host
	if host[len(host)-1] != "/"[0] && a.service.ControlURL[0] != "/"[0] {
		host += "/"
	}
	req, err := http.NewRequest("POST", "http://"+host+a.service.ControlURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("CONTENT-TYPE", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPACTION", "\""+a.service.ServiceType+"#"+a.action.Name+"\"")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	a.lastResponse = string(b)
	if err != nil {
		return nil, err
	}
	var xmlHead XMLRepActionHead
	err = xml.Unmarshal(b, &xmlHead)
	if err != nil {
		return nil, err
	}
	response := Response{}
	response.resultXML = &xmlHead
	scpd, err := a.service.GetSCPD()
	if err != nil {
		return nil, err
	}
	response.scpd = scpd
	response.action = a
	response.codeHTTP = resp.StatusCode
	return &response, nil
}

// GetLastRequest return last requet xml for debug
func (a *Action) GetLastRequest() string {
	return a.lastRequest
}

// GetLastResponse return last response xml
func (a *Action) GetLastResponse() string {
	return a.lastResponse
}

// XMLActionHead Header
type XMLActionHead struct {
	XMLName       xml.Name `xml:"s:Envelope"`
	Xmlns         string   `xml:"xmlns:s,attr"`
	EncodingStyle string   `xml:"s:encodingStyle,attr"`
	Body          XMLBody  `xml:"s:Body"`
}

// XMLBody body
type XMLBody struct {
	Action XMLAction `xml:",any"`
}

// XMLAction action
type XMLAction struct {
	XMLName   xml.Name
	Xmlns     string        `xml:"xmlns:u,attr"`
	Variables []XMLVariable `xml:",any"`
}

// XMLVariable variable
type XMLVariable struct {
	XMLName xml.Name
	Value   string `xml:",innerxml"`
}
