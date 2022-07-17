package main

import (
	"log"
	"time"
)

func main() {

	log.Println("Hello, World!")

	//Design a scheduling rule engine

	//Design event sourcing system which would be able to determine 1) oxide update, 2) whether it was a facepunch monthly update wipe or not

	//Have a mechanism to define "plugins"

	//Plugins are c# sources that can have files associated with them
	// 		- The file associations are of two kind:

	r := &Scheduler{}

	r.Register(WipeRule{
		Name:                    "",
		Days:                    nil,
		Hour:                    0,
		Minute:                  0,
		FullWipe:                false,
		WipeOnForced:            false,
		StartTimestamp:          0,
		EndTimestamp:            0,
		MinDaysSinceLastTrigger: 0,
	})

	//use sqlite to store the log?

}

/*
Scheduler holds all the wipe rules
*/
type Scheduler struct {
	Rules        []WipeRule
	triggerTimes map[string]int64
	triggerLog   *TriggerLog
}

/*
NewSchedulerRegistry Constructor for SchedulerRegistry
*/
func NewSchedulerRegistry() *Scheduler {
	return &Scheduler{
		Rules:        make([]WipeRule, 0),
		triggerTimes: make(map[string]int64),
		triggerLog:   &TriggerLog{},
	}
}

/*
Schedule This schedule function will need to be called once per minute, as the granularity will be no better than that
*/
func (s *Scheduler) Schedule(input int64) []*WipeTrigger {
	triggers := make([]*WipeTrigger, 0)

	for _, rule := range s.Rules {
		if trigger := s.tryApply(&rule, input); trigger != nil {
			triggers = append(triggers, trigger)
		}
	}

	return triggers
}

/**
tryApply tries to apply the rule, if successful, then returns wipeTrigger, otherwise nil
*/
func (s *Scheduler) tryApply(wr *WipeRule, timestamp int64) *WipeTrigger {
	if _, ok := s.triggerTimes[wr.Name]; !ok {
		s.triggerTimes[wr.Name] = 0
	}

	if !wr.apply(timestamp, s.triggerTimes[wr.Name]) {
		return nil
	}

	//update timestamp
	s.triggerTimes[wr.Name] = timestamp
	//store in log

	//return value
	trigger := &WipeTrigger{
		Name:      wr.Name,
		Timestamp: timestamp,
		FullWipe:  wr.FullWipe,
	}

	s.triggerLog.Log(trigger)

	return trigger

}

func (s *Scheduler) Register(rule WipeRule) error {
	//Check whether the name is unique?
	s.Rules = append(s.Rules, rule)

	return nil
}

type WipeTrigger struct {
	Name        string
	FullWipe    bool
	Timestamp   int64
	LastTrigger int64
}

type TriggerLog struct {
	Name        string
	Timestamp   int64
	LastTrigger int64
}

func (t *TriggerLog) Log(wt *WipeTrigger) {
	//Store in the table
}

func (t *TriggerLog) Get(start int64, end int64, limit int64) []*WipeTrigger {

	return nil
}

//How would you define interval of 2 weeks or 3 weeks?
//Let's say wipe every monday after 3 weeks?
type WipeRule struct {
	Name           string
	Days           []int
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
