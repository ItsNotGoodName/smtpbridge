package forms

type EnvelopeCreate struct {
	Subject string
	From    string
	To      string
	ToSlice []string `schema:"-"`
	Body    string
}

type RuleCreate struct {
	Name       string
	Expression string
	Endpoints  []int64
}

type RuleUpdate struct {
	Name       string
	Expression string
	Endpoints  []int64
}
