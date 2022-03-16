package endpoint

import "context"

type EndpointService struct{}

func NewEndpointService() *EndpointService {
	return &EndpointService{}
}

func (EndpointService) Send(ctx context.Context, req []SendRequest) <-chan SendResponse {
	count := len(req)
	resC := make(chan SendResponse, count)
	doneC := make(chan struct{})

	for i := range req {
		go func(msg Envelope, f *Facade) {
			resC <- SendResponse{Envelope: msg, Facade: f, Error: f.Send(ctx, msg)}
			doneC <- struct{}{}
		}(req[i].Envelope, req[i].Facade)
	}

	go func() {
		for i := 0; i < count; i++ {
			<-doneC
		}
		close(resC)
	}()

	return resC
}
