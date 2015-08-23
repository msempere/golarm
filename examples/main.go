package main

import (
	"fmt"

	"github.com/msempere/golarm"
)

func main() {
	err := golarm.AddAlarm(golarm.SystemLoad(golarm.OneMinPeriod).Above(0.8).Run(func() {
		fmt.Println("System load >0.8 !!")
	}))

	if err != nil {
		fmt.Println(err)
	}

	err = golarm.AddAlarm(golarm.SystemMemory().Below(50).Percent().Run(func() {
		fmt.Println("System memory <50% !!")
	}))

	if err != nil {
		fmt.Println(err)
	}

	err = golarm.AddAlarm(golarm.SystemProc(2453).MemoryUsed().Below(50).Percent().Run(func() {
		fmt.Println("<50% mem used by process 2453!!")
	}))

	if err != nil {
		fmt.Println(err)
	}

	err = golarm.AddAlarm(golarm.SystemProc(2453).RunningTime().Above(45).Run(func() {
		fmt.Println("Process 2453 running for more than 45 minutes")
	}))

	if err != nil {
		fmt.Println(err)
	}

	select {}
}
