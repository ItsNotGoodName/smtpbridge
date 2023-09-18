// routes lists all available HTTP routes.
package routes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

type Route string

func (r Route) String() string {
	return string(r)
}

// URL generates a template safe URL.
func (r Route) URL() templ.SafeURL {
	return templ.URL(string(r))
}

// URLString generates a template safe URL.
func (r Route) URLString() string {
	return string(r.URL())
}

// URLQuery generates a template safe URL with a query.
func (r Route) URLQuery(q string) templ.SafeURL {
	return templ.URL(string(r) + "?" + q)
}

// URLQueryString generates a template safe URL with a query.
func (r Route) URLQueryString(q string) string {
	return string(r.URLQuery(q))
}

// ChildOf checks if the current route is child of given route.
// Only works for root level routes (e.g. "/", "/profile", ...)
func (r Route) ChildOf(route Route) bool {
	if route == "/" {
		return route == r
	}
	return strings.HasPrefix(string(r), string(route))
}

func str(s any) string {
	switch t := s.(type) {
	case string:
		return t
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.Itoa(t)
	default:
		panic(fmt.Sprintf("invalid type: %T", t))
	}
}

// Routes

func Index() Route {
	return "/"
}

func envelope(extras ...string) string {
	s := "/envelopes"
	for _, extra := range extras {
		s += "/" + extra
	}
	return s
}

func Envelope(id any) Route {
	return Route(envelope(str(id)))
}

func EnvelopeHTML(id any) Route {
	return Route(envelope(str(id), "html"))
}

func EnvelopeEndpointSend(id any) Route {
	return Route(envelope(str(id), "endpoint-send"))
}

func EnvelopeList() Route {
	return "/envelopes"
}

func EnvelopeCreate() Route {
	return "/envelopes/create"
}

func AttachmentList() Route {
	return "/attachments"
}

func AttachmentTrim() Route {
	return "/attachments/trim"
}

func AttachmentFile(fileName string) Route {
	return Route("/files/" + fileName)
}

func Login() Route {
	return "/login"
}

func Logout() Route {
	return "/logout"
}

func StorageStatsComponent() Route {
	return "/c/storage-stats"
}

func NullComponent() Route {
	return "/c/null"
}

func RetentionPolicyRun() Route {
	return "/retention-policy/run"
}

type EnvelopeTab string

func (t EnvelopeTab) String() string {
	return string(t)
}

const (
	EnvelopeTabText        EnvelopeTab = ""
	EnvelopeTabHTML        EnvelopeTab = "html"
	EnvelopeTabAttachments EnvelopeTab = "attachments"
)

func EnvelopeTabComponent(id any, tab EnvelopeTab) Route {
	if tab == "" {
		tab = "text"
	}
	return Route("/c/envelopes/" + str(id) + "/tab/" + tab.String())
}

func EndpointList() Route {
	return "/endpoints"
}

func endpoint(extras ...string) string {
	s := "/endpoints"
	for _, extra := range extras {
		s += "/" + extra
	}
	return s
}

func Endpoint(id any) Route {
	return Route(endpoint(str(id)))
}

func EndpointCreate() Route {
	return "/endpoints/create"
}

func EndpointTest(id any) Route {
	return Route(endpoint(str(id), "test"))
}

func TraceList() Route {
	return "/traces"
}

func RuleList() Route {
	return "/rules"
}

func rule(extras ...string) string {
	s := "/rules"
	for _, extra := range extras {
		s += "/" + extra
	}
	return s
}

func Rule(id any) Route {
	return Route(rule(str(id)))
}

func RuleToggle(id any) Route {
	return Route(rule(str(id), "toggle"))
}

func RecentEnvelopeListComponent() Route {
	return "/c/recent-envelope-list"
}

func RuleCreate() Route {
	return "/rules/create"
}

func RuleExpressionCheck() Route {
	return "/rules/expression-check"
}

func EndpointFormConfigComponent() Route {
	return "/c/endpoint-form-config"
}
