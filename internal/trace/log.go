package trace

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/rs/zerolog/log"
)

type LogStore struct{}

// Save implements Store.
func (LogStore) Save(ctx context.Context, trace models.Trace) error {
	switch trace.Level {
	case LevelInfo:
		log.Info().Msg(fmt.Sprintf("%+v", trace))
	case LevelError:
		log.Error().Msg(fmt.Sprintf("%+v", trace))
	default:
		return fmt.Errorf("unknown trace level: %s", trace.Level)
	}

	return nil
}
