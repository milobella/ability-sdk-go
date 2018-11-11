package ability

import (
	"milobella/oratio/pkg/cerebro"
)

type Request struct {
	Nlu 	cerebro.NLU		`json:"nlu,omitempty"`
}
