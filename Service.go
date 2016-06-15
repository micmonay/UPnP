package UPnP

import (
	"encoding/xml"
	"errors"
	"log"

	scpd "github.com/micmonay/UPnP/SCPD"
)

// Service for manipulation service
type Service struct {
	unupRoot    *Root
	scpd        *scpd.SCPD
	XMLName     xml.Name `xml:"service"`
	ServiceType string   `xml:"serviceType"`
	ServiceID   string   `xml:"serviceId"`
	ControlURL  string   `xml:"controlURL"`
	EventSubURL string   `xml:"eventSubURL"`
	SCPDURL     string   `xml:"SCPDURL"`
}

// GetSCPD return service description
func (s *Service) GetSCPD() (*scpd.SCPD, error) {
	var err error
	if s.unupRoot == nil {
		return nil, errors.New("Error root file description")
	}
	host := s.unupRoot.Location.Host
	if host[len(host)-1] != "/"[0] && s.SCPDURL[0] != "/"[0] {
		host += "/"
	}
	if s.scpd == nil {
		s.scpd, err = scpd.GetDefinitionService("http://" + host + s.SCPDURL)
	}
	if err != nil {
		return nil, errors.New("Error load service description file at " + "http://" + host + s.SCPDURL + " find in root file :" + s.unupRoot.Location.String())
	}
	return s.scpd, nil
}
func (s *Service) newAction(actionSCPD *scpd.Action) *Action {
	scpd, err := s.GetSCPD()
	if err != nil {
		log.Print(err)
		return nil
	}
	action := Action{}
	action.service = s
	action.action = actionSCPD
	action.stateVariables = scpd.StateVariables
	return &action
}

// GetActions return all actions for the service
func (s *Service) GetActions() []*Action {
	scpd, err := s.GetSCPD()
	if err != nil {
		log.Print(err)
		return nil
	}
	actions := []*Action{}
	for _, actionSCPD := range scpd.Actions {
		actions = append(actions, s.newAction(actionSCPD))
	}
	return actions
}

// GetAction return action selected by name
func (s *Service) GetAction(name string) *Action {
	scpd, err := s.GetSCPD()
	if err != nil {
		log.Print(err)
		return nil
	}
	for _, actionSCPD := range scpd.Actions {
		if actionSCPD.Name == name {
			return s.newAction(actionSCPD)
		}
	}
	return nil
}
