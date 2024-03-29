// Code generated by templ@v0.2.334 DO NOT EDIT.

package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	c "github.com/ItsNotGoodName/smtpbridge/web/components"
	"github.com/ItsNotGoodName/smtpbridge/web/icons"
	"github.com/ItsNotGoodName/smtpbridge/web/meta"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"strconv"
)

type endpointListViewProps struct {
	Endpoints []models.Endpoint
}

func endpointListView(m meta.Meta, props endpointListViewProps) templ.Component {
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
			_, err = templBuffer.WriteString("<div class=\"border-base-200 breadcrumbs border-b p-4 text-xl font-bold\"><ul><li>")
			if err != nil {
				return err
			}
			var_3 := `Endpoints`
			_, err = templBuffer.WriteString(var_3)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li></ul></div> <div class=\"flex flex-col gap-4 p-4\"><div class=\"join flex items-center justify-end\"><a title=\"Create Endpoint\" class=\"btn btn-sm btn-success\" href=\"")
			if err != nil {
				return err
			}
			var var_4 templ.SafeURL = routes.EndpointCreate().URL()
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_4)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = icons.Add("w-5 h-5").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div></div> <div class=\"mx-auto flex flex-col\">")
			if err != nil {
				return err
			}
			for _, end := range props.Endpoints {
				_, err = templBuffer.WriteString("<div class=\"hover:bg-base-200 border-base-200 flex items-center justify-between gap-2 border-b first:border-t\" id=\"rule-row\"><a class=\"flex-1 truncate py-4 pl-4\" href=\"")
				if err != nil {
					return err
				}
				var var_5 templ.SafeURL = routes.Endpoint(end.ID).URL()
				_, err = templBuffer.WriteString(templ.EscapeString(string(var_5)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var var_6 string = end.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_6))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a><div class=\"flex items-center gap-2 pr-4\">")
				if err != nil {
					return err
				}
				if !end.Internal {
					_, err = templBuffer.WriteString("<div data-loading-states><button title=\"Delete\" class=\"btn btn-error btn-sm\" hx-delete=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(routes.Endpoint(end.ID).URLString()))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\" hx-confirm=\"Are you sure you wish to delete this endpoint?\" hx-target=\"closest #rule-row\" hx-swap=\"outerHTML\" data-loading-disable><span data-loading-class=\"loading loading-spinner loading-xs\">")
					if err != nil {
						return err
					}
					err = icons.Trash("h-4 w-4").Render(ctx, templBuffer)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</span></button></div>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("<div data-loading-states><button class=\"btn btn-sm btn-success\" hx-post=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(routes.EndpointTest(end.ID).URLString()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" data-loading-disable><span data-loading-class=\"loading loading-spinner loading-sm\">")
				if err != nil {
					return err
				}
				var_7 := `Test`
				_, err = templBuffer.WriteString(var_7)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span></button></div></div></div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
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

type endpointCreateProps struct {
	EndpointFormProps c.EndpointFormProps
}

func endpointCreate(m meta.Meta, props endpointCreateProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_8 := templ.GetChildren(ctx)
		if var_8 == nil {
			var_8 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_9 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div class=\"border-base-200 breadcrumbs border-b p-4 text-xl font-bold\"><ul><li><a href=\"")
			if err != nil {
				return err
			}
			var var_10 templ.SafeURL = routes.EndpointList().URL()
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_10)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_11 := `Endpoints`
			_, err = templBuffer.WriteString(var_11)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></li><li>")
			if err != nil {
				return err
			}
			var_12 := `Create`
			_, err = templBuffer.WriteString(var_12)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li></ul></div> <div class=\"mx-auto max-w-lg p-4\">")
			if err != nil {
				return err
			}
			err = c.EndpointForm(props.EndpointFormProps).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = c.LayoutDefault(m).Render(templ.WithChildren(ctx, var_9), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

type endpointViewProps struct {
	Endpoint          models.Endpoint
	EndpointFormProps c.EndpointFormProps
}

func endpointView(m meta.Meta, props endpointViewProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_13 := templ.GetChildren(ctx)
		if var_13 == nil {
			var_13 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_14 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div class=\"border-base-200 breadcrumbs border-b p-4 text-xl font-bold\"><ul><li><a href=\"")
			if err != nil {
				return err
			}
			var var_15 templ.SafeURL = routes.EndpointList().URL()
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_15)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_16 := `Rules`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></li><li>")
			if err != nil {
				return err
			}
			var var_17 string = strconv.FormatInt(props.Endpoint.ID, 10)
			_, err = templBuffer.WriteString(templ.EscapeString(var_17))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li></ul></div> <div class=\"mx-auto flex max-w-lg flex-col gap-4 p-4\">")
			if err != nil {
				return err
			}
			err = c.EndpointForm(props.EndpointFormProps).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("<div data-loading-states><button class=\"btn btn-sm btn-success w-full\" hx-post=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(routes.EndpointTest(props.Endpoint.ID).URLString()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" data-loading-disable><span data-loading-class=\"loading loading-spinner loading-sm\">")
			if err != nil {
				return err
			}
			var_18 := `Test`
			_, err = templBuffer.WriteString(var_18)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span></button></div></div>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = c.LayoutDefault(m).Render(templ.WithChildren(ctx, var_14), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
