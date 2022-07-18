package Scheduler

/*
Scheduler holds all the wipe rules
*/
type Scheduler struct {
	Rules []WipeRule
	//Todo: TriggerTimes need to be persisted as well
	triggerTimes map[string]int64
	triggerLog   *TriggerLog
}

/*
NewScheduler Constructor for SchedulerRegistry
*/
func NewScheduler() *Scheduler {
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
