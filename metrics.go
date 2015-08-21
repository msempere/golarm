package golarm

type period int
type metric int
type state int
type procType int

var (
	states = map[string]float64{
		"S": 1.0,
		"R": 2.0,
		"T": 3.0,
		"Z": 4.0,
		"D": 5.0,
	}
)

const (
	free_ metric = iota + 1
	used_
	time
	status
)

const (
	Sleeping state = iota + 1
	Running
	Stopped
	Zombie
	Idle
	Unknown
)

type load struct {
	period period
}

type proc struct {
	state state
	pid   uint
}

type stats struct {
	period period
	proc   proc
	metric metric
}

const (
	OneMinPeriod period = iota + 1
	FiveMinPeriod
	FifteenMinPeriod
)

func getLoadAverage(p period, manager SigarMetrics, percentage bool) float64 {
	average, err := manager.GetLoadAverage()
	value := 0.0

	if err != nil {
		return value
	}

	switch p {
	case OneMinPeriod:
		value = float64(average.One)
	case FiveMinPeriod:
		value = float64(average.Five)
	case FifteenMinPeriod:
		value = float64(average.Fifteen)
	}

	if percentage {
		value *= 10
	}
	return value
}

func getPidState(pid uint, manager SigarMetrics) float64 {
	value, err := manager.GetProcState(int(pid))
	if err != nil {
		return 6.0
	}
	return states[string(value.State)]
}

func getPidMemory(pid uint, manager SigarMetrics, percentage bool) float64 {
	memory, err := manager.GetProcMem(int(pid))

	if err != nil {
		return 0.0
	}

	value := float64(memory.Resident / 1048576)

	if percentage {
		return 100.0 * (value / (getTotalMemory(manager) / 1048576))
	}
	return value
}

// get running time for PID in minutes
func getPidTime(pid uint, manager SigarMetrics) float64 {
	value, err := manager.GetProcTime(int(pid))

	if err != nil {
		return 0.0
	}
	return float64(value.Total / 1000)
}

func getTotalMemory(manager SigarMetrics) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}
	return float64(mem.Total)
}

func getTotalSwap(manager SigarMetrics) float64 {
	swap, err := manager.GetSwap()

	if err != nil {
		return 0.0
	}
	return float64(swap.Total)
}

func getUsedMemory(manager SigarMetrics) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}
	return float64(mem.Used)
}

func getUsedSwap(manager SigarMetrics) float64 {
	swap, err := manager.GetSwap()

	if err != nil {
		return 0.0
	}
	return float64(swap.Used)
}

func getFreeMemory(manager SigarMetrics) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}
	return float64(mem.Free)
}

func getFreeSwap(manager SigarMetrics) float64 {
	swap, err := manager.GetSwap()

	if err != nil {
		return 0.0
	}
	return float64(swap.Free)
}

func getActualUsedMemory(manager SigarMetrics, percentage bool) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}

	value := float64(mem.ActualUsed) / 1048576

	if percentage {
		return 100.0 * (value / getTotalMemory(manager))
	}
	return value
}

func getActualFreeMemory(manager SigarMetrics, percentage bool) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}
	value := float64(mem.ActualFree) / 1048576

	if percentage {
		return 100.0 * (value / getTotalMemory(manager))
	}
	return value
}

func getActualFreeSwap(manager SigarMetrics, percentage bool) float64 {
	value := float64(getFreeSwap(manager)) / 1048576

	if percentage {
		return 100.0 * (value / getTotalSwap(manager))
	}
	return value
}

func getActualUsedSwap(manager SigarMetrics, percentage bool) float64 {
	value := float64(getUsedSwap(manager)) / 1048576

	if percentage {
		return 100.0 * (value / (getTotalSwap(manager) / 1048576))
	}
	return value
}

func getUptime(manager SigarMetrics) float64 {
	value, err := manager.GetUpTime()
	if err != nil {
		return 0.0
	}
	return value.Length
}
