// pages is used by http package.
package pages

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	c "github.com/ItsNotGoodName/smtpbridge/web/components"
	"github.com/ItsNotGoodName/smtpbridge/web/events"
	"github.com/ItsNotGoodName/smtpbridge/web/forms"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/ItsNotGoodName/smtpbridge/web/meta"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"github.com/ItsNotGoodName/smtpbridge/web/session"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/kballard/go-shellquote"
)

// utils

const ParamID = "{id}"

func withID(ct Controller, fn func(w http.ResponseWriter, r *http.Request, id int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}
		fn(w, r, id)
	}
}

func getCode(err error) int {
	if errors.Is(err, models.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

// Controller

type Controller interface {
	// Meta returns the meta data for the request.
	Meta(r *http.Request) meta.Meta
	// Page renders a html page.
	Page(w http.ResponseWriter, r *http.Request, body templ.Component)
	// Component renders a html component.
	Component(w http.ResponseWriter, r *http.Request, body templ.Component)
	// Error renders error.
	Error(w http.ResponseWriter, r *http.Request, err error, code int)
}

// Pages / Components

func NullComponent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func IndexView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storage, err := app.StorageGet(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		envelopeList, err := app.EnvelopeList(ctx, pagination.NewPage(1, 5), models.DTOEnvelopeListRequest{})
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		policy := app.RetentionPolicyGet(ctx)

		ct.Page(w, r, indexView(ct.Meta(r), indexViewProps{
			StorageStatsProps: c.StorageStatsProps{
				Storage: storage,
			},
			Envelopes:       envelopeList.Envelopes,
			RetentionPolicy: policy,
		}))
	}
}

func RecentEnvelopeListComponent(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		envelopeList, err := app.EnvelopeList(ctx, pagination.NewPage(1, 5), models.DTOEnvelopeListRequest{})
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		ct.Component(w, r, c.RecentEnvelopeList(ct.Meta(r), c.RecentEnvelopeListProps{
			Envelopes: envelopeList.Envelopes,
		}))
	}
}

func StorageStatsComponent(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storage, err := app.StorageGet(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		ct.Component(w, r, c.StorageStats(c.StorageStatsProps{
			Storage: storage,
		}))
	}
}

func EnvelopeListView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()

		pagination, err := helpers.Pagination(query)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		req := models.DTOEnvelopeListRequest{
			Ascending:     query.Get("ascending") != "",
			Search:        query.Get("search"),
			SearchSubject: helpers.Checkbox(r, "search-subject"),
			SearchText:    helpers.Checkbox(r, "search-text"),
			Order:         models.NewDTOEnvelopeField(query.Get("order")),
		}

		res, err := app.EnvelopeList(ctx, pagination, req)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		if res.PageResult.Overflow() {
			// Page does not exist
			url := routes.EnvelopeList().URLQueryString(helpers.Query(query, "page", res.PageResult.TotalPages))
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		ct.Page(w, r, envelopeListView(ct.Meta(r), envelopeListViewProps{
			Query:                  query,
			EnvelopeRequestRequest: req,
			EnvelopeRequestResult:  res,
		}))
	}
}

func EnvelopeListDrop(ct Controller, app core.App, view http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := app.EnvelopeDrop(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		htmx.SetRetarget(w, "body")
		view.ServeHTTP(w, r)
	}
}

func AttachmentListView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()

		page, err := helpers.Pagination(query)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		req := models.DTOAttachmentListRequest{
			Ascending: query.Get("ascending") != "",
		}

		res, err := app.AttachmentList(ctx, page, req)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		if res.PageResult.Overflow() {
			// Page does not exist
			url := routes.AttachmentList().URLQueryString(helpers.Query(query, "page", res.PageResult.TotalPages))
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		ct.Page(w, r, attachmentListView(ct.Meta(r), attachmentListViewProps{
			Query:             query,
			AttachmentRequest: req,
			AttachmentResult:  res,
		}))
	}
}

func EnvelopeCreateView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct.Page(w, r, envelopeCreateView(ct.Meta(r), envelopeCreateViewProps{}))
	}
}

func EnvelopeCreate(ct Controller, app core.App) http.HandlerFunc {
	handleErr := func(w http.ResponseWriter, r *http.Request, err error, form forms.EnvelopeCreate) {
		ct.Component(w, r, c.EnvelopeForm(c.EnvelopeFormProps{
			Subject: form.Subject,
			To:      form.To,
			From:    form.From,
			Body:    form.Body,
			Flash:   c.Flash(c.FlashTypeError, c.FlashMessage(err.Error())),
		}))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.ParseMultipartForm(r)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		var form forms.EnvelopeCreate
		if err := helpers.DecodeForm(w, r, &form); err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}
		form.ToSlice = strings.Split(form.To, ",")

		msg := models.DTOMessageCreate{
			Subject: form.Subject,
			From:    form.From,
			To:      form.ToSlice,
			Text:    form.Body,
			Date:    time.Now(),
		}

		var datts []models.DTOAttachmentCreate
		for _, fh := range r.MultipartForm.File["attachments"] {
			a, err := fh.Open()
			if err != nil {
				handleErr(w, r, err, form)
				return
			}
			defer a.Close()

			datts = append(datts, models.DTOAttachmentCreate{
				Name: fh.Filename,
				Data: a,
			})
		}

		ctx := r.Context()

		id, err := app.EnvelopeCreate(ctx, msg, datts)
		if err != nil {
			handleErr(w, r, err, form)
			return
		}
		helpers.Tracer(app, r).Trace(ctx, trace.ActionEnvelopeCreated, trace.WithEnvelope(id))

		events.EnvelopeCreated.SetTrigger(w)
		htmx.SetLocation(w, routes.Envelope(id).URLString())
	}
}

func LoginView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct.Page(w, r, loginView(ct.Meta(r)))
	}
}

func Login(ct Controller, app core.App, ss sessions.Store) http.HandlerFunc {
	handleErr := func(w http.ResponseWriter, r *http.Request, err error, form forms.Login) {
		ct.Component(w, r, c.LoginForm(c.LoginFormProps{
			Flash:    c.Flash(c.FlashTypeError, c.FlashMessage(err.Error())),
			Username: form.Username,
			Password: form.Password,
		}))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.Login
		if err := helpers.DecodeForm(w, r, &form); err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		user, err := app.AuthHTTPLogin(ctx, form.Username, form.Password)
		if err != nil {
			handleErr(w, r, err, form)
			return
		}

		err = session.AuthLogin(w, r, ss, user.ID)
		if err != nil {
			handleErr(w, r, err, form)
			return
		}

		htmx.SetRedirect(w, routes.Index().URLString())
	}
}

func Logout(ct Controller, app core.App, ss sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := session.AuthLogout(w, r, ss)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		htmx.SetRedirect(w, routes.Login().URLString())
	}
}

func EnvelopeView(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		env, err := app.EnvelopeGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		ends, err := app.EndpointList(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		tab := r.URL.Query().Get("tab")

		ct.Page(w, r, envelopeView(ct.Meta(r), envelopeViewProps{
			Envelope:  env,
			Endpoints: ends,
			Tab:       routes.EnvelopeTab(tab),
		}))
	})
}

func EnvelopeDelete(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		err := app.EnvelopeDelete(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func EnvelopeHTMLView(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		html, err := app.MessageHTMLGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})
}

func EnvelopeTabComponent(ct Controller, app core.App, tab routes.EnvelopeTab) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()
		query := r.URL.Query()

		env, err := app.EnvelopeGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		htmx.SetReplaceURL(w, routes.Envelope(id).URLQueryString(helpers.Query(query, "tab", tab.String())))
		ct.Component(w, r, c.EnvelopeTab(c.EnvelopeTabProps{
			Envelope: env,
			Tab:      tab,
		}))
	})
}

func TraceListView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()

		page, err := helpers.Pagination(query)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		listRequest := models.DTOTraceListRequest{
			Ascending: query.Get("ascending") != "",
		}

		listResult, err := app.TraceList(ctx, page, listRequest)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		// Page requested does not exist
		if listResult.PageResult.Overflow() {
			url := routes.TraceList().URLQueryString(helpers.Query(query, "page", listResult.PageResult.TotalPages))
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		ct.Page(w, r, traceListView(ct.Meta(r), traceListViewProps{
			TraceListRequest: listRequest,
			TraceListResult:  listResult,
			Query:            query,
		}))
	}
}

func TraceListDrop(ct Controller, app core.App, view http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := app.TraceDrop(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		helpers.Reroute(w, r, view)
	}
}
func RuleListView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rules, err := app.RuleList(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		ct.Page(w, r, ruleListView(ct.Meta(r), ruleListViewProps{
			Rules: rules,
		}))
	}
}

func AttachmentTrim(ct Controller, app core.App, view http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := helpers.Tracer(app, r)

		err := app.AttachmentOrphanDelete(ctx, tracer)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		helpers.Reroute(w, r, view)
	}
}

func RetentionPolicyRun(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := helpers.Tracer(app, r)

		err := app.RetentionPolicyRun(ctx, tracer)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		events.RetentionPolicyRun.SetTrigger(w)
		w.WriteHeader(http.StatusNoContent)
	}
}

func EnvelopeEndpointSend(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		endpointID, err := strconv.ParseInt(r.FormValue("endpoint"), 10, 64)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		err = app.EnvelopeSend(ctx, id, endpointID)
		if err != nil {
			ct.Component(w, r, c.Flash(c.FlashTypeError, c.FlashMessage(err.Error())))
			return
		}

		ct.Component(w, r, c.Flash(c.FlashTypeSuccess, c.FlashMessage("Sent")))
	})
}

func Files(ct Controller, app core.App, fs fs.FS) http.HandlerFunc {
	idFromFilename := func(r *http.Request) (int64, error) {
		path := r.URL.Path
		i := strings.Index(path, ".")
		if i == -1 {
			return 0, fmt.Errorf("invalid filename")
		}

		idStr := path[:i]
		return strconv.ParseInt(idStr, 10, 64)
	}

	fsHandler := http.FileServer(http.FS(fs))

	return http.StripPrefix("/files/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := idFromFilename(r)
		if err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		att, err := app.AttachmentGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		download, _ := strconv.ParseBool(r.URL.Query().Get("download"))
		if download {
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, shellquote.Join(att.Name)))
		}

		fsHandler.ServeHTTP(w, r)
	})).ServeHTTP
}
