// Code generated by templ@v0.2.334 DO NOT EDIT.

package icons

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func _(class any) templ.Component {
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
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/attachment-line

func Attachment(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_2 := templ.GetChildren(ctx)
		if var_2 == nil {
			var_2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_3 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_3...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_3).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M14 13.5V8C14 5.79086 12.2091 4 10 4C7.79086 4 6 5.79086 6 8V13.5C6 17.0899 8.91015 20 12.5 20C16.0899 20 19 17.0899 19 13.5V4H21V13.5C21 18.1944 17.1944 22 12.5 22C7.80558 22 4 18.1944 4 13.5V8C4 4.68629 6.68629 2 10 2C13.3137 2 16 4.68629 16 8V13.5C16 15.433 14.433 17 12.5 17C10.567 17 9 15.433 9 13.5V8H11V13.5C11 14.3284 11.6716 15 12.5 15C13.3284 15 14 14.3284 14 13.5Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/delete-bin-line

func Trash(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_4 := templ.GetChildren(ctx)
		if var_4 == nil {
			var_4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_5 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_5...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_5).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M17 6H22V8H20V21C20 21.5523 19.5523 22 19 22H5C4.44772 22 4 21.5523 4 21V8H2V6H7V3C7 2.44772 7.44772 2 8 2H16C16.5523 2 17 2.44772 17 3V6ZM18 8H6V20H18V8ZM9 11H11V17H9V11ZM13 11H15V17H13V11ZM9 4V6H15V4H9Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/mail-line

func Mail(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_6 := templ.GetChildren(ctx)
		if var_6 == nil {
			var_6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_7 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_7...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_7).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M3 3H21C21.5523 3 22 3.44772 22 4V20C22 20.5523 21.5523 21 21 21H3C2.44772 21 2 20.5523 2 20V4C2 3.44772 2.44772 3 3 3ZM20 7.23792L12.0718 14.338L4 7.21594V19H20V7.23792ZM4.51146 5L12.0619 11.662L19.501 5H4.51146Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/add-line

func Add(class any) templ.Component {
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
		var var_9 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_9...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_9).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M11 11V5H13V11H19V13H13V19H11V13H5V11H11Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/more-line

func More(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_10 := templ.GetChildren(ctx)
		if var_10 == nil {
			var_10 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_11 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_11...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_11).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M4.5 10.5C3.675 10.5 3 11.175 3 12C3 12.825 3.675 13.5 4.5 13.5C5.325 13.5 6 12.825 6 12C6 11.175 5.325 10.5 4.5 10.5ZM19.5 10.5C18.675 10.5 18 11.175 18 12C18 12.825 18.675 13.5 19.5 13.5C20.325 13.5 21 12.825 21 12C21 11.175 20.325 10.5 19.5 10.5ZM12 10.5C11.175 10.5 10.5 11.175 10.5 12C10.5 12.825 11.175 13.5 12 13.5C12.825 13.5 13.5 12.825 13.5 12C13.5 11.175 12.825 10.5 12 10.5Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/send-plane-line

func SendPlane(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_12 := templ.GetChildren(ctx)
		if var_12 == nil {
			var_12 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_13 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_13...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_13).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M21.7264 2.95706L16.2732 22.0433C16.1222 22.5718 15.7976 22.5958 15.5561 22.1127L10.9998 13.0002L1.92266 9.36931C1.41298 9.16544 1.41929 8.86034 1.9567 8.6812L21.0429 2.31913C21.5714 2.14297 21.8745 2.43878 21.7264 2.95706ZM19.0351 5.0966L6.81197 9.17097L12.4486 11.4256L15.4893 17.507L19.0351 5.0966Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/home-4-line

func Home(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_14 := templ.GetChildren(ctx)
		if var_14 == nil {
			var_14 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_15 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_15...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_15).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M19 21.0001H5C4.44772 21.0001 4 20.5524 4 20.0001V11.0001L1 11.0001L11.3273 1.61162C11.7087 1.26488 12.2913 1.26488 12.6727 1.61162L23 11.0001L20 11.0001V20.0001C20 20.5524 19.5523 21.0001 19 21.0001ZM13 19.0001H18V9.15757L12 3.70302L6 9.15757V19.0001H11V13.0001H13V19.0001Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/ruler-line

func Ruler(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_16 := templ.GetChildren(ctx)
		if var_16 == nil {
			var_16 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_17 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_17...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_17).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M6.34323 14.7278L3.5148 17.5563L7.05033 21.0918L20.4854 7.65678L16.9498 4.12124L14.8285 6.24257L16.2427 7.65678L14.8285 9.07099L13.4143 7.65678L11.293 9.7781L13.4143 11.8994L12.0001 13.3136L9.87876 11.1923L7.75744 13.3136L9.17165 14.7278L7.75744 16.1421L6.34323 14.7278ZM17.6569 1.99992L22.6067 6.94967C22.9972 7.3402 22.9972 7.97336 22.6067 8.36389L7.75744 23.2131C7.36692 23.6037 6.73375 23.6037 6.34323 23.2131L1.39348 18.2634C1.00295 17.8729 1.00295 17.2397 1.39348 16.8492L16.2427 1.99992C16.6332 1.6094 17.2664 1.6094 17.6569 1.99992Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/sketching

func Sketching(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_18 := templ.GetChildren(ctx)
		if var_18 == nil {
			var_18 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_19 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_19...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_19).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M10.0357 7.69802C8.38492 9.55932 6.5134 12.2442 4.89465 15.4817C4.64766 15.9757 4.04698 16.1759 3.55301 15.9289C3.05903 15.6819 2.8588 15.0812 3.10579 14.5873C4.79739 11.2041 6.76494 8.37171 8.53943 6.37095C9.4251 5.37234 10.2797 4.56162 11.0449 3.99131C11.4272 3.7063 11.8049 3.46806 12.1677 3.29756C12.5193 3.13234 12.921 3 13.3336 3C13.5496 3 13.7872 3.0535 14.007 3.19476C14.2233 3.33371 14.3629 3.51925 14.4495 3.69083C14.6066 4.00215 14.624 4.33473 14.6201 4.55938C14.6118 5.03651 14.4847 5.6328 14.3216 6.23975C13.9874 7.48318 13.3994 9.13104 12.8149 10.7577L12.7329 10.9858C12.1671 12.5598 11.6101 14.1093 11.248 15.3466C11.1505 15.68 11.0706 15.9792 11.0094 16.2414C11.7035 15.6835 12.5581 14.8454 13.466 13.9534L13.4956 13.9243C14.3772 13.0581 15.3098 12.1418 16.0967 11.5127C16.4872 11.2006 16.9082 10.904 17.3138 10.7322C17.6544 10.5878 18.4343 10.3532 19.0407 10.9596C19.4251 11.344 19.5318 11.8438 19.5594 12.2164C19.5883 12.6064 19.5429 13.0267 19.4725 13.4261C19.3315 14.2258 19.0483 15.159 18.7894 16.0009L18.7478 16.136C18.5165 16.8874 18.3102 17.5577 18.1926 18.0965C18.4529 17.8352 18.7734 17.4216 19.1475 16.811C19.436 16.34 20.0517 16.1921 20.5226 16.4806C20.9935 16.7691 21.1414 17.3848 20.8529 17.8557C20.3099 18.7422 19.748 19.4622 19.1519 19.9092C18.5283 20.377 17.7121 20.6407 16.8863 20.2278C16.2779 19.9235 16.1398 19.3173 16.1091 18.9819C16.0759 18.6192 16.1284 18.2233 16.1979 17.8667C16.3288 17.1944 16.5829 16.3698 16.823 15.5907L16.8777 15.4129C17.1447 14.5451 17.3873 13.734 17.5028 13.0789C17.5117 13.0284 17.5196 12.9802 17.5266 12.9341C17.4697 12.977 17.4094 13.0239 17.3455 13.0749C16.6477 13.6328 15.785 14.4788 14.8677 15.38L14.8381 15.4091C13.9566 16.2752 13.024 17.1915 12.2371 17.8206C11.8466 18.1328 11.4255 18.4293 11.02 18.6012C10.6794 18.7455 9.89947 18.9801 9.29311 18.3738C8.9843 18.065 8.9052 17.6753 8.87972 17.4382C8.8515 17.1755 8.86901 16.8971 8.90269 16.6351C8.9706 16.1069 9.12934 15.4656 9.32855 14.7849C9.70829 13.4872 10.2842 11.8852 10.8411 10.3362L10.9327 10.0814C11.5263 8.42931 12.082 6.8674 12.3901 5.72074C12.4172 5.61968 12.4418 5.52435 12.4638 5.43468C12.3924 5.48361 12.3178 5.53695 12.2401 5.59489C11.6173 6.05907 10.8627 6.76559 10.0357 7.69802Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/telegram-line

func Telegram(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_20 := templ.GetChildren(ctx)
		if var_20 == nil {
			var_20 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_21 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_21...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_21).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12.001 20C16.4193 20 20.001 16.4183 20.001 12C20.001 7.58172 16.4193 4 12.001 4C7.5827 4 4.00098 7.58172 4.00098 12C4.00098 16.4183 7.5827 20 12.001 20ZM12.001 22C6.47813 22 2.00098 17.5228 2.00098 12C2.00098 6.47715 6.47813 2 12.001 2C17.5238 2 22.001 6.47715 22.001 12C22.001 17.5228 17.5238 22 12.001 22ZM8.89113 13.1708L6.39308 12.3911C5.85317 12.2255 5.84996 11.854 6.51397 11.5871L16.247 7.82754C16.8115 7.59644 17.1325 7.88842 16.9495 8.61641L15.2917 16.4377C15.1763 16.9941 14.8407 17.1266 14.3759 16.8698L11.825 14.9823L10.6361 16.1297C10.5139 16.2476 10.4145 16.3489 10.2273 16.3743C10.04 16.3996 9.88564 16.3442 9.77337 16.0337L8.90378 13.1628L8.89113 13.1708Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/terminal-line

func Terminal(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_22 := templ.GetChildren(ctx)
		if var_22 == nil {
			var_22 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_23 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_23...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_23).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M10.9999 11.9998L3.92886 19.0708L2.51465 17.6566L8.1715 11.9998L2.51465 6.34292L3.92886 4.92871L10.9999 11.9998ZM10.9999 18.9998H20.9999V20.9998H10.9999V18.9998Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/logout-box-r-line

func LogoutBoxR(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_24 := templ.GetChildren(ctx)
		if var_24 == nil {
			var_24 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_25 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_25...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_25).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M5 22C4.44772 22 4 21.5523 4 21V3C4 2.44772 4.44772 2 5 2H19C19.5523 2 20 2.44772 20 3V6H18V4H6V20H18V18H20V21C20 21.5523 19.5523 22 19 22H5ZM18 16V13H11V11H18V8L23 12L18 16Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/recycle-line

func Recycle(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_26 := templ.GetChildren(ctx)
		if var_26 == nil {
			var_26 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_27 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_27...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_27).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M19.5605 12.0975L21.092 14.7501C22.0585 16.4241 21.4849 18.5647 19.8109 19.5312C19.2788 19.8384 18.6753 20.0001 18.0609 20.0001L15.9985 19.9995L15.9987 22.5001L10.9987 19.0001L15.9987 15.5001L15.9985 17.9995L18.0609 18.0001C18.3242 18.0001 18.5829 17.9308 18.8109 17.7991C19.5283 17.3849 19.7742 16.4675 19.3599 15.7501L17.8285 13.0975L19.5605 12.0975ZM7.30257 9.13407L7.83366 15.2142L5.66848 13.9635L4.63751 15.7501C4.50586 15.9781 4.43655 16.2368 4.43655 16.5001C4.43655 17.3285 5.10812 18.0001 5.93655 18.0001L8.99948 18.0002V20.0001H5.93655C4.00355 20.0001 2.43655 18.4331 2.43655 16.5001C2.43655 15.8857 2.59827 15.2822 2.90546 14.7501L3.93648 12.9635L1.77148 11.7142L7.30257 9.13407ZM13.7487 2.96901C14.2808 3.2762 14.7226 3.71803 15.0298 4.2501L16.0605 6.03551L18.226 4.78599L17.6949 10.8661L12.1638 8.28599L14.3285 7.03551L13.2978 5.2501C13.1661 5.02207 12.9768 4.83271 12.7487 4.70106C12.0313 4.28684 11.1139 4.53266 10.6997 5.2501L9.16818 7.90261L7.43615 6.90266L8.96764 4.2501C9.93414 2.57607 12.0747 2.00251 13.7487 2.96901Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/search-line

func Search(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_28 := templ.GetChildren(ctx)
		if var_28 == nil {
			var_28 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_29 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_29...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_29).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M18.031 16.6168L22.3137 20.8995L20.8995 22.3137L16.6168 18.031C15.0769 19.263 13.124 20 11 20C6.032 20 2 15.968 2 11C2 6.032 6.032 2 11 2C15.968 2 20 6.032 20 11C20 13.124 19.263 15.0769 18.031 16.6168ZM16.0247 15.8748C17.2475 14.6146 18 12.8956 18 11C18 7.1325 14.8675 4 11 4C7.1325 4 4 7.1325 4 11C4 14.8675 7.1325 18 11 18C12.8956 18 14.6146 17.2475 15.8748 16.0247L16.0247 15.8748Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/eye-line

func Eye(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_30 := templ.GetChildren(ctx)
		if var_30 == nil {
			var_30 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_31 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_31...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_31).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12.0003 3C17.3924 3 21.8784 6.87976 22.8189 12C21.8784 17.1202 17.3924 21 12.0003 21C6.60812 21 2.12215 17.1202 1.18164 12C2.12215 6.87976 6.60812 3 12.0003 3ZM12.0003 19C16.2359 19 19.8603 16.052 20.7777 12C19.8603 7.94803 16.2359 5 12.0003 5C7.7646 5 4.14022 7.94803 3.22278 12C4.14022 16.052 7.7646 19 12.0003 19ZM12.0003 16.5C9.51498 16.5 7.50026 14.4853 7.50026 12C7.50026 9.51472 9.51498 7.5 12.0003 7.5C14.4855 7.5 16.5003 9.51472 16.5003 12C16.5003 14.4853 14.4855 16.5 12.0003 16.5ZM12.0003 14.5C13.381 14.5 14.5003 13.3807 14.5003 12C14.5003 10.6193 13.381 9.5 12.0003 9.5C10.6196 9.5 9.50026 10.6193 9.50026 12C9.50026 13.3807 10.6196 14.5 12.0003 14.5Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/arrow-down-s-line

func ArrowUpS(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_32 := templ.GetChildren(ctx)
		if var_32 == nil {
			var_32 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_33 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_33...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_33).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M11.9997 10.8284L7.04996 15.7782L5.63574 14.364L11.9997 8L18.3637 14.364L16.9495 15.7782L11.9997 10.8284Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/arrow-down-s-line

func ArrowDownS(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_34 := templ.GetChildren(ctx)
		if var_34 == nil {
			var_34 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_35 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_35...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_35).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M11.9997 13.1714L16.9495 8.22168L18.3637 9.63589L11.9997 15.9999L5.63574 9.63589L7.04996 8.22168L11.9997 13.1714Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/download-line

func Download(class any) templ.Component {
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
		var var_37 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_37...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_37).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M3 19H21V21H3V19ZM13 13.1716L19.0711 7.1005L20.4853 8.51472L12 17L3.51472 8.51472L4.92893 7.1005L11 13.1716V2H13V13.1716Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/checkbox-circle-line

func CheckboxCircle(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_38 := templ.GetChildren(ctx)
		if var_38 == nil {
			var_38 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_39 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_39...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_39).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12C22 17.5228 17.5228 22 12 22ZM12 20C16.4183 20 20 16.4183 20 12C20 7.58172 16.4183 4 12 4C7.58172 4 4 7.58172 4 12C4 16.4183 7.58172 20 12 20ZM11.0026 16L6.75999 11.7574L8.17421 10.3431L11.0026 13.1716L16.6595 7.51472L18.0737 8.92893L11.0026 16Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/close-circle-line

func CloseCircle(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_40 := templ.GetChildren(ctx)
		if var_40 == nil {
			var_40 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_41 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_41...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_41).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" viewBox=\"0 0 24 24\"><path d=\"M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12C22 17.5228 17.5228 22 12 22ZM12 20C16.4183 20 20 16.4183 20 12C20 7.58172 16.4183 4 12 4C7.58172 4 4 7.58172 4 12C4 16.4183 7.58172 20 12 20ZM12 10.5858L14.8284 7.75736L16.2426 9.17157L13.4142 12L16.2426 14.8284L14.8284 16.2426L12 13.4142L9.17157 16.2426L7.75736 14.8284L10.5858 12L7.75736 9.17157L9.17157 7.75736L12 10.5858Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/menu-line

func Menu(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_42 := templ.GetChildren(ctx)
		if var_42 == nil {
			var_42 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_43 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_43...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_43).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M3 4H21V6H3V4ZM3 11H21V13H3V11ZM3 18H21V20H3V18Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/error-warning-line

func ErrorWarning(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_44 := templ.GetChildren(ctx)
		if var_44 == nil {
			var_44 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_45 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_45...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_45).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12C22 17.5228 17.5228 22 12 22ZM12 20C16.4183 20 20 16.4183 20 12C20 7.58172 16.4183 4 12 4C7.58172 4 4 7.58172 4 12C4 16.4183 7.58172 20 12 20ZM11 15H13V17H11V15ZM11 7H13V13H11V7Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/alert-line

func Alert(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_46 := templ.GetChildren(ctx)
		if var_46 == nil {
			var_46 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_47 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_47...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_47).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12.865 3.00017L22.3912 19.5002C22.6674 19.9785 22.5035 20.5901 22.0252 20.8662C21.8732 20.954 21.7008 21.0002 21.5252 21.0002H2.47266C1.92037 21.0002 1.47266 20.5525 1.47266 20.0002C1.47266 19.8246 1.51886 19.6522 1.60663 19.5002L11.1329 3.00017C11.4091 2.52187 12.0206 2.358 12.4989 2.63414C12.651 2.72191 12.7772 2.84815 12.865 3.00017ZM4.20471 19.0002H19.7932L11.9989 5.50017L4.20471 19.0002ZM10.9989 16.0002H12.9989V18.0002H10.9989V16.0002ZM10.9989 9.00017H12.9989V14.0002H10.9989V9.00017Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/filter-3-line

func Filter(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_48 := templ.GetChildren(ctx)
		if var_48 == nil {
			var_48 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_49 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_49...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_49).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M10 18H14V16H10V18ZM3 6V8H21V6H3ZM6 13H18V11H6V13Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

// https://remixicon.com/icon/question-mark

func QuestionMark(class any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_50 := templ.GetChildren(ctx)
		if var_50 == nil {
			var_50 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var var_51 = []any{class, "fill-current"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_51...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<svg class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_51).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12 19C12.8284 19 13.5 19.6716 13.5 20.5C13.5 21.3284 12.8284 22 12 22C11.1716 22 10.5 21.3284 10.5 20.5C10.5 19.6716 11.1716 19 12 19ZM12 2C15.3137 2 18 4.68629 18 8C18 10.1646 17.2474 11.2907 15.3259 12.9231C13.3986 14.5604 13 15.2969 13 17H11C11 14.526 11.787 13.3052 14.031 11.3989C15.5479 10.1102 16 9.43374 16 8C16 5.79086 14.2091 4 12 4C9.79086 4 8 5.79086 8 8V9H6V8C6 4.68629 8.68629 2 12 2Z\"></path></svg>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
