// helpers are used by pages/components packages.
package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/schema"
	"github.com/samber/lo"
)

var TimeHourFormat12 string = "12"
var TimeHourFormat24 string = "24"

func BytesHumanize(bytes int64) string {
	return humanize.Bytes(uint64(bytes))
}

func TimeHumanize(date time.Time) string {
	return humanize.Time(date)
}

func JSON(data any) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}

// Query merges current URL query with new values.
func Query(queries url.Values, vals ...any) string {
	tuple := []lo.Tuple2[string, any]{}
	for i := 0; i < len(vals); i += 2 {
		tuple = append(tuple, lo.Tuple2[string, any]{A: (vals[i]).(string), B: vals[i+1]})
	}

	var newQueries []string
	for _, t := range tuple {
		newQueries = append(newQueries, t.A+"="+fmt.Sprint(t.B))
	}

	for key := range queries {
		_, found := queryTupleFind(tuple, key)
		if !found {
			newQueries = append(newQueries, key+"="+queries.Get(key))
		}
	}

	sort.Slice(newQueries, func(i, j int) bool {
		return newQueries[i] < newQueries[j]
	})

	return strings.Join(newQueries, "&")
}

func queryTupleFind(tuple []lo.Tuple2[string, any], key string) (int, bool) {
	for i, t := range tuple {
		if t.A == key {
			return i, true
		}
	}

	return 0, false
}

func Checkbox(r *http.Request, key string) bool {
	q := r.URL.Query()
	if q.Get("-"+key) != "" {
		return q.Get(key) != ""
	}

	return true
}

func Pagination(q url.Values) (pagination.Page, error) {
	page := 1
	if str := q.Get("page"); str != "" {
		var err error
		page, err = strconv.Atoi(q.Get("page"))
		if err != nil {
			return pagination.Page{}, err
		}
	}

	perPage := 1
	if str := q.Get("perPage"); str != "" {
		var err error
		perPage, err = strconv.Atoi(q.Get("perPage"))
		if err != nil {
			return pagination.Page{}, err
		}
	}

	return pagination.NewPage(page, perPage), nil
}

func EndpointSchema() models.EndpointSchema {
	return endpoint.Schema
}

func EndpointsSelections(selected []models.Endpoint, all []models.Endpoint) []bool {
	var endpointsSelections []bool
	for _, end := range all {
		has := lo.ContainsBy(selected, func(e models.Endpoint) bool { return e.ID == end.ID })
		endpointsSelections = append(endpointsSelections, has)
	}
	return endpointsSelections
}

func Tracer(app core.App, r *http.Request) trace.Tracer {
	return app.Tracer(trace.SourceHTTP).Sticky(trace.WithAddress(r.RemoteAddr))
}

var decoder = schema.NewDecoder()

func DecodeForm(w http.ResponseWriter, r *http.Request, form any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(form, r.PostForm); err != nil {
		return err
	}

	return nil
}

func ParseMultipartForm(r *http.Request) error {
	return r.ParseMultipartForm(32 << 20)
}

func Reroute(w http.ResponseWriter, r *http.Request, view http.Handler) {
	htmx.SetRetarget(w, "body")
	view.ServeHTTP(w, r)
}
