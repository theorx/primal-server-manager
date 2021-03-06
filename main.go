package main

import (
	"log"
	"primal-server-manager/pkg/Scheduler"
	"time"
)

func main() {

	log.Println("Hello, World!")

	//Design a scheduling rule engine

	//Design event sourcing system which would be able to determine 1) oxide update, 2) whether it was a facepunch monthly update wipe or not

	//Have a mechanism to define "plugins"

	//Plugins are c# sources that can have files associated with them
	// 		- The file associations are of two kind:

	r := Scheduler.NewScheduler(nil, nil, nil)

	r.Register(Scheduler.WipeRule{
		Name:                    "Example rule",
		Days:                    []time.Weekday{time.Monday},
		Hour:                    0,
		Minute:                  0,
		WipeOnForced:            false,
		StartTimestamp:          0,
		EndTimestamp:            0,
		MinDaysSinceLastTrigger: 0,
	})

	for _, t := range r.Schedule(0) {
		log.Println("Trigger:", t)
	}

	//use sqlite to store the log?

}
