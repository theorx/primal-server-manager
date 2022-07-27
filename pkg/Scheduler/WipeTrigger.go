package Scheduler

import "fmt"

type WipeTrigger struct {
	Name          string
	MapWipe       bool
	BlueprintWipe bool
	ForcedUpdate  bool
	Timestamp     int64
	LastTrigger   int64
}

func (w WipeTrigger) String() string {
	return fmt.Sprintf("Name: %s, MapWipe: %v, BPWipe: %v, ForcedUpdate: %v, Timestamp: %d, LastTrigger: %d", w.Name, w.MapWipe, w.BlueprintWipe, w.ForcedUpdate, w.Timestamp, w.LastTrigger)
}
