package ufacter

import (
	"strings"
)

// Fact represents a reported fact
type Fact struct {
	// Slice of structured names, nil indicate end of facts (use instead of channel close).
	Name []string
	// Reported value
	Value interface{}
	// Represents facts which change too often
	Volatile bool
	// Native (PuppetLabs facter) or extended (extra) fact
	Native bool
}

// TODO use references instead of copying in channels?

// NewFact creates new generic fact
func NewFact(value interface{}, volatile bool, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: volatile,
		Native:   true,
	}
}

// NewStableFact creates fact which is not volatile
func NewStableFact(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: false,
		Native:   true,
	}
}

// NewFactEx creates a non-native (extended) fact
func NewFactEx(value interface{}, volatile bool, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: volatile,
		Native:   false,
	}
}

// NewStableFactEx creates non-volatile non-native fact
func NewStableFactEx(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: false,
		Native:   false,
	}
}

// NewVolatileFact creates a volatile native fact
func NewVolatileFact(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: true,
		Native:   true,
	}
}

// NewVolatileFactEx creates a volatile non-native fact
func NewVolatileFactEx(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: true,
		Native:   false,
	}
}

// SendVolatileFactEx creates a fact via NewVolatileFactEx and sends it
func SendVolatileFactEx(facts chan<- Fact, value interface{}, keys ...string) {
	facts <- NewVolatileFactEx(value, keys...)
}

// NewLastFact creates a fact that needs to be reported as the final one into the channel
func NewLastFact() Fact {
	return Fact{
		Name:  nil,
		Value: nil,
	}
}

// SendLastFact creates a fact via NewLastFact and sends it
func SendLastFact(facts chan<- Fact) {
	facts <- NewLastFact()
}

// NameDots returns fact name in dot format
func (fact Fact) NameDots() string {
	return strings.Join(fact.Name, ".")
}
