// events are used by pages/components packages.
package events

import (
	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
)

var (
	EnvelopeCreated    htmx.Event = htmx.NewEventString("envelope-created")
	RetentionPolicyRun htmx.Event = htmx.NewEventString("retention-policy-run")
)

func CSRFToken(csrfToken string) htmx.Event {
	return htmx.NewEvent("csrfToken", csrfToken)
}

func ToastSuccess(toast string) htmx.Event {
	return htmx.NewEvent("toast", toast)
}
