package ptool

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetCPUPercent() interface{} {
	cc, err := cpu.Percent(time.Second, false)
	if len(cc) < 1 || err != nil {
		return float64(0)
	}

	return cc[0]
}

func GetMemoryPercent() interface{} {
	v, err := mem.VirtualMemory()
	if err != nil {
		return float64(0)
	}

	return v.UsedPercent
}

func GetSwapPercent() interface{} {
	s, err := mem.SwapMemory()
	if err != nil {
		return float64(0)
	}

	return s.UsedPercent
}
