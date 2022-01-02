package domain

import "testing"

func TestFilter(t *testing.T) {
	type FilterTest struct {
		Filter Filter
		Match  bool
	}

	msg := Message{
		From: "foo",
		To:   map[string]struct{}{"bar": {}},
	}

	tests := []FilterTest{
		{Filter: NewFilter("bar", "", "", ""), Match: true},
		{Filter: NewFilter("bar", "foo", "", ""), Match: true},
		{Filter: NewFilter("bar", "foorr", "", ""), Match: false},
		{Filter: NewFilter("barr", "", "", ""), Match: false},
		{Filter: NewFilter("", "barrr", "", ""), Match: false},
		{Filter: NewFilter("", "", "", "f.$"), Match: false},
		{Filter: NewFilter("", "", "", "f"), Match: true},
		{Filter: NewFilter("", "", "b", ""), Match: true},
		{Filter: NewFilter("bar", "", "f", ""), Match: false},
		{Filter: NewFilter("bar", "", "b", ""), Match: true},
		{Filter: NewFilter("bar", "", "b", "x"), Match: false},
	}

	for _, test := range tests {
		if test.Filter.Match(&msg) != test.Match {
			t.Errorf("Filter.Match(%v) != %v", msg, test.Match)
		}
	}
}
