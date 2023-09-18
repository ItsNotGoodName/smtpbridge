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

type Login struct {
	Username string
	Password string
}

type EndpointCreate struct {
	Name              string
	TextDisable       bool
	AttachmentDisable bool
	TitleTemplate     string
	BodyTemplate      string
	Kind              string
	Config            []endpointConfig
}

type EndpointUpdate struct {
	Name              string
	TextDisable       bool
	AttachmentDisable bool
	TitleTemplate     string
	BodyTemplate      string
	Kind              string
	Config            []endpointConfig
}

type endpointConfig struct {
	Key   string
	Value string
}
