package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

var (
	Duration = 5
	Alerts   = make([]*Alert, 0)
)

type Period int
type AlertType int
type Comparison int

const (
	Above Comparison = iota
	AboveEqual
	Equal
	Below
	BelowEqual
)

const (
	OneMinPeriod Period = iota
	FiveMinPeriod
	FifteenMinPeriod
)

const (
	Load AlertType = iota
	Memory
	Uptime
)

func AddAlert(alert *Alert) {
	Alerts = append(Alerts, alert)
}

func StartAlerts() {
	for _, alert := range Alerts {
		go func(load sigar.LoadAverage, period Period, value float32, comparison Comparison, result chan bool) {
			for {
				select {
				case <-time.After(time.Duration(Duration) * time.Second):
					if alert.jobType == Load {
						CheckLoadAverage(&load, period, value, comparison, result)
					}
				}
			}
		}(alert.Load, alert.period, alert.value, alert.comparison, alert.Result)

		go func(result chan bool, task func()) {
			for {
				select {
				case r := <-result:
					if r {
						task()
					}
				}
			}
		}(alert.Result, alert.Task)
	}
}

func Compare(value1, value2 float32, c Comparison) bool {
	switch c {
	case Above:
		return value1 > value2
	case AboveEqual:
		return value1 >= value2
	case Equal:
		return value1 == value2
	case Below:
		return value1 < value2
	case BelowEqual:
		return value1 <= value2
	}
	return false
}

func CheckLoadAverage(load *sigar.LoadAverage, p Period, value2 float32, c Comparison, channel chan<- bool) {
	channel <- Compare(GetLoadAverage(p, load), value2, c)
}

func GetLoadAverage(p Period, load *sigar.LoadAverage) float32 {
	load.Get()
	switch p {
	case OneMinPeriod:
		return float32(load.One)
	case FiveMinPeriod:
		return float32(load.Five)
	case FifteenMinPeriod:
		return float32(load.Fifteen)
	}
	return 0.0
}

func GetTotalMemory() uint64 {
	mem := sigar.Mem{}
	mem.Get()
	return mem.Total / 1000000
}

func GetUsedMemory() uint64 {
	mem := sigar.Mem{}
	mem.Get()
	return mem.Used / 1000000
}

func GetFreeMemory() uint64 {
	mem := sigar.Mem{}
	mem.Get()
	return mem.Free / 1000000
}

func GetActualUsedMemory(mem *sigar.Mem) uint64 {
	mem.Get()
	return mem.ActualUsed / 1000000
}

func GetActualFreeMemory(mem *sigar.Mem) uint64 {
	mem.Get()
	return mem.ActualFree / 1000000
}

type Alert struct {
	Quit       chan bool
	Result     chan bool
	err        error
	Task       func()
	jobType    AlertType
	period     Period
	comparison Comparison
	value      float32
	Load       sigar.LoadAverage
	Memory     sigar.Mem
}

func SystemLoad(p Period) *Alert {
	return &Alert{jobType: Load,
		period: p,
		Result: make(chan bool),
		Load:   sigar.LoadAverage{},
	}
}

func SystemMemory() *Alert {
	return &Alert{jobType: Memory,
		Result: make(chan bool),
		Memory: sigar.Mem{},
	}
}

func (j *Alert) Above(value float32) *Alert {
	(*j).value = value
	(*j).comparison = Above
	return j
}

func (j *Alert) Below(value float32) *Alert {
	(*j).value = value
	(*j).comparison = Below
	return j
}

func (j *Alert) Run(f func()) *Alert {
	(*j).Task = f
	return j
}

func main() {
	callback := func(str string) func() {
		return func() {
			fmt.Printf("Alarm! %s 0.8\n", str)
		}
	}
	AddAlert(SystemLoad(OneMinPeriod).Above(0.8).Run(callback("Above")))
	AddAlert(SystemLoad(OneMinPeriod).Below(0.8).Run(callback("Below")))
	StartAlerts()
	select {}
}
