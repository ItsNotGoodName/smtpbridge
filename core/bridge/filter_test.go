package bridge

import (
	"testing"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	type FilterTest struct {
		Filter Filter
		Match  bool
	}

	newFilter := func(to, from, toRegex, fromRegex string) Filter {
		filter, err := NewFilter(from, to, fromRegex, toRegex)
		assert.Nil(t, err)
		return filter
	}

	msg := envelope.Message{
		From: "foo",
		To:   map[string]struct{}{"bar": {}},
	}

	tests := []FilterTest{
		{Filter: newFilter("bar", "", "", ""), Match: true},
		{Filter: newFilter("bar", "foo", "", ""), Match: true},
		{Filter: newFilter("bar", "foorr", "", ""), Match: false},
		{Filter: newFilter("barr", "", "", ""), Match: false},
		{Filter: newFilter("", "barrr", "", ""), Match: false},
		{Filter: newFilter("", "", "", "f.$"), Match: false},
		{Filter: newFilter("", "", "", "f"), Match: true},
		{Filter: newFilter("", "", "b", ""), Match: true},
		{Filter: newFilter("bar", "", "f", ""), Match: false},
		{Filter: newFilter("bar", "", "b", ""), Match: true},
		{Filter: newFilter("bar", "", "b", "x"), Match: false},
	}

	for _, test := range tests {
		assert.Equal(t, test.Filter.Match(&msg), test.Match)
	}
}
