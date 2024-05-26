package modules

import (
	"syscall"

	"github.com/5N41P4/raspberry/internal/data"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
)

// Functions to get system variables and information

// GetDiskSpace returns the total and free disk space in GB and the free space in percentage
func GetDiskSpace() (data.DiskUsage, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs("/", &stat)
	if err != nil {
		return data.DiskUsage{}, err
	}

	// Calculate total, free and used space
	total := float64(stat.Blocks * uint64(stat.Bsize))
	free := float64(stat.Bfree * uint64(stat.Bsize))
	used := float64(total - free)

	// Calculate total free and used in Gigabytes
	total /= 1 << 30
	free /= 1 << 30
	used /= 1 << 30

	// Calculate percentage (avoid division by zero)
	percent := 0.0
	if total > 0 {
		percent = float64(used) / float64(total) * 100.0
	}

	return data.DiskUsage{
		Total:   total,
		Free:    free,
		Used:    used,
		Percent: percent,
	}, nil
}

// GetCpuUsage returns the total and free CPU usage in percentage

func GetCpuUsage() (data.CpuUsage, error) {
	cpu, err := loadavg.Get()
	if err != nil {
		return data.CpuUsage{}, err
	}

	return data.CpuUsage{
		AvgLoad1:  cpu.Loadavg1,
		AvgLoad5:  cpu.Loadavg5,
		AvgLoad15: cpu.Loadavg15,
	}, nil
}

// GetMemUsage returns the memory usage in percent
func GetMemUsage() (data.MemUsage, error) {
	memStats, err := memory.Get()
	if err != nil {
		return data.MemUsage{}, err
	}
	var mem = (float64(memStats.Used) / float64(memStats.Total)) * 100

	return data.MemUsage{
		Memory: mem,
	}, nil
}
