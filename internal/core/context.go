package core

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/uptrace/bun"
)

type Context struct {
	context.Context
	*bun.DB
	Actor  Actor
	Bus    *Bus
	File   FileStore
	Config *models.Config
}
