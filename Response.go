package UPnP

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"log"

	scpd "github.com/micmonay/UPnP/SCPD"
)

// XMLRepActionHead Header
type XMLRepActionHead struct {
	XMLName       xml.Name   `xml:"Envelope"`
	Xmlns         string     `xml:"xmlns:s,attr"`
	EncodingStyle string     `xml:"encodingStyle,attr"`
	Body          XMLRepBody `xml:"Body"`
}

// XMLRepBody Body
type XMLRepBody struct {
	Action XMLRepAction `xml:",any"`
}

// XMLRepAction Action
type XMLRepAction struct {
	XMLName   xml.Name
	Xmlns     string           `xml:"xmlns:u,attr"`
	Variables []XMLRepVariable `xml:",any"`
}

// XMLRepVariable variable
type XMLRepVariable struct {
	XMLName xml.Name

	Value string `xml:",innerxml"`
}

// Error format upnp error
type Error struct {
	XMLName          xml.Name `xml:"UPnPError"`
	ErrorCode        string   `xml:"errorCode"`
	ErrorDescription string   `xml:"errorDescription"`
}

// Response contain a response upnp
type Response struct {
	resultXML *XMLRepActionHead
	scpd      *scpd.SCPD
	action    *Action
	codeHTTP  int
}

// ToString build a string for debug
func (r *Response) ToString() string {
	str := ""
	for _, variable := range r.resultXML.Body.Action.Variables {
		switch r.scpd.GetTypeValue(r.scpd.GetAction(r.action.GetName()).GetArgumentsStateVariable(variable.XMLName.Local)) {
		case "bin.base64":
			base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(variable.Value)))
			_, err := base64.StdEncoding.Decode(base64Text, []byte(variable.Value))
			if err != nil {
				log.Println("Decode base64 error : ", err)
			}
			str += variable.XMLName.Local + " = " + string(base64Text) + "\n"
			break
		default:
			str += variable.XMLName.Local + " = " + variable.Value + "\n"
			break
		}
	}
	return str
}

// Success true if ok
func (r *Response) Success() bool {
	if r.codeHTTP == 200 {
		return true
	}
	return false
}

// GetValueArgument get value of argument
func (r *Response) GetValueArgument(nameArgument string) (string, error) {
	if r == nil || r.resultXML == nil {
		return "", errors.New("Invalide reponse")
	}
	for _, variable := range r.resultXML.Body.Action.Variables {
		if variable.XMLName.Local == nameArgument {
			return variable.Value, nil
		}
	}
	return "", errors.New("not found argument :" + nameArgument)
}

// GetError if response have an error upnp
func (r *Response) GetError() *Error {
	strXML, _ := r.GetValueArgument("detail")
	errorUpnp := Error{}
	err := xml.Unmarshal([]byte(strXML), &errorUpnp)
	if err != nil {
		log.Println(err)
	}
	return &errorUpnp
}
