// Code generated by templ@v0.2.334 DO NOT EDIT.

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/ItsNotGoodName/smtpbridge/web/icons"
	"github.com/ItsNotGoodName/smtpbridge/web/meta"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

func Header(m meta.Meta) templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"navbar bg-base-100 border-b-base-200 sticky top-0 z-10 border-b\"><div class=\"flex-none\"><label for=\"my-drawer-2\" class=\"drawer-button btn btn-square btn-ghost lg:hidden\">")
		if err != nil {
			return err
		}
		err = icons.Menu("h-5 w-5 inline-block").Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label></div><div class=\"flex-1\"><a class=\"btn btn-ghost text-xl normal-case\" href=\"")
		if err != nil {
			return err
		}
		var var_2 templ.SafeURL = routes.Index().URL()
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_2)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_3 := `SMTPBridge`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></div>")
		if err != nil {
			return err
		}
		if !m.Anonymous {
			_, err = templBuffer.WriteString("<div class=\"flex-none\"><div class=\"dropdown dropdown-end\"><label tabindex=\"0\" class=\"btn btn-square btn-ghost\">")
			if err != nil {
				return err
			}
			err = icons.More("h-5 w-5").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><ul class=\"menu dropdown-content bg-base-100 rounded-box z-50 w-52 p-2 shadow-lg\"><li><a href=\"#\" hx-delete=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(routes.Logout().URLString()))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = icons.LogoutBoxR("h-4 w-4").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_4 := `Logout`
			_, err = templBuffer.WriteString(var_4)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></li></ul></div></div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
