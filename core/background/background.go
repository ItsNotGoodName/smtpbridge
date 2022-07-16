package background

import "context"

type Background interface {
	Run(ctx context.Context, doneC chan<- struct{})
}

func Run(ctx context.Context, backgrounds []Background) {
	done := make(chan struct{}, len(backgrounds))
	running := 0

	// Start backgrounds
	for _, background := range backgrounds {
		go background.Run(ctx, done)
		running++
	}

	// Wait for context
	<-ctx.Done()

	// Wait for backgrounds
	for i := 0; i < running; i++ {
		<-done
	}
}
