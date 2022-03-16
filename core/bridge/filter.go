package bridge

import (
	"fmt"
	"regexp"

	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

var ErrInvalidRegex = fmt.Errorf("invalid regex")

type Filter struct {
	To         string
	From       string
	toRegexp   *regexp.Regexp
	fromRegexp *regexp.Regexp
}

func NewFilter(to, from, toRegex, fromRegex string) (Filter, error) {
	var toRegexp, fromRegexp *regexp.Regexp
	var err error
	if toRegex != "" {
		toRegexp, err = regexp.Compile(toRegex)
		if err != nil {
			return Filter{}, fmt.Errorf("%w: %s", ErrInvalidRegex, err)
		}
	}
	if fromRegex != "" {
		fromRegexp, err = regexp.Compile(fromRegex)
		if err != nil {
			return Filter{}, fmt.Errorf("%w: %s", ErrInvalidRegex, err)
		}
	}

	return Filter{
		To:         to,
		From:       from,
		toRegexp:   toRegexp,
		fromRegexp: fromRegexp,
	}, nil
}

func (f *Filter) Match(msg *message.Message) bool {
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
