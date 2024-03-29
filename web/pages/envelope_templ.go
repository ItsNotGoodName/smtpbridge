// Code generated by templ@v0.2.334 DO NOT EDIT.

package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"net/url"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	c "github.com/ItsNotGoodName/smtpbridge/web/components"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/ItsNotGoodName/smtpbridge/web/icons"
	"github.com/ItsNotGoodName/smtpbridge/web/meta"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"strings"
)

type envelopeViewProps struct {
	Envelope  models.Envelope
	Endpoints []models.Endpoint
	Tab       routes.EnvelopeTab
}

func envelopeView(m meta.Meta, props envelopeViewProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_2 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div class=\"border-base-200 breadcrumbs border-b p-4 text-xl font-bold\"><ul><li><a href=\"")
			if err != nil {
				return err
			}
			var var_3 templ.SafeURL = routes.EnvelopeList().URL()
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_3)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_4 := `Envelopes`
			_, err = templBuffer.WriteString(var_4)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></li><li>")
			if err != nil {
				return err
			}
			var var_5 string = strconv.FormatInt(props.Envelope.Message.ID, 10)
			_, err = templBuffer.WriteString(templ.EscapeString(var_5))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li></ul></div> <div class=\"overflow-x-auto\"><table class=\"table\"><tbody><tr><th class=\"w-0 whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_6 := `From`
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><td>")
			if err != nil {
				return err
			}
			var var_7 string = props.Envelope.Message.From
			_, err = templBuffer.WriteString(templ.EscapeString(var_7))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td></tr><tr><th class=\"w-0 whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_8 := `Subject`
			_, err = templBuffer.WriteString(var_8)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><td>")
			if err != nil {
				return err
			}
			var var_9 string = props.Envelope.Message.Subject
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td></tr><tr><th class=\"w-0 whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_10 := `To`
			_, err = templBuffer.WriteString(var_10)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><td class=\"flex items-center gap-2\">")
			if err != nil {
				return err
			}
			for _, to := range props.Envelope.Message.To {
				_, err = templBuffer.WriteString("<span class=\"badge\">")
				if err != nil {
					return err
				}
				var var_11 string = to
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</td></tr><tr><th class=\"w-0 whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_12 := `Date`
			_, err = templBuffer.WriteString(var_12)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><td>")
			if err != nil {
				return err
			}
			err = c.FormatDate(m, props.Envelope.Message.Date.Time()).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td></tr><tr><th class=\"w-0 whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_13 := `Created At`
			_, err = templBuffer.WriteString(var_13)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><td>")
			if err != nil {
				return err
			}
			err = c.FormatDate(m, props.Envelope.Message.CreatedAt.Time()).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td></tr></tbody></table></div> <div class=\"flex flex-col gap-4 p-4\" data-loading-states><form class=\"join\" hx-post=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(routes.EnvelopeEndpointSend(props.Envelope.Message.ID).URLString()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-target=\"next div\" hx-swap=\"outerHTML\"><select name=\"endpoint\" id=\"endpoint\" class=\"select select-bordered select-sm join-item w-full\"><option disabled selected>")
			if err != nil {
				return err
			}
			var_14 := `Select Endpoint`
			_, err = templBuffer.WriteString(var_14)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option>")
			if err != nil {
				return err
			}
			for _, end := range props.Endpoints {
				_, err = templBuffer.WriteString("<option value=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(strconv.FormatInt(end.ID, 10)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var var_15 string = end.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_15))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</select><button type=\"submit\" class=\"btn-sm btn btn-primary join-item\" data-loading-disable><span data-loading-class=\"loading loading-spinner loading-sm\">")
			if err != nil {
				return err
			}
			var_16 := `Send`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span></button></form><div class=\"hidden\"></div></div> ")
			if err != nil {
				return err
			}
			err = c.EnvelopeTab(c.EnvelopeTabProps{Envelope: props.Envelope, Tab: props.Tab}).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = c.LayoutDefault(m).Render(templ.WithChildren(ctx, var_2), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

type envelopeListViewProps struct {
	EnvelopeRequestRequest models.DTOEnvelopeListRequest
	EnvelopeRequestResult  models.DTOEnvelopeListResult
	Query                  url.Values
}

func envelopeListView(m meta.Meta, props envelopeListViewProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_17 := templ.GetChildren(ctx)
		if var_17 == nil {
			var_17 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_18 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div class=\"border-base-200 breadcrumbs border-b p-4 text-xl font-bold\"><ul><li>")
			if err != nil {
				return err
			}
			var_19 := `Envelopes`
			_, err = templBuffer.WriteString(var_19)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li></ul></div> <div class=\"flex flex-col gap-4 p-4\"><div class=\"flex flex-col-reverse justify-between gap-4 sm:flex-row\"><form class=\"flex gap-2\" action=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(routes.EnvelopeList().URLString()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			for k := range props.Query {
				if !strings.HasPrefix(k, "search") {
					_, err = templBuffer.WriteString("<input type=\"hidden\" name=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(k))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(props.Query.Get(k)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
				}
			}
			_, err = templBuffer.WriteString("<div class=\"join\"><div class=\"dropdown join-item\"><label tabindex=\"0\" class=\"btn btn-sm join-item\">")
			if err != nil {
				return err
			}
			err = icons.Filter("w-5 h-5").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><ul tabindex=\"0\" class=\"dropdown-content bg-base-100 rounded-box z-[1] w-52 p-2 shadow-lg\"><li><div class=\"form-control\"><label class=\"label cursor-pointer\"><span class=\"label-text\">")
			if err != nil {
				return err
			}
			var_20 := `Subject`
			_, err = templBuffer.WriteString(var_20)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span><input type=\"hidden\" value=\"0\" name=\"-search-subject\"><input name=\"search-subject\" type=\"checkbox\" class=\"checkbox checkbox-sm\"")
			if err != nil {
				return err
			}
			if props.EnvelopeRequestRequest.SearchSubject {
				_, err = templBuffer.WriteString(" checked")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("></label></div></li><li><div class=\"form-control\"><label class=\"label cursor-pointer\"><span class=\"label-text\">")
			if err != nil {
				return err
			}
			var_21 := `Text`
			_, err = templBuffer.WriteString(var_21)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span><input type=\"hidden\" value=\"0\" name=\"-search-text\"><input name=\"search-text\" type=\"checkbox\" class=\"checkbox checkbox-sm\"")
			if err != nil {
				return err
			}
			if props.EnvelopeRequestRequest.SearchText {
				_, err = templBuffer.WriteString(" checked")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("></label></div></li></ul></div><input name=\"search\" type=\"text\" placeholder=\"Search\" class=\"input input-sm input-bordered join-item w-full max-w-xs\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(props.EnvelopeRequestRequest.Search))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"><button title=\"Search Envelopes\" type=\"submit\" class=\"btn btn-sm btn-primary join-item\">")
			if err != nil {
				return err
			}
			err = icons.Search("w-5 h-5").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></div></form><div class=\"join flex items-center justify-end\" data-loading-states><a title=\"Create Envelope\" class=\"btn btn-sm join-item btn-success\" href=\"")
			if err != nil {
				return err
			}
			var var_22 templ.SafeURL = routes.EnvelopeCreate().URL()
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_22)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" data-loading-states>")
			if err != nil {
				return err
			}
			err = icons.Add("w-5 h-5").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a><button title=\"Delete All Envelopes\" class=\"btn btn-sm btn-error join-item\" hx-delete=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(routes.EnvelopeList().URLString()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-confirm=\"Are you sure you wish to delete all envelopes?\" data-loading-disable><span data-loading-class=\"loading loading-spinner loading-sm\">")
			if err != nil {
				return err
			}
			err = icons.Trash("w-5 h-5").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span></button></div></div>")
			if err != nil {
				return err
			}
			err = c.PaginateHeader(c.PaginateHeaderProps{
				Route:      routes.EnvelopeList(),
				Query:      props.Query,
				PageResult: props.EnvelopeRequestResult.PageResult,
				Ascending:  props.EnvelopeRequestRequest.Ascending,
			}).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div> <div class=\"overflow-x-auto\"><table class=\"table-pin-cols table\"><thead><tr><th></th><td>")
			if err != nil {
				return err
			}
			var_23 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
				templBuffer, templIsBuffer := w.(*bytes.Buffer)
				if !templIsBuffer {
					templBuffer = templ.GetBuffer()
					defer templ.ReleaseBuffer(templBuffer)
				}
				var_24 := `Created At`
				_, err = templBuffer.WriteString(var_24)
				if err != nil {
					return err
				}
				if !templIsBuffer {
					_, err = io.Copy(w, templBuffer)
				}
				return err
			})
			err = c.TableOrderTH(c.TableOrderTHProps{
				Query:     props.Query,
				Ascending: props.EnvelopeRequestRequest.Ascending,
				Order:     string(props.EnvelopeRequestRequest.Order),
				Field:     models.DTOEnvelopeFieldCreatedAt,
			}).Render(templ.WithChildren(ctx, var_23), templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td>")
			if err != nil {
				return err
			}
			var_25 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
				templBuffer, templIsBuffer := w.(*bytes.Buffer)
				if !templIsBuffer {
					templBuffer = templ.GetBuffer()
					defer templ.ReleaseBuffer(templBuffer)
				}
				var_26 := `From`
				_, err = templBuffer.WriteString(var_26)
				if err != nil {
					return err
				}
				if !templIsBuffer {
					_, err = io.Copy(w, templBuffer)
				}
				return err
			})
			err = c.TableOrderTH(c.TableOrderTHProps{
				Query:     props.Query,
				Ascending: props.EnvelopeRequestRequest.Ascending,
				Order:     string(props.EnvelopeRequestRequest.Order),
				Field:     models.DTOEnvelopeFieldFrom,
			}).Render(templ.WithChildren(ctx, var_25), templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td>")
			if err != nil {
				return err
			}
			var_27 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
				templBuffer, templIsBuffer := w.(*bytes.Buffer)
				if !templIsBuffer {
					templBuffer = templ.GetBuffer()
					defer templ.ReleaseBuffer(templBuffer)
				}
				var_28 := `Subject`
				_, err = templBuffer.WriteString(var_28)
				if err != nil {
					return err
				}
				if !templIsBuffer {
					_, err = io.Copy(w, templBuffer)
				}
				return err
			})
			err = c.TableOrderTH(c.TableOrderTHProps{
				Query:     props.Query,
				Ascending: props.EnvelopeRequestRequest.Ascending,
				Order:     string(props.EnvelopeRequestRequest.Order),
				Field:     models.DTOEnvelopeFieldSubject,
			}).Render(templ.WithChildren(ctx, var_27), templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><th></th></tr></thead><tbody>")
			if err != nil {
				return err
			}
			for _, env := range props.EnvelopeRequestResult.Envelopes {
				_, err = templBuffer.WriteString("<tr><th class=\"z-10 w-0 whitespace-nowrap\">")
				if err != nil {
					return err
				}
				var var_29 string = strconv.FormatInt(env.Message.ID, 10)
				_, err = templBuffer.WriteString(templ.EscapeString(var_29))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</th><td class=\"w-0 whitespace-nowrap\"><div class=\"tooltip tooltip-right\" data-tip=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(helpers.TimeHumanize(env.Message.CreatedAt.Time())))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"><a href=\"")
				if err != nil {
					return err
				}
				var var_30 templ.SafeURL = routes.Envelope(env.Message.ID).URL()
				_, err = templBuffer.WriteString(templ.EscapeString(string(var_30)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				err = c.FormatDate(m, time.Time(env.Message.CreatedAt)).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a></div></td><td class=\"w-0 whitespace-nowrap\"><a href=\"")
				if err != nil {
					return err
				}
				var var_31 templ.SafeURL = routes.Envelope(env.Message.ID).URL()
				_, err = templBuffer.WriteString(templ.EscapeString(string(var_31)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var var_32 string = env.Message.From
				_, err = templBuffer.WriteString(templ.EscapeString(var_32))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a></td><td class=\"whitespace-nowrap\"><a href=\"")
				if err != nil {
					return err
				}
				var var_33 templ.SafeURL = routes.Envelope(env.Message.ID).URL()
				_, err = templBuffer.WriteString(templ.EscapeString(string(var_33)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var var_34 string = env.Message.Subject
				_, err = templBuffer.WriteString(templ.EscapeString(var_34))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a></td><th class=\"w-0 whitespace-nowrap\"><div class=\"flex flex-row justify-end gap-2\">")
				if err != nil {
					return err
				}
				if len(env.Attachments) != 0 {
					_, err = templBuffer.WriteString("<a title=\"Attachments\" href=\"")
					if err != nil {
						return err
					}
					var var_35 templ.SafeURL = routes.Envelope(env.Message.ID).URLQuery("tab=" + routes.EnvelopeTabAttachments.String())
					_, err = templBuffer.WriteString(templ.EscapeString(string(var_35)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" class=\"tooltip tooltip-left flex items-center fill-current\" data-tip=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(len(env.Attachments)) + " Attachment(s)"))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					err = icons.Attachment("h-4 w-4").Render(ctx, templBuffer)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</a>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("<button title=\"Delete\" class=\"btn btn-error btn-xs join-item\" hx-delete=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(routes.Envelope(env.Message.ID).URLString()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" hx-target=\"closest tr\" hx-confirm=\"Are you sure you wish to delete this envelope?\">")
				if err != nil {
					return err
				}
				err = icons.Trash("h-4 w-4").Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</button></div></th></tr>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</tbody></table></div> ")
			if err != nil {
				return err
			}
			if len(props.EnvelopeRequestResult.Envelopes) != 0 {
				err = c.PaginateFooter(c.PaginateFooterProps{
					Route:      routes.EnvelopeList(),
					Query:      props.Query,
					PageResult: props.EnvelopeRequestResult.PageResult,
				}).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = c.LayoutDefault(m).Render(templ.WithChildren(ctx, var_18), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

type envelopeCreateViewProps struct {
}

func envelopeCreateView(m meta.Meta, props envelopeCreateViewProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_36 := templ.GetChildren(ctx)
		if var_36 == nil {
			var_36 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_37 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div><div class=\"border-base-200 breadcrumbs border-b p-4 text-xl font-bold\"><ul><li><a href=\"")
			if err != nil {
				return err
			}
			var var_38 templ.SafeURL = routes.EnvelopeList().URL()
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_38)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_39 := `Envelopes`
			_, err = templBuffer.WriteString(var_39)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></li><li>")
			if err != nil {
				return err
			}
			var_40 := `Create`
			_, err = templBuffer.WriteString(var_40)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li></ul></div><div class=\"mx-auto flex max-w-lg flex-col gap-4 p-4\">")
			if err != nil {
				return err
			}
			err = c.EnvelopeForm(c.EnvelopeFormProps{}).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></div>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = c.LayoutDefault(m).Render(templ.WithChildren(ctx, var_37), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
