package golarm

import (
	"errors"
	"syscall"

	"github.com/carlescere/scheduler"
	"github.com/cloudfoundry/gosigar"
)

// check interval set to 5 seconds as default
var Duration = 5

var Alarms = make([]*Alarm, 0)
var notSet = 123456.123456

// Error codes returned by failures when trying to create the alert chain
var (
	ErrAlarmTypeNotDefined           = errors.New("Bad chain. Alarm type not defined")
	ErrComparisonNotDefined          = errors.New("Bad chain. Alarm comparison not defined")
	ErrExpectedNumWhenPercentage     = errors.New("Bad chain. A number is needed for applying a percentage")
	ErrIncorrectTypeForFree          = errors.New("Alarm type not set or trying to use Free with something different than memory or swap memory")
	ErrIncorrectTypeForUsed          = errors.New("Alarm type not set or trying to use Free with something different than memory, swap memory or memory used by a proc")
	ErrIncorrectTypeForTime          = errors.New("Alarm type not set or trying to use RunningTime with something different than uptime or uptime by a proc")
	ErrIncorrectTypeForStatus        = errors.New("Alarm type not set or trying to use Status with something different than SysteProcAlarm")
	ErrMultipleComparisonDefined     = errors.New("Alarm comparison already defined")
	ErrIncorrectTypeForAbove         = errors.New("Alarm type not set or trying to use Above with a status metric")
	ErrIncorrectTypeForBelow         = errors.New("Alarm type not set or trying to use Below with a status metric")
	ErrIncorrectTypeForPercentage    = errors.New("Couldn't apply percentage to uptime/status Alarms")
	ErrIncorrectValuesWithPercentage = errors.New("Couldn't apply percentage")
	ErrInexistentPid                 = errors.New("Pid does not exist")
	ErrIncorrectTypeForComparison    = errors.New("Alarm type not set or trying to use an incorrect comparison with this type of Alarm")
	ErrIncorrectTypeForMetric        = errors.New("Alarm type not set or trying to use an incorrect metric with this type of Alarm")
)

type AlarmType int
type comparison int

type value struct {
	value      float64
	percentage bool
}

type Alarm struct {
	metricsManager sigarMetrics
	quit           chan bool
	result         chan bool
	Err            error
	task           func()
	jobType        AlarmType
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
	alertTypeNotDefined AlarmType = iota
	loadAlarm
	memoryAlarm
	swapAlarm
	uptimeAlarm
	procAlarm
)

type sigarMetrics interface {
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

func AddAlarm(a *Alarm) error {
	if a.Err == nil {
		go func(b *Alarm) {
			scheduler.Every(Duration).Seconds().NotImmediately().Run(func() {
				check(b)
			})
		}(a)

		go func(b *Alarm) {
			for {
				select {
				case fired := <-b.result:
					if fired {
						b.execute()
					}
				}
			}
		}(a)

		Alarms = append(Alarms, a)
		return nil
	}
	return a.Err
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

func (j *Alarm) SetMetricsManager(m sigarMetrics) {
	(*j).metricsManager = m
}

func pidExists(pid int) bool {
	killErr := syscall.Kill(pid, syscall.Signal(0))
	return killErr == nil
}

func check(Alarm *Alarm) {
	if Alarm.Err == nil {
		switch Alarm.jobType {
		case loadAlarm:
			Alarm.result <- compare(
				getLoadAverage(
					Alarm.stats.period,
					Alarm.metricsManager,
					Alarm.value.percentage),
				Alarm.value.value,
				Alarm.comparison)

		case uptimeAlarm:
			Alarm.result <- compare(
				getUptime(
					Alarm.metricsManager),
				Alarm.value.value,
				Alarm.comparison)

		case procAlarm:
			switch Alarm.stats.metric {
			case used_:
				Alarm.result <- compare(
					getPidMemory(Alarm.stats.proc.pid,
						Alarm.metricsManager,
						Alarm.value.percentage),
					Alarm.value.value,
					Alarm.comparison)
			case time_:
				Alarm.result <- compare(
					getPidTime(Alarm.stats.proc.pid,
						Alarm.metricsManager),
					Alarm.value.value,
					Alarm.comparison)
			case status:
				Alarm.result <- compare(
					float64(getPidState(Alarm.stats.proc.pid,
						Alarm.metricsManager)),
					float64(Alarm.stats.proc.state),
					equal)
			}

		case memoryAlarm:
			switch Alarm.stats.metric {
			case free_:
				Alarm.result <- compare(
					float64(getActualFreeMemory(
						Alarm.metricsManager,
						Alarm.value.percentage,
					)),
					Alarm.value.value,
					Alarm.comparison)
			case used_:
				Alarm.result <- compare(
					float64(getActualUsedMemory(
						Alarm.metricsManager,
						Alarm.value.percentage,
					)),
					Alarm.value.value,
					Alarm.comparison)
			}

		case swapAlarm:
			switch Alarm.stats.metric {
			case free_:
				Alarm.result <- compare(
					float64(getActualFreeSwap(
						Alarm.metricsManager,
						Alarm.value.percentage,
					)),
					Alarm.value.value,
					Alarm.comparison)
			case used_:
				Alarm.result <- compare(
					float64(getActualUsedSwap(
						Alarm.metricsManager,
						Alarm.value.percentage,
					)),
					Alarm.value.value,
					Alarm.comparison)
			}
		}
	}
}

func SystemLoad(p period) *Alarm {
	a := &Alarm{
		jobType: loadAlarm,
		value: value{
			value:      notSet,
			percentage: false},
		result: make(chan bool),
		stats: stats{
			period: p,
			metric: 0,
		},
	}
	a.SetMetricsManager(&ConcreteSigar{})
	return a
}

func SystemProc(pid uint) *Alarm {
	a := &Alarm{
		jobType: procAlarm,
		value: value{
			value:      notSet,
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
		a.Err = ErrInexistentPid
	}
	a.SetMetricsManager(&ConcreteSigar{})
	return a
}

func SystemMemory() *Alarm {
	a := &Alarm{
		jobType: memoryAlarm,
		value: value{
			value:      notSet,
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

func SystemSwap() *Alarm {
	a := &Alarm{
		jobType: swapAlarm,
		value: value{
			value:      notSet,
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

func SystemUptime() *Alarm {
	a := &Alarm{
		jobType: uptimeAlarm,
		value: value{
			value:      notSet,
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

func (j *Alarm) execute() {
	if j.Err == nil {
		j.task()
	}
}

func (j *Alarm) Run(f func()) *Alarm {
	if j.Err == nil {
		if (j.comparison == comparisonNotDefined) && j.stats.metric != status {
			(*j).Err = ErrComparisonNotDefined
			return j
		}
		(*j).task = f
	}
	return j
}
