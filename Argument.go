package UPnP

import scpd "github.com/micmonay/UPnP/SCPD"

// Argument arguement
type Argument struct {
	argSCPD *scpd.Argument
	scpd    *scpd.SCPD
	sv      *scpd.StateVariable
}

// GetName return name
func (a *Argument) GetName() string {
	return a.argSCPD.Name
}

// GetDirection return in or out
func (a *Argument) GetDirection() string {
	return a.argSCPD.Direction
}
func (a *Argument) getStateVariable() *scpd.StateVariable {
	if a.sv != nil {
		return a.sv
	}
	for _, sv := range a.scpd.StateVariables {
		if sv.Name == a.argSCPD.RelatedStateVariable {
			a.sv = sv
		}
	}
	return a.sv
}

// GetType return type for argument
func (a *Argument) GetType() string {
	sv := a.getStateVariable()
	if sv == nil {
		return ""
	}
	return sv.DataType
}

// GetAllowedValues return value possible
func (a *Argument) GetAllowedValues() []string {
	sv := a.getStateVariable()
	strs := []string{}
	if sv == nil {
		return strs
	}
	for _, str := range sv.AllowedValues {
		strs = append(strs, str.Value)
	}
	return strs
}

// GetDefault return default string or ""
func (a *Argument) GetDefault() string {
	sv := a.getStateVariable()
	if sv == nil {
		return ""
	}
	return sv.Default
}

// GetMaximum return maximum or ""
func (a *Argument) GetMaximum() string {
	sv := a.getStateVariable()
	if sv == nil {
		return ""
	}
	return sv.AllowedValueRanges.Maximum
}

// GetMinimum return minimum or ""
func (a *Argument) GetMinimum() string {
	sv := a.getStateVariable()
	if sv == nil {
		return ""
	}
	return sv.AllowedValueRanges.Minimum
}

// GetStep return step or ""
func (a *Argument) GetStep() string {
	sv := a.getStateVariable()
	if sv == nil {
		return ""
	}
	return sv.AllowedValueRanges.Step
}
