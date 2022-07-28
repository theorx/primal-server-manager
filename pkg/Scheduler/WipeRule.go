package Scheduler

import (
	"time"
)

const ForceWipeHourUtc = 18 //18:00 UTC rust force wipe

type WipeRule struct {
	Name                    string
	Server                  string //Uuidv4 of server?
	Days                    []time.Weekday
	Hour                    int
	Minute                  int
	MapWipe                 bool
	BlueprintWipe           bool
	WipeOnForced            bool //If WipeOnForced = true, then the rule does not apply at all, and is only triggered when force wipe is detected on that day
	StartTimestamp          int64
	EndTimestamp            int64
	MinDaysSinceLastTrigger int
}

/**
apply It will determine whether the rule can be applied given the timestamp and when it was lastApplied
timestamp - lastApplied has to be greater than MinDaysSinceLastTrigger is, otherwise not going to apply.
Also, if the MinDaysSinceLastTrigger is equal to 0, then lastApplied is totally ignored, but still updated.
*/
func (w *WipeRule) apply(timestamp int64, lastApplied int64) bool {

	//Determine whether the rule is active. If the EndTimestamp is 0, then the rule never expires
	if w.StartTimestamp >= timestamp || (w.EndTimestamp != 0 && w.EndTimestamp <= timestamp) {
		return false
	}

	if w.WipeOnForced && w.isForcedUpdate(timestamp) {
		if time.Unix(timestamp, 0).UTC().Hour() == ForceWipeHourUtc && time.Unix(timestamp, 0).Minute() == 0 {
			return true
		}
		return false
	}

	if lastApplied+(int64(w.MinDaysSinceLastTrigger)*86400) >= timestamp {
		return false
	}

	if !w.matchWeekday(timestamp) {
		return false
	}

	return w.matchHourAndMinute(timestamp)
}

/*
matchWeekday matches the weekday of the given timestamp based on the days configured for the rule
If the currentDay is not present within the rule, false is returned.
*/
func (w *WipeRule) matchWeekday(timestamp int64) bool {

	currentDay := time.Unix(timestamp, 0).Weekday()
	for _, day := range w.Days {
		if currentDay == day {
			return true
		}
	}

	return false
}

/**
matchHourAndMinute
*/
func (w *WipeRule) matchHourAndMinute(timestamp int64) bool {
	t := time.Unix(timestamp, 0).UTC()

	if t.Hour() == w.Hour && t.Minute() == w.Minute {
		return true
	}

	return false
}

/**
isForcedUpdate checks whether it's the first thursday of the month as per Facepunch's update policy
*/
func (w *WipeRule) isForcedUpdate(timestamp int64) bool {
	t := time.Unix(timestamp, 0)

	if t.Weekday() == time.Thursday && t.Day() <= 7 {
		return true
	}

	return false
}
