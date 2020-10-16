package cpu

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	"github.com/shirou/gopsutil/cpu"
)

// ReportFacts gathers facts related to CPU
func ReportFacts(facts chan<- ufacter.Fact, volatile bool, extended bool) {
	start := time.Now()
	defer ufacter.SendLastFact(facts)

	totalCount, err := cpu.Counts(true)
	if err == nil {
		facts <- ufacter.NewStableFact(totalCount, "processors", "count")
	} else {
		c.LogError(facts, err, "cpu", "total count")
	}

	CPUs, err := cpu.Info()
	if err == nil {
		physIDs := make(map[uint64]string)
		maxSpeed := 0.0
		for _, v := range CPUs {
			physID, err := strconv.ParseUint(v.PhysicalID, 10, 32)
			if err == nil {
				physIDs[physID] = v.ModelName
			}
			if v.Mhz > maxSpeed {
				maxSpeed = v.Mhz
			}
		}
		models := []string{}
		for _, value := range physIDs {
			models = append(models, value)
		}
		sort.Strings(models)
		facts <- ufacter.NewStableFact(models, "processors", "models")
		facts <- ufacter.NewStableFact(len(physIDs), "processors", "physicalcount")
		// facter4 reports speed this as volatile fact but we don't
		facts <- ufacter.NewStableFact(fmt.Sprintf("%.2f MHz", maxSpeed), "processors", "speed")
	} else {
		c.LogError(facts, err, "cpu", "info")
	}
	ufacter.SendVolatileFactEx(facts, time.Since(start), "ufacter", "stats", "cpu")
}
