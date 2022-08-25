package model

type Response struct {
	Nlg          NLG         `json:"nlg,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	Actions      []Action    `json:"actions,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
	Context      Context     `json:"context,omitempty"`
}
