package golarm

import (
	"errors"
	"syscall"

	"github.com/carlescere/scheduler"
	"github.com/cloudfoundry/gosigar"
)

var (
	Duration                         = 5 // default check interval set to 5 seconds
	alarms                           = make([]*alarm, 0)
	not_set                          = 123456.123456
	ErrAlarmTypeNotDefined           = errors.New("Bad chain. alarm type not defined")
	ErrComparisonNotDefined          = errors.New("Bad chain. alarm comparison not defined")
	ErrExpectedNumWhenPercentage     = errors.New("Bad chain. A number is needed for applying a percentage")
	ErrIncorrectTypeForFree          = errors.New("Alarm type not set or trying to use Free with something different than memory or swap memory")
	ErrIncorrectTypeForUsed          = errors.New("Alarm type not set or trying to use Free with something different than memory, swap memory or memory used by a proc")
	ErrIncorrectTypeForTime          = errors.New("Alarm type not set or trying to use RunningTime with something different than uptime or uptime by a proc")
	ErrIncorrectTypeForStatus        = errors.New("Alarm type not set or trying to use Status with something different than SysteProcAlarm")
	ErrMultipleComparisonDefined     = errors.New("Alarm comparison already defined")
	ErrIncorrectTypeForAbove         = errors.New("Alarm type not set or trying to use Above with a status metric")
	ErrIncorrectTypeForBelow         = errors.New("Alarm type not set or trying to use Below with a status metric")
	ErrIncorrectTypeForPercentage    = errors.New("Couldn't apply percentage to uptime/status alarms")
	ErrIncorrectValuesWithPercentage = errors.New("Couldn't apply percentage")
	ErrInexistentPid                 = errors.New("Pid does not exist")
)

type alarmType int
type comparison int

type value struct {
	value      float64
	percentage bool
}

type alarm struct {
	metricsManager SigarMetrics
	procMem        sigar.ProcMem
	quit           chan bool
	result         chan bool
	err            error
	task           func()
	jobType        alarmType
	comparison     comparison
	value          value
	stats          stats
}

const (
	comparisonNotDefined comparison = iota
	above
	aboveEqual
	equal
	below
	belowEqual
)

const (
	alertTypeNotDefined alarmType = iota
	loadAlarm
	memoryAlarm
	swapAlarm
	uptimeAlarm
	procAlarm
)

type SigarMetrics interface {
	sigar.Sigar
	GetProcState(int) (sigar.ProcState, error)
	GetProcMem(int) (sigar.ProcMem, error)
	GetProcTime(int) (sigar.ProcTime, error)
	GetUpTime() (sigar.Uptime, error)
}

type ConcreteSigar struct {
	sigar.ConcreteSigar
}

func (c *ConcreteSigar) GetUpTime() (sigar.Uptime, error) {
	p := sigar.Uptime{}
	err := p.Get()
	return p, err
}

func (c *ConcreteSigar) GetProcState(pid int) (sigar.ProcState, error) {
	p := sigar.ProcState{}
	err := p.Get(pid)
	return p, err
}

func (c *ConcreteSigar) GetProcMem(pid int) (sigar.ProcMem, error) {
	p := sigar.ProcMem{}
	err := p.Get(pid)
	return p, err
}

func (c *ConcreteSigar) GetProcTime(pid int) (sigar.ProcTime, error) {
	p := sigar.ProcTime{}
	err := p.Get(pid)
	return p, err
}

func AddAlarm(a *alarm) error {
	if a.err == nil {
		go func(b *alarm) {
			scheduler.Every(Duration).Seconds().NotImmediately().Run(func() {
				Check(b)
			})
		}(a)

		go func(b *alarm) {
			for {
				select {
				case fired := <-b.result:
					if fired {
						b.execute()
					}
				}
			}
		}(a)

		alarms = append(alarms, a)
		return nil
	}
	return a.err
}

func compare(value1, value2 float64, c comparison) bool {
	switch c {
	case above:
		return value1 > value2
	case aboveEqual:
		return value1 >= value2
	case equal:
		return value1 == value2
	case below:
		return value1 < value2
	case belowEqual:
		return value1 <= value2
	}
	return false
}

func (j *alarm) SetMetricsManager(m SigarMetrics) {
	(*j).metricsManager = m
}

func pidExists(pid int) bool {
	killErr := syscall.Kill(pid, syscall.Signal(0))
	return killErr == nil
}

func Check(alarm *alarm) {
	if alarm.err == nil {
		switch alarm.jobType {
		case loadAlarm:
			alarm.result <- compare(
				getLoadAverage(
					alarm.stats.period,
					alarm.metricsManager,
					alarm.value.percentage),
				alarm.value.value,
				alarm.comparison)

		case uptimeAlarm:
			alarm.result <- compare(
				getUptime(
					alarm.metricsManager),
				alarm.value.value,
				alarm.comparison)

		case procAlarm:
			switch alarm.stats.metric {
			case used_:
				alarm.result <- compare(
					getPidMemory(alarm.stats.proc.pid,
						alarm.metricsManager,
						alarm.value.percentage),
					alarm.value.value,
					alarm.comparison)
			case time:
				alarm.result <- compare(
					getPidTime(alarm.stats.proc.pid,
						alarm.metricsManager),
					alarm.value.value,
					alarm.comparison)
			case status:
				alarm.result <- compare(
					float64(getPidState(alarm.stats.proc.pid,
						alarm.metricsManager)),
					alarm.value.value,
					alarm.comparison)
			}

		case memoryAlarm:
			switch alarm.stats.metric {
			case free_:
				alarm.result <- compare(
					float64(getActualFreeMemory(
						alarm.metricsManager,
						alarm.value.percentage,
					)),
					alarm.value.value,
					alarm.comparison)
			case used_:
				alarm.result <- compare(
					float64(getActualUsedMemory(
						alarm.metricsManager,
						alarm.value.percentage,
					)),
					alarm.value.value,
					alarm.comparison)
			}

		case swapAlarm:
			switch alarm.stats.metric {
			case free_:
				alarm.result <- compare(
					float64(getActualFreeSwap(
						alarm.metricsManager,
						alarm.value.percentage,
					)),
					alarm.value.value,
					alarm.comparison)
			case used_:
				alarm.result <- compare(
					float64(getActualUsedSwap(
						alarm.metricsManager,
						alarm.value.percentage,
					)),
					alarm.value.value,
					alarm.comparison)
			}
		}
	}
}

func SystemLoad(p period) *alarm {
	return &alarm{
		jobType: loadAlarm,
		value: value{
			value:      not_set,
			percentage: false},
		result: make(chan bool),
		stats: stats{
			period: p,
			metric: 0,
		},
	}
}

func SystemProc(pid uint) *alarm {
	alarm := &alarm{
		jobType: procAlarm,
		value: value{
			value:      not_set,
			percentage: false},
		result: make(chan bool),
		stats: stats{
			metric: 0,
			period: 0,
			proc: proc{
				pid:   pid,
				state: Unknown,
			}},
	}
	if !pidExists(int(pid)) {
		alarm.err = errors.New("Pid does not exist")
	}
	return alarm
}

func SystemMemory() *alarm {
	a := &alarm{
		jobType: memoryAlarm,
		value: value{
			value:      not_set,
			percentage: false},
		result: make(chan bool),
		stats: stats{
			metric: 0,
			period: 0,
		},
	}
	a.SetMetricsManager(&ConcreteSigar{})
	return a
}

func SystemSwap() *alarm {
	return &alarm{
		jobType: swapAlarm,
		value: value{
			value:      not_set,
			percentage: false},
		result: make(chan bool),
		stats: stats{
			metric: 0,
			period: 0,
		},
	}
}

func SystemUptime() *alarm {
	return &alarm{
		jobType: uptimeAlarm,
		value: value{
			value:      not_set,
			percentage: false},
		result: make(chan bool),
		stats: stats{
			metric: 0,
			period: 0,
		},
	}
}

func (j *alarm) execute() {
	if j.err == nil {
		j.task()
	}
}

func (j *alarm) Run(f func()) *alarm {
	if j.err == nil {
		if j.comparison == comparisonNotDefined {
			(*j).err = ErrComparisonNotDefined
			return j
		}

		(*j).task = f
	}
	return j
}
