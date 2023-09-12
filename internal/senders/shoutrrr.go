package senders

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type Shoutrrr struct {
	router *router.ServiceRouter
}

func NewShoutrrr(router *router.ServiceRouter) Shoutrrr {
	return Shoutrrr{
		router: router,
	}
}

func (s Shoutrrr) Send(ctx context.Context, env models.Envelope, tr Transformer) error {
	params := types.Params{}
	body, err := tr.Body(ctx, env)
	if err != nil {
		return err
	}

	title, err := tr.Title(ctx, env)
	if err != nil {
		return err
	}
	if title != "" {
		params.SetTitle(title)
	} else if body == "" {
		return nil
	}

	if errs := s.router.Send(body, &params); errs != nil {
		return errors.Join(errs...)
	}

	return nil
}
