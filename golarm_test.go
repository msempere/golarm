package golarm

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddingWrongAlert(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Percent().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrExpectedNumWhenPercentage)
	assert.NotNil(test, a.Err)
	err := AddAlarm(a)
	assert.NotNil(test, err)
	assert.Equal(test, a.Err, err)
}

func TestWrongPercentage(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Above(101).Percent().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectValuesWithPercentage)
	assert.NotNil(test, a.Err)
}

func TestBadChain1(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Percent().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrExpectedNumWhenPercentage)
	assert.NotNil(test, a.Err)
}

func TestBadChain2(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrComparisonNotDefined)
	assert.NotNil(test, a.Err)
}

func TestComparisonUndefined(test *testing.T) {
	a := SystemMemory().Free().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrComparisonNotDefined)
	assert.NotNil(test, a.Err)
}

func TestIncorrecTypeWithFree(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Free().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectTypeForMetric)
	assert.NotNil(test, a.Err)
}

func TestIncorrecTypeWithUsed(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Used().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectTypeForMetric)
	assert.NotNil(test, a.Err)
}

func TestIncorrecTypeWithStatus(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Status(Running).Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectTypeForMetric)
	assert.NotNil(test, a.Err)
}

func TestAboveWithStatus(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Status(Running).Above(5).Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectTypeForComparison)
	assert.NotNil(test, a.Err)
}

func TestIncorrectPid(test *testing.T) {
	a := SystemProc(9999999999).RunningTime().Above(5).Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrInexistentPid)
	assert.NotNil(test, a.Err)
}

func TestBelowWithStatus(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Status(Running).Below(5).Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectTypeForComparison)
	assert.NotNil(test, a.Err)
}

func TestMultipleComparisons(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Above(5).Below(1).Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrMultipleComparisonDefined)
	assert.NotNil(test, a.Err)
}

func TestIncorrecTypeWithTime(test *testing.T) {
	a := SystemLoad(OneMinPeriod).RunningTime().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrIncorrectTypeForMetric)
	assert.NotNil(test, a.Err)
}

func TestIncorrecTypeWithPercentaje(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Status(Running).Percent().Run(func() {})
	go check(a)
	assert.Equal(test, a.Err, ErrExpectedNumWhenPercentage)
	assert.NotNil(test, a.Err)
}

func TestSystemProc(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Used().Equal(95).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).Used().Below(50).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).Used().Below(50).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).Used().Below(100).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).Status(Running).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).RunningTime().Above(5).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).RunningTime().Below(500).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemProc(uint(os.Getpid())).RunningTime().Below(1).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.Err, nil)
}

func TestSystemLoad(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Above(0).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemLoad(OneMinPeriod).AboveEqual(5).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.Err, nil)

	a = SystemLoad(FifteenMinPeriod).BelowEqual(5).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.Err, nil)

	a = SystemLoad(OneMinPeriod).Below(0).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.Err, nil)

	a = SystemLoad(OneMinPeriod).Below(0).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.Err, nil)
}

func TestSystemMemory(test *testing.T) {
	a := SystemMemory().Free().Above(70).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Free().Above(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Free().Above(90).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Free().Below(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Free().Below(10).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Used().Above(10).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Used().Above(10).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Used().Above(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Used().Below(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemMemory().Used().Below(10).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)
}

func TestSystemSwapMemory(test *testing.T) {
	a := SystemSwap().Free().Above(70).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Free().Above(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Free().Above(90).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Free().Below(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Free().Below(10).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Used().Above(10).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Used().Above(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Used().Below(90).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Used().Below(10).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)

	a = SystemSwap().Used().Below(1).Percent().Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)
}

func TestSystemUptime(test *testing.T) {
	a := SystemUptime().Above(70).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.Err, nil)

	a = SystemUptime().Below(70).Run(func() {})
	a.SetMetricsManager(&fakeSigar{})
	go check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.Err, nil)
}

func TestRealSystemUptime(test *testing.T) {
	value := false
	a := SystemUptime().Above(1).Run(func() { value = true })
	Duration = 1
	err := AddAlarm(a)
	time.Sleep(1250 * time.Millisecond)
	assert.Equal(test, value, true)
	assert.Nil(test, a.Err, nil)
	assert.Nil(test, err, nil)
}
