package ability

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

func (req *Request) IsInSlotFillingAction(action string) bool {
	return req.Context.SlotFilling.Action == action
}

type NLU struct {
	BestIntent string
	Intents    []Intent
	Entities   []Entity
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

// Device information
type Device struct {
	// Some dynamic information sent with each request
	State map[string]interface{} `json:"state,omitempty"`
}

// Context that will be sent back to us in the next request
type Context struct {
	SlotFilling SlotFilling `json:"slot_filling,omitempty"`
}

// SlotFilling state
type SlotFilling struct {
	Action       string   `json:"action,omitempty"`
	MissingSlots []string `json:"missing_slots,omitempty"`
	FilledSlots  []Slot   `json:"filled_slots,omitempty"`
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

type Response struct {
	Nlg          NLG         `json:"nlg,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	Actions      []Action    `json:"actions,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
	Context      Context     `json:"context,omitempty"`
}
