package core

import "testing"

func TestFilter(t *testing.T) {
	type FilterTest struct {
		Filter Filter
		Match  bool
	}

	newFilter := func(to, from, toRegex, fromRegex string) Filter {
		filter, err := NewFilter(to, from, toRegex, fromRegex)
		if err != nil {
			t.Fatalf("newFilter(%s, %s, %s, %s): %s", to, from, toRegex, fromRegex, err)
		}
		return filter
	}

	msg := Message{
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
		if test.Filter.Match(&msg) != test.Match {
			t.Errorf("Filter.Match(%v) != %v", msg, test.Match)
		}
	}
}
