package interpreters

import (
	"github.com/milobella/ability-sdk-go/pkg/model"
)

const instrumentSlot = "instrument_name"
const instrumentNotFound = "I didn't find any instrument in the device matching your request."

type InstrumentInterpreter struct {
	kind   model.InstrumentKind
	action string
}

func FromInstrument(kind model.InstrumentKind, action string) *InstrumentInterpreter {
	return &InstrumentInterpreter{
		kind:   kind,
		action: action,
	}
}

// Interpret get the name of the device's instrument that corresponds to the given parameters.
// If it doesn't find any instrument matching, it will return a pre-computed response with an error message.
// If it finds several instruments matching, it will even return a reprompt answer,
//
//	but only if the instrument can't be found in the NLU and if the previous response was not already a reprompt. (Prevent from infinite loop)
func (i *InstrumentInterpreter) Interpret(req *model.Request) (*string, func(response *model.Response)) {
	instruments := req.Device.CanDo(i.kind, i.action)

	if len(instruments) == 0 {
		// No instrument found, we return an error.
		return nil, func(resp *model.Response) {
			resp.Nlg.Sentence = instrumentNotFound
		}
	} else if len(instruments) > 1 {
		// Several instruments found, we apply a disambiguation algorithm

		if instrument := req.InterpretInstrumentFromNLU(i.kind, i.action); instrument != nil {
			// We found the instrument in NLU, and it is able to do the action
			return &instrument.Name, nil
		}

		if req.IsInSlotFillingAction(i.action, instrumentSlot) {
			// We are in slot filling context and no answer has been found in the NLU,
			//   we stop here to prevent from infinite loop.
			return nil, func(response *model.Response) {
				response.Nlg.Sentence = instrumentNotFound
			}
		}

		// Build a reprompt answer to ask user about instrument
		return nil, func(resp *model.Response) {
			var instrumentsNames []string
			for _, instrument := range instruments {
				instrumentsNames = append(instrumentsNames, instrument.Name)
			}
			resp.Nlg.Sentence = "I found several instruments in the device matching your request : {{ instruments }}. Which one do you want to use ?"
			resp.Nlg.Params = []model.NLGParam{{
				Name:  "instruments",
				Value: instrumentsNames,
				Type:  "enumerated_list",
			}}
			resp.Context.SlotFilling = &model.SlotFilling{
				Action:       i.action,
				MissingSlots: []string{instrumentSlot},
				FilledSlots:  nil,
			}
			resp.AutoReprompt = true
		}
	}

	return &instruments[0].Name, nil
}

type NLUInterpreter struct {
	entity      string
	action      string
	notFoundMsg model.NLG
}

func FromNLU(entity string, action string) *NLUInterpreter {
	defaultNotFoundMsg := model.NLG{
		Sentence: "I didn't find any {{ entity }} in your request. Can you precise ?",
		Params: []model.NLGParam{{
			Name:  "entity",
			Value: entity,
			Type:  "string",
		}},
	}
	return &NLUInterpreter{
		entity:      entity,
		action:      action,
		notFoundMsg: defaultNotFoundMsg,
	}
}

func (i *NLUInterpreter) OverridingNotFoundMsg(msg model.NLG) *NLUInterpreter {
	i.notFoundMsg = msg
	return i
}

func (i *NLUInterpreter) InterpretFirst(req *model.Request) (*string, func(response *model.Response)) {
	values, stopper := i.InterpretAll(req)
	if stopper != nil {
		return nil, stopper
	}
	return &values[0], nil
}

func (i *NLUInterpreter) InterpretAll(req *model.Request) ([]string, func(response *model.Response)) {
	values := req.GetEntitiesByLabel(i.entity)

	if len(values) == 0 {
		// No entity found, we return an error.
		return nil, func(resp *model.Response) {
			resp.Nlg = i.notFoundMsg
			resp.Context.SlotFilling = &model.SlotFilling{
				Action:       i.action,
				MissingSlots: []string{i.entity},
				FilledSlots:  nil,
			}
			resp.AutoReprompt = true
		}
	}
	return values, nil
}
