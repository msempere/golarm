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
	freeMetric metric = iota + 1
	usedMetric
	timeMetric
	statusMetric
)

// Linux process states to be used with status alarms
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

// Load average can be calculated for the last one minute, five minutes and fifteen minutes respectively. Load average is an indication of whether the system resources (mainly the CPU) are adequately available for the processes (system load) that are running, runnable or in uninterruptible sleep states during the previous n minutes.
const (
	OneMinPeriod period = iota + 1
	FiveMinPeriod
	FifteenMinPeriod
)

func getLoadAverage(p period, manager sigarMetrics, percentage bool) float64 {
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

func getPidState(pid uint, manager sigarMetrics) float64 {
	value, err := manager.getProcState(int(pid))
	if err != nil {
		return 6.0
	}
	return states[string(value.State)]
}

func getPidMemory(pid uint, manager sigarMetrics, percentage bool) float64 {
	memory, err := manager.getProcMem(int(pid))

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
func getPidTime(pid uint, manager sigarMetrics) float64 {
	value, err := manager.getProcTime(int(pid))

	if err != nil {
		return 0.0
	}
	return float64(value.Total / 1000)
}

func getTotalMemory(manager sigarMetrics) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}
	return float64(mem.Total)
}

func getTotalSwap(manager sigarMetrics) float64 {
	swap, err := manager.GetSwap()

	if err != nil {
		return 0.0
	}
	return float64(swap.Total)
}

func getUsedSwap(manager sigarMetrics) float64 {
	swap, err := manager.GetSwap()

	if err != nil {
		return 0.0
	}
	return float64(swap.Used)
}

func getFreeSwap(manager sigarMetrics) float64 {
	swap, err := manager.GetSwap()

	if err != nil {
		return 0.0
	}
	return float64(swap.Free)
}

func getActualUsedMemory(manager sigarMetrics, percentage bool) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}

	value := float64(mem.ActualUsed) / 1048576

	if percentage {
		return 100.0 * (value / (getTotalMemory(manager) / 1048576))
	}
	return value
}

func getActualFreeMemory(manager sigarMetrics, percentage bool) float64 {
	mem, err := manager.GetMem()

	if err != nil {
		return 0.0
	}
	value := float64(mem.ActualFree) / 1048576

	if percentage {
		return 100.0 * (value / (getTotalMemory(manager) / 1048576))
	}
	return value
}

func getActualFreeSwap(manager sigarMetrics, percentage bool) float64 {
	value := float64(getFreeSwap(manager)) / 1048576

	if percentage {
		return 100.0 * (value / getTotalSwap(manager) / 1048576)
	}
	return value
}

func getActualUsedSwap(manager sigarMetrics, percentage bool) float64 {
	value := float64(getUsedSwap(manager)) / 1048576

	if percentage {
		return 100.0 * (value / (getTotalSwap(manager) / 1048576))
	}
	return value
}

func getUptime(manager sigarMetrics) float64 {
	value, err := manager.getUpTime()
	if err != nil {
		return 0.0
	}
	return value.Length
}
