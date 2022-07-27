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
		if (&WipeRule{StartTimestamp: c.Start, EndTimestamp: c.End}).apply(c.Timestamp, 0) {
			t.Errorf("Expected apply() to return false, got true. StartTimestamp: %d, Timestamp: %d", c.Start, c.Timestamp)
		}
	}
}

func TestWipeRuleLastAppliedMinDaysSinceLastTriggerRules(t *testing.T) {
	const day = 86400

	tt := []struct {
		MinDaysSinceLastTrigger int
		Timestamp               int64
		LastApplied             int64
	}{
		{Timestamp: 100, LastApplied: 2134},
		{Timestamp: 100, LastApplied: 100},
		{MinDaysSinceLastTrigger: 1, Timestamp: 100, LastApplied: 0},
		{MinDaysSinceLastTrigger: 1, Timestamp: day, LastApplied: 0},
		{MinDaysSinceLastTrigger: 1, Timestamp: day * 2, LastApplied: day},
	}

	for _, c := range tt {
		if (&WipeRule{StartTimestamp: 1, MinDaysSinceLastTrigger: c.MinDaysSinceLastTrigger}).apply(c.Timestamp, c.LastApplied) {
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

func TestWipeRuleMatchHourAndMinute(t *testing.T) {
	tt := []struct {
		Hour   int
		Minute int
		Result bool
	}{
		{Hour: 0, Minute: 0, Result: true},
		{Hour: 24, Minute: 0, Result: false},
		{Hour: 1, Minute: 1, Result: true},
		{Hour: 12, Minute: 0, Result: true},
		{Hour: 23, Minute: 59, Result: true},
		{Hour: 25, Minute: 159, Result: false},
		{Hour: 13, Minute: 30, Result: true},
		{Hour: 12, Minute: 00, Result: true},
		{Hour: 11, Minute: 59, Result: true},
		{Hour: 11, Minute: 60, Result: false},
		{Hour: -22, Minute: 60, Result: false},
	}

	for _, c := range tt {
		if (&WipeRule{
			Hour:   c.Hour,
			Minute: c.Minute,
		}).matchHourAndMinute(
			time.Date(2000, 1, 1, c.Hour, c.Minute, 0, 0, time.UTC).Unix(),
		) != c.Result {
			t.Errorf("matchHourAndMinute failed with input: H: %d M: %d", c.Hour, c.Minute)
		}
	}
}

func TestWipeRuleIsForcedWipe(t *testing.T) {
	tt := []struct {
		Month  int
		Day    int
		Result bool
	}{
		{Month: 7, Day: 1, Result: false},
		{Month: 7, Day: 6, Result: false},
		{Month: 7, Day: 7, Result: true},
		{Month: 7, Day: 8, Result: false},
		{Month: 7, Day: 14, Result: false},
		{Month: 7, Day: 7, Result: true},
		{Month: 8, Day: 2, Result: false},
		{Month: 1, Day: 1, Result: false},
		{Month: 1, Day: 2, Result: false},
		{Month: 1, Day: 3, Result: false},
		{Month: 1, Day: 4, Result: false},
		{Month: 1, Day: 5, Result: false},
		{Month: 1, Day: 6, Result: true},
		{Month: 1, Day: 7, Result: false},
		{Month: 1, Day: 8, Result: false},
		{Month: 1, Day: 8, Result: false},
		{Month: 6, Day: 2, Result: true},
		{Month: 6, Day: 9, Result: false},
		{Month: 6, Day: 16, Result: false},
		{Month: 6, Day: 23, Result: false},
	}

	for _, c := range tt {
		if got := (&WipeRule{}).isForcedUpdate(time.Date(2022, time.Month(c.Month), c.Day, 12, 0, 0, 0, time.UTC).Unix()); got != c.Result {
			t.Errorf("isForcedWipe() failed, expected: %v, got: %v, case: %v", c.Result, got, c)
		}
	}
}

func TestWipeRuleApplyHandlesForcedUpdateOnTime(t *testing.T) {

	tt := []struct {
		Timestamp int64
		Result    bool
	}{
		{Timestamp: time.Date(2022, 7, 7, ForceWipeHourUtc, 0, 0, 0, time.UTC).Unix(), Result: true},
		{Timestamp: time.Date(2022, 6, 2, ForceWipeHourUtc, 0, 0, 0, time.UTC).Unix(), Result: true},
		{Timestamp: time.Date(2022, 6, 2, 14, 0, 0, 0, time.UTC).Unix(), Result: false},
		{Timestamp: time.Date(2022, 7, 7, ForceWipeHourUtc, 1, 0, 0, time.UTC).Unix(), Result: false},
		{Timestamp: time.Date(2022, 7, 7, ForceWipeHourUtc-1, 59, 0, 0, time.UTC).Unix(), Result: false},
	}

	for _, c := range tt {
		if (&WipeRule{
			WipeOnForced: true,
		}).apply(c.Timestamp, 0) != c.Result {
			t.Errorf("WipeRule.apply failed handling forced update, case: %v", c)
		}
	}
}

func TestWipeRuleApplyHandlesMatchWeekDay(t *testing.T) {
	tt := []struct {
		Days      []time.Weekday
		Timestamp int64
		Result    bool
	}{
		{
			Days:      []time.Weekday{time.Monday},
			Timestamp: time.Date(2022, 7, 25, 0, 0, 0, 0, time.UTC).Unix(),
			Result:    true,
		},
		{
			Days:      []time.Weekday{time.Tuesday},
			Timestamp: time.Date(2022, 7, 26, 0, 0, 0, 0, time.UTC).Unix(),
			Result:    true,
		},
		{
			Days:      []time.Weekday{time.Tuesday},
			Timestamp: time.Date(2022, 7, 26, 0, 45, 0, 0, time.UTC).Unix(),
			Result:    false,
		},
		{
			Days:      []time.Weekday{time.Tuesday},
			Timestamp: time.Date(2022, 7, 26, 2, 0, 0, 0, time.UTC).Unix(),
			Result:    false,
		},
		{
			Days:      []time.Weekday{time.Wednesday},
			Timestamp: time.Date(2022, 7, 26, 0, 0, 0, 0, time.UTC).Unix(),
			Result:    false,
		},
		{
			Days:      []time.Weekday{},
			Timestamp: time.Date(2022, 7, 4, 0, 0, 0, 0, time.UTC).Unix(),
			Result:    false,
		},
	}

	for _, c := range tt {
		if (&WipeRule{Days: c.Days, WipeOnForced: false}).apply(c.Timestamp, 0) != c.Result {
			t.Errorf("WipeRule.apply() failed to handle matchWeekday(), case: %v", c)
		}
	}
}
