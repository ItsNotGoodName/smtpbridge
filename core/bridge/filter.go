package bridge

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"text/template"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Filter struct {
	To            string
	From          string
	toRegexp      *regexp.Regexp
	fromRegexp    *regexp.Regexp
	matchTemplate *template.Template
}

func NewFilter(from, to, fromRegex, toRegex, templateStr string) (Filter, error) {
	var fromRegexp, toRegexp *regexp.Regexp
	var err error
	if fromRegex != "" {
		fromRegexp, err = regexp.Compile(fromRegex)
		if err != nil {
			return Filter{}, err
		}
	}
	if toRegex != "" {
		toRegexp, err = regexp.Compile(toRegex)
		if err != nil {
			return Filter{}, err
		}
	}

	var matchTmpl *template.Template
	if templateStr != "" {
		matchTmpl, err = template.New("").Parse(templateStr)
		if err != nil {
			return Filter{}, err
		}
	}

	return Filter{
		From:          from,
		To:            to,
		fromRegexp:    fromRegexp,
		toRegexp:      toRegexp,
		matchTemplate: matchTmpl,
	}, nil
}

func (f *Filter) Match(env *envelope.Envelope) bool {
	if f.toRegexp != nil {
		found := false
		for to := range env.Message.To {
			if f.toRegexp.MatchString(to) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	} else if f.To != "" {
		if _, ok := env.Message.To[f.To]; !ok {
			return false
		}
	}

	if f.fromRegexp != nil {
		if !f.fromRegexp.MatchString(env.Message.From) {
			return false
		}
	} else if f.From != "" {
		if env.Message.From != f.From {
			return false
		}
	}

	if f.matchTemplate != nil {
		var buffer bytes.Buffer
		if err := f.matchTemplate.Execute(&buffer, env); err != nil {
			log.Println("bridge.Filter.Match:", err)
			return false
		}

		truthy, err := strconv.ParseBool(buffer.String())
		if err != nil {
			log.Println("bridge.Filter.Match:", err)
			return false
		}

		if !truthy {
			return false
		}
	}

	return true
}
