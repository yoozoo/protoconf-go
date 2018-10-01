package protoconf

// Configuration protoconf configuration object interface
type Configuration interface {
	//retrieve application name
	ApplicationName() string
	//retrieve all keys
	ValidKeys() []string
	//set values inside the java config class
	SetValue(key string, value string) error
	//get default values from the java config class
	DefaultValue(key string) *string
	// add key change to the change list
	NotifyValueChange(key string, newValue string)
}
