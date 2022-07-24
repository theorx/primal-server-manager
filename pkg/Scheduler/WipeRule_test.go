package Scheduler

import (
	"testing"
	"time"
)

func TestWipeRuleApplyInActiveTimeRange(t *testing.T) {
	tt := []struct {
		Start     int64
		End       int64
		Timestamp int64
	}{
		{Start: 100, Timestamp: 100},
		{Start: 10, Timestamp: 9},
		{Start: 10, Timestamp: 10},
		{Start: 0, End: 1, Timestamp: 1},
		{Start: 0, End: 100, Timestamp: 100},
		{Start: 99, End: 100, Timestamp: 99},
		{Start: 99, End: 100, Timestamp: 100},
	}

	for _, c := range tt {
		//All the cases are expected to return false, otherwise fail
		if (&WipeRule{StartTimestamp: c.Start, EndTimestamp: c.End}).apply(c.Timestamp, 0) {
			t.Errorf("Expected apply() to return false, got true. StartTimestamp: %d, Timestamp: %d", c.Start, c.Timestamp)
		}
	}
}

func TestWipeRuleLastAppliedMinDaysSinceLastTriggerRules(t *testing.T) {
	startTime := int64(1)
	const day = 86400

	tt := []struct {
		MinDaysSinceLastTrigger int
		Timestamp               int64
		LastApplied             int64
	}{
		{
			Timestamp:   100,
			LastApplied: 2134,
		},
		{
			Timestamp:   100,
			LastApplied: 100,
		},
		{
			MinDaysSinceLastTrigger: 1,
			Timestamp:               100,
			LastApplied:             0,
		},
		{
			MinDaysSinceLastTrigger: 1,
			Timestamp:               day,
			LastApplied:             0,
		},
		{
			MinDaysSinceLastTrigger: 1,
			Timestamp:               day * 2,
			LastApplied:             day,
		},
	}

	for _, c := range tt {
		if (&WipeRule{StartTimestamp: startTime, MinDaysSinceLastTrigger: c.MinDaysSinceLastTrigger}).apply(c.Timestamp, c.LastApplied) {
			t.Errorf("Expected apply() to return false when timestamp is less than lastApplied + minDaysSinceLastTrigger in seconds")
		}
	}
}

func TestWipeRuleApplyMatchesDays(t *testing.T) {
	instance := &WipeRule{}

	tt := []struct {
		Days      []time.Weekday
		Timestamp int64
		Result    bool
	}{
		{Days: nil, Timestamp: 0, Result: false},
		{Days: []time.Weekday{}, Timestamp: 0, Result: false},
		{Days: []time.Weekday{time.Sunday}, Timestamp: 1658674601, Result: true},
		{Days: []time.Weekday{time.Monday}, Timestamp: 1658761001, Result: true},
		{Days: []time.Weekday{time.Tuesday}, Timestamp: 1658847401, Result: true},
		{Days: []time.Weekday{time.Wednesday}, Timestamp: 1658933801, Result: true},
		{Days: []time.Weekday{time.Thursday}, Timestamp: 1659020201, Result: true},
		{Days: []time.Weekday{time.Friday}, Timestamp: 1659106601, Result: true},
		{Days: []time.Weekday{time.Saturday}, Timestamp: 1659193001, Result: true},
		{Days: []time.Weekday{time.Friday}, Timestamp: 1659193001, Result: false},
		{Days: []time.Weekday{time.Sunday}, Timestamp: 1661093801, Result: true},
		{Days: []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}, Timestamp: 4, Result: true},
		{Days: []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Thursday, time.Friday, time.Saturday}, Timestamp: 1661353001, Result: false},
	}

	for _, c := range tt {
		instance.Days = c.Days
		if got := instance.matchWeekday(c.Timestamp); got != c.Result {
			t.Errorf("WipeRule.matchWeekday failed, expected: %v, Got: %v", c.Result, got)
		}
	}
}
