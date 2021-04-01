package ability

type Response struct {
	Nlg          NLG         `json:"nlg,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
	Context      Context     `json:"context,omitempty"`
}
