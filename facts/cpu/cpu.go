package cpu

import (
	"sort"
	"strconv"

	"github.com/lzap/ufacter/lib/ufacter"
	c "github.com/shirou/gopsutil/cpu"
)

// Gathers facts related to CPU
func ReportFacts(facts chan<- ufacter.Fact) error {
	totalCount, err := c.Counts(true)
	if err != nil {
		return err
	}

	facts <- ufacter.NewStableFact(totalCount, "processors", "count")

	CPUs, err := c.Info()
	if err != nil {
		return err
	}
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

	facts <- ufacter.NewLastFact()
	return nil
}
