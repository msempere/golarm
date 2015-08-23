# golarm
[![Build Status](https://travis-ci.org/msempere/golarm.svg?branch=master)](https://travis-ci.org/msempere/golarm)
[![Coverage Status](https://coveralls.io/repos/msempere/golarm/badge.svg?branch=master)](https://coveralls.io/r/msempere/golarm?branch=master)

General purpose system event callbacks

## Usage
```go
golarm.AddAlarm(golarm.SystemLoad(golarm.OneMinPeriod).Above(0.8).Run(func() {
		fmt.Println("System load >0.8 !!")
		smtp.SendMail(smtpHost, emailConf.Port, "System load >0.8 !!")
	}))
```

## Options
 - SystemLoad
 
 ```go
golarm.AddAlarm(golarm.SystemLoad(golarm.OneMinPeriod).AboveEqual(0.5).Run(func() {
		fmt.Println("System load >=0.5 !!")
	}))
```
 - SystemUptime
 
  ```go
golarm.AddAlarm(golarm.SystemUptime().Below(1).Run(func() {
		fmt.Println("System just started !!")
	}))
```
 - SystemMemory / SystemSwap [Free, Used]
 
 ```go
golarm.AddAlarm(golarm.SystemMemory().Used().Above(90).Percent().Run(func() {
		fmt.Println("Used system memory > 90% !!")
	}))
```
 - SystemProc [Status, RunningTime, Used (Memory)]

  ```go
golarm.AddAlarm(golarm.SystemProc(72332).Status(golarm.Zombie).Run(func() {
		fmt.Println("Our process with PID 72332 became Zombie !!")
	}))
```

## License
Distributed under MIT license. See `LICENSE` for more information.
