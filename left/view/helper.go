package view

import (
	"html/template"
	"time"
)

var helperMap template.FuncMap = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("Jan _2 2006 15:04:05")
	},
	"formatSubject": func(subject string) string {
		if subject == "" {
			return "(no subject)"
		}

		return subject
	},
}
