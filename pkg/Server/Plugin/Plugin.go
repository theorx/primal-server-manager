package Plugin

type Plugin struct {
	Name string
	//These are the keys used by the plugin, used for validation?
	ConfigKeys []string
	ConfigFile string
	//Has to be parsed before writing config file, so config keys are substituted
	ConfigTemplate string
	//Files / folders to remove on wipe
	DataFiles []DataPath
	//Determines whether all data files are removed in the event of a wipe
	WipeDataOnWipe bool
	//Source of plugin, the csharp code
	PluginSource string

	//This has to depend on the actual config, which has to be a global key-val list used across plugins
}
