package domain

import "testing"

func TestFilter(t *testing.T) {
	type FilterTest struct {
		Filter Filter
		Match  bool
	}

	msg := Message{
		From: "foo",
		To:   map[string]bool{"bar": true},
	}

	tests := []FilterTest{
		{Filter: Filter{To: "bar"}, Match: true},
		{Filter: Filter{From: "foo", To: "bar"}, Match: true},
		{Filter: Filter{From: "foorr", To: "bar"}, Match: false},
		{Filter: Filter{To: "barr"}, Match: false},
		{Filter: Filter{From: "barrr"}, Match: false},
	}

	for _, test := range tests {
		if test.Filter.Match(&msg) != test.Match {
			t.Errorf("Filter.Match(%v) != %v", msg, test.Match)
		}
	}
}
