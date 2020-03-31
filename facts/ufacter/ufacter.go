package ufacter

import (
	"github.com/lzap/ufacter/lib/ufacter"
)

// ReportFacts reports facts related to ufacter itself
func ReportFacts(facts chan<- ufacter.Fact) {

	facts <- ufacter.NewStableFact(UFACTER_VERSION, "ufacter", "version")
	facts <- ufacter.NewStableFact("3.0.0", "facterversion")

	facts <- ufacter.NewLastFact()
}
