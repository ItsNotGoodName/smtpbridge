package app

import (
	"errors"
)

var ErrBridgeNotFound = errors.New("bridge not found")

type Bridge struct {
	Name      string   `json:"name" yaml:"name"`
	Filters   []Filter `json:"filters" yaml:"filters"`
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
}

func (b *Bridge) Match(msg *Message) bool {
	for _, f := range b.Filters {
		if f.Match(msg) {
			return true
		}
	}
	return false
}
