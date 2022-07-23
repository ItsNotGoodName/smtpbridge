package endpoint

import (
	"context"
	"fmt"
)

type Console struct{}

func (c *Console) Send(ctx context.Context, text string, atts []Attachment) error {
	fmt.Println(text)
	return nil
}
