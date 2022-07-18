package event

type (
	Event struct {
		Topic Topic
		Data  interface{}
	}

	Topic string
)

const (
	TopicEnvelopeCreated Topic = "envelope.created" // *envelope.Envelope
)
