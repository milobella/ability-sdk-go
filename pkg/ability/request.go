package ability

// Request in the format in which we receive it.
type Request struct {
	Nlu     NLU     `json:"nlu,omitempty"`
	Context Context `json:"context,omitempty"`
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