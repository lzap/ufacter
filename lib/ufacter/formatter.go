package ufacter

// Formatter prints data to standard output
type Formatter interface {
	// Store or print fact
	Add(Fact)

	// Finish printing, some implementations do nothing
	Finish()
}
