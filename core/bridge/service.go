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

	bridgesMu sync.Mutex
	bridges   []Bridge
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

	filters := make([]Filter, 0, len(req.Filters))
	for _, filterReq := range req.Filters {
		filter, err := NewFilter(filterReq.From, filterReq.To, filterReq.From, filterReq.ToRegex, filterReq.MatchTemplate)
		if err != nil {
			return err
		}

		filters = append(filters, filter)
	}

	bridge := NewBridge(filters, req.Endpoints)

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
				log.Printf("bridge.BridgeService.send: envelope %d: %s", env.Message.ID, err)
			} else {
				if err := end.Send(ctx, text, atts); err != nil {
					log.Printf("bridge.BridgeService.send: envelope %d: name '%s': type '%s': %s", env.Message.ID, end.Name, end.Type, err)
				}
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

		for bridgeIndex, brid := range bridges {
			// Skip if empty
			if len(brid.Endpoints) == 0 {
				continue
			}

			// Skip if no match
			filterIndex, match := brid.Match(env)
			if !match {
				continue
			}

			log.Printf("bridge.BridgeService.send: envelope %d: bridge '%d': filter '%d': match", env.Message.ID, bridgeIndex, filterIndex)

			for _, endName := range brid.Endpoints {
				// Do not send to a duplicate endpoint
				if _, ok := endNameMemo[endName]; ok {
					log.Printf("bridge.BridgeService.send: envelope %d: bridge '%d': filter '%d': duplicate: skipping", env.Message.ID, bridgeIndex, filterIndex)
					continue
				}
				endNameMemo[endName] = struct{}{}

				end, err := bs.endpointService.GetEndpoint(endName)
				if err != nil {
					log.Printf("bridge.BridgeService.send: envelope %d: bridge '%d': filter '%d': %s", env.Message.ID, bridgeIndex, filterIndex, err)
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
			doneC <- struct{}{}
			return
		case event := <-eventChan:
			env, ok := event.Data.(*envelope.Envelope)
			if !ok {
				log.Println("bridge.BridgeService.Run: could not cast envelope")
				continue
			}

			if err := bs.send(ctx, env); err != nil {
				log.Printf("bridge.BridgeService.Run: envelope %d: %s", env.Message.ID, err)
			}
		}
	}
}
