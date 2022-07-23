package Scheduler

import "testing"

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
