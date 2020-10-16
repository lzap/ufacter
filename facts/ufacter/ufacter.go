package ufacter

import (
	"github.com/lzap/ufacter/lib/ufacter"
)

// ReportFacts reports facts related to ufacter itself
func ReportFacts(facts chan<- ufacter.Fact, volatile bool, extended bool) {
	defer ufacter.SendLastFact(facts)

	facts <- ufacter.NewStableFactEx(UFACTER_VERSION, "ufacter", "version")
	facts <- ufacter.NewStableFact("3.0.0", "facterversion")
}
