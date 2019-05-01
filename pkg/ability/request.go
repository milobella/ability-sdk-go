package ability

import (
    "gitlab.milobella.com/milobella/oratio/pkg/cerebro"
)

type Request struct {
    Nlu             cerebro.NLU     `json:"nlu,omitempty"`,
    AutoReprompt    bool            `json:"auto_reprompt,omitempty"`
}
