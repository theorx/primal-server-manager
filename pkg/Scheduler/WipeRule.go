package Scheduler

import "time"

//How would you define interval of 2 weeks or 3 weeks?
//Let's say wipe every monday after 3 weeks?
type WipeRule struct {
	Name           string
	Days           []time.Weekday
	Hour           int   //Value between 0-24
	Minute         int   //Value between 0-59
	FullWipe       bool  //False = Map wipe
	WipeOnForced   bool  //If WipeOnForced = true, then the rule does not apply at all, and is only triggered when force wipe is detected on that day
	StartTimestamp int64 //Rule is only applied if the startTimestamp is > current unix time
	EndTimestamp   int64 //To make sure, that the wipeRule can also expire. 0 = it will never expire
	//MinDaysSinceLastTrigger can be used to implement rules with longer than 1 week frequency
	//For example if Days only has [1] = Monday, and MinDaysSinceLastTrigger = 13, then
	//it should trigger every 2 weeks. TODO: Create unit tests for that as well
	MinDaysSinceLastTrigger int //Minimum number of days since last trigger.
}

/**
apply It will determine whether the rule can be applied given the timestamp and when it was lastApplied
timestamp - lastApplied has to be greater than MinDaysSinceLastTrigger is, otherwise not going to apply.
Also, if the MinDaysSinceLastTrigger is equal to 0, then lastApplied is totally ignored, but still updated.
*/
func (w *WipeRule) apply(timestamp int64, lastApplied int64) bool {

	//Determine whether the rule is active. If the EndTimestamp is 0, then the rule never expires
	if w.StartTimestamp < timestamp || (w.EndTimestamp != 0 && w.EndTimestamp > timestamp) {
		return false
	}

	//First determine if lastApplied + minDaysSinceLastTrigger * 86400 is greater than timestamp
	//if so, then return false
	if lastApplied+(int64(w.MinDaysSinceLastTrigger)*int64(time.Hour)*24) >= timestamp {
		return false
	}

	t := time.Unix(timestamp, 0)

	t.Day()

	//Get day number from timestamp
	//Fetch hour from timestamp
	//Fetch minute from

	return true
}

func (w *WipeRule) isForcedUpdate(timestamp int64) bool {
	//Determine whether this is the first thursday of the month

	return true
}
