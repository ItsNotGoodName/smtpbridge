package endpoint

import (
	"context"
	"fmt"
)

type Console struct{}

func (c *Console) Send(ctx context.Context, text string, atts []Attachment) error {
	output := text
	for _, a := range atts {
		output += "\n" + a.Name
	}
	fmt.Println(output)
	return nil
}
