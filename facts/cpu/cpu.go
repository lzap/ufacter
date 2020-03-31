package cpu

import (
	"sort"
	"strconv"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	"github.com/shirou/gopsutil/cpu"
)

// ReportFacts gathers facts related to CPU
func ReportFacts(facts chan<- ufacter.Fact) {
	totalCount, err := cpu.Counts(true)
	if err == nil {
		facts <- ufacter.NewStableFact(totalCount, "processors", "count")
	} else {
		c.LogError(facts, err, "cpu", "total count")
	}

	CPUs, err := cpu.Info()
	if err == nil {
		physIDs := make(map[uint64]string)
		for _, v := range CPUs {
			physID, err := strconv.ParseUint(v.PhysicalID, 10, 32)
			if err == nil {
				physIDs[physID] = v.ModelName
			}
		}
		models := []string{}
		for _, value := range physIDs {
			models = append(models, value)
		}
		sort.Strings(models)
		facts <- ufacter.NewStableFact(models, "processors", "models")
		facts <- ufacter.NewStableFact(len(physIDs), "processors", "physicalcount")
	} else {
		c.LogError(facts, err, "cpu", "info")
	}

	facts <- ufacter.NewLastFact()
}
