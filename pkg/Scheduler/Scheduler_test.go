package Scheduler

import (
	"log"
	"reflect"
	"testing"
	"time"
)

var fullWeek = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

type triggerTimeMemory struct {
	data map[string]int64
}

func (t *triggerTimeMemory) Store(m map[string]int64) error {
	panic("implement me")
}

func (t *triggerTimeMemory) Load() (map[string]int64, error) {
	panic("implement me")
}

func (t *triggerTimeMemory) Get(s string) int64 {
	if _, ok := t.data[s]; ok {
		return t.data[s]
	}

	return 0
}

func (t *triggerTimeMemory) Update(s string, i int64) {
	t.data[s] = i
}

type ruleRegistryMemory struct {
	data []WipeRule
}

func (r *ruleRegistryMemory) Store(rules []WipeRule) error {
	panic("implement me")
}

func (r *ruleRegistryMemory) Load() ([]WipeRule, error) {
	panic("implement me")
}

func (r *ruleRegistryMemory) Update(rule WipeRule) error {
	panic("implement me")
}

func (r *ruleRegistryMemory) Insert(rule WipeRule) error {
	r.data = append(r.data, rule)
	return nil
}

func (r *ruleRegistryMemory) List() []WipeRule {
	return r.data
}

type triggerLogMemory struct {
}

func (t *triggerLogMemory) Log(trigger *WipeTrigger) {

}

func (t *triggerLogMemory) Get(start int64, end int64, limit int64) []*WipeTrigger {

	return nil
}

func NewTestScheduler() *Scheduler {

	return NewScheduler(&triggerTimeMemory{data: make(map[string]int64)}, &triggerLogMemory{}, &ruleRegistryMemory{data: make([]WipeRule, 0)})
}

func TestNewScheduler(t *testing.T) {

	want := &Scheduler{}
	got := NewTestScheduler()

	if got == nil {
		t.Errorf("failed to create new scheduler, got: %v, wanted %v", got, want)
	}
}

func TestRegisterRegistersRules(t *testing.T) {
	s := NewTestScheduler()

	input := []WipeRule{
		{Name: "First"},
		{Name: "Second"},
		{Name: "Third"},
	}

	for _, r := range input {
		if err := s.Register(r); err != nil {
			t.Errorf("Failed to register, err: %v", err)
		}
	}

	if got := s.Rules(); !reflect.DeepEqual(input, got) {
		t.Errorf("Output %v not equal to %v", got, input)
	}
}

func TestTriggeringUpdatesLastTriggeredCorrectly(t *testing.T) {
	s := NewTestScheduler()
	if err := s.Register(WipeRule{
		Name:   "test_last_triggered",
		Days:   fullWeek,
		Hour:   14,
		Minute: 30,
	}); err != nil {
		t.Fatal("Failed to register wipe rule")
	}

	firstTime := time.Date(2022, 1, 17, 14, 30, 0, 0, time.Local)
	secondTime := time.Date(2022, 1, 18, 14, 30, 0, 0, time.Local)

	triggers := s.Schedule(firstTime.Unix())
	log.Println(triggers)
	if len(triggers) != 1 {
		t.Fatalf("first trigger has failed, expected len of 1, but got: %d", len(triggers))
	}

	trigger1 := triggers[0]

	if trigger1.LastTrigger != 0 {
		t.Errorf("trigger of firstTime's LastTrigger timestamp is not 0, got %d", trigger1.LastTrigger)
	}

	triggers = s.Schedule(secondTime.Unix())
	log.Println(triggers)
	if len(triggers) != 1 {
		t.Fatalf("second trigger has failed, expected len of 1, but got: %d", len(triggers))
	}

	trigger2 := triggers[0]

	if trigger2.LastTrigger != firstTime.Unix() {
		t.Errorf("trigger of firstTime's LastTrigger timestamp is not 0, got %d", trigger2.LastTrigger)
	}
}
