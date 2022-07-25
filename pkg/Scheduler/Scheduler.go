package Scheduler

/*
Scheduler holds all the wipe rules
*/
type Scheduler struct {
	registry    RuleRegistryDAO
	triggerTime TriggerTimeDAO
	triggerLog  TriggerLogDAO
}

/*
NewScheduler Constructor for SchedulerRegistry
*/
func NewScheduler(triggerTime TriggerTimeDAO, triggerLog TriggerLogDAO, registry RuleRegistryDAO) *Scheduler {
	return &Scheduler{
		registry:    registry,
		triggerTime: triggerTime,
		triggerLog:  triggerLog,
	}
}

/*
Schedule This schedule function will need to be called once per minute, as the granularity will be no better than that
*/
func (s *Scheduler) Schedule(input int64) []*WipeTrigger {
	triggers := make([]*WipeTrigger, 0)

	for _, rule := range s.registry.List() {
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
	if len(wr.Name) == 0 {
		return nil
	}

	if !wr.apply(timestamp, s.triggerTime.Get(wr.Name)) {
		return nil
	}

	//return value
	trigger := &WipeTrigger{
		Name:        wr.Name,
		Timestamp:   timestamp,
		FullWipe:    wr.FullWipe,
		LastTrigger: s.triggerTime.Get(wr.Name),
	}

	//update timestamp
	s.triggerTime.Update(wr.Name, timestamp)
	s.triggerLog.Log(trigger)

	return trigger

}

func (s *Scheduler) Register(rule WipeRule) error {

	return s.registry.Insert(rule)
}

func (s *Scheduler) Rules() []WipeRule {
	return s.registry.List()
}

type TriggerTimeDAO interface {
	Store(map[string]int64) error
	Load() (map[string]int64, error)
	Get(string) int64
	Update(string, int64)
}

type TriggerLogDAO interface {
	Log(trigger *WipeTrigger)
	Get(start int64, end int64, limit int64) []*WipeTrigger
}

type RuleRegistryDAO interface {
	Store([]WipeRule) error
	Load() ([]WipeRule, error)
	Update(WipeRule) error
	Insert(WipeRule) error
	List() []WipeRule
}
