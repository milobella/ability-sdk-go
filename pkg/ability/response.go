package ability

import (
	"gitlab.milobella.com/milobella/oratio/pkg/anima"
)

type Response struct {
	Nlg 	anima.NLG	`json:"nlg,omitempty"`
	Visu	interface{}	`json:"visu,omitempty"`
}
