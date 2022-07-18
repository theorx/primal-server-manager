package Scheduler

import (
	"reflect"
	"testing"
)

func TestNewScheduler(t *testing.T) {

	want := &Scheduler{}
	got := NewScheduler()

	if got == nil {
		t.Errorf("failed to create new scheduler, got: %v, wanted %v", got, want)
	}
}

func TestRegisterRegistersRules(t *testing.T) {
	s := NewScheduler()

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

	got := s.Rules()

	if !reflect.DeepEqual(input, got) {
		t.Errorf("Output %v not equal to %v", got, input)
	}
}
