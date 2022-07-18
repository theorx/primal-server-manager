package Scheduler

import "fmt"

type WipeTrigger struct {
	Name        string
	FullWipe    bool
	Timestamp   int64
	LastTrigger int64
}

func (w WipeTrigger) String() string {
	return fmt.Sprintf("Name: %s, FullWipe: %v, Timestamp: %d, LastTrigger: %d", w.Name, w.FullWipe, w.Timestamp, w.LastTrigger)
}
