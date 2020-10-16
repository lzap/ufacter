package mem

import (
	"fmt"
	"time"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	m "github.com/shirou/gopsutil/mem"
)

func reportMemory(facts chan<- ufacter.Fact, volatile bool, value uint64, rootKey string, bytesKey string, totalKey string) {
	human, unit, err := c.ConvertBytes(value)
	if err != nil {
		c.LogError(facts, err, "cpu", "convert bytes")
		return
	}
	facts <- ufacter.NewFact(value, volatile, "memory", rootKey, bytesKey)
	facts <- ufacter.NewFact(fmt.Sprintf("%.2f %v", human, unit), volatile, "memory", rootKey, totalKey)
}

// ReportFacts gathers facts related to memory
func ReportFacts(facts chan<- ufacter.Fact, volatile bool, extended bool) {
	start := time.Now()
	defer ufacter.SendLastFact(facts)

	hostVM, err := m.VirtualMemory()
	if err == nil {
		reportMemory(facts, false, hostVM.Total, "system", "total_bytes", "total")
		reportMemory(facts, true, hostVM.Used, "system", "used_bytes", "used")
		facts <- ufacter.NewVolatileFact(fmt.Sprintf("%.2f%%", hostVM.UsedPercent), "memory", "system", "capacity")
		reportMemory(facts, true, hostVM.Available, "system", "available_bytes", "available")
	} else {
		c.LogError(facts, err, "cpu", "virtual memory")
	}

	// Get the swap information from gopsutil
	hostSwap, err := m.SwapMemory()
	if err == nil {
		reportMemory(facts, false, hostSwap.Total, "swap", "total_bytes", "total")
		reportMemory(facts, true, hostSwap.Used, "swap", "used_bytes", "used")
		facts <- ufacter.NewVolatileFact(fmt.Sprintf("%.2f%%", hostSwap.UsedPercent), "memory", "swap", "capacity")
		reportMemory(facts, true, hostSwap.Free, "swap", "available_bytes", "available")
	}

	ufacter.SendVolatileFactEx(facts, time.Since(start), "ufacter", "stats", "mem")
}
