# golarm
[![GoDoc](https://godoc.org/github.com/msempere/golarm?status.svg)](https://godoc.org/github.com/msempere/golarm)
[![Build Status](https://travis-ci.org/msempere/golarm.svg?branch=master)](https://travis-ci.org/msempere/golarm)
[![Coverage Status](https://coveralls.io/repos/msempere/golarm/badge.svg?branch=master)](https://coveralls.io/r/msempere/golarm?branch=master)

Fire alarms with system events

## Usage
```go
golarm.AddAlarm(golarm.SystemLoad(golarm.OneMinPeriod).Above(0.8).Run(func() {
		fmt.Println("System load >0.8 !!")
		smtp.SendMail(smtpHost, emailConf.Port, "System load >0.8 !!")
	}))
```

![Usage example](http://i.imgur.com/FybUkVg.gif)

## Options
 - SystemLoad
 
 ```go
// checks if the system load is lower or equal to 0.5
golarm.AddAlarm(golarm.SystemLoad(golarm.OneMinPeriod).AboveEqual(0.5).Run(func() {
		fmt.Println("System load >=0.5 !!")
	}))
```
 - SystemUptime
 
  ```go
// checks if the system has been running for less than 1 minute
golarm.AddAlarm(golarm.SystemUptime().Below(1).Run(func() {
		fmt.Println("System just started !!")
	}))
```
 - SystemMemory / SystemSwap [Free, Used]
 
 ```go
// checks if used memory is higher that 90%
golarm.AddAlarm(golarm.SystemMemory().Used().Above(90).Percent().Run(func() {
		fmt.Println("Used system memory > 90% !!")
	}))
```

 ```go
// checks if free memory is lower that 500MB
golarm.AddAlarm(golarm.SystemMemory().Free().BelowEqual(500).Run(func() {
		fmt.Println("Free memory <= 500MB !!")
	}))
```
 - SystemProc [Status, RunningTime, Used (Memory)]

  ```go
// checks if the process 72332 has changed to zombie status
golarm.AddAlarm(golarm.SystemProc(72332).Status(golarm.Zombie).Run(func() {
		fmt.Println("Our process with PID 72332 became Zombie !!")
	}))
```

  ```go
// checks if the process 72332 has been running for more than 20 minutes
golarm.AddAlarm(golarm.SystemProc(72332).RunningTime().Above(20).Run(func() {
		fmt.Println("Our process with PID 72332 exceeded 20 minutes running !!")
	}))
```

## License
Distributed under MIT license. See `LICENSE` for more information.
