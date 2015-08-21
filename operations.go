package golarm

// Free can only work with:
//	- MemoryAlarm
//	- SwapAlarm
func (j *alarm) Free() *alarm {
	if j.err == nil {
		if j.jobType == memoryAlarm || j.jobType == swapAlarm {
			(*j).value = value{value: not_set, percentage: false}
			(*j).stats.metric = free_
		} else {
			(*j).err = ErrIncorrectTypeForFree
		}
	}
	return j
}

// Used can only work with:
//	- MemoryAlarm
//	- SwapAlarm
//	- Procalarm -> used
func (j *alarm) Used() *alarm {
	if j.err == nil {
		if j.jobType == memoryAlarm || j.jobType == swapAlarm || j.jobType == procAlarm {
			(*j).value = value{value: not_set, percentage: false}
			(*j).stats.metric = used_
		} else {
			(*j).err = ErrIncorrectTypeForUsed
		}
	}
	return j
}

// Used can only work with:
//	- Uptime
//	- Procalarm -> time
func (j *alarm) RunningTime() *alarm {
	if j.err == nil {
		if j.jobType == uptimeAlarm || j.jobType == procAlarm {
			(*j).value = value{value: not_set, percentage: false}
			(*j).stats.metric = time
		} else {
			(*j).err = ErrIncorrectTypeForTime
		}
	}
	return j
}

// Used can only work with:
//	- Procalarm -> status
func (j *alarm) Status(s state) *alarm {
	if j.err == nil {
		if j.jobType != alertTypeNotDefined && j.jobType == procAlarm {
			(*j).value = value{value: not_set, percentage: false}
			(*j).stats.proc.state = s
			(*j).stats.metric = status
		} else {
			(*j).err = ErrIncorrectTypeForStatus
		}
	}
	return j
}

// Time = Minutes
// Memory = MB
func (j *alarm) Above(v float64) *alarm {
	if j.err == nil {
		if j.jobType != alertTypeNotDefined && j.stats.metric != status {
			if j.comparison == comparisonNotDefined {
				(*j).value = value{value: v, percentage: false}
				(*j).comparison = above
			} else {
				(*j).err = ErrMultipleComparisonDefined
			}
		} else {
			(*j).err = ErrIncorrectTypeForAbove
		}
	}
	return j
}

// Time = Minutes
// Memory = MB
func (j *alarm) Below(v float64) *alarm {
	if j.err == nil {
		if j.jobType != alertTypeNotDefined && j.stats.metric != status {
			if j.comparison == comparisonNotDefined {
				(*j).value = value{value: v, percentage: false}
				(*j).comparison = below
			} else {
				(*j).err = ErrMultipleComparisonDefined
			}
		} else {
			(*j).err = ErrIncorrectTypeForBelow
		}
	}
	return j
}

func parsePercentage(percent float64) (float64, error) {
	if percent > 100 || percent < 0 {
		return 0.0, ErrIncorrectValuesWithPercentage
	}
	return percent, nil
}

func (j *alarm) Percent() *alarm {
	if j.err == nil {
		if j.value.value == not_set {
			(*j).err = ErrExpectedNumWhenPercentage
			return j
		}
		if j.jobType == uptimeAlarm || j.stats.metric == status {
			(*j).err = ErrIncorrectTypeForPercentage
			return j
		}

		val, err := parsePercentage(j.value.value)

		if err != nil {
			(*j).err = err
		} else {
			(*j).value = value{value: val, percentage: true}
		}
	}
	return j
}
