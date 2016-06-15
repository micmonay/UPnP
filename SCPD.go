/*
 dependence for https://github.com/micmonay/UPnP
*/
package SCPD

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

// SCPD service description
type SCPD struct {
	URL            string
	XMLName        xml.Name         `xml:"scpd"`
	Actions        []*Action        `xml:"actionList>action"`
	StateVariables []*StateVariable `xml:"serviceStateTable>stateVariable"`
}

// Action action
type Action struct {
	parentSCPD *SCPD
	XMLName    xml.Name    `xml:"action"`
	Name       string      `xml:"name"`
	Arguments  []*Argument `xml:"argumentList>argument"`
}

// Argument Argument
type Argument struct {
	XMLName              xml.Name `xml:"argument"`
	Name                 string   `xml:"name"`
	Direction            string   `xml:"direction"`
	RelatedStateVariable string   `xml:"relatedStateVariable"`
}

// StateVariable StateVariable
type StateVariable struct {
	parentSCPD         *SCPD
	XMLName            xml.Name           `xml:"stateVariable"`
	SendEvents         string             `xml:"sendEvents,attr"`
	Name               string             `xml:"name"`
	DataType           string             `xml:"dataType"`
	Default            string             `xml:"defaultValue"`
	AllowedValues      []*AllowedValue    `xml:"allowedValueList>allowedValue"`
	AllowedValueRanges *AllowedValueRange `xml:"allowedValueRange"`
}

// AllowedValueRange value range
type AllowedValueRange struct {
	XMLName xml.Name `xml:"allowedValueRange"`
	Minimum string   `xml:"minimum"`
	Maximum string   `xml:"maximum"`
	Step    string   `xml:"step"`
}

// AllowedValue value
type AllowedValue struct {
	XMLName xml.Name `xml:"allowedValue"`
	Value   string   `xml:",innerxml"`
}

// GetDefinitionService get service description
func GetDefinitionService(url string) (*SCPD, error) {
	rep, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(rep.Body)
	r := SCPD{URL: url}
	err = xml.Unmarshal(bytes, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// GetInArguments return IN argument only
func (a *Action) GetInArguments() []*Argument {
	var listArgument []*Argument
	for _, argument := range a.Arguments {
		if strings.Contains(argument.Direction, "in") {
			listArgument = append(listArgument, argument)
		}
	}
	return listArgument
}

// GetArgumentsStateVariable return argument with name
func (a *Action) GetArgumentsStateVariable(name string) string {

	for _, argument := range a.Arguments {
		if name == argument.Name {
			return argument.RelatedStateVariable
		}
	}
	return ""
}

// GetTypeValue get type for variable
func (a *SCPD) GetTypeValue(_stateVariable string) string {
	for _, stateVariable := range a.StateVariables {
		if stateVariable.Name == _stateVariable {
			return stateVariable.DataType
		}
	}
	return ""
}

// GetStateVariable get variable
func (a *SCPD) GetStateVariable(_stateVariable string) *StateVariable {
	for _, stateVariable := range a.StateVariables {
		if stateVariable.Name == _stateVariable {
			stateVariable.parentSCPD = a
			return stateVariable
		}
	}
	return nil
}

// GetAction get action
func (a *SCPD) GetAction(_actionName string) *Action {
	for _, action := range a.Actions {
		if action.Name == _actionName {
			action.parentSCPD = a
			return action
		}
	}
	return nil
}
