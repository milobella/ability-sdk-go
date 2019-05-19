package ability

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
