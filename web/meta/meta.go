package meta

import (
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

type Meta struct {
	Route          routes.Route
	Anonymous      bool
	TimeHourFormat string
}
