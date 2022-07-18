package Scheduler

type WipeTrigger struct {
	Name        string
	FullWipe    bool
	Timestamp   int64
	LastTrigger int64
}

