package event

type EventService struct {
	eventRepository Repository
}

func NewEventService(eventRepository Repository) *EventService {
	return &EventService{
		eventRepository: eventRepository,
	}
}

func (es *EventService) Create(ev *Event) error {
	return es.eventRepository.Create(ev)
}
