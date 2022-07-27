package Scheduler

import (
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

	firstTime := time.Date(2022, 1, 17, 14, 30, 0, 0, time.UTC)
	secondTime := time.Date(2022, 1, 18, 14, 30, 0, 0, time.UTC)

	triggers := s.Schedule(firstTime.Unix())

	if len(triggers) != 1 {
		t.Fatalf("first trigger has failed, expected len of 1, but got: %d", len(triggers))
	}

	trigger1 := triggers[0]

	if trigger1.LastTrigger != 0 {
		t.Errorf("trigger of firstTime's LastTrigger timestamp is not 0, got %d", trigger1.LastTrigger)
	}

	triggers = s.Schedule(secondTime.Unix())

	if len(triggers) != 1 {
		t.Fatalf("second trigger has failed, expected len of 1, but got: %d", len(triggers))
	}

	trigger2 := triggers[0]

	if trigger2.LastTrigger != firstTime.Unix() {
		t.Errorf("trigger of firstTime's LastTrigger timestamp is not 0, got %d", trigger2.LastTrigger)
	}
}

func TestTryApplyReturnsNilWhenNoNameSpecified(t *testing.T) {
	s := NewTestScheduler()
	if got := s.tryApply(&WipeRule{Name: ""}, 0); got != nil {
		t.Errorf("tryApply expected to return nil with empty Name on the wipe rule, got: %v", got)
	}
}

func TestMonthlyForceWipes(t *testing.T) {
	s := NewTestScheduler()

	//First test with a monthly wipe that is happening only on forced
	if err := s.Register(WipeRule{
		Name:                    "Monthly only on forced and on time",
		Server:                  "",
		Days:                    nil,
		Hour:                    0,
		Minute:                  0,
		BlueprintWipe:           false,
		MapWipe:                 false,
		WipeOnForced:            true,
		StartTimestamp:          0,
		EndTimestamp:            0,
		MinDaysSinceLastTrigger: 0,
	}); err != nil {
		t.Fatalf("Failed to register wipe rule, error: %v", err)
	}

	//Loop over 12 months and detect 12 wipes
	triggers := make([]*WipeTrigger, 0)
	for ts := int64(0); ts < 86400*365; ts += 60 {
		triggers = append(triggers, s.Schedule(ts)...)
	}

	if len(triggers) != 12 {
		t.Errorf("Monthly wipes failed, expected 12 but got: %d", len(triggers))
	}

	//Verify that every single wipe is at the correct hour
	for _, tr := range triggers {
		dateTime := time.Unix(tr.Timestamp, 0).UTC()

		if dateTime.Hour() != ForceWipeHourUtc {
			t.Errorf("wipe hour not correct. Expected, %d, got: %d", ForceWipeHourUtc, dateTime.Hour())
		}

		if dateTime.Weekday() != time.Thursday {
			t.Errorf("wipe day not correct. Expected, Thursday, got: %d", dateTime.Weekday())
		}

		if tr.ForcedUpdate != true {
			t.Errorf("ForcedUpdate flag expected to be true, but got false")
		}
	}
}

/**
Bi-weekly schedule to wipe on monday and thursday
On the weeks when there is forced wipe on thursday, the wipe only has to happen once - when the forced is
*/
func TestBiWeeklyWipeScheduleForSameServer(t *testing.T) {
	s := NewTestScheduler()

	const wipeHour = 14

	//First test with a monthly wipe that is happening only on forced
	if err := s.Register(WipeRule{
		Name:                    "monday_wipe",
		Server:                  "test_server_1",
		Days:                    []time.Weekday{time.Monday},
		Hour:                    wipeHour,
		Minute:                  0,
		BlueprintWipe:           false,
		MapWipe:                 true,
		WipeOnForced:            false,
		StartTimestamp:          0,
		EndTimestamp:            0,
		MinDaysSinceLastTrigger: 0,
	}); err != nil {
		t.Fatalf("Failed to register wipe rule, error: %v", err)
	}

	//First test with a monthly wipe that is happening only on forced
	if err := s.Register(WipeRule{
		Name:                    "thursday_wipe",
		Server:                  "test_server_1",
		Days:                    []time.Weekday{time.Thursday},
		Hour:                    wipeHour,
		Minute:                  0,
		BlueprintWipe:           true,
		MapWipe:                 true,
		WipeOnForced:            true,
		StartTimestamp:          0,
		EndTimestamp:            0,
		MinDaysSinceLastTrigger: 0,
	}); err != nil {
		t.Fatalf("Failed to register wipe rule, error: %v", err)
	}

	triggers := make([]*WipeTrigger, 0)

	for ts := int64(0); ts < 86400*180; ts += 60 {
		triggers = append(triggers, s.Schedule(ts)...)
	}

	if len(triggers) != 52 {
		t.Errorf("Expected 52 wipes, got: %d", len(triggers))
	}

	for _, tr := range triggers {
		if hour := time.Unix(tr.Timestamp, 0).UTC().Hour(); hour != wipeHour && tr.ForcedUpdate == false {
			t.Errorf("Wipe not at %d, got: %d", wipeHour, hour)
		}

		if hour := time.Unix(tr.Timestamp, 0).UTC().Hour(); hour != ForceWipeHourUtc && tr.ForcedUpdate == true {
			t.Errorf("Forced Wipe not at %d, got: %d", ForceWipeHourUtc, hour)
		}

		if (tr.MapWipe != true || tr.BlueprintWipe != true) && time.Unix(tr.Timestamp, 0).UTC().Weekday() == time.Thursday {
			t.Errorf("Full wipe failed on a thursday, trigger: %v", tr)
		}

		if (tr.MapWipe != true || tr.BlueprintWipe != false) && time.Unix(tr.Timestamp, 0).UTC().Weekday() == time.Monday {
			t.Errorf("Map wipe failed on a monday, trigger: %v", tr)
		}
	}
}
