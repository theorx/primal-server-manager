package Scheduler

/*
Scheduler holds all the wipe rules
*/
type Scheduler struct {
	rules []WipeRule
	//Todo: TriggerTimes need to be persisted as well
	triggerTimes map[string]int64
	triggerLog   *TriggerLog
}

/*
NewScheduler Constructor for SchedulerRegistry
*/
func NewScheduler() *Scheduler {
	return &Scheduler{
		rules:        make([]WipeRule, 0),
		triggerTimes: make(map[string]int64),
		triggerLog:   &TriggerLog{},
	}
}

/*
Schedule This schedule function will need to be called once per minute, as the granularity will be no better than that
*/
func (s *Scheduler) Schedule(input int64) []*WipeTrigger {
	triggers := make([]*WipeTrigger, 0)

	for _, rule := range s.rules {
		if trigger := s.tryApply(&rule, input); trigger != nil {
			triggers = append(triggers, trigger)
		}
	}

	return triggers
}

/*
NextTrigger predicts the timestamp for the next future trigger for a given rule
Input is the int64 unix timestamp of the current time
*/
func (s *Scheduler) NextTrigger(input int64, rule *WipeRule) int64 {

	return 0
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

/*
getTriggerTime is called to determine when the specific trigger was triggered successfully
*/
func (s *Scheduler) getTriggerTime(name string) int64 {
	if _, ok := s.triggerTimes[name]; !ok {
		s.triggerTimes[name] = 0
	}

	return s.triggerTimes[name]
}

/*
Updates triggered time
*/
func (s *Scheduler) updateTriggerTime(name string, value int64) {
	s.triggerTimes[name] = value
}

/*
Registers new rule, has to be persisted TODO: Design persistence layer separately from scheduler?
*/
func (s *Scheduler) Register(rule WipeRule) error {
	//Check whether the name is unique?
	s.rules = append(s.rules, rule)

	return nil
}

/*
Returns the rules loaded to memory
*/
func (s *Scheduler) Rules() []WipeRule {
	return s.rules
}
