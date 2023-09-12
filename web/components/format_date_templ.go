// Code generated by templ@v0.2.334 DO NOT EDIT.

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "time"
import "github.com/ItsNotGoodName/smtpbridge/web/meta"

func FormatDate(m meta.Meta, t time.Time) templ.Component {
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
		_, err = templBuffer.WriteString("<sl-format-date month=\"numeric\" day=\"numeric\" year=\"numeric\" hour=\"numeric\" minute=\"numeric\" hour-format=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(m.TimeHourFormat))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" date=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(t.Format(time.RFC3339)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></sl-format-date>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}