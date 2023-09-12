// events are used by pages/components packages.
package events

import "github.com/ItsNotGoodName/smtpbridge/pkg/htmx"

const (
	EnvelopeCreated    htmx.Event = "envelope-created"
	RetentionPolicyRun htmx.Event = "retention-policy-run"
)
