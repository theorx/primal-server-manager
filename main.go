package main

import "log"

func main() {

	log.Println("Hello, World!")

	//Design a scheduling rule engine

	//Design event sourcing system which would be able to determine 1) oxide update, 2) whether it was a facepunch monthly update wipe or not

	//Have a mechanism to define "plugins"

	//Plugins are c# sources that can have files associated with them
	// 		- The file associations are of two kind:
}

const (
	Monday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

/*
SchedulerRegistry holds all the wipe rules
*/
type SchedulerRegistry struct {
	Rules []WipeRule
}

//How would you define interval of 2 weeks or 3 weeks?
//Let's say wipe every monday after 3 weeks?
type WipeRule struct {
	Name                    string
	Days                    []int
	Hour                    int   //Value between 0-24
	Minute                  int   //Value between 0-59
	FullWipe                bool  //False = Map wipe
	WipeOnForced            bool  //If WipeOnForced = true, then the rule does not apply at all, and is only triggered when force wipe is detected on that day
	StartTimestamp          int64 //Rule is only applied if the startTimestamp is > current unix time
	EndTimestamp            int64 //To make sure, that the wipeRule can also expire. 0 = it will never expire
	//MinDaysSinceLastTrigger can be used to implement rules with longer than 1 week frequency
	//For example if Days only has [1] = Monday, and MinDaysSinceLastTrigger = 13, then
	//it should trigger every 2 weeks. TODO: Create unit tests for that as well
	MinDaysSinceLastTrigger int   //Minimum number of days since last trigger.
}

func isForceUpdateToday() {
	//Determine whether this is the first thursday of the month
}
