package model

import "github.com/milobella/ability-sdk-go/pkg/utils"

type NLU struct {
	BestIntent string
	Intents    []Intent
	Entities   []Entity
	Text       string
}

func (nlu *NLU) GetFirstEntityOf(label string) *Entity {
	for _, ent := range nlu.Entities {
		if ent.Label == label {
			return &ent
		}
	}
	return nil
}

type Intent struct {
	Label string
	Score float32
}

type Entity struct {
	Label string
	Text  string
}

type NLG struct {
	Sentence string     `json:"sentence,omitempty"`
	Params   []NLGParam `json:"params,omitempty"`
}

type NLGParam struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Type  string      `json:"type,omitempty"`
}

type InstrumentKind string

const (
	InstrumentKindChromeCast InstrumentKind = "chromecast"
	InstrumentStateUnknown                  = "UNKNOWN"
)

type InstrumentState struct {
	Status string
}

type Instrument struct {
	Kind    InstrumentKind
	Actions []string
	Name    string
	State   InstrumentState
}

func (i *Instrument) IsActive() bool {
	return len(i.State.Status) > 0 && i.State.Status != InstrumentStateUnknown
}

// Device information
type Device struct {
	// Some dynamic information sent with each request
	State       map[string]interface{} `json:"state,omitempty"`
	Instruments []Instrument           `json:"instruments,omitempty"`
}

func (d *Device) CanDo(kind InstrumentKind, action string) (result []Instrument) {
	for _, instrument := range d.Instruments {
		if instrument.Kind == kind && utils.StringSliceContains(instrument.Actions, action) {
			result = append(result, instrument)
		}
	}
	return
}

// Context that will be sent back to us in the next request
type Context struct {
	SlotFilling *SlotFilling `json:"slot_filling,omitempty"`
}

// SlotFilling state
type SlotFilling struct {
	Action       string   `json:"action,omitempty"`
	MissingSlots []string `json:"missing_slots,omitempty"`
	FilledSlots  []Slot   `json:"filled_slots,omitempty"`
}

func (sf *SlotFilling) GetSlot() *Slot {
	return nil //TODO
}

// Slot used for slot filling mechanism
type Slot struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type ActionParameter struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Action struct {
	Identifier string            `json:"identifier,omitempty"`
	Params     []ActionParameter `json:"params,omitempty"`
}
