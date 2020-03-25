package mem

import (
	"fmt"

	"github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	m "github.com/shirou/gopsutil/mem"
)

func reportMemory(facts chan<- ufacter.Fact, volatile bool, value uint64, root_key string, bytes_key string, total_key string) error {
	human, unit, err := common.ConvertBytes(value)
	if err != nil {
		return err
	}
	facts <- ufacter.NewFact(value, volatile, "memory", root_key, bytes_key)
	facts <- ufacter.NewFact(fmt.Sprintf("%.2f %v", human, unit), volatile, "memory", root_key, total_key)
	return nil
}

// Gathers facts related to memory
func ReportFacts(facts chan<- ufacter.Fact) error {
	hostVM, err := m.VirtualMemory()
	if err != nil {
		return err
	}

	err = reportMemory(facts, false, hostVM.Total, "system", "total_bytes", "total")
	if err != nil {
		return err
	}

	err = reportMemory(facts, true, hostVM.Used, "system", "used_bytes", "used")
	if err != nil {
		return err
	}
	facts <- ufacter.NewVolatileFact(fmt.Sprintf("%.2f%%", hostVM.UsedPercent), "memory", "system", "capacity")

	err = reportMemory(facts, true, hostVM.Available, "system", "available_bytes", "available")
	if err != nil {
		return err
	}

	// Get the swap information from gopsutil
	hostSwap, err := m.SwapMemory()
	if err != nil {
		return err
	}

	err = reportMemory(facts, false, hostSwap.Total, "swap", "total_bytes", "total")
	if err != nil {
		return err
	}

	err = reportMemory(facts, true, hostSwap.Used, "swap", "used_bytes", "used")
	if err != nil {
		return err
	}
	facts <- ufacter.NewVolatileFact(fmt.Sprintf("%.2f%%", hostSwap.UsedPercent), "memory", "swap", "capacity")

	err = reportMemory(facts, true, hostSwap.Free, "swap", "available_bytes", "available")
	if err != nil {
		return err
	}

	facts <- ufacter.NewLastFact()
	return nil
}
