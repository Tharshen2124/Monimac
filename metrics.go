package main

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CPUStats holds CPU usage percentage.
type CPUStats struct {
	Percent float64
}

// MemStats holds memory usage data.
type MemStats struct {
	Percent float64
	Used    uint64
	Total   uint64
}

// GetCPUStats returns current CPU usage averaged across all cores.
func GetCPUStats() (CPUStats, error) {
	percents, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		return CPUStats{}, err
	}
	if len(percents) == 0 {
		return CPUStats{}, nil
	}
	return CPUStats{Percent: percents[0]}, nil
}

// GetMemStats returns current virtual memory statistics.
func GetMemStats() (MemStats, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemStats{}, err
	}
	return MemStats{
		Percent: v.UsedPercent,
		Used:    v.Used,
		Total:   v.Total,
	}, nil
}
