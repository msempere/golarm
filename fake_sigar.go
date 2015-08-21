package golarm

import (
	t "time"

	"github.com/cloudfoundry/gosigar"
)

type FakeSigar struct {
}

func (f *FakeSigar) CollectCpuStats(collectionInterval t.Duration) (<-chan sigar.Cpu, chan<- struct{}) {
	return nil, nil
}

func (f *FakeSigar) GetLoadAverage() (sigar.LoadAverage, error) {
	return sigar.LoadAverage{
		One:     1,
		Five:    1,
		Fifteen: 1,
	}, nil
}

func (f *FakeSigar) GetMem() (sigar.Mem, error) {
	return sigar.Mem{
		Total:      100000000,
		Used:       20000000,
		Free:       80000000,
		ActualFree: 80000000,
		ActualUsed: 20000000,
	}, nil
}

func (f *FakeSigar) GetSwap() (sigar.Swap, error) {
	return sigar.Swap{
		Total: 100000000,
		Used:  20000000,
		Free:  80000000,
	}, nil
}

func (f *FakeSigar) GetFileSystemUsage(string) (sigar.FileSystemUsage, error) {
	return sigar.FileSystemUsage{
		Total:     500,
		Used:      250,
		Free:      250,
		Avail:     250,
		Files:     1,
		FreeFiles: 1,
	}, nil
}

func (f *FakeSigar) GetProcState(pid int) (sigar.ProcState, error) {
	return sigar.ProcState{
		Name:      "fakeProc",
		State:     sigar.RunStateRun,
		Ppid:      500,
		Tty:       69,
		Priority:  1,
		Nice:      0,
		Processor: 1,
	}, nil
}

func (f *FakeSigar) GetProcMem(pid int) (sigar.ProcMem, error) {
	return sigar.ProcMem{
		Size:        100000000,
		Resident:    200000000,
		Share:       0,
		MinorFaults: 0,
		MajorFaults: 0,
		PageFaults:  0,
	}, nil
}

func (f *FakeSigar) GetProcTime(pid int) (sigar.ProcTime, error) {
	return sigar.ProcTime{
		StartTime: 123456,
		User:      123456,
		Sys:       123456,
		Total:     123456,
	}, nil
}

func (f *FakeSigar) GetUpTime() (sigar.Uptime, error) {
	return sigar.Uptime{
		Length: 120,
	}, nil
}
