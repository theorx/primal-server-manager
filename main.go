package main

import (
	"log"
	"primal-server-manager/pkg/Scheduler"
)

func main() {

	log.Println("Hello, World!")

	//Design a scheduling rule engine

	//Design event sourcing system which would be able to determine 1) oxide update, 2) whether it was a facepunch monthly update wipe or not

	//Have a mechanism to define "plugins"

	//Plugins are c# sources that can have files associated with them
	// 		- The file associations are of two kind:

	r := Scheduler.NewScheduler()

	r.Register(Scheduler.WipeRule{
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
