package background

import "context"

type Background interface {
	Run(ctx context.Context, done chan<- struct{})
}
