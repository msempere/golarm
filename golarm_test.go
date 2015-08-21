package golarm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddingWrongAlert(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Percent().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrExpectedNumWhenPercentage)
	assert.NotNil(test, a.err)
	err := AddAlarm(a)
	assert.NotNil(test, err)
	assert.Equal(test, a.err, err)
}

func TestWrongPercentage(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Above(101).Percent().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectValuesWithPercentage)
	assert.NotNil(test, a.err)
}

func TestBadChain1(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Percent().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrExpectedNumWhenPercentage)
	assert.NotNil(test, a.err)
}

func TestBadChain2(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrComparisonNotDefined)
	assert.NotNil(test, a.err)
}

func TestComparisonUndefined(test *testing.T) {
	a := SystemMemory().Free().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrComparisonNotDefined)
	assert.NotNil(test, a.err)
}

func TestIncorrecTypeWithFree(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Free().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectTypeForFree)
	assert.NotNil(test, a.err)
}

func TestIncorrecTypeWithUsed(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Used().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectTypeForUsed)
	assert.NotNil(test, a.err)
}

func TestIncorrecTypeWithStatus(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Status(Running).Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectTypeForStatus)
	assert.NotNil(test, a.err)
}

func TestAboveWithStatus(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Status(Running).Above(5).Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectTypeForAbove)
	assert.NotNil(test, a.err)
}

func TestIncorrectPid(test *testing.T) {
	a := SystemProc(9999999999).RunningTime().Above(5).Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrInexistentPid)
	assert.NotNil(test, a.err)
}

func TestBelowWithStatus(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Status(Running).Below(5).Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectTypeForBelow)
	assert.NotNil(test, a.err)
}

func TestMultipleComparisons(test *testing.T) {
	a := SystemLoad(OneMinPeriod).Above(5).Below(1).Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrMultipleComparisonDefined)
	assert.NotNil(test, a.err)
}

func TestIncorrecTypeWithTime(test *testing.T) {
	a := SystemLoad(OneMinPeriod).RunningTime().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrIncorrectTypeForTime)
	assert.NotNil(test, a.err)
}

func TestIncorrecTypeWithPercentaje(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Status(Running).Percent().Run(func() {})
	go Check(a)
	assert.Equal(test, a.err, ErrExpectedNumWhenPercentage)
	assert.NotNil(test, a.err)
}

func TestSystemProc(test *testing.T) {
	a := SystemProc(uint(os.Getpid())).Used().Above(50).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.err, nil)

	a = SystemProc(uint(os.Getpid())).Used().Below(50).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.err, nil)

	a = SystemProc(uint(os.Getpid())).RunningTime().Above(5).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.err, nil)

	a = SystemProc(uint(os.Getpid())).RunningTime().Below(500).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.err, nil)

	a = SystemProc(uint(os.Getpid())).RunningTime().Below(1).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.err, nil)
}

func TestSystemLoad(test *testing.T) {
	a := SystemLoad(FiveMinPeriod).Above(0).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.err, nil)

	a = SystemLoad(OneMinPeriod).Above(5).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.err, nil)

	a = SystemLoad(FifteenMinPeriod).Below(5).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Equal(test, a.err, nil)

	a = SystemLoad(OneMinPeriod).Below(0).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Equal(test, a.err, nil)
}

func TestSystemMemory(test *testing.T) {
	a := SystemMemory().Free().Above(70).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Free().Above(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Free().Below(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Free().Below(10).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Used().Above(10).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Used().Above(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Used().Below(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemMemory().Used().Below(10).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)
}

func TestSystemSwapMemory(test *testing.T) {
	a := SystemSwap().Free().Above(70).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Free().Above(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Free().Below(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Free().Below(10).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Used().Above(10).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Used().Above(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Used().Below(90).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemSwap().Used().Below(10).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)
}

func TestSystemUptime(test *testing.T) {
	a := SystemUptime().Above(70).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, true)
	assert.Nil(test, a.err, nil)

	a = SystemUptime().Below(70).Run(func() {})
	a.SetMetricsManager(&FakeSigar{})
	go Check(a)
	assert.Equal(test, <-a.result, false)
	assert.Nil(test, a.err, nil)
}
