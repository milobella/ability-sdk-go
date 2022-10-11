package model

import (
	"github.com/milobella/ability-sdk-go/pkg/utils"
)

// Request in the format in which we receive it.
type Request struct {
	Nlu     NLU     `json:"nlu,omitempty"`
	Context Context `json:"context,omitempty"`
	Device  Device  `json:"device,omitempty"`
}

// GetEntitiesByLabel : Extract all entities which have the given label.
func (req *Request) GetEntitiesByLabel(label string) (items []string) {
	for _, ent := range req.Nlu.Entities {
		if ent.Label == label {
			items = append(items, ent.Text)
		}
	}
	return
}

func (req *Request) IsInSlotFillingAction(action string, withAnyOfTheseMissingSlots ...string) bool {
	if req.Context.SlotFilling == nil {
		return false
	}
	isInSlotFillingAction := req.Context.SlotFilling.Action == action
	if !isInSlotFillingAction {
		return false
	}
	oneOfTheSlotIsMissing := false
	for _, slot := range withAnyOfTheseMissingSlots {
		oneOfTheSlotIsMissing = oneOfTheSlotIsMissing || utils.StringSliceContains(req.Context.SlotFilling.MissingSlots, slot)
	}
	return isInSlotFillingAction && (len(withAnyOfTheseMissingSlots) == 0 || oneOfTheSlotIsMissing)
}

// InterpretInstrumentFromNLU Search in the device's instruments the one matching the NLU of the request
//
//	Working examples :
//	  The user has following instruments of a kind : "living room", "parents' bedroom".
//	  User: "I want the living room" / "living room" / "bedroom" > ok
//
// TODO: This algorithm should be improved as it is too naïve.
//
//	Non working example :
//	  The user has following instruments of a kind : "living room", "parents' bedroom".
//	  User: "I want the bedroom" / "bedroom of the parents" > nok
func (req *Request) InterpretInstrumentFromNLU(kind InstrumentKind, action string) *Instrument {
	for _, instrument := range req.Device.Instruments {
		// Verify that the instrument correspond to the predicates
		if instrument.Kind != kind || !utils.StringSliceContains(instrument.Actions, action) {
			continue
		}

		// Search instrument in entities
		entity := req.Nlu.GetFirstEntityOf("instrument")
		if entity != nil && utils.StringFuzzyMatch(instrument.Name, entity.Text) {
			return &instrument
		}

		// Search instrument directly in the text
		if utils.StringFuzzyMatch(instrument.Name, req.Nlu.Text) {
			return &instrument
		}
	}
	return nil
}
