package core

import (
	"log"
	"regexp"
)

type Filter struct {
	To         string
	From       string
	toRegexp   *regexp.Regexp
	fromRegexp *regexp.Regexp
}

func NewFilter(to, from, toRegex, fromRegex string) Filter {
	var toRegexp, fromRegexp *regexp.Regexp
	var err error
	// TODO: move error handling to config
	if toRegex != "" {
		toRegexp, err = regexp.Compile(toRegex)
		if err != nil {
			log.Fatalln("core.NewFilter: bad to regex:", err)
		}
	}
	if fromRegex != "" {
		fromRegexp, err = regexp.Compile(fromRegex)
		if err != nil {
			log.Fatalln("core.NewFilter: bad from regex:", err)
		}
	}

	return Filter{
		To:         to,
		From:       from,
		toRegexp:   toRegexp,
		fromRegexp: fromRegexp,
	}
}

func (f *Filter) Match(msg *Message) bool {
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
