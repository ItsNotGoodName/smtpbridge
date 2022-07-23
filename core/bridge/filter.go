package bridge

import (
	"regexp"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Filter struct {
	To         string
	From       string
	toRegexp   *regexp.Regexp
	fromRegexp *regexp.Regexp
}

func NewFilter(from, to, fromRegex, toRegex string) (Filter, error) {
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

	return Filter{
		From:       from,
		To:         to,
		fromRegexp: fromRegexp,
		toRegexp:   toRegexp,
	}, nil
}

func (f *Filter) Match(msg *envelope.Message) bool {
	if f.toRegexp != nil {
		found := false
		for to := range msg.To {
			if f.toRegexp.MatchString(to) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	} else if f.To != "" {
		if _, ok := msg.To[f.To]; !ok {
			return false
		}
	}

	if f.fromRegexp != nil {
		if !f.fromRegexp.MatchString(msg.From) {
			return false
		}
	} else if f.From != "" {
		if msg.From != f.From {
			return false
		}
	}

	return true
}
