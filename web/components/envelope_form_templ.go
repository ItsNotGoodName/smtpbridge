// Code generated by templ@v0.2.334 DO NOT EDIT.

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

type EnvelopeFormProps struct {
	Flash   templ.Component
	Subject string
	From    string
	To      string
	Body    string
}

func EnvelopeForm(props EnvelopeFormProps) templ.Component {
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
		_, err = templBuffer.WriteString("<form class=\"flex flex-col gap-4\" hx-post=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(routes.EnvelopeCreate().URLString()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" enctype=\"multipart/form-data\"><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">")
		if err != nil {
			return err
		}
		var_2 := `Subject`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></label><input name=\"subject\" type=\"text\" placeholder=\"Subject\" class=\"input input-bordered\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(props.Subject))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">")
		if err != nil {
			return err
		}
		var_3 := `From`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></label><input name=\"from\" type=\"text\" placeholder=\"From\" class=\"input input-bordered\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(props.From))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">")
		if err != nil {
			return err
		}
		var_4 := `To`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></label><input name=\"to\" type=\"text\" placeholder=\"To\" class=\"input input-bordered\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(props.To))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></div><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">")
		if err != nil {
			return err
		}
		var_5 := `Body`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></label><textarea name=\"body\" class=\"textarea textarea-bordered h-24\" placeholder=\"Body\">")
		if err != nil {
			return err
		}
		var var_6 string = props.Body
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</textarea></div><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">")
		if err != nil {
			return err
		}
		var_7 := `Attachments`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></label><input name=\"attachments\" type=\"file\" class=\"file-input file-input-bordered\" multiple></div><button type=\"submit\" class=\"btn btn-primary btn-block\">")
		if err != nil {
			return err
		}
		var_8 := `Create Envelope`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		if props.Flash != nil {
			err = props.Flash.Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</form>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
