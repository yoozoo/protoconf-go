package protoconf_go

// Configuration protoconf configuration object interface
type Configuration interface {
	//retrieve application name
	GetApplicationName() string
	//retrieve all keys
	GetValidKeys() []string
	//set values inside the java config class
	SetValue(key string, value string) error
	//get default values from the java config class
	GetDefaultValue(key string) *string
	// add key change to the change list
	NotifyValueChange(key string, newValue string)
}

/*
// WatchKey protoconf watch key interface to add key to watch list
type WatchKey interface {
	WatchKey(key string, callback func(newValue string))
}
*/
