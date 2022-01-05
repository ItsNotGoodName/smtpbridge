package core

import "context"

type BackgroundPort interface {
	Run(ctx context.Context, done chan<- struct{})
}
