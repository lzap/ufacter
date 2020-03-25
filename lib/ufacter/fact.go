package ufacter

import "strings"

// Reported fact
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

// Create new fact
func NewFact(value interface{}, volatile bool, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: volatile,
		Native:   true,
	}
}

// Create new stable native fact
func NewStableFact(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: false,
		Native:   true,
	}
}

// Create new fact
func NewFactEx(value interface{}, volatile bool, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: volatile,
		Native:   false,
	}
}

// Create new stable native fact
func NewStableFactEx(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: false,
		Native:   false,
	}
}

// Create new volatile native fact
func NewVolatileFact(value interface{}, keys ...string) Fact {
	return Fact{
		Name:     keys,
		Value:    value,
		Volatile: true,
	}
}

// Creates a fact that needs to be reported as the final one into the channel
func NewLastFact() Fact {
	return Fact{
		Name:  nil,
		Value: nil,
	}
}

func (fact Fact) NameDots() string {
	return strings.Join(fact.Name, ".")
}
