// events are used by pages/components packages.
package events

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
)

const (
	EnvelopeCreated    htmx.Event = "envelope-created"
	RetentionPolicyRun htmx.Event = "retention-policy-run"
)

func CSRFToken(csrfToken string) htmx.Event {
	return htmx.Event(fmt.Sprintf(`{ "csrfToken": "%s" }`, csrfToken))
}
