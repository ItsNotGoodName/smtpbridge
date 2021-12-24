package app

import "testing"

func TestBridge(t *testing.T) {
	type BridgeTest struct {
		Bridge Bridge
		Match  bool
	}

	msg := Message{
		From: "foo",
		To:   map[string]bool{"bar": true},
	}

	tests := []BridgeTest{
		{Bridge: Bridge{EmailTo: "bar"}, Match: true},
		{Bridge: Bridge{EmailFrom: "foo", EmailTo: "bar"}, Match: true},
		{Bridge: Bridge{EmailFrom: "foorr", EmailTo: "bar"}, Match: false},
		{Bridge: Bridge{EmailTo: "barr"}, Match: false},
		{Bridge: Bridge{EmailFrom: "barrr"}, Match: false},
	}

	for _, test := range tests {
		if test.Bridge.Match(&msg) != test.Match {
			t.Errorf("Bridge.Match(%v) != %v", msg, test.Match)
		}
	}
}
