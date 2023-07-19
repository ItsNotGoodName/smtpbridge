package endpoints

import (
	"errors"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
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

func (s Shoutrrr) Send(cc *core.Context, env envelope.Envelope, config Config) error {
	body, err := GetBody(env, config)
	if err != nil {
		return err
	}

	params := types.Params{}
	if title := GetTitle(env, config); title != "" {
		params.SetTitle(title)
	} else if body == "" {
		return nil
	}

	if errs := s.router.Send(body, &params); errs != nil {
		return errors.Join(errs...)
	}

	return nil
}
