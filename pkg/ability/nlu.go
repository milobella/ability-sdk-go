package ability

// NLU nlu
type NLU struct {
	BestIntent string
	Intents    []Intent
	Entities   []Entity
}

// Intent intent
type Intent struct {
	Label string
	Score float32
}

// Entity entity
type Entity struct {
	Label string
	Text  string
}
