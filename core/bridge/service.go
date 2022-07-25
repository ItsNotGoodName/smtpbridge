package bridge

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
)

type BridgeService struct {
	eventPub        *event.Pub
	envelopeService envelope.Service
	endpointService endpoint.Service
	bridgesMu       sync.Mutex
	bridges         []Bridge
}

func NewBridgeService(pub *event.Pub, envelopeService envelope.Service, endpointService endpoint.Service) *BridgeService {
	return &BridgeService{
		eventPub:        pub,
		envelopeService: envelopeService,
		endpointService: endpointService,
	}
}

func (bs *BridgeService) ListBridge() []Bridge {
	bridges := []Bridge{}

	bs.bridgesMu.Lock()
	bridges = append(bridges, bs.bridges...)
	bs.bridgesMu.Unlock()

	return bridges
}

func (bs *BridgeService) CreateBridge(req *CreateBridgeRequest) error {
	for _, end := range req.Endpoints {
		if _, err := bs.endpointService.GetEndpoint(end); err != nil {
			return fmt.Errorf("%w: '%s'", err, end)
		}
	}

	filter, err := NewFilter(req.To, req.From, req.ToRegex, req.FromRegex, req.MatchTemplate)
	if err != nil {
		return err
	}

	bridge := NewBridge(filter, req.Endpoints)

	bs.bridgesMu.Lock()
	bs.bridges = append(bs.bridges, bridge)
	bs.bridgesMu.Unlock()

	return nil
}

func (bs *BridgeService) send(ctx context.Context, env *envelope.Envelope) error {
	atts, err := endpoint.ConvertAttachments(ctx, bs.envelopeService, env)
	if err != nil {
		return err
	}

	// Concurrent endpoint send
	var wg sync.WaitGroup
	send := func(end endpoint.Endpoint) {
		wg.Add(1)
		go func() {
			text, err := end.Text(env)
			if err != nil {
				log.Println("bridge.BridgeService.send:", err)
			} else {
				end.Send(ctx, text, atts)
			}

			wg.Done()
		}()
	}

	bridges := bs.ListBridge()

	// Send to all endpoints there are no bridges
	if len(bridges) == 0 {
		for _, end := range bs.endpointService.ListEndpoint() {
			send(end)
		}
	} else { // Send to bridge's endpoints
		endNameMemo := make(map[string]struct{})

		for _, brid := range bridges {
			if len(brid.Endpoints) == 0 || !brid.Filter.Match(env) {
				continue
			}

			for _, endName := range brid.Endpoints {
				// Do not send to a duplicate endpoint
				if _, ok := endNameMemo[endName]; ok {
					continue
				}
				endNameMemo[endName] = struct{}{}

				end, err := bs.endpointService.GetEndpoint(endName)
				if err != nil {
					log.Println("bridge.BridgeService.send:", err)
					continue
				}

				send(end)
			}
		}
	}

	wg.Wait()

	return nil
}

func (bs *BridgeService) Run(ctx context.Context, doneC chan<- struct{}) {
	log.Println("bridge.BridgeService.Run: started")

	eventChan := make(chan event.Event, 100)
	bs.eventPub.Subscribe(event.TopicEnvelopeCreated, eventChan)
	for {
		select {
		case <-ctx.Done():
			log.Println("bridge.BridgeService.Run: stopped")
			doneC <- struct{}{}
			return
		case event := <-eventChan:
			if err := bs.send(ctx, event.Data.(*envelope.Envelope)); err != nil {
				log.Println("bridge.BridgeService.Run:", err)
			}
		}
	}
}
