package core

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/uptrace/bun"
)

type Context struct {
	Actor  Actor
	Bus    *Bus
	DB     *bun.DB
	File   FileStore
	Config *models.Config
	ctx    context.Context
}

func (c Context) Context() context.Context {
	return c.ctx
}
