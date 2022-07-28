package Config

type ParamType string

const (
	// VariableType is used for substitution within config files, templates
	VariableType ParamType = "variable"
	//ServerFlagType params are used to modify / add command line flags for server startup
	ServerFlagType ParamType = "flag"
)

type Param struct {
	Key     string
	Value   string
	Default string
	Type    ParamType
}

func (p Param) GetValue() string {
	if len(p.Value) == 0 {
		return p.Default
	}
	return p.Value
}

func (p Param) GetKey() string {
	return p.Key
}

func (p Param) GetType() ParamType {
	return p.Type
}
