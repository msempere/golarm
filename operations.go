package golarm

func setMetric(a *Alarm, v float64, m metric) {
	a.value = value{value: v, percentage: false}
	a.stats.metric = m
}

func setComparison(a *Alarm, v float64, c comparison) {
	a.value = value{value: v, percentage: false}
	a.comparison = c
}

func isMetricCorrect(a *Alarm, v float64, m metric) bool {
	if a.Err == nil {
		switch m {
		case free_:
			if a.jobType == memoryAlarm || a.jobType == swapAlarm {
				return true
			}
		case used_:
			if a.jobType == memoryAlarm || a.jobType == swapAlarm || a.jobType == procAlarm {
				return true
			}
		case time_:
			if a.jobType == uptimeAlarm || a.jobType == procAlarm {
				return true
			}
		case status:
			if a.jobType != alertTypeNotDefined && a.jobType == procAlarm {
				return true
			}
		}
		a.Err = ErrIncorrectTypeForMetric
	}
	return false
}

func isComparisonCorrect(a *Alarm, v float64, c comparison) bool {
	if a.Err == nil {
		if a.comparison == comparisonNotDefined {
			switch c {
			case above, below, equal, belowEqual, aboveEqual:
				if a.jobType != alertTypeNotDefined && a.stats.metric != status {
					return true
				}
				a.Err = ErrIncorrectTypeForComparison

			}
		} else {
			a.Err = ErrMultipleComparisonDefined
		}
	}
	return false
}

// Free allows to specify that the created alarm will use the free memory as main metric
func (j *Alarm) Free() *Alarm {
	if isMetricCorrect(j, notSet, free_) {
		setMetric(j, notSet, free_)
	}
	return j
}

// Used allows to specify that the created alarm will use the used memory as main metric
func (j *Alarm) Used() *Alarm {
	if isMetricCorrect(j, notSet, used_) {
		setMetric(j, notSet, used_)
	}
	return j
}

// RunningTime gets the time a process has been running
func (j *Alarm) RunningTime() *Alarm {
	if isMetricCorrect(j, notSet, time_) {
		setMetric(j, notSet, time_)
	}
	return j
}

// Status allows to specify the state for a given process
func (j *Alarm) Status(s state) *Alarm {
	if isMetricCorrect(j, notSet, status) {
		setMetric(j, notSet, status)
		(*j).stats.proc.state = s
	}
	return j
}

// Above compares if the specified alarm is greater than the number set
func (j *Alarm) Above(v float64) *Alarm {
	if isComparisonCorrect(j, v, above) {
		setComparison(j, v, above)
	}
	return j
}

// AboveEqual compares if the specified alarm is greater or equal than the number set
func (j *Alarm) AboveEqual(v float64) *Alarm {
	if isComparisonCorrect(j, v, aboveEqual) {
		setComparison(j, v, aboveEqual)
	}
	return j
}

// BelowEqual compares if the specified alarm is lower or equal than the number set
func (j *Alarm) BelowEqual(v float64) *Alarm {
	if isComparisonCorrect(j, v, belowEqual) {
		setComparison(j, v, belowEqual)
	}
	return j
}

// Equal compares if the specified alarm is equal than the number set
func (j *Alarm) Equal(v float64) *Alarm {
	if isComparisonCorrect(j, v, equal) {
		setComparison(j, v, equal)
	}
	return j
}

// Below compares if the specified alarm is lower than the number set
func (j *Alarm) Below(v float64) *Alarm {
	if isComparisonCorrect(j, v, below) {
		setComparison(j, v, below)
	}
	return j
}

func parsePercentage(percent float64) (float64, error) {
	if percent > 100 || percent < 0 {
		return 0.0, ErrIncorrectValuesWithPercentage
	}
	return percent, nil
}

// Percent allows using the value specified as a percentage
func (j *Alarm) Percent() *Alarm {
	if j.Err == nil {
		if j.value.value == notSet {
			(*j).Err = ErrExpectedNumWhenPercentage
			return j
		}
		if j.jobType == uptimeAlarm || j.stats.metric == status {
			(*j).Err = ErrIncorrectTypeForPercentage
			return j
		}

		val, err := parsePercentage(j.value.value)

		if err != nil {
			(*j).Err = err
		} else {
			(*j).value = value{value: val, percentage: true}
		}
	}
	return j
}
