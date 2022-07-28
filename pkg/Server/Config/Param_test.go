package Config

import (
	"testing"
)

func TestParamGetters(t *testing.T) {
	key := "test key"
	value := "test value"
	defaultVal := "default-val"

	p := Param{
		Key:     key,
		Value:   value,
		Default: defaultVal,
		Type:    VariableType,
	}

	if p.GetKey() != key {
		t.Fatalf("GetKey() expected to retrieve %s but got %s", key, p.GetKey())
	}

	if p.GetValue() != value {
		t.Fatalf("GetValue() expected to retrieve %s but got %s", value, p.GetValue())
	}

	if p.GetType() != VariableType {
		t.Fatalf("GetType() expected to retrieve %s but got %s", VariableType, p.GetType())
	}

	pDefault := Param{
		Default: defaultVal,
		Type:    ServerFlagType,
	}

	if pDefault.GetValue() != defaultVal {
		t.Fatalf("GetValue() expected to retrieve the default: %s but got %s", defaultVal, pDefault.GetType())
	}

	if pDefault.GetType() != ServerFlagType {
		t.Fatalf("GetType() expected to retrieve %s but got %s", ServerFlagType, pDefault.GetType())
	}
}
